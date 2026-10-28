import contextlib
import importlib.machinery
import importlib.util
import json
import os
from pathlib import Path
import shlex
import shutil
import subprocess
import tempfile
import unittest
from unittest.mock import patch

import test_strata_router as helpers


GUARD = helpers.SCRIPTS / "strata-main-guard"
loader = importlib.machinery.SourceFileLoader("strata_main_guard", str(GUARD))
spec = importlib.util.spec_from_loader(loader.name, loader)
guard = importlib.util.module_from_spec(spec)
loader.exec_module(guard)


class GuardTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name).resolve()

    def event(self, tool="Edit", path="src/file.py", **fields):
        key = "notebook_path" if tool == "NotebookEdit" else "file_path"
        return dict({"hook_event_name": "PreToolUse", "tool_name": tool,
                     "cwd": str(self.root), "tool_input": {key: path}}, **fields)

    def assert_denied(self, event):
        output = guard.hook_output(event)["hookSpecificOutput"]
        self.assertEqual(output["hookEventName"], "PreToolUse")
        self.assertEqual(output["permissionDecision"], "deny")
        for text in ["Claude subagent", "goal", "files", "constraints", "check"]:
            self.assertIn(text, output["permissionDecisionReason"])
        self.assertIn("Write scratch and tool-input files in the scratchpad directory.",
                      output["permissionDecisionReason"])

    def test_main_thread_edit_tools_denied_and_subagents_silent(self):
        for tool in ["Edit", "Write", "NotebookEdit"]:
            for path in ["src/file.py", str(self.root / "src/file.py")]:
                with self.subTest(tool=tool, path=path):
                    self.assert_denied(self.event(tool, path))
                    self.assertIsNone(guard.hook_output(self.event(tool, path, agent_id="agent-test",
                                                                  agent_type="general-purpose")))
            for agent_id in [None, "", 0, {}]:
                self.assert_denied(self.event(tool, agent_id=agent_id))

    def test_plan_and_current_uid_scratchpad_exceptions(self):
        scratch = Path("/private/tmp/claude-" + str(os.getuid()))
        for tool in ["Edit", "Write", "NotebookEdit"]:
            for path in ["docs/plans/task.md", "nested/docs/plans/subdir/task.md",
                         str(self.root / "docs/plans/task.md"), str(scratch / "session/scratchpad/probe.txt")]:
                with self.subTest(tool=tool, path=path):
                    self.assertIsNone(guard.hook_output(self.event(tool, path)))

    def test_path_boundaries_and_parent_traversal(self):
        scratch = "/private/tmp/claude-" + str(os.getuid())
        for path in [scratch + "-other/file.txt", "/private/tmp/claude-" + str(os.getuid() + 1) + "/file.txt",
                     scratch + "/../repo/file.txt", "docs/plans-old/task.md", "otherdocs/plans/task.md",
                     "docs/plans/../../src/file.py", "docs/plans", "docs/plans/task.md/.."]:
            with self.subTest(path=path):
                self.assert_denied(self.event(path=path))

    def test_native_plan_directory_default_empty_and_custom_config(self):
        with patch.dict(os.environ, {"HOME": str(self.root)}):
            for config in [None, "", str(self.root / "custom-config")]:
                if config is None:
                    os.environ.pop("CLAUDE_CONFIG_DIR", None)
                else:
                    os.environ["CLAUDE_CONFIG_DIR"] = config
                plans = (Path(config) if config else self.root / ".claude") / "plans"
                plans.mkdir(parents=True, exist_ok=True)
                if not (plans / "escape").exists():
                    (plans / "escape").symlink_to(self.root, target_is_directory=True)
                for tool in ["Edit", "Write", "NotebookEdit"]:
                    self.assertIsNone(guard.hook_output(self.event(tool, str(plans / "task.md"))))
                    for path in [plans, plans.parent / "plans-old/task.md", plans / "../source.txt",
                                 plans / "escape/source.txt"]:
                        self.assert_denied(self.event(tool, str(path)))

    def test_symlinks_follow_the_actual_target(self):
        plans = self.root / "docs/plans"
        plans.mkdir(parents=True)
        outside = self.root / "src"
        outside.mkdir()
        (plans / "escape").symlink_to(outside, target_is_directory=True)
        (self.root / "plan-link").symlink_to(plans, target_is_directory=True)
        (self.root / "loop").symlink_to("loop")
        self.assert_denied(self.event(path="docs/plans/escape/file.py"))
        self.assert_denied(self.event(path="loop/file.py"))
        self.assertIsNone(guard.hook_output(self.event(path="plan-link/task.md")))

    def test_missing_or_invalid_targets_deny(self):
        for value in [None, "", 12, [], "invalid\x00path"]:
            self.assert_denied(self.event(path=value))
        event = self.event()
        for value in [None, [], {}, {"notebook_path": "docs/plans/task.md"}]:
            event["tool_input"] = value
            self.assert_denied(event)
        self.assert_denied(self.event(cwd=None))

    def test_other_tools_and_events_are_silent(self):
        for tool in ["Read", "Bash", "Agent", None, []]:
            self.assertIsNone(guard.hook_output(self.event(tool)))
        self.assertIsNone(guard.hook_output({"hook_event_name": "Stop"}))

    def test_prompt_context_is_the_entire_sibling_prompt(self):
        prompt = (GUARD.parent / "claude-strata.prompt.md").read_text()
        output = guard.hook_output({"hook_event_name": "UserPromptSubmit"})
        self.assertEqual(output, {"hookSpecificOutput": {
            "hookEventName": "UserPromptSubmit", "additionalContext": prompt,
        }})

    def test_prompt_comes_from_the_guards_resolved_directory(self):
        directory = self.root / "scripts"
        directory.mkdir()
        copy = directory / GUARD.name
        shutil.copyfile(GUARD, copy)
        prompt = "# Local-engine main session\n\nA sibling prompt with Unicode: \u2713\n"
        (directory / "claude-strata.prompt.md").write_text(prompt)
        link = self.root / "guard-link"
        link.symlink_to(copy)
        result = subprocess.run(["/usr/bin/python3", str(link)], cwd=str(self.root),
                                input='{"hook_event_name":"UserPromptSubmit"}',
                                capture_output=True, text=True, check=True)
        self.assertEqual(json.loads(result.stdout)["hookSpecificOutput"]["additionalContext"], prompt)

    def test_settings_command_handles_spaces_and_quotes_and_cli_is_silent_on_allow(self):
        copy = self.root / "guard with a 'quote'"
        shutil.copyfile(GUARD, copy)
        shutil.copyfile(GUARD.parent / "claude-strata.prompt.md", self.root / "claude-strata.prompt.md")
        result = subprocess.run(["/usr/bin/python3", str(copy), "--settings"],
                                capture_output=True, text=True, check=True)
        hooks = json.loads(result.stdout)["hooks"]
        self.assertEqual(hooks["PreToolUse"][0]["matcher"], "Edit|Write|NotebookEdit")
        for name in ["PreToolUse", "UserPromptSubmit"]:
            command = shlex.split(hooks[name][0]["hooks"][0]["command"])
            self.assertEqual(command, ["/usr/bin/python3", str(copy)])
        command = shlex.split(hooks["PreToolUse"][0]["hooks"][0]["command"])
        for event, silent in [(self.event(), False), (self.event(agent_id="agent-test"), True),
                              ({"hook_event_name": "UserPromptSubmit"}, False)]:
            result = subprocess.run(command, input=json.dumps(event), capture_output=True, text=True, check=True)
            self.assertEqual(result.stderr, "")
            if silent:
                self.assertEqual(result.stdout, "")
            else:
                self.assertEqual(json.loads(result.stdout), guard.hook_output(event))

    def test_invalid_hook_json_fails_without_echoing_input(self):
        for value in ["secret-invalid", "[]", "null"]:
            result = subprocess.run(["/usr/bin/python3", str(GUARD)], input=value,
                                    capture_output=True, text=True)
            self.assertEqual(result.returncode, 2)
            self.assertEqual(result.stdout, "")
            self.assertNotIn(value, result.stderr)
            self.assertNotIn("Traceback", result.stderr)


