# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Standalone runnable code examples for the Beluga AI framework (`github.com/lookatitude/beluga-ai/v2`). Each example is a self-contained Go program that demonstrates a specific feature or pattern and can be built and run independently. Examples are meant to be linked from the framework guides and the website's tutorials; as of 2026-10-25 no guide or website page links to a specific example yet (only release posts link to the repo), so add the link on the docs side when you ship or rename one.

## Project scope

This repo owns the **demonstration layer** of the Beluga project: runnable, copy-paste-ready programs that prove the framework works. It does **not** own framework code (that's `beluga-ai`), docs content (that's `beluga-website`), or project strategy (that's `beluga.git` workspace).

When working on cross-repo concerns (a framework API change that breaks examples, a new tutorial that needs a companion example), consult the **workspace wiki** at `../.guild/wiki/context/workspace-router.md` and the workspace `CLAUDE.md` at `../CLAUDE.md` — they're only visible when this repo is checked out inside the multi-repo Beluga workspace.

## Structure

Each example lives in its own directory with its own `go.mod`:

```
examples/
├── CLAUDE.md
├── .guild/wiki/
│   └── index.md
├── .claude/
│   ├── agents/
│   ├── commands/
│   └── rules/
├── .github/workflows/ci.yml   # root CI: the only workflow GitHub runs
├── first-agent/               # hand-written example (agent + typed tool)
│   ├── go.mod / go.sum
│   ├── main.go / main_test.go
│   └── README.md
├── scaffolded-basic/          # `beluga init` output (see its README)
│   ├── .beluga/               # project.yaml, eval.yaml, eval.smoke.json
│   ├── .github/workflows/ci.yml  # shipped for copy-paste; NOT run by this repo
│   ├── Dockerfile / Makefile / .env.example / .gitignore
│   ├── go.mod / go.sum / main.go
│   └── README.md
├── dev-loop-demo/             # scaffolded shape + `beluga dev` walk-through
└── eval-demo/                 # scaffolded shape + 4-row dataset + `beluga eval` walk-through
```

Scaffolded examples (those with `.beluga/project.yaml`) must stay identical to `beluga init` output at the pinned framework version apart from timestamps and the hand-edits listed in `scaffolded-basic/README.md`. Their nested `.github/workflows/ci.yml` never runs here; root CI runs the equivalent `make check`, `make eval-ci` and `docker build` for each of them.

Each example directory MUST:
- Have its own `go.mod` importing `github.com/lookatitude/beluga-ai/v2`
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

# Bump all examples to a new framework version (<tag> = e.g. the release tag)
# Normally done by the CI bump-on-release job, which also updates the README
# `cmd/beluga@<tag>` lines, .beluga/project.yaml and the Dockerfile label.
for d in */; do [ -f "$d/go.mod" ] && (cd "$d" && go get github.com/lookatitude/beluga-ai/v2@<tag> && go mod tidy); done

# Scaffolded example gates (what root CI runs per scaffolded example)
cd scaffolded-basic && \
  go install "github.com/lookatitude/beluga-ai/v2/cmd/beluga@$(go list -m -f '{{.Version}}' github.com/lookatitude/beluga-ai/v2)" && \
  make check && make eval-ci && docker build .
```

## Naming conventions

- Directory name = example name = slug used in website tutorial links
- Use kebab-case: `first-agent`, `dev-loop-demo`, `multi-agent-team`
- Match the framework guide name where a companion guide exists (e.g., `first-agent` matches `docs/guides/first-agent.md`; the companions of `dev-loop-demo` and `eval-demo` are `docs/guides/dev-loop.md` and `docs/guides/evaluation.md`). Those guides don't link back to the examples yet.

## Code style

- Follow the framework's Go conventions: `context.Context` first parameter, explicit error handling, `slog` for logging.
- Use the registry pattern (`llm.New("anthropic", cfg)`) not direct construction.
- Import providers for side-effect registration: `import _ "github.com/lookatitude/beluga-ai/v2/llm/providers/anthropic"`.
- Keep examples minimal — demonstrate one feature, not all features.
- Add a file-top comment block documenting what the example demonstrates and what env vars are needed.

## Related repos

- `beluga-ai` (framework) — the library these examples import. API changes there may break examples here.
- `beluga-website` — tutorials and guides link to examples here. When renaming example directories, update the website links too.
- `beluga.git` (workspace) — shared knowledge, roadmap, coordination.

## Prompting baseline and example agents

Agents under `.claude/agents/` include a **Prompting baseline** (current Claude 5 family models). For this repo, that means **literal** scope per example (one feature per example), **explicit** file-top doc comments and env var documentation, and **no assumed** handoff to the website team—link and label what the example proves. When an example’s `main.go` uses the user’s real API keys, keep comments focused on **defensive, local use**; avoid celebratory offensive-security phrasing in sample text.

## Wiki

`.guild/wiki/` holds examples-specific knowledge (currently a glossary, log and architecture map; per-example conventions and testing patterns live in this file and `.claude/rules/examples.md`). See `.guild/wiki/index.md` for routing.

## Security review

Every new or modified example triggers `example-security-reviewer` (`.claude/agents/example-security-reviewer.md`) via the `/new-example` command. The reviewer runs gosec + govulncheck + manual checks for hardcoded secrets, unsafe patterns, and missing validation.

Root CI runs `govulncheck ./...` and `gosec -quiet ./...` on every example, so the CVE and gosec rules below are enforced on every PR, not only via `/new-example`.

Blocking findings (any gosec high-severity, any govulncheck CVE, hardcoded secrets, unbounded external reads, user input → exec/SQL/file without validation) must be fixed before merge.

Non-blocking advisory findings (style, docs, low-severity) are surfaced in the review but don't block.

Manual invocation: `@agent-example-security-reviewer` on any example diff.
