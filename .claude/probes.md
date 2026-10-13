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

Provenance:
- dismissal = backfill incident (commit 9fb5dbd)
- yagni and prose-route = verbose-agents episode (commit 0fd7ce6). The 2026-10-13 audit (prose#7) moved prose-route into a `/pr` dry run.
- deslop-route = the nanaism/yomiyasu integration (2026-10-01)
- web-tool and text-tool = silent-success overrides that can never be lazy-loaded
