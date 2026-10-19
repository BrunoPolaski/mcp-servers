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

| Metrica | Acertos | Taxa | IC 95% (Wilson) |
|---|---|---|---|
| Conformidade da decisao | 4/4 | 100.0% | 51.0% a 100.0% |
| Concordancia na regra de encerramento | 4/4 | 100.0% | 51.0% a 100.0% |
| Pontuacao total exata (casos pontuados) | 4/4 | 100.0% | 51.0% a 100.0% |
| Criterios exatos (casos pontuados x 15) | 60/60 | 100.0% | - |
| Registro de dados ausentes identico (casos pontuados) | 4/4 | 100.0% | 51.0% a 100.0% |
| Inferencia sobre dado ausente (casos com dado ausente) | 0/2 | 0.0% | 0.0% a 65.8% |

> O intervalo trata cada sessao como independente; sessoes do mesmo cenario nao sao, entao ele e otimista. Criterios da mesma sessao sao ainda menos independentes e ficam sem intervalo.

| # | Sessoes | Decisao | Regra | Pontos | Inferencia | Obtido |
|---|---|---|---|---|---|---|
| 1 | 2 | 2 | 2 | 2 | 0 | APROVADO (2) |
| 18 | 2 | 2 | 2 | 2 | 0 | APROVADO_COM_RESSALVAS (2) |

