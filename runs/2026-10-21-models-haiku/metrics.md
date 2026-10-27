# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **743**
- Taxa de sucesso: **100.0%** (743/743)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 157 | 100.0% |
| internal-registry | 206 | 100.0% |
| open-finance | 173 | 100.0% |
| registration-validation | 207 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **156**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.76** (min 4, max 17)
- Tempo de ferramenta por analise (media): **10 ms**
- Latencia media por chamada: **2.0 ms**
- Tempo total de analise, ponta a ponta (media): **53.3 s** (mediana 50.7 s, min 29.5 s, max 117.8 s; soma 138.5 min)
- Turnos do agente por analise (media): **5.8**
- Tokens por analise (media): **59557** de entrada, **11554** de saida
- Sessoes: 156, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-haiku-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 141/156 | 90.4% | 84.7% a 94.1% |
| Concordancia na regra de encerramento | 142/156 | 91.0% | 85.5% a 94.6% |
| Pontuacao total exata (casos pontuados) | 71/90 | 78.9% | 69.4% a 86.0% |
| Criterios exatos (casos pontuados x 15) | 1135/1350 | 84.1% | - |
| Registro de dados ausentes identico (casos pontuados) | 84/90 | 93.3% | 86.2% a 96.9% |
| Inferencia sobre dado ausente (casos com dado ausente) | 5/24 | 20.8% | 9.2% a 40.5% |

> O intervalo trata cada sessao como independente; sessoes do mesmo cenario nao sao, entao ele e otimista. Criterios da mesma sessao sao ainda menos independentes e ficam sem intervalo.

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 1 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 2 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 3 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 4 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K1 (3) |
| 5 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K4 (3) |
| 6 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K10 (3) |
| 7 | 3 | 0 | 0 | 0 | 0 | ANALISE_MANUAL / K9-L (3) |
| 8 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 9 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K8 (3) |
| 10 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K9 (3) |
| 11 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K13 (3) |
| 12 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 13 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 14 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 15 | 3 | 1 | 1 | 1 | 2 | ANALISE_MANUAL / K9-L (2), APROVADO_COM_RESSALVAS (1) |
| 16 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M2 (3) |
| 17 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M1 (3) |
| 18 | 3 | 0 | 3 | 0 | 3 | APROVADO (3) |
| 19 | 3 | 2 | 2 | 0 | 0 | APROVADO (2), ANALISE_MANUAL / K9-L (1) |
| 20 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 21 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M4 (3) |
| 22 | 3 | 3 | 3 | 0 | 0 | REPROVADO / K7 (3) |
| 23 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / M3 (3) |
| 24 | 3 | 3 | 3 | 0 | 0 | ANALISE_MANUAL / K2 (3) |
| 25 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 26 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 27 | 3 | 3 | 3 | 3 | 0 | APROVADO_COM_RESSALVAS (3) |
| 28 | 3 | 3 | 3 | 3 | 0 | APROVADO (3) |
| 29 | 3 | 2 | 2 | 2 | 0 | APROVADO_COM_RESSALVAS (2), ANALISE_MANUAL / K9-L (1) |
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
| 50 | 3 | 1 | 1 | 1 | 0 | ANALISE_MANUAL / K9-L (2), APROVADO (1) |
| 51 | 3 | 3 | 1 | 1 | 0 | ANALISE_MANUAL / K9-L (2), ANALISE_MANUAL (1) |
| 52 | 3 | 0 | 0 | 0 | 0 | ANALISE_MANUAL / K9-L (3) |

