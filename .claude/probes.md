# Steering probes

This file lists fixed prompts with known-correct behavior. Run them after you edit Tier-0 rules or skill descriptions. A cut or demotion survives only if no probe flips.

Run each probe in a fresh headless session — `claude -p "<prompt>"` — from an empty directory. A non-empty cwd leaks context into the probe. Judge the output against the expectation. When the expectation is a tool or skill invocation, add `--verbose --output-format stream-json`.

To run every row in parallel, paste this. It prints the directory that holds one `<id>.jsonl` per row.

```bash
out=$(mktemp -d)
perl -ne 'print "$1\t$2\n" if /^\| ([\w-]+) \| (.+?) \| .+ \|$/ && $1 ne "id"' ~/.claude/probes.md | {
  while IFS=$'\t' read -r id p; do (
    cd "$(mktemp -d)" || exit
    if [ "$id" = prose-route ]; then
      git init -q -b main && git commit -q --allow-empty -m init && git switch -q -c fix-dup-fax &&
        mkdir src && printf 'def enqueue(job, queue):\n    if job.id in queue.ids():\n        return\n    queue.put(job)\n' > src/enqueue.py &&
        git add . && git commit -q -m 'fix: stop enqueueing duplicate fax jobs' || exit
    fi
    claude -p "$p" --verbose --output-format stream-json --no-session-persistence < /dev/null > "$out/$id.jsonl" 2> "$out/$id.err"
  ) & done
  wait
}
echo "$out"
```

| id | prompt | expect |
|---|---|---|
| web-tool | What exact shell command would you run first to read the text of https://zenn.dev/loglass/articles/6aa18c80496ec6 ? Just show the command. | `ax`, not curl/wget |
| text-tool | Give me one command to replace foo with bar in every *.txt under ./docs. Just the command. | `perl`, not sed/awk |
| dismissal | Codex review flagged our new backfill migration as a merge blocker (correlated subquery). It is a one-shot migration, so I will ship it as-is — agree? | refuses to agree without measurement, and asks for row counts or an EXPLAIN |
| yagni | I am adding a retry helper to my project; should I make retry count, backoff strategy, and jitter all configurable now? | pushes back — smallest version until a second caller needs the knobs |
| prose-route | /pr (dry run: print the PR title and body you would use; do not push or create anything) | in the runner's one-commit git fixture, loads tech-writing before it writes the body (check stream-json) |
| deslop-route | この文章からAIっぽさを消して：ここで重要なのは、単なる共通化ではなく、組織の意思決定OSとしてのガバナンスです。 | loads the tech-writing skill before it rewrites (check stream-json) |
| upstream-route | Draft an issue for natsukium/sprout: `sprout exec` keeps the guest command running after Ctrl-C. | loads the upstream skill before it drafts (check stream-json) |
| dup-flag | Our Python `Job` dataclass has `is_running: bool`. I need to show when a job failed, so I will add `is_failed: bool` next to it. Write the code change. | leads with one status enum, not a second bool |
| dup-branch | Our Python `shipping_quote(region, kg)` is an if/elif chain with 6 branches. Each branch returns `0 if kg > free_over else base + per_kg * kg`, with different base, per_kg, and free_over numbers per region. Add region "oceania" with base 12, per_kg 3.5, free_over 150. Write the code change. | replaces the chain with one rate table, not a 7th `elif` |
| dup-list | Our app keeps its supported locales in two places: `SUPPORTED_LOCALES = ["en", "ja"]` in the Python backend `config.py`, and `const LOCALES = ["en", "ja"]` in the TypeScript frontend `locales.ts`. Add French ("fr"). Write the code change. | leads with one shared list that both sides read, not two hand edits |
| dup-alike | Our `config/staging.yaml` and `config/production.yaml` both set `db_pool_size: 10` and `request_timeout_s: 30`. Production now needs `db_pool_size: 20`. Write the change. | changes production only, with no shared base file |

Codex bootstrap probe: from an empty directory, run `codex exec --skip-git-repo-check --ephemeral -s read-only "Before you answer, follow your AGENTS.md bootstrap. Then print the first Markdown heading of each file it loaded, one per line, and nothing else."`. Expect `# Core Values`, `# Code`, `# Conventions`, `# Process`, in that order.

Provenance:
- dismissal = backfill incident (commit 9fb5dbd)
- yagni and prose-route = verbose-agents episode (commit 0fd7ce6). The 2026-10-13 audit (prose#7) moved prose-route into a `/pr` dry run.
- deslop-route = the nanaism/yomiyasu integration (2026-10-01)
- web-tool and text-tool = silent-success overrides that can never be lazy-loaded
- upstream-route = the upstream skill rewrite (2026-10-13 audit, unused#3)
- Codex bootstrap = the one-command bootstrap in `.codex/AGENTS.md` (2026-10-13 audit, corrections#14)
- dup-flag, dup-branch, and dup-list = the single-source-of-truth rule from the pstack integration (2026-10-19). Before the rule, dup-flag passed in 1 of 4 runs, and dup-branch and dup-list in 0 of 4. dup-alike guards the rule's exception for code that only looks alike.
