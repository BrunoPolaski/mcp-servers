# Servidor MCP de Registro Interno — Design

**Data:** 2026-09-01
**Cronograma:** `tcc/src/5_next-steps.tex`, atividade "Desenvolver servidor MCP simulando banco de dados interno da instituição financeira"
**Repositório:** `internal-registry/`

## 1. Objetivo

Tornar o `internal-registry` funcional e integrável ao host Claude Web como
terceiro conector do protótipo, expondo os dados que a instituição financeira
mantém internamente sobre o cliente: relacionamento (tempo de casa, segmento,
score interno, risco de churn), produtos contratados, histórico de pagamento
interno, limites pré-aprovados e declarações de renda.

Como consequência da integração de uma terceira fonte, a política de crédito
ganha uma versão 1.2 que incorpora critérios de registro interno (relacionamento,
score interno e comportamento de pagamento interno), uma regra eliminatória de
inadimplência interna e um piso de limite pela pré-aprovação vigente.

## 2. Estado atual

O repositório já contém, no `internal-registry`:

- Entidades do domínio: `CustomerRelationship`, `ContractedProduct`,
  `InternalPaymentRecord`, `PreApprovedLimit`, `IncomeDeclaration`
- DTOs completos em `internal/infra/controllers/dto/internal_registry_dto.go` e
  `income_declaration_dto.go`, com `PersonDTO` já agregando as cinco dimensões
- `gormPersonRepository` com os `Preload` corretos das associações internas
- Toda a camada compartilhada (auth, usuários, analistas, endereços, middlewares,
  logger, factory de terceiros)
- `go build ./...` passa

Lacunas que este trabalho fecha (paralelas às L1–L7 do Open Finance):

| # | Lacuna | Evidência |
|---|---|---|
| L1 | Migration inicial é cópia do birô | `000001_init.up.sql` cria `credit_scores`, `fraud_alerts`, `debts` e mais 13 tabelas sem entidade; não cria nenhuma tabela de registro interno |
| L2 | Fixtures do domínio vazias | `contracted_products.json`, `customer_relationships.json`, `internal_payment_records.json`, `pre_approved_limits.json`, `income_declarations.json` contêm `[]` |
| L3 | Clientes sintéticos desalinhados | Os 10 do registro interno (Ana Souza…) não correspondem aos 10 do birô (Felipe Pereira Santos…) |
| L4 | Tools genéricas | `get_person_by_id`, `get_person_by_document`, `get_all_persons` com descrição herdada; sem `PersonSummaryDTO`/`GetAllSummary`; sem tools por dimensão |
| L5 | Sem repositórios das dimensões | Só existe `PersonRepository` |
| L6 | Fora da infraestrutura | Ausente de `docker-compose.yaml`, `init.sql`, `infra/main.tf` e dos workflows |
| L7 | Sem testes | Nenhum `*_test.go` no repositório |

## 3. Esquema de dados

Reescrever `internal/infra/thirdparty/database/migrations/000001_init.{up,down}.sql`
em vez de empilhar uma migration corretiva — o serviço nunca foi provisionado.

**Remover** (herança do birô, sem entidade correspondente): `credit_scores`,
`financial_profiles`, `credit_accounts`, `credit_inquiries`, `debts`,
`payment_histories`, `negative_records`, `employment_records`, `legal_records`,
`compliance_checks`, `fraud_alerts`, `risk_assessments`, `person_relationships`,
`data_sources`, `person_data_sources`.

**Manter**: `files`, `addresses`, `personal_informations`, `person_addresses`,
`person_documents`, `admins`, `analysts`, `users`, `sessions`, `api_keys`,
`tokens`, `income_declarations` (esta última já tem o esquema correto da entidade
`IncomeDeclaration`).

**Criar**, conforme as entidades já existentes:

- `customer_relationships` — `person_id` (índice único), `customer_since`,
  `relationship_months` (índice), `segment`, `branch`, `is_active` (índice),
  `churn_risk`, `internal_score` (índice)
- `contracted_products` — `person_id` (índice), `product_type` (índice),
  `product_name`, `contract_number` (índice), `contracted_date`, `status`
  (índice), `balance`, `monthly_value`
