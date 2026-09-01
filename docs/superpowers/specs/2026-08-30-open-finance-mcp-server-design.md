# Servidor MCP de Open Finance — Design

**Data:** 2026-08-30
**Cronograma:** `tcc/src/5_next-steps.tex`, atividade "Desenvolver servidor MCP simulando Open Finance"
**Repositório:** `open-finance/`

## 1. Objetivo

Tornar o `open-finance` funcional e integrável ao host Claude Web como
segundo conector do protótipo, expondo os dados que a metodologia atribui a essa
fonte: consentimento de compartilhamento, extratos bancários dos últimos 90 dias,
análise de fluxo de caixa e identificação de receitas recorrentes e despesas fixas.

Como consequência da integração de uma segunda fonte, a política de crédito
parametrizada ganha uma versão 1.1 que incorpora critérios de Open Finance e uma
regra eliminatória de divergência entre fontes.

## 2. Estado atual

O repositório já contém, no `open-finance`:

- Entidades do domínio: `BankAccountProfile`, `BankStatement`, `CashFlowAnalysis`,
  `RecurringTransaction`, `DataSharingConsent`
- DTOs completos em `internal/infra/controllers/dto/open_finance_dto.go`, com
  `PersonDTO` já agregando as cinco dimensões
- `gormPersonRepository` com os `Preload` corretos das associações de Open Finance
- Toda a camada compartilhada (auth, usuários, analistas, endereços, middlewares,
  logger, factory de terceiros)
- `go build ./...` passa

Lacunas que este trabalho fecha:

| # | Lacuna | Evidência |
|---|---|---|
| L1 | Migration inicial é cópia do birô | `000001_init.up.sql` cria `credit_scores`, `fraud_alerts`, `debts` e mais 13 tabelas sem entidade; não cria nenhuma tabela de Open Finance |
| L2 | Fixtures do domínio vazias | `bank_account_profiles.json`, `bank_statements.json`, `cash_flow_analyses.json`, `recurring_transactions.json`, `data_sharing_consents.json` contêm `[]` |
| L3 | Clientes sintéticos desalinhados | Os 10 do Open Finance (Ana Souza…) não correspondem aos 10 do birô (Felipe Pereira Santos…) |
| L4 | Tools genéricas | `get_person_by_id` com descrição herdada do birô; sem `PersonSummaryDTO`/`GetAllSummary`; sem tools por dimensão |
| L5 | Sem repositórios das dimensões | Só existe `PersonRepository` |
| L6 | Fora da infraestrutura | Ausente de `docker-compose.yaml`, `init.sql`, `infra/main.tf` e dos workflows |
| L7 | Sem testes | Nenhum `*_test.go` no repositório, embora a metodologia afirme que há testes unitários das funções críticas de cada servidor |

## 3. Esquema de dados

Reescrever `internal/infra/thirdparty/database/migrations/000001_init.{up,down}.sql`
em vez de empilhar uma migration corretiva — o serviço nunca foi provisionado, então
não há banco em produção para migrar.

**Remover** (herança do birô, sem entidade correspondente): `credit_scores`,
`financial_profiles`, `credit_accounts`, `credit_inquiries`, `debts`,
`payment_histories`, `negative_records`, `employment_records`,
`income_declarations`, `legal_records`, `compliance_checks`, `fraud_alerts`,
`risk_assessments`, `person_relationships`, `data_sources`, `person_data_sources`.

**Manter**: `files`, `addresses`, `personal_informations`, `person_addresses`,
`person_documents`, `admins`, `analysts`, `users`, `sessions`, `api_keys`, `tokens`.

**Criar**, conforme as entidades já existentes:

- `bank_account_profiles` — `person_id`, `profile_date` (índice único composto),
  `banking_relationships`, `account_age_average`, `has_checking_account`,
  `has_savings_account`, `has_investment_account`, `investments_value`
- `bank_statements` — `person_id`, `institution`, `institution_document`,
  `account_type`, `period_start`, `period_end`, `opening_balance`,
  `closing_balance`, `total_credits`, `total_debits`, `transaction_count`,
  `currency`
- `cash_flow_analyses` — `person_id`, `analysis_date`, `period_days`,
  `average_monthly_inflow`, `average_monthly_outflow`, `net_cash_flow`,
  `inflow_volatility`, `negative_balance_days`, `has_recurring_income`
