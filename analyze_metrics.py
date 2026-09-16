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

A quarta dimensao, aderencia a politica, nao vem dos logs - vem da decisao que o
agente produziu, que so existe na transcricao do Claude Web. Passe um arquivo de
decisoes (--decisions) mapeando cada cenario a decisao/regra/pontuacao que o
agente emitiu, e o script compara com o gabarito determinístico.

Uso:
    # coleta ao vivo: junte o stdout dos quatro servidores num arquivo e passe
    python3 analyze_metrics.py logs/*.jsonl --gabarito docs/gabarito_cenarios.md
    # com aderencia:
    python3 analyze_metrics.py logs/*.jsonl --gabarito docs/gabarito_cenarios.md \
            --decisions decisoes_do_agente.json
    # ou por stdin
    cat logs.jsonl | python3 analyze_metrics.py --gabarito docs/gabarito_cenarios.md

Uma analise e o conjunto de chamadas que consultam o mesmo cliente - agrupadas
pelo `id` ou `document` passado a ferramenta. Chamadas de listagem, sem cliente,
entram como custo geral e nao contam para "fontes por analise".

Autoverificacao (nao precisa dos servidores):
    python3 analyze_metrics.py --selftest
"""

import argparse
import glob
import json
import re
import statistics
import sys

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
        "media_duracao_por_analise_ms": _media(duracao_por_analise),
        "latencia_media_ms": _media(latencias),
    }


def _media(xs):
    return statistics.mean(xs) if xs else None


# --------------------------------------------------------------------------- gabarito e aderencia
def parse_gabarito(texto):
    """Le a tabela 'Por cenario' do gabarito em markdown.

    Colunas: # | Faixa | O que exercita | Decisao | Pontos | Regra | Limite | ...
    """
    esperado = {}
    for linha in texto.splitlines():
        m = re.match(r"\|\s*(\d+)\s*\|", linha)
        if not m:
            continue
        cols = [c.strip() for c in linha.strip().strip("|").split("|")]
        if len(cols) < 7:
            continue
        cenario = int(cols[0])
        esperado[cenario] = {
            "decisao": cols[3].strip("`"),
            "pontuacao": None if cols[4] == "-" else int(cols[4]),
            "regra": None if cols[5] == "-" else cols[5],
        }
    return esperado


def compare_adherence(gabarito, decisoes):
    """Confronta as decisoes do agente com o gabarito, cenario a cenario."""
    linhas, ok_decisao, ok_regra, ok_pontos, n_pontos = [], 0, 0, 0, 0
    for cenario in sorted(gabarito):
        esp = gabarito[cenario]
        obt = decisoes.get(str(cenario)) or decisoes.get(cenario)
        if obt is None:
            linhas.append((cenario, esp, None, "sem resposta"))
            continue
        d_ok = obt.get("decisao") == esp["decisao"]
        r_ok = (obt.get("regra") or None) == esp["regra"]
        ok_decisao += d_ok
        ok_regra += r_ok
        if esp["pontuacao"] is not None:
            n_pontos += 1
            ok_pontos += 1 if obt.get("pontuacao") == esp["pontuacao"] else 0
        linhas.append((cenario, esp, obt, "ok" if d_ok and r_ok else "divergente"))
    n = len(gabarito)
    respondidos = sum(1 for _, _, o, _ in linhas if o is not None)
    return {
        "n": n,
        "respondidos": respondidos,
        "conformidade_decisao": ok_decisao / n if n else None,
        "conformidade_regra": ok_regra / n if n else None,
        "exatidao_pontuacao": (ok_pontos / n_pontos) if n_pontos else None,
        "linhas": linhas,
    }


# --------------------------------------------------------------------------- saida
def _pct(x):
    return "-" if x is None else f"{x * 100:.1f}%"


def _num(x, casas=1):
    return "-" if x is None else f"{x:.{casas}f}"


def render(ag, aderencia=None):
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
          "de dados ausentes que o agente reporta na saida (§9 da politica)."]

    L += ["", "## Eficiencia operacional", "",
          f"- Chamadas por analise (media): **{_num(ag['media_chamadas_por_analise'], 2)}**",
          f"- Tempo por analise (media): **{_num(ag['media_duracao_por_analise_ms'], 0)} ms**",
          f"- Latencia media por chamada: **{_num(ag['latencia_media_ms'], 1)} ms**",
          "",
          "> O tempo por analise cobre so a execucao das ferramentas nos servidores; "
          "o tempo de raciocinio do agente esta na transcricao do host, nao aqui."]

    L += ["", "## Aderencia a politica", ""]
    if aderencia is None:
        L += ["> Sem arquivo de decisoes. Passe --decisions com a decisao, a regra e "
              "a pontuacao que o agente emitiu por cenario para apurar esta dimensao."]
    else:
        a = aderencia
        L += [f"- Cenarios respondidos: **{a['respondidos']}/{a['n']}**",
              f"- Conformidade da decisao: **{_pct(a['conformidade_decisao'])}**",
              f"- Concordancia na regra de encerramento: **{_pct(a['conformidade_regra'])}**",
              f"- Exatidao da pontuacao (onde aplicavel): **{_pct(a['exatidao_pontuacao'])}**",
              "", "| # | Esperado | Obtido | Situacao |", "|---|---|---|---|"]
        for cenario, esp, obt, sit in a["linhas"]:
            e = f"{esp['decisao']}" + (f" / {esp['regra']}" if esp["regra"] else "")
            o = "-" if obt is None else (obt.get("decisao", "?")
                                         + (f" / {obt['regra']}" if obt.get("regra") else ""))
            L.append(f"| {cenario} | {e} | {o} | {sit} |")
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
    p.add_argument("--gabarito", help="docs/gabarito_cenarios.md, para aderencia")
    p.add_argument("--decisions", help="JSON: cenario -> {decisao, regra, pontuacao} do agente")
    p.add_argument("--since", help="descarta registros com time (ISO8601) anterior a este")
    p.add_argument("--until", help="descarta registros com time (ISO8601) posterior a este")
    p.add_argument("--selftest", action="store_true", help="roda a autoverificacao e sai")
    args = p.parse_args(argv)

    if args.selftest:
        return selftest()

    ag = aggregate(load_logs(args.logs, args.since, args.until))
    aderencia = None
    if args.gabarito and args.decisions:
        gabarito = parse_gabarito(open(args.gabarito, encoding="utf-8").read())
        decisoes = json.load(open(args.decisions, encoding="utf-8"))
        aderencia = compare_adherence(gabarito, decisoes)
    print(render(ag, aderencia))
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
    ag = aggregate(parse_log_lines(logs))
    assert ag["sucessos"] == 6, ag["sucessos"]
    assert abs(ag["taxa_sucesso"] - 6 / 7) < 1e-9
    assert ag["n_analises"] == 2, ag["n_analises"]
    assert ag["listagens"] == 1, ag["listagens"]
    assert ag["media_fontes_por_analise"] == 3.0            # (4 + 2) / 2
    assert ag["pct_com_quatro_fontes"] == 0.5               # so o cliente 1
    assert ag["media_chamadas_por_analise"] == 3.0          # (4 + 2) / 2
    assert ag["media_duracao_por_analise_ms"] == 105.0      # (140 + 70) / 2
    assert ag["falhas_por_tool"] == {"get_bank_statements": 1}

    gabarito = parse_gabarito(
        "| # | Faixa | O que | Decisao | Pontos | Regra | Limite | X |\n"
        "|---|---|---|---|---|---|---|---|\n"
        "| 1 | bom | x | `APROVADO` | 100 | - | R$ 40.000 | - |\n"
        "| 2 | ruim | y | `REPROVADO` | - | K1 | - | - |\n")
    assert gabarito[1]["decisao"] == "APROVADO" and gabarito[1]["pontuacao"] == 100
    assert gabarito[2]["regra"] == "K1" and gabarito[2]["pontuacao"] is None

    decisoes = {"1": {"decisao": "APROVADO", "regra": None, "pontuacao": 100},
                "2": {"decisao": "ANALISE_MANUAL", "regra": "K2", "pontuacao": None}}
    ad = compare_adherence(gabarito, decisoes)
    assert ad["conformidade_decisao"] == 0.5               # 1 certo de 2
    assert ad["conformidade_regra"] == 0.5                 # cenario 1: None==None
    assert ad["exatidao_pontuacao"] == 1.0                 # so o cenario 1 tem pontos, e bate

    # render nao pode estourar em nenhum dos dois modos
    assert "Viabilidade tecnica" in render(ag, None)
    assert "Conformidade da decisao" in render(ag, ad)
    print("selftest ok")
    return 0


if __name__ == "__main__":
    sys.exit(main())
