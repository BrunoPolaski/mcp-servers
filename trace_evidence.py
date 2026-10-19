#!/usr/bin/env python3
import argparse
import glob
import json
import os
import statistics
import sys

CRITERION_SOURCE = {f"C{i}": "bureau" for i in range(1, 7)}
CRITERION_SOURCE |= {f"C{i}": "open-finance" for i in range(7, 10)}
CRITERION_SOURCE |= {f"C{i}": "internal-registry" for i in range(10, 13)}
CRITERION_SOURCE |= {f"C{i}": "registration-validation" for i in range(13, 16)}

RULE_EVIDENCE = {
    "K1": {"bureau": "fraud_alerts"},
    "K2": {"bureau": "fraud_alerts"},
    "K3": {"bureau": None},
    "K4": {"bureau": "credit_score"},
    "K5": {"bureau": "negative_records"},
    "K6": {"bureau": "negative_records"},
    "K7": {"bureau": "debts"},
    "K8": {"bureau": "financial_profile", "open-finance": "recurring_transactions"},
    "K9": {"internal-registry": "internal_payment_records"},
    "K9-L": {"internal-registry": "internal_payment_records"},
    "K10": {"registration-validation": "document_validations"},
    "K11": {"registration-validation": "document_validations"},
    "K12": {"registration-validation": "compliance_checks"},
    "K13": {"registration-validation": "compliance_checks"},
    "M1": {"bureau": None},
    "M2": {"bureau": "financial_profile"},
    "M3": {"bureau": "personal_information"},
    "M4": {"bureau": "negative_records"},
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


def evaluate(session_calls, block, run_id, trace):
    linked = {}
    for call in session_calls:
        payload = link(call, run_id, trace)
        if payload is not None and call["tool"] == "get_customer_by_id":
            linked.setdefault(call["server"], payload)
    rule = block.get("regra")
    if rule:
        needs = RULE_EVIDENCE.get(rule, {})
        ok = bool(needs) and all(s in linked and (k is None or k in linked[s]) for s, k in needs.items())
        return {"rule": rule, "rule_ok": ok, "rule_calls": len(needs),
                "criteria": 0, "criteria_ok": 0, "calls": len(session_calls)}
    crit = block.get("criterios") or {}
    return {"rule": None, "criteria": len(crit),
            "criteria_ok": sum(CRITERION_SOURCE.get(c) in linked for c in crit),
            "evidence_calls": len({CRITERION_SOURCE[c] for c in crit if c in CRITERION_SOURCE}),
            "calls": len(session_calls)}


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
                rows.append(evaluate(transcript_calls(t), r.get("block") or {}, r.get("run_id"), trace))
    return rows


def summarize(rows):
    scored = [r for r in rows if r["rule"] is None and r["criteria"]]
    closed = [r for r in rows if r["rule"]]
    return {
        "sessions": len(rows),
        "scored": len(scored),
        "criteria": sum(r["criteria"] for r in scored),
        "criteria_ok": sum(r["criteria_ok"] for r in scored),
        "evidence_calls_per_scored": statistics.mean(r["evidence_calls"] for r in scored) if scored else None,
        "calls_per_session": statistics.mean(r["calls"] for r in rows) if rows else None,
        "closed": len(closed),
        "closed_ok": sum(r["rule_ok"] for r in closed),
        "rule_calls": sorted({r["rule_calls"] for r in closed}),
    }


def selftest():
    payload = lambda src, tid, **kw: json.dumps({"metadata": {"source": src, "transaction_id": tid}, **kw})
    texts = {s: payload(s, f"t{i}") for i, s in enumerate(["bureau", "open-finance", "internal-registry"])}
    texts["registration-validation"] = payload("registration-validation", "t3", compliance_checks=[])
    trace_lines = []
    for i, (src, text) in enumerate(texts.items()):
        args = {"id": 9} if src != "internal-registry" else {"id": 8}
        trace_lines.append(json.dumps({"message": "mcp tool call", "source": src, "tool": "get_customer_by_id",
                                       "transaction_id": f"t{i}", "arguments": args, "run_id": "9-1"}))
        trace_lines.append(json.dumps({"message": "[MCP RESPONSE]", "result": {"content": [{"type": "text", "text": text}]}}))
    trace_lines.append("not json")
    trace = load_trace(trace_lines)
    transcript = []
    for i, (src, text) in enumerate(texts.items()):
        transcript.append(json.dumps({"message": {"content": [
            {"type": "tool_use", "id": f"u{i}", "name": f"mcp__{src}__get_customer_by_id", "input": {"id": 9}}]}}))
        transcript.append(json.dumps({"message": {"content": [
            {"type": "tool_result", "tool_use_id": f"u{i}", "content": [{"type": "text", "text": text}]}]}}))
    calls = transcript_calls(transcript)
    assert len(calls) == 4
    scored = evaluate(calls, {"criterios": {f"C{i}": 1 for i in range(1, 16)}}, "9-1", trace)
    assert scored["criteria_ok"] == 12 and scored["evidence_calls"] == 4, scored
    assert evaluate(calls, {"criterios": {"C1": 1}}, "9-2", trace)["criteria_ok"] == 0
    assert evaluate(calls, {"regra": "K13"}, "9-1", trace)["rule_ok"] is True
    assert evaluate(calls, {"regra": "K10"}, "9-1", trace)["rule_ok"] is False
    assert evaluate(calls, {"regra": "K9"}, "9-1", trace)["rule_ok"] is False
    tampered = [dict(c) for c in calls]
    tampered[0]["text"] = payload("bureau", "t0", extra=1)
    assert evaluate(tampered, {"criterios": {"C1": 1}}, "9-1", trace)["criteria_ok"] == 0
    s = summarize([scored, evaluate(calls, {"regra": "K13"}, "9-1", trace)])
    assert (s["criteria"], s["criteria_ok"], s["closed"], s["closed_ok"]) == (15, 12, 1, 1), s
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
