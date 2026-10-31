# scaffolded-basic

**What this is:** the output of `beluga init scaffolded-basic` from the released
CLI at the framework version this example pins (the `require` line in
`go.mod`). Run it as-is, or treat it as the canonical starting shape
when you bootstrap a new Beluga AI project.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.12
beluga init scaffolded-basic \
  --module github.com/lookatitude/beluga-examples/scaffolded-basic
cd scaffolded-basic && go mod tidy
```

The regenerated tree is identical to this one **apart from**:

- **Timestamps.** `beluga init` stamps the current time into
  `.beluga/project.yaml` (`scaffolded-at`), the `Dockerfile` header comment and
  the `main.go` package comment, so two runs never match byte-for-byte.
- **`go.mod` / `go.sum`.** `beluga init` writes a bare `require` line; `go mod
  tidy` (and the CI bump job) fill in the indirect requirements and may raise the
  `go` directive to the framework's minimum.
- **Documented hand-edits** (each marked `[beluga-examples hand-edit: …]` in the
  file):
  - `.github/workflows/ci.yml` — the `Install beluga CLI` step installs
    `cmd/beluga@<go.mod version>`. The v2.13.x template runs a bare
    `go install …/cmd/beluga`, which fails inside the module with
    `missing go.sum entry … github.com/spf13/cobra`. golangci-lint is pinned to
    `v2.15.0` instead of `latest`.
  - `Makefile` — the `eval-ci` comment gives the same versioned install command.
- **This README.**

The version stamps (`beluga-version` in `.beluga/project.yaml`, the
`org.opencontainers.image.version` label in the `Dockerfile`, and the
`go install …@vX.Y.Z` lines in these READMEs) are kept in step with `go.mod` by
the repo's framework-bump CI job.

## What you get

| Path                         | Purpose                                                           |
|------------------------------|-------------------------------------------------------------------|
| `main.go`                    | Agent + echo tool + OTel bootstrap + `BELUGA_ENV=eval` branch     |
| `go.mod`                     | Pins `github.com/lookatitude/beluga-ai/v2 v2.13.12`                |
| `.beluga/project.yaml`       | Project metadata (schema-version, name, template, beluga-version) |
| `.beluga/eval.yaml`          | `beluga eval` config — metrics, row timeout, parallelism          |
| `.beluga/eval.smoke.json`    | 3-row mock-provider smoke dataset for `make eval-ci`              |
| `Makefile`                   | `build`, `test`, `test-ci`, `eval-ci`, `check`, `lint` targets    |
| `.github/workflows/ci.yml`   | `check` + `eval-smoke` jobs on push and pull_request              |
| `Dockerfile`                 | Multi-stage build producing a distroless runtime image            |
| `.env.example`               | `OPENAI_API_KEY` + `BELUGA_OTEL_STDOUT=1` starter                 |
| `.gitignore`                 | `.env` exclusions, build output, editor detritus                  |

## Run the agent

The scaffolded `main.go` wires the `openai` provider by default. `go run .` does
**not** read `.env`; only the `beluga` CLI (`beluga run`, `beluga dev`) loads it
into the child's environment. Either use the CLI:

```bash
cp .env.example .env          # then edit OPENAI_API_KEY
beluga run                    # loads .env, builds, runs once
```

…or export the variables yourself and use plain Go:

```bash
export OPENAI_API_KEY=sk-...
export BELUGA_OTEL_STDOUT=1   # optional: pretty-print OTel spans to stdout
go run .
```

Expected output: the agent invokes the echo tool and returns `hello world`. OTel
spans pretty-print to stdout when `BELUGA_OTEL_STDOUT=1` is in the environment —
which `.env.example` sets, so it applies under `beluga run`, or under `go run .`
only if you export it.

## Run the smoke eval (no API key required)

`make eval-ci` executes `.beluga/eval.smoke.json` through the `mock` LLM provider —
the fixture queue in the dataset drives the assistant trajectory, so there are no
paid calls and no flakiness:

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.12
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

…and update the provider name / model in both `llm.New(...)` calls (the
interactive path and the eval path). The full provider matrix lives in the
framework reference docs.

## Related examples

- [`dev-loop-demo`](../dev-loop-demo/) — same scaffolded shape, walks through
  the `beluga dev` file-watch loop and its dev playground.
- [`eval-demo`](../eval-demo/) — same scaffolded shape, expands
  `.beluga/eval.smoke.json` to 4 rows and walks through the four common
  `beluga eval` use cases (smoke, regression, PR gate, real-provider).
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference; new projects should start from this scaffolded-basic shape.
