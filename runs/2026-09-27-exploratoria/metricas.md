# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **133**
- Taxa de sucesso: **100.0%** (133/133)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 30 | 100.0% |
| internal-registry | 35 | 100.0% |
| open-finance | 31 | 100.0% |
| registration-validation | 37 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **30**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.43** (min 4, max 9)
- Tempo de ferramenta por analise (media): **6 ms**
- Latencia media por chamada: **1.3 ms**
- Tempo total de analise, ponta a ponta (media): **26.8 s** (mediana 27.5 s, min 17.6 s, max 35.7 s; soma 13.4 min)
- Turnos do agente por analise (media): **5.4**
- Tokens por analise (media): **53594** de entrada, **2949** de saida
- Sessoes: 30, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

- Cenarios com bloco estruturado: **30/30**
- Conformidade da decisao: **93.3%**
- Concordancia na regra de encerramento: **100.0%**
- Exatidao da pontuacao total (18 casos pontuados): **44.4%**
- Exatidao por criterio (casos pontuados x 15 criterios): **95.9%**
- Registro de dados ausentes identico ao gabarito (casos pontuados): **94.4%**

| # | Esperado | Obtido | Pontos esp. | Pontos obt. | Situacao |
|---|---|---|---|---|---|
| 1 | APROVADO | APROVADO | 100 | 100 | ok |
| 2 | APROVADO | APROVADO | 100 | 100 | ok |
| 3 | APROVADO | APROVADO | 100 | 100 | ok |
| 4 | REPROVADO / K1 | REPROVADO / K1 | - | - | ok |
| 5 | REPROVADO / K4 | REPROVADO / K4 | - | - | ok |
| 6 | REPROVADO / K10 | REPROVADO / K10 | - | - | ok |
| 7 | APROVADO | APROVADO | 78 | 81 | ok |
| 8 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 69 | 69 | ok |
| 9 | ANALISE_MANUAL / K8 | ANALISE_MANUAL / K8 | - | - | ok |
| 10 | REPROVADO / K9 | REPROVADO / K9 | - | - | ok |
| 11 | ANALISE_MANUAL / K13 | ANALISE_MANUAL / K13 | - | - | ok |
| 12 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 67 | 67 | ok |
| 13 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 64 | 63 | ok |
| 14 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 67 | 70 | ok |
| 15 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 63 | 63 | ok |
| 16 | ANALISE_MANUAL / M2 | ANALISE_MANUAL / M2 | - | - | ok |
| 17 | ANALISE_MANUAL / M1 | ANALISE_MANUAL / M1 | - | - | ok |
| 18 | APROVADO_COM_RESSALVAS | APROVADO | 67 | 78 | divergente |
| 19 | APROVADO | APROVADO | 96 | 95 | ok |
| 20 | APROVADO | APROVADO | 76 | 79 | ok |
| 21 | ANALISE_MANUAL / M4 | ANALISE_MANUAL / M4 | - | - | ok |
| 22 | REPROVADO / K7 | REPROVADO / K7 | - | - | ok |
| 23 | ANALISE_MANUAL / M3 | ANALISE_MANUAL / M3 | - | - | ok |
| 24 | ANALISE_MANUAL / K2 | ANALISE_MANUAL / K2 | - | - | ok |
| 25 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 67 | 69 | ok |
| 26 | APROVADO | APROVADO | 76 | 79 | ok |
| 27 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 68 | 68 | ok |
| 28 | APROVADO | APROVADO | 95 | 94 | ok |
| 29 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 64 | 64 | ok |
| 30 | APROVADO_COM_RESSALVAS | APROVADO | 67 | 70 | divergente |

### Divergencias por criterio

| Criterio | Casos | Cenario: esperado -> obtido |
|---|---|---|
| C2 | 1 | 18: 0 -> 11 |
| C3 | 6 | 7: 8 -> 11, 14: 8 -> 11, 20: 8 -> 11, 25: 4 -> 8, 26: 8 -> 11, 30: 8 -> 11 |
| C12 | 4 | 13: 5 -> 4, 19: 5 -> 4, 25: 4 -> 2, 28: 5 -> 4 |