- `internal_payment_records` — `person_id` (índice), `contracted_product_id`
  (índice), `reference_month` (índice), `due_date`, `payment_date`, `amount_due`,
  `amount_paid`, `status` (índice), `days_late` (índice)
- `pre_approved_limits` — `person_id` (índice), `product_type` (índice),
  `approved_amount`, `interest_rate`, `calculated_date`, `valid_until` (índice),
  `policy_version`, `is_active` (índice)

**Alterar `persons`**: substituir `credit_score_id` e `financial_profile_id` por
`customer_relationship_id`; remover quaisquer colunas herdadas do birô sem
correspondência na entidade `Person` atual.

Todas as tabelas seguem a convenção do birô: `id BIGSERIAL PRIMARY KEY`,
`created_at`, `updated_at`, `deleted_at` (soft delete do Gorm) e índices nos
campos marcados com `index` nas entidades.

## 4. Camada de acesso a dados

As tools por dimensão consultam cada agregado sem carregar a `Person` inteira.
Novos repositórios, seguindo o padrão de `gorm_person_repository.go` (erros via
`rest_err`, `NotFoundError` explícito, `Where` por `person_id`):

| Interface | Métodos |
|---|---|
| `CustomerRelationshipRepository` | `GetByPersonID(ctx, personID)` |
| `ContractedProductRepository` | `GetByPersonID(ctx, personID, productType *string, status *string)` |
| `InternalPaymentRecordRepository` | `GetByPersonID(ctx, personID, status *string, productID *uint)` |
| `PreApprovedLimitRepository` | `GetByPersonID(ctx, personID, onlyActive bool)` |
| `IncomeDeclarationRepository` | `GetByPersonID(ctx, personID, verifiedOnly bool)` |

Registrar as cinco no `RepositoryFactory`, com os respectivos getters.

**`internal/services/internal_registry_service.go`** orquestra os cinco
repositórios e resolve a identificação do cliente: toda tool por dimensão aceita
`customer_id` **ou** `document`; quando recebe `document`, o serviço resolve para
`person_id` via `PersonRepository.GetByDocument`. Exatamente um dos dois deve ser
informado.

**`PersonSummaryDTO`** e **`PersonService.GetAllSummary`** replicam o birô, para
que a listagem devolva apenas `id`, `name` e `document` e obrigue o agente a
recuperar o cadastro completo antes de analisar.

## 5. Ferramentas MCP

Oito tools, registradas em `internal/interfaces/mcp/tools/`.

### 5.1 Consolidadas (espelham o birô e o Open Finance)

| Tool | Entrada | Saída |
|---|---|---|
| `get_all_customers` | `limit`, `offset`, `params` | `PaginatedResponse[PersonSummaryDTO]` |
| `get_customer_by_id` | `id` | `PersonDTO` (agregado completo) |
| `get_customer_by_document` | `document` | `PersonDTO` (agregado completo) |

Os nomes são idênticos aos das outras fontes por decisão de projeto: o host
distingue as fontes pelo conector. As descrições declaram explicitamente
"Internal Registry" na primeira linha, para que o modelo não confunda a origem.

### 5.2 Por dimensão

Todas aceitam `customer_id` (inteiro) **ou** `document` (string).

| Tool | Parâmetros adicionais | Saída |
|---|---|---|
| `get_customer_relationship` | — | `CustomerRelationshipDTO` |
| `get_contracted_products` | `product_type` (opcional), `status` (opcional) | `[]ContractedProductDTO` |
| `get_internal_payment_records` | `status` (opcional), `product_id` (opcional) | `[]InternalPaymentRecordDTO` |
| `get_pre_approved_limits` | `only_active` (padrão `true`) | `[]PreApprovedLimitDTO` |
| `get_income_declarations` | `verified_only` (padrão `false`) | `[]IncomeDeclarationDTO` |

Erros seguem o padrão existente: argumento ausente ou inválido devolve erro de
validação; cliente inexistente devolve `NotFound`; falha de banco devolve
`InternalServerError` com causa registrada no log. `customer_id` e `document`
informados simultaneamente, ou nenhum dos dois, é erro de validação.

## 6. Dados sintéticos

### 6.1 Alinhamento com o birô

