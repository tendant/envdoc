# envdoc — Safe Environment Variable Troubleshooting for Go

## Overview

**envdoc** is a small, embeddable Go component (or standalone binary) designed to help troubleshoot environment variable issues **without cluster access** and **without leaking secrets**.

It focuses on answering questions like:

- Does an environment variable exist?
- Is it empty or unexpectedly short/long?
- Does it parse as the expected type?
- Did the value change between deployments?
- Is the app running with the config version I think it is?

The core design principle is:

> **Never expose raw environment variable values.**

Instead, envdoc exposes **metadata only**, with strict gating and redaction.

---

## Goals

- Debug env vars in locked-down environments (no kubectl, no node access)
- Detect presence / absence of env vars
- Validate type and basic structure
- Surface misconfiguration early (fail-fast)
- Minimize information leakage
- Be safe to run in production when properly gated

## Non-Goals

- Dumping raw environment variable values
- Acting as a general-purpose secrets viewer
- Replacing secret managers or config systems

---

## Design Principles

1. **Allow-list by default**
   - Only explicitly declared keys are checked
   - Prevents accidental leakage

2. **Dump-all is metadata-only and opt-in**
   - No raw values
   - Heavy redaction
   - Explicit enable flags

3. **Fail fast when required config is broken**
   - Misconfiguration should crash early, not fail silently

4. **App-level introspection beats infra-level access**
   - The process itself is the source of truth

---

## Modes of Operation

### Mode 1: Allow-List (Default)

Only environment variables defined in rules are inspected.

Use cases:
- Normal production operation
- Compliance-sensitive systems
- CI / startup validation

### Mode 2: Dump-All Metadata (Opt-in)

Enumerates **all environment variables**, but exposes **metadata only**.

Required flags:

```text
ENVDOC_DUMP_ALL=true
ENVDOC_ENABLE_HTTP=true
```

Optional safety controls:

```text
ENVDOC_TOKEN=secret-token
ENVDOC_EXPIRES_AT=2026-02-05T20:00:00Z
```

### Mode 3: Dump-All Metadata + Fingerprints (Highly Restricted)

Adds short cryptographic fingerprints for **non-secret-like keys only**.

Required flags:

```text
ENVDOC_DUMP_ALL_FINGERPRINT=true
```

Fingerprints allow detecting config drift without revealing values.

---

## Information Exposed (Safe by Design)

For each environment variable, envdoc may expose:

| Field | Description | Risk |
|------|-------------|------|
| key | Env var name | Low |
| present | Exists or not | Low |
| length | Length of value | Low–Medium |
| required | Required by rules | Low |
| valid | Passed validation | Low |
| problems | Validation errors | Low |
| secret_like | Heuristic classification | Low |
| fingerprint | Short hash prefix (opt-in) | Medium |

**Never exposed:**
- Raw values
- Partial values
- Prefixes / suffixes

---

## Secret Classification

In dump-all mode, envdoc automatically classifies variables as `secret_like` based on name patterns:

```regex
(?i)(password|passwd|pwd|secret|token|apikey|api[_-]?key|private[_-]?key|ssh|pem|cert|jwt|session|cookie|authorization|bearer|kms|encrypt|signing)
```

Rules:
- `secret_like=true` → fingerprints disabled by default
- Length may still be shown

---

## Validation Rules (Allow-List Mode)

Rules are defined in YAML:

```yaml
rules:
  - key: DB_HOST
    required: true
    type: string
    min_len: 1

  - key: DB_PORT
    required: true
    type: int

  - key: DB_PASSWORD
    required: true
    type: string
    secret: true
    min_len: 16
    fingerprint: true
```

Supported types:
- string
- int
- bool
- duration
- url
- json

Additional checks:
- min/max length
- regex match
- allowed values
- whitespace trimming detection

---

## Output Channels

### 1. Startup Logs

One-line-per-variable summary at application startup:

```text
envdoc: key=DB_PASSWORD present=true len=32 fp=9f2c1a2b valid=true
envdoc: key=FEATURE_X present=false required=false valid=true
```

Safe for centralized logging when secrets are redacted.

---

### 2. HTTP Debug Endpoint

Endpoint:

```text
GET /debug/env
```

Features:
- JSON report
- Token-protected (optional)
- Can be disabled entirely
- Intended for internal access only

---

### 3. Fail-Fast Mode

When enabled:

```text
ENVDOC_FAIL_FAST=true
```

The app exits non-zero if:
- Required env vars are missing
- Type validation fails

This prevents pods from running with broken configuration.

---

## Security Controls

Recommended protections:

1. **Explicit enable flags** (off by default)
2. **Token-based auth** for HTTP endpoint
3. **Internal-only binding** (e.g. 127.0.0.1)
4. **Expiry time** for debug endpoints
5. **Rate limiting** (optional)
6. **Audit logging of access (not contents)**

---

## Deployment Patterns

### Pattern A: Embedded Library (Recommended)

- envdoc runs inside your app
- Exact same env context
- Most accurate results

### Pattern B: Sidecar (Limited)

- Sidecar has its own env
- Requires duplicated env injection
- Useful mainly for standardizing endpoints

### Pattern C: Init Container (Validation Only)

- Validates before app starts
- Same env duplication caveat

---

## Risks & Tradeoffs

| Risk | Mitigation |
|-----|-----------|
| Metadata leakage | Disable dump-all by default |
| Fingerprint correlation | Short hashes, opt-in only |
| Accidental exposure | Token + internal bind |
| Overconfidence | Logs + fail-fast + metrics |

---

## Recommended Defaults

```text
ENVDOC_MODE=allowlist
ENVDOC_ENABLE_HTTP=false
ENVDOC_FAIL_FAST=true
```

Enable dump-all only during incident response, with expiry.

---

## Future Enhancements

- Prometheus metrics (`env_present{key=...}`)
- Config signature hash for rollout verification
- Auto-rule generation from struct tags
- OpenTelemetry attributes
- Policy integration (OPA-style rules)

---

## Summary

envdoc treats environment variables as **untrusted input** and makes their state observable **without revealing secrets**.

It replaces `kubectl exec env` in locked-down production systems with:

- explicit validation
- safe introspection
- controlled debug surfaces