- `recurring_transactions` — `person_id`, `transaction_type`, `category`,
  `description`, `amount`, `frequency`, `counterparty`, `first_detected_date`,
  `last_occurrence_date`, `is_active`
- `data_sharing_consents` — `person_id`, `consent_id` (único), `institution`,
  `status`, `scope` (JSON), `granted_at`, `expires_at`, `revoked_at`

**Alterar `persons`**: substituir `credit_score_id` e `financial_profile_id` por
`bank_account_profile_id`; remover `consent_status` e `consent_granted_at`, já que
o consentimento passa a viver em `data_sharing_consents` (como a entidade
`Person` já modela).

Todas as tabelas seguem a convenção do birô: `id BIGSERIAL PRIMARY KEY`,
`created_at`, `updated_at`, `deleted_at` (soft delete do Gorm) e índices nos
campos marcados com `index` nas entidades.

## 4. Camada de acesso a dados

As tools por dimensão precisam consultar cada agregado sem carregar a `Person`
inteira. Novos repositórios, seguindo o padrão de `gorm_person_repository.go`
(erros via `rest_err`, `NotFoundError` explícito, `Where` por `person_id`):

| Interface | Métodos |
|---|---|
| `BankStatementRepository` | `GetByPersonID(ctx, personID, accountType *string)` |
| `CashFlowAnalysisRepository` | `GetLatestByPersonID(ctx, personID)`, `GetByPersonID(ctx, personID, limit)` |
| `RecurringTransactionRepository` | `GetByPersonID(ctx, personID, transactionType *string, onlyActive bool)` |
| `DataSharingConsentRepository` | `GetByPersonID(ctx, personID)` |

Registrar as quatro no `RepositoryFactory`, com os respectivos getters.

**`internal/services/open_finance_service.go`** orquestra os quatro repositórios e
resolve a identificação do cliente: toda tool por dimensão aceita `customer_id`
**ou** `document`; quando recebe `document`, o serviço resolve para `person_id` via
`PersonRepository.GetByDocument`. Exatamente um dos dois deve ser informado.

**`PersonSummaryDTO`** e **`PersonService.GetAllSummary`** replicam o birô
(`bureau/internal/infra/controllers/dto/person_dto.go:268` e
`person_service.go:94`), para que a listagem devolva apenas `id`, `name` e
`document` e obrigue o agente a recuperar o cadastro completo antes de analisar.

## 5. Ferramentas MCP

Sete tools, registradas em `internal/interfaces/mcp/tools/`.

### 5.1 Consolidadas (espelham o birô)

| Tool | Entrada | Saída |
|---|---|---|
| `get_all_customers` | `limit`, `offset`, `params` | `PaginatedResponse[PersonSummaryDTO]` |
| `get_customer_by_id` | `id` | `PersonDTO` (agregado completo) |
| `get_customer_by_document` | `document` | `PersonDTO` (agregado completo) |

Os nomes são idênticos aos do birô por decisão de projeto: o host distingue as
fontes pelo conector, e a paridade mantém a simetria descrita no capítulo de
metodologia. As descrições declaram explicitamente "Open Finance" na primeira
linha, para que o modelo não confunda a origem dos dados.

### 5.2 Por dimensão

Todas aceitam `customer_id` (inteiro) **ou** `document` (string).

| Tool | Parâmetros adicionais | Saída |
|---|---|---|
| `get_bank_statements` | `account_type` (opcional: `checking`, `savings`, `payment`) | `[]BankStatementDTO` |
| `get_cash_flow_analysis` | `limit` (padrão 1; 0 devolve todo o histórico) | `[]CashFlowAnalysisDTO` |
| `get_recurring_transactions` | `transaction_type` (opcional: `income`, `expense`), `only_active` (padrão `true`) | `[]RecurringTransactionDTO` |
| `get_data_sharing_consents` | — | `[]DataSharingConsentDTO` |

Erros seguem o padrão já existente: argumento ausente ou inválido devolve erro de
validação; cliente inexistente devolve `NotFound`; falha de banco devolve
`InternalServerError` com causa registrada no log.

## 6. Dados sintéticos

### 6.1 Alinhamento com o birô

`personal_informations.json`, `persons.json`, `addresses.json`,
`person_addresses.json` e `person_documents.json` passam a replicar os 10 clientes
do birô — mesmos `id`, nomes e CPFs. Isso satisfaz a exigência de consistência
entre fontes da metodologia e permite executar análises cruzadas já nesta etapa.