`personal_informations.json`, `persons.json`, `addresses.json`,
`person_addresses.json` e `person_documents.json` passam a replicar os 10 clientes
do birô — mesmos `id`, nomes e CPFs. Os campos exclusivos do birô
(`credit_score_id`, `financial_profile_id`) não são copiados; `persons.json`
passa a referenciar `customer_relationship_id`.

### 6.2 Perfis de registro interno

Renda considerada = renda declarada do birô (§3 da política). Cada perfil é
coerente com o quadro de crédito do cliente. Três casos existem para exercitar as
regras novas:

- **Cliente 6** (Henrique Almeida Ribeiro — já o "caixa-preta": consentimento de
  Open Finance revogado e renda indeterminável) → **não é cliente da instituição**:
  sem `customer_relationship`, sem produtos, sem pagamentos internos, sem
  pré-aprovação. Aciona as regras de dado ausente de C10, C11 e C12.
- **Cliente 1** (Felipe Pereira Santos — score 254, fluxo de caixa negativo) →
  `internal_payment_record` com `status = "missed"` e `days_late > 90` em produto
  ativo. Aciona **K9 → REPROVADO**, reforçando a reprovação já indicada pelas
  demais fontes.
- **Cliente 5** (Gabriela Ribeiro Barbosa — sinais contraditórios: score de birô
  809, mas fluxo apertado no Open Finance) → **score interno baixo** (a instituição
  a conhece melhor: relacionamento com atritos), aprofundando o tema de fontes que
  discordam.

Demais parâmetros por cliente:

- `customer_relationships`: um registro por cliente (exceto o 6), com
  `relationship_months` entre 6 e 180 coerente com o perfil, `segment` em
  {`retail`, `private`, `business`}, `internal_score` (0–1000) correlacionado —
  mas não idêntico — ao score de birô, e `churn_risk` coerente.
- `contracted_products`: de 1 a 4 produtos por cliente (exceto o 6), com
  `product_type` variado (`checking_account`, `credit_card`, `loan`, `insurance`,
  `investment`) e `status` majoritariamente `active`.
- `internal_payment_records`: ~12 meses de registros por cliente com produto,
  encerrando em 2026-08, `status` coerente com o perfil (pontualidade alta para os
  bons pagadores; um `missed` >90d para o cliente 1).
- `pre_approved_limits`: um limite pré-aprovado ativo por cliente elegível
  (`valid_until` futuro), com `approved_amount` coerente com renda e score; o
  cliente 6 não tem pré-aprovação.
- `income_declarations`: uma declaração por cliente (exceto o 6), com
  `monthly_amount` coerente com a renda declarada no birô; `verified` verdadeiro
  para a maioria.

## 7. Política de crédito v1.2

A v1.1 permanece intacta em `docs/politica_credito_agente_v1.1.md`. A v1.2 nasce
em `docs/politica_credito_agente_v1.2.md`, por passar a reger três servidores.

### 7.1 Instruções gerais

Acrescentar à §1: para analisar um cliente, o agente deve recuperar o cadastro do
birô, os dados de Open Finance **e** os dados de registro interno. A ausência de
relacionamento interno (cliente novo ou não-cliente) **não** impede a análise; os
critérios de registro interno passam a ser pontuados pelas regras de dado ausente.
Registrar em `fontes_consultadas` os servidores efetivamente consultados.

### 7.2 Scorecard rebalanceado (total 100)

O bloco de registro interno (C10–C12, 18 pontos) entra abatendo 18 pontos dos
critérios de birô e Open Finance. As bandas de decisão da §6 (≥70 `APROVADO`,
50–69 `APROVADO_COM_RESSALVAS`, 35–49 `ANALISE_MANUAL`, <35 `REPROVADO`) e a regra
de degradação permanecem, pois o total continua 100.

