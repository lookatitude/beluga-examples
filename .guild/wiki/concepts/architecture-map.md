---
type: concept
title: "Architecture map (stub)"
description: "Cheap-scan architecture stub for examples. Confidence-tagged candidate; the deep knowledge graph is deferred."
owner: orchestrator
confidence: medium
trust_tier: synthesized
status: candidate
importance: low
source_refs: [codebase-map]
created_at: 2026-10-25
updated_at: 2026-10-25
expires_at: null
supersedes: null
sensitivity: internal
generated_from_commit: eb667cac765e1525fc9a6f9b32ba207a1c8d2771
recall_importance: 3
---
# Architecture map (stub)

Built from the cheap scan (codebase map, a derived cache outside the repo; 48 files at
commit `eb667cac`; languages: dockerfile, go, json, makefile, markdown, yaml; frameworks: docker) plus this repo's
`CLAUDE.md`. Candidate page: every line carries a confidence tag. The deep knowledge
graph and onboarding tour were **not** built at init; run `/guild:learn` to build them.

## What this repo is

- **[high]** `beluga-examples`: standalone runnable Go programs demonstrating the Beluga AI framework. The demonstration layer of the Beluga project. Source: `CLAUDE.md` §What this is.
- **[high]** Every example is its own directory with its own `go.mod` importing `github.com/lookatitude/beluga-ai/v2`, builds with `go build .`, runs with `go run .`, and shares no code with other examples. Source: `CLAUDE.md` §Structure.

## Examples (from the scan)

- **[high]** `dev-loop-demo/main.go`
- **[high]** `eval-demo/main.go`
- **[high]** `first-agent/main.go`
- **[high]** `scaffolded-basic/main.go`

## Modules by size (from the scan)

| Path | Files |
|---|---|
| `dev-loop-demo` | 11 |
| `eval-demo` | 11 |
| `scaffolded-basic` | 11 |
| `.claude` | 8 |
| `first-agent` | 4 |
| `.github` | 1 |
| `.wiki` | 1 |
| `CLAUDE.md` | 1 |

## Flows

- **[high]** CI builds and vets every example. A framework release dispatches `framework-release` here; the bump job updates each `go.mod`. Source: `CLAUDE.md`, workspace release-flow page.

## Knowledge that stays where it is

- **[high]** This repo's wiki lives in `.guild/wiki/` (moved from `.wiki/` on 2026-10-25, structure kept).
