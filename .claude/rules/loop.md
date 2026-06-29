# Operating Loop

- Default flow: Explore → Plan → Implement → Verify → Commit → Ship. Auto-advance through Implement, Verify, and Commit on your own; Ship — pushing and opening PRs — stays explicitly user-triggered.
- Verify before declaring done: run the available tests/build/lint, then `feedback` (static review) and `qa-review` (dynamic QA, when there is runnable behavior). Report the commands run and what they returned, not just "done".
- Ask only when blocked by: a value or preference you can't infer, missing secrets/credentials, a scope boundary you can't infer, or a destructive or irreversible action (data loss, force-push, deploy, publishing to an external service).
