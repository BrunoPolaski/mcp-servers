# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **668**
- Taxa de sucesso: **100.0%** (668/668)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 154 | 100.0% |
| internal-registry | 174 | 100.0% |
| open-finance | 155 | 100.0% |
| registration-validation | 185 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **150**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.45** (min 4, max 9)
- Tempo de ferramenta por analise (media): **8 ms**
- Latencia media por chamada: **1.8 ms**
- Tempo total de analise, ponta a ponta (media): **27.4 s** (mediana 27.2 s, min 13.7 s, max 41.2 s; soma 68.6 min)
- Turnos do agente por analise (media): **5.5**
- Tokens por analise (media): **54608** de entrada, **2758** de saida
- Sessoes: 150, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 150/150 | 100.0% | 97.5% a 100.0% |
| Concordancia na regra de encerramento | 150/150 | 100.0% | 97.5% a 100.0% |
| Pontuacao total exata (casos pontuados) | 90/90 | 100.0% | 95.9% a 100.0% |
| Criterios exatos (casos pontuados x 15) | 1350/1350 | 100.0% | 99.7% a 100.0% |
| Registro de dados ausentes identico (casos pontuados) | 90/90 | 100.0% | 95.9% a 100.0% |
| Inferencia sobre dado ausente (casos com dado ausente) | 0/25 | 0.0% | 0.0% a 13.3% |

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 1 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 2 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 3 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 4 | 5 | 5 | 5 | 0 | 0 | REPROVADO / K1 (5) |
| 5 | 5 | 5 | 5 | 0 | 0 | REPROVADO / K4 (5) |
| 6 | 5 | 5 | 5 | 0 | 0 | REPROVADO / K10 (5) |
| 7 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 8 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 9 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / K8 (5) |
| 10 | 5 | 5 | 5 | 0 | 0 | REPROVADO / K9 (5) |
| 11 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / K13 (5) |
| 12 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 13 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 14 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 15 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 16 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / M2 (5) |
| 17 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / M1 (5) |
| 18 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 19 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 20 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 21 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / M4 (5) |
| 22 | 5 | 5 | 5 | 0 | 0 | REPROVADO / K7 (5) |
| 23 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / M3 (5) |
| 24 | 5 | 5 | 5 | 0 | 0 | ANALISE_MANUAL / K2 (5) |
| 25 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 26 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 27 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 28 | 5 | 5 | 5 | 5 | 0 | APROVADO (5) |
| 29 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |
| 30 | 5 | 5 | 5 | 5 | 0 | APROVADO_COM_RESSALVAS (5) |

