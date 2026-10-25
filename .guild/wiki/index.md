---
title: "Beluga Examples Wiki"
description: "Wiki router and catalog for this repo (legacy .wiki index, migrated)."
recall_importance: 2
---
# Beluga Examples Wiki

> Migrated from `.wiki/index.md` on 2026-10-25 (Guild 3.x; structure kept). Verbatim original: `.guild/knowledge/sources/wiki-index/original.md` (sha256 `36dd7c5bb1d8…`).

**Scope: examples-local.** This wiki holds knowledge specific to the runnable examples: per-example conventions, framework version compatibility notes, testing patterns, and known gotchas. It does **not** hold framework patterns (those are in `../../../framework/.guild/wiki/`) or cross-repo coordination material (that is in the workspace wiki at `./`).

## What lives here (and when)

| Topic | Status |
|---|---|
| Per-example compatibility notes (which framework version each example targets) | *on-demand* |
| Testing patterns (how CI validates examples) | *pending — when CI is added* |
| Known gotchas (API keys, rate limits, provider-specific quirks in examples) | *on-demand* |
| Example-to-tutorial linking conventions | *pending — Phase 5+* |

## Where everything else lives

- **Framework Go patterns, APIs, architecture** → `../../../framework/.guild/wiki/`
- **Release flow, dispatch chains, cross-repo incidents** → workspace `./`
- **Website content conventions, design tokens** → `../../../website/.guild/wiki/`

## Guild pages

- [Architecture map (stub)](concepts/architecture-map.md): cheap-scan, confidence-tagged map of this repo.
- [Glossary](glossary.md): Guild operating terms, plus project terms.
