# OpenClaw agents served by the one gateway in ./default.nix. One block per
# agent; ./README.md "Add an agent" lists the steps outside this file. Each
# agent's Slack channel and member IDs live outside git, in
# ~/.secrets/openclaw/agents.nix.
{
  main = {
    # Merged into agents.entries.<id>, e.g. { model.primary = "anthropic/claude-sonnet-5-5"; }.
    entry = { };
  };
}
