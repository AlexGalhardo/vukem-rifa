# Security audit, 2026-11-07

Target: `vukem-rifa-vuejs` at branch `build/vukem-v0.39.0`.
Method: the vendored Cloudflare `security-audit` skill (`.claude/skills/security-audit`), full audit mode, profile `quick`.

Machine-readable records: [findings](security-audit-2026-11-07.findings.json), [coverage ledger](security-audit-2026-11-07.coverage.json).

## Scope and coverage statement

Partial coverage by design. Profile `quick`: one reviewer, coarse ledger units (7), one hunting pass, no follow-up waves. Source review only; no live endpoints were probed and no dependency was added. Scanners: `go tool govulncheck ./...` (0 reachable vulnerabilities) and `bun audit --prod` (none).

In scope: payment webhook (secret, signature, dedupe), raffle draw algorithm, public reservation and checkout, API keys and idempotency, CSV export, compose stack, Dockerfile and the GitHub workflows.

Not reviewed in depth: the Vue pages, mail templates, AbacatePay client internals, admin payout and refund flows, and the vukem framework itself.

## Summary

| Severity | Confirmed | Needs validation | Rejected |
| --- | --- | --- | --- |
| Critical | 0 | - | - |
| High | 0 | - | - |
| Medium | 0 | - | - |
| Low | 2 | - | - |
| Informational | 1 | - | - |
| Unrated | - | 0 | 2 |

## Findings

### SA-001 Dev compose ports published on all interfaces with default credentials (low)

- Fingerprint: `vukem-rifa-vuejs:infra:compose-ports-all-interfaces`
- CWE: CWE-1327 (binding to an unrestricted IP address)
- Location: `compose.yaml:31`, `:45`, `:57`, `:58`
- Impact: Postgres (password `rifa`), the app with a default webhook secret and Mailpit are reachable from the local network while the dev stack runs. The same mapping comes from the vukem scaffold (workspace `vukem-feedback.md`, VE-8).
- Fix: Every published port is prefixed with `127.0.0.1:`. Regression: `TestComposePublishesPortsOnLoopbackOnly`. Status: fixed in `f5bed55`.

### SA-002 Workflow checkouts persist the GitHub token (informational)

- Fingerprint: `vukem-rifa-vuejs:ci:checkout-persists-credentials`
- CWE: CWE-522 (insufficiently protected credentials)
- Location: `.github/workflows/ci.yml:22`, `release.yml:15`
- Impact: The job token stays in `.git/config` while later steps run repository code. Token scope is limited by each workflow's `permissions`.
- Fix: `persist-credentials: false` on every checkout. Status: fixed in `83e912f`.

### SA-003 Third-party Actions pinned by tag (low)

- Fingerprint: `vukem-rifa-vuejs:ci:actions-pinned-by-tag`
- CWE: CWE-829 (inclusion of functionality from an untrusted control sphere)
- Location: `.github/workflows/ci.yml`, `release.yml`
- Impact: A moved upstream tag would run attacker code with the job token; release workflows hold `packages: write`.
- Fix: Every `uses:` is pinned to a commit SHA with the tag in a comment. Status: fixed in `eb11bb6`.

## Rejected after validation

- `vukem-rifa-vuejs:payments:webhook-forgery`: Fails closed on every check.
- `vukem-rifa-vuejs:draws:biased-or-predictable-draw`: No modulo bias and no predictable seed.

Hardening notes, not findings: the dev compose file ships `ABACATEPAY_WEBHOOK_SECRET: local-webhook-secret`; it is only reachable on loopback after the compose fix.

## Verification

Findings and ledger passed the skill validators (`validate-findings.cjs`: 5 records; `validate-coverage-ledger.cjs`: 7 units), run on Windows Node with a preload shim for the missing POSIX no-follow constants. The skill's own CLI tests fail 7 cases on Windows for the same reason, as in `vukem-email`.
