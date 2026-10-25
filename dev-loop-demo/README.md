# dev-loop-demo

**What this is:** a scaffolded project (structurally identical to
[`scaffolded-basic`](../scaffolded-basic/)) used as the walk-through target for
`beluga dev` — the file-watch auto-restart loop introduced in framework `v2.12.0`.
The code is the canonical `beluga init` output; the README is the walk-through.

No hand-written Go. Everything here regenerates from `beluga init` (see
[Regenerate locally](#regenerate-locally) for the few documented differences).
The walk-through below is checked against the framework's `beluga dev` source at
the version pinned in `go.mod`.

## What `beluga dev` does

`beluga dev` is the inner-loop command for Beluga AI projects. It:

1. Builds the project and runs it once (same build + exec path as `beluga run`,
   so `.env` entries are loaded into the child's environment).
2. Watches the project tree for **`.go` file** writes, creates and renames, with a
   500 ms debounce so a burst of saves coalesces into one restart. Nothing else
   triggers a restart: non-Go files (including `.beluga/*.yaml` / `*.json`) and
   anything under a hidden directory (`.beluga/`, `.git/`, …), `vendor/` or
   `node_modules/` are ignored.
3. On an accepted change, sends `SIGTERM` to the running child, waits up to 3 s
   for it to exit (then `SIGKILL`), rebuilds, and starts a fresh child. A child
   that crashes is logged (`devloop: child exited: …`) and the watcher keeps
   running, so saving the fix recovers the session.
4. Handles `Ctrl+C` cleanly: the child gets the same `SIGTERM` → grace → `SIGKILL`
   sequence, then `beluga dev` exits.

The effect is the Go equivalent of `air` / `nodemon` — but aware of the Beluga
project layout and wired to the same persona/tool/provider surface your production
binary uses. The authoritative reference is the framework's
[`docs/reference/cli.md` § `beluga dev`](https://github.com/lookatitude/beluga-ai/blob/main/docs/reference/cli.md#beluga-dev).

## Walk-through — persona-edit loop

This is the canonical inner-loop demo. You edit the agent's persona, save, and see
the agent's next answer shift without ever restarting `beluga dev` yourself.

### 1. Start the dev loop

```bash
export OPENAI_API_KEY=sk-...
beluga dev
```

`beluga dev` prints `beluga dev: playground at http://127.0.0.1:8089` (the
playground is on by default — see below), builds the project, runs the agent
once, and then parks on the file watcher. The first run asks the agent to echo
`hello world`; you'll see the echo tool fire, followed by the agent's final
answer. Pass `--playground=off` if you only want the watcher.

### 2. Edit the persona

Open `main.go` and find the persona block inside `newAgent(...)`:

```go
agent.WithPersona(agent.Persona{
    Role:      "helpful assistant",
    Goal:      "answer questions accurately",
    Backstory: "You are concise and direct. When asked to echo something, always use the echo tool.",
    Traits:    []string{"concise", "accurate"},
}),
```

Change `Traits` to `[]string{"concise", "accurate", "cheerful"}` and save.

### 3. Observe the restart

About 500 ms after the save, `beluga dev` stops the previous child, rebuilds and
starts a new one. It does not print a "change detected" line of its own: what you
see is the new run's output — the trait change flows through the persona template
into the LLM's system prompt, so the tone shifts. No manual rebuild, no manual
restart. (With the playground enabled, each restart is also marked in its stderr
panel as `--- restart #N ---`.)

### 4. Iterate

Tighten the traits, add a new tool, swap the provider (see
[`scaffolded-basic/README.md`](../scaffolded-basic/README.md#swap-the-provider)),
change the prompt in `Invoke(...)`. Every save of a `.go` file triggers another
rebuild. `Ctrl+C` when done.

## The dev playground

The playground is **on by default**: `--playground` is a port flag (default
`8089`), not a boolean switch, so a bare `beluga dev --playground` fails with
`flag needs an argument`. Use:

```bash
beluga dev                     # playground on http://127.0.0.1:8089
beluga dev --playground=9000   # different port
beluga dev --playground=0      # ephemeral port (printed at startup)
beluga dev --playground=off    # watcher only, no dev UI
```

The server binds to `127.0.0.1` only and exports its URL to the child as
`BELUGA_PLAYGROUND_URL`. The page has three read-only panels — there is no
prompt input; the agent still runs the prompt hard-coded in `main.go`:

- **Traces** — a table of spans (name, model, input/output tokens, duration,
  status).
- **Tokens & cost** — input/output token and span counters, plus the current
  restart number.
- **stderr tail** — the last 100 batches of the child's stderr, with a
  `--- restart #N ---` marker at each restart.

> **Current limitation (verified against the framework source):** the CLI does
> not yet forward the child's OTel spans to the playground, so the **Traces** and
> **Tokens & cost** panels stay empty and only the stderr tail and the restart
> counter update. To see spans today, run with `BELUGA_OTEL_STDOUT=1` (pretty-
> printed to the child's stdout) or point `OTEL_EXPORTER_OTLP_ENDPOINT` at a
> collector.

## Try the smoke eval too

The project also ships the default `.beluga/eval.smoke.json` from the scaffolder.
`make eval-ci` runs it against the mock LLM provider — no API keys required — and
exits non-zero on any row failure:

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.10
make eval-ci
```

See [`eval-demo/README.md`](../eval-demo/README.md) for the full eval walk-through.

## What you *don't* need to do

- **Kill the previous process.** `beluga dev` does it.
- **Type `go run .` again.** `beluga dev` does it.
- **Add an entry to a config file.** The watcher is on by default; no config.

What you *do* still need to do: restart `beluga dev` yourself after editing
anything that isn't a `.go` file (for example `.beluga/eval.yaml` or `.env`) —
the watcher ignores those.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.10
beluga init dev-loop-demo \
  --module github.com/lookatitude/beluga-examples/dev-loop-demo
cd dev-loop-demo && go mod tidy
```

The regenerated tree matches this one apart from the `beluga init` timestamps,
this README, and the documented `ci.yml` / `Makefile` hand-edits shared with
scaffolded-basic (see
[`scaffolded-basic/README.md`](../scaffolded-basic/README.md#regenerate-locally)).

## Related examples

- [`scaffolded-basic`](../scaffolded-basic/) — the `beluga init` output this
  project is structurally copied from.
- [`eval-demo`](../eval-demo/) — same shape, walks through the four common
  `beluga eval` use cases.
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference.
