# Prompting current frontier models — house findings

Distilled 2026-07-17 from the vendor guides below plus an OpenAI Codex cross-review, during the house-asset optimization pass (see the dotfiles commits of that date); findings 12–13 added 2026-09-02 from the Claude Fable 5.1 guide. Read this before writing or trimming any skill or rule.

## Sources

- Claude prompting best practices: https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices
- Prompting Claude Fable 5: https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-fable-5
- Prompting Claude Fable 5.1: https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-fable-5-1
  (append `.md` to these three URLs for raw markdown fetchable with `ax --body`)
- Codex prompting guide: https://developers.openai.com/cookbook/examples/gpt-5/codex_prompting_guide
- GPT-5.6 model guide: https://developers.openai.com/api/docs/guides/latest-model?model=gpt-5.6
- Codex reference starter prompt: https://github.com/openai/codex/blob/main/codex-rs/core/gpt-5.1-codex-max_prompt.md
- Anthropic frontend-design skill (canonical wording): https://github.com/anthropics/claude-code/blob/main/plugins/frontend-design/skills/frontend-design/SKILL.md
- Codex customization and AGENTS.md docs (cited by the cross-review): https://learn.chatgpt.com/docs/customization/overview, https://learn.chatgpt.com/docs/agent-configuration/agents-md

## Findings

1. **Leaner prompts measurably win.** OpenAI reports 10–15% eval gains with 41–66% fewer tokens from removing repeated instructions; Anthropic says skills written for prior models are "often too prescriptive" and can degrade output. Both treat the numbers as directional — validate on your own tasks.
2. **State each rule once.** Repeated caution rules ("ask first", "do not mutate") cause spurious approval pauses; a sentence living in two always-loaded files is a future contradiction.
3. **Brief steering beats enumeration.** One instruction with its why now replaces a list of named edge cases. Give intent, hard constraints, and success criteria — not step-by-step prescriptions.
4. **Dial back aggressive language.** ALL-CAPS, "CRITICAL: you MUST", and trigger-spam descriptions ("prefer this over everything") now cause overtriggering. Describe what the skill does plus its genuine triggers, plainly.
5. **Keep what encodes a product requirement or corrects a measured gap; cut trained-in textbook content** and anti-laziness prompting ("if in doubt, use X").
6. **Fan out subagents proportionally.** Both vendors warn against blanket delegation: inline analysis for small scope, 1–3 agents when independence or specialization helps, full panels only for audits.
7. **Never instruct the model to echo its reasoning in the response** — on Fable 5 this can trigger `reasoning_extraction` refusals. Ask for findings and evidence, not thought transcripts.
8. **No fabricated precision.** 0–100 confidence scores and rigid multi-section report templates are ceremony; keep an output template only where a downstream consumer parses it.
9. **Vendor-endorsed keepers:** ground progress claims in tool results from the session; separate assessment from change ("report and stop" vs "implement and verify"); overengineering guardrails; one explicit parallel-tool-calls line (Codex harnesses still need it).
10. **Optimize by loading level.** Description is charged to every session of every agent; SKILL.md body per invocation; references/ only when read. Trim in that order — file length alone is not context cost.
11. **Change one group at a time** and compare fixed representative prompts before and after (`~/.claude/META.md` states the same principle; OpenAI's gains were measured this way).
12. **Remove update suppressors and anti-formatting rules before adding anything.** Claude Fable 5.1 narrates less between tool calls and formats less than Fable 5, so "hold findings for the final response", "don't narrate", and "no bullets/headers/bold" now strip output the reader wanted. Delete them first; if more narration is still wanted, say *when* user-facing text is wanted, not how much.
13. **A parallel-tool-calls line belongs next to the request, not in a standing file.** Fable 5.1 batches implied reads less than Fable 5, and one sentence appended per turn moves the rate far more than the same text in a rule or skill; the Claude Code harness already appends that sentence. Keep the one line in `rules/conventions.md` for Codex and add no copies to skills.

## Cross-agent notes (assets shared between Claude Code and Codex)

- Claude-only skills must not be exposed to Codex: `codex-work` and `codex-advise` delegate TO Codex (recursive when Codex is the caller), and names can collide with Codex's built-in `.system` skills. `sync-claude-skills` owns the exclusion list.
- Codex loads skills progressively like Claude (metadata → SKILL.md → resources), and closely adheres to AGENTS.md, which it injects as user-role messages once per run.
- Codex-side environments lack Claude-specific tools (Context7 MCP, Agent tool types); shared skills should name capabilities, not harness-specific tool names, or qualify them with "when available".
