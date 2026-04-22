# eval-demo

**What this is:** a scaffolded project (structurally identical to
[`scaffolded-basic`](../scaffolded-basic/)) used as the walk-through target for
`beluga eval` — the per-row subprocess evaluation runner introduced in framework
`v2.13.0`. The code is the canonical `beluga init` output; the README walks
through four common eval use cases.

The only content difference vs. scaffolded-basic: `.beluga/eval.smoke.json` is
expanded from 3 to 4 rows, all exercising the echo tool.

## How `beluga eval` works

`beluga eval <dataset.json>` runs one subprocess per dataset row. For each row:

1. The parent serialises the `EvalSample` (input + expected output + recorded turn
   trajectory) and hands it to the child via `BELUGA_EVAL_SAMPLE_JSON`.
2. The child (your agent's `main.go`) detects `BELUGA_ENV=eval`, routes through
   `runEvalMode`, invokes the agent against the sample's input, and emits the
   populated sample (with `Output` + `Metadata["latency_ms"]`) on stdout.
3. The parent reads stdout, joins rows into `eval-report.json`, and applies the
   metrics configured in `.beluga/eval.yaml` (`exact_match`, `latency` by default;
   `pass_rate` is available too).
4. Non-zero row failures propagate to a non-zero exit from `beluga eval` so the CI
   step fails visibly.

Exec-once-per-row is the isolation boundary: no shared process state, no fixture
bleed between rows, and one crashed row does not take down the rest of the run.

## The four common use cases

### 1. Smoke test — PR gate without API keys

The scaffolded `.github/workflows/ci.yml` runs `make eval-ci` on every PR:

```bash
make eval-ci
```

The `eval-ci` Makefile target sets `BELUGA_LLM_PROVIDER=mock` +
`BELUGA_DETERMINISTIC=1` + `BELUGA_SEED=42` + `OTEL_SDK_DISABLED=true`, so the
child subprocess replays fixtures from `Turns` instead of calling a real LLM.
Zero API keys, zero cost, ~1 ms per row. Failure exit code is surfaced to CI.

Expected output (4/4 rows pass):

```
IDX  ROW_ID                 INPUT                       OUTPUT            EXPECTED          exact_match  latency
0    ...                    Please echo: hello world    hello world       hello world       1.00         1.00
1    ...                    Please echo: beluga         beluga            beluga            1.00         1.00
2    ...                    Please echo: 42             42                42                1.00         1.00
3    ...                    Please echo: quick brown…   quick brown fox   quick brown fox   1.00         1.00
aggregate:
  exact_match: 1.0000
  latency: 1.0000
```

### 2. Regression — detect output drift

Add a new row to `.beluga/eval.smoke.json` capturing the fixture you *expect* your
agent to produce for a given input. If a later code change alters the trajectory,
the replay mismatches the captured `Turns` and the row fails:

```json
{
  "Input": "Please echo: new-fixture",
  "ExpectedOutput": "new-fixture",
  "Turns": [
    {"Role": "assistant", "ToolCalls": [{"ID": "call_1", "Name": "echo", "Arguments": "{\"message\":\"new-fixture\"}"}]},
    {"Role": "tool", "Content": "new-fixture"},
    {"Role": "assistant", "Content": "new-fixture"}
  ],
  "ExpectedTools": ["echo"]
}
```

Run `make eval-ci` before and after your change — any drift in the assistant's
decided tool calls will surface as a row failure.

### 3. Debug a failing row — `--parallel 1`

When a row fails and the logs are muddled because of concurrent output, force
serial execution for clarity:

```bash
BELUGA_LLM_PROVIDER=mock \
BELUGA_DETERMINISTIC=1 \
OTEL_SDK_DISABLED=true \
beluga eval --parallel 1 .beluga/eval.smoke.json
```

`--parallel 1` runs rows sequentially; combined with `OTEL_SDK_DISABLED=true`, the
child subprocess's stdout is pure IPC JSON, so you can tee it to a file and inspect
the exact payload that failed. This is the recommended posture when
`.beluga/eval.yaml` already sets `parallel: 1` (the default) but you want to be
explicit on the command line while varying other flags.

### 4. Real-provider eval — opt-in, cost-capped

The scaffolded `ci.yml` ships a commented-out `eval-real` job — uncomment it,
paste an `OPENAI_API_KEY` secret into the repo settings, and trigger it manually
or on a schedule:

```bash
beluga eval --max-rows 20 --max-cost 5.00 .beluga/eval.smoke.json
```

`--max-rows` caps the number of rows processed (newest first); `--max-cost` caps
the total spend and aborts mid-run if the estimate crosses the threshold. The
default `eval.yaml` config still applies — the dataset is the same, only the
provider and the invocation flags change.

Because the mock fixtures are aligned with the real model's tool-call shape, the
real-provider run exercises the same trajectory structure and metrics, giving you
a paid canary on top of the free smoke gate.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.0
beluga init eval-demo \
  --module github.com/lookatitude/beluga-examples/eval-demo
```

Files hand-authored relative to the scaffolded output:

- This README.
- `.beluga/eval.smoke.json` — extended from 3 rows (scaffolded default) to 4 rows
  covering `hello world`, `beluga`, `42`, `quick brown fox`.
- `name:` fields in `.beluga/project.yaml` + project-name substitutions in
  `go.mod`, `main.go`, `Dockerfile`, `.gitignore` (all trivially derivable from
  the `--module` arg).

## Related examples

- [`scaffolded-basic`](../scaffolded-basic/) — the verbatim `beluga init` output
  this project is structurally copied from.
- [`dev-loop-demo`](../dev-loop-demo/) — same shape, walks through `beluga dev`.
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference.
