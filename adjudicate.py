#!/usr/bin/env python3
import argparse
import json
import re
import sys
from datetime import date
from pathlib import Path


def bands(x, cuts):
    for limit, points in cuts:
        if x >= limit:
            return points
    return 0


def year_before(ref):
    return date(ref.year - 1, ref.month, 28 if (ref.month, ref.day) == (2, 29) else ref.day)


def day(s):
    return date.fromisoformat(s[:10])


def income(bureau):
    fp = bureau.get("financial_profile") or {}
    for key in ("declared_monthly_income", "estimated_monthly_income"):
        if fp.get(key) is not None:
            return fp[key]
    salaries = [e["salary"] for e in bureau.get("employment_records") or []
                if e.get("is_current") and e.get("salary") is not None]
    return max(salaries) if salaries else None


def c2(records, ref):
    bureau = records.get("bureau") or {}
    fp = bureau.get("financial_profile") or {}
    ratio = fp.get("debt_to_income_ratio")
    if ratio is None and fp.get("total_monthly_payments") is not None and income(bureau):
        ratio = fp["total_monthly_payments"] / income(bureau)
    if ratio is None:
        return 0
    ratio = round(ratio, 2)
    return 11 if ratio <= 0.30 else 8 if ratio <= 0.40 else 4 if ratio <= 0.50 else 0


def on_time_share(rows, ref):
    start = year_before(ref)
    window = [r for r in rows if r.get("due_date") and start < day(r["due_date"]) <= ref]
    if not window:
        return None
    return sum(r.get("status") == "on_time" for r in window) / len(window)


def c3(records, ref):
    share = on_time_share((records.get("bureau") or {}).get("payment_histories") or [], ref)
    return 4 if share is None else bands(share, [(0.95, 11), (0.85, 8), (0.70, 4)])


def c12(records, ref):
    share = on_time_share((records.get("internal-registry") or {}).get("internal_payment_records") or [], ref)
    return 2 if share is None else bands(share, [(0.95, 5), (0.85, 4), (0.70, 2)])


RECOMPUTE = {"C2": c2, "C3": c3, "C12": c12}


def records_from_transcript(path):
    merged = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        event = json.loads(line)
        if event.get("type") != "user" or not isinstance(event["message"]["content"], list):
            continue
        for item in event["message"]["content"]:
            if not isinstance(item, dict) or item.get("type") != "tool_result":
                continue
            content = item["content"]
            text = content if isinstance(content, str) else "".join(c.get("text", "") for c in content)
            try:
                payload = json.loads(text)
            except ValueError:
                continue
            if not isinstance(payload, dict) or "metadata" not in payload:
                continue
            merge(merged.setdefault(payload["metadata"].get("source"), {}), payload)
    return merged


MERGED_TOOLS = {"get_internal_payment_records": "internal_payment_records",
              "get_customer_by_id": None, "get_customer_by_document": None}


def merge(target, payload):
    tool = payload["metadata"].get("tool")
    if tool not in MERGED_TOOLS:
        return
    for key, value in payload.items():
        if key == "metadata" or (MERGED_TOOLS[tool] and key != "items"):
            continue
        key = MERGED_TOOLS[tool] or key
        if isinstance(value, list):
            rows = target.setdefault(key, [])
            rows.extend(v for v in value if v not in rows)
        elif target.get(key) is None:
            target[key] = value
        elif value is not None and target[key] != value:
            raise ValueError(f"{payload['metadata'].get('source')}: responses disagree on {key}")


def reference_dates(session, records, fixed):
    if fixed:
        return [fixed]
    found = re.search(r"(\d{2})/(\d{2})/(\d{4})", session.get("prompt") or "")
    if found:
        d, m, y = map(int, found.groups())
        return [date(y, m, d)]
    snapshot = day((records.get("bureau") or {}).get("created_at") or session["started_at"])
    return sorted({day(session["started_at"]), snapshot})


def adjudicate(criterion, key, agent, records, refs):
    values = {RECOMPUTE[criterion](records, ref) for ref in refs}
    if len(values) > 1:
        return sorted(values), "politica"
    value = values.pop()
    if value == agent:
        return value, "gabarito"
    if value == key:
        return value, "agente"
    return value, "politica"


def divergences_from_key(sessions, key):
    for s in sessions:
        expected = key.get(str(s["scenario"]))
        if not expected or expected["pontuacao"] is None:
            continue
        got = (s.get("block") or {}).get("criterios") or {}
        for criterion, points in expected["criterios"].items():
            if got.get(criterion) != points:
                yield s, criterion, points, got.get(criterion)


def divergences_from_metrics(sessions, metrics):
    by_scenario = {str(s["scenario"]): s for s in sessions}
    for line in metrics.read_text(encoding="utf-8").split("### Divergencias por criterio")[-1].splitlines():
        cells = [c.strip() for c in line.strip("|").split("|")]
        if len(cells) != 3 or not re.fullmatch(r"C\d+", cells[0]):
            continue
        for scenario, expected, got in re.findall(r"(\d+): ([^,\s]+) -> ([^,\s]+)", cells[2]):
            yield by_scenario[scenario], cells[0], int(expected), None if got == "None" else int(got)


