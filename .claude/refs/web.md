# Web fetching (`ax`)

- `ax <url> --md --all --budget <n>` gives the page's exact text. `--all` matters: without it extract modes stop at 50 blocks. `--budget` caps tokens; truncation is announced.
- `ax` accepts curl-style flags.
- Empty output = JS-rendered page: append `.md` to the URL (docs sites often serve raw markdown — read it with `--body`, the raw-response mode), else WebFetch.
- WebFetch answers through a sub-model summary — use it for a single fact in a huge page, or claude.ai URLs only it can reach.
