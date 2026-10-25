---
name: new-example
description: Create a new runnable example from a Linear sub-issue + brief + merged framework PR; Codex adversarial review + example-security-reviewer gate it, then the coordinator opens the PR.
---

Scaffold example: $ARGUMENTS

$ARGUMENTS should be a Linear sub-issue ID (e.g., `LOO-48`) labeled `layer:examples`. For examples not tracked in Linear, pass a brief path instead: `/new-example ../research/briefs/<slug>.md`.

## Invocation notes

- **Who does what:** the main session running this command is the **coordinator**. It reads Linear, creates the branch, dispatches the agents (as teammates or inline agents), runs the reviews, and does every outward action: push, PR, Linear. The dispatched agents only prepare local commits and reports.
- **Conventions:** the agents follow `.claude/rules/examples.md` and the Guild skills `.guild/skills/examples-conventions/SKILL.md` and `.guild/skills/examples-ci-gates/SKILL.md`.

- **Scope:** One competency per example directory—enumerate required env vars and failure modes explicitly in the header comment.
- **Copy:** Comments and README excerpts should stay defensive/teaching; avoid sensational security wording.

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
gh pr diff <PR_NUMBER> --repo lookatitude/beluga-ai > <scratchpad>/framework-pr-diff.patch
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

`@agent-example-author` (Guild: `example-author`) receives: sub-issue + parent + brief + framework PR summary/diff. Creates:

- `<example-name>/` directory
- `<example-name>/go.mod` importing `github.com/lookatitude/beluga-ai/v2` pinned to the framework PR's release tag (or `@<commit-sha>` if no release tag yet)
- `<example-name>/main.go` — minimal, focused demonstration; one feature per example; full imports; explicit error handling; `os.Getenv` for secrets
- `<example-name>/README.md` — what it demonstrates, prerequisites, how to run, expected output
- For a scaffolded example: the `beluga init` output at the pinned version (Dockerfile `golang:1.27-alpine` builder, Makefile, nested CI) with only documented hand-edits

The author commits locally on the branch and hands back; it does not push.

Per `.claude/rules/examples.md`:

- `go.mod` per example (not shared)
- Compiles with `go build .`
- Runs with `go run .`
- No `_` for error returns
- Registry pattern (`llm.New("provider", cfg)`), not direct construction
- Side-effect imports for providers

### 3. Invoke example-tester

`@agent-example-tester` (Guild: `qa`) runs the root-CI gate chain from `examples-ci-gates` on the new directory:

- `go build -o /dev/null ./...`, `go vet ./...`, `go test -race -count=1 ./...`
- `go mod tidy -diff` reports no diff
- `govulncheck ./...` (v1.9.0) and `gosec -quiet ./...` (v2.30.0)
- Scaffolded examples: `make check`, `make eval-ci`, `docker build .`
- If the example doesn't require API keys, `go run .` produces the documented expected output

Any `NOT RUN` gate is stated in the PR body. It is never reported as passing.

### 4. General code review (pre-merge gate, A1)

Review the example for: idiomatic Go, consistent style with other examples, complete imports, proper error handling, no hardcoded secrets, README accuracy. Address any high-confidence findings before proceeding. Use the first reviewer that is available, and say in the PR body which one ran:

1. **Codex adversarial review** — `/codex:adversarial-review` on the branch (needs the `codex@openai-codex` plugin and a logged-in `codex` CLI; see the workspace `AGENTS.md` § Adversarial review). This is the default reviewer.
2. **`@agent-pr-review-toolkit:code-reviewer`** — only if the `pr-review-toolkit` plugin is installed (`claude plugin list`). It is not installed by default.
3. **Neither available** — do a self-review against the checklist above and record "general review: self-review only (no Codex / pr-review-toolkit)" in the PR body. Do not describe it as an independent review.

### 5. Invoke example-security-reviewer (pre-merge gate, A2)

`@agent-example-security-reviewer` (Guild: `security`) runs gosec + govulncheck + manual review for hardcoded secrets, unsafe patterns, and missing validation. It writes a report and returns a verdict. It never posts to GitHub.

- Verdict `blocking`: send the findings back to `example-author`, then re-run the reviewer after the fixes. Do not open the PR while the verdict is blocking.
- Verdict `advisory` or `pass`: proceed. The coordinator pastes the report summary into the PR body.

Blocking rules: any gosec high-severity finding, any govulncheck CVE, any hardcoded secret, any unbounded external read, user input → exec/SQL/file without validation.

### 6. Open the PR (coordinator only)

Only the main session runs this step. A teammate that reaches it stops after its local commit and reports back.

```bash
git branch --show-current   # must not be main
git status --short          # the author's commit(s) should already be on the branch
git push -u origin <branch-name>
gh pr create --title "feat(example): <example-name> (LOO-NN)" --body "<PR body: Linear sub-issue, brief, framework PR, gate report, which general reviewer ran, security verdict>"
```

If anything is still uncommitted, commit it first as `feat(example): <example-name> (LOO-NN)`. Never use `--no-verify`.

Include `LOO-NN` in the PR title so Linear auto-links and auto-closes on merge.

### 7. Human review + merge

Linear auto-closes the sub-issue on merge. If the example references a tutorial (`layer:website` sibling sub-issue under the same parent), ensure the tutorial's link to the example directory is updated in the website PR for that sibling (or in a follow-up).

## Framework version pinning

When `go.mod` is generated, pin the framework version to the exact tag shipped by the merged framework PR (`github.com/lookatitude/beluga-ai/v2 vX.Y.Z`). If the PR is merged but no release tag exists yet, use `@<commit-sha>` until a release ships. Never use `@latest` — pin explicitly for reproducibility.
