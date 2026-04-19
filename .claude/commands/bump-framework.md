---
name: bump-framework
description: Bump all examples to a new framework version, fix breakage, and open a PR.
---

Bump framework to: $ARGUMENTS

$ARGUMENTS should be a framework release tag (e.g., `v2.8.5`). If omitted, bumps to the latest release.

## Workflow

### Step 1 — Resolve tag
If no tag provided, look up the latest `lookatitude/beluga-ai` release:
```bash
gh release view --repo lookatitude/beluga-ai --json tagName --jq .tagName
```

### Step 2 — Bump and test
`@agent-example-maintainer` runs:
1. Bumps all examples to the resolved tag (`go get github.com/lookatitude/beluga-ai/v2@<tag>`)
2. Runs `go build .` and `go vet .` on each
3. Fixes breakage where possible (minimal fixes — API renames, new required parameters)
4. Reports results: OK / FIXED / BROKEN per example

### Step 3 — Output
Open a single PR with all bumps and fixes. PR title: `chore(deps): bump beluga-ai to <tag>`. PR body includes the per-example report.

If any examples are BROKEN (unable to fix automatically), list them in the PR body and file issues in `lookatitude/beluga-ai` asking for migration guidance.
