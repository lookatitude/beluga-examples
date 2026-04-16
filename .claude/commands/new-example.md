---
name: new-example
description: Create a new runnable example for a framework feature, with its own go.mod and self-contained main.go.
---

Create example for: $ARGUMENTS

## Workflow

### Step 1 — Source material
Identify the framework feature being demonstrated. Read:
- The relevant framework docs (`../framework/docs/guides/`, `../framework/docs/patterns/`)
- The package source (`../framework/<package>/`)
- Any existing website tutorial this example should companion (`../website/src/content/docs/docs/guides/`)

Check existing examples to ensure no overlap.

### Step 2 — Create
`@agent-example-author` creates the example:
- Directory: `<example-name>/` (kebab-case, matching the framework guide slug if one exists)
- `go.mod` importing `github.com/lookatitude/beluga-ai@latest`
- `main.go` with file-top comment block, full imports, explicit error handling

### Step 3 — Verify
```bash
cd <example-name>
go build .
go vet .
go mod tidy
```
All three must succeed.

### Step 4 — Output
Open a PR against `main` per branch discipline. Link to the companion website tutorial (or note that one should be created via `/tutorial` in the website repo).
