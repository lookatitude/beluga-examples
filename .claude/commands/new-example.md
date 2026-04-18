---
name: new-example
description: Create a new runnable example from a Linear sub-issue + brief + merged framework PR; open a PR with code-reviewer + security-reviewer pre-merge.
---

Scaffold example: $ARGUMENTS

$ARGUMENTS should be a Linear sub-issue ID (e.g., `LOO-48`) labeled `layer:examples`. For examples not tracked in Linear, pass a brief path instead: `/new-example ../research/briefs/<slug>.md`.

## Pre-flight (Linear-integrated, A1)

### 1. Fetch Linear context

If $ARGUMENTS matches `^LOO-\d+$` (case-insensitive):

```
mcp__claude_ai_Linear__get_issue(id="$ARGUMENTS")
```

Capture: sub-issue title, description, labels, parentId.

Then fetch the parent:

```
mcp__claude_ai_Linear__get_issue(id="<parentId>")
```

Capture: parent title, description, any brief path referenced in the description.

If $ARGUMENTS is a brief path, skip Linear and read the brief directly.

### 2. Read the workspace brief

From the parent's description, find the brief path at `../research/briefs/<slug>.md` (relative to the examples repo). Read it fully — it defines the feature being demonstrated.

### 3. Read the merged framework PR

The parent issue should have a linked merged PR in `lookatitude/beluga-ai`. Find its number. Then:

```bash
gh pr view <PR_NUMBER> --repo lookatitude/beluga-ai --json title,body,files,additions,deletions
gh pr diff <PR_NUMBER> --repo lookatitude/beluga-ai > /tmp/framework-pr-diff.patch
```

Examples reference the actual shipped API. The diff is the canonical reference for imports, function signatures, and idiomatic usage.

### 4. On Linear MCP failure

Retry twice with short backoff. On persistent failure, proceed with brief + framework PR only; flag the gap in the PR description. Do not silently skip.

## Scaffold

### 1. Derive branch name

From Linear issue labels:

- `Feature` label → `feat/`
- `Bug` label → `fix/`
- `Improvement` label → `refactor/`
- Otherwise → `chore/` (workspace `/plan-feature` deliberately does NOT create a `documentation` type label; doc-only example changes flow through untyped to `chore/`)

Branch: `<type>/loo-NN-<example-name>` where `<example-name>` is the kebab-case example directory name.

```bash
git checkout main
git pull
git checkout -b <branch-name>
```

### 2. Invoke example-author

`@agent-example-author` receives: sub-issue + parent + brief + framework PR summary/diff. Creates:

- `<example-name>/` directory
- `<example-name>/go.mod` importing `github.com/lookatitude/beluga-ai` pinned to the framework PR's release tag (or `@<commit-sha>` if no release tag yet)
- `<example-name>/main.go` — minimal, focused demonstration; one feature per example; full imports; explicit error handling; `os.Getenv` for secrets
- `<example-name>/README.md` — what it demonstrates, prerequisites, how to run, expected output

Per `.claude/rules/examples.md`:

- `go.mod` per example (not shared)
- Compiles with `go build .`
- Runs with `go run .`
- No `_` for error returns
- Registry pattern (`llm.New("provider", cfg)`), not direct construction
- Side-effect imports for providers

### 3. Invoke example-tester

`@agent-example-tester` verifies:

- `go build .` succeeds
- `go vet .` clean
- `go mod tidy` produces no diff
- If the example doesn't require API keys, `go run .` produces the documented expected output

### 4. Invoke pr-review-toolkit:code-reviewer (pre-merge gate, A1)

`@agent-pr-review-toolkit:code-reviewer` reviews the scaffolded example for: idiomatic Go, consistent style with other examples, complete imports, proper error handling, no hardcoded secrets, README accuracy. Address any high-confidence findings before proceeding.

### 5. Invoke example-security-reviewer (pre-merge gate, A2)

`@agent-example-security-reviewer` runs gosec + govulncheck + manual review for hardcoded secrets, unsafe patterns, and missing validation.

- If the reviewer posts `--request-changes`: address the blocking findings before opening the PR. Re-run the reviewer after fixes.
- If the reviewer posts `--comment` or `--approve`: proceed.

Blocking rules: any gosec high-severity finding, any govulncheck CVE, any hardcoded secret, any unbounded external read, user input → exec/SQL/file without validation.

### 6. Open the PR

```bash
git add <example-name>/
git commit -m "feat(example): <example-name> (LOO-NN)"
git push -u origin <branch-name>
gh pr create --fill --title "feat(example): <example-name> (LOO-NN)" --body "<PR body referencing Linear sub-issue, brief, and framework PR>"
```

Include `LOO-NN` in the PR title so Linear auto-links and auto-closes on merge.

### 7. Human review + merge

Linear auto-closes the sub-issue on merge. If the example references a tutorial (`layer:website` sibling sub-issue under the same parent), ensure the tutorial's link to the example directory is updated in the website PR for that sibling (or in a follow-up).

## Framework version pinning

When `go.mod` is generated, pin the framework version to the exact tag shipped by the merged framework PR (`github.com/lookatitude/beluga-ai vX.Y.Z`). If the PR is merged but no release tag exists yet, use `@<commit-sha>` until a release ships. Never use `@latest` — pin explicitly for reproducibility.
