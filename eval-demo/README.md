# eval-demo

**What this is:** a scaffolded project (structurally identical to
[`scaffolded-basic`](../scaffolded-basic/)) used as the walk-through target for
`beluga eval` — the per-row subprocess evaluation runner introduced in framework
`v2.13.0`. The code is the canonical `beluga init` output; the README walks
through four common eval use cases. Verified against the `beluga` CLI at the
version pinned in `go.mod`.

The only content difference vs. scaffolded-basic: `.beluga/eval.smoke.json` is
expanded from 3 to 4 rows, all exercising the echo tool.

## How `beluga eval` works

`beluga eval <dataset.json>` builds your project once, then runs one subprocess
per dataset row. For each row:

1. The parent serialises the `EvalSample` (input + expected output + recorded turn
   trajectory) and hands it to the child via `BELUGA_EVAL_SAMPLE_JSON`, together
   with `BELUGA_ENV=eval`.
2. The child (your agent's `main.go`) detects `BELUGA_ENV=eval`, routes through
   `runEvalMode`, prints the protocol probe `{"beluga_eval_protocol":1}`, invokes
   the agent against the sample's input, and emits the populated sample (with
   `Output` + `Metadata["latency_ms"]`) on stdout.
3. The parent reads stdout, scores each row with the metrics configured in
   `.beluga/eval.yaml`, prints a summary table and writes `eval-report.json`.
   The metric set is closed: `exact_match` and `latency` are the only built-in
   names (`--metric` with any other name is rejected).
4. `beluga eval` exits non-zero when any row **errors** — the child crashed,
   timed out (`--row-timeout`, default 30s), broke the stdout protocol, or a
   metric could not be computed. A row that runs cleanly but scores
   `exact_match: 0.00` is reported in the table and in `eval-report.json` but
   does **not** change the exit code (verified against the v2.13.6 CLI). If you
   need a score threshold to fail CI, check `aggregate` in `eval-report.json`
   yourself (for example with `jq`).

Exec-once-per-row is the isolation boundary: no shared process state, no fixture
bleed between rows, and one crashed row is recorded as that row's error without
aborting the rest of the run.

## The four common use cases

### 1. Smoke test — PR gate without API keys

The scaffolded `.github/workflows/ci.yml` runs `make eval-ci` on every PR (and
this repo's root CI runs it for each scaffolded example):

```bash
make eval-ci
```

The `eval-ci` Makefile target sets `BELUGA_LLM_PROVIDER=mock` +
`BELUGA_DETERMINISTIC=1` + `BELUGA_SEED=42` + `OTEL_SDK_DISABLED=true`, so the
child subprocess replays fixtures from `Turns` instead of calling a real LLM.
Zero API keys, zero cost; after the one-off build the rows finish in well under a
second. The gate fails on row errors (see above), not on a low `exact_match`.

Expected output (4/4 rows pass; row IDs, run ID and duration vary per run):

```
wrote <project-root>/eval-report.json
IDX  ROW_ID                 INPUT                         OUTPUT           EXPECTED         exact_match  latency  ERROR
0    ...-0000               Please echo: hello world      hello world      hello world      1.00         1.00
1    ...-0001               Please echo: beluga           beluga           beluga           1.00         1.00
2    ...-0002               Please echo: 42               42               42               1.00         1.00
3    ...-0003               Please echo: quick brown fox  quick brown fox  quick brown fox  1.00         1.00

run_id:  ...
dataset: eval-demo-smoke (.beluga/eval.smoke.json)
samples: 4
duration: ...
aggregate:
  exact_match: 1.0000
  latency: 1.0000
```

### 2. Regression — pin the expected answer

Add a row to `.beluga/eval.smoke.json` with the input, the answer you expect, and
the assistant trajectory the mock provider should replay:

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

Know what the mock run can and cannot catch. The mock provider replays the
assistant turns from `Turns` in order, whatever your code does, so `make eval-ci`
catches changes that break the run (the agent no longer builds, a tool errors or
is missing, the eval branch crashes) and changes to how your code turns the
model's final message into `Output` (those show up as `exact_match: 0.00`). It
does not catch a change in which tools a real model would choose:
`ExpectedTools` is carried in the dataset but is not scored by the built-in
metrics. Use the real-provider run (use case 4) for behavioural drift.

### 3. Debug a failing row

Rows already run one at a time by default (`parallel: 1` in `.beluga/eval.yaml`;
keep it at 1 in mock mode, because the mock fixture queue is shared). To narrow
down a failure:

```bash
export BELUGA_LLM_PROVIDER=mock BELUGA_DETERMINISTIC=1 BELUGA_SEED=42 OTEL_SDK_DISABLED=true
beluga eval --dry-run .beluga/eval.smoke.json      # row count + row timeout, no subprocesses
beluga eval --max-rows 1 .beluga/eval.smoke.json   # only the first row
beluga eval --format junit .beluga/eval.smoke.json # also writes a JUnit XML next to the report
```

The child's stdout is consumed by `beluga eval` as the IPC channel, so put debug
output on **stderr** (`log.Printf` does this) — the child's stderr is forwarded
to your terminal. Each failing row's message is in the `ERROR` column and in
`samples[i].error` / `errors` in `eval-report.json`.

### 4. Real-provider eval — opt-in

The scaffolded `ci.yml` ships a commented-out `eval-real` job — uncomment it, add
an `OPENAI_API_KEY` repository secret, and trigger it manually or on a schedule.
Without `BELUGA_LLM_PROVIDER=mock`, the eval branch in `main.go` calls the real
`openai` provider (it fails fast if `OPENAI_API_KEY` is unset) and ignores the
recorded `Turns`:

```bash
beluga eval --max-rows 20 .beluga/eval.smoke.json
```

`--max-rows N` runs only the **first** N rows of the dataset. There is no cost
cap yet: `--max-cost` (used in the commented template) is accepted but reserved
for a future release and currently has no effect, so bound spend with
`--max-rows`, `--row-timeout` and the provider's own budget limits. The default
`eval.yaml` config still applies — the dataset is the same, only the provider and
the invocation flags change.

Expect `exact_match` to be stricter against a real model than against the
replayed fixtures (a trailing period or extra word scores 0), so treat the real
run as a canary to read, not a hard gate.

## Regenerate locally

```bash
go install github.com/lookatitude/beluga-ai/v2/cmd/beluga@v2.13.11
beluga init eval-demo \
  --module github.com/lookatitude/beluga-examples/eval-demo
cd eval-demo && go mod tidy
```

Differences from the scaffolded output (otherwise identical apart from the
`beluga init` timestamps — see
[`scaffolded-basic/README.md`](../scaffolded-basic/README.md#regenerate-locally)):

- This README.
- `.beluga/eval.smoke.json` — extended from 3 rows (scaffolded default) to 4 rows
  covering `hello world`, `beluga`, `42`, `quick brown fox`.
- The documented `ci.yml` / `Makefile` hand-edits shared with scaffolded-basic
  (versioned CLI install, pinned golangci-lint).

## Related examples

- [`scaffolded-basic`](../scaffolded-basic/) — the `beluga init` output this
  project is structurally copied from.
- [`dev-loop-demo`](../dev-loop-demo/) — same shape, walks through `beluga dev`.
- [`first-agent`](../first-agent/) — pre-v2.13.0 hand-written example, kept as a
  historical reference.
