import contextlib
from concurrent.futures import ThreadPoolExecutor
import gzip
import http.client
import importlib.machinery
import importlib.util
import io
import json
import os
from pathlib import Path
import queue
import shlex
import shutil
import socket
import struct
import subprocess
import tempfile
import threading
import time
from types import SimpleNamespace
import unittest
from unittest.mock import patch
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


SCRIPTS = Path(__file__).resolve().parents[1]
ROUTER = SCRIPTS / "strata-router"
LAUNCHER = SCRIPTS / "claude-strata"
loader = importlib.machinery.SourceFileLoader("strata_router", str(ROUTER))
spec = importlib.util.spec_from_loader(loader.name, loader)
router = importlib.util.module_from_spec(spec)
loader.exec_module(router)


def wait_for(predicate, timeout=5):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if predicate():
            return
        time.sleep(0.02)
    raise AssertionError("condition did not become true")


def unused_port():
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        return sock.getsockname()[1]


def respond(handler, body=b"OK", status=200, headers=()):
    handler.send_response_only(status)
    for name, value in headers:
        handler.send_header(name, value)
    handler.send_header("Content-Length", str(len(body)))
    handler.send_header("Connection", "close")
    handler.end_headers()
    if handler.command != "HEAD" and status not in {204, 304}:
        handler.wfile.write(body)
        handler.wfile.flush()
    handler.close_connection = True


class FakeHandler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def handle_request(self):
        body = self.rfile.read(int(self.headers.get("Content-Length", "0")))
        self.server.requests.put({"method": self.command, "path": self.path,
                                  "headers": list(self.headers.items()), "body": body})
        self.server.reply(self, body)

    do_GET = do_HEAD = do_POST = do_PUT = do_PATCH = do_DELETE = do_OPTIONS = handle_request
    do_TRACE = handle_request

    def log_message(self, *args):
        pass


@contextlib.contextmanager
def running_server(server):
    thread = threading.Thread(target=server.serve_forever, kwargs={"poll_interval": 0.02})
    thread.start()
    try:
        yield server
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=3)


def fake_server(stack, label):
    server = ThreadingHTTPServer(("127.0.0.1", 0), FakeHandler)
    server.daemon_threads = True
    server.requests = queue.Queue()
    server.reply = lambda handler, body: respond(handler, label.encode())
    server.url = "http://127.0.0.1:" + str(server.server_port)
    return stack.enter_context(running_server(server))


