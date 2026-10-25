---
schema_version: guild.init_record.v1
slug: examples
run_id: run-20261025-043828-init
mode: brownfield
repo_kind: regular
created_at: 2026-10-25
generated_from_commit: eb667cac765e1525fc9a6f9b32ba207a1c8d2771
guild_version: 3.1.0
storage_layout_version: 2
---
# Init record: examples

## Detection
- `workspace-detect`: **regular** (no immediate child with `.git/` or `.guild/`).
  Registered as a sub-guild of the Beluga workspace (`../.guild/workspace.json`).

## Scan (cheap tier = Init-DONE)
- `learn-scan` → codebase map (derived cache under `~/.cache/guild/roots/`, not committed; rebuild with `/guild:learn map`): 48 files; languages dockerfile, go, json, makefile, markdown, yaml.
- `wiki/concepts/architecture-map.md`: confidence-tagged stub, status candidate.
  Grounded in the scan plus `CLAUDE.md`.
- **Deferred, not produced:** `knowledge-graph.json`, `onboarding-tour.md`. Run
  `/guild:learn` to build them.

## Existing knowledge
- The legacy `.wiki/` is **unchanged** and still used by this repo's agents, hooks and
  commands. It was not migrated at init; that is a separate, opt-in change.

## Reviews
- G-init cross-host review: **skipped** by policy (`review: local`).