| Critério | Peso v1.1 | Peso v1.2 | Faixas v1.2 |
|---|---|---|---|
| C1 Score de birô | 30 | **25** | ≥800 → 25 · 700–799 → 21 · 600–699 → 16 · 500–599 → 11 · 400–499 → 6 · 300–399 → 2 |
| C2 Comprometimento de renda | 15 | **12** | ≤0,30 → 12 · 0,31–0,40 → 9 · 0,41–0,50 → 5 · >0,50 → 0 |
| C3 Histórico de pagamento | 15 | **12** | ≥0,95 → 12 · 0,85–0,949 → 9 · 0,70–0,849 → 5 · <0,70 → 0 |
| C4 Negativações ativas | 12 | **10** | 0 → 10 · 1 com soma ≤ R$ 1.000 → 5 · 1 com soma > R$ 1.000 → 2 · 2 → 1 |
| C5 Utilização de crédito | 4 | **3** | <0,30 → 3 · 0,30–0,599 → 2 · 0,60–0,899 → 1 · ≥0,90 → 0 |
| C6 Vínculo empregatício | 4 | **3** | vigente verificado → 3 · vigente não verificado → 1 · sem vínculo → 0 |
| C7 Fluxo de caixa líquido | 10 | **8** | ≥0,20 → 8 · 0,10–0,199 → 6 · 0,00–0,099 → 3 · <0 → 0 |
| C8 Renda recorrente detectada | 5 | **4** | ≥80% da renda → 4 · qualquer receita ativa → 2 · nenhuma → 0 |
| C9 Dias com saldo negativo | 5 | **5** | 0 → 5 · 1–5 → 3 · 6–15 → 1 · >15 → 0 |
| **C10 Tempo de relacionamento** | — | **6** | `relationship_months` ≥60 → 6 · 24–59 → 4 · 6–23 → 2 · <6 ou `is_active=false` → 0 |
| **C11 Score interno** | — | **6** | `internal_score` ≥800 → 6 · 600–799 → 4 · 400–599 → 2 · <400 → 0 |
| **C12 Comportamento de pagamento interno** | — | **6** | `pct_on_time` ≥0,95 → 6 · 0,85–0,949 → 4 · 0,70–0,849 → 2 · <0,70 → 0 |

Soma: 82 (C1–C9) + 18 (C10–C12) = **100**.

**Origem dos novos critérios.** C10 e C11 usam a `customer_relationships` do
cliente (`relationship_months`, `is_active`, `internal_score`). C12 usa
`internal_payment_records` dos últimos 12 meses: `pct_on_time = on_time ÷ total`,
onde `on_time` conta registros com `status = "on_time"`. `churn_risk = "high"` é
mencionado na `justificativa`, mas não altera pontuação.

**Regras de dado ausente dos novos critérios** (registrar em `dados_ausentes`):

| Situação | Tratamento |
|---|---|
| Sem `customer_relationship` (não-cliente) | C10 → 2 pontos, C11 → 2 pontos |
| `internal_score` nulo, com relacionamento existente | C11 → 2 pontos |
| Sem `internal_payment_records` nos últimos 12 meses (*thin internal file*) | C12 → 3 pontos |

A degradação por dados ausentes (§6) passa a considerar C2, C3, C5, C7, C8, C9,
C10, C11 e C12: dois ou mais critérios pontuados por regra de nulo rebaixam
`APROVADO` para `APROVADO_COM_RESSALVAS`.

### 7.3 Nova regra eliminatória K9

Avaliada **após** K1–K8, como última regra do estágio:

> **K9 — Inadimplência interna.** Havendo `internal_payment_record` com
> `status = "missed"` e `days_late > 90` vinculado a um `contracted_product` com
> `status = "active"`, a decisão é `REPROVADO`. Havendo `internal_payment_record`
> com `status = "missed"` e `days_late ≤ 90` cujo `reference_month` está nos
> últimos 6 meses (sem nenhum caso >90d), a decisão é `ANALISE_MANUAL`.

Sem registro `missed`, K9 não dispara. A inadimplência interna é um sinal forte:
o cliente deixou de pagar a própria instituição.

### 7.4 Limite de crédito — piso pela pré-aprovação

A fórmula da v1.1 (base × fator × (1 − DTI), com teto de `net_cash_flow × 12`)
ganha um piso pela pré-aprovação vigente, **somente para `APROVADO`**:

```
limite = max(limite_pos_teto, approved_amount_ativo)   [se decisão = APROVADO e houver
                                                        pre_approved_limit ativo e não expirado]
limite = limite_pos_teto                                [caso contrário]
```

