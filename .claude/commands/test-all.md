---
name: test-all
description: Validate that every example in the repo compiles and passes go vet.
---

Test all examples.

## Workflow

### Step 1 — Discover
Find all directories with a `go.mod` file.

### Step 2 — Test
`@agent-example-tester` runs build + vet + tidy checks on each example.

### Step 3 — Report
Print results per example (OK / FAIL with error), plus a summary line: `N/M examples passing`.

If any example fails, exit with a non-zero status so CI can gate on this command.
