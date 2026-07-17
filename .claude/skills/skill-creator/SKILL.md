---
name: skill-creator
description: Create, improve, and evaluate skills. Use when the user wants to create a skill from scratch, edit or optimize an existing one, run evals or benchmarks on it, or optimize its description for better triggering accuracy.
---

# Skill Creator

The loop: capture intent → draft → run test prompts with and without the skill → the user reviews the outputs → improve → repeat until satisfied. Figure out where in this loop the user already is and help from there; if they say "just vibe with me", skip the eval machinery.

Before writing or revising any skill, read `references/model-prompting-insights.md` — house findings on prompting current frontier models, with sources.

**Toolkit**: the eval scripts, grader/analyzer agent instructions, schemas, and viewer live in the upstream plugin copy at
`~/.claude/plugins/marketplaces/claude-plugins-official/plugins/skill-creator/skills/skill-creator/` (referenced below as `<toolkit>`). If that path is missing, refetch it from https://github.com/anthropics/claude-plugins before running evals.

## Creating a skill

### Capture intent

The conversation may already contain the workflow ("turn this into a skill") — extract the tools, steps, corrections, and formats from history first, and have the user confirm the gaps. Establish:

1. What should the skill enable Claude to do?
2. When should it trigger?
3. Expected output format?
4. Test cases? Objectively verifiable outputs (file transforms, code generation, fixed workflows) benefit; subjective ones (writing style, art) usually don't. Suggest a default, let the user decide.

Ask about edge cases, formats, example files, success criteria, and dependencies before writing test prompts. Research via available MCPs or subagents in parallel where that reduces the burden on the user.

### Anatomy and loading costs

```
skill-name/
├── SKILL.md          # YAML frontmatter (name, description) + instructions
└── optional: scripts/ (executable), references/ (loaded as needed), assets/ (used in output)
```

Each level costs context at a different rate: name + description are charged to every session; the SKILL.md body to every invocation (<500 lines, ideally far less); references and scripts only when read. Write the least that carries the meaning: keep the description to what the skill does plus its triggers, keep the body to operative rules, and push detail down a level — cut restatements of the description, justification prose, and summaries of content a reference file already owns. When a skill supports multiple domains, one reference file per variant, read selectively.

### Writing style

- The description is the triggering mechanism — cover what the skill does AND when to use it. Name genuine triggers concretely; don't pad with keyword spam or "prefer this over everything" — current models overtrigger on pushy descriptions.
- Explain why instead of heavy-handed MUSTs; ALL-CAPS is a yellow flag. Current models generalize well from a stated reason.
- Imperative form. Examples earn their place when they encode a requirement prose can't.
- No surprises: a skill's contents must match what its description claims; never malware or misleading behavior.

### Test cases

Write 2–3 realistic prompts a real user would actually type; confirm them with the user; save to `evals/evals.json` (schema: `<toolkit>/references/schemas.md`). Assertions come later, while the runs execute.

## Running and evaluating

One continuous sequence — don't stop partway. Workspace: `<skill-name>-workspace/` as a sibling of the skill directory, `iteration-<N>/<eval-name>/` per test case, created as you go.

1. **Spawn all runs in one turn** — per test case, one with-skill subagent and one baseline (creating → no skill; improving → a pre-edit snapshot: `cp -r <skill> <workspace>/skill-snapshot/`). Give each: the skill path (or none), the task, input files, the output directory (`with_skill/outputs/`, `without_skill/` or `old_skill/outputs/`), and which outputs to save. Write an `eval_metadata.json` per case (descriptive `eval_name`, empty `assertions` for now).
2. **Draft assertions while runs execute** — objectively verifiable, descriptively named; leave subjective qualities to human review. Update the metadata files and `evals/evals.json`.
3. **Capture timing from each completion notification** — save `total_tokens` and `duration_ms` to `timing.json` immediately; the notification is the only source.
4. **Grade and aggregate** — grade each run per `<toolkit>/agents/grader.md` into `grading.json`; its expectations array uses exactly the fields `text`, `passed`, `evidence` (the viewer depends on them). Script programmatic checks instead of eyeballing. Then:
   ```bash
   cd <toolkit> && python -m scripts.aggregate_benchmark <workspace>/iteration-N --skill-name <name>
   ```
   Do an analyst pass per `<toolkit>/agents/analyzer.md`: non-discriminating assertions, flaky evals, time/token trade-offs.
5. **Launch the viewer** — always `generate_review.py`, never hand-rolled HTML; add `--previous-workspace` from iteration 2 on; use `--static <path>` on headless setups:
   ```bash
   nohup python <toolkit>/eval-viewer/generate_review.py <workspace>/iteration-N \
     --skill-name <name> --benchmark <workspace>/iteration-N/benchmark.json > /dev/null 2>&1 &
   ```
   Tell the user: the Outputs tab collects per-case feedback, the Benchmark tab shows the stats; come back when done.
6. **Read `feedback.json`** when the user says they're done — empty feedback means fine. Kill the viewer.

## Improving

- **Generalize from the feedback.** The skill must work far beyond the few dev examples; prefer reframing or a different working pattern over fiddly overfit constraints.
- **Keep it lean.** Read the transcripts, not just the outputs — cut instructions that cause wasted work.
- **Explain the why.** Transmit understanding of what the user actually needs, not rote rules.
- **Bundle repeated work.** If every run hand-wrote the same helper script, ship it once in `scripts/`.

Then iterate: apply improvements → rerun all cases (with baselines) into `iteration-<N+1>/` → viewer with `--previous-workspace` → feedback → repeat, until the user is happy, the feedback comes back empty, or progress stalls.

For a rigorous "is v2 actually better?", run a blind comparison per `<toolkit>/agents/comparator.md` and `analyzer.md`: judge two outputs without revealing which is which.

## Description optimization

After the skill works, offer to optimize its triggering.

1. **Generate ~20 eval queries**, mixed should-trigger and should-not-trigger, as JSON (`{"query", "should_trigger"}`). Make them concrete and realistic (file names, casual phrasing, typos). Positives: varied phrasings, cases that never name the skill, cases where it competes with another tool but should win. Negatives: near-misses sharing keywords — not obviously-irrelevant strawmen.
2. **Review with the user** via `<toolkit>/assets/eval_review.html` (fill the `__EVAL_DATA_PLACEHOLDER__`, `__SKILL_NAME_PLACEHOLDER__`, `__SKILL_DESCRIPTION_PLACEHOLDER__` placeholders; the export lands in `~/Downloads/eval_set.json`). Bad queries make bad descriptions.
3. **Run the loop in the background**, reporting progress as it iterates:
   ```bash
   cd <toolkit> && python -m scripts.run_loop --eval-set <eval-set.json> \
     --skill-path <skill> --model <model-id-powering-this-session> --max-iterations 5 --verbose
   ```
4. **Apply `best_description`** (chosen on held-out test score); show the user before/after and the scores.

Triggering context: Claude consults a skill from its name + description alone, and only for tasks it can't trivially handle — eval queries must be substantive enough that the skill would actually help.

## Reference files

- `references/model-prompting-insights.md` — house findings on prompting current models; read before writing or trimming any skill
- `<toolkit>/references/schemas.md` — JSON schemas for evals.json, grading.json, benchmark.json
- `<toolkit>/agents/{grader,comparator,analyzer}.md` — subagent instructions for grading and comparison
