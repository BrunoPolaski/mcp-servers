#!/usr/bin/env python3
"""Consolida as metricas de avaliacao dos experimentos a partir dos logs.

Os quatro servidores MCP carimbam cada invocacao de ferramenta com um registro
JSON em stdout (ver tracing.go): mensagem "mcp tool call", com source, tool,
transaction_id, duration_ms, os argumentos da chamada e o nivel (info em caso de
sucesso, error em caso de falha). Este script le esses registros e emite, em
markdown, as tres dimensoes que a metodologia deriva dos logs:

  - viabilidade tecnica: total de chamadas e taxa de sucesso, por fonte e tool;
  - completude informacional: numero medio de fontes consultadas por analise e
    percentual de analises que tocaram as quatro fontes;
  - eficiencia operacional: tempo total, numero de chamadas e latencia media por
    analise, alem das medias gerais.

The fourth dimension, policy adherence, comes from the structured block the
agent emits at the end of each analysis (policy §8, item 6). run_evaluation.py
collects those blocks, with each session's time window and host-side stats, in
results.jsonl; pass it with --results and the JSON answer key with
--gabarito.

Uso:
    python3 analyze_metrics.py logs/*.jsonl --gabarito docs/gabarito_cenarios.json \
            --results runs/2026-09-27-final/results.jsonl
    # so os logs, sem aderencia (agrupa por cliente, ver abaixo)
    cat logs.jsonl | python3 analyze_metrics.py

With --results, an analysis is the set of log records tagged with the run_id of
one agent session; records without run_id (older rounds) fall back to the
session time window. Records matching no session are dropped. Without it, an analysis
is the set of calls about the same client, grouped by the `id` or `document`
passed to the tool; listing calls, with no client, count as general cost.

Autoverificacao (nao precisa dos servidores):
    python3 analyze_metrics.py --selftest
"""

import argparse
import glob
import json
import statistics
import sys
from datetime import datetime

TRACE_MESSAGE = "mcp tool call"
EXPECTED_SOURCES = ("bureau", "open-finance", "internal-registry", "registration-validation")


# --------------------------------------------------------------------------- logs
def parse_log_lines(lines, since=None, until=None):
    """Extrai os registros de invocacao de ferramenta de um fluxo de linhas.

    Ignora linhas que nao sao JSON e registros que nao sao "mcp tool call"
    (os middlewares de HTTP/MCP tambem logam, com outra mensagem, e sao ruido).

    since/until sao carimbos ISO8601 (o campo `time` do zap); fora da janela, o
    registro e descartado. Servem para isolar a rodada definitiva dos 30 casos
    numa exportacao maior de log, ja que datas ISO8601 ordenam lexicograficamente.
    """
    registros = []
    for linha in lines:
        linha = linha.strip()
        if not linha:
            continue
        try:
            obj = json.loads(linha)
        except (ValueError, TypeError):
            continue
        if obj.get("message") != TRACE_MESSAGE:
            continue
        t = obj.get("time")
        if since and (t is None or t < since):
            continue
        if until and (t is None or t > until):
            continue
        registros.append({
            "source": obj.get("source"),
            "tool": obj.get("tool"),
            "transaction_id": obj.get("transaction_id"),
            "duration_ms": obj.get("duration_ms"),
            "sucesso": obj.get("level") != "error",
            "cliente": client_key(obj.get("arguments")),
            "time": t,
            "run_id": obj.get("run_id"),
        })
    return registros


def client_key(arguments):
    """Identidade do cliente consultado, a partir dos argumentos da chamada.

    As ferramentas de consulta recebem `id` (inteiro) ou `document` (CPF). As de
    listagem nao recebem cliente e devolvem None.
    """
    if not isinstance(arguments, dict):
        return None
    if arguments.get("document") not in (None, ""):
        return f"doc:{arguments['document']}"
    if arguments.get("id") not in (None, ""):
        return f"id:{arguments['id']}"
    return None


