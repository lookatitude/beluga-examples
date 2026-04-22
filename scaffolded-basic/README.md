# scaffolded-basic

**What this is:** the verbatim output of `beluga init scaffolded-basic` on framework
`v2.13.0`. Zero hand-edits. Run it as-is, or treat it as the canonical starting shape
when you bootstrap a new Beluga AI project.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.0
beluga init scaffolded-basic \
  --module github.com/lookatitude/beluga-examples/scaffolded-basic
```

The output is deterministic — the same command produces a byte-identical tree.

## What you get

| Path                         | Purpose                                                           |
|------------------------------|-------------------------------------------------------------------|
| `main.go`                    | Agent + echo tool + OTel bootstrap + `BELUGA_ENV=eval` branch     |
| `go.mod`                     | Pins `github.com/lookatitude/beluga-ai/v2 v2.13.0`                |
| `.beluga/project.yaml`       | Project metadata (schema-version, name, template, beluga-version) |
| `.beluga/eval.yaml`          | `beluga eval` config — metrics, row timeout, parallelism          |
| `.beluga/eval.smoke.json`    | 3-row mock-provider smoke dataset for `make eval-ci`              |
| `Makefile`                   | `build`, `test`, `test-ci`, `eval-ci`, `check`, `lint` targets    |
| `.github/workflows/ci.yml`   | `check` + `eval-smoke` jobs on push and pull_request              |
| `Dockerfile`                 | Multi-stage build producing a distroless runtime image            |
| `.env.example`               | `OPENAI_API_KEY` + `BELUGA_OTEL_STDOUT=1` starter                 |
| `.gitignore`                 | `.env` exclusions, build output, editor detritus                  |

## Run the agent

The scaffolded `main.go` wires the `openai` provider by default. Set an API key and
run:

```bash
cp .env.example .env          # then edit OPENAI_API_KEY
export OPENAI_API_KEY=sk-...
go run .
```

Expected output: the agent invokes the echo tool and returns `hello world`. OTel
spans pretty-print to stdout because `.env.example` ships `BELUGA_OTEL_STDOUT=1`.

## Run the smoke eval (no API key required)

`make eval-ci` executes `.beluga/eval.smoke.json` through the `mock` LLM provider —
the fixture queue in the dataset drives the assistant trajectory, so there are no
paid calls and no flakiness:

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.0
make eval-ci
```

Each row runs the agent in a subprocess (`BELUGA_ENV=eval`) that receives the sample
via `BELUGA_EVAL_SAMPLE_JSON`, populates `Output` + `Metadata["latency_ms"]`, and
emits the result on stdout. The parent joins rows into `eval-report.json`, which the
CI workflow uploads as an artefact.

## Swap the provider

The blank-import line registers `openai`; swap it for any of:

```go
_ "github.com/lookatitude/beluga-ai/v2/llm/providers/anthropic"
_ "github.com/lookatitude/beluga-ai/v2/llm/providers/ollama"
_ "github.com/lookatitude/beluga-ai/v2/llm/providers/google"
```

…and update the provider name / model in the `llm.New(...)` call. The full provider
matrix lives in the framework reference docs.

## Related examples

- [`dev-loop-demo`](../dev-loop-demo/) — same scaffolded shape, walks through
  `beluga dev` file-watch + `beluga dev --playground` three-panel UI.
- [`eval-demo`](../eval-demo/) — same scaffolded shape, expands
  `.beluga/eval.smoke.json` to 4 rows and walks through the four common
  `beluga eval` use cases (smoke, regression, PR gate, real-provider).
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference; new projects should start from this scaffolded-basic shape.
