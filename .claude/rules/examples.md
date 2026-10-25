---
description: Rules for all example code in this repo.
globs: "**/*"
alwaysApply: true
---

# Examples Rules

## Every example MUST

- Have its own `go.mod` — no workspace-level go.mod, no shared modules.
- Compile with `go build .` from its own directory.
- Have a file-top comment block documenting: what it demonstrates, prerequisites (Go version, env vars), and how to run it.
- Use full imports (`github.com/lookatitude/beluga-ai/v2/...`).
- Handle every error explicitly — no `_` for error returns.
- Use the registry pattern (`llm.New("anthropic", cfg)`) not direct construction.
- Import providers via side-effect: `import _ "github.com/lookatitude/beluga-ai/v2/llm/providers/anthropic"`.
- Use `os.Getenv()` for API keys — never hardcoded secrets.

## Every example MUST NOT

- Import from other example directories.
- Share helper packages or utility code across examples.
- Commit API keys, tokens, or credentials.
- Be committed directly to `main` — branch discipline applies.

## Directory naming

- kebab-case: `first-agent`, `custom-provider`, `scatter-gather`
- Match the framework guide slug when a companion guide exists
- One feature per directory; use `<feature>-basic` and `<feature>-advanced` for graduated examples

## Toolchain

- `go.mod` floor: `go 1.26` (current examples declare 1.26.x). Root CI builds with Go **1.27.3**, using the `GO_VERSION` pin in `.github/workflows/ci.yml`.
- Scanners pinned in CI: govulncheck **v1.9.0**, gosec **v2.30.0**, and golangci-lint **v2.15.0** (scaffolded `make check` only).
- Scaffolded Dockerfiles build on `golang:1.27-alpine` and come from `beluga init`. Do not hand-tune them.
- Full gate list and the release-bump stamps: `.guild/skills/examples-ci-gates/SKILL.md`. The longer authoring checklist: `.guild/skills/examples-conventions/SKILL.md`.

## Before commit

Run what root CI runs (CI parity, with Go 1.27.x and `GOTOOLCHAIN=local`):

```bash
cd <example-dir>
go build -o /dev/null ./...
go vet ./...
go test -race -count=1 ./...
go mod tidy -diff     # non-zero exit + diff if go.mod/go.sum are not tidy
govulncheck ./...
gosec -quiet ./...
# scaffolded examples (.beluga/project.yaml) also: make check && make eval-ci && docker build .
```

If a tool is missing locally, say so in the commit or PR. Never claim green for a gate that did not run.

## Who commits and pushes

Agents and teammates commit on their branch only. The main session (coordinator) pushes, opens PRs, posts reviews, and writes to Linear. Codex (`/codex:adversarial-review`) is the general adversarial reviewer before a PR opens.
