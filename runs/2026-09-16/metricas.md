# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **129**
- Taxa de sucesso: **100.0%** (129/129)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 33 | 100.0% |
| internal-registry | 32 | 100.0% |
| open-finance | 32 | 100.0% |
| registration-validation | 32 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **30**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 9

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§9 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.00**
- Tempo por analise (media): **6 ms**
- Latencia media por chamada: **2.7 ms**

> O tempo por analise cobre so a execucao das ferramentas nos servidores; o tempo de raciocinio do agente esta na transcricao do host, nao aqui.

## Aderencia a politica

- Cenarios respondidos: **30/30**
- Conformidade da decisao: **100.0%**
- Concordancia na regra de encerramento: **100.0%**
- Exatidao da pontuacao (onde aplicavel): **55.6%**

| # | Esperado | Obtido | Situacao |
|---|---|---|---|
| 1 | APROVADO | APROVADO | ok |
| 2 | APROVADO | APROVADO | ok |
| 3 | APROVADO | APROVADO | ok |
| 4 | REPROVADO / K1 | REPROVADO / K1 | ok |
| 5 | REPROVADO / K4 | REPROVADO / K4 | ok |
| 6 | REPROVADO / K10 | REPROVADO / K10 | ok |
| 7 | APROVADO | APROVADO | ok |
| 8 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 9 | ANALISE_MANUAL / K8 | ANALISE_MANUAL / K8 | ok |
| 10 | REPROVADO / K9 | REPROVADO / K9 | ok |
| 11 | ANALISE_MANUAL / K13 | ANALISE_MANUAL / K13 | ok |
| 12 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 13 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 14 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 15 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 16 | ANALISE_MANUAL / M2 | ANALISE_MANUAL / M2 | ok |
| 17 | ANALISE_MANUAL / M1 | ANALISE_MANUAL / M1 | ok |
| 18 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 19 | APROVADO | APROVADO | ok |
| 20 | APROVADO | APROVADO | ok |
| 21 | ANALISE_MANUAL / M4 | ANALISE_MANUAL / M4 | ok |
| 22 | REPROVADO / K7 | REPROVADO / K7 | ok |
| 23 | ANALISE_MANUAL / M3 | ANALISE_MANUAL / M3 | ok |
| 24 | ANALISE_MANUAL / K2 | ANALISE_MANUAL / K2 | ok |
| 25 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 26 | APROVADO | APROVADO | ok |
| 27 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 28 | APROVADO | APROVADO | ok |
| 29 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |
| 30 | APROVADO_COM_RESSALVAS | APROVADO_COM_RESSALVAS | ok |

