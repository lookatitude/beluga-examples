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
- Use full imports (`github.com/lookatitude/beluga-ai/...`).
- Handle every error explicitly — no `_` for error returns.
- Use the registry pattern (`llm.New("anthropic", cfg)`) not direct construction.
- Import providers via side-effect: `import _ "github.com/lookatitude/beluga-ai/llm/providers/anthropic"`.
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

## Before commit

```bash
cd <example-dir>
go build .
go vet .
go mod tidy && git diff --exit-code go.mod go.sum
```
