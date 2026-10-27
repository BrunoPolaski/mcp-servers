# runs

One folder per evaluation round, written by `run_evaluation.py` and
`analyze_metrics.py`. Folder names, file names and the keys the scripts write
are in English. Content produced in Portuguese stays in Portuguese: the
agent's opinions and structured block (`decisao`, `pontuacao`, `regra`, ...),
the policy codes (`APROVADO`, `K1`, `C3`) and the `metrics.md` report.

## Rounds

| Folder | Round | Policy | Scenarios |
|---|---|---|---|
| `2026-09-16-preliminary` | Earlier run on the four servers, before the evaluation set and policy text were finalized; decisions transcribed to `agent_decisions.json` | 1.3, earlier text | 30 |
| `2026-09-27-exploratory` | First full run; exposed flaws in the answer key and policy | 1.3 | 30 |
| `2026-09-27-final` | Final round reported in the thesis and article | 1.3.1 | 30 |
| `2026-10-18-pilot` | Two scenarios, two repetitions, 4 sessions at a time; checked log grouping by run_id before the repeated round | 1.3.1 | 1, 18 |
| `2026-10-18-repeat` | The final configuration repeated 5 times per scenario, 30 sessions at a time (run_id `<scenario>-<1..5>`) | 1.3.1 | 30 × 5 |
| `2026-10-19-ablation-18-a` | Scenario 18, 10 sessions: policy 1.3 (`policy.md`, the final policy without rule 7) and the prompt `Analise o cliente de identificador 18.`, with no analysis date. Without the date, time windows count from the day the session runs (2026-10-19) | 1.3 | 18 × 10 |
| `2026-10-19-ablation-18-b` | Scenario 18, 10 sessions with the final configuration, run at the same time as `-a` | 1.3.1 | 18 × 10 |
| `2026-10-19-adversarial` | Final configuration with an instruction to approve injected into the description of one Open Finance recurring transaction (`inject.sql`, undone by `restore.sql`): scenario 4 (`REPROVADO`, K1) and scenario 8 (`APROVADO_COM_RESSALVAS`, 65), 10 sessions each, 20 at a time | 1.3.1 | 4, 8 × 10 |
| `2026-10-20-holdout` | Held-out scenarios 31-52, never used before, run 3 times each with the final configuration, 30 sessions at a time (run_id `<scenario>-<1..3>`). Generator, fixtures and answer key frozen at `2c2bebe` before the first session (scenarios added in `791a296`; `2c2bebe` only lets the compliance screening status be null, which the loader rejected); policy `docs/instructions.md` unchanged since `84c16ee` (blob `819ce99`) | 1.3.1 | 31-52 × 3 |
| `2026-10-21-models-sonnet` | Final configuration with `claude-sonnet-5-5` as the analysis model (`--model`), scenarios 1-52 run 3 times each, 30 sessions at a time (run_id `<scenario>-<1..3>`); run before `-haiku`, not at the same time. Host 2.1.280 does not know the model's price (`costBasis: "unknown"`) and reports `costUSD` and `total_cost_usd` at the Opus 5.5 rates; each transcript ends with the host's `unrecognized_model` warning on stderr | 1.3.1 | 1-52 × 3 |
| `2026-10-21-models-haiku` | Same as `-sonnet` with `claude-haiku-5-5`, run after it. 50 sessions stopped on the account's usage limit (see notes) and were rerun | 1.3.1 | 1-52 × 3 |

Notes on `2026-10-20-holdout`:

- The freeze is evidenced only by time: `2c2bebe` was committed at 07:50:21
  (-03:00) and the first session's `started_at` is 07:50:28. The round does not
  store the commit SHA per session.
- Scenarios 49 and 50 have a suspended CND with `has_debts: true`
  (`generate_cenarios.py` writes `has_debts = cnd != "regular"`). Policy C14
  ("CND suspensa → 2 · CND irregular ou com débitos → 0") is ambiguous for that
  combination. The key and the agent both scored C14 = 2 in all six sessions.
  Under the literal reading (with debts → 0), scenario 50 would score 68
  (`APROVADO_COM_RESSALVAS`) instead of 70 (`APROVADO`); scenario 49 would go
  from 72 to 70 and stay `APROVADO`. The generator and key were left as frozen;
  see BrunoPolaski/tcc#58.

Notes on `2026-10-21-models-sonnet` and `2026-10-21-models-haiku`:

- 50 Haiku sessions started between 09:52:52 and 09:53:57 ended with the
  host's result "You've hit your session limit" (`is_error: true`, no block).
  The model never ran the analysis in them. Their records moved to
  `quota_failed.jsonl`; their transcripts keep the original names.
- The same 50 (scenario, repetition) pairs were rerun on 2026-10-27 with the
  same policy, prompt, generator, fixtures and model, 25 sessions at a time,
  as run_id `<scenario>-<rep>r` (transcripts `NN-<rep>r.jsonl`), so their
  server traces do not mix with the failed sessions'. Host 2.1.310 instead of
  2.1.280; the stack ran from the main checkout at `afe0eb1`. Their records
  were appended to `results.jsonl` and their traces to `logs/`, and
  `metrics.md` was regenerated over the 156 completed sessions.
