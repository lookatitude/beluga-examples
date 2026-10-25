---
derived_from_template: guild.skill_template.v1
name: examples-conventions
description: "The beluga-examples authoring contract: one self-contained Go module per example directory (own go.mod, no cross-example imports or shared helpers), kebab-case naming that matches the framework guide slug, the mandatory file-top comment block, explicit error handling, registry-pattern construction with side-effect provider imports, os.Getenv-only secrets, exact framework version pins, and the beluga-init-identical rule for scaffolded examples. Pulled by example-author, example-maintainer, security and qa. TRIGGER for \"write a new example\", \"does this example follow the repo conventions\", \"what should the example directory be called\", \"what goes in the file-top comment\", \"can two examples share a helper package\". DO NOT TRIGGER for: what root CI runs, toolchain/scanner versions or local CI-parity commands (examples-ci-gates); gosec/govulncheck triage and secret-scan verdicts (security); framework API design or framework code (beluga-ai repo); website tutorial content (beluga-website repo)."
when_to_use: "Any lane that creates, renames, edits or reviews an example directory in beluga-examples and needs the per-example rules before writing or judging code."
type: specialist
---

# examples-conventions

Binding sources, in precedence order: `.claude/rules/examples.md` (always loaded by Claude Code), then `AGENTS.md` §Structure / §Naming conventions / §Code style. This skill is the checklist form of both; when they disagree, the rules file wins and this skill is the one to fix.

# When to use it

- Creating a new example directory (with `/new-example` or ad hoc).
- Fixing an example after a framework bump without changing its scope.
- Reviewing an example diff for convention drift.
- Choosing or changing an example's directory name.

# When not to use it

- To decide which CI gates run, or which Go, gosec or golangci-lint version to use: load `examples-ci-gates`.
- To triage scanner findings: the `security` specialist owns verdicts.
- For framework or website changes: those belong to the other repos.

# Required inputs

- The example directory (existing or planned) and the feature it demonstrates.
- The framework release tag it pins (from the merged framework PR or the release). Use a commit SHA only when no tag exists yet.
- For scaffolded examples (those with `.beluga/project.yaml`): the `beluga init` output at the pinned version, and the hand-edit list in `scaffolded-basic/README.md`.

# Output format

A pass/fail checklist over the rules below, one line per rule (`OK` or `VIOLATION: <file>:<line> — <rule>`), or code that satisfies every rule.

# Workflow steps

1. **Layout.** `<example>/go.mod` (module `github.com/lookatitude/beluga-examples/<example>`), `main.go`, `README.md`, optional `main_test.go`. No root `go.mod`, no `go.work`, no imports from sibling example directories, no shared helper packages. Duplicate a few lines instead.
2. **Naming.** The directory name is the example name and the slug used in links. Use kebab-case. Match the framework guide slug when a companion guide exists (`first-agent` ↔ `framework/docs/guides/first-agent.md`; `dev-loop-demo` ↔ `dev-loop.md`; `eval-demo` ↔ `evaluation.md`). Graduated pairs are `<feature>-basic` / `<feature>-advanced`. Each directory demonstrates one feature.
3. **Version pin.** `go get github.com/lookatitude/beluga-ai/v2@<tag>`, with the exact tag. Never use `@latest`. The `go` line is 1.26 or later (the repo floor; current examples use 1.26.x).
4. **File-top comment** in `main.go`:
   ```go
   // Example: <example-name>
   //
   // Demonstrates: <one sentence>
   //
   // Prerequisites:
   //   - Go 1.26+ (see go.mod)
   //   - <ENV_VAR>=<what it is>
   //
   // Run:
   //
   //	go run .
   ```
5. **Code style.** `context.Context` comes first. Handle every error, with no `_` for error returns. Log with `slog`. Construct through the registry (`llm.New("anthropic", cfg)`) with side-effect provider imports (`import _ "github.com/lookatitude/beluga-ai/v2/llm/providers/anthropic"`). Use full `github.com/lookatitude/beluga-ai/v2/...` imports. Write the shortest program that demonstrates the feature correctly.
6. **Secrets.** Read API keys only through `os.Getenv`. Never hardcode a key, call `os.Setenv` with a literal key, or log a config struct that holds one. Keep the comment and README wording defensive and focused on local use.
7. **Scaffolded examples.** Keep them identical to `beluga init` output at the pinned version, apart from timestamps and the hand-edits marked `[beluga-examples hand-edit: …]` and listed in `scaffolded-basic/README.md`. Their `Dockerfile` (a `golang:1.27-alpine` builder on a distroless runtime), `Makefile` and nested `.github/workflows/ci.yml` come from the generator. Do not hand-tune them.
8. **Cross-repo links.** When you add or rename a directory, record a followup for the website and framework guide links. The examples layer does not edit those repos.

# Evidence requirements

Cite a file and line for every violation. A claim that the example compiles or is tidy needs the command output from `examples-ci-gates`, not this checklist.

# Escalation rules

- If a feature cannot be shown without a second package or a shared helper, stop and raise it with the coordinator. Do not create one.
- If a scaffolded example needs a change that `beluga init` does not produce, propose a documented hand-edit and do not apply it silently.

# Safety constraints

- Never commit credentials. `.env` stays gitignored; ship a `.env.example` only.
- Never weaken error handling to shorten an example.

# Eval cases

See `evals.json` next to this file.