Um `pre_approved_limit` é considerado ativo se `is_active = true` e
`valid_until` for futuro em relação à data da análise. Se houver mais de um do
mesmo tipo, usar o de maior `approved_amount`. Para `APROVADO_COM_RESSALVAS`,
`ANALISE_MANUAL` e `REPROVADO`, o piso é ignorado. Arredondamento para baixo em
múltiplos de R$ 100, como nas versões anteriores.

### 7.5 Fora do scorecard determinístico

`contracted_products` e `income_declarations` são expostos pelas tools e ficam
disponíveis como contexto ao agente, mas **não** são pontuados de forma
independente nem alteram a `renda considerada` (que segue a precedência da §3:
declarada → estimada → salário). Isso mantém o gabarito estável e evita
redundância com K8. `contracted_products` serve para vincular os
`internal_payment_records`; `income_declarations` pode corroborar a renda na
`justificativa`.

### 7.6 Formato de saída

O bloco JSON da §9 ganha as chaves `C10_relacionamento`, `C11_score_interno` e
`C12_pagamento_interno` em `criterios`; `politica_versao` passa a `"1.2"`; e
`fontes_consultadas` passa a admitir `internal-registry-mcp`.

**Fora de escopo deste trabalho:** derivar o gabarito dos dez clientes sob a v1.2.
Isso pertence à atividade "Executar testes de integração entre os servidores e o
agente de IA" do cronograma (2ª quinzena de setembro).

## 8. Testes

Primeiros testes do repositório, cobrindo as funções críticas:

- **Handlers das 8 tools** — testes de tabela: argumento ausente, tipo inválido,
  `customer_id` e `document` simultâneos, nenhum dos dois, resolução
  `document → person_id`, propagação de `NotFound`, filtros opcionais
  (`product_type`, `status`, `product_id`, `only_active`, `verified_only`).
- **Mapeamento entidade → DTO** — `NewPersonDTO` com associações vazias, parciais
  e completas; ponteiros nulos não podem gerar panic; `ToEntity` como inverso nos
  campos preenchidos.
- **`InternalRegistryService`** — com repositórios substituídos por dublês,
  verificando a resolução de identificação e a propagação de erros.

Os repositórios Gorm ficam de fora dos testes unitários; sua verificação ocorre na
execução das fixtures contra o Postgres local.

## 9. Infraestrutura

- **`init.sql`** — acrescentar `CREATE DATABASE "internal-registry";`
- **`docker-compose.yaml`** — serviço `internal-registry-mcp`, build de
  `./internal-registry` com `target: mcp`, porta `8083:8080`, mesmas redes e
  `depends_on` das demais fontes
- **`infra/main.tf`** — acrescentar `"internal-registry-mcp"` ao `toset`
  `local.mcp_services`, já consumido por `for_each` em todos os recursos
  (Artifact Registry, service account, Cloud Run e `iam_member` público); o mapa
  `var.images` ganha a imagem do novo serviço (`var.images[each.key]`)
- **`infra/secrets.tf`** e **`variables.tf`** — secret `DATABASE_URL` próprio do
  registro interno
- **`.github/workflows/deploy-internal-registry.yml`** — espelha
  `deploy-open-finance.yml`, com `paths` filtrando `internal-registry/**` e
  `infra/**`

Os arquivos ficam prontos e versionados. **Nenhum deploy ou push é executado.**

## 10. Fora de escopo

- Alterações no documento do TCC (`tcc/`).
- Derivação do gabarito dos dez clientes sob a v1.2 (atividade de setembro).
- O servidor de validação cadastral (`registration-validation`), que tem as mesmas
  lacunas L1–L7 e será tratado na atividade seguinte do cronograma.
- Os 30 cenários definitivos, distintos dos dez preliminares.

## 11. Ordem de execução

1. Migrations reescritas
2. Repositórios das cinco dimensões e registro no factory
3. `InternalRegistryService`, `PersonSummaryDTO`, `GetAllSummary`
4. As oito tools e o registro no `NewMCPServer`
5. Fixtures alinhadas e populadas
6. Política v1.2
7. Testes
8. Infraestrutura