def assign_sessions(registros, resultados):
    """Attributes each log record to an agent session.

    A record carrying `run_id` (the X-Run-Id header of the session, see
    run_evaluation.py) belongs to the session with that run_id, so sessions may
    overlap in time, and each repetition of a scenario is its own analysis. A
    record without it, from a run made before the tag existed, belongs to the
    untagged session whose time window contains it; those sessions ran one at a
    time. An untagged record in a tagged run (a server image built before the
    tag) matches no window and is dropped, so it shows up as missing calls
    instead of calls in the wrong session. This is sturdier than client_key: an agent that
    looks a client up by id in one server and by CPF in another still yields
    one analysis. Records that match no session (a pilot run, a manual check)
    are dropped.
    """
    windows = [(datetime.fromisoformat(r["started_at"]), datetime.fromisoformat(r["ended_at"]), r["scenario"])
               for r in resultados if not r.get("run_id")]
    by_run_id = {r["run_id"]: f'{r["scenario"]}:{r["run_id"]}' for r in resultados if r.get("run_id")}
    kept = []
    for r in registros:
        if r.get("run_id"):
            cenario = by_run_id.get(r["run_id"])
        elif r["time"]:
            t = datetime.fromisoformat(r["time"])
            cenario = next((n for start, end, n in windows if start <= t <= end), None)
        else:
            cenario = None
        if cenario is not None:
            kept.append({**r, "cliente": f"cenario:{cenario}"})
    return kept


def aggregate(registros):
    """Calcula as tres dimensoes derivaveis dos logs."""
    total = len(registros)
    sucessos = sum(1 for r in registros if r["sucesso"])

    # viabilidade tecnica: taxa de sucesso por fonte e falhas por tool
    por_fonte = {}
    falhas_por_tool = {}
    for r in registros:
        f = por_fonte.setdefault(r["source"], {"total": 0, "sucesso": 0})
        f["total"] += 1
        f["sucesso"] += 1 if r["sucesso"] else 0
        if not r["sucesso"]:
            falhas_por_tool[r["tool"]] = falhas_por_tool.get(r["tool"], 0) + 1

    # agrupa por analise (cliente)
    analises = {}
    listagens = 0
    for r in registros:
        if r["cliente"] is None:
            listagens += 1
            continue
        a = analises.setdefault(r["cliente"], {"fontes": set(), "chamadas": 0, "duracao_ms": 0})
        a["fontes"].add(r["source"])
        a["chamadas"] += 1
        a["duracao_ms"] += r["duration_ms"] or 0

    fontes_por_analise = [len(a["fontes"]) for a in analises.values()]
    com_quatro = sum(1 for n in fontes_por_analise if n >= len(EXPECTED_SOURCES))
    chamadas_por_analise = [a["chamadas"] for a in analises.values()]
    duracao_por_analise = [a["duracao_ms"] for a in analises.values()]
    latencias = [r["duration_ms"] for r in registros if r["duration_ms"] is not None]

    return {
        "total_chamadas": total,
        "sucessos": sucessos,
        "taxa_sucesso": sucessos / total if total else None,
        "por_fonte": por_fonte,
        "falhas_por_tool": falhas_por_tool,
        "n_analises": len(analises),
        "listagens": listagens,
        "media_fontes_por_analise": _media(fontes_por_analise),
        "pct_com_quatro_fontes": (com_quatro / len(analises)) if analises else None,
        "media_chamadas_por_analise": _media(chamadas_por_analise),
        "min_chamadas_por_analise": min(chamadas_por_analise, default=None),
        "max_chamadas_por_analise": max(chamadas_por_analise, default=None),
        "media_duracao_por_analise_ms": _media(duracao_por_analise),
        "latencia_media_ms": _media(latencias),
    }


def _media(xs):
    return statistics.mean(xs) if xs else None


