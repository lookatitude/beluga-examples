---
derived_from_template: guild.skill_template.v1
name: examples-ci-gates
description: "What beluga-examples root CI (.github/workflows/ci.yml) enforces and how to reproduce it locally with CI parity: the per-example build/vet/race-test/tidy/govulncheck/gosec chain, the scaffolded-example matrix (make check, make eval-ci, docker build), the release-dispatch bump-on-release job and the version stamps it rewrites, and the pinned toolchain (Go 1.27.3 in CI over a 1.26 go.mod floor, govulncheck v1.9.0, gosec v2.30.0, golangci-lint v2.15.0, golang:1.27-alpine builders). Pulled by qa, example-maintainer, security and example-author. TRIGGER for \"run what CI runs\", \"test all examples\", \"is this example green\", \"which gosec version does CI use\", \"bump the examples to the new framework tag\", \"why did scaffolded-checks fail\". DO NOT TRIGGER for: per-example layout, naming, comment-block or code-style rules (examples-conventions); judging whether a scanner finding blocks merge (security); editing .github/workflows (out of scope for agent lanes — raise with the coordinator)."
when_to_use: "Before claiming any example is green, when bumping the framework version, or when diagnosing a root CI failure in beluga-examples."
type: specialist
---

# examples-ci-gates

Source of truth: `.github/workflows/ci.yml` (the only workflow GitHub runs here). If this skill and the workflow disagree, the workflow wins. Report the drift to the coordinator; do not edit the workflow from an agent lane.

# When to use it

- Running `/test-all` or validating a single example before commit.
- Bumping examples to a new framework tag (`/bump-framework`, or fixing a red `bump-on-release` PR).
- Diagnosing a failed `test-examples`, `scaffolded-checks` or `bump-on-release` job.

# When not to use it

- For the authoring rules (layout, naming, comments, style): load `examples-conventions`.
- For deciding whether a gosec or govulncheck finding blocks merge: the `security` specialist decides. This skill only says how the scanner is run.

# Required inputs

- The example directories to check: every `*/go.mod`, or the ones in the diff.
- Optional: a framework tag matching `^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$` when bumping.

# Output format

One line per example: `<example>  OK`, `FAIL: <gate> — <first error line>`, or `NOT RUN: <gate> (<reason>)` when a tool could not run. End with a summary line: `N/M examples passing`. An example with a NOT RUN gate does not count as passing.

# Workflow steps

**Toolchain (CI pins in `ci.yml` `env:`).** CI runs Go `1.27.3` (`setup-go` `go-version`, without `check-latest`) over example `go.mod` files that declare `go 1.26.x` (the floor). It installs govulncheck `v1.9.0`, gosec `v2.30.0` and golangci-lint `v2.15.0` (the scaffolded jobs only). The nested scaffolded workflows, which this repo never runs, use `go-version-file: go.mod` with `check-latest: true`. Scaffolded Dockerfiles build on `golang:1.27-alpine`. Locally, use `GOTOOLCHAIN=local` with Go 1.27.x and install the pinned scanner versions. If your local versions differ, say so in the report.

```bash
go install golang.org/x/vuln/cmd/govulncheck@v1.9.0
go install github.com/securego/gosec/v2/cmd/gosec@v2.30.0
```

**1. Per-example chain (`test-examples` job).** For each `*/go.mod`, in order, stopping at the first failure:

```bash
cd <example>
go build -o /dev/null ./...
go vet ./...
go test -race -count=1 ./...
go mod tidy -diff        # skipped in CI only on a manual workflow_dispatch bump
govulncheck ./...
gosec -quiet ./...
```

**2. Scaffolded matrix (`scaffolded-checks` job).** This runs for each directory with `.beluga/project.yaml`:

```bash
cd <example>
go install "github.com/lookatitude/beluga-ai/v2/cmd/beluga@$(go list -m -f '{{.Version}}' github.com/lookatitude/beluga-ai/v2)"
make check        # build + test + lint (golangci-lint v2.15.0)
make eval-ci      # mock-provider smoke eval
docker build .
```

`make check` builds a binary inside the example directory. A read-only lane (qa) runs this block in a temporary copy (`cp -r <example> "$TMPDIR"/ && cd "$TMPDIR"/<example>`) so no build output lands in the repo.

**3. Release bump (`bump-on-release` job, on `repository_dispatch: framework-release` only).** The job validates the payload tag, then for every example runs `go get github.com/lookatitude/beluga-ai/v2@<tag>` (retrying while the proxy lags), `go mod tidy`, build and vet. It rewrites the version stamps that track the pin:
- the `README.md` lines `beluga-ai/v2/cmd/beluga@<tag>` and ``Pins `github.com/lookatitude/beluga-ai/v2 <tag>``;
- `beluga-version:` in `.beluga/project.yaml`;
- the `org.opencontainers.image.version="<tag>"` label in `Dockerfile`.

It then opens `chore/bump-framework-<tag>` through `peter-evans/create-pull-request`, closes superseded bump PRs, and fails when any example is BROKEN. A manual bump must make the same stamp edits so the result matches what CI would produce.

# Evidence requirements

- Every `OK` needs the command that produced it. Never write "green" for a gate whose tool did not run. Use `NOT RUN` instead.
- Every `FAIL` quotes the first error line verbatim.

# Escalation rules

- A FAIL that needs framework intent rather than an API rename is BROKEN. Hand it to the coordinator, who may file a `lookatitude/beluga-ai` issue.
- If CI and this skill disagree, report the drift to the coordinator.

# Safety constraints

- Validate a framework tag against the pattern above before using it in a shell command.
- Do not hand-edit `.github/workflows/` or nested scaffolded workflows from an agent lane. The one exception: `/new-example` may generate a scaffolded example's nested workflows unchanged via `beluga init`.
- Do not push, open PRs or post reviews. The coordinator does that.

# Eval cases

See `evals.json` next to this file.
