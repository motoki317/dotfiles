# Web fetching (`ax`)

- `ax <url> --md --all --budget <n>` gives the page's exact text. Without `--all`, extract modes stop at 50 blocks. `--budget` caps tokens, and `ax` announces truncation.
- `ax` accepts curl-style flags.
- Empty output = JS-rendered page: append `.md` to the URL (docs sites often serve raw markdown — read it with `--body`, the raw-response mode), else WebFetch.
- WebFetch answers through a sub-model summary. Use it for a single fact in a huge page, or for claude.ai URLs that only it can reach.
