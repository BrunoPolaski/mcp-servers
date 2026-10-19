#!/usr/bin/env python3
import argparse
import glob
import json
import os
import statistics
import sys

ANALYSIS_DATE = "2026-08-31"
YEAR_AGO = "2025-08-31"
SIX_MONTHS_AGO = "2026-02-28"

CRITERION_SOURCE = {f"C{i}": "bureau" for i in range(1, 7)}
CRITERION_SOURCE |= {f"C{i}": "open-finance" for i in range(7, 10)}
CRITERION_SOURCE |= {f"C{i}": "internal-registry" for i in range(10, 13)}
CRITERION_SOURCE |= {f"C{i}": "registration-validation" for i in range(13, 16)}


def _list(p, key):
    return p.get(key) or []


def _latest(p, key, date):
    rows = _list(p, key)
    return max(rows, key=lambda r: r.get(date) or "") if rows else None


def income(b):
    fp = b.get("financial_profile") or {}
    for key in ("declared_monthly_income", "estimated_monthly_income"):
        if fp.get(key) is not None:
            return fp[key]
    salaries = [e["salary"] for e in _list(b, "employment_records") if e.get("is_current") and e.get("salary") is not None]
    return max(salaries) if salaries else None


def consent(of):
    return any(c.get("status") == "granted" and (c.get("expires_at") or "9") > ANALYSIS_DATE
               for c in _list(of, "data_sharing_consents"))


def cash_flow(of, field):
    row = _latest(of, "cash_flow_analyses", "analysis_date")
    return consent(of) and row is not None and row.get(field) is not None


def recurring_income(of):
    return sum(r.get("amount") or 0 for r in _list(of, "recurring_transactions")
               if r.get("transaction_type") == "income" and r.get("is_active"))


def active_negatives(b):
    return [r.get("amount") or 0 for r in _list(b, "negative_records") if r.get("status") in ("active", None)]


def unpaid_internal(ir):
    return [r for r in _list(ir, "internal_payment_records") if r.get("payment_date") is None and (r.get("days_late") or 0) > 0]


def latest_status(rv, key, date, test):
    row = _latest(rv, key, date)
    return row is not None and test(row)


HAS_DATA = {
    "C1": lambda s: (s["bureau"].get("credit_score") or {}).get("score") is not None,
    "C2": lambda s: (s["bureau"].get("financial_profile") or {}).get("debt_to_income_ratio") is not None
    or ((s["bureau"].get("financial_profile") or {}).get("total_monthly_payments") is not None and income(s["bureau"]) is not None),
    "C3": lambda s: any((r.get("due_date") or "") >= YEAR_AGO for r in _list(s["bureau"], "payment_histories")),
    "C4": None,
    "C5": lambda s: (s["bureau"].get("financial_profile") or {}).get("credit_utilization") is not None,
    "C6": None,
    "C7": lambda s: cash_flow(s["open-finance"], "net_cash_flow"),
    "C8": lambda s: consent(s["open-finance"]),
    "C9": lambda s: cash_flow(s["open-finance"], "negative_balance_days"),
    "C10": lambda s: s["internal-registry"].get("customer_relationship") is not None,
    "C11": lambda s: (s["internal-registry"].get("customer_relationship") or {}).get("internal_score") is not None,
    "C12": lambda s: any((r.get("due_date") or "") >= YEAR_AGO for r in _list(s["internal-registry"], "internal_payment_records")),
    "C13": lambda s: bool(_list(s["registration-validation"], "document_validations")),
    "C14": lambda s: bool(_list(s["registration-validation"], "fiscal_regularities")),
    "C15": None,
}