# --------------------------------------------------------------------------- host
def host_stats(resultados):
    """End-to-end cost of each analysis as the host saw it: time, turns, tokens."""
    duracoes = [r["duration_ms"] / 1000 for r in resultados if r.get("duration_ms")]
    turnos = [r["turns"] for r in resultados if r.get("turns")]
    usos = [r["usage"] for r in resultados if r.get("usage")]
    entrada = [u["inputTokens"] + u["cacheCreationInputTokens"] + u["cacheReadInputTokens"] for u in usos]
    return {
        "sessoes": len(resultados),
        "erros": sum(1 for r in resultados if r.get("error")),
        "modelos": sorted({m for r in resultados for m in r.get("models") or []}),
        "tempo_total_s": sum(duracoes),
        "tempo_medio_s": _media(duracoes),
        "tempo_mediano_s": statistics.median(duracoes) if duracoes else None,
        "tempo_min_s": min(duracoes, default=None),
        "tempo_max_s": max(duracoes, default=None),
        "turnos_medio": _media(turnos),
        "tokens_entrada_medio": _media(entrada),
        "tokens_saida_medio": _media([u["outputTokens"] for u in usos]),
    }


# --------------------------------------------------------------------------- aderencia
def compare_adherence(gabarito, blocos):
    """Confronta os blocos estruturados do agente com o gabarito, cenario a cenario.

    gabarito and blocos are keyed by scenario number as a string. Criteria and
    the missing-data list are only compared where the answer key has a score:
    a case closed by M/K has neither.
    """
    linhas, ok_decisao, ok_regra, ok_pontos, n_pontos = [], 0, 0, 0, 0
    crit_ok, crit_n, crit_div = 0, 0, {}
    aus_ok, aus_n = 0, 0
    for cenario in sorted(gabarito, key=int):
        esp = gabarito[cenario]
        obt = blocos.get(cenario)
        if obt is None:
            linhas.append((cenario, esp, None, "sem bloco"))
            continue
        d_ok = obt.get("decisao") == esp["decisao"]
        r_ok = (obt.get("regra") or None) == esp["regra"]
        ok_decisao += d_ok
        ok_regra += r_ok
        if esp["pontuacao"] is not None:
            n_pontos += 1
            ok_pontos += obt.get("pontuacao") == esp["pontuacao"]
            obtidos = obt.get("criterios") or {}
            for criterio, pontos in esp["criterios"].items():
                crit_n += 1
                if obtidos.get(criterio) == pontos:
                    crit_ok += 1
                else:
                    crit_div.setdefault(criterio, []).append((cenario, pontos, obtidos.get(criterio)))
            aus_n += 1
            aus_ok += sorted(obt.get("criterios_sem_dado") or []) == sorted(esp["criterios_sem_dado"])
        linhas.append((cenario, esp, obt, "ok" if d_ok and r_ok else "divergente"))
    n = len(gabarito)
    return {
        "n": n,
        "respondidos": sum(1 for _, _, o, _ in linhas if o is not None),
        "conformidade_decisao": ok_decisao / n if n else None,
        "conformidade_regra": ok_regra / n if n else None,
        "exatidao_pontuacao": (ok_pontos / n_pontos) if n_pontos else None,
        "n_pontuados": n_pontos,
        "exatidao_criterios": (crit_ok / crit_n) if crit_n else None,
        "divergencias_por_criterio": crit_div,
        "registro_ausentes": (aus_ok / aus_n) if aus_n else None,
        "linhas": linhas,
    }


# --------------------------------------------------------------------------- saida
def _pct(x):
    return "-" if x is None else f"{x * 100:.1f}%"


def _num(x, casas=1):
    return "-" if x is None else f"{x:.{casas}f}"


def _criterio_key(c):
    return int(c[1:])


