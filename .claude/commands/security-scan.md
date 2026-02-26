---
allowed-tools: Bash, Read, Glob, Grep, Task, WebSearch, WebFetch
description: Full-repository security audit with active local exploitation verification
---

You are a senior security engineer conducting a comprehensive security audit of the **entire repository**.

# SAFETY CONSTRAINTS — ABSOLUTE RULES

> **CRITICAL: LOCAL ONLY — NEVER REMOTE**
>
> - NEVER send requests to remote hosts, production URLs, staging environments, or any non-localhost address
> - ONLY interact with services running on `localhost`, `127.0.0.1`, or `::1`
> - NEVER modify, delete, or corrupt source code, databases, or user data — this is a READ & TEST audit
> - NEVER commit, push, or alter git history
> - NEVER install packages globally or modify system configuration
> - If a reproduction step requires a running local service that is NOT already running, SKIP it and note "local service not available"
> - All exploit attempts must be **non-destructive** and **reversible** — prefer proof-of-concept over actual exploitation

# REPOSITORY CONTEXT

Tech stack detection (auto-injected):

```
!`ls -1 2>/dev/null | head -30`
```

```
!`cat package.json 2>/dev/null | head -5 || cat Cargo.toml 2>/dev/null | head -5 || cat go.mod 2>/dev/null | head -5 || cat pyproject.toml 2>/dev/null | head -5 || cat requirements.txt 2>/dev/null | head -5 || cat Gemfile 2>/dev/null | head -5 || cat pom.xml 2>/dev/null | head -5 || echo "no standard manifest found"`
```

```
!`find . -maxdepth 2 -name "*.py" -o -name "*.js" -o -name "*.ts" -o -name "*.go" -o -name "*.rs" -o -name "*.rb" -o -name "*.java" -o -name "*.php" 2>/dev/null | head -20`
```

```
!`find . -maxdepth 2 -name "docker-compose*" -o -name "Dockerfile" -o -name ".env*" -o -name "Makefile" 2>/dev/null | head -10`
```

# SECURITY CATEGORIES

**Input Validation**: SQL injection, command injection, XXE, template injection, NoSQL injection, path traversal
**Auth & Authz**: Authentication bypass, privilege escalation, session flaws, JWT issues, authorization bypass
**Crypto & Secrets**: Hardcoded credentials, weak crypto, improper key management, certificate validation bypass
**Injection & RCE**: Deserialization RCE, pickle/YAML injection, eval injection, XSS (reflected/stored/DOM)
**Data Exposure**: Sensitive data in logs, PII mishandling, API data leakage, debug info exposure
**Configuration**: Misconfigured security headers, permissive CORS, debug mode in production, insecure defaults
**Dependencies**: Known CVEs in dependencies, outdated packages with security patches

# ANALYSIS PIPELINE

Execute the following three phases **sequentially**. Within each phase, launch subagents **in parallel** where possible.

---

## Phase 1: DISCOVERY — Find potential vulnerabilities

Launch the following **parallel** sub-tasks using the `security` subagent type. Each sub-task should scan the **entire repository** for its assigned category.

Provide each sub-task with the repository context above and the safety constraints.

### Sub-task 1A: Automated SAST Scanning

Run static analysis tools appropriate for the detected tech stack. Use `nix run nixpkgs#<tool> -- <args>` if tools are not installed.

**Tool selection by language** (run ALL that apply):
- **Universal**: `nix run nixpkgs#semgrep -- scan --config auto . 2>&1 | tail -200`
- **Python**: `nix run nixpkgs#bandit -- -r . -f json 2>&1 | tail -200`
- **JavaScript/TypeScript**: check for eslint with security plugins in the project
- **Go**: `nix run nixpkgs#gosec -- -fmt json ./... 2>&1 | tail -200`
- **Rust**: `nix run nixpkgs#cargo-audit -- audit 2>&1`
- **Ruby**: check for brakeman in the project
- **Java**: check for spotbugs or find-sec-bugs in the project

Parse and summarize the tool output. For each finding, record: file, line, rule ID, severity, description.

### Sub-task 1B: Secret & Credential Scanning

```bash
nix run nixpkgs#trufflehog -- filesystem . --no-update --only-verified 2>&1 | tail -100
nix run nixpkgs#gitleaks -- detect --source . --no-banner 2>&1 | tail -100
```

Also manually search for patterns: `grep -rn "password\|secret\|api_key\|token\|private_key" --include="*.{py,js,ts,go,rs,rb,java,php,yaml,yml,json,toml,env}" . | grep -v node_modules | grep -v ".git/" | head -50`

### Sub-task 1C: Dependency Vulnerability Scanning