def transcript_path(round_dir, session):
    name = session.get("run_id") or f"{session['scenario']:02d}"
    if "-" in name:
        n, rep = name.split("-")
        name = f"{int(n):02d}-{rep}"
    return round_dir / "transcripts" / f"{name}.jsonl"


def run(round_dir, key=None, fixed=None):
    sessions = [json.loads(l) for l in (round_dir / "results.jsonl").read_text(encoding="utf-8").splitlines() if l.strip()]
    found = (divergences_from_key(sessions, key) if key is not None
             else divergences_from_metrics(sessions, round_dir / "metrics.md"))
    rows = []
    for session, criterion, expected, got in found:
        if criterion not in RECOMPUTE:
            rows.append((session, criterion, expected, got, None, "sem recalculo", None))
            continue
        records = records_from_transcript(transcript_path(round_dir, session))
        refs = reference_dates(session, records, fixed)
        value, origin = adjudicate(criterion, expected, got, records, refs)
        rows.append((session, criterion, expected, got, value, origin, refs))
    return rows


def render(rows):
    lines = ["| Sessao | Criterio | Gabarito | Agente | Recalculo | Data de referencia | Origem |",
             "|---|---|---|---|---|---|---|"]
    for session, criterion, expected, got, value, origin, refs in rows:
        shown = "-" if value is None else " ou ".join(map(str, value)) if isinstance(value, list) else value
        dates = "-" if not refs else ", ".join(r.isoformat() for r in refs)
        lines.append(f"| {session.get('run_id') or session['scenario']} | {criterion} | {expected} | {got} "
                     f"| {shown} | {dates} | {origin} |")
    counts = {}
    for row in rows:
        counts[row[5]] = counts.get(row[5], 0) + 1
    lines += ["", "Origem: " + (", ".join(f"{k} {v}" for k, v in sorted(counts.items())) or "nenhuma divergencia")]
    return "\n".join(lines) + "\n"


def selftest():
    hist = [{"due_date": f"2025-{m:02d}-05", "status": "late"} for m in (9, 10)]
    hist += [{"due_date": f"2026-{m:02d}-05" if m <= 8 else f"2025-{m:02d}-05", "status": "on_time"}
             for m in (11, 12, 1, 2, 3, 4, 5, 6, 7, 8)]
    records = {"bureau": {"payment_histories": hist, "financial_profile": {"declared_monthly_income": 4000}}}
    assert c3(records, date(2026, 8, 31)) == 4
    assert c3(records, date(2026, 9, 27)) == 8
    assert c3({}, date(2026, 8, 31)) == 4
    assert c2(records, None) == 0
    records["bureau"]["financial_profile"]["total_monthly_payments"] = 1200
    assert c2(records, None) == 11
    assert adjudicate("C3", 4, 8, records, [date(2026, 8, 31)]) == (4, "agente")
    assert adjudicate("C3", 4, 8, records, [date(2026, 9, 27)]) == (8, "gabarito")
    assert adjudicate("C3", 0, 8, records, [date(2026, 8, 31)]) == (4, "politica")
    assert adjudicate("C3", 4, 8, records, [date(2026, 8, 31), date(2026, 9, 27)]) == ([4, 8], "politica")
    session = {"started_at": "2026-09-27T23:00:00-03:00"}
    assert reference_dates(session, {"bureau": {"created_at": "2026-08-31 00:00:00"}}, None) == [date(2026, 8, 31), date(2026, 9, 27)]
    assert reference_dates({**session, "prompt": "Data da análise: 31/08/2026."}, {}, None) == [date(2026, 8, 31)]
    meta = {"source": "internal-registry", "tool": "get_customer_by_id"}
    row = {"id": 1, "due_date": "2026-08-05", "status": "on_time"}
    late = {"id": 2, "due_date": "2026-07-05", "status": "late"}
    merged = {}
    merge(merged, {"metadata": meta, "id": 7, "internal_payment_records": [row]})
    merge(merged, {"metadata": {**meta, "tool": "get_internal_payment_records"}, "items": [row, late]})
    merge(merged, {"metadata": {**meta, "tool": "get_pre_approved_limits"}, "items": [{"id": 9}]})
    assert merged == {"id": 7, "internal_payment_records": [row, late]}
    try:
        merge(merged, {"metadata": meta, "id": 8})
        raise AssertionError("conflicting responses accepted")
    except ValueError:
        pass
    print("selftest ok")
    return 0


def main(argv=None):
    p = argparse.ArgumentParser(description="Atribui cada divergencia entre agente e gabarito recalculando o criterio a partir dos registros brutos.")
    p.add_argument("round", nargs="?", type=Path, help="pasta da rodada em runs/")
    p.add_argument("--gabarito", type=Path, help="gabarito JSON; sem ele, as divergencias vem do metrics.md da rodada")
    p.add_argument("--date", type=date.fromisoformat, help="data da analise (AAAA-MM-DD), para rodadas que nao a registram em results.jsonl")
    p.add_argument("--selftest", action="store_true")
    args = p.parse_args(argv)
    if args.selftest:
        return selftest()
    if args.round is None:
        p.error("informe a pasta da rodada")
    key = json.loads(args.gabarito.read_text(encoding="utf-8")) if args.gabarito else None
    sys.stdout.write(render(run(args.round, key, args.date)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
