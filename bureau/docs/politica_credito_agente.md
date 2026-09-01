# Política de Crédito — Agente de Análise (LLM)

> **Versão:** 1.0 · **Vigência:** 2026-07-02 · **Moeda:** BRL
>
> Este documento parametriza a análise de crédito executada pelo agente (LLM)
> sobre os dados retornados pelas ferramentas `get_all_customers`,
> `get_customer_by_id` e `get_customer_by_document` do Bureau MCP Server.
> A política é **determinística**: dado o mesmo conjunto de dados de entrada,
> existe exatamente uma decisão correta.

---

## 1. Instruções gerais ao agente

1. Para analisar um cliente, **sempre** recupere o cadastro completo via
   `get_customer_by_id` ou `get_customer_by_document`. A listagem
   (`get_all_customers`) retorna apenas `id`, `name` e `document` e **não é
   suficiente** para análise.
2. Aplique as seções na ordem: **§3 (dados mínimos) → §4 (regras
   eliminatórias) → §5 (scorecard) → §6 (decisão) → §7 (limite)**.
3. **Nunca invente ou estime valores para campos ausentes.** Campo nulo,
   vazio ou inexistente deve ser tratado exclusivamente pela regra de nulo
   correspondente (coluna "Se nulo" de cada critério e §8).
4. `null` **não é** `0`. Ausência de informação nunca deve ser interpretada
   como valor zero, exceto quando a regra do critério disser explicitamente.
5. Listas vazias de ocorrências negativas (`negative_records`,
   `fraud_alerts`, `legal_records`, `debts`) significam **ausência de
   ocorrência** (informação positiva), e não dado faltante.
6. Todo campo ausente que influenciou a análise deve ser listado em
   `dados_ausentes` na saída (§9).

## 2. Decisões possíveis

| Decisão | Código | Significado |
|---|---|---|
| Aprovado | `APROVADO` | Crédito concedido com limite integral (§7) |
| Aprovado com ressalvas | `APROVADO_COM_RESSALVAS` | Crédito concedido com limite reduzido (§7) |
| Análise manual | `ANALISE_MANUAL` | Encaminhar a analista humano; agente não decide |
| Reprovado | `REPROVADO` | Crédito negado |

## 3. Dados mínimos (pré-condição)

Se **qualquer** condição abaixo ocorrer, a decisão é `ANALISE_MANUAL`
imediatamente (não aplicar §4–§7):

| # | Condição | Campos verificados |
|---|---|---|
| M1 | Score de crédito ausente (sem registro em `credit_score` ou `score` nulo) | `credit_score.score` |
| M2 | Renda indeterminável: `declared_monthly_income` nulo **e** `estimated_monthly_income` nulo **e** nenhum vínculo empregatício com `is_current = true` e `salary` não nulo | `financial_profile.*`, `employment_records[*]` |
| M3 | Documento (CPF) ausente ou vazio | `personal_information.document` |
| M4 | Existe `negative_record` com `is_disputed = true` e `status = "active"` | `negative_records[*]` |

**Renda considerada** (usada em §4, §5 e §7) — primeira alternativa disponível:
1. `financial_profile.declared_monthly_income`;
2. senão, `financial_profile.estimated_monthly_income`;
3. senão, `salary` do vínculo com `is_current = true` (se houver mais de um, o de maior salário).

## 4. Regras eliminatórias (knock-out)

Avaliadas em ordem. A primeira acionada encerra a análise.

