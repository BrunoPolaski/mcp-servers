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

Lower case, words joined by `-`. The date is the day the round started.

## Layout

```
<YYYY-MM-DD>-<round>/
  mcp/<run_id>.json MCP servers the host connected to, with the session's X-Run-Id header
                    (rounds before run_id have a single mcp.json without it)
  results.jsonl     one line per session (keys below)
  progress.txt      console output of run_evaluation.py
  metrics.md        report from analyze_metrics.py
  transcripts/NN.jsonl   raw stream-json of session NN (scenario number, 2 digits)
  logs/<server>.jsonl    server traces, one file per MCP server (not versioned)
  workdir/          empty working directory of the sessions (not versioned)
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
```

Each session tags its MCP requests with `X-Run-Id`, and `analyze_metrics.py`
groups server traces by that `run_id`, so sessions may overlap in time. Traces
without `run_id` (the three rounds above) are grouped by session window, which
assumes sessions ran one at a time. `run_evaluation.py` itself still runs
sessions sequentially; `--rep` sets the repetition part of the run id.