class RouterTests(unittest.TestCase):
    def setUp(self):
        self.stack = contextlib.ExitStack()
        self.addCleanup(self.stack.close)
        self.engine = fake_server(self.stack, "engine")
        self.upstream = fake_server(self.stack, "upstream")
        self.log = io.StringIO()
        server = router.RouterServer(("127.0.0.1", 0), "strata-test", self.engine.url,
                                     self.upstream.url, self.log)
        self.server = self.stack.enter_context(running_server(server))

    def request(self, method="POST", path="/v1/messages?beta=true", body=b' {"model": "opus"} ',
                headers=()):
        connection = http.client.HTTPConnection("127.0.0.1", self.server.server_port, timeout=2)
        try:
            connection.putrequest(method, path, skip_accept_encoding=True)
            connection.putheader("Content-Length", str(len(body)))
            for name, value in headers:
                connection.putheader(name, value)
            connection.endheaders(body)
            response = connection.getresponse()
            return response.status, response.getheaders(), response.read()
        finally:
            connection.close()

    def raw(self, request):
        with socket.create_connection(("127.0.0.1", self.server.server_port), timeout=2) as sock:
            sock.sendall(request)
            chunks = []
            while True:
                chunk = sock.recv(65536)
                if not chunk:
                    return b"".join(chunks)
                chunks.append(chunk)

    def records(self):
        wait_for(lambda: bool(self.log.getvalue()))
        return [json.loads(line) for line in self.log.getvalue().splitlines()]

    def test_model_and_method_routing_preserves_bytes_and_path(self):
        cases = [("POST", b' {"model": "strata-test", "extra": [1,2]}\n', self.engine),
                 ("POST", b'{"model":"opus"}', self.upstream),
                 ("POST", b"not JSON\xff", self.upstream),
                 ("POST", b"[" * 1500 + b"]" * 1500, self.upstream),
                 ("POST", b"[]", self.upstream),
                 ("GET", b'{"model":"strata-test"}', self.upstream),
                 ("TRACE", b"", self.upstream),
                 ("HEAD", b"", self.upstream)]
        for method, body, destination in cases:
            with self.subTest(method=method, body=body):
                status, _, _ = self.request(method=method, body=body)
                self.assertEqual(status, 200)
                received = destination.requests.get(timeout=1)
                self.assertEqual(received["method"], method)
                self.assertEqual(received["path"], "/v1/messages?beta=true")
                self.assertEqual(received["body"], body)
        self.assertTrue(self.engine.requests.empty())
        self.assertTrue(self.upstream.requests.empty())

    def test_auth_and_duplicate_end_to_end_headers(self):
        headers = [("Authorization", "Bearer test-secret"), ("aUtHoRiZaTiOn", "other-secret"),
                   ("X-API-Key", "api-secret"), ("x-api-KEY", "second-key"),
                   ("anthropic-beta", "oauth-test,feature-test"), ("Origin", "https://example.test"),
                   ("X-Duplicate", "first"), ("X-Duplicate", "second"),
                   ("Connection", "keep-alive, X-Remove"), ("cOnNeCtIoN", "X-Also-Remove"),
                   ("X-Remove", "hidden"), ("x-also-remove", "hidden"),
                   ("Keep-Alive", "timeout=9"), ("Proxy-Authorization", "proxy-secret"),
                   ("TE", "trailers"), ("Trailer", "X-Later"), ("Upgrade", "websocket")]
        for model, target in [("strata-test", self.engine), ("opus", self.upstream)]:
            self.request(body=json.dumps({"model": model}).encode(), headers=headers)
            received = target.requests.get(timeout=1)["headers"]
            expected = [(name, value) for name, value in headers
                        if name.lower() in {"authorization", "x-api-key", "anthropic-beta",
                                            "origin", "x-duplicate"}
                        and not (target is self.engine and name.lower() in {"authorization", "x-api-key"})]
            self.assertEqual([(name, value) for name, value in received
                              if name.lower() not in {"host", "content-length"}], expected)
            self.assertIn(("Host", "127.0.0.1:" + str(target.server_port)), received)
        for secret in ["test-secret", "other-secret", "api-secret", "second-key", "proxy-secret",
                       "oauth-test", "example.test", "hidden"]:
            self.assertNotIn(secret, self.log.getvalue())

    def test_request_framing_validation(self):
        cases = [(b"POST / HTTP/1.1\r\nHost: test\r\n\r\n", 411),
                 (b"POST / HTTP/1.1\r\nTransfer-Encoding: chunked\r\nContent-Length: 0\r\n\r\n", 411),
                 (b"POST / HTTP/1.1\r\nContent-Length: -1\r\n\r\n", 400),
                 (b"POST / HTTP/1.1\r\nContent-Length: 0\r\nContent-Length: 0\r\n\r\n", 400),
                 (b"POST / HTTP/1.1\r\nContent-Length: " + str(router.MAX_BODY_BYTES + 1).encode() + b"\r\n\r\n", 413)]
        for raw, status in cases:
            with self.subTest(raw=raw):
                self.assertTrue(self.raw(raw).startswith(("HTTP/1.1 " + str(status)).encode()))
        self.assertTrue(self.upstream.requests.empty())

    def test_body_read_deadline_and_early_eof(self):
        original = router.BODY_TIMEOUT
        router.BODY_TIMEOUT = 0.15
        try:
            data = self.raw(b"POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\nx")
            self.assertTrue(data.startswith(b"HTTP/1.1 408"))
        finally:
            router.BODY_TIMEOUT = original
        with socket.create_connection(("127.0.0.1", self.server.server_port), timeout=2) as sock:
            sock.sendall(b"POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\nx")
            sock.shutdown(socket.SHUT_WR)
            self.assertTrue(sock.recv(4096).startswith(b"HTTP/1.1 400"))

    def test_expect_continue(self):
        body = b'{"model":"opus"}'
        with socket.create_connection(("127.0.0.1", self.server.server_port), timeout=2) as sock:
            sock.sendall(b"POST /v1/messages HTTP/1.1\r\nHost: test\r\nExpect: 100-continue\r\nContent-Length: "
                         + str(len(body)).encode() + b"\r\n\r\n")
            interim = b""
            while not interim.endswith(b"\r\n\r\n"):
                interim += sock.recv(1)
            self.assertEqual(interim, b"HTTP/1.1 100 Continue\r\n\r\n")
            sock.sendall(body)
            response = http.client.HTTPResponse(sock)
            response.begin()
            self.assertEqual(response.read(), b"upstream")
        received = self.upstream.requests.get(timeout=1)
        self.assertNotIn("expect", [name.lower() for name, _ in received["headers"]])

    def test_response_framing_and_compression(self):
        payload = b"first\x00second" * 100
        for framing in ["fixed", "chunked", "close", "gzip"]:
            with self.subTest(framing=framing):
                wire = gzip.compress(payload) if framing == "gzip" else payload

                def reply(handler, body):
                    handler.send_response_only(201, "Created")
                    handler.send_header("X-Duplicate", "one")
                    handler.send_header("X-Duplicate", "two")
                    handler.send_header("Connection", "close, X-Remove")
                    handler.send_header("X-Remove", "hidden")
                    if framing == "chunked":
                        handler.send_header("Transfer-Encoding", "chunked")
                        handler.send_header("Trailer", "X-Final")
                    elif framing in {"fixed", "gzip"}:
                        handler.send_header("Content-Length", str(len(wire)))
                    if framing == "gzip":
                        handler.send_header("Content-Encoding", "gzip")
                    handler.end_headers()
                    if framing == "chunked":
                        handler.wfile.write(("%x\r\n" % len(wire)).encode() + wire + b"\r\n0\r\nX-Final: yes\r\n\r\n")
                    else:
                        handler.wfile.write(wire)
                    handler.close_connection = True

                self.upstream.reply = reply
                status, headers, received = self.request()
                self.assertEqual(status, 201)
                self.assertEqual(received, wire)
                self.assertEqual([value for name, value in headers if name == "X-Duplicate"], ["one", "two"])
                lowered = {name.lower(): value for name, value in headers}
                self.assertEqual(lowered["connection"], "close")
                self.assertNotIn("x-remove", lowered)
                self.assertNotIn("transfer-encoding", lowered)
                self.assertNotIn("trailer", lowered)
                if framing == "gzip":
                    self.assertEqual(lowered["content-encoding"], "gzip")

    def test_head_and_bodyless_statuses(self):
        for method, status in [("HEAD", 200), ("GET", 204), ("GET", 304)]:
            self.upstream.reply = lambda handler, body: respond(handler, b"metadata", status=status)
            response = self.raw((method + " / HTTP/1.1\r\nHost: test\r\n\r\n").encode())
            headers, body = response.split(b"\r\n\r\n", 1)
            self.assertEqual(body, b"")
            if status != 204:
                self.assertIn(b"Content-Length: 8", headers)
            else:
                self.assertNotIn(b"Content-Length", headers)

    def test_sse_first_event_arrives_before_last_is_sent(self):
        release = threading.Event()
        finished = threading.Event()

        def reply(handler, body):
            handler.send_response_only(200)
            handler.send_header("Content-Type", "text/event-stream")
            handler.send_header("Connection", "close")
            handler.end_headers()
            handler.wfile.write(b"data: first\n\n")
            handler.wfile.flush()
            release.wait(3)
            handler.wfile.write(b"data: last\n\n")
            handler.wfile.flush()
            finished.set()
            handler.close_connection = True

        self.upstream.reply = reply
        connection = http.client.HTTPConnection("127.0.0.1", self.server.server_port, timeout=1)
        try:
            connection.request("POST", "/", b'{}')
            response = connection.getresponse()
            self.assertEqual(response.read1(65536), b"data: first\n\n")
            self.assertFalse(finished.is_set())
            release.set()
            self.assertEqual(response.read(), b"data: last\n\n")
        finally:
            release.set()
            connection.close()

    def test_down_target_returns_anthropic_error(self):
        address = "127.0.0.1:" + str(unused_port())
        self.server.identity["upstream_url"] = "http://" + address
        status, _, body = self.request()
        self.assertEqual(status, 502)
        error = json.loads(body)
        self.assertEqual(error["type"], "error")
        self.assertEqual(error["error"]["type"], "api_error")
        self.assertEqual(error["error"]["message"], "strata-router: cannot reach " + address)

    def test_down_engine_error_identifies_url_and_remedy(self):
        url = "http://127.0.0.1:" + str(unused_port())
        self.server.identity["engine_url"] = url
        status, _, body = self.request(body=b'{"model":"strata-test"}')
        self.assertEqual(status, 502)
        message = json.loads(body)["error"]["message"]
        self.assertIn(url, message)
        self.assertIn("start it first", message)

    def assert_truncated_response(self, chunked=False, reset=False):
        release = threading.Event()

        def reply(handler, body):
            handler.send_response_only(200)
            handler.send_header("Transfer-Encoding" if chunked else "Content-Length", "chunked" if chunked else "99")
            handler.end_headers()
            handler.wfile.write(b"7\r\npartial\r\n" if chunked else b"partial")
            handler.wfile.flush()
            release.wait(3)
            if reset:
                handler.connection.setsockopt(socket.SOL_SOCKET, socket.SO_LINGER, struct.pack("ii", 1, 0))
                handler.connection.close()
            handler.close_connection = True

        self.upstream.reply = reply
        try:
            with socket.create_connection(("127.0.0.1", self.server.server_port), timeout=2) as sock:
                sock.sendall(b"POST / HTTP/1.1\r\nContent-Length: 2\r\n\r\n{}")
                response = b""
                while not response.endswith(b"partial"):
                    piece = sock.recv(65536)
                    self.assertTrue(piece, "response ended before the first bytes")
                    response += piece
                self.assertTrue(response.startswith(b"HTTP/1.1 200"))
                release.set()
                with self.assertRaises(ConnectionResetError):
                    while sock.recv(65536):
                        pass
        finally:
            release.set()
        self.assertEqual(self.records()[-1]["outcome"], "truncated")
        self.assertEqual(self.upstream.requests.qsize(), 1)

    def test_truncated_response_is_logged_without_retry(self):
        self.assert_truncated_response()

    def test_target_reset_is_logged_as_truncated(self):
        self.assert_truncated_response(reset=True)

    def test_incomplete_chunked_response_resets_client(self):
        self.assert_truncated_response(chunked=True)

    def test_redirect_is_forwarded_and_length_representation_is_preserved(self):
        def reply(handler, body):
            handler.send_response_only(302)
            handler.send_header("Location", self.engine.url + "/redirect-target")
            handler.send_header("Content-Length", "0006")
            handler.end_headers()
            handler.wfile.write(b"moved!")
            handler.close_connection = True

        self.upstream.reply = reply
        status, headers, body = self.request()
        self.assertEqual((status, body), (302, b"moved!"))
        self.assertIn(("Content-Length", "0006"), headers)
        self.assertTrue(self.engine.requests.empty())

    def test_disconnect_closes_target(self):
        release = threading.Event()
        closed = threading.Event()

        def reply(handler, body):
            handler.send_response_only(200)
            handler.send_header("Content-Length", "1000000")
            handler.end_headers()
            handler.wfile.write(b"first")
            handler.wfile.flush()
            release.wait(3)
            try:
                handler.wfile.write(b"x" * 65536)
                handler.wfile.flush()
                handler.connection.settimeout(2)
                if not handler.rfile.read(1):
                    closed.set()
            except (OSError, ConnectionError):
                closed.set()
            handler.close_connection = True

        self.upstream.reply = reply
        sock = socket.create_connection(("127.0.0.1", self.server.server_port), timeout=2)
        try:
            sock.sendall(b"POST / HTTP/1.1\r\nContent-Length: 2\r\n\r\n{}")
            self.assertTrue(sock.recv(4096).startswith(b"HTTP/1.1 200"))
            sock.setsockopt(socket.SOL_SOCKET, socket.SO_LINGER, struct.pack("ii", 1, 0))
        finally:
            sock.close()
            release.set()
        self.assertTrue(closed.wait(3), "router did not close the target after client disconnect")

    def test_identity_and_log_schema(self):
        status, _, body = self.request(method="GET", path=router.IDENTITY_PATH, body=b"")
        self.assertEqual(status, 200)
        self.assertEqual(json.loads(body), router.identity("strata-test", self.engine.url, self.upstream.url))
        self.assertEqual(json.loads(body)["upstream_url"], self.upstream.url)
        self.assertTrue(self.upstream.requests.empty())
        record = self.records()[-1]
        self.assertEqual(set(record), {"method", "path", "model", "route", "status", "outcome", "duration"})
        self.assertEqual(record["route"], "router")


