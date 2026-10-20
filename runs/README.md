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

The `metrics.md` of the exploratory round was generated with the answer key of
that time. Regenerating it with the current `docs/gabarito_cenarios.json` changes
the expected-score column.

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

Lower case, words joined by `-`. The date is the day the round started.

## Layout

```
<YYYY-MM-DD>-<round>/
  mcp/<run_id>.json MCP servers the host connected to, with the session's X-Run-Id header
                    (rounds before run_id have a single mcp.json without it)
  results.jsonl     one line per session (keys below)
  progress.txt      console output of run_evaluation.py
  metrics.md        report from analyze_metrics.py
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
