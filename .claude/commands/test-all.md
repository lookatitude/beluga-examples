---
name: test-all
description: Run the root-CI gate chain on every example locally (CI parity) and report OK / FAIL / NOT RUN per example.
---

Test all examples.

This is a local, read-only check. GitHub CI does not call this command; it runs `.github/workflows/ci.yml` directly. This command reproduces those gates so local claims match CI. The gate chain and pinned versions are in `.guild/skills/examples-ci-gates/SKILL.md`.

## Workflow

### Step 1: Discover
Find every directory with a `go.mod`, and among them the scaffolded ones (with `.beluga/project.yaml`).

### Step 2: Test
`@agent-example-tester` (Guild: `qa`) runs the following on each example, stopping at the first failure: `go build -o /dev/null ./...`, `go vet ./...`, `go test -race -count=1 ./...`, `go mod tidy -diff`, `govulncheck ./...` (v1.9.0) and `gosec -quiet ./...` (v2.30.0). Scaffolded examples also run `make check`, `make eval-ci` and `docker build .`. Use Go 1.27.x with `GOTOOLCHAIN=local` (CI uses 1.27.3).

### Step 3: Report
Print results per example (OK, FAIL with the first error line, or NOT RUN with the missing tool), the Go and tool versions used, and a summary line: `N/M examples passing`.

If any example fails or any gate did not run, say so first in the report. Do not describe the repo as green.
