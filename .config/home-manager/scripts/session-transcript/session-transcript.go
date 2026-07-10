// session-transcript — dump the current Claude Code session as a redacted, compacted
// transcript for codex-consult's --context mode (which runs this; standalone: `… | less`).
// It reads the session's on-disk JSONL — found via $CLAUDE_CODE_SESSION_ID — and prints the
// visible conversation (task, messages, tool calls + results). Thinking is NOT recoverable:
// Claude Code persists it as an opaque signature, not plaintext.
//
// Redaction is load-bearing because --context egresses to OpenAI: it strips <system-reminder>
// blocks (your private CLAUDE.md, memory, MCP/org data), drops images, skips sub-agent
// sidechains, truncates long blocks, and caps size (keeping first + most recent turns) — a
// lossy approximation, not the byte-exact context the server-side advisor sees.
//
// Built and installed on PATH by home-manager, like codex-consult. Malformed/partial JSON
// lines (the last may be mid-write) are skipped; stdlib only.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Budgets. Chosen to give Codex enough context to advise without egressing or paying for
// the whole session: blockMax tames a single giant tool output (head+tail kept), total
// caps the transcript at a few tens of thousands of tokens. Both overridable via flags.
const (
	defaultTotalBudget = 120000 // max output chars (~30k tokens)
	defaultBlockMax    = 3000   // max chars per block before head/tail truncation
)

// systemReminder matches a harness-injected <system-reminder>…</system-reminder> block.
// These are the largest privacy surface in the transcript (they embed CLAUDE.md, the
// memory index, and MCP instructions), so they are stripped before anything leaves the box.
var systemReminder = regexp.MustCompile(`(?s)<system-reminder>.*?</system-reminder>`)

func main() { os.Exit(run()) }

type options struct {
	file     string
	total    int
	blockMax int
	thinking bool
}

func run() int {
	opt := options{total: defaultTotalBudget, blockMax: defaultBlockMax, thinking: true}
	argv := os.Args[1:]
	needVal := func(i int, flag string) (string, bool) {
		if i+1 >= len(argv) {
			fmt.Fprintf(os.Stderr, "session-transcript: %s needs a value\n", flag)
			return "", false
		}
		return argv[i+1], true
	}
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		switch {
		case a == "-f" || a == "--file":
			v, ok := needVal(i, "-f/--file")
			if !ok {
				return 2
			}
			opt.file, i = v, i+1
		case a == "-b" || a == "--budget":
			v, ok := needVal(i, "-b/--budget")
			if !ok {
				return 2
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "session-transcript: -b/--budget: %v\n", err)
				return 2
			}
			opt.total, i = n, i+1
		case a == "--block-max":
			v, ok := needVal(i, "--block-max")
			if !ok {
				return 2
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "session-transcript: --block-max: %v\n", err)
				return 2
			}
			opt.blockMax, i = n, i+1
		case a == "--no-thinking":
			opt.thinking = false
		case a == "-h" || a == "--help":
			usage(os.Stdout)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "session-transcript: unknown option: %s\n", a)
			usage(os.Stderr)
			return 2
		}
	}

	path, err := locate(opt.file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "session-transcript: %v\n", err)
		return 1
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "session-transcript: cannot open %s: %v\n", path, err)
		return 1
	}
	defer f.Close()

	turns, err := parse(f, opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "session-transcript: reading %s: %v\n", path, err)
		return 1
	}
	if len(turns) == 0 {
		fmt.Fprintf(os.Stderr, "session-transcript: no conversation turns found in %s\n", path)
		return 1
	}
	out := assemble(turns, opt.total)
	if !strings.HasSuffix(out, "\n") {
		out += "\n" // terminate the stream so a concatenated brief starts on its own line
	}
	fmt.Fprintf(os.Stderr, "session-transcript: %s → %d turns, %d chars\n", path, len(turns), len(out))
	io.WriteString(os.Stdout, out)
	return 0
}

// locate resolves the transcript path: an explicit -f wins; otherwise the current
// session's log is found by its id ($CLAUDE_CODE_SESSION_ID) under ~/.claude/projects.
// Globbing by the id filename avoids reproducing Claude Code's cwd→slug encoding.
func locate(fileFlag string) (string, error) {
	if fileFlag != "" {
		return fileFlag, nil
	}
	id := os.Getenv("CLAUDE_CODE_SESSION_ID")
	if id == "" {
		return "", fmt.Errorf("CLAUDE_CODE_SESSION_ID not set — run inside a Claude Code session, or pass -f <transcript.jsonl>")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve home dir: %v", err)
	}
	matches, _ := filepath.Glob(filepath.Join(home, ".claude", "projects", "*", id+".jsonl"))
	if len(matches) == 0 {
		return "", fmt.Errorf("no transcript found for session %s under %s", id, filepath.Join(home, ".claude", "projects"))
	}
	// A session id is globally unique, so >1 match is not expected; if it happens, the
	// freshest file is the live one.
	if len(matches) > 1 {
		sort.Slice(matches, func(i, j int) bool { return mtime(matches[i]).After(mtime(matches[j])) })
	}
	return matches[0], nil
}

