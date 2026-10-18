# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **134**
- Taxa de sucesso: **100.0%** (134/134)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 31 | 100.0% |
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

- Chamadas por analise (media): **4.47** (min 4, max 9)
- Tempo de ferramenta por analise (media): **7 ms**
- Latencia media por chamada: **1.5 ms**
- Tempo total de analise, ponta a ponta (media): **23.9 s** (mediana 23.4 s, min 13.0 s, max 34.8 s; soma 12.0 min)
- Turnos do agente por analise (media): **5.5**
- Tokens por analise (media): **54883** de entrada, **2661** de saida
- Sessoes: 30, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

- Cenarios com bloco estruturado: **30/30**
- Conformidade da decisao: **100.0%**
- Concordancia na regra de encerramento: **100.0%**
- Exatidao da pontuacao total (18 casos pontuados): **100.0%**
- Exatidao por criterio (casos pontuados x 15 criterios): **100.0%**
- Registro de dados ausentes identico ao gabarito (casos pontuados): **100.0%**

| # | Esperado | Obtido | Pontos esp. | Pontos obt. | Situacao |
|---|---|---|---|---|---|
| 1 | APROVADO | APROVADO | 100 | 100 | ok |
| 2 | APROVADO | APROVADO | 100 | 100 | ok |
| 3 | APROVADO | APROVADO | 100 | 100 | ok |
| 4 | REPROVADO / K1 | REPROVADO / K1 | - | - | ok |
| 5 | REPROVADO / K4 | REPROVADO / K4 | - | - | ok |
| 6 | REPROVADO / K10 | REPROVADO / K10 | - | - | ok |
| 7 | APROVADO | APROVADO | 78 | 78 | ok |
| 8 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 65 | 65 | ok |
| 9 | ANALISE_MANUAL / K8 | ANALISE_MANUAL / K8 | - | - | ok |
| 10 | REPROVADO / K9 | REPROVADO / K9 | - | - | ok |
| 11 | ANALISE_MANUAL / K13 | ANALISE_MANUAL / K13 | - | - | ok |
| 12 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 63 | 63 | ok |
| 13 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 59 | 59 | ok |
| 14 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 67 | 67 | ok |
| 15 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 63 | 63 | ok |
| 16 | ANALISE_MANUAL / M2 | ANALISE_MANUAL / M2 | - | - | ok |
| 17 | ANALISE_MANUAL / M1 | ANALISE_MANUAL / M1 | - | - | ok |
| 18 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 63 | 63 | ok |
| 19 | APROVADO | APROVADO | 92 | 92 | ok |
| 20 | APROVADO | APROVADO | 76 | 76 | ok |
| 21 | ANALISE_MANUAL / M4 | ANALISE_MANUAL / M4 | - | - | ok |
| 22 | REPROVADO / K7 | REPROVADO / K7 | - | - | ok |
| 23 | ANALISE_MANUAL / M3 | ANALISE_MANUAL / M3 | - | - | ok |
| 24 | ANALISE_MANUAL / K2 | ANALISE_MANUAL / K2 | - | - | ok |
| 25 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 65 | 65 | ok |
| 26 | APROVADO | APROVADO | 76 | 76 | ok |
| 27 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 64 | 64 | ok |
| 28 | APROVADO | APROVADO | 91 | 91 | ok |
| 29 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 64 | 64 | ok |
| 30 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 67 | 67 | ok |

