# first-agent

**Historical reference.** This was the original hand-written Beluga AI example —
a calculator-tool agent that computes `17 × 42 − 19`. It predates the
`beluga init` scaffolder introduced in framework `v2.11.0` and the eval /
dev-loop surface added in `v2.12.0`–`v2.13.0`.

**New projects should start from [`scaffolded-basic`](../scaffolded-basic/)
instead.** That example is the verbatim `beluga init` output on `v2.13.0` —
same agent shape, plus the `.beluga/` config directory, Makefile with
`eval-ci` / `test-ci` targets, ci.yml, Dockerfile, and the `BELUGA_ENV=eval`
branch in `main.go` that makes `beluga eval` work out of the box.

This example is kept in-tree so the pre-scaffolder narrative (blog posts,
tutorial screenshots, prior changelogs) continues to resolve. It still builds
and runs; it's just no longer the recommended starting point.

## Run it

```bash
export OPENAI_API_KEY=sk-...
go run .
```

## Related examples

- [`scaffolded-basic`](../scaffolded-basic/) — the canonical `beluga init` output.
- [`dev-loop-demo`](../dev-loop-demo/) — `beluga dev` walk-through.
- [`eval-demo`](../eval-demo/) — `beluga eval` walk-through with 4-row dataset.