def render(ag, aderencia=None, host=None):
    L = ["# Metricas de avaliacao", ""]

    L += ["## Viabilidade tecnica", "",
          f"- Chamadas de ferramenta: **{ag['total_chamadas']}**",
          f"- Taxa de sucesso: **{_pct(ag['taxa_sucesso'])}** "
          f"({ag['sucessos']}/{ag['total_chamadas']})",
          "", "| Fonte | Chamadas | Sucesso |", "|---|---|---|"]
    for fonte in sorted(ag["por_fonte"]):
        f = ag["por_fonte"][fonte]
        L.append(f"| {fonte} | {f['total']} | {_pct(f['sucesso'] / f['total'])} |")
    if ag["falhas_por_tool"]:
        L += ["", "Falhas por ferramenta: "
              + ", ".join(f"{t} ({n})" for t, n in sorted(ag["falhas_por_tool"].items()))]
    else:
        L += ["", "Nenhuma falha registrada."]

    L += ["", "## Completude informacional", "",
          f"- Analises identificadas: **{ag['n_analises']}**",
          f"- Fontes consultadas por analise (media): **{_num(ag['media_fontes_por_analise'], 2)}** "
          f"de {len(EXPECTED_SOURCES)}",
          f"- Analises que consultaram as quatro fontes: **{_pct(ag['pct_com_quatro_fontes'])}**",
          f"- Chamadas de listagem (sem cliente): {ag['listagens']}",
          "",
          "> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco "
          "de dados ausentes que o agente reporta na saida (§8 da politica)."]

    L += ["", "## Eficiencia operacional", "",
          f"- Chamadas por analise (media): **{_num(ag['media_chamadas_por_analise'], 2)}** "
          f"(min {ag['min_chamadas_por_analise']}, max {ag['max_chamadas_por_analise']})",
          f"- Tempo de ferramenta por analise (media): **{_num(ag['media_duracao_por_analise_ms'], 0)} ms**",
          f"- Latencia media por chamada: **{_num(ag['latencia_media_ms'], 1)} ms**"]
    if host:
        L += [f"- Tempo total de analise, ponta a ponta (media): **{_num(host['tempo_medio_s'], 1)} s** "
              f"(mediana {_num(host['tempo_mediano_s'], 1)} s, min {_num(host['tempo_min_s'], 1)} s, "
              f"max {_num(host['tempo_max_s'], 1)} s; soma {_num(host['tempo_total_s'] / 60, 1)} min)",
              f"- Turnos do agente por analise (media): **{_num(host['turnos_medio'], 1)}**",
              f"- Tokens por analise (media): **{_num(host['tokens_entrada_medio'], 0)}** de entrada, "
              f"**{_num(host['tokens_saida_medio'], 0)}** de saida",
              f"- Sessoes: {host['sessoes']}, com erro: {host['erros']}; modelos no host: {', '.join(host['modelos']) or '-'}",
              "",
              "> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, "
              "chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores."]
    else:
        L += ["", "> O tempo por analise cobre so a execucao das ferramentas nos servidores; "
              "o tempo de raciocinio do agente esta no host, nao aqui."]

    L += ["", "## Aderencia a politica", ""]
    if aderencia is None:
        L += ["> Sem resultados. Passe --results (results.jsonl de run_evaluation.py) "
              "e --gabarito (docs/gabarito_cenarios.json) para apurar esta dimensao."]
    else:
        a = aderencia
        L += [f"- Cenarios com bloco estruturado: **{a['respondidos']}/{a['n']}**",
              f"- Conformidade da decisao: **{_pct(a['conformidade_decisao'])}**",
              f"- Concordancia na regra de encerramento: **{_pct(a['conformidade_regra'])}**",
              f"- Exatidao da pontuacao total ({a['n_pontuados']} casos pontuados): "
              f"**{_pct(a['exatidao_pontuacao'])}**",
              f"- Exatidao por criterio (casos pontuados x 15 criterios): **{_pct(a['exatidao_criterios'])}**",
              f"- Registro de dados ausentes identico ao gabarito (casos pontuados): "
              f"**{_pct(a['registro_ausentes'])}**",
              "", "| # | Esperado | Obtido | Pontos esp. | Pontos obt. | Situacao |", "|---|---|---|---|---|---|"]
        for cenario, esp, obt, sit in a["linhas"]:
            e = f"{esp['decisao']}" + (f" / {esp['regra']}" if esp["regra"] else "")
            o = "-" if obt is None else (obt.get("decisao", "?")
                                         + (f" / {obt['regra']}" if obt.get("regra") else ""))
            pe = "-" if esp["pontuacao"] is None else esp["pontuacao"]
            po = "-" if obt is None or obt.get("pontuacao") is None else obt["pontuacao"]
            L.append(f"| {cenario} | {e} | {o} | {pe} | {po} | {sit} |")
        if a["divergencias_por_criterio"]:
            L += ["", "### Divergencias por criterio", "",
                  "| Criterio | Casos | Cenario: esperado -> obtido |", "|---|---|---|"]
            for criterio in sorted(a["divergencias_por_criterio"], key=_criterio_key):
                casos = a["divergencias_por_criterio"][criterio]
                L.append(f"| {criterio} | {len(casos)} | "
                         + ", ".join(f"{c}: {e} -> {o}" for c, e, o in casos) + " |")
    L.append("")
    return "\n".join(L)


