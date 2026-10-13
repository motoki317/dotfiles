---
name: upstream
description: Draft and file an issue or PR on a third-party repository (not attmcojp or motoki317). Use when a bug or a needed change belongs to an upstream project.
argument-hint: "[owner/repo]"
---

# Upstream issue or PR

1. Search open and closed issues and PRs (`gh search issues --include-prs -R <owner/repo> "<keywords>"`), the docs, and the default branch for an existing report or fix. If one exists, report it instead of drafting.
2. Follow the repo's issue or PR template and CONTRIBUTING.
3. Draft with `/tech-writing`. Give the minimal facts and a minimal reproduction against the upstream project alone. Include no names, URLs, code, or data from attmcojp repos.
4. Show the draft. Post only after the user approves posting.
5. Before you post, run `cd ~ && gh api user -q .login` as a separate call. If it does not print `motoki317`, stop and tell the user. Post from `~`.