Os campos exclusivos do birô (`credit_score_id`, `financial_profile_id`) não são
copiados; `persons.json` passa a referenciar `bank_account_profile_id`.

### 6.2 Perfis de Open Finance

Renda considerada = renda declarada do birô (§3 da política). Cada perfil é
construído para ser coerente com o quadro de crédito do cliente:

| # | Cliente | Perfil no birô | Fluxo (entrada/saída/líquido) | Dias neg. | Renda recorrente ativa | Consentimento |
|---|---|---|---|---|---|---|
| 1 | Felipe Pereira Santos | score 254, DTI 1,10, 0/7 pontuais | 3.550 / 4.300 / **−750** | 22 | nenhuma (salário encerrado em mar/2026) | concedido |
| 2 | Henrique Martins Barbosa | score 726, DTI 0,31, 6/6 | 7.600 / 6.100 / **+1.500** | 0 | salário 6.950 | concedido |
| 3 | Fernanda Costa Barbosa | score 448, 2 negativações | 2.380 / 2.610 / **−230** | 14 | salário 2.150 | concedido |
| 4 | Fernanda Rodrigues Pereira | score 349, alerta de fraude | 4.100 / 3.980 / **+120** | 9 | nenhuma (volatilidade 0,55) | concedido |
| 5 | Gabriela Ribeiro Barbosa | score 809, DTI 0,59 | 7.100 / 6.980 / **+120** | 7 | salário 6.400 | concedido |
| 6 | Henrique Almeida Ribeiro | renda indeterminável, score 294 | — | — | — | **revogado** |
| 7 | Igor Souza Martins | score 931, DTI 0,21 | 17.200 / 12.100 / **+5.100** | 0 | salário 15.800 | concedido |
| 8 | Lucas Martins Souza | score 461, renda 12.043, DTI 0,69 | 5.900 / 5.750 / **+150** | 11 | salário **4.200** | concedido |
| 9 | Eduardo Barbosa Almeida | score 663, DTI 0,39 | 10.800 / 9.500 / **+1.300** | 3 | salário 9.600 | concedido |
| 10 | Eduardo Ribeiro Ribeiro | score 800, DTI 0,26 | 8.300 / 6.600 / **+1.700** | 0 | salário 7.500 | concedido |

Três casos existem para exercitar o que a metodologia promete:

- **Cliente 5** — sinais contraditórios: score 809 (excelente) com fluxo de caixa
  apertado e sete dias de saldo negativo. Obriga o agente a ponderar fontes que
  discordam.
- **Cliente 6** — consentimento revogado: nenhum dado de Open Finance disponível,
  acionando as regras de dado ausente dos critérios C7, C8 e C9.
- **Cliente 8** — divergência de renda entre fontes: renda declarada de 12.043 no
  birô contra 4.200 de receita recorrente detectada. Aciona a nova regra K8.

Demais parâmetros por cliente:

- `bank_account_profiles`: um registro por cliente (exceto o 6), com
  `banking_relationships` entre 2 e 4, `account_age_average` coerente com o perfil
  (15 a 132 meses), e `investments_value` preenchido apenas para os clientes 2, 5,
  7 e 10.
- `bank_statements`: três extratos mensais consecutivos por cliente cobrindo os 90
  dias, encerrando em 2026-08-31. Clientes com três ou mais relacionamentos
  bancários recebem extratos de duas instituições.
- `cash_flow_analyses`: um registro por cliente, `analysis_date` 2026-08-31,
  `period_days` 90, valores conforme a tabela acima.
- `recurring_transactions`: uma receita (quando houver) e de duas a quatro despesas
  fixas por cliente (aluguel, financiamento, utilidades, assinaturas).
- `data_sharing_consents`: um consentimento por cliente. O do cliente 6 tem
  `status: "revoked"`, `granted_at` 2026-01-15 e `revoked_at` 2026-04-02.

## 7. Política de crédito v1.1

A v1.0 permanece intacta em `bureau/docs/politica_credito_agente.md` —
é o documento sob o qual os pré-resultados do capítulo 4 foram apurados. A v1.1
nasce em `docs/politica_credito_agente_v1.1.md`, na raiz do repositório, por passar
a reger dois servidores.

### 7.1 Instruções gerais

Acrescentar à §1: para analisar um cliente, o agente deve recuperar o cadastro do
birô **e** os dados de Open Finance. A ausência de consentimento ativo não impede a
análise; os critérios de Open Finance passam a ser pontuados pelas regras de dado
ausente.

