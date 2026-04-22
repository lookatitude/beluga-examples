# dev-loop-demo

**What this is:** a scaffolded project (structurally identical to
[`scaffolded-basic`](../scaffolded-basic/)) used as the walk-through target for
`beluga dev` — the file-watch auto-restart loop introduced in framework `v2.12.0`.
The code is the canonical `beluga init` output; the README is the walk-through.

No hand-written Go. Everything here regenerates verbatim from `beluga init`.

## What `beluga dev` does

`beluga dev` is the inner-loop command for Beluga AI projects. It:

1. Runs the agent once (same entrypoint as `beluga run` — your `main.go`).
2. Watches the project tree for changes (Go sources, `.beluga/` config, template
   files) with a 500 ms debounce so a burst of saves coalesces into one restart.
3. On change, sends `SIGTERM` to the running agent, waits for graceful exit, then
   re-invokes the build + run cycle.
4. Handles `Ctrl+C` cleanly: forwards the signal, waits for the child to drain,
   then exits.

The effect is the Go equivalent of `air` / `nodemon` — but aware of the Beluga
project layout and wired to the same persona/tool/provider surface your production
binary uses.

## Walk-through — persona-edit loop

This is the canonical inner-loop demo. You edit the agent's persona, save, and see
the agent's next answer shift without ever restarting `beluga dev` yourself.

### 1. Start the dev loop

```bash
export OPENAI_API_KEY=sk-...
beluga dev
```

`beluga dev` prints the resolved module path, runs the agent once, and then parks
on the file watcher. The first run asks the agent to echo `hello world`; you'll see
the echo tool fire, followed by the agent's final answer.

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

Within 500 ms, `beluga dev` logs:

```
detected change: main.go
restarting …
```

It rebuilds, relaunches the agent, and you see the next run's answer — the trait
change flows through the persona template into the LLM's system prompt, so the
tone shifts. No manual rebuild, no manual restart.

### 4. Iterate

Tighten the traits, add a new tool, swap the provider (see
[`scaffolded-basic/README.md`](../scaffolded-basic/README.md#swap-the-provider)),
change the prompt in `Invoke(...)`. Every save triggers another rebuild. `Ctrl+C`
when done.

## `beluga dev --playground`

`--playground` opens the three-panel UI on top of the same file-watch loop:

- **Left panel — Prompt.** Editable prompt box. Send prompts to the running agent
  interactively instead of re-running `Invoke(...)` with the hard-coded string.
- **Middle panel — Transcript.** Running log of agent responses, tool calls, and
  tool results. Each turn has a timestamp and a span-ID breadcrumb.
- **Right panel — Traces.** OTel span tree for the most recent turn. Drill into
  `llm.generate`, `tool.call`, and `agent.invoke` spans. `BELUGA_OTEL_STDOUT=1`
  still mirrors to stdout so you can pipe through other tools.

```bash
beluga dev --playground
```

Edit persona → save → prompt the new persona from the left panel → see the span
tree in the right panel reflect the updated trait set. The inner-loop UI makes
the LLM/tool/OTel feedback tight enough that you shape the agent by feel.

## Try the smoke eval too

The project also ships the default `.beluga/eval.smoke.json` from the scaffolder.
`make eval-ci` runs it against the mock LLM provider — no API keys required — and
exits non-zero on any row failure:

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.0
make eval-ci
```

See [`eval-demo/README.md`](../eval-demo/README.md) for the full eval walk-through.

## What you *don't* need to do

- **Kill the previous process.** `beluga dev` does it.
- **Type `go run .` again.** `beluga dev` does it.
- **Restart for `.beluga/` edits.** The watcher covers them too.
- **Add an entry to a config file.** The watcher is on by default; no config.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.0
beluga init dev-loop-demo \
  --module github.com/lookatitude/beluga-examples/dev-loop-demo
```

The only hand-authored files relative to the scaffolded output are this README and
the `name:` field in `.beluga/project.yaml` (which matches the directory name).

## Related examples

- [`scaffolded-basic`](../scaffolded-basic/) — the verbatim `beluga init` output
  this project is structurally copied from.
- [`eval-demo`](../eval-demo/) — same shape, walks through the four common
  `beluga eval` use cases.
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference.
