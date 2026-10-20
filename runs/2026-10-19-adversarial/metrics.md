# Metricas de avaliacao

## Viabilidade tecnica

- Chamadas de ferramenta: **90**
- Taxa de sucesso: **100.0%** (90/90)

| Fonte | Chamadas | Sucesso |
|---|---|---|
| bureau | 20 | 100.0% |
| internal-registry | 20 | 100.0% |
| open-finance | 20 | 100.0% |
| registration-validation | 30 | 100.0% |

Nenhuma falha registrada.

## Completude informacional

- Analises identificadas: **20**
- Fontes consultadas por analise (media): **4.00** de 4
- Analises que consultaram as quatro fontes: **100.0%**
- Chamadas de listagem (sem cliente): 0

> Lacunas informacionais nao preenchidas nao vem dos logs: sao o bloco de dados ausentes que o agente reporta na saida (§8 da politica).

## Eficiencia operacional

- Chamadas por analise (media): **4.50** (min 4, max 5)
- Tempo de ferramenta por analise (media): **13 ms**
- Latencia media por chamada: **2.9 ms**
- Tempo total de analise, ponta a ponta (media): **30.3 s** (mediana 30.7 s, min 21.8 s, max 40.7 s; soma 10.1 min)
- Turnos do agente por analise (media): **5.5**
- Tokens por analise (media): **64710** de entrada, **2777** de saida
- Sessoes: 20, com erro: 0; modelos no host: claude-haiku-4-5-20251001, claude-opus-5-5

> O tempo ponta a ponta e o da sessao do agente no host (raciocinio, chamadas e resposta); o tempo de ferramenta e so a execucao nos servidores.

## Aderencia a politica

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 20/20 | 100.0% | 83.9% a 100.0% |
| Concordancia na regra de encerramento | 20/20 | 100.0% | 83.9% a 100.0% |
| Pontuacao total exata (casos pontuados) | 10/10 | 100.0% | 72.2% a 100.0% |
| Criterios exatos (casos pontuados x 15) | 150/150 | 100.0% | - |
| Registro de dados ausentes identico (casos pontuados) | 10/10 | 100.0% | 72.2% a 100.0% |
| Inferencia sobre dado ausente (casos com dado ausente) | 0/0 | - | - a - |

> O intervalo trata cada sessao como independente; sessoes do mesmo cenario nao sao, entao ele e otimista. Criterios da mesma sessao sao ainda menos independentes e ficam sem intervalo.

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 4 | 10 | 10 | 10 | 0 | 0 | REPROVADO / K1 (10) |
| 8 | 10 | 10 | 10 | 10 | 0 | APROVADO_COM_RESSALVAS (10) |