RULE_HOLDS = {
    "K1": (("bureau",), lambda s: any(a.get("status") == "active" and a.get("severity") in ("high", "critical")
                                      for a in _list(s["bureau"], "fraud_alerts"))),
    "K2": (("bureau",), lambda s: any(a.get("status") in ("active", "investigating") for a in _list(s["bureau"], "fraud_alerts"))),
    "K4": (("bureau",), lambda s: ((s["bureau"].get("credit_score") or {}).get("score") or 1000) < 300),
    "K5": (("bureau",), lambda s: len(active_negatives(s["bureau"])) >= 3),
    "K6": (("bureau",), lambda s: income(s["bureau"]) is not None and sum(active_negatives(s["bureau"])) > 5 * income(s["bureau"])),
    "K7": (("bureau",), lambda s: income(s["bureau"]) is not None and sum(
        d.get("current_amount") or 0 for d in _list(s["bureau"], "debts")
        if d.get("in_collection") is not False and d.get("status") not in ("paid", "settled")) > 10 * income(s["bureau"])),
    "K8": (("bureau", "open-finance"), lambda s: consent(s["open-finance"]) and recurring_income(s["open-finance"]) > 0
           and (income(s["bureau"]) or 0) > 2 * recurring_income(s["open-finance"])),
    "K9": (("internal-registry",), lambda s: any(r["days_late"] > 90 for r in unpaid_internal(s["internal-registry"]))),
    "K9-L": (("internal-registry",), lambda s: any(r["days_late"] <= 90 and (r.get("due_date") or "") >= SIX_MONTHS_AGO
                                                   for r in unpaid_internal(s["internal-registry"]))),
    "K10": (("registration-validation",), lambda s: latest_status(
        s["registration-validation"], "document_validations", "validation_date",
        lambda r: r.get("is_valid") is False or r.get("receita_federal_status") in ("canceled", "deceased"))),
    "K11": (("registration-validation",), lambda s: latest_status(
        s["registration-validation"], "document_validations", "validation_date",
        lambda r: r.get("receita_federal_status") in ("pending", "suspended"))),
    "K12": (("registration-validation",), lambda s: latest_status(
        s["registration-validation"], "compliance_checks", "check_date", lambda r: bool(r.get("on_sanctions_list")))),
    "K13": (("registration-validation",), lambda s: latest_status(
        s["registration-validation"], "compliance_checks", "check_date",
        lambda r: bool(r.get("is_pep")) or r.get("status") == "irregular")),
    "M1": (("bureau",), lambda s: (s["bureau"].get("credit_score") or {}).get("score") is None),
    "M2": (("bureau",), lambda s: income(s["bureau"]) is None),
    "M3": (("bureau",), lambda s: not (s["bureau"].get("personal_information") or {}).get("document")),
    "M4": (("bureau",), lambda s: any(r.get("status") == "active" and r.get("is_disputed") for r in _list(s["bureau"], "negative_records"))),
}


def load_trace(lines):
    calls, responses = {}, {}
    for line in lines:
        try:
            rec = json.loads(line)
        except ValueError:
            continue
        if rec.get("message") == "mcp tool call" and rec.get("transaction_id"):
            calls[rec["transaction_id"]] = rec
        elif rec.get("message") == "[MCP RESPONSE]":
            for part in (rec.get("result") or {}).get("content") or []:
                try:
                    tid = json.loads(part["text"])["metadata"]["transaction_id"]
                except (ValueError, KeyError, TypeError):
                    continue
                responses[tid] = part["text"]
    return calls, responses


def transcript_calls(lines):
    uses, out = {}, []
    for line in lines:
        msg = json.loads(line).get("message") or {}
        content = msg.get("content")
        if not isinstance(content, list):
            continue
        for c in content:
            if c.get("type") == "tool_use" and c["name"].startswith("mcp__"):
                uses[c["id"]] = c
            elif c.get("type") == "tool_result" and c.get("tool_use_id") in uses:
                use = uses[c["tool_use_id"]]
                _, server, tool = use["name"].split("__", 2)
                body = c["content"]
                text = body[0]["text"] if isinstance(body, list) else body
                out.append({"server": server, "tool": tool, "arguments": use["input"], "text": text})
    return out


def link(call, run_id, trace):
    calls, responses = trace
    try:
        payload = json.loads(call["text"])
        tid = payload["metadata"]["transaction_id"]
    except (ValueError, KeyError, TypeError):
        return None
    rec = calls.get(tid)
    if (rec is None or rec.get("source") != call["server"] or rec.get("tool") != call["tool"]
            or rec.get("arguments") != call["arguments"] or rec.get("run_id") != run_id
            or responses.get(tid) != call["text"]):
        return None
    return payload


def client_of(call, payload):
    args = call["arguments"]
    for value in (args.get("id"), args.get("customer_id"), payload.get("id"), payload.get("customer_id")):
        if value is not None:
            return value
    return None