# --------------------------------------------------------------------------- cli
def load_logs(paths, since=None, until=None):
    if not paths:
        return parse_log_lines(sys.stdin, since, until)
    linhas = []
    for pattern in paths:
        for path in glob.glob(pattern):
            with open(path, encoding="utf-8") as fh:
                linhas.extend(fh.readlines())
    return parse_log_lines(linhas, since, until)


def main(argv=None):
    p = argparse.ArgumentParser(description="Consolida metricas dos logs dos servidores MCP.")
    p.add_argument("logs", nargs="*", help="arquivos de log (glob); vazio le do stdin")
    p.add_argument("--gabarito", help="docs/gabarito_cenarios.json, para aderencia")
    p.add_argument("--results", help="results.jsonl de run_evaluation.py: blocos, janelas e custo no host")
    p.add_argument("--since", help="descarta registros com time (ISO8601) anterior a este")
    p.add_argument("--until", help="descarta registros com time (ISO8601) posterior a este")
    p.add_argument("--selftest", action="store_true", help="roda a autoverificacao e sai")
    args = p.parse_args(argv)

    if args.selftest:
        return selftest()

    registros = load_logs(args.logs, args.since, args.until)
    aderencia = host = None
    if args.results:
        with open(args.results, encoding="utf-8") as fh:
            resultados = [json.loads(linha) for linha in fh if linha.strip()]
        run_ids = [r["run_id"] for r in resultados if r.get("run_id")]
        repetidos = sorted({x for x in run_ids if run_ids.count(x) > 1})
        if repetidos:
            raise SystemExit(f"run_id repetido em {args.results}: {', '.join(repetidos)}; use outro --rep ao retomar uma rodada")
        registros = assign_sessions(registros, resultados)
        host = host_stats(resultados)
        if args.gabarito:
            with open(args.gabarito, encoding="utf-8") as fh:
                gabarito = json.load(fh)
            blocos = {str(r["scenario"]): r["block"] for r in resultados if r.get("block")}
            aderencia = compare_adherence(gabarito, blocos)
    print(render(aggregate(registros), aderencia, host))
    return 0