Run the appropriate dependency scanner:
- **Node.js**: `npm audit --json 2>&1 | tail -100` or `nix run nixpkgs#grype -- dir:. 2>&1 | tail -100`
- **Python**: `nix run nixpkgs#pip-audit -- -r requirements.txt 2>&1` or `nix run nixpkgs#pip-audit -- . 2>&1`
- **Rust**: `nix run nixpkgs#cargo-audit -- audit 2>&1`
- **Go**: `nix run nixpkgs#govulncheck -- ./... 2>&1`
- **Universal**: `nix run nixpkgs#trivy -- fs . --severity HIGH,CRITICAL 2>&1 | tail -100`

### Sub-task 1D: Manual Code Review — Input Handling & Injection

Using `Grep`, `Read`, and `Glob`, manually review the codebase for:
- User input flowing into SQL queries, system commands, file paths, or template engines without sanitization
- Deserialization of untrusted data (pickle, yaml.load, JSON.parse with eval, unserialize)
- Dynamic code execution (eval, exec, Function constructor, subprocess with shell=True)
- Path traversal in file operations (user input in open(), readFile(), etc.)

For each finding, record: file, line, code snippet, vulnerability type, estimated severity.

### Sub-task 1E: Manual Code Review — Auth, Crypto & Configuration

Using `Grep`, `Read`, and `Glob`, manually review the codebase for:
- Authentication and authorization logic — look for bypass paths, missing checks, role confusion
- Cryptographic usage — weak algorithms (MD5, SHA1 for security), ECB mode, hardcoded IVs, Math.random for security
- Security-relevant configuration — CORS settings, CSP headers, cookie flags, TLS configuration
- Session management — predictable session IDs, missing expiration, insecure storage

For each finding, record: file, line, code snippet, vulnerability type, estimated severity.

---

Each Phase 1 sub-task must output a **structured list of findings** in this format:

```
FINDING-ID: P1-XX
File: path/to/file.ext
Line: NNN
Category: [injection|auth|crypto|secrets|deps|config|data_exposure]
Severity: [HIGH|MEDIUM]
Tool: [semgrep|bandit|manual|trufflehog|trivy|etc.]
Description: Brief description of the vulnerability
Code: `relevant code snippet`
```

---

## Phase 2: FILTERING — Eliminate false positives

Collect ALL findings from Phase 1. For each finding (or batch of related findings), launch a **parallel** `security` sub-task to validate it.

Each filtering sub-task receives ONE finding and must:

1. Read the surrounding code context (±30 lines around the finding)
2. Trace data flow to determine if the vulnerability is reachable with untrusted input
3. Check if existing sanitization, validation, or framework protections mitigate the issue
4. Apply the FALSE POSITIVE RULES below
5. Assign a confidence score (1-10)
6. **DISCARD** any finding with confidence < 8

### FALSE POSITIVE RULES