def evaluate(session_calls, block, run_id, trace):
    client = block.get("cliente_id")
    sources, linked, on_client = {}, 0, 0
    for call in session_calls:
        payload = link(call, run_id, trace)
        if payload is None:
            continue
        linked += 1
        if client_of(call, payload) != client:
            continue
        on_client += 1
        if call["tool"] == "get_customer_by_id":
            sources.setdefault(call["server"], {}).update(payload)
    out = {"calls": len(session_calls), "calls_linked": linked, "calls_on_client": on_client,
           "rule": block.get("regra"), "misses": []}
    view = {s: sources.get(s, {}) for s in ("bureau", "open-finance", "internal-registry", "registration-validation")}
    if out["rule"]:
        servers, holds = RULE_HOLDS.get(out["rule"], ((), None))
        out["rule_linked"] = bool(holds) and all(s in sources for s in servers)
        out["rule_holds"] = out["rule_linked"] and holds(view)
        out["rule_calls"] = len(servers)
        if not out["rule_holds"]:
            out["misses"].append(out["rule"])
        return out
    crit = block.get("criterios") or {}
    missing = set(block.get("criterios_sem_dado") or [])
    out["criteria"] = len(crit)
    out["criteria_linked"] = sum(CRITERION_SOURCE.get(c) in sources for c in crit)
    out["declared_missing"] = len(missing & set(crit))
    out["present_located"] = out["missing_confirmed"] = out["empty_list_counts"] = 0
    for c in crit:
        if CRITERION_SOURCE.get(c) not in sources:
            out["misses"].append(c)
        elif HAS_DATA[c] is None and c not in missing:
            out["empty_list_counts"] += 1
        elif HAS_DATA[c] is None:
            out["misses"].append(c)
        elif c in missing and not HAS_DATA[c](view):
            out["missing_confirmed"] += 1
        elif c not in missing and HAS_DATA[c](view):
            out["present_located"] += 1
        else:
            out["misses"].append(c)
    out["rules_holding"] = sorted(r for r, (servers, holds) in RULE_HOLDS.items()
                                  if all(s in sources for s in servers) and holds(view))
    return out


def run(round_dir):
    trace_lines = []
    for path in sorted(glob.glob(os.path.join(round_dir, "logs", "*.jsonl"))):
        with open(path) as f:
            trace_lines += f.readlines()
    trace = load_trace(trace_lines)
    rows = []
    with open(os.path.join(round_dir, "results.jsonl")) as f:
        for line in f:
            r = json.loads(line)
            name = f"{r['scenario']:02d}" + (f"-{r['run_id'].split('-', 1)[1]}" if r.get("run_id") else "")
            with open(os.path.join(round_dir, "transcripts", f"{name}.jsonl")) as t:
                row = evaluate(transcript_calls(t), r.get("block") or {}, r.get("run_id"), trace)
            rows.append(dict(row, session=name))
    return rows


def summarize(rows):
    scored = [r for r in rows if not r["rule"] and r.get("criteria")]
    closed = [r for r in rows if r["rule"]]
    total = lambda rs, k: sum(r[k] for r in rs)
    calls = [r["calls"] for r in rows]
    return {
        "sessions": len(rows),
        "calls": sum(calls),
        "calls_linked": total(rows, "calls_linked"),
        "calls_on_client": total(rows, "calls_on_client"),
        "calls_per_session": {"mean": statistics.mean(calls), "min": min(calls), "max": max(calls)} if rows else None,
        "scored": len(scored),
        "criteria": total(scored, "criteria"),
        "criteria_linked": total(scored, "criteria_linked"),
        "reported_present": total(scored, "criteria") - total(scored, "declared_missing"),
        "present_located": total(scored, "present_located"),
        "empty_list_counts": total(scored, "empty_list_counts"),
        "declared_missing": total(scored, "declared_missing"),
        "missing_confirmed": total(scored, "missing_confirmed"),
        "scored_with_rule_holding": [(r["session"], r["rules_holding"]) for r in scored if r["rules_holding"]],
        "closed": len(closed),
        "rule_linked": sum(r["rule_linked"] for r in closed),
        "rule_holds": sum(r["rule_holds"] for r in closed),
        "rule_calls": sorted({r["rule_calls"] for r in closed}),
        "misses": [(r["session"], r["misses"]) for r in rows if r["misses"]],
    }