class LauncherTests(unittest.TestCase):
    def setUp(self):
        self.stack = contextlib.ExitStack()
        self.addCleanup(self.stack.close)
        self.temp = Path(self.stack.enter_context(tempfile.TemporaryDirectory()))
        self.engine = fake_server(self.stack, "engine")
        self.upstream = fake_server(self.stack, "upstream")
        self.port = unused_port()
        self.log = self.temp / "router.log"
        self.processes = []
        self.stack.callback(self.stop_processes)
        scripts = self.temp / "scripts"
        scripts.mkdir()
        (scripts / "strata-router").symlink_to(ROUTER)
        (scripts / "strata-main-guard").symlink_to(SCRIPTS / "strata-main-guard")
        (scripts / "claude-strata.prompt.md").symlink_to(SCRIPTS / "claude-strata.prompt.md")
        launcher = LAUNCHER.read_text().replace("router_port=8098", "router_port=" + str(self.port))
        launcher = launcher.replace("engine_url=http://127.0.0.1:8099", "engine_url=" + self.engine.url)
        self.upstream_option = "--test-upstream-url " + shlex.quote(self.upstream.url) + " "
        launcher = launcher.replace('--port "$router_port"', self.upstream_option + '--port "$router_port"')
        (scripts / "claude-strata").write_text(launcher)
        (scripts / "claude-strata").chmod(0o755)
        self.bin = self.temp / "bin"
        self.bin.mkdir()
        (self.bin / "claude-strata").symlink_to("../scripts/claude-strata")
        self.stub = self.bin / "claude"
        self.stub.write_text('''#!/usr/bin/python3
import http.client, json, os, sys
from urllib.parse import urlsplit
if os.environ.get("PROBE_REQUEST"):
    target = urlsplit(os.environ["ANTHROPIC_BASE_URL"])
    connection = http.client.HTTPConnection(target.hostname, target.port)
    headers = {}
    if "ANTHROPIC_AUTH_TOKEN" in os.environ:
        headers["Authorization"] = "Bearer " + os.environ["ANTHROPIC_AUTH_TOKEN"]
    if "ANTHROPIC_API_KEY" in os.environ:
        headers["x-api-key"] = os.environ["ANTHROPIC_API_KEY"]
    connection.request("POST", "/v1/messages", b'{"model":"opus"}', headers)
    assert connection.getresponse().read() == b"upstream"
    connection.close()
print(json.dumps({"env": dict(os.environ), "args": sys.argv[1:]}))
''')
        self.stub.chmod(0o755)
        self.env = {"PATH": str(self.bin) + ":/usr/bin:/bin:/usr/sbin", "HOME": str(self.temp)}

    def stop_processes(self):
        for process in self.processes:
            if process.poll() is None:
                process.terminate()
            process.wait(timeout=5)

    def start_router(self, model="strata-local"):
        process = subprocess.Popen(["/usr/bin/python3", str(ROUTER), "--port", str(self.port),
                                    "--model", model, "--engine-url", self.engine.url,
                                    "--test-upstream-url", self.upstream.url,
                                    "--log-file", str(self.log)],
                                   stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        self.processes.append(process)
        wait_for(lambda: router.port_open(self.port))
        return process

    def launch(self, *args):
        return subprocess.run([str(self.bin / "claude-strata")] + list(args), env=self.env,
                              capture_output=True, text=True, timeout=10)

    def test_failed_realpath_stops_before_running_sibling_scripts(self):
        realpath = self.bin / "realpath"
        realpath.write_text("#!/bin/sh\necho 'realpath failed' >&2\nexit 1\n")
        realpath.chmod(0o755)
        for name in ["strata-router", "strata-main-guard"]:
            (self.temp / name).write_text('from pathlib import Path\n'
                                          'Path("unexpected-helper").write_text("called")\n'
                                          'print("{}")\n')
        result = subprocess.run([str(self.bin / "claude-strata")], env=self.env, cwd=self.temp,
                                capture_output=True, text=True, timeout=10)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("realpath failed", result.stderr)
        self.assertEqual(result.stdout, "")
        self.assertFalse((self.temp / "unexpected-helper").exists())

    def test_matching_router_reused_and_inherited_environment_cleared(self):
        process = self.start_router()
        cleared = ["ANTHROPIC_AUTH_TOKEN", "ANTHROPIC_API_KEY", "ANTHROPIC_SMALL_FAST_MODEL",
                   "ANTHROPIC_DEFAULT_HAIKU_MODEL", "ANTHROPIC_DEFAULT_SONNET_MODEL",
                   "ANTHROPIC_DEFAULT_OPUS_MODEL"]
        self.env.update({name: "inherited-secret" for name in cleared})
        self.env["PROBE_REQUEST"] = "1"
        for _ in range(2):
            result = self.launch("-p", "probe")
            self.assertEqual(result.returncode, 0, result.stderr)
            output = json.loads(result.stdout)
            for name in cleared:
                self.assertNotIn(name, output["env"])
            expected = {"ANTHROPIC_BASE_URL": "http://127.0.0.1:" + str(self.port),
                        "ANTHROPIC_MODEL": "strata-local", "CLAUDE_CODE_SUBAGENT_MODEL": "opus",
                        "CLAUDE_CODE_SUBAGENT_MODEL_FORCE": "1", "CLAUDE_CODE_FORK_SUBAGENT": "0",
                        "CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION": "0",
                        "CLAUDE_CODE_MAX_CONTEXT_TOKENS": "200000",
                        "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"}
            for name, value in expected.items():
                self.assertEqual(output["env"][name], value)
            args = output["args"]
            self.assertEqual(args[:2], ["--dangerously-skip-permissions", "--settings"])
            hooks = json.loads(args[2])["hooks"]
            self.assertEqual(hooks["PreToolUse"][0]["matcher"], "Edit|Write|NotebookEdit")
            for name in ["PreToolUse", "UserPromptSubmit"]:
                self.assertEqual(shlex.split(hooks[name][0]["hooks"][0]["command"]),
                                 ["/usr/bin/python3", str(SCRIPTS / "strata-main-guard")])
            self.assertEqual(args[3:], ["--append-system-prompt-file",
                                       str(self.temp.resolve() / "scripts/claude-strata.prompt.md"), "-p", "probe"])
            received = self.upstream.requests.get(timeout=1)
            self.assertFalse(any(name.lower() in {"authorization", "x-api-key"}
                                 for name, _ in received["headers"]))
        self.assertIsNone(process.poll())

    def test_stale_router_refused_with_pid(self):
        process = self.start_router(model="old-model")
        result = self.launch()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("PID " + str(process.pid), result.stderr)
        self.assertIn("kill <PID>", result.stderr)
        self.assertIsNone(process.poll())

    def test_launcher_without_override_refuses_test_upstream_router(self):
        process = self.start_router()
        script = self.temp / "scripts/claude-strata"
        script.write_text(script.read_text().replace(self.upstream_option, ""))
        result = self.launch()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("PID " + str(process.pid), result.stderr)
        self.assertIsNone(process.poll())
        self.assertTrue(self.upstream.requests.empty())

    def test_foreign_listener_refused_with_pid(self):
        self.port = self.upstream.server_port
        script = self.temp / "scripts/claude-strata"
        content = LAUNCHER.read_text().replace("router_port=8098", "router_port=" + str(self.port))
        script.write_text(content)
        result = self.launch()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("PID " + str(os.getpid()), result.stderr)
        self.assertIn("kill <PID>", result.stderr)

    def test_source_digest_mismatch_refused(self):
        self.upstream.reply = lambda handler, body: respond(handler, json.dumps({
            "source_sha256": "stale", "model": "strata-local", "engine_url": self.engine.url,
            "upstream_url": self.upstream.url,
        }).encode())
        self.assertFalse(router.matches_listener(self.upstream.server_port,
                                                router.identity("strata-local", self.engine.url, self.upstream.url)))

    def test_startup_detaches_and_concurrent_launches_reuse_winner(self):
        args = SimpleNamespace(port=self.port, model="strata-local", engine_url=self.engine.url,
                               log_file=str(self.log), test_upstream_url=self.upstream.url)
        real_popen = subprocess.Popen

        def track_child(*args, **kwargs):
            child = real_popen(*args, **kwargs)
            self.processes.append(child)
            return child

        try:
            with patch.object(router.subprocess, "Popen", side_effect=track_child):
                with ThreadPoolExecutor(max_workers=2) as pool:
                    results = list(pool.map(router.ensure_router, [args, args]))
            self.assertEqual(results, [0, 0])
            self.assertTrue(router.matches_listener(self.port, router.identity("strata-local", self.engine.url,
                                                                               self.upstream.url)))
            winners = [process for process in self.processes if process.poll() is None]
            wait_for(lambda: sum(process.poll() is None for process in self.processes) == 1)
            pid = next(process.pid for process in winners if process.poll() is None)
            self.assertEqual(os.getsid(pid), pid)
            self.assertNotEqual(os.getsid(pid), os.getsid(0))
        finally:
            self.stop_processes()
            wait_for(lambda: not router.port_open(self.port))

    def test_bound_port_startup_fails_and_test_override_is_loopback_only(self):
        self.start_router()
        command = ["/usr/bin/python3", str(ROUTER), "--port", str(self.port),
                   "--model", "strata-local", "--engine-url", self.engine.url,
                   "--log-file", str(self.temp / "loser.log")]
        result = subprocess.run(command, capture_output=True, timeout=5)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("bind listener", (self.temp / "loser.log").read_text())
        result = subprocess.run(command + ["--test-upstream-url", "https://example.test"],
                                capture_output=True, text=True, timeout=5)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("loopback HTTP origin", result.stderr)

    def test_failed_child_exits_without_waiting_for_startup_deadline(self):
        blocker = self.temp / "not-a-directory"
        blocker.write_text("blocked")
        command = ["/usr/bin/python3", str(ROUTER), "--ensure", "--port", str(self.port),
                   "--model", "strata-local", "--engine-url", self.engine.url,
                   "--test-upstream-url", self.upstream.url, "--log-file", str(blocker / "router.log")]
        result = subprocess.run(command, capture_output=True, text=True, timeout=3)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("startup failed", result.stderr)
        self.assertFalse(router.port_open(self.port))

    @unittest.skipUnless(os.environ.get("STRATA_ROUTER_E2E") == "1", "set STRATA_ROUTER_E2E=1 for Claude Code probe")
    def test_claude_subagent_routes_upstream_and_fork_is_refused(self):
        claude = shutil.which("claude")
        self.assertIsNotNone(claude)
        self.start_router()
        launch = self.launch()
        self.assertEqual(launch.returncode, 0, launch.stderr)
        env = json.loads(launch.stdout)["env"]
        env["ANTHROPIC_AUTH_TOKEN"] = "x"
        env["CLAUDE_CONFIG_DIR"] = str(self.temp / "claude-config")
        marker = "STRATA_ROUTER_FAKE_PROBE"

        def sse(event, data):
            return ("event: " + event + "\ndata: " + json.dumps(data) + "\n\n").encode()

        for agent_type in ["general-purpose", "fork"]:
            with self.subTest(agent_type=agent_type):
                requests = []
                spawned = []

                def reply(handler, body):
                    if handler.command == "HEAD":
                        respond(handler, b"")
                        return
                    req = json.loads(body)
                    if "count_tokens" in handler.path:
                        respond(handler, b'{"input_tokens":10}')
                        return
                    requests.append(req)
                    system = req.get("system", "")
                    system = system if isinstance(system, str) else json.dumps(system)
                    tools = [tool.get("name") for tool in req.get("tools", [])]
                    spawn = not spawned and marker in system and "Agent" in tools
                    message = {"id": "msg_probe", "type": "message", "role": "assistant",
                               "model": req.get("model"), "content": [], "stop_reason": None,
                               "stop_sequence": None, "usage": {"input_tokens": 10, "output_tokens": 1}}
                    events = [sse("message_start", {"type": "message_start", "message": message})]
                    if spawn:
                        spawned.append(True)
                        args = {"description": "probe", "prompt": "Reply with DONE.",
                                "subagent_type": agent_type}
                        block = {"type": "tool_use", "id": "toolu_probe1", "name": "Agent", "input": {}}
                        delta = {"type": "input_json_delta", "partial_json": json.dumps(args)}
                    else:
                        block = {"type": "text", "text": ""}
                        delta = {"type": "text_delta", "text": "DONE"}
                    events += [
                        sse("content_block_start", {"type": "content_block_start", "index": 0, "content_block": block}),
                        sse("content_block_delta", {"type": "content_block_delta", "index": 0, "delta": delta}),
                        sse("content_block_stop", {"type": "content_block_stop", "index": 0}),
                        sse("message_delta", {"type": "message_delta", "delta": {
                            "stop_reason": "tool_use" if spawn else "end_turn", "stop_sequence": None,
                        }, "usage": {"output_tokens": 1}}),
                        sse("message_stop", {"type": "message_stop"}),
                    ]
                    respond(handler, b"".join(events), headers=[("Content-Type", "text/event-stream")])

                self.engine.reply = self.upstream.reply = reply
                result = subprocess.run([claude, "--dangerously-skip-permissions", "--setting-sources", "",
                                         "--strict-mcp-config", "--mcp-config", '{"mcpServers":{}}',
                                         "--no-session-persistence", "--append-system-prompt", marker,
                                         "--verbose", "--output-format", "stream-json", "-p", "Run the probe."],
                                        env=env, cwd=str(self.temp), capture_output=True, text=True, timeout=60)
                self.assertEqual(result.returncode, 0, result.stderr + result.stdout[-4000:])
                self.assertTrue(spawned, result.stdout[-4000:])
                models = [req["model"] for req in requests]
                self.assertIn("strata-local", models)
                if agent_type == "general-purpose":
                    self.assertIn("claude-opus-5-5", models)
                    records = [json.loads(line) for line in self.log.read_text().splitlines()
                               if line.startswith("{")]
                    for model, route in [("strata-local", "engine"), ("claude-opus-5-5", "upstream")]:
                        self.assertTrue(any(row["model"] == model and row["route"] == route and row["status"] == 200
                                            for row in records), records)
                    final = json.loads(result.stdout.splitlines()[-1])
                    self.assertIn("claude-opus-5-5", final["modelUsage"])
                else:
                    self.assertTrue(all(model == "strata-local" for model in models), models)
                    self.assertIn("Agent type 'fork' not found", json.dumps(requests))


if __name__ == "__main__":
    unittest.main()
