# Política de Crédito — Agente de Análise (LLM)

> **Versão:** 1.1 · **Vigência:** 2026-09-01 · **Moeda:** BRL
>
> Esta política parametriza a análise de crédito executada pelo agente (LLM)
> sobre os dados retornados por dois servidores MCP: o **Birô de Crédito**
> (`get_all_customers`, `get_customer_by_id`, `get_customer_by_document`) e o
> **Open Finance** (as mesmas três ferramentas consolidadas, mais
> `get_bank_statements`, `get_cash_flow_analysis`,
> `get_recurring_transactions` e `get_data_sharing_consents`).
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
2-A. Para analisar um cliente, recupere o cadastro do **birô** e os dados de
   **Open Finance**. As duas fontes expõem ferramentas de mesmo nome; identifique
   a origem pelo conector. Registre em `fontes_consultadas` quais servidores
   foram efetivamente consultados.
2-B. A ausência de consentimento ativo de compartilhamento **não** impede a
   análise. Os critérios de Open Finance passam a ser pontuados pelas regras de
   dado ausente da §8.
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
| K8 | Divergência de renda entre fontes | Havendo consentimento com `status = "granted"` e não expirado **e** ao menos uma `recurring_transaction` com `transaction_type = "income"` e `is_active = true`, se a renda considerada (birô) > **2 ×** a soma dos `amount` dessas receitas | `ANALISE_MANUAL` |

A dupla condição de K8 é necessária: sem consentimento ativo ou sem receita
recorrente detectada, a soma é zero e a regra dispararia para todo cliente sem
dados de Open Finance. Nesses casos, a ausência é tratada por C8, não por K8.

## 5. Scorecard (0 a 100 pontos)

| Critério | Peso | Campo(s) | Faixas → pontos | Se nulo |
|---|---|---|---|---|
| **C1. Score de bureau** | 30 | `credit_score.score` (0–1000; se houver mais de um registro, usar o de `score_date` mais recente) | ≥ 800 → 30 · 700–799 → 26 · 600–699 → 19 · 500–599 → 13 · 400–499 → 7 · 300–399 → 3 | Não ocorre (coberto por M1) |
| **C2. Comprometimento de renda (DTI)** | 15 | `financial_profile.debt_to_income_ratio`; se nulo, calcular `total_monthly_payments ÷ renda considerada` | ≤ 0,30 → 15 · 0,31–0,40 → 11 · 0,41–0,50 → 6 · > 0,50 → 0 | Se incalculável: **0 pontos** + registrar em `dados_ausentes` |
| **C3. Histórico de pagamento (12 meses)** | 15 | `payment_histories[*]` com `payment_date` nos últimos 12 meses: `pct = on_time ÷ total` | ≥ 0,95 → 15 · 0,85–0,949 → 11 · 0,70–0,849 → 6 · < 0,70 → 0 | Sem pagamentos no período (*thin file*): **6 pontos** + registrar em `dados_ausentes` |
| **C4. Negativações ativas** | 12 | `negative_records[*]` com `status = "active"` | 0 registros → 12 · 1 registro com soma ≤ R$ 1.000 → 6 · 1 registro com soma > R$ 1.000 → 3 · 2 registros → 2 | Lista vazia = 0 registros (12 pontos); **não** é dado faltante |
| **C5. Utilização de crédito** | 4 | `financial_profile.credit_utilization` (**fração de 0 a 1**) | < 0,30 → 4 · 0,30–0,599 → 2 · 0,60–0,899 → 1 · ≥ 0,90 → 0 | **1 ponto** (conservador) + registrar em `dados_ausentes` |
| **C6. Vínculo empregatício** | 4 | `employment_records[*]` com `is_current = true` e seu `verification_status` | vínculo atual `verified` → 4 · vínculo atual não verificado → 2 · sem vínculo atual → 0 | Lista vazia = sem vínculo (0 pontos); **não** é dado faltante |
| **C7. Fluxo de caixa líquido** | 10 | `cash_flow_analyses` de `analysis_date` mais recente: `net_cash_flow ÷ renda considerada` | ≥ 0,20 → 10 · 0,10–0,199 → 7 · 0,00–0,099 → 4 · < 0 → 0 | **4 pontos** + registrar em `dados_ausentes` |
| **C8. Renda recorrente detectada** | 5 | `recurring_transactions` com `transaction_type = "income"` e `is_active = true` | soma ≥ 80% da renda considerada → 5 · soma > 0 → 3 · nenhuma → 0 | Sem consentimento ativo: **3 pontos** + registrar em `dados_ausentes` |
| **C9. Dias com saldo negativo** | 5 | `negative_balance_days` da análise mais recente | 0 → 5 · 1–5 → 3 · 6–15 → 1 · > 15 → 0 | **1 ponto** + registrar em `dados_ausentes` |

**Total máximo:** 100 pontos.

**Correção de unidade em C5.** A v1.0 descrevia `credit_utilization` como
percentual de 0 a 100, mas o campo armazena a fração (0,16 a 0,80 nos dados
sintéticos). Lida ao pé da letra, a v1.0 concedia pontuação máxima a todos os
clientes. A v1.1 declara o campo como fração e ajusta as faixas.

