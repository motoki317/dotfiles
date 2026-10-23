Read `memory/todo.md` in your workspace. If it is missing or lists no items, reply NO_REPLY.

For each item, check its current state at its source: the Slack thread, the pull request, or the ClickUp task. Delete the item if it is done or cancelled. Judge done by the done condition written in the item.

Report an item marked `[urgent]` on every run until a person starts handling it. Report any other item that is blocked, waits on someone, or has had no progress for a day. If the item has a Slack thread, post the report in that thread with the message tool. Otherwise, put it in your reply, which goes to the report channel. Note the report date in the item. Except for unhandled `[urgent]` items, do not report an item again on the same day unless its state changed.

If nothing needs a report, reply NO_REPLY.