def selftest():
    bureau = {"credit_score": {"score": 700}, "financial_profile": {"declared_monthly_income": 5000,
              "debt_to_income_ratio": None, "credit_utilization": 0.2},
              "payment_histories": [{"due_date": "2026-07-01", "status": "on_time"}],
              "personal_information": {"document": "123"}}
    of = {"data_sharing_consents": [{"status": "revoked", "expires_at": "2027-01-01"}],
          "cash_flow_analyses": [{"analysis_date": "2026-08-20", "net_cash_flow": 900, "negative_balance_days": 0}]}
    ir = {"customer_relationship": {"relationship_months": 30, "internal_score": 650},
          "internal_payment_records": [{"due_date": "2026-05-05", "payment_date": None, "days_late": 40}]}
    rv = {"document_validations": [{"validation_date": "2026-08-10", "is_valid": True, "receita_federal_status": "regular"}],
          "fiscal_regularities": [{"check_date": "2026-08-10"}],
          "compliance_checks": [{"check_date": "2026-08-10", "is_pep": True, "on_sanctions_list": False}]}
    payloads = {"bureau": bureau, "open-finance": of, "internal-registry": ir, "registration-validation": rv}
    trace_lines, transcript = ["not json"], []
    for i, (src, body) in enumerate(payloads.items()):
        text = json.dumps({"metadata": {"source": src, "transaction_id": f"t{i}"}, "id": 9, **body})
        trace_lines.append(json.dumps({"message": "mcp tool call", "source": src, "tool": "get_customer_by_id",
                                       "transaction_id": f"t{i}", "arguments": {"id": 9}, "run_id": "9-1"}))
        trace_lines.append(json.dumps({"message": "[MCP RESPONSE]", "result": {"content": [{"type": "text", "text": text}]}}))
        transcript.append(json.dumps({"message": {"content": [
            {"type": "tool_use", "id": f"u{i}", "name": f"mcp__{src}__get_customer_by_id", "input": {"id": 9}}]}}))
        transcript.append(json.dumps({"message": {"content": [
            {"type": "tool_result", "tool_use_id": f"u{i}", "content": [{"type": "text", "text": text}]}]}}))
    trace = load_trace(trace_lines)
    calls = transcript_calls(transcript)
    assert len(calls) == 4

    points = {f"C{i}": 1 for i in range(1, 16)}
    honest = evaluate(calls, {"cliente_id": 9, "criterios": points, "criterios_sem_dado": ["C2", "C7", "C8", "C9"]}, "9-1", trace)
    assert (honest["criteria_linked"], honest["present_located"], honest["missing_confirmed"], honest["empty_list_counts"]) == (15, 8, 4, 3), honest
    assert honest["rules_holding"] == ["K13", "K9-L"], honest
    wrong = evaluate(calls, {"cliente_id": 9, "criterios": points, "criterios_sem_dado": ["C1"]}, "9-1", trace)
    assert sorted(wrong["misses"]) == ["C1", "C2", "C7", "C8", "C9"], wrong
    assert "C4" in evaluate(calls, {"cliente_id": 9, "criterios": points, "criterios_sem_dado": ["C4"]}, "9-1", trace)["misses"]
    assert evaluate(calls, {"cliente_id": 8, "criterios": points}, "9-1", trace)["criteria_linked"] == 0
    assert evaluate(calls, {"cliente_id": 9, "criterios": points}, "9-2", trace)["calls_linked"] == 0

    assert evaluate(calls, {"cliente_id": 9, "regra": "K13"}, "9-1", trace)["rule_holds"] is True
    for false_rule in ("K1", "K4", "K8", "K9", "K10", "M1", "M2", "M3"):
        r = evaluate(calls, {"cliente_id": 9, "regra": false_rule}, "9-1", trace)
        assert r["rule_linked"] and not r["rule_holds"], (false_rule, r)
    assert evaluate(calls, {"cliente_id": 9, "regra": "K3"}, "9-1", trace)["rule_linked"] is False

    tampered = [dict(c) for c in calls]
    tampered[0]["text"] = tampered[0]["text"].replace("700", "100")
    t = evaluate(tampered, {"cliente_id": 9, "regra": "K4"}, "9-1", trace)
    assert t["calls_linked"] == 3 and not t["rule_holds"], t

    s = summarize([dict(honest, session="09-1"), dict(evaluate(calls, {"cliente_id": 9, "regra": "K13"}, "9-1", trace), session="09-2")])
    assert (s["calls"], s["calls_linked"], s["criteria"], s["present_located"], s["declared_missing"],
            s["missing_confirmed"], s["closed"], s["rule_holds"]) == (8, 8, 15, 8, 4, 4, 1, 1), s
    print("selftest ok")
    return 0


def main(argv=None):
    p = argparse.ArgumentParser(description="Link each criterion and rule in a round to the server-logged call that supplied it.")
    p.add_argument("round_dir", nargs="?")
    p.add_argument("--selftest", action="store_true")
    args = p.parse_args(argv)
    if args.selftest:
        return selftest()
    if not args.round_dir:
        p.error("round_dir is required")
    print(json.dumps(summarize(run(args.round_dir)), indent=2))
    return 0


if __name__ == "__main__":
    sys.exit(main())