| # | Regra | Condição exata | Decisão |
|---|---|---|---|
| K1 | Alerta de fraude grave | Existe `fraud_alert` com `status = "active"` e `severity ∈ {"high", "critical"}` | `REPROVADO` |
| K2 | Alerta de fraude em apuração | Existe `fraud_alert` com `status ∈ {"active", "investigating"}` (qualquer outra severidade, inclusive nula) | `ANALISE_MANUAL` |
| K3 | Falência/recuperação ativa | Existe `legal_record` com `record_type = "bankruptcy"` e `status = "active"` | `REPROVADO` |
| K4 | Score crítico | `credit_score.score < 300` | `REPROVADO` |
| K5 | Excesso de negativações | 3 ou mais `negative_records` com `status = "active"` | `REPROVADO` |
| K6 | Negativação de alto valor | Soma de `amount` dos `negative_records` com `status = "active"` > **5 × renda considerada** | `REPROVADO` |
| K7 | Dívida em cobrança de alto valor | Soma de `current_amount` das `debts` com `in_collection = true` e `status ≠ "settled"` > **10 × renda considerada** | `REPROVADO` |

## 5. Scorecard (0 a 100 pontos)

| Critério | Peso | Campo(s) | Faixas → pontos | Se nulo |
|---|---|---|---|---|
| **C1. Score de bureau** | 35 | `credit_score.score` (0–1000; se houver mais de um registro, usar o de `score_date` mais recente) | ≥ 800 → 35 · 700–799 → 30 · 600–699 → 22 · 500–599 → 15 · 400–499 → 8 · 300–399 → 3 | Não ocorre (coberto por M1) |
| **C2. Comprometimento de renda (DTI)** | 20 | `financial_profile.debt_to_income_ratio`; se nulo, calcular `total_monthly_payments ÷ renda considerada` | ≤ 0,30 → 20 · 0,31–0,40 → 14 · 0,41–0,50 → 8 · > 0,50 → 0 | Se incalculável (razão e parcelas nulas): **0 pontos** + registrar em `dados_ausentes` |
| **C3. Histórico de pagamento (12 meses)** | 20 | `payment_histories[*]` com `payment_date` nos últimos 12 meses: `pct = on_time ÷ total` (status `on_time`) | ≥ 0,95 → 20 · 0,85–0,949 → 15 · 0,70–0,849 → 8 · < 0,70 → 0 | Sem pagamentos no período (*thin file*): **8 pontos** + registrar em `dados_ausentes` |
| **C4. Negativações ativas** | 15 | `negative_records[*]` com `status = "active"` (quantidade e soma de `amount`) | 0 registros → 15 · 1 registro com soma ≤ R$ 1.000 → 8 · 1 registro com soma > R$ 1.000 → 4 · 2 registros → 2 | Lista vazia = 0 registros (15 pontos); **não** é dado faltante |
| **C5. Utilização de crédito** | 5 | `financial_profile.credit_utilization` (percentual 0–100) | < 30 → 5 · 30–59,9 → 3 · 60–89,9 → 1 · ≥ 90 → 0 | **1 ponto** (conservador) + registrar em `dados_ausentes` |
| **C6. Vínculo empregatício** | 5 | `employment_records[*]` com `is_current = true` e seu `verification_status` | vínculo atual com `verification_status = "verified"` → 5 · vínculo atual não verificado → 3 · sem vínculo atual → 0 | Lista vazia = sem vínculo (0 pontos); **não** é dado faltante |

**Total máximo:** 100 pontos. Arredondamento: nenhum (todas as faixas produzem inteiros).

**Regras de desempate de faixas:** limites são inclusivos conforme escritos
(ex.: DTI = 0,30 → 20 pontos; DTI = 0,31 → 14 pontos; score = 700 → 30 pontos).

## 6. Decisão pela pontuação

| Pontuação total | Decisão |
|---|---|
| ≥ 70 | `APROVADO` |
| 50–69 | `APROVADO_COM_RESSALVAS` |
| 35–49 | `ANALISE_MANUAL` |
| < 35 | `REPROVADO` |

**Regra de degradação por dados ausentes:** se **2 ou mais** critérios do
scorecard tiverem sido pontuados por regra de nulo (C2, C3 ou C5 com dado
faltante), a decisão `APROVADO` é rebaixada para `APROVADO_COM_RESSALVAS`.

## 7. Limite de crédito sugerido

Somente para decisões `APROVADO` ou `APROVADO_COM_RESSALVAS`:

