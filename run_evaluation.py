#!/usr/bin/env python3
"""Runs the evaluation analyses headless, one fresh agent session per scenario.

Each session gets docs/instructions.md as its whole system prompt and only the
four MCP servers as tools: no built-in tools (so no file access), no user or
project settings, hooks or CLAUDE.md, and an empty working directory. The agent
cannot see the answer key, the generator or a previous session.

Sessions run one at a time, so every MCP log record falls inside exactly one
session window [started_at, ended_at], which is how analyze_metrics.py groups them.

Writes to <out>/ (layout and naming in runs/README.md):
  results.jsonl          one line per scenario: the structured block the agent
                         emitted (policy §8, item 6) plus host-side stats
  transcripts/NN.jsonl   raw stream-json of each session

Usage (servers up on :8081-8084):
    python3 run_evaluation.py --out runs/2026-09-27-final | tee runs/2026-09-27-final/progress.txt
    python3 run_evaluation.py --out runs/2026-10-20-pilot --scenarios 1 23
"""

import argparse
import json
import re
import subprocess
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SERVERS = {"bureau": 8081, "open-finance": 8082, "internal-registry": 8083, "registration-validation": 8084}
# the fixtures are built around this date (HOJE in generate_cenarios.py); the
# policy's time windows count from it, whatever day the run happens
PROMPT = "Data da análise: 31/08/2026. Analise o cliente de identificador {n}."
JSON_BLOCK = re.compile(r"```json\s*(.*?)```", re.S)


def last_block(text):
    """The last ```json block of the answer, parsed; None if missing or invalid."""
    blocks = JSON_BLOCK.findall(text or "")
    if not blocks:
        return None
    try:
        return json.loads(blocks[-1])
    except ValueError:
        return None


def result_event(stream):
    for line in reversed(stream.splitlines()):
        try:
            event = json.loads(line)
        except ValueError:
            continue
        if event.get("type") == "result":
            return event
    return {}


def run(n, model, out):
    cmd = ["claude", "-p", PROMPT.format(n=n), "--model", model,
           "--system-prompt-file", str(ROOT / "docs" / "instructions.md"),
           "--mcp-config", str(out / "mcp.json"), "--strict-mcp-config",
           "--tools", "", "--allowedTools", ",".join(f"mcp__{s}" for s in SERVERS),
           "--setting-sources", "", "--disable-slash-commands", "--no-session-persistence",
           "--output-format", "stream-json", "--verbose"]
    started_at = datetime.now().astimezone()
    proc = subprocess.run(cmd, cwd=out / "workdir", stdin=subprocess.DEVNULL, capture_output=True,
                          text=True, timeout=1800)
    ended_at = datetime.now().astimezone()
    (out / "transcripts" / f"{n:02d}.jsonl").write_text(proc.stdout + proc.stderr)
    res = result_event(proc.stdout)
    return {
        "scenario": n,
        "started_at": started_at.isoformat(),
        "ended_at": ended_at.isoformat(),
        "duration_ms": res.get("duration_ms"),
        "api_duration_ms": res.get("duration_api_ms"),
        "turns": res.get("num_turns"),
        # the host also calls a small model for housekeeping; count only the
        # tokens of the model that ran the analysis
        "usage": (res.get("modelUsage") or {}).get(model),
        "models": sorted(res.get("modelUsage") or {}),
        "error": bool(res.get("is_error", True)) or proc.returncode != 0,
        "block": last_block(res.get("result")),
    }


def main():
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("--out", required=True, type=Path)
    p.add_argument("--scenarios", nargs="*", type=int, default=list(range(1, 31)))
    p.add_argument("--model", default="claude-opus-5-5")
    args = p.parse_args()

    out = args.out.resolve()
    (out / "workdir").mkdir(parents=True, exist_ok=True)
    (out / "transcripts").mkdir(exist_ok=True)
    (out / "mcp.json").write_text(json.dumps({"mcpServers": {
        name: {"type": "http", "url": f"http://localhost:{port}/mcp"} for name, port in SERVERS.items()}}, indent=1))

    with open(out / "results.jsonl", "a", encoding="utf-8") as fh:
        for n in args.scenarios:
            r = run(n, args.model, out)
            fh.write(json.dumps(r, ensure_ascii=False) + "\n")
            fh.flush()
            b = r["block"] or {}
            print(f"scenario {n:2d}: {b.get('decisao', 'NO BLOCK')} {b.get('regra') or ''} "
                  f"{b.get('pontuacao') if b.get('pontuacao') is not None else ''} "
                  f"({(r['duration_ms'] or 0) / 1000:.0f}s, {r['turns']} turns{', ERROR' if r['error'] else ''})",
                  flush=True)


if __name__ == "__main__":
    main()
