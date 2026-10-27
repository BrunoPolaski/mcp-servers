# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **710**
- Taxa de sucesso: **100.0%** (710/710)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 159 | 100.0% |
| internal-registry | 182 | 100.0% |
| open-finance | 172 | 100.0% |
| registration-validation | 197 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **156**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.55** (min 4, max 13)
- Tempo de ferramenta por analise (media): **11 ms**
- Latencia media por chamada: **2.3 ms**
- Tempo total de analise, ponta a ponta (media): **22.7 s** (mediana 21.9 s, min 11.3 s, max 40.7 s; soma 59.1 min)
- Turnos do agente por analise (media): **5.6**
- Tokens por analise (media): **55367** de entrada, **2698** de saida
- Sessoes: 156, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-sonnet-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 153/156 | 98.1% | 94.5% a 99.3% |
| Concordancia na regra de encerramento | 156/156 | 100.0% | 97.6% a 100.0% |
| Pontuacao total exata (casos pontuados) | 87/90 | 96.7% | 90.7% a 98.9% |
| Criterios exatos (casos pontuados x 15) | 1347/1350 | 99.8% | - |
| Registro de dados ausentes identico (casos pontuados) | 87/90 | 96.7% | 90.7% a 98.9% |
| Inferencia sobre dado ausente (casos com dado ausente) | 3/24 | 12.5% | 4.3% a 31.0% |

> O intervalo trata cada sessao como independente; sessoes do mesmo cenario nao sao, entao ele e otimista. Criterios da mesma sessao sao ainda menos independentes e ficam sem intervalo.

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 1 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 2 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 3 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 4 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K1 (3) |
| 5 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K4 (3) |
| 6 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K10 (3) |
| 7 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 8 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 9 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K8 (3) |
| 10 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K9 (3) |
| 11 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K13 (3) |
| 12 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 13 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 14 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 15 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 16 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M2 (3) |
| 17 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M1 (3) |
| 18 | 3 | 0 | 3 | 0 | 3 | APROVADO (3) |
| 19 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 20 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 21 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M4 (3) |
| 22 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K7 (3) |
| 23 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M3 (3) |
| 24 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K2 (3) |
| 25 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 26 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 27 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 28 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 29 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 30 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
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