```
limite = renda_considerada × fator × (1 − DTI_efetivo)
```

- `fator` = **3,0** para `APROVADO`; **1,5** para `APROVADO_COM_RESSALVAS`;
- `DTI_efetivo` = valor usado em C2; se C2 foi pontuado por regra de nulo, usar **0,50**;
- Arredondar o resultado para baixo em múltiplos de R$ 100;
- Para `REPROVADO` e `ANALISE_MANUAL`, `limite_sugerido` deve ser `null`.

## 8. Consolidação das regras de dado ausente

| Situação | Tratamento | Onde |
|---|---|---|
| `credit_score` inexistente ou `score` nulo | `ANALISE_MANUAL` | M1 |
| Nenhuma fonte de renda | `ANALISE_MANUAL` | M2 |
| CPF ausente | `ANALISE_MANUAL` | M3 |
| `severity` nula em alerta de fraude ativo | Tratar pela K2 (`ANALISE_MANUAL`) — nunca presumir gravidade | K2 |
| `debt_to_income_ratio` nulo | Calcular por `total_monthly_payments ÷ renda`; se impossível, 0 pontos em C2 | C2 |
| Sem `payment_histories` nos últimos 12 meses | 8 pontos em C3 (*thin file*) | C3 |
| `credit_utilization` nulo | 1 ponto em C5 | C5 |
| Listas vazias de ocorrências negativas | Ausência de ocorrência (pontuação cheia no critério) | §1.5 |
| `status` nulo em `negative_record` | Tratar como **`active`** (conservador) | C4/K5/K6 |
| `status` nulo em `debt` | Tratar como **não quitada** (conservador) | K7 |
| Qualquer outro campo nulo não coberto acima | Ignorar o campo; **nunca** substituir por valor presumido | §1.3 |

## 9. Formato de saída obrigatório

O agente **deve** responder com um único bloco JSON, sem texto fora dele:

```json
{
  "customer_id": 1,
  "document": "12345678909",
  "politica_versao": "1.0",
  "decisao": "APROVADO | APROVADO_COM_RESSALVAS | ANALISE_MANUAL | REPROVADO",
  "regras_eliminatorias_acionadas": ["K1"],
  "pontuacao_total": 0,
  "criterios": {
    "C1_score_bureau":        { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C2_dti":                 { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C3_historico_pagamento": { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C4_negativacoes":        { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C5_utilizacao_credito":  { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C6_vinculo_empregaticio":{ "valor_observado": null, "pontos": 0, "dado_ausente": false }
  },
  "renda_considerada": null,
  "fonte_renda": "declared | estimated | salary | null",
  "dados_ausentes": ["financial_profile.credit_utilization"],
  "limite_sugerido": null,
  "justificativa": "Resumo objetivo em até 3 frases citando as regras aplicadas."
}
```

Observações:
- Quando a análise termina em §3 ou §4 (dados mínimos / knock-out), preencher
  `pontuacao_total` com `null` e os critérios não avaliados com `pontos: null`.
- `regras_eliminatorias_acionadas` lista apenas a **primeira** regra acionada
  (a análise encerra nela); array vazio se nenhuma.
- `valor_observado` reporta o dado bruto usado (ex.: score 742, DTI 0.28,
  "2 negativações ativas somando R$ 3.400").

## 10. Exemplo resolvido

Cliente hipotético: score 715 (C1 → 30), DTI nulo mas
`total_monthly_payments = 1.200` e renda declarada `4.000` → DTI 0,30
(C2 → 20), 12 pagamentos no período com 11 `on_time` → 91,7 % (C3 → 15),
nenhuma negativação (C4 → 15), `credit_utilization` nulo (C5 → 1, dado
ausente), vínculo atual não verificado (C6 → 3).

Total = **84** → `APROVADO`. Apenas 1 critério pontuado por regra de nulo
(C5), portanto sem degradação. Limite = `4.000 × 3,0 × (1 − 0,30)` =
`8.400` → **R$ 8.400,00**.
