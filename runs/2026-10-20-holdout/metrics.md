# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **287**
- Taxa de sucesso: **100.0%** (287/287)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 66 | 100.0% |
| internal-registry | 72 | 100.0% |
| open-finance | 71 | 100.0% |
| registration-validation | 78 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **66**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.35** (min 4, max 8)
- Tempo de ferramenta por analise (media): **9 ms**
- Latencia media por chamada: **2.1 ms**
- Tempo total de analise, ponta a ponta (media): **27.2 s** (mediana 25.8 s, min 18.2 s, max 78.1 s; soma 29.9 min)
- Turnos do agente por analise (media): **5.3**
- Tokens por analise (media): **56334** de entrada, **2689** de saida
- Sessoes: 66, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 66/66 | 100.0% | 94.5% a 100.0% |
| Concordancia na regra de encerramento | 66/66 | 100.0% | 94.5% a 100.0% |
| Pontuacao total exata (casos pontuados) | 36/36 | 100.0% | 90.4% a 100.0% |
| Criterios exatos (casos pontuados x 15) | 540/540 | 100.0% | - |
| Registro de dados ausentes identico (casos pontuados) | 36/36 | 100.0% | 90.4% a 100.0% |
| Inferencia sobre dado ausente (casos com dado ausente) | 0/9 | 0.0% | 0.0% a 29.9% |

> O intervalo trata cada sessao como independente; sessoes do mesmo cenario nao sao, entao ele e otimista. Criterios da mesma sessao sao ainda menos independentes e ficam sem intervalo.

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 31 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 32 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 33 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K6 (3) |
| 34 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K12 (3) |
| 35 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K3 (3) |
| 36 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K5 (3) |
| 37 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K2 (3) |
| 38 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K7 (3) |
| 39 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K11 (3) |
| 40 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K9-L (3) |
| 41 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K13 (3) |
| 42 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K13 (3) |
| 43 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 44 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 45 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 46 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 47 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 48 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 49 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 50 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 51 | 3 | 3 | 3 | 3 | 0 | ANALISE_MANUAL (3) |
| 52 | 3 | 3 | 3 | 3 | 0 | REPROVADO (3) |

