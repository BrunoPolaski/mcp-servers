# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **16**
- Taxa de sucesso: **100.0%** (16/16)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 4 | 100.0% |
| internal-registry | 4 | 100.0% |
| open-finance | 4 | 100.0% |
| registration-validation | 4 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **4**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.00** (min 4, max 4)
- Tempo de ferramenta por analise (media): **38 ms**
- Latencia media por chamada: **9.4 ms**
- Tempo total de analise, ponta a ponta (media): **29.3 s** (mediana 27.8 s, min 22.7 s, max 38.7 s; soma 2.0 min)
- Turnos do agente por analise (media): **5.0**
- Tokens por analise (media): **48904** de entrada, **3124** de saida
- Sessoes: 4, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

- Cenarios com bloco estruturado: **2/30**
- Conformidade da decisao: **6.7%**
- Concordancia na regra de encerramento: **6.7%**
- Exatidao da pontuacao total (2 casos pontuados): **100.0%**
- Exatidao por criterio (casos pontuados x 15 criterios): **100.0%**
- Registro de dados ausentes identico ao gabarito (casos pontuados): **100.0%**

| # | Esperado | Obtido | Pontos esp. | Pontos obt. | Situacao |
|---|---|---|---|---|---|
| 1 | APROVADO | APROVADO | 100 | 100 | ok |
| 2 | APROVADO | - | 100 | - | sem bloco |
| 3 | APROVADO | - | 100 | - | sem bloco |
| 4 | REPROVADO / K1 | - | - | - | sem bloco |
| 5 | REPROVADO / K4 | - | - | - | sem bloco |
| 6 | REPROVADO / K10 | - | - | - | sem bloco |
| 7 | APROVADO | - | 78 | - | sem bloco |
| 8 | APROVADO_COM_RESSALVAS | - | 65 | - | sem bloco |
| 9 | ANALISE_MANUAL / K8 | - | - | - | sem bloco |
| 10 | REPROVADO / K9 | - | - | - | sem bloco |
| 11 | ANALISE_MANUAL / K13 | - | - | - | sem bloco |
| 12 | APROVADO_COM_RESSALVAS | - | 63 | - | sem bloco |
| 13 | APROVADO_COM_RESSALVAS | - | 59 | - | sem bloco |
| 14 | APROVADO_COM_RESSALVAS | - | 67 | - | sem bloco |
| 15 | APROVADO_COM_RESSALVAS | - | 63 | - | sem bloco |
| 16 | ANALISE_MANUAL / M2 | - | - | - | sem bloco |
| 17 | ANALISE_MANUAL / M1 | - | - | - | sem bloco |
| 18 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | 63 | 63 | ok |
| 19 | APROVADO | - | 92 | - | sem bloco |
| 20 | APROVADO | - | 76 | - | sem bloco |
| 21 | ANALISE_MANUAL / M4 | - | - | - | sem bloco |
| 22 | REPROVADO / K7 | - | - | - | sem bloco |
| 23 | ANALISE_MANUAL / M3 | - | - | - | sem bloco |
| 24 | ANALISE_MANUAL / K2 | - | - | - | sem bloco |
| 25 | APROVADO_COM_RESSALVAS | - | 65 | - | sem bloco |
| 26 | APROVADO | - | 76 | - | sem bloco |
| 27 | APROVADO_COM_RESSALVAS | - | 64 | - | sem bloco |
| 28 | APROVADO | - | 91 | - | sem bloco |
| 29 | APROVADO_COM_RESSALVAS | - | 64 | - | sem bloco |
| 30 | APROVADO_COM_RESSALVAS | - | 67 | - | sem bloco |