**Unidade de `net_cash_flow`.** O campo é mensal e vale
`average_monthly_inflow − average_monthly_outflow`.

**Origem dos dados.** C1 a C6 vêm do servidor de birô; C7 a C9, do servidor de
Open Finance.

## 6. Decisão pela pontuação

| Pontuação total | Decisão |
|---|---|
| ≥ 70 | `APROVADO` |
| 50–69 | `APROVADO_COM_RESSALVAS` |
| 35–49 | `ANALISE_MANUAL` |
| < 35 | `REPROVADO` |

**Regra de degradação por dados ausentes:** se **2 ou mais** critérios do
scorecard tiverem sido pontuados por regra de nulo (C2, C3, C5, C7, C8 ou C9),
a decisão `APROVADO` é rebaixada para `APROVADO_COM_RESSALVAS`.

## 7. Limite de crédito sugerido

Somente para decisões `APROVADO` ou `APROVADO_COM_RESSALVAS`:

```
limite_base = renda_considerada × fator × (1 − DTI_efetivo)
limite      = min(limite_base, net_cash_flow × 12)   [se net_cash_flow > 0]
limite      = limite_base                            [caso contrário]
```

- `fator` = **3,0** para `APROVADO`; **1,5** para `APROVADO_COM_RESSALVAS`;
- `DTI_efetivo` = valor usado em C2; se C2 foi pontuado por regra de nulo, usar **0,50**;
- `net_cash_flow` é o da análise de fluxo de caixa mais recente. Sem consentimento
  ativo ou com fluxo nulo ou negativo, não há teto: vale `limite_base`;
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
| Sem `payment_histories` nos últimos 12 meses | 6 pontos em C3 (*thin file*) | C3 |
| `credit_utilization` nulo | 1 ponto em C5 | C5 |
| Listas vazias de ocorrências negativas | Ausência de ocorrência (pontuação cheia no critério) | §1.5 |
| `status` nulo em `negative_record` | Tratar como **`active`** (conservador) | C4/K5/K6 |
| `status` nulo em `debt` | Tratar como **não quitada** (conservador) | K7 |
| Sem consentimento `granted` e não expirado | C7 → 4 pontos, C8 → 3 pontos, C9 → 1 ponto | C7/C8/C9 |
| Consentimento ativo, sem `cash_flow_analyses` | C7 → 4 pontos, C9 → 1 ponto | C7/C9 |
| Consentimento ativo, sem receita recorrente | C8 → 0 pontos (ausência de ocorrência, não é dado faltante) | C8 |
| `inflow_volatility` nulo | Ignorar; não integra o scorecard | — |
| Qualquer outro campo nulo não coberto acima | Ignorar o campo; **nunca** substituir por valor presumido | §1.3 |

## 9. Formato de saída obrigatório

O agente **deve** responder com um único bloco JSON, sem texto fora dele:

```json
{
  "customer_id": 1,
  "document": "12345678909",
  "politica_versao": "1.1",
  "fontes_consultadas": ["bureau-mcp", "open-finance-mcp"],
  "decisao": "APROVADO | APROVADO_COM_RESSALVAS | ANALISE_MANUAL | REPROVADO",
  "regras_eliminatorias_acionadas": ["K1"],
  "pontuacao_total": 0,
  "criterios": {
    "C1_score_bureau":        { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C2_dti":                 { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C3_historico_pagamento": { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C4_negativacoes":        { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C5_utilizacao_credito":  { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C6_vinculo_empregaticio":{ "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C7_fluxo_caixa":          { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C8_renda_recorrente":     { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C9_dias_saldo_negativo":  { "valor_observado": null, "pontos": 0, "dado_ausente": false }
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

Cliente hipotético. Birô: score 715 (C1 → 26); DTI nulo, mas
`total_monthly_payments = 1.200` e renda declarada `4.000` → DTI 0,30 (C2 → 15);
12 pagamentos no período, 11 `on_time` → 91,7 % (C3 → 11); nenhuma negativação
(C4 → 12); `credit_utilization` nulo (C5 → 1, dado ausente); vínculo atual não
verificado (C6 → 2). Open Finance: consentimento concedido, `net_cash_flow`
mensal de `600` → razão 0,15 (C7 → 7); receita recorrente ativa de `3.600`,
90 % da renda (C8 → 5); 2 dias com saldo negativo (C9 → 3).

K8 não dispara: `4.000` não excede `2 × 3.600 = 7.200`.

Total = 26 + 15 + 11 + 12 + 1 + 2 + 7 + 5 + 3 = **82** → `APROVADO`. Apenas um
critério pontuado por regra de nulo (C5), portanto sem degradação.

Limite: `limite_base = 4.000 × 3,0 × (1 − 0,30) = 8.400`; teto de fluxo de caixa
`600 × 12 = 7.200`. Como o teto é menor, o limite é **R$ 7.200,00**.