func mtime(p string) time.Time {
	fi, err := os.Stat(p)
	if err != nil {
		return time.Time{}
	}
	return fi.ModTime()
}

// JSONL shapes. Only the fields this tool reads are declared; everything else is ignored.
type line struct {
	Type        string          `json:"type"`
	IsSidechain bool            `json:"isSidechain"`
	Message     json.RawMessage `json:"message"`
}

type message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"` // a string OR an array of blocks
}

type block struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`     // text
	Thinking string          `json:"thinking"` // thinking
	Name     string          `json:"name"`     // tool_use
	Input    json.RawMessage `json:"input"`    // tool_use
	Content  json.RawMessage `json:"content"`  // tool_result: a string OR an array of blocks
}

// parse reads the JSONL in file order (which is chronological) and returns one formatted
// string per conversation turn. Non-conversation lines, sub-agent sidechains, and lines
// that fail to parse (e.g. a partially written trailing line) are skipped.
func parse(r io.Reader, opt options) ([]string, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20) // lines embed tool output and can be large
	var turns []string
	n := 0
	dropped := 0             // non-empty records skipped because their JSON would not decode
	lastLinePartial := false // a top-level decode failure on the final line = expected live-append
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) == "" {
			continue // blank line — not a record, and not a decode failure
		}
		var ln line
		if json.Unmarshal(sc.Bytes(), &ln) != nil {
			dropped++
			lastLinePartial = true
			continue
		}
		lastLinePartial = false
		if ln.IsSidechain || (ln.Type != "user" && ln.Type != "assistant") {
			continue
		}
		role, body, ok := formatTurn(ln, opt)
		if !ok {
			// The line decoded, but its nested message/content did not — a malformed or
			// schema-changed record. Unlike a truncated final line, this is never an expected
			// partial write, so it always counts toward the incompleteness warning.
			dropped++
			continue
		}
		if strings.TrimSpace(body) == "" {
			continue
		}
		n++
		turns = append(turns, fmt.Sprintf("## [%d] %s\n%s", n, role, body))
	}
	if err := sc.Err(); err != nil {
		// A scan error (I/O, or a record over the buffer cap) truncates the transcript; report
		// it rather than hand --context a silent partial.
		return turns, err
	}
	// Tolerate one expected artifact: a top-level decode failure on the very last line is the
	// transcript's partially written live-append. Any other dropped record means corruption or
	// a JSONL schema change silently lost turns — warn so --context never passes a truncated
	// transcript off as complete.
	if lastLinePartial {
		dropped--
	}
	if dropped > 0 {
		fmt.Fprintf(os.Stderr, "session-transcript: warning — %d record(s) failed to decode and were skipped; transcript may be incomplete (corruption or schema change)\n", dropped)
	}
	return turns, nil
}

func formatTurn(ln line, opt options) (role, body string, ok bool) {
	var m message
	if json.Unmarshal(ln.Message, &m) != nil {
		return "", "", false // malformed message — signal the drop so parse can flag it
	}
	role = m.Role
	if role == "" {
		role = ln.Type
	}
	body, ok = formatContent(m.Content, opt)
	return role, strings.TrimRight(body, "\n"), ok
}

// formatContent renders a message's content, which the JSONL stores either as a plain string
// or as an array of typed blocks. ok is false only when content is present but fits neither
// shape (a malformed or schema-changed record) so the caller can flag the drop; absent
// content is a valid empty turn.
func formatContent(raw json.RawMessage, opt options) (string, bool) {
	if len(raw) == 0 {
		return "", true
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return truncate(redactText(s), opt.blockMax), true
	}
	var blocks []block
	if json.Unmarshal(raw, &blocks) != nil {
		return "", false
	}
	var parts []string
	for _, bl := range blocks {
		if p := formatBlock(bl, opt); p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, "\n\n"), true
}

func formatBlock(bl block, opt options) string {
	switch bl.Type {
	case "text":
		return truncate(redactText(bl.Text), opt.blockMax)
	case "thinking":
		if !opt.thinking {
			return ""
		}
		// Thinking is persisted as an opaque signature with an empty text field, so this
		// is normally "" (skipped); it renders only in the rare case plaintext is present.
		t := truncate(redactText(bl.Thinking), opt.blockMax)
		if t == "" {
			return ""
		}
		return "<thinking>\n" + t + "\n</thinking>"
	case "image":
		return "[image omitted]"
	}
	// Tool calls and their results arrive under several block types (tool_use /
	// server_tool_use; tool_result / advisor_tool_result / web_search_tool_result / …).
	// Rather than enumerate every one, key off the fields present: a name or input reads
	// as a call, a content payload reads as a result. Unknown or empty blocks fall to "".
	if bl.Name != "" || len(bl.Input) > 0 {
		label := bl.Name
		if label == "" {
			label = bl.Type
		}
		// redactText the input too: a tool's arguments (a command, a file body, a prompt)
		// can embed a <system-reminder> or other stripped content, so the guarantee must
		// hold on this path as it does for text and results — not only on message text.
		return fmt.Sprintf("→ tool_use: %s\n%s", label, truncate(redactText(compactJSON(bl.Input)), opt.blockMax))
	}
	if len(bl.Content) > 0 {
		out := truncate(strings.TrimSpace(resultContent(bl.Content)), opt.blockMax)
		if out == "" {
			out = "(empty)"
		}
		return "← tool_result:\n" + out
	}
	return ""
}

// resultContent renders a tool_result's content (a string, or an array whose text blocks
// are kept and image blocks are dropped).
//
// Unlike formatContent, this does not signal decode failure to the incompleteness warning:
// a result that fits neither form falls back to compactJSON below rather than being dropped,
// and the warning is scoped to whole-turn losses, not sub-blocks within a turn that is present.
func resultContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return redactText(s)
	}
	var blocks []block
	if json.Unmarshal(raw, &blocks) == nil {
		var parts []string
		for _, bl := range blocks {
			switch bl.Type {
			case "text":
				parts = append(parts, redactText(bl.Text))
			case "image":
				parts = append(parts, "[image omitted]")
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "\n")
		}
	}
	// Some results (e.g. advisor_tool_result) carry a structured object rather than a
	// string or block array; render it compactly instead of dropping it.
	return redactText(compactJSON(raw))
}

func redactText(s string) string {
	return strings.TrimSpace(systemReminder.ReplaceAllString(s, ""))
}

func compactJSON(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var buf bytes.Buffer
	if json.Compact(&buf, raw) != nil {
		return string(raw)
	}
	return buf.String()
}

// truncate keeps the head and tail of an over-long block with a marker between, so both
// what a step started with and how it ended survive. Rune-based, so it never splits a
// multi-byte character.
func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	head := max * 6 / 10
	tail := max - head
	omitted := len(r) - head - tail
	return string(r[:head]) + fmt.Sprintf("\n… [%d chars truncated] …\n", omitted) + string(r[len(r)-tail:])
}

// assemble joins the turns under a header, within a total-character budget. Under budget,
// everything is kept. Over it, the first turn (the original task) is always retained and
// the most recent turns are added until the budget fills, with the dropped middle marked —
// recent turns hold the live problem, the first holds the task, and the middle is the most
// expendable.
func assemble(turns []string, budget int) string {
	const header = "# Session transcript (redacted: <system-reminder> blocks, images, and sub-agent sidechains removed; long blocks truncated)\n\n"
	joined := strings.Join(turns, "\n\n")
	if len(header)+len(joined) <= budget || len(turns) == 1 {
		return header + joined
	}
	const marker = "_[… earlier turns omitted to fit budget …]_"
	room := budget - len(header)
	first := turns[0]
	if len(first) > room-len(marker) {
		first = truncate(first, room-len(marker))
	}
	remaining := room - len(first) - len(marker)
	var tail []string
	used := 0
	for i := len(turns) - 1; i >= 1; i-- {
		if used+len(turns[i])+2 > remaining {
			break
		}
		tail = append([]string{turns[i]}, tail...)
		used += len(turns[i]) + 2
	}
	parts := []string{first}
	if len(tail) < len(turns)-1 {
		parts = append(parts, marker)
	}
	parts = append(parts, tail...)
	return header + strings.Join(parts, "\n\n")
}

func usage(w io.Writer) {
	fmt.Fprint(w, `session-transcript — dump the current Claude Code session as a redacted transcript.

  session-transcript | less        # inspect; codex-consult --context runs this for you

  -f, --file <path>    transcript JSONL (default: auto-locate via $CLAUDE_CODE_SESSION_ID)
  -b, --budget <n>     max output chars; keeps first + most recent turns (default: 120000)
      --block-max <n>  max chars per block before head/tail truncation (default: 3000)
      --no-thinking    exclude the agent's own reasoning (thinking) blocks
  -h, --help           this help

Strips <system-reminder> blocks, images, and sub-agent sidechains; truncates long blocks.
stdout = the transcript.   stderr = source path + turn/char counts.
`)
}