### 7.2 Scorecard rebalanceado (total 100)

| Critério | Peso v1.0 | Peso v1.1 | Faixas v1.1 |
|---|---|---|---|
| C1 Score de birô | 35 | **30** | ≥800 → 30 · 700–799 → 26 · 600–699 → 19 · 500–599 → 13 · 400–499 → 7 · 300–399 → 3 |
| C2 Comprometimento de renda | 20 | **15** | ≤0,30 → 15 · 0,31–0,40 → 11 · 0,41–0,50 → 6 · >0,50 → 0 |
| C3 Histórico de pagamento | 20 | **15** | ≥0,95 → 15 · 0,85–0,949 → 11 · 0,70–0,849 → 6 · <0,70 → 0 |
| C4 Negativações ativas | 15 | **12** | 0 → 12 · 1 com soma ≤ R$ 1.000 → 6 · 1 com soma > R$ 1.000 → 3 · 2 → 2 |
| C5 Utilização de crédito | 5 | **4** | <0,30 → 4 · 0,30–0,599 → 2 · 0,60–0,899 → 1 · ≥0,90 → 0 |
| C6 Vínculo empregatício | 5 | **4** | vigente verificado → 4 · vigente não verificado → 2 · sem vínculo → 0 |
| **C7 Fluxo de caixa líquido** | — | **10** | `net_cash_flow ÷ renda` ≥0,20 → 10 · 0,10–0,199 → 7 · 0,00–0,099 → 4 · <0 → 0 |
| **C8 Renda recorrente detectada** | — | **5** | receita recorrente ativa ≥80% da renda → 5 · qualquer receita recorrente ativa → 3 · nenhuma → 0 |
| **C9 Dias com saldo negativo** | — | **5** | 0 → 5 · 1–5 → 3 · 6–15 → 1 · >15 → 0 |

**Correção de unidade em C5.** A v1.0 descreve `credit_utilization` como
"percentual 0–100", mas as fixtures do birô armazenam a fração (0,16 a 0,80). Lida
ao pé da letra, a v1.0 atribui a pontuação máxima a todos os clientes,
neutralizando o critério. A v1.1 declara o campo como fração de 0 a 1 e ajusta as
faixas. Isso altera o gabarito de alguns dos dez casos preliminares sob a v1.1;
os resultados do capítulo 4, apurados sob a v1.0, não são afetados.

**Fontes dos critérios de Open Finance.** C7 e C9 usam a `cash_flow_analyses` de
`analysis_date` mais recente. C8 usa `recurring_transactions` com
`transaction_type = "income"` e `is_active = true`, somando os valores mensais.

**Unidade de `net_cash_flow`.** O campo é mensal, coerente com
`average_monthly_inflow` e `average_monthly_outflow`, e vale
`average_monthly_inflow − average_monthly_outflow`. As fixtures respeitam essa
identidade. A v1.1 declara isso explicitamente para que C7 e o teto do limite (§7.4)
não fiquem ambíguos quanto ao período.

**Regras de dado ausente** (registrar em `dados_ausentes` em todos os casos):

| Situação | Tratamento |
|---|---|
| Sem consentimento com `status = "granted"` e não expirado | C7 → 4 pontos, C8 → 3 pontos, C9 → 1 ponto |
| Consentimento ativo mas sem `cash_flow_analyses` | C7 → 4 pontos, C9 → 1 ponto |
| Consentimento ativo mas sem `recurring_transactions` de receita | C8 → 0 pontos (ausência de ocorrência, **não** é dado faltante) |
| `inflow_volatility` nulo | Ignorar; não integra o scorecard |

A degradação por dados ausentes (§6 da v1.0) passa a considerar C2, C3, C5, C7,
C8 e C9: dois ou mais critérios pontuados por regra de nulo rebaixam `APROVADO`
para `APROVADO_COM_RESSALVAS`.

### 7.3 Nova regra eliminatória K8

Avaliada **após** K1–K7, como última regra do estágio:

> **K8 — Divergência de renda entre fontes.** Havendo consentimento de
> compartilhamento ativo **e** ao menos uma receita recorrente ativa no Open
> Finance, se a renda considerada (birô) for maior que **2×** a soma das receitas
> recorrentes ativas mensais, a decisão é `ANALISE_MANUAL`.