@unittest.skipUnless(os.environ.get("STRATA_GUARD_E2E") == "1", "set STRATA_GUARD_E2E=1 for Claude Code probe")
class GuardIntegrationTests(unittest.TestCase):
    def test_bypass_denies_main_edit_allows_subagent_edit_and_merges_user_hooks(self):
        claude = shutil.which("claude")
        self.assertIsNotNone(claude)
        with contextlib.ExitStack() as stack:
            root = Path(stack.enter_context(tempfile.TemporaryDirectory())).resolve()
            server = helpers.fake_server(stack, "probe")
            main_file, sub_file = root / "main.txt", root / "subagent.txt"
            main_file.write_text("original\n")
            sub_file.write_text("original\n")
            config = root / "config"
            config.mkdir()
            hook_log = root / "user-hooks.jsonl"
            user_hook = root / "user-hook.py"
            user_hook.write_text('''import json, sys
event = json.load(sys.stdin)
with open(sys.argv[1], "a") as log:
    log.write(json.dumps(event) + "\\n")
if event["hook_event_name"] == "UserPromptSubmit":
    print(json.dumps({"hookSpecificOutput": {"hookEventName": "UserPromptSubmit",
                                          "additionalContext": "USER_SETTINGS_CONTEXT_MARKER"}}))
''')
            command = shlex.join(["/usr/bin/python3", str(user_hook), str(hook_log)])
            (config / "settings.json").write_text(json.dumps({"hooks": {
                "PreToolUse": [{"matcher": "Edit", "hooks": [{"type": "command", "command": command}]}],
                "UserPromptSubmit": [{"hooks": [{"type": "command", "command": command}]}],
            }}))
            requests = []
            counts = {"strata-local": 0, "claude-opus-5-5": 0}

            def sse(event, data):
                return ("event: " + event + "\ndata: " + json.dumps(data) + "\n\n").encode()

            def reply(handler, body):
                if handler.command != "POST":
                    helpers.respond(handler, b"")
                    return
                req = json.loads(body)
                if "count_tokens" in handler.path:
                    helpers.respond(handler, b'{"input_tokens":10}')
                    return
                requests.append(req)
                model = req["model"]
                index = counts[model]
                counts[model] += 1
                target = main_file if model == "strata-local" else sub_file
                name, args = None, None
                if index == 0:
                    name, args = "Read", {"file_path": str(target)}
                elif index == 1:
                    name, args = "Edit", {"file_path": str(target), "old_string": "original",
                                          "new_string": "main-edited" if model == "strata-local" else "subagent-edited"}
                elif index == 2 and model == "strata-local":
                    name, args = "Agent", {"description": "guard probe", "prompt": "Edit subagent.txt as requested.",
                                           "subagent_type": "general-purpose"}
                message = {"id": "msg_probe", "type": "message", "role": "assistant", "model": model,
                           "content": [], "stop_reason": None, "stop_sequence": None,
                           "usage": {"input_tokens": 10, "output_tokens": 1}}
                block = ({"type": "tool_use", "id": "toolu_" + model + str(index), "name": name, "input": {}}
                         if name else {"type": "text", "text": ""})
                delta = ({"type": "input_json_delta", "partial_json": json.dumps(args)}
                         if name else {"type": "text_delta", "text": "DONE"})
                events = [
                    sse("message_start", {"type": "message_start", "message": message}),
                    sse("content_block_start", {"type": "content_block_start", "index": 0, "content_block": block}),
                    sse("content_block_delta", {"type": "content_block_delta", "index": 0, "delta": delta}),
                    sse("content_block_stop", {"type": "content_block_stop", "index": 0}),
                    sse("message_delta", {"type": "message_delta", "delta": {
                        "stop_reason": "tool_use" if name else "end_turn", "stop_sequence": None,
                    }, "usage": {"output_tokens": 1}}),
                    sse("message_stop", {"type": "message_stop"}),
                ]
                helpers.respond(handler, b"".join(events), headers=[("Content-Type", "text/event-stream")])

            server.reply = reply
            env = {"HOME": str(root), "PATH": os.environ.get("PATH", "/usr/bin:/bin"),
                   "CLAUDE_CONFIG_DIR": str(config), "ANTHROPIC_BASE_URL": server.url,
                   "ANTHROPIC_AUTH_TOKEN": "x", "ANTHROPIC_MODEL": "strata-local",
                   "CLAUDE_CODE_SUBAGENT_MODEL": "opus", "CLAUDE_CODE_SUBAGENT_MODEL_FORCE": "1",
                   "CLAUDE_CODE_FORK_SUBAGENT": "0", "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
                   "CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION": "0", "CLAUDE_CODE_MAX_CONTEXT_TOKENS": "200000"}
            result = subprocess.run([claude, "--dangerously-skip-permissions", "--setting-sources", "user",
                                     "--settings", json.dumps(guard.settings()), "--append-system-prompt-file",
                                     str(helpers.SCRIPTS / "claude-strata.prompt.md"), "--strict-mcp-config",
                                     "--mcp-config", '{"mcpServers":{}}', "--no-session-persistence",
                                     "--verbose", "--output-format", "stream-json", "-p", "Run the edit guard probe."],
                                    env=env, cwd=str(root), capture_output=True, text=True, timeout=60)
            self.assertEqual(result.returncode, 0, result.stderr + result.stdout[-6000:])
            main_requests = [req for req in requests if req["model"] == "strata-local"]
            sub_requests = [req for req in requests if req["model"] == "claude-opus-5-5"]
            self.assertGreaterEqual(len(main_requests), 4, result.stdout[-6000:])
            self.assertGreaterEqual(len(sub_requests), 3, result.stdout[-6000:])
            self.assertEqual(main_file.read_text(), "original\n")
            self.assertEqual(sub_file.read_text(), "subagent-edited\n")
            self.assertIn(guard.DENY_REASON, json.dumps(main_requests[2]["messages"]))
            prompt = (GUARD.parent / "claude-strata.prompt.md").read_text().rstrip()
            self.assertIn(json.dumps(prompt)[1:-1], json.dumps(main_requests[0]["messages"]))
            self.assertIn(json.dumps(prompt)[1:-1], json.dumps(main_requests[0]["system"]))
            self.assertIn("USER_SETTINGS_CONTEXT_MARKER", json.dumps(main_requests[0]["messages"]))
            hooks = [json.loads(line) for line in hook_log.read_text().splitlines()]
            edits = [event for event in hooks if event.get("tool_name") == "Edit"]
            self.assertTrue(any(not event.get("agent_id") for event in edits), edits)
            self.assertTrue(any(event.get("agent_id") and event.get("agent_type") == "general-purpose"
                                for event in edits), edits)


if __name__ == "__main__":
    unittest.main()
