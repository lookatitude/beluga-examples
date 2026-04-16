# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Standalone runnable code examples for the Beluga AI framework (`github.com/lookatitude/beluga-ai`). Each example is a self-contained Go program that demonstrates a specific feature or pattern and can be built and run independently. Examples are linked from the website's tutorials and guides.

## Project scope

This repo owns the **demonstration layer** of the Beluga project: runnable, copy-paste-ready programs that prove the framework works. It does **not** own framework code (that's `beluga-ai`), docs content (that's `beluga-website`), or project strategy (that's `beluga.git` workspace).

When working on cross-repo concerns (a framework API change that breaks examples, a new tutorial that needs a companion example), consult the **workspace wiki** at `../.wiki/index.md` and the workspace `CLAUDE.md` at `../CLAUDE.md` — they're only visible when this repo is checked out inside the multi-repo Beluga workspace.

## Structure

Each example lives in its own directory with its own `go.mod`:

```
examples/
├── CLAUDE.md
├── .wiki/
│   └── index.md
├── .claude/
│   ├── agents/
│   ├── commands/
│   └── rules/
├── first-agent/
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── custom-provider/
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── ...
```

Each example directory MUST:
- Have its own `go.mod` importing `github.com/lookatitude/beluga-ai`
- Compile with `go build .`
- Run with `go run .` (may require API keys via env vars — document in a comment at the top of `main.go`)
- Be self-contained — no shared libraries or helper packages across examples
- Handle errors explicitly (no `_` for error returns)

## Common commands

```bash
# Build a single example
cd first-agent && go build .

# Run a single example (may need API keys)
cd first-agent && go run .

# Test all examples compile
for d in */; do [ -f "$d/go.mod" ] && (cd "$d" && go build . && echo "OK: $d" || echo "FAIL: $d"); done

# Bump all examples to a new framework version
for d in */; do [ -f "$d/go.mod" ] && (cd "$d" && go get github.com/lookatitude/beluga-ai@v2.8.5 && go mod tidy); done
```

## Naming conventions

- Directory name = example name = slug used in website tutorial links
- Use kebab-case: `first-agent`, `custom-provider`, `multi-agent-team`
- Match the framework guide name where a companion guide exists (e.g., `first-agent` matches `docs/guides/first-agent.md`)

## Code style

- Follow the framework's Go conventions: `context.Context` first parameter, explicit error handling, `slog` for logging.
- Use the registry pattern (`llm.New("anthropic", cfg)`) not direct construction.
- Import providers for side-effect registration: `import _ "github.com/lookatitude/beluga-ai/llm/providers/anthropic"`.
- Keep examples minimal — demonstrate one feature, not all features.
- Add a file-top comment block documenting what the example demonstrates and what env vars are needed.

## Related repos

- `beluga-ai` (framework) — the library these examples import. API changes there may break examples here.
- `beluga-website` — tutorials and guides link to examples here. When renaming example directories, update the website links too.
- `beluga.git` (workspace) — shared knowledge, roadmap, coordination.

## Wiki

`.wiki/` holds examples-specific knowledge: per-example conventions, known framework version compatibility, testing patterns. See `.wiki/index.md` for routing.