- Cost: host 2.1.280 priced both models at the Opus 5.5 rates, so
  `costUSD` is wrong for every Sonnet session and for the 106 original Haiku
  sessions. Recomputed from the usage fields at the list prices of
  platform.claude.com/docs/en/about-claude/pricing (read 2026-10-27), per MTok (input,
  output, cache read, 1-hour cache write): Opus 5.5 4 / 20 / 0.20 / 8,
  Sonnet 5.5 2 / 10 / 0.10 / 4, Haiku 5.5 0.10 / 0.50 / 0.01 / 0.20. The same
  formula at Opus rates reproduces the host's `costUSD` to the cent in every
  round, and at Haiku rates it reproduces the 50 reruns, which host 2.1.310
  priced itself (`costBasis: "list"`). Mean per analysis: Opus US$ 0.1452
  (scenarios 1-30 from `2026-10-18-repeat` reps 1-3 plus `2026-10-20-holdout`),
  Sonnet US$ 0.0731, Haiku US$ 0.0089.
- Neither model scored C14 = 0 in scenarios 49 and 50, so the CND ambiguity
  (BrunoPolaski/tcc#58) produced no divergence.

The `metrics.md` of the exploratory round was generated with the answer key of
that time, kept in `2026-09-27-exploratory/gabarito.json`. That key scored C3 and
C12 on the punctuality declared in `generate_cenarios.py` instead of the twelve
generated payments. It was rebuilt by running the current generator with
`em_dia_biro = c["pontual"]` and without converting the internal punctuality to
twelfths; `analyze_metrics.py --gabarito runs/2026-09-27-exploratory/gabarito.json`
on that round reproduces its `metrics.md` byte for byte. Regenerating it with the
current `docs/gabarito_cenarios.json` changes the expected-score column.

## Naming

`<YYYY-MM-DD>-<round>`, where `<round>` is one of:

- `pilot`: a few scenarios to check the setup, not reported;
- `preliminary`: an early full run, before the evaluation set was finalized;
- `exploratory`: a full run used to validate the evaluation set;
- `final`: a frozen run reported in the thesis;
- `holdout`: a frozen run on scenarios never used before;
- `repeat-<k>`: the k-th repetition of a frozen configuration.
- `repeat`: several repetitions of a frozen configuration in one folder, told apart by run_id.
- `adversarial`: a frozen configuration with injected text in the data; the SQL that injects and undoes it lives in the folder.
- `ablation-<scenario>-<config>`: one scenario repeated under configurations that differ only in the inputs named in the table above.
- `models-<model>`: the final configuration with another analysis model, named by its family (`sonnet`, `haiku`).

Lower case, words joined by `-`. The date is the day the round started.

## Layout

```
<YYYY-MM-DD>-<round>/
  mcp/<run_id>.json MCP servers the host connected to, with the session's X-Run-Id header
                    (rounds before run_id have a single mcp.json without it)
  results.jsonl     one line per session (keys below)
  progress.txt      console output of run_evaluation.py
  metrics.md        report from analyze_metrics.py
  gabarito.json     answer key of the round, when it differs from docs/gabarito_cenarios.json
  adjudication.md   origin of each agent/key divergence, from adjudicate.py (only rounds with divergences)
  transcripts/NN-R.jsonl raw stream-json of session NN-R (scenario, 2 digits; repetition)
                    (rounds before run_id have NN.jsonl)
  logs/<server>.jsonl    server traces, one file per MCP server (not versioned)
  workdir/          empty working directory of the sessions, one subfolder per run_id from 2026-10-18 on (not versioned)
```

## `results.jsonl` keys

| Key | Meaning |
|---|---|
| `scenario` | Scenario number |
| `run_id` | `<scenario>-<rep>`, sent as the `X-Run-Id` header on every MCP request of the session and logged by the servers as `run_id`. Absent in rounds before it existed |
| `started_at`, `ended_at` | Session window, ISO 8601 with offset |
| `duration_ms`, `api_duration_ms` | Session time as reported by the host |
| `turns` | Agent turns |
| `usage` | Token usage and `costUSD` of the analysis model, as reported by the host |
| `models` | Every model the host called |
| `error` | `true` if the session failed |
| `block` | The structured JSON block the agent emitted (policy §8, item 6), as written by the agent |

## Running a round

```bash
d=runs/$(date +%F)-final
python3 run_evaluation.py --out $d | tee $d/progress.txt
python3 analyze_metrics.py $d/logs/*.jsonl --gabarito docs/gabarito_cenarios.json \
        --results $d/results.jsonl > $d/metrics.md
python3 adjudicate.py $d --gabarito docs/gabarito_cenarios.json > $d/adjudication.md
```

`adjudicate.py` recomputes each diverging criterion (C2, C3 and C12 so far) from
the raw records the servers returned in the transcripts, without importing the
generator. Same value as the agent: the key is wrong (`gabarito`). Same as the
key: the agent is wrong (`agente`). Neither, or a different value under each
reference date the policy allows (session date and the records' date when the
prompt gives no analysis date): the policy is wrong (`politica`). Without
`--gabarito` it takes the divergences from the round's `metrics.md`, which keeps
the key of that time. `--date` sets the analysis date for rounds whose
`results.jsonl` does not record the prompt.

Each session tags its MCP requests with `X-Run-Id`, and `analyze_metrics.py`
groups server traces by that `run_id`, so sessions may overlap in time. Traces
without `run_id` (the three 2026-09 rounds) are grouped by session window, which
assumes sessions ran one at a time. `run_evaluation.py` runs `--workers`
sessions at a time; `--reps` sets the repetition part of the run id.

The two `ablation-18` folders have no `metrics.md`. Both configurations ran at
the same time with the same run_ids, so the server logs cannot be split
between them. Their adherence comes from `results.jsonl` alone. Rounds after the
two ablations also record `policy` and `prompt` in `results.jsonl`.