# --------------------------------------------------------------------------- autoverificacao
def selftest():
    logs = [
        # analise do cliente id:1 - toca as quatro fontes, tudo ok
        '{"level":"info","message":"mcp tool call","source":"bureau","tool":"get_customer_by_id","transaction_id":"t1","duration_ms":30,"arguments":{"id":1}}',
        '{"level":"info","message":"mcp tool call","source":"open-finance","tool":"get_cash_flow_analysis","transaction_id":"t2","duration_ms":50,"arguments":{"id":1}}',
        '{"level":"info","message":"mcp tool call","source":"internal-registry","tool":"get_customer_relationship","transaction_id":"t3","duration_ms":20,"arguments":{"id":1}}',
        '{"level":"info","message":"mcp tool call","source":"registration-validation","tool":"get_fiscal_regularity","transaction_id":"t4","duration_ms":40,"arguments":{"id":1}}',
        # analise do cliente id:2 - so duas fontes, e uma falhou
        '{"level":"info","message":"mcp tool call","source":"bureau","tool":"get_customer_by_id","transaction_id":"t5","duration_ms":60,"arguments":{"id":2}}',
        '{"level":"error","message":"mcp tool call","source":"open-finance","tool":"get_bank_statements","transaction_id":"t6","duration_ms":10,"arguments":{"id":2},"error":"timeout"}',
        # listagem, sem cliente - nao conta como analise
        '{"level":"info","message":"mcp tool call","source":"bureau","tool":"get_all_customers","transaction_id":"t7","duration_ms":15,"arguments":{}}',
        # ruido: outra mensagem, deve ser ignorada
        '{"level":"info","message":"[MCP RESPONSE]","name":"get_customer_by_id","duration":"30ms"}',
        'linha que nao e json',
    ]
    ag = aggregate(parse_log_lines(logs))
    assert ag["total_chamadas"] == 7, ag["total_chamadas"]          # 8 linhas uteis - 1 ruido json
    assert ag["sucessos"] == 6, ag["sucessos"]

    # janela de tempo: so a rodada de 2026-09-16 conta, a de 2026-09-10 fica fora
    logs_datados = [
        '{"level":"info","message":"mcp tool call","time":"2026-09-10T08:00:00Z","source":"bureau","tool":"get_customer_by_id","transaction_id":"v","duration_ms":99,"arguments":{"id":9}}',
        '{"level":"info","message":"mcp tool call","time":"2026-09-16T10:00:00Z","source":"bureau","tool":"get_customer_by_id","transaction_id":"w","duration_ms":30,"arguments":{"id":9}}',
    ]
    janela = aggregate(parse_log_lines(logs_datados, since="2026-09-16T00:00:00Z"))
    assert janela["total_chamadas"] == 1, janela["total_chamadas"]
    assert janela["latencia_media_ms"] == 30.0, janela["latencia_media_ms"]
    assert abs(ag["taxa_sucesso"] - 6 / 7) < 1e-9
    assert ag["n_analises"] == 2, ag["n_analises"]
    assert ag["listagens"] == 1, ag["listagens"]
    assert ag["media_fontes_por_analise"] == 3.0            # (4 + 2) / 2
    assert ag["pct_com_quatro_fontes"] == 0.5               # so o cliente 1
    assert ag["media_chamadas_por_analise"] == 3.0          # (4 + 2) / 2
    assert ag["media_duracao_por_analise_ms"] == 105.0      # (140 + 70) / 2
    assert ag["falhas_por_tool"] == {"get_bank_statements": 1}

    # session windows: a CPF lookup and an id lookup in the same session are one
    # analysis; the record outside every window (a pilot) is dropped
    sessao = parse_log_lines([
        '{"level":"info","message":"mcp tool call","time":"2026-09-27T10:00:05.000-0300","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":23}}',
        '{"level":"info","message":"mcp tool call","time":"2026-09-27T10:00:09.000-0300","source":"open-finance","tool":"get_customer_by_document","duration_ms":5,"arguments":{"document":"123"}}',
        '{"level":"info","message":"mcp tool call","time":"2026-09-27T09:00:00.000-0300","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":1}}',
    ])
    resultados = [{"scenario": 23, "started_at": "2026-09-27T10:00:00-03:00", "ended_at": "2026-09-27T10:01:00-03:00",
                   "duration_ms": 60000, "turns": 8, "error": False, "models": ["claude-opus-5-5"],
                   "usage": {"inputTokens": 10, "cacheCreationInputTokens": 100,
                           "cacheReadInputTokens": 1000, "outputTokens": 50}}]
    por_sessao = aggregate(assign_sessions(sessao, resultados))
    assert por_sessao["n_analises"] == 1 and por_sessao["total_chamadas"] == 2, por_sessao
    assert por_sessao["media_fontes_por_analise"] == 2
    h = host_stats(resultados)
    assert h["tempo_medio_s"] == 60.0 and h["tokens_entrada_medio"] == 1110 and h["erros"] == 0

    sobrepostas = parse_log_lines([
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:05.000-0300","run_id":"5-1","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":5}}',
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:06.000-0300","run_id":"7-1","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":7}}',
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:07.000-0300","run_id":"5-1","source":"open-finance","tool":"get_cash_flow_analysis","duration_ms":5,"arguments":{"id":5}}',
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:08.000-0300","run_id":"9-1","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":9}}',
    ])
    janela = {"started_at": "2026-10-20T10:00:00-03:00", "ended_at": "2026-10-20T10:02:00-03:00"}
    sessoes = assign_sessions(sobrepostas, [{"scenario": 5, "run_id": "5-1", **janela},
                                            {"scenario": 7, "run_id": "7-1", **janela}])
    assert sorted((r["cliente"], r["source"]) for r in sessoes) == [
        ("cenario:5:5-1", "bureau"), ("cenario:5:5-1", "open-finance"), ("cenario:7:7-1", "bureau")], sessoes
    assert aggregate(sessoes)["n_analises"] == 2

    # repetitions of one scenario stay apart, and an untagged record in a tagged
    # run is dropped instead of landing in whichever window contains it
    repeticoes = parse_log_lines([
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:05.000-0300","run_id":"5-1","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":5}}',
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:06.000-0300","run_id":"5-2","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":5}}',
        '{"level":"info","message":"mcp tool call","time":"2026-10-20T10:00:07.000-0300","source":"bureau","tool":"get_customer_by_id","duration_ms":5,"arguments":{"id":7}}',
    ])
    sessoes = assign_sessions(repeticoes, [{"scenario": 5, "run_id": "5-1", **janela},
                                           {"scenario": 5, "run_id": "5-2", **janela}])
    assert sorted(r["cliente"] for r in sessoes) == ["cenario:5:5-1", "cenario:5:5-2"], sessoes
    assert aggregate(sessoes)["n_analises"] == 2

    todos = {f"C{i}": 1 for i in range(1, 16)}
    gabarito = {
        "1": {"decisao": "APROVADO", "regra": None, "pontuacao": 15, "criterios": todos,
              "criterios_sem_dado": ["C5"]},
        "2": {"decisao": "REPROVADO", "regra": "K1", "pontuacao": None, "criterios": {},
              "criterios_sem_dado": []},
    }
    blocos = {"1": {"decisao": "APROVADO", "regra": None, "pontuacao": 14,
                    "criterios": {**todos, "C3": 0}, "criterios_sem_dado": ["C5"]},
              "2": {"decisao": "ANALISE_MANUAL", "regra": "K2", "pontuacao": None, "criterios": {}}}
    ad = compare_adherence(gabarito, blocos)
    assert ad["conformidade_decisao"] == 0.5               # 1 certo de 2
    assert ad["conformidade_regra"] == 0.5                 # cenario 1: None==None
    assert ad["exatidao_pontuacao"] == 0.0                 # so o cenario 1 tem pontos, e erra por 1
    assert ad["exatidao_criterios"] == 14 / 15
    assert ad["divergencias_por_criterio"] == {"C3": [("1", 1, 0)]}
    assert ad["registro_ausentes"] == 1.0
    assert compare_adherence(gabarito, {})["respondidos"] == 0

    # render nao pode estourar em nenhum dos modos
    assert "Viabilidade tecnica" in render(ag, None)
    texto = render(ag, ad, h)
    assert "Conformidade da decisao" in texto and "| C3 | 1 |" in texto and "ponta a ponta" in texto
    print("selftest ok")
    return 0


if __name__ == "__main__":
    sys.exit(main())
