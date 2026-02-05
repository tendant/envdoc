# envdoc

Safe environment variable inspection for Go. Debug env vars in locked-down environments without leaking secrets.

envdoc exposes **metadata only** (presence, length, type validity, fingerprints) — never raw values.

## Quick Start

```bash
# Run instantly — dumps metadata for all env vars
go run ./cmd/envdoc

# With validation rules
go run ./cmd/envdoc -rules rules.yaml

# Docker
docker run --rm wang/envdoc
```

### Output

```
envdoc: key=DB_HOST present=true len=9 valid=true
envdoc: key=DB_PORT present=true len=4 valid=true
envdoc: key=DB_PASSWORD present=true len=32 fp=9f2c1a2b required=true valid=true
envdoc: key=FEATURE_X present=false valid=true
```

No raw values. No secrets. Safe for centralized logging.

## Install

```bash
# As a Go library
go get github.com/tendant/envdoc

# As a CLI binary
go install github.com/tendant/envdoc/cmd/envdoc@latest

# Docker
docker pull wang/envdoc
```

## Usage

### CLI

```bash
# Dump all env var metadata (default, no config needed)
envdoc

# Validate specific vars with a rules file
envdoc -rules rules.yaml

# Print version
envdoc -version
```

### Go Library

```go
package main

import "github.com/tendant/envdoc"

func main() {
    // Dump all env var metadata
    report, err := envdoc.Run()

    // Or with rules
    rules, _ := envdoc.LoadRulesFile("rules.yaml")
    report, err = envdoc.Run(envdoc.WithRules(rules))
}
```

## Rules File

Define validation rules in YAML:

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

  - key: APP_URL
    required: false
    type: url

  - key: FEATURE_FLAGS
    type: json

  - key: LOG_LEVEL
    allowed: [debug, info, warn, error]

  - key: SERVICE_ID
    regex: "^svc-[a-z0-9]+$"
```

### Supported Types

| Type | Validates |
|------|-----------|
| `string` | Any string (default) |
| `int` | Integer (e.g. `8080`) |
| `bool` | Boolean (`true`, `false`, `1`, `0`) |
| `duration` | Go duration (e.g. `30s`, `5m`) |
| `url` | URL with scheme and host |
| `json` | Valid JSON |

### Rule Options

| Field | Type | Description |
|-------|------|-------------|
| `key` | string | Environment variable name (required) |
| `required` | bool | Fail if missing |
| `type` | string | Expected type |
| `min_len` | int | Minimum value length |
| `max_len` | int | Maximum value length |
| `regex` | string | Regex the value must match |
| `allowed` | list | Allowed values |
| `secret` | bool | Override secret classification |
| `fingerprint` | bool | Override fingerprint behavior |

## Configuration

All configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ENVDOC_FAIL_FAST` | `false` | Exit non-zero if required vars are missing or invalid |
| `ENVDOC_ENABLE_HTTP` | `false` | Start HTTP debug endpoint |
| `ENVDOC_TOKEN` | | Bearer token for HTTP endpoint |
| `ENVDOC_EXPIRES_AT` | | Endpoint expiry time (RFC3339) |
| `ENVDOC_LISTEN_ADDR` | `127.0.0.1:9090` | HTTP listen address |
| `ENVDOC_DUMP_ALL` | `true`* | Dump all env var metadata |
| `ENVDOC_DUMP_ALL_FINGERPRINT` | `false` | Add fingerprints for non-secret vars |

\* Dump-all is automatic when no rules file is provided.

## HTTP Debug Endpoint

```bash
# Enable the endpoint
export ENVDOC_ENABLE_HTTP=true
export ENVDOC_TOKEN=my-secret-token
export ENVDOC_EXPIRES_AT=2026-02-06T00:00:00Z
envdoc -rules rules.yaml

# Query it
curl -H "Authorization: Bearer my-secret-token" http://127.0.0.1:9090/debug/env
```

Returns a JSON report with live inspection data on each request.

## Modes

| Mode | When | Behavior |
|------|------|----------|
| **Dump-all** | No rules file | Inspect all env vars, metadata only |
| **Allowlist** | Rules file provided | Only inspect declared vars with validation |

## Kubernetes Deployment

See [`deploy/k8s/`](deploy/k8s/) for sample manifests:

- **Sidecar** — runs alongside your app, shares env vars via `envFrom`
- **Init container** — validates env vars before app starts

```bash
kubectl apply -f deploy/k8s/sidecar.yaml
```

## Development

```bash
make              # show all targets
make build        # compile CLI binary
make run          # run with dump-all (default)
make run RULES=rules.yaml  # run with rules
make test         # run all tests
make test-cover   # run tests with coverage
make docker-build # build Docker image
make deploy       # build and push to Docker Hub
```

## Security

envdoc is designed to be safe for production use:

- **Never exposes raw values** — only metadata (presence, length, validity)
- **Secret classification** — auto-detects secret-like variable names
- **Fingerprints gated** — disabled for secret-like vars by default
- **Token auth** — optional Bearer token for HTTP endpoint
- **Expiry** — time-limited debug endpoints
- **Internal binding** — defaults to `127.0.0.1`

## License

MIT