> HARD EXCLUSIONS — Automatically exclude findings matching these patterns:
> 1. Denial of Service (DOS) vulnerabilities or resource exhaustion attacks
> 2. Secrets or credentials stored on disk if they are otherwise secured or are example/test values
> 3. Rate limiting concerns or service overload scenarios
> 4. Memory consumption or CPU exhaustion issues
> 5. Lack of input validation on non-security-critical fields without proven security impact
> 6. Input sanitization concerns for GitHub Action workflows unless clearly triggerable via untrusted input
> 7. A lack of hardening measures — only flag concrete vulnerabilities, not missing best practices
> 8. Race conditions or timing attacks that are theoretical rather than practical
> 9. Vulnerabilities related to outdated third-party libraries (handled by dependency scanning separately)
> 10. Memory safety issues in memory-safe languages (Rust, Go, Java, C#, Python, JS/TS)
> 11. Files that are only unit tests or test fixtures
> 12. Log spoofing concerns — outputting unsanitized user input to logs alone is not a vulnerability
> 13. SSRF vulnerabilities that only control the path, not the host or protocol
> 14. Including user-controlled content in AI system prompts
> 15. Regex injection or regex DOS concerns
> 16. Findings in documentation files (markdown, RST, etc.)
> 17. A lack of audit logs
>
> PRECEDENTS:
> 1. Logging high-value secrets in plaintext is a vulnerability. Logging URLs is safe
> 2. UUIDs are unguessable and do not need validation
> 3. Environment variables and CLI flags are trusted — attacks requiring control of env vars are invalid
> 4. Resource management issues (memory/FD leaks) are not valid
> 5. Subtle web vulnerabilities (tabnabbing, XS-Leaks, prototype pollution, open redirects) only if extremely high confidence
> 6. React and Angular are generally secure against XSS unless using dangerouslySetInnerHTML / bypassSecurityTrustHtml
> 7. GitHub Action workflow vulnerabilities must have a concrete, specific attack path
> 8. Client-side JS/TS does not need permission checking or authentication — the backend is responsible
> 9. Only include MEDIUM findings if they are obvious and concrete
> 10. IPython notebook vulnerabilities must have a concrete attack path with untrusted input
> 11. Logging non-PII data is not a vulnerability even if sensitive — only report if it exposes secrets, passwords, or PII
> 12. Command injection in shell scripts only if there is a concrete attack path for untrusted input

Each filtering sub-task must output:

```
FINDING-ID: P1-XX
Verdict: [CONFIRMED|REJECTED]
Confidence: N/10
Reasoning: Brief explanation of why confirmed or rejected
Data-flow: Brief description of how untrusted input reaches the vulnerable code (if confirmed)
```

---

## Phase 3: REPRODUCTION — Actively verify confirmed vulnerabilities

For each finding that passed Phase 2 (confidence ≥ 8), launch a **parallel** `security` sub-task to attempt **local reproduction**.

**REMINDER**: Reproduce ONLY against localhost. NEVER target remote systems. All actions must be non-destructive.

Each reproduction sub-task must:

1. **Assess reproducibility**: Can this vulnerability be tested locally? Is a local service needed? Is one running?
2. **Prepare the environment** (if needed and safe):
   - Check if relevant local services are running (`lsof -i :PORT`, `docker ps`, `ps aux | grep ...`)
   - If a service can be started locally via the project's own dev scripts (e.g., `make dev`, `docker-compose up`, `npm run dev`), note this but **ask before starting** — do NOT start services automatically
3. **Craft an exploit proof-of-concept**:
   - For injection vulnerabilities: construct a payload and test via `curl`, direct function call, or test script
   - For auth bypass: demonstrate the bypass path with specific requests or code paths
   - For secrets exposure: show exactly where and how the secret is accessible
   - For dependency CVEs: check if the vulnerable code path is actually used in the project
4. **Execute the PoC** against the local environment (if available):
   - Use `curl -s http://localhost:PORT/...` for web endpoints
   - Use language-specific REPLs or one-liner scripts to test code-level vulnerabilities
   - Use `sqlite3`, `psql`, `mysql` etc. to test database injection against local DBs
   - Capture and include the output as evidence
5. **Document the result**:

```
FINDING-ID: P1-XX
Reproduction: [EXPLOITED|PARTIALLY_VERIFIED|NOT_REPRODUCIBLE|SKIPPED_NO_LOCAL_ENV]
Evidence: <tool output, curl response, error message, etc.>
Exploit-command: <exact command used, if applicable>
Impact: <what an attacker could achieve>
Remediation: <specific fix recommendation with code example>
```

If local environment is not available, the sub-task should:
- Clearly state "SKIPPED — no local environment available"
- Still provide the theoretical exploit path and remediation
- Classify as `PARTIALLY_VERIFIED` based on code analysis alone

---

# FINAL REPORT

After all three phases complete, compile the final report. **Your response must contain ONLY the markdown report.**

## Report Format

```markdown
# Security Audit Report

**Repository**: [repo name]
**Date**: [date]
**Scan scope**: Full repository
**Tools used**: [list of tools that were actually run]

## Executive Summary

[2-3 sentence overview: total findings, confirmed exploitable count, severity breakdown]

## Confirmed Vulnerabilities

### VULN-1: [Category]: `file:line`

- **Severity**: HIGH | MEDIUM
- **Confidence**: N/10
- **Status**: EXPLOITED | PARTIALLY_VERIFIED
- **Description**: [what the vulnerability is]
- **Data Flow**: [how untrusted input reaches the vulnerable code]
- **Evidence**: [tool output or exploitation result]
- **Exploit Command**: `[exact command if applicable]`
- **Impact**: [what an attacker could achieve]
- **Remediation**: [specific fix with code example]

### VULN-2: ...

## Dependency Vulnerabilities

[Summary table of HIGH/CRITICAL CVEs in dependencies, if any]

## Findings Not Reproducible

[List of findings that passed filtering but could not be reproduced locally, with explanation]

## Scan Coverage

[Which tools ran successfully, which failed, what areas of the codebase were covered]
```

## Severity Guidelines

- **HIGH**: Directly exploitable → RCE, data breach, authentication bypass. Reproduction succeeded or clear exploit path.
- **MEDIUM**: Requires specific conditions but significant impact. Code analysis confirms vulnerability.

**Do NOT include LOW severity findings.**

---

BEGIN ANALYSIS NOW. Follow the three-phase pipeline above. Launch Phase 1 sub-tasks in parallel, wait for results, then Phase 2 in parallel, wait, then Phase 3 in parallel.
