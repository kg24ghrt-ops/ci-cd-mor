# ci-cd-mor

An agent-focused CLI for managing GitHub Actions workflows, built with Go and the Cobra library. Wraps the `gh` CLI to provide a streamlined interface for starting, watching, stopping, and inspecting workflow runs.

## Requirements

- Go 1.22+
- GitHub CLI (`gh`) installed and authenticated (`gh auth login`)
- Network access to GitHub

## Installation

```sh
go install ./cmd/agent-ci
```

Or build locally:

```sh
go build -o agent-ci ./cmd/agent-ci
```

## Commands

### `agent-ci start`

Starts a GitHub Actions workflow.

```sh
agent-ci start --workflow .github/workflows/ci.yml --branch main
```

Flags:
- `--workflow` (required) — Path to the workflow YAML file.
- `--branch` — Git branch to run the workflow against. Default: `main`.

Output: JSON object containing `runId`, `status`, `url`, `conclusion`, `headSha`, and `headBranch`.

---

### `agent-ci watch`

Polls a workflow run until it completes or a timeout is reached.

```sh
agent-ci watch --run-id 1234567 --interval 10 --max-wait 60
```

Flags:
- `--run-id` (required) — The workflow run ID.
- `--interval` — Polling interval in seconds. Default: `5`.
- `--max-wait` — Maximum time to wait in minutes. Default: `30`.

The command prints status updates as the run progresses and exits with the run's URL when complete. Exits non-zero on timeout.

---

### `agent-ci logs`

Fetches logs for a workflow run.

```sh
agent-ci logs --run-id 1234567 --errors-only --limit 200
```

Flags:
- `--run-id` (required) — The workflow run ID.
- `--errors-only` — Filter output to lines containing error/fail/fatal/panic/exception/traceback/exit code.
- `--limit` — Maximum number of log lines to return (from the tail).

---

### `agent-ci stop`

Cancels a running workflow.

```sh
agent-ci stop --run-id 1234567
```

Flags:
- `--run-id` (required) — The workflow run ID to cancel.

Output: JSON object with run details, or a minimal `{"runId": ..., "status": "cancelled"}` response if the run was already terminal.

## Architecture

```
cmd/agent-ci/main.go     — entry point, registers all subcommands
pkg/cmd/
  start.go               — workflow trigger via `gh workflow run`
  watch.go               — polling loop with interval and max-wait
  logs.go                — log retrieval with error filtering and line limiting
  stop.go                — workflow cancellation via `gh run cancel`
```

All commands shell out to the `gh` CLI and parse its JSON output.

## Development

```sh
go mod tidy      # sync dependencies
go build ./...   # build
go test ./...    # run tests
```

## CI/CD

- Push to `main` or open a PR triggers the `ci` workflow (build, lint, test).
- Push a `v*` tag triggers the `release` workflow (GoReleaser).