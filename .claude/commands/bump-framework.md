---
name: bump-framework
description: Bump all examples to a new framework version (go.mod plus the version stamps CI tracks), fix breakage, and hand the branch to the coordinator for the PR.
---

Bump framework to: $ARGUMENTS

$ARGUMENTS should be a framework release tag (e.g. `v2.8.5`). If it is omitted, bump to the latest release.

**Normally automatic.** On every framework release, root CI's `bump-on-release` job (`.github/workflows/ci.yml`) bumps every example, rewrites the version stamps, and opens `chore/bump-framework-<tag>`. Use this command when that PR is red (BROKEN examples), or to try a specific tag by hand. Bump steps, stamps and gates: `.guild/skills/examples-ci-gates/SKILL.md`.

## Workflow

### Step 1: Resolve and validate the tag
If no tag is given, look up the latest `lookatitude/beluga-ai` release:
```bash
gh release view --repo lookatitude/beluga-ai --json tagName --jq .tagName
```
Reject anything that does not match `^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$`.

### Step 2: Branch
If CI already opened `chore/bump-framework-<tag>`, check it out. Otherwise create it from an up-to-date `main`. Never commit to `main`.

### Step 3: Bump, test, fix
`@agent-example-maintainer` (Guild: `example-maintainer`):
1. Runs `go get github.com/lookatitude/beluga-ai/v2@<tag>` and `go mod tidy` for every example, then updates the version stamps (README `cmd/beluga@` and `Pins` lines, `.beluga/project.yaml` `beluga-version:`, the Dockerfile `org.opencontainers.image.version` label).
2. Runs the CI-parity gate chain on each example (Go 1.27.x, `GOTOOLCHAIN=local`, govulncheck v1.9.0, gosec v2.30.0; scaffolded examples also run `make check`, `make eval-ci` and `docker build .`).
3. Applies minimal fixes only (API renames, new required parameters).
4. Commits locally and reports OK / FIXED / BROKEN per example.

### Step 4: Review and PR (coordinator only)
The main session, never the maintainer:
1. Runs `/codex:adversarial-review` on the branch (see the workspace `AGENTS.md` § Adversarial review). If Codex is unavailable, records "general review: self-review only" in the PR body.
2. Pushes and opens the PR, or updates the CI one. The title is `chore(deps): bump beluga-ai to <tag>` and the body carries the per-example report.
3. For BROKEN examples (unable to fix automatically), lists them in the PR body and files `lookatitude/beluga-ai` issues asking for migration guidance, using the issue text the maintainer drafted.