A dupla condição é necessária: sem consentimento ativo ou sem receita recorrente
detectada, a soma é zero e a regra dispararia para todos os clientes sem dados de
Open Finance. Nesses casos, a ausência é tratada por C8, não por K8.

K8 é a única regra da política que exige duas fontes simultâneas. É o argumento
central do trabalho materializado em regra de negócio.

### 7.4 Limite de crédito

A fórmula da v1.0 ganha um teto por capacidade de geração de caixa:

```
limite_base = renda_considerada × fator × (1 − DTI_efetivo)
limite      = min(limite_base, net_cash_flow × 12)          [se net_cash_flow > 0]
limite      = limite_base                                    [caso contrário]
```

Arredondamento para baixo em múltiplos de R$ 100, como na v1.0. Assim o Open
Finance influencia o valor concedido, não apenas a decisão.

### 7.5 Formato de saída

O bloco JSON da §9 ganha as chaves `C7_fluxo_caixa`, `C8_renda_recorrente` e
`C9_dias_saldo_negativo` em `criterios`, e o campo `politica_versao` passa a
`"1.1"`. Acrescentar `fontes_consultadas` (array com os servidores efetivamente
consultados), útil para a métrica de completude informacional do capítulo 3.

**Fora de escopo deste trabalho:** derivar o gabarito dos dez clientes sob a v1.1.
Isso pertence à atividade "Executar testes de integração entre os servidores e o
agente de IA" do cronograma (2ª quinzena de setembro).

## 8. Testes

Primeiros testes do repositório, cobrindo as funções críticas do servidor de Open
Finance:

- **Handlers das 7 tools** — testes de tabela: argumento ausente, tipo inválido,
  `customer_id` e `document` informados simultaneamente, nenhum dos dois informado,
  resolução `document → person_id`, propagação de `NotFound`, filtros opcionais
  (`account_type`, `transaction_type`, `only_active`, `limit`).
- **Mapeamento entidade → DTO** — `NewPersonDTO` com associações vazias, parciais e
  completas; ponteiros nulos não podem gerar panic; `ToEntity` como inverso nos
  campos preenchidos.
- **`OpenFinanceService`** — com repositórios substituídos por dublês, verificando
  a resolução de identificação e a propagação de erros.

Os repositórios Gorm ficam de fora dos testes unitários; sua verificação ocorre na
execução das fixtures contra o Postgres local.

## 9. Infraestrutura

- **`init.sql`** — acrescentar `CREATE DATABASE "open-finance";`
- **`docker-compose.yaml`** — serviço `open-finance-mcp`, build de
  `./open-finance` com `target: mcp`, porta `8082:8080`, mesmas redes e
  `depends_on` do birô
- **`infra/main.tf`** — segundo Artifact Registry (`open-finance-mcp`), segunda
  service account, segundo `google_cloud_run_v2_service` e o `iam_member` público
  correspondente. Os recursos duplicados são parametrizados por `for_each` sobre um
  mapa de serviços, evitando copiar o bloco inteiro; `var.image` passa a ser um mapa
  de imagens por serviço.
- **`infra/secrets.tf`** e **`variables.tf`** — secret `DATABASE_URL` próprio do
  Open Finance
- **`.github/workflows/deploy-open-finance.yml`** — espelha `deploy-bureau.yml`,
  com `paths` filtrando `open-finance/**` e `infra/**`

Os arquivos ficam prontos e versionados. **Nenhum deploy ou push é executado** —
isso fica a cargo do autor.

## 10. Fora de escopo

- Alterações no documento do TCC (`tcc/`). O capítulo 3 precisará da tabela de
  critérios atualizada e da menção às sete tools, e o capítulo 6 dos resultados,
  mas só depois que os testes de integração produzirem números.
- Derivação do gabarito dos dez clientes sob a v1.1 (atividade de setembro).
- Os servidores de cadastro interno e de validação cadastral, que têm as mesmas
  lacunas L1–L7 e serão tratados nas atividades seguintes do cronograma.
- Os 30 cenários definitivos, que são um conjunto distinto dos dez preliminares.

## 11. Ordem de execução

1. Migrations reescritas
2. Repositórios das quatro dimensões e registro no factory
3. `OpenFinanceService`, `PersonSummaryDTO`, `GetAllSummary`
4. As sete tools e o registro no `NewMCPServer`
5. Fixtures alinhadas e populadas
6. Política v1.1
7. Testes
8. Infraestrutura
