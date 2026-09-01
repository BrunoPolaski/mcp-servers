# Servidor MCP de Open Finance — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deixar o `open-finance` funcional como segundo conector MCP do protótipo, com esquema próprio, sete ferramentas, dados sintéticos coerentes com o birô e a política de crédito v1.1.

**Architecture:** Arquitetura em camadas já estabelecida pelo `bureau` — entidades Gorm em `internal/core/entities`, repositórios em `internal/infra/repositories` atrás de interfaces, serviços em `internal/services`, ferramentas MCP em `internal/interfaces/mcp/tools`. O servidor expõe três ferramentas consolidadas (espelhando o birô) e quatro por dimensão de dado. As dependências dos serviços e das ferramentas passam a ser declaradas como interfaces, para permitir dublês nos testes.

**Tech Stack:** Go 1.25.5, `github.com/mark3labs/mcp-go` v0.54.0, Gorm v1.31.1 sobre PostgreSQL, `github.com/BrunoPolaski/go-rest-err` para erros, `golang-migrate` para migrations, Docker Compose local, Terraform + Cloud Run em produção.

**Spec:** `docs/superpowers/specs/2026-08-30-open-finance-design.md`

## Global Constraints

- Módulo Go: `github.com/BrunoPolaski/open-finance`. Todo import interno usa esse prefixo.
- Erros de domínio sempre como `*rest_err.RestErr`: `rest_err.NewBadRequestError`, `NewNotFoundError`, `NewInternalServerError(msg).WithCause(err)`.
- Repositórios devolvem `(valor, *rest_err.RestErr)`; nunca `error` cru.
- Nomes de ferramentas: `get_all_customers`, `get_customer_by_id`, `get_customer_by_document`, `get_bank_statements`, `get_cash_flow_analysis`, `get_recurring_transactions`, `get_data_sharing_consents`.
- Toda descrição de ferramenta começa declarando "Open Finance" para o modelo não confundir a fonte.
- Datas das fixtures encerram em `2026-08-31`; janela de análise de 90 dias.
- `net_cash_flow` é mensal e vale `average_monthly_inflow − average_monthly_outflow`.
- Moeda `BRL`, CPFs sem pontuação (11 dígitos), CNPJs sem pontuação (14 dígitos).
- Nenhum `git push` e nenhum `terraform apply` são executados por este plano.
- Commits em português, prefixo `feat:`, `fix:`, `test:`, `docs:` ou `chore:`. **Nunca** adicionar rodapé de co-autoria.

---

### Task 1: Esquema de dados do Open Finance

Reescreve a migration inicial, que hoje é cópia do birô, para refletir as entidades que o servidor realmente tem.

**Files:**
- Modify: `open-finance/internal/infra/thirdparty/database/migrations/000001_init.up.sql`
- Modify: `open-finance/internal/infra/thirdparty/database/migrations/000001_init.down.sql`
- Modify: `init.sql`

**Interfaces:**
- Consumes: nada.
- Produces: tabelas `bank_account_profiles`, `bank_statements`, `cash_flow_analyses`, `recurring_transactions`, `data_sharing_consents`; coluna `persons.bank_account_profile_id`. As Tasks 2 e 7 dependem desses nomes.

- [ ] **Step 1: Subir o Postgres local e criar o banco**

Acrescentar a linha ao final de `init.sql`:

```sql
CREATE DATABASE bureau;
CREATE DATABASE "open-finance";
```

Se o volume do Postgres já existe, o `init.sql` não roda de novo. Criar o banco à mão:

```bash
docker compose up -d postgres
docker compose exec postgres psql -U "$DB_USERNAME" -d postgres -c 'CREATE DATABASE "open-finance";'
```

- [ ] **Step 2: Reescrever `000001_init.up.sql`**

Manter, **sem alteração**, os blocos `CREATE TABLE` de: `api_keys`, `files`, `addresses`, `person_addresses`, `person_documents`, `admins`, `analysts`, `users`, `sessions`.

Em `personal_informations`, remover as três colunas que pertencem ao domínio de validação cadastral e não existem na entidade: `document_validated`, `biometric_validated`, `receita_federal_status`.

Remover integralmente os blocos de: `credit_scores`, `financial_profiles`, `credit_accounts`, `credit_inquiries`, `debts`, `payment_histories`, `negative_records`, `employment_records`, `income_declarations`, `legal_records`, `compliance_checks`, `fraud_alerts`, `risk_assessments`, `person_relationships`, `data_sources`, `person_data_sources`.

Substituir o bloco de `persons` por:

```sql
CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    bank_account_profile_id BIGINT,
    last_verified_at TIMESTAMP
);

CREATE INDEX idx_persons_personal_information_id ON persons (personal_information_id);
CREATE INDEX idx_persons_bank_account_profile_id ON persons (bank_account_profile_id);
```

Acrescentar, **antes** do bloco de `persons` (a tabela é referenciada por ele):

```sql
CREATE TABLE bank_account_profiles (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL,
    profile_date TIMESTAMP NOT NULL,
    banking_relationships INTEGER NOT NULL DEFAULT 0,
    account_age_average INTEGER,
    has_checking_account BOOLEAN NOT NULL DEFAULT FALSE,
    has_savings_account BOOLEAN NOT NULL DEFAULT FALSE,
    has_investment_account BOOLEAN NOT NULL DEFAULT FALSE,
    investments_value NUMERIC(15, 2)
);

CREATE UNIQUE INDEX idx_person_bank_profile_date
    ON bank_account_profiles (person_id, profile_date);
```

E, **depois** do bloco de `persons`:

```sql
CREATE TABLE bank_statements (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    institution VARCHAR(255) NOT NULL,
    institution_document VARCHAR(14),
    account_type VARCHAR(50) NOT NULL,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    opening_balance NUMERIC(15, 2) NOT NULL,
    closing_balance NUMERIC(15, 2) NOT NULL,
    total_credits NUMERIC(15, 2) NOT NULL,
    total_debits NUMERIC(15, 2) NOT NULL,
    transaction_count INTEGER DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'BRL'
);

CREATE INDEX idx_bank_statements_person_id ON bank_statements (person_id);
CREATE INDEX idx_bank_statements_period_start ON bank_statements (period_start);
CREATE INDEX idx_bank_statements_period_end ON bank_statements (period_end);

CREATE TABLE cash_flow_analyses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    analysis_date TIMESTAMP NOT NULL,
    period_days INTEGER NOT NULL,
    average_monthly_inflow NUMERIC(15, 2) NOT NULL,
    average_monthly_outflow NUMERIC(15, 2) NOT NULL,
    net_cash_flow NUMERIC(15, 2) NOT NULL,
    inflow_volatility NUMERIC(6, 4),
    negative_balance_days INTEGER DEFAULT 0,
    has_recurring_income BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_cash_flow_analyses_person_id ON cash_flow_analyses (person_id);
CREATE INDEX idx_cash_flow_analyses_analysis_date ON cash_flow_analyses (analysis_date);
CREATE INDEX idx_cash_flow_analyses_net_cash_flow ON cash_flow_analyses (net_cash_flow);
CREATE INDEX idx_cash_flow_analyses_has_recurring_income ON cash_flow_analyses (has_recurring_income);

CREATE TABLE recurring_transactions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    transaction_type VARCHAR(20) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    amount NUMERIC(15, 2) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    counterparty VARCHAR(255),
    first_detected_date TIMESTAMP NOT NULL,
    last_occurrence_date TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_recurring_transactions_person_id ON recurring_transactions (person_id);
CREATE INDEX idx_recurring_transactions_type ON recurring_transactions (transaction_type);
CREATE INDEX idx_recurring_transactions_last_occurrence ON recurring_transactions (last_occurrence_date);
CREATE INDEX idx_recurring_transactions_is_active ON recurring_transactions (is_active);

CREATE TABLE data_sharing_consents (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    consent_id VARCHAR(100) NOT NULL UNIQUE,
    institution VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    scope JSON,
    granted_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP
);

CREATE INDEX idx_data_sharing_consents_person_id ON data_sharing_consents (person_id);
CREATE INDEX idx_data_sharing_consents_status ON data_sharing_consents (status);
```

- [ ] **Step 3: Reescrever `000001_init.down.sql`**

```sql
DROP TABLE IF EXISTS data_sharing_consents;
DROP TABLE IF EXISTS recurring_transactions;
DROP TABLE IF EXISTS cash_flow_analyses;
DROP TABLE IF EXISTS bank_statements;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS bank_account_profiles;
DROP TABLE IF EXISTS analysts;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS person_documents;
DROP TABLE IF EXISTS person_addresses;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS personal_informations;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS files;
```

- [ ] **Step 4: Aplicar e reverter, para provar que as duas direções funcionam**

```bash
cd open-finance
make migration-up
make migration-down
make migration-up
```

Esperado: nenhuma das três chamadas emite erro. Se `make` reclamar de `.env` ausente, copiar o `.env` do birô ajustando `DB_NAME=open-finance`.

- [ ] **Step 5: Conferir o esquema aplicado**

```bash
docker compose exec postgres psql -U "$DB_USERNAME" -d open-finance -c '\dt'
```

Esperado: as 17 tabelas listadas (as 12 mantidas mais as 5 de Open Finance) e **nenhuma** tabela de birô (`credit_scores`, `debts`, `fraud_alerts`…).

- [ ] **Step 6: Commit**

```bash
git add init.sql open-finance/internal/infra/thirdparty/database/migrations/
git commit -m "feat: esquema de dados próprio do servidor de Open Finance"
```

---

### Task 2: Repositórios das quatro dimensões

Hoje só existe `PersonRepository`. As ferramentas por dimensão precisam consultar cada agregado isoladamente.

**Files:**
- Create: `open-finance/internal/infra/repositories/interfaces/bank_statement_repository.go`
- Create: `open-finance/internal/infra/repositories/interfaces/cash_flow_analysis_repository.go`
- Create: `open-finance/internal/infra/repositories/interfaces/recurring_transaction_repository.go`
- Create: `open-finance/internal/infra/repositories/interfaces/data_sharing_consent_repository.go`
- Create: `open-finance/internal/infra/repositories/gorm_bank_statement_repository.go`
- Create: `open-finance/internal/infra/repositories/gorm_cash_flow_analysis_repository.go`
- Create: `open-finance/internal/infra/repositories/gorm_recurring_transaction_repository.go`
- Create: `open-finance/internal/infra/repositories/gorm_data_sharing_consent_repository.go`
- Modify: `open-finance/internal/infra/repositories/factory.go`
- Test: `open-finance/internal/infra/repositories/contract_test.go`

**Interfaces:**
- Consumes: as tabelas da Task 1; `entities.BankStatement`, `entities.CashFlowAnalysis`, `entities.RecurringTransaction`, `entities.DataSharingConsent`.
- Produces:
  - `interfaces.BankStatementRepository.GetByPersonID(ctx context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr)` — `accountType` vazio significa sem filtro
  - `interfaces.CashFlowAnalysisRepository.GetByPersonID(ctx context.Context, personID uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr)` — ordenado por `analysis_date` decrescente; `limit` 0 devolve tudo
  - `interfaces.RecurringTransactionRepository.GetByPersonID(ctx context.Context, personID uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr)`
  - `interfaces.DataSharingConsentRepository.GetByPersonID(ctx context.Context, personID uint) ([]entities.DataSharingConsent, *rest_err.RestErr)`
  - `RepositoryFactory.BankStatementRepository()`, `.CashFlowAnalysisRepository()`, `.RecurringTransactionRepository()`, `.DataSharingConsentRepository()`

Estes repositórios devolvem lista vazia sem erro quando não há registros: ausência de dado é informação válida para a política, não falha.

- [ ] **Step 1: Escrever o teste de contrato que falha**

Os repositórios Gorm não são testados por unidade (exigiriam banco); o que se verifica em compilação é que cada implementação satisfaz sua interface e que a factory as expõe.

Criar `internal/infra/repositories/contract_test.go`:

```go
package repositories

import (
	"testing"

	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
)

// As asserções abaixo falham em compilação se alguma implementação divergir da
// interface, que é o erro mais provável ao adicionar um repositório novo.
var (
	_ interfaces.BankStatementRepository       = (*gormBankStatementRepository)(nil)
	_ interfaces.CashFlowAnalysisRepository    = (*gormCashFlowAnalysisRepository)(nil)
	_ interfaces.RecurringTransactionRepository = (*gormRecurringTransactionRepository)(nil)
	_ interfaces.DataSharingConsentRepository  = (*gormDataSharingConsentRepository)(nil)
)

func TestFactoryExposesOpenFinanceRepositories(t *testing.T) {
	f := &RepositoryFactory{}

	if f.BankStatementRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.CashFlowAnalysisRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.RecurringTransactionRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.DataSharingConsentRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
}
```

- [ ] **Step 2: Rodar e confirmar que falha**

Run: `cd open-finance && go test ./internal/infra/repositories/...`
Expected: FAIL na compilação, com `undefined: gormBankStatementRepository` e os demais.

- [ ] **Step 3: Escrever as quatro interfaces**

`interfaces/bank_statement_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type BankStatementRepository interface {
	// GetByPersonID devolve os extratos da pessoa, do mais recente para o mais
	// antigo. accountType vazio não filtra por tipo de conta.
	GetByPersonID(ctx context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr)
}
```

`interfaces/cash_flow_analysis_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type CashFlowAnalysisRepository interface {
	// GetByPersonID devolve as análises da pessoa ordenadas por analysis_date
	// decrescente. limit menor ou igual a zero devolve todas.
	GetByPersonID(ctx context.Context, personID uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr)
}
```

`interfaces/recurring_transaction_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type RecurringTransactionRepository interface {
	// GetByPersonID devolve as transações recorrentes da pessoa.
	// transactionType vazio não filtra; onlyActive restringe a is_active = true.
	GetByPersonID(ctx context.Context, personID uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr)
}
```

`interfaces/data_sharing_consent_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type DataSharingConsentRepository interface {
	GetByPersonID(ctx context.Context, personID uint) ([]entities.DataSharingConsent, *rest_err.RestErr)
}
```

- [ ] **Step 4: Implementar os quatro repositórios Gorm**

`gorm_bank_statement_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormBankStatementRepository struct {
	db *gorm.DB
}

func NewGormBankStatementRepository(db *gorm.DB) interfaces.BankStatementRepository {
	return &gormBankStatementRepository{db: db}
}

func (g *gormBankStatementRepository) GetByPersonID(ctx context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr) {
	query := gorm.G[entities.BankStatement](g.db).Where("person_id = ?", personID)
	if accountType != "" {
		query = query.Where("account_type = ?", accountType)
	}

	statements, err := query.Order("period_end DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching bank statements").WithCause(err)
	}

	return statements, nil
}
```

`gorm_cash_flow_analysis_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormCashFlowAnalysisRepository struct {
	db *gorm.DB
}

func NewGormCashFlowAnalysisRepository(db *gorm.DB) interfaces.CashFlowAnalysisRepository {
	return &gormCashFlowAnalysisRepository{db: db}
}

func (g *gormCashFlowAnalysisRepository) GetByPersonID(ctx context.Context, personID uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr) {
	query := gorm.G[entities.CashFlowAnalysis](g.db).
		Where("person_id = ?", personID).
		Order("analysis_date DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	analyses, err := query.Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching cash flow analyses").WithCause(err)
	}

	return analyses, nil
}
```

`gorm_recurring_transaction_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormRecurringTransactionRepository struct {
	db *gorm.DB
}

func NewGormRecurringTransactionRepository(db *gorm.DB) interfaces.RecurringTransactionRepository {
	return &gormRecurringTransactionRepository{db: db}
}

func (g *gormRecurringTransactionRepository) GetByPersonID(ctx context.Context, personID uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr) {
	query := gorm.G[entities.RecurringTransaction](g.db).Where("person_id = ?", personID)
	if transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	transactions, err := query.Order("last_occurrence_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching recurring transactions").WithCause(err)
	}

	return transactions, nil
}
```

`gorm_data_sharing_consent_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormDataSharingConsentRepository struct {
	db *gorm.DB
}

func NewGormDataSharingConsentRepository(db *gorm.DB) interfaces.DataSharingConsentRepository {
	return &gormDataSharingConsentRepository{db: db}
}

func (g *gormDataSharingConsentRepository) GetByPersonID(ctx context.Context, personID uint) ([]entities.DataSharingConsent, *rest_err.RestErr) {
	consents, err := gorm.G[entities.DataSharingConsent](g.db).
		Where("person_id = ?", personID).
		Order("granted_at DESC").
		Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching data sharing consents").WithCause(err)
	}

	return consents, nil
}
```

- [ ] **Step 5: Registrar os quatro na factory**

Em `factory.go`, acrescentar os campos ao struct `RepositoryFactory`:

```go
	bankStatementRepository       interfaces.BankStatementRepository
	cashFlowAnalysisRepository    interfaces.CashFlowAnalysisRepository
	recurringTransactionRepository interfaces.RecurringTransactionRepository
	dataSharingConsentRepository  interfaces.DataSharingConsentRepository
```

Preenchê-los em `NewRepositoryFactory`:

```go
		bankStatementRepository:        NewGormBankStatementRepository(tpf.DB()),
		cashFlowAnalysisRepository:     NewGormCashFlowAnalysisRepository(tpf.DB()),
		recurringTransactionRepository: NewGormRecurringTransactionRepository(tpf.DB()),
		dataSharingConsentRepository:   NewGormDataSharingConsentRepository(tpf.DB()),
```

E acrescentar os getters:

```go
func (f *RepositoryFactory) BankStatementRepository() interfaces.BankStatementRepository {
	return f.bankStatementRepository
}

func (f *RepositoryFactory) CashFlowAnalysisRepository() interfaces.CashFlowAnalysisRepository {
	return f.cashFlowAnalysisRepository
}

func (f *RepositoryFactory) RecurringTransactionRepository() interfaces.RecurringTransactionRepository {
	return f.recurringTransactionRepository
}

func (f *RepositoryFactory) DataSharingConsentRepository() interfaces.DataSharingConsentRepository {
	return f.dataSharingConsentRepository
}
```

- [ ] **Step 6: Rodar o teste e confirmar que passa**

Run: `cd open-finance && go test ./internal/infra/repositories/... && go build ./...`
Expected: PASS e build sem erro.

- [ ] **Step 7: Commit**

```bash
git add open-finance/internal/infra/repositories/
git commit -m "feat: repositórios das dimensões de Open Finance"
```

---

### Task 3: DTOs de listagem e de resultado

`PersonSummaryDTO` força o agente a buscar o cadastro completo. Os DTOs de resultado dão às ferramentas por dimensão um objeto de topo, exigido pela validação de esquema de saída do MCP, e carregam a identificação do cliente como metadado de rastreabilidade.

**Files:**
- Modify: `open-finance/internal/infra/controllers/dto/person_dto.go`
- Create: `open-finance/internal/infra/controllers/dto/open_finance_result_dto.go`
- Test: `open-finance/internal/infra/controllers/dto/person_dto_test.go`

**Interfaces:**
- Consumes: `entities.Person`, `entities.PersonalInformation` e os DTOs já existentes em `open_finance_dto.go`.
- Produces:
  - `dto.PersonSummaryDTO{ID uint; Name string; Document string}` e `dto.NewPersonSummaryDTO(*entities.Person) *dto.PersonSummaryDTO`
  - `dto.BankStatementsResultDTO`, `dto.CashFlowAnalysesResultDTO`, `dto.RecurringTransactionsResultDTO`, `dto.DataSharingConsentsResultDTO`, cada um com `CustomerID uint`, `Document string` e `Items []<DTO da dimensão>`

- [ ] **Step 1: Escrever os testes que falham**

Criar `dto/person_dto_test.go`:

```go
package dto

import (
	"testing"

	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonSummaryDTO(t *testing.T) {
	tests := []struct {
		name         string
		entity       *entities.Person
		wantID       uint
		wantName     string
		wantDocument string
	}{
		{
			name: "pessoa com informações pessoais",
			entity: &entities.Person{
				Model: gorm.Model{ID: 7},
				PersonalInformation: &entities.PersonalInformation{
					FullName: "Igor Souza Martins",
					Document: valueobjects.Document("40161087990"),
				},
			},
			wantID:       7,
			wantName:     "Igor Souza Martins",
			wantDocument: "40161087990",
		},
		{
			name:         "pessoa sem informações pessoais não deve causar panic",
			entity:       &entities.Person{Model: gorm.Model{ID: 3}},
			wantID:       3,
			wantName:     "",
			wantDocument: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPersonSummaryDTO(tt.entity)

			if got.ID != tt.wantID {
				t.Errorf("ID = %d, esperado %d", got.ID, tt.wantID)
			}
			if got.Name != tt.wantName {
				t.Errorf("Name = %q, esperado %q", got.Name, tt.wantName)
			}
			if got.Document != tt.wantDocument {
				t.Errorf("Document = %q, esperado %q", got.Document, tt.wantDocument)
			}
		})
	}
}

func TestNewPersonDTOComAssociacoesVazias(t *testing.T) {
	got := NewPersonDTO(&entities.Person{Model: gorm.Model{ID: 1}})

	if got.ID != 1 {
		t.Errorf("ID = %d, esperado 1", got.ID)
	}
	if got.PersonalInformation != nil {
		t.Error("PersonalInformation deveria ser nil quando a entidade não a carrega")
	}
	if got.BankAccountProfile != nil {
		t.Error("BankAccountProfile deveria ser nil quando a entidade não o carrega")
	}
	if len(got.BankStatements) != 0 {
		t.Errorf("BankStatements = %d itens, esperado 0", len(got.BankStatements))
	}
	if len(got.DataSharingConsents) != 0 {
		t.Errorf("DataSharingConsents = %d itens, esperado 0", len(got.DataSharingConsents))
	}
}

func TestNewPersonDTOAgregaDimensoes(t *testing.T) {
	entity := &entities.Person{
		Model: gorm.Model{ID: 2},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Henrique Martins Barbosa",
			Document: valueobjects.Document("15393772025"),
		},
		BankAccountProfile: &entities.BankAccountProfile{BankingRelationships: 3},
		BankStatements:     []entities.BankStatement{{Institution: "Banco Sintético Beta"}},
		CashFlowAnalyses:   []entities.CashFlowAnalysis{{NetCashFlow: 1500}},
		RecurringTransactions: []entities.RecurringTransaction{
			{TransactionType: "income", Amount: 6950},
		},
		DataSharingConsents: []entities.DataSharingConsent{{Status: "granted"}},
	}

	got := NewPersonDTO(entity)

	if got.PersonalInformation == nil || got.PersonalInformation.FullName != "Henrique Martins Barbosa" {
		t.Error("PersonalInformation não foi mapeada")
	}
	if got.BankAccountProfile == nil || got.BankAccountProfile.BankingRelationships != 3 {
		t.Error("BankAccountProfile não foi mapeado")
	}
	if len(got.BankStatements) != 1 || got.BankStatements[0].Institution != "Banco Sintético Beta" {
		t.Error("BankStatements não foram mapeados")
	}
	if len(got.CashFlowAnalyses) != 1 || got.CashFlowAnalyses[0].NetCashFlow != 1500 {
		t.Error("CashFlowAnalyses não foram mapeadas")
	}
	if len(got.RecurringTransactions) != 1 || got.RecurringTransactions[0].Amount != 6950 {
		t.Error("RecurringTransactions não foram mapeadas")
	}
	if len(got.DataSharingConsents) != 1 || got.DataSharingConsents[0].Status != "granted" {
		t.Error("DataSharingConsents não foram mapeados")
	}
}
```

- [ ] **Step 2: Rodar e confirmar que falha**

Run: `cd open-finance && go test ./internal/infra/controllers/dto/...`
Expected: FAIL na compilação, com `undefined: NewPersonSummaryDTO`.

- [ ] **Step 3: Implementar `PersonSummaryDTO`**

Acrescentar ao final de `dto/person_dto.go`:

```go
// PersonSummaryDTO é uma projeção enxuta de um cliente, usada pela listagem.
// Expõe apenas o necessário para identificá-lo, de modo que o chamador seja
// obrigado a buscar o cadastro completo por get_customer_by_id ou
// get_customer_by_document antes de qualquer análise.
type PersonSummaryDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

func NewPersonSummaryDTO(entity *entities.Person) *PersonSummaryDTO {
	dto := &PersonSummaryDTO{
		ID: entity.ID,
	}

	if entity.PersonalInformation != nil {
		dto.Name = entity.PersonalInformation.FullName
		dto.Document = entity.PersonalInformation.Document.String()
	}

	return dto
}
```

- [ ] **Step 4: Implementar os DTOs de resultado**

Criar `dto/open_finance_result_dto.go`:

```go
package dto

// Os DTOs abaixo embrulham as listas devolvidas pelas ferramentas por dimensão.
// O MCP valida o esquema de saída contra um objeto de topo, e os campos de
// identificação servem de metadado de rastreabilidade da consulta.

type BankStatementsResultDTO struct {
	CustomerID uint               `json:"customer_id"`
	Document   string             `json:"document"`
	Items      []BankStatementDTO `json:"items"`
}

type CashFlowAnalysesResultDTO struct {
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []CashFlowAnalysisDTO `json:"items"`
}

type RecurringTransactionsResultDTO struct {
	CustomerID uint                      `json:"customer_id"`
	Document   string                    `json:"document"`
	Items      []RecurringTransactionDTO `json:"items"`
}

type DataSharingConsentsResultDTO struct {
	CustomerID uint                    `json:"customer_id"`
	Document   string                  `json:"document"`
	Items      []DataSharingConsentDTO `json:"items"`
}
```

- [ ] **Step 5: Rodar os testes e confirmar que passam**

Run: `cd open-finance && go test ./internal/infra/controllers/dto/... -v`
Expected: PASS nos três testes.

- [ ] **Step 6: Commit**

```bash
git add open-finance/internal/infra/controllers/dto/
git commit -m "feat: DTOs de listagem resumida e de resultado das dimensões"
```

---

### Task 4: OpenFinanceService

Orquestra os quatro repositórios e resolve a identificação do cliente, que pode chegar por id ou por CPF.

**Files:**
- Create: `open-finance/internal/services/open_finance_service.go`
- Modify: `open-finance/internal/services/person_service.go`
- Test: `open-finance/internal/services/open_finance_service_test.go`

**Interfaces:**
- Consumes: os repositórios da Task 2, os DTOs da Task 3, `interfaces.PersonRepository`.
- Produces:
  - `services.CustomerRef{ID uint; Document string}`
  - `services.NewOpenFinanceService(rf *repositories.RepositoryFactory) *OpenFinanceService`
  - `(*OpenFinanceService).GetBankStatements(ctx, ref CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr)`
  - `(*OpenFinanceService).GetCashFlowAnalyses(ctx, ref CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr)`
  - `(*OpenFinanceService).GetRecurringTransactions(ctx, ref CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr)`
  - `(*OpenFinanceService).GetDataSharingConsents(ctx, ref CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr)`
  - `(*PersonService).GetAllSummary(ctx, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)`

- [ ] **Step 1: Escrever os testes que falham**

Criar `internal/services/open_finance_service_test.go`. Os dublês são declarados no próprio arquivo de teste; como o teste vive no pacote `services`, pode preencher os campos não exportados do serviço diretamente.

```go
package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type fakePersonRepository struct {
	byID       map[uint]*entities.Person
	byDocument map[string]*entities.Person
}

func (f *fakePersonRepository) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byID[id]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byDocument[document]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetAll(_ context.Context, _, _ int, _ map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	return nil, 0, nil
}

func (f *fakePersonRepository) Delete(_ context.Context, _ uint) *rest_err.RestErr { return nil }

type fakeBankStatementRepository struct {
	gotPersonID    uint
	gotAccountType string
	statements     []entities.BankStatement
	err            *rest_err.RestErr
}

func (f *fakeBankStatementRepository) GetByPersonID(_ context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr) {
	f.gotPersonID = personID
	f.gotAccountType = accountType
	return f.statements, f.err
}

type fakeCashFlowAnalysisRepository struct {
	gotLimit int
	analyses []entities.CashFlowAnalysis
}

func (f *fakeCashFlowAnalysisRepository) GetByPersonID(_ context.Context, _ uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr) {
	f.gotLimit = limit
	return f.analyses, nil
}

type fakeRecurringTransactionRepository struct {
	gotType       string
	gotOnlyActive bool
	transactions  []entities.RecurringTransaction
}

func (f *fakeRecurringTransactionRepository) GetByPersonID(_ context.Context, _ uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr) {
	f.gotType = transactionType
	f.gotOnlyActive = onlyActive
	return f.transactions, nil
}

type fakeDataSharingConsentRepository struct {
	consents []entities.DataSharingConsent
}

func (f *fakeDataSharingConsentRepository) GetByPersonID(_ context.Context, _ uint) ([]entities.DataSharingConsent, *rest_err.RestErr) {
	return f.consents, nil
}

func personFixture() *entities.Person {
	return &entities.Person{
		Model: gorm.Model{ID: 8},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Lucas Martins Souza",
			Document: valueobjects.Document("18979021232"),
		},
	}
}

func newTestService(bs *fakeBankStatementRepository) (*OpenFinanceService, *fakeCashFlowAnalysisRepository, *fakeRecurringTransactionRepository, *fakeDataSharingConsentRepository) {
	person := personFixture()
	cf := &fakeCashFlowAnalysisRepository{}
	rt := &fakeRecurringTransactionRepository{}
	dc := &fakeDataSharingConsentRepository{}

	svc := &OpenFinanceService{
		personRepository: &fakePersonRepository{
			byID:       map[uint]*entities.Person{8: person},
			byDocument: map[string]*entities.Person{"18979021232": person},
		},
		bankStatementRepository:        bs,
		cashFlowAnalysisRepository:     cf,
		recurringTransactionRepository: rt,
		dataSharingConsentRepository:   dc,
	}

	return svc, cf, rt, dc
}

func TestResolveCustomerRef(t *testing.T) {
	tests := []struct {
		name         string
		ref          CustomerRef
		wantErr      bool
		wantMessage  string
		wantCustomer uint
	}{
		{name: "por id", ref: CustomerRef{ID: 8}, wantCustomer: 8},
		{name: "por documento", ref: CustomerRef{Document: "18979021232"}, wantCustomer: 8},
		{
			name:        "os dois informados",
			ref:         CustomerRef{ID: 8, Document: "18979021232"},
			wantErr:     true,
			wantMessage: "informe apenas um entre customer_id e document",
		},
		{
			name:        "nenhum informado",
			ref:         CustomerRef{},
			wantErr:     true,
			wantMessage: "informe customer_id ou document",
		},
		{
			name:        "id inexistente",
			ref:         CustomerRef{ID: 99},
			wantErr:     true,
			wantMessage: "person not found",
		},
		{
			name:        "documento inexistente",
			ref:         CustomerRef{Document: "00000000000"},
			wantErr:     true,
			wantMessage: "person not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, _, _ := newTestService(&fakeBankStatementRepository{})

			person, err := svc.resolveCustomer(context.Background(), tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				if err.Message != tt.wantMessage {
					t.Errorf("mensagem = %q, esperado %q", err.Message, tt.wantMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if person.ID != tt.wantCustomer {
				t.Errorf("ID = %d, esperado %d", person.ID, tt.wantCustomer)
			}
		})
	}
}

func TestGetBankStatementsRepassaFiltroEIdentificacao(t *testing.T) {
	bs := &fakeBankStatementRepository{
		statements: []entities.BankStatement{{Institution: "Banco Sintético Gama", AccountType: "checking"}},
	}
	svc, _, _, _ := newTestService(bs)

	got, err := svc.GetBankStatements(context.Background(), CustomerRef{Document: "18979021232"}, "checking")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if bs.gotPersonID != 8 {
		t.Errorf("person_id repassado = %d, esperado 8", bs.gotPersonID)
	}
	if bs.gotAccountType != "checking" {
		t.Errorf("account_type repassado = %q, esperado \"checking\"", bs.gotAccountType)
	}
	if got.CustomerID != 8 || got.Document != "18979021232" {
		t.Errorf("identificação = (%d, %q), esperado (8, \"18979021232\")", got.CustomerID, got.Document)
	}
	if len(got.Items) != 1 || got.Items[0].Institution != "Banco Sintético Gama" {
		t.Error("extratos não foram mapeados para o DTO")
	}
}

func TestGetBankStatementsPropagaErroDoRepositorio(t *testing.T) {
	bs := &fakeBankStatementRepository{err: rest_err.NewInternalServerError("boom")}
	svc, _, _, _ := newTestService(bs)

	_, err := svc.GetBankStatements(context.Background(), CustomerRef{ID: 8}, "")

	if err == nil {
		t.Fatal("esperado erro, veio nil")
	}
	if err.Message != "boom" {
		t.Errorf("mensagem = %q, esperado \"boom\"", err.Message)
	}
}

func TestGetBankStatementsSemRegistrosDevolveListaVazia(t *testing.T) {
	svc, _, _, _ := newTestService(&fakeBankStatementRepository{})

	got, err := svc.GetBankStatements(context.Background(), CustomerRef{ID: 8}, "")

	if err != nil {
		t.Fatalf("ausência de extrato não é erro, mas veio: %v", err)
	}
	if got.Items == nil {
		t.Error("Items deve ser lista vazia, nunca nil, para serializar como [] no JSON")
	}
	if len(got.Items) != 0 {
		t.Errorf("Items = %d, esperado 0", len(got.Items))
	}
}

func TestGetCashFlowAnalysesRepassaLimite(t *testing.T) {
	svc, cf, _, _ := newTestService(&fakeBankStatementRepository{})
	cf.analyses = []entities.CashFlowAnalysis{{NetCashFlow: 150}}

	got, err := svc.GetCashFlowAnalyses(context.Background(), CustomerRef{ID: 8}, 1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if cf.gotLimit != 1 {
		t.Errorf("limite repassado = %d, esperado 1", cf.gotLimit)
	}
	if len(got.Items) != 1 || got.Items[0].NetCashFlow != 150 {
		t.Error("análises não foram mapeadas para o DTO")
	}
}

func TestGetRecurringTransactionsRepassaFiltros(t *testing.T) {
	svc, _, rt, _ := newTestService(&fakeBankStatementRepository{})
	rt.transactions = []entities.RecurringTransaction{{TransactionType: "income", Amount: 4200}}

	got, err := svc.GetRecurringTransactions(context.Background(), CustomerRef{ID: 8}, "income", true)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if rt.gotType != "income" {
		t.Errorf("tipo repassado = %q, esperado \"income\"", rt.gotType)
	}
	if !rt.gotOnlyActive {
		t.Error("only_active repassado = false, esperado true")
	}
	if len(got.Items) != 1 || got.Items[0].Amount != 4200 {
		t.Error("transações não foram mapeadas para o DTO")
	}
}

func TestGetDataSharingConsents(t *testing.T) {
	svc, _, _, dc := newTestService(&fakeBankStatementRepository{})
	dc.consents = []entities.DataSharingConsent{{ConsentID: "urn:openfinance:consent:008", Status: "granted"}}

	got, err := svc.GetDataSharingConsents(context.Background(), CustomerRef{ID: 8})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].Status != "granted" {
		t.Error("consentimentos não foram mapeados para o DTO")
	}
}
```

- [ ] **Step 2: Rodar e confirmar que falha**

Run: `cd open-finance && go test ./internal/services/...`
Expected: FAIL na compilação, com `undefined: OpenFinanceService`.

- [ ] **Step 3: Implementar o serviço**

Criar `internal/services/open_finance_service.go`:

```go
package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
)

// CustomerRef identifica o cliente de uma consulta por dimensão. Exatamente um
// dos dois campos deve estar preenchido.
type CustomerRef struct {
	ID       uint
	Document string
}

type OpenFinanceService struct {
	personRepository               interfaces.PersonRepository
	bankStatementRepository        interfaces.BankStatementRepository
	cashFlowAnalysisRepository     interfaces.CashFlowAnalysisRepository
	recurringTransactionRepository interfaces.RecurringTransactionRepository
	dataSharingConsentRepository   interfaces.DataSharingConsentRepository
}

func NewOpenFinanceService(rf *repositories.RepositoryFactory) *OpenFinanceService {
	return &OpenFinanceService{
		personRepository:               rf.PersonRepository(),
		bankStatementRepository:        rf.BankStatementRepository(),
		cashFlowAnalysisRepository:     rf.CashFlowAnalysisRepository(),
		recurringTransactionRepository: rf.RecurringTransactionRepository(),
		dataSharingConsentRepository:   rf.DataSharingConsentRepository(),
	}
}

// resolveCustomer traduz a referência recebida pela ferramenta em uma pessoa
// existente, recusando referências ambíguas ou vazias.
func (s *OpenFinanceService) resolveCustomer(ctx context.Context, ref CustomerRef) (*entities.Person, *rest_err.RestErr) {
	hasID := ref.ID > 0
	hasDocument := ref.Document != ""

	switch {
	case hasID && hasDocument:
		return nil, rest_err.NewBadRequestError("informe apenas um entre customer_id e document")
	case !hasID && !hasDocument:
		return nil, rest_err.NewBadRequestError("informe customer_id ou document")
	case hasID:
		return s.personRepository.GetById(ctx, ref.ID)
	default:
		return s.personRepository.GetByDocument(ctx, ref.Document)
	}
}

func documentOf(person *entities.Person) string {
	if person.PersonalInformation == nil {
		return ""
	}
	return person.PersonalInformation.Document.String()
}

func (s *OpenFinanceService) GetBankStatements(ctx context.Context, ref CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	statements, err := s.bankStatementRepository.GetByPersonID(ctx, person.ID, accountType)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BankStatementDTO, 0, len(statements))
	for i := range statements {
		items = append(items, *dto.NewBankStatementDTO(&statements[i]))
	}

	return &dto.BankStatementsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetCashFlowAnalyses(ctx context.Context, ref CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	analyses, err := s.cashFlowAnalysisRepository.GetByPersonID(ctx, person.ID, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.CashFlowAnalysisDTO, 0, len(analyses))
	for i := range analyses {
		items = append(items, *dto.NewCashFlowAnalysisDTO(&analyses[i]))
	}

	return &dto.CashFlowAnalysesResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetRecurringTransactions(ctx context.Context, ref CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	transactions, err := s.recurringTransactionRepository.GetByPersonID(ctx, person.ID, transactionType, onlyActive)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RecurringTransactionDTO, 0, len(transactions))
	for i := range transactions {
		items = append(items, *dto.NewRecurringTransactionDTO(&transactions[i]))
	}

	return &dto.RecurringTransactionsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetDataSharingConsents(ctx context.Context, ref CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	consents, err := s.dataSharingConsentRepository.GetByPersonID(ctx, person.ID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DataSharingConsentDTO, 0, len(consents))
	for i := range consents {
		items = append(items, *dto.NewDataSharingConsentDTO(&consents[i]))
	}

	return &dto.DataSharingConsentsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}
```

- [ ] **Step 4: Acrescentar `GetAllSummary` ao `PersonService`**

Ao final de `internal/services/person_service.go`:

```go
func (as *PersonService) GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr) {
	persons, count, err := as.personRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginated := dto.NewPaginatedResponse(count, make([]*dto.PersonSummaryDTO, len(persons)))

	for i := range persons {
		paginated.Items[i] = dto.NewPersonSummaryDTO(&persons[i])
	}

	return paginated, nil
}
```

- [ ] **Step 5: Rodar os testes e confirmar que passam**

Run: `cd open-finance && go test ./internal/services/... -v && go build ./...`
Expected: PASS em todos os casos de `TestResolveCustomerRef` e nos demais; build limpo.

- [ ] **Step 6: Commit**

```bash
git add open-finance/internal/services/
git commit -m "feat: serviço de Open Finance com resolução de cliente por id ou documento"
```

---

### Task 5: Ferramentas consolidadas

Renomeia as três ferramentas herdadas, reescreve as descrições para o domínio de Open Finance e passa a listagem a usar o resumo. As dependências do `Server` viram interfaces, para que os handlers possam ser testados com dublês.

**Files:**
- Modify: `open-finance/internal/interfaces/mcp/tools/server.go`
- Modify: `open-finance/internal/interfaces/mcp/tools/person.go`
- Modify: `open-finance/internal/interfaces/mcp/mcp.go`
- Test: `open-finance/internal/interfaces/mcp/tools/person_test.go`

**Interfaces:**
- Consumes: `services.PersonService` (Task 4), `dto.PersonSummaryDTO` (Task 3).
- Produces:
  - `tools.PersonService` — interface com `GetById`, `GetByDocument`, `GetAllSummary`
  - `tools.OpenFinanceService` — interface com os quatro métodos da Task 4 (usada pela Task 6)
  - `tools.NewMCPServer(userService *services.UserService, addressService *services.AddressService, analystService *services.AnalystService, personService PersonService, openFinanceService OpenFinanceService) *server.MCPServer`
  - Handlers `HandleGetPersonByID`, `HandleGetPersonByDocument`, `HandleGetAllPersons`

- [ ] **Step 1: Escrever os testes que falham**

Criar `internal/interfaces/mcp/tools/person_test.go`:

```go
package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
)

type fakePersonService struct {
	person    *entities.Person
	summary   *dto.PaginatedResponse[dto.PersonSummaryDTO]
	err       *rest_err.RestErr
	gotID     uint
	gotDoc    string
	gotLimit  int
	gotOffset int
	gotParams map[string]any
}

func (f *fakePersonService) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	f.gotID = id
	return f.person, f.err
}

func (f *fakePersonService) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	f.gotDoc = document
	return f.person, f.err
}

func (f *fakePersonService) GetAllSummary(_ context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr) {
	f.gotLimit, f.gotOffset, f.gotParams = limit, offset, params
	return f.summary, f.err
}

func requestWith(args map[string]any) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Arguments = args
	return r
}

func personEntity() *entities.Person {
	return &entities.Person{
		Model: gorm.Model{ID: 5},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Gabriela Ribeiro Barbosa",
			Document: valueobjects.Document("35509139404"),
		},
	}
}

func TestHandleGetPersonByID(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		wantID  uint
	}{
		{name: "id válido", args: map[string]any{"id": float64(5)}, wantID: 5},
		{name: "id ausente", args: map[string]any{}, wantErr: true},
		{name: "id zero", args: map[string]any{"id": float64(0)}, wantErr: true},
		{name: "id negativo", args: map[string]any{"id": float64(-3)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakePersonService{person: personEntity()}
			s := &Server{personService: svc}

			got, err := s.HandleGetPersonByID(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if svc.gotID != tt.wantID {
				t.Errorf("id repassado = %d, esperado %d", svc.gotID, tt.wantID)
			}
			if got.ID != 5 {
				t.Errorf("DTO.ID = %d, esperado 5", got.ID)
			}
		})
	}
}

func TestHandleGetPersonByIDPropagaNotFound(t *testing.T) {
	svc := &fakePersonService{err: rest_err.NewNotFoundError("person not found")}
	s := &Server{personService: svc}

	_, err := s.HandleGetPersonByID(context.Background(), requestWith(map[string]any{"id": float64(99)}), mcp.CallToolParams{})

	if err == nil {
		t.Fatal("esperado erro, veio nil")
	}
}

func TestHandleGetPersonByDocument(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		wantDoc string
	}{
		{name: "documento válido", args: map[string]any{"document": "35509139404"}, wantDoc: "35509139404"},
		{name: "documento ausente", args: map[string]any{}, wantErr: true},
		{name: "documento vazio", args: map[string]any{"document": ""}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakePersonService{person: personEntity()}
			s := &Server{personService: svc}

			_, err := s.HandleGetPersonByDocument(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if svc.gotDoc != tt.wantDoc {
				t.Errorf("documento repassado = %q, esperado %q", svc.gotDoc, tt.wantDoc)
			}
		})
	}
}

func TestHandleGetAllPersonsUsaPadroesEFiltros(t *testing.T) {
	svc := &fakePersonService{
		summary: dto.NewPaginatedResponse(1, []*dto.PersonSummaryDTO{
			{ID: 5, Name: "Gabriela Ribeiro Barbosa", Document: "35509139404"},
		}),
	}
	s := &Server{personService: svc}

	got, err := s.HandleGetAllPersons(context.Background(), requestWith(map[string]any{
		"params": map[string]any{"id": float64(5)},
	}), mcp.CallToolParams{})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if svc.gotLimit != 10 || svc.gotOffset != 0 {
		t.Errorf("padrões = (limit %d, offset %d), esperado (10, 0)", svc.gotLimit, svc.gotOffset)
	}
	if svc.gotParams == nil {
		t.Error("params não foi repassado ao serviço")
	}
	if got.Total != 1 || len(got.Items) != 1 {
		t.Errorf("resposta = (total %d, %d itens), esperado (1, 1)", got.Total, len(got.Items))
	}
}

func TestHandleGetAllPersonsRecusaPaginacaoNegativa(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
	}{
		{name: "limit negativo", args: map[string]any{"limit": float64(-1)}},
		{name: "offset negativo", args: map[string]any{"offset": float64(-1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{personService: &fakePersonService{}}

			if _, err := s.HandleGetAllPersons(context.Background(), requestWith(tt.args), mcp.CallToolParams{}); err == nil {
				t.Fatal("esperado erro, veio nil")
			}
		})
	}
}

func TestNomesDasFerramentasConsolidadas(t *testing.T) {
	s := &Server{personService: &fakePersonService{}}

	tests := []struct {
		got  string
		want string
	}{
		{s.GetPersonByIDTool().Name, "get_customer_by_id"},
		{s.GetPersonByDocumentTool().Name, "get_customer_by_document"},
		{s.GetAllPersonsTool().Name, "get_all_customers"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("nome = %q, esperado %q", tt.got, tt.want)
		}
	}
}
```

- [ ] **Step 2: Rodar e confirmar que falha**

Run: `cd open-finance && go test ./internal/interfaces/mcp/tools/...`
Expected: FAIL — `cannot use svc (variable of type *fakePersonService) as *services.PersonService`, e nomes de ferramenta divergentes.

- [ ] **Step 3: Trocar as dependências do `Server` por interfaces**

Substituir o conteúdo de `tools/server.go` por:

```go
package tools

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PersonService e OpenFinanceService são declarados aqui, no consumidor, para
// que os handlers possam ser exercitados com dublês nos testes.
type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)
}

type OpenFinanceService interface {
	GetBankStatements(ctx context.Context, ref services.CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr)
	GetCashFlowAnalyses(ctx context.Context, ref services.CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr)
	GetRecurringTransactions(ctx context.Context, ref services.CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr)
	GetDataSharingConsents(ctx context.Context, ref services.CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr)
}

type Server struct {
	userService        *services.UserService
	addressService     *services.AddressService
	analystService     *services.AnalystService
	personService      PersonService
	openFinanceService OpenFinanceService
}

func NewMCPServer(
	userService *services.UserService,
	addressService *services.AddressService,
	analystService *services.AnalystService,
	personService PersonService,
	openFinanceService OpenFinanceService,
) *server.MCPServer {
	s := &Server{
		userService:        userService,
		addressService:     addressService,
		analystService:     analystService,
		personService:      personService,
		openFinanceService: openFinanceService,
	}
	mcpSrv := server.NewMCPServer(
		"open-finance-mcp",
		"1.0.0",
		server.WithRecovery(),
		server.WithOutputSchemaValidation(),
		server.WithToolHandlerMiddleware(middlewares.MCPLogMiddleware),
	)
	s.registerTools(mcpSrv)
	return mcpSrv
}

func (s *Server) registerTools(mcpSrv *server.MCPServer) {
	mcpSrv.AddTool(s.GetPersonByIDTool(), mcp.NewStructuredToolHandler(s.HandleGetPersonByID))
	mcpSrv.AddTool(s.GetPersonByDocumentTool(), mcp.NewStructuredToolHandler(s.HandleGetPersonByDocument))
	mcpSrv.AddTool(s.GetAllPersonsTool(), mcp.NewStructuredToolHandler(s.HandleGetAllPersons))
}
```

- [ ] **Step 4: Reescrever as três ferramentas consolidadas**

Substituir o conteúdo de `tools/person.go` por:

```go
package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) GetPersonByIDTool() mcp.Tool {
	return mcp.NewTool(
		"get_customer_by_id",
		mcp.WithDescription(
			`
			Open Finance: get a customer's shared banking data by their ID.
			Returns the consolidated Open Finance record: banking profile, the last
			90 days of bank statements, cash flow analysis, recurring incomes and
			fixed expenses, and the data sharing consents that authorize access.
			This is NOT credit bureau data; for score, debts and negative records,
			use the credit bureau server.
			Example usage:
			{
				"id": 123
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithInteger(
			"id",
			mcp.Description("The ID of the customer to retrieve"),
		),
	)
}

func (s *Server) HandleGetPersonByID(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PersonDTO, error) {
	id, err := request.RequireInt("id")
	if err != nil {
		return nil, err
	} else if id <= 0 {
		return nil, fmt.Errorf("invalid ID: must be a positive integer")
	}

	person, restErr := s.personService.GetById(ctx, uint(id))
	if restErr != nil {
		return nil, restErr
	}

	return dto.NewPersonDTO(person), nil
}

func (s *Server) GetPersonByDocumentTool() mcp.Tool {
	return mcp.NewTool(
		"get_customer_by_document",
		mcp.WithDescription(
			`
			Open Finance: get a customer's shared banking data by their document (CPF).
			Returns the consolidated Open Finance record: banking profile, the last
			90 days of bank statements, cash flow analysis, recurring incomes and
			fixed expenses, and the data sharing consents that authorize access.
			This is NOT credit bureau data; for score, debts and negative records,
			use the credit bureau server.
			Example usage:
			{
				"document": "12345678900"
			}
			`,
		),
		mcp.WithOutputSchema[dto.PersonDTO](),
		mcp.WithString(
			"document",
			mcp.Description("The document number (CPF) of the customer to retrieve"),
		),
	)
}

func (s *Server) HandleGetPersonByDocument(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PersonDTO, error) {
	document, err := request.RequireString("document")
	if err != nil {
		return nil, err
	} else if document == "" {
		return nil, fmt.Errorf("invalid document: cannot be empty")
	}

	person, restErr := s.personService.GetByDocument(ctx, document)
	if restErr != nil {
		return nil, restErr
	}

	return dto.NewPersonDTO(person), nil
}

func (s *Server) GetAllPersonsTool() mcp.Tool {
	return mcp.NewTool(
		"get_all_customers",
		mcp.WithDescription(
			`
			Open Finance: list customers with pagination.
			This returns only a summary of each customer: their id, name and document.
			It does NOT return their Open Finance data.
			To retrieve the complete record for a customer, call get_customer_by_id
			(using the returned "id") or get_customer_by_document (using the
			returned "document").

			Example usage:
			{
				"limit": 10,
				"offset": 0
			}

			To fetch all customers, set "limit" to 0.
			`,
		),
		mcp.WithOutputSchema[dto.PaginatedResponse[dto.PersonSummaryDTO]](),
		mcp.WithInteger(
			"limit",
			mcp.Description("Limit the number of results"),
			mcp.DefaultNumber(10),
			mcp.Min(0),
		),
		mcp.WithInteger(
			"offset",
			mcp.Description("Offset for pagination"),
			mcp.DefaultNumber(0),
			mcp.Min(0),
		),
		mcp.WithObject(
			"params",
			mcp.Description("Additional filters for querying customers. e.g. {\"id\": 5}"),
		),
	)
}

func (s *Server) HandleGetAllPersons(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PaginatedResponse[dto.PersonSummaryDTO], error) {
	arguments := map[string]any{}
	if request.Params.Arguments != nil {
		var ok bool
		arguments, ok = request.Params.Arguments.(map[string]any)
		if !ok {
			return nil, errors.New("invalid arguments")
		}
	}

	params, _ := arguments["params"].(map[string]any)
	limit := request.GetInt("limit", 10)
	offset := request.GetInt("offset", 0)
	if limit < 0 {
		return nil, fmt.Errorf("invalid limit: must be a non-negative integer")
	}
	if offset < 0 {
		return nil, fmt.Errorf("invalid offset: must be a non-negative integer")
	}

	persons, restErr := s.personService.GetAllSummary(ctx, limit, offset, params)
	if restErr != nil {
		return nil, restErr
	}

	return persons, nil
}
```

- [ ] **Step 5: Injetar o novo serviço na composição**

Em `internal/interfaces/mcp/mcp.go`, acrescentar `services.NewOpenFinanceService(rf)` como quinto argumento:

```go
	s := tools.NewMCPServer(
		services.NewUserService(rf),
		services.NewAddressService(rf),
		services.NewAnalystService(rf),
		services.NewPersonService(rf),
		services.NewOpenFinanceService(rf),
	)
```

- [ ] **Step 6: Rodar os testes e confirmar que passam**

Run: `cd open-finance && go test ./... && go build ./...`
Expected: PASS em todos os pacotes; build limpo.

- [ ] **Step 7: Commit**

```bash
git add open-finance/internal/interfaces/mcp/
git commit -m "feat: ferramentas consolidadas de Open Finance com listagem resumida"
```

---

### Task 6: Ferramentas por dimensão

Quatro ferramentas que expõem cada dimensão isoladamente, para que o agente possa aprofundar sem recarregar o agregado inteiro.

**Files:**
- Create: `open-finance/internal/interfaces/mcp/tools/open_finance.go`
- Modify: `open-finance/internal/interfaces/mcp/tools/server.go:registerTools`
- Test: `open-finance/internal/interfaces/mcp/tools/open_finance_test.go`

**Interfaces:**
- Consumes: `tools.OpenFinanceService` e `services.CustomerRef` (Task 5 e 4).
- Produces: `GetBankStatementsTool`/`HandleGetBankStatements`, `GetCashFlowAnalysisTool`/`HandleGetCashFlowAnalysis`, `GetRecurringTransactionsTool`/`HandleGetRecurringTransactions`, `GetDataSharingConsentsTool`/`HandleGetDataSharingConsents`, e o helper `customerRefFrom(request) (services.CustomerRef, error)`.

- [ ] **Step 1: Escrever os testes que falham**

Criar `internal/interfaces/mcp/tools/open_finance_test.go`:

```go
package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeOpenFinanceService struct {
	gotRef         services.CustomerRef
	gotAccountType string
	gotLimit       int
	gotType        string
	gotOnlyActive  bool
	err            *rest_err.RestErr
}

func (f *fakeOpenFinanceService) GetBankStatements(_ context.Context, ref services.CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotAccountType = ref, accountType
	if f.err != nil {
		return nil, f.err
	}
	return &dto.BankStatementsResultDTO{CustomerID: 8, Items: []dto.BankStatementDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetCashFlowAnalyses(_ context.Context, ref services.CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotLimit = ref, limit
	if f.err != nil {
		return nil, f.err
	}
	return &dto.CashFlowAnalysesResultDTO{CustomerID: 8, Items: []dto.CashFlowAnalysisDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetRecurringTransactions(_ context.Context, ref services.CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotType, f.gotOnlyActive = ref, transactionType, onlyActive
	if f.err != nil {
		return nil, f.err
	}
	return &dto.RecurringTransactionsResultDTO{CustomerID: 8, Items: []dto.RecurringTransactionDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetDataSharingConsents(_ context.Context, ref services.CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.DataSharingConsentsResultDTO{CustomerID: 8, Items: []dto.DataSharingConsentDTO{}}, nil
}

func TestCustomerRefFrom(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		want    services.CustomerRef
	}{
		{
			name: "por customer_id",
			args: map[string]any{"customer_id": float64(8)},
			want: services.CustomerRef{ID: 8},
		},
		{
			name: "por document",
			args: map[string]any{"document": "18979021232"},
			want: services.CustomerRef{Document: "18979021232"},
		},
		{
			name:    "nenhum dos dois",
			args:    map[string]any{},
			wantErr: true,
		},
		{
			name:    "customer_id negativo",
			args:    map[string]any{"customer_id": float64(-1)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRefFrom(requestWith(tt.args))

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("ref = %+v, esperado %+v", got, tt.want)
			}
		})
	}
}

func TestCustomerRefFromComOsDoisDeixaOServicoDecidir(t *testing.T) {
	// A recusa de referência ambígua é responsabilidade do serviço, que é quem
	// conhece a regra. O helper apenas transporta o que veio na chamada.
	got, err := customerRefFrom(requestWith(map[string]any{
		"customer_id": float64(8),
		"document":    "18979021232",
	}))

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got.ID != 8 || got.Document != "18979021232" {
		t.Errorf("ref = %+v, esperado ambos preenchidos", got)
	}
}

func TestHandleGetBankStatements(t *testing.T) {
	tests := []struct {
		name            string
		args            map[string]any
		wantAccountType string
	}{
		{
			name:            "sem filtro de conta",
			args:            map[string]any{"customer_id": float64(8)},
			wantAccountType: "",
		},
		{
			name:            "com filtro de conta",
			args:            map[string]any{"customer_id": float64(8), "account_type": "savings"},
			wantAccountType: "savings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeOpenFinanceService{}
			s := &Server{openFinanceService: svc}

			got, err := s.HandleGetBankStatements(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if svc.gotAccountType != tt.wantAccountType {
				t.Errorf("account_type = %q, esperado %q", svc.gotAccountType, tt.wantAccountType)
			}
			if svc.gotRef.ID != 8 {
				t.Errorf("ref.ID = %d, esperado 8", svc.gotRef.ID)
			}
			if got.CustomerID != 8 {
				t.Errorf("CustomerID = %d, esperado 8", got.CustomerID)
			}
		})
	}
}

func TestHandleGetBankStatementsPropagaErro(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{err: rest_err.NewNotFoundError("person not found")}}

	if _, err := s.HandleGetBankStatements(context.Background(), requestWith(map[string]any{"customer_id": float64(99)}), mcp.CallToolParams{}); err == nil {
		t.Fatal("esperado erro, veio nil")
	}
}

func TestHandleGetCashFlowAnalysisUsaLimitePadraoUm(t *testing.T) {
	svc := &fakeOpenFinanceService{}
	s := &Server{openFinanceService: svc}

	if _, err := s.HandleGetCashFlowAnalysis(context.Background(), requestWith(map[string]any{"customer_id": float64(8)}), mcp.CallToolParams{}); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if svc.gotLimit != 1 {
		t.Errorf("limite = %d, esperado 1 (a análise mais recente)", svc.gotLimit)
	}
}

func TestHandleGetCashFlowAnalysisRecusaLimiteNegativo(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{}}

	if _, err := s.HandleGetCashFlowAnalysis(context.Background(), requestWith(map[string]any{
		"customer_id": float64(8),
		"limit":       float64(-1),
	}), mcp.CallToolParams{}); err == nil {
		t.Fatal("esperado erro, veio nil")
	}
}

func TestHandleGetRecurringTransactions(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]any
		wantType      string
		wantOnlyActiv bool
	}{
		{
			name:          "padrões",
			args:          map[string]any{"customer_id": float64(8)},
			wantType:      "",
			wantOnlyActiv: true,
		},
		{
			name:          "somente receitas, incluindo inativas",
			args:          map[string]any{"customer_id": float64(8), "transaction_type": "income", "only_active": false},
			wantType:      "income",
			wantOnlyActiv: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeOpenFinanceService{}
			s := &Server{openFinanceService: svc}

			if _, err := s.HandleGetRecurringTransactions(context.Background(), requestWith(tt.args), mcp.CallToolParams{}); err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if svc.gotType != tt.wantType {
				t.Errorf("transaction_type = %q, esperado %q", svc.gotType, tt.wantType)
			}
			if svc.gotOnlyActive != tt.wantOnlyActiv {
				t.Errorf("only_active = %v, esperado %v", svc.gotOnlyActive, tt.wantOnlyActiv)
			}
		})
	}
}

func TestHandleGetDataSharingConsents(t *testing.T) {
	svc := &fakeOpenFinanceService{}
	s := &Server{openFinanceService: svc}

	got, err := s.HandleGetDataSharingConsents(context.Background(), requestWith(map[string]any{"document": "18979021232"}), mcp.CallToolParams{})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if svc.gotRef.Document != "18979021232" {
		t.Errorf("ref.Document = %q, esperado \"18979021232\"", svc.gotRef.Document)
	}
	if got.CustomerID != 8 {
		t.Errorf("CustomerID = %d, esperado 8", got.CustomerID)
	}
}

func TestNomesDasFerramentasPorDimensao(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{}}

	tests := []struct {
		got  string
		want string
	}{
		{s.GetBankStatementsTool().Name, "get_bank_statements"},
		{s.GetCashFlowAnalysisTool().Name, "get_cash_flow_analysis"},
		{s.GetRecurringTransactionsTool().Name, "get_recurring_transactions"},
		{s.GetDataSharingConsentsTool().Name, "get_data_sharing_consents"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("nome = %q, esperado %q", tt.got, tt.want)
		}
	}
}
```

- [ ] **Step 2: Rodar e confirmar que falha**

Run: `cd open-finance && go test ./internal/interfaces/mcp/tools/...`
Expected: FAIL na compilação, com `undefined: customerRefFrom` e `s.HandleGetBankStatements undefined`.

- [ ] **Step 3: Implementar as quatro ferramentas**

Criar `internal/interfaces/mcp/tools/open_finance.go`:

```go
package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

const customerRefDescription = `Identify the customer by "customer_id" OR by "document" (CPF). Provide exactly one of them.`

// customerRefFrom apenas transporta a identificação recebida. A recusa de
// referência ambígua ou vazia cabe ao serviço, que é onde a regra vive.
func customerRefFrom(request mcp.CallToolRequest) (services.CustomerRef, error) {
	id := request.GetInt("customer_id", 0)
	if id < 0 {
		return services.CustomerRef{}, fmt.Errorf("invalid customer_id: must be a positive integer")
	}

	document := request.GetString("document", "")
	if id == 0 && document == "" {
		return services.CustomerRef{}, fmt.Errorf("informe customer_id ou document")
	}

	return services.CustomerRef{ID: uint(id), Document: document}, nil
}

func withCustomerRefParams() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithInteger("customer_id", mcp.Description("The ID of the customer")),
		mcp.WithString("document", mcp.Description("The document number (CPF) of the customer")),
	}
}

func (s *Server) GetBankStatementsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's bank statements for the last 90 days.
				Each statement covers one month of one account at one institution and
				reports opening and closing balances, total credits and debits, and the
				number of transactions.
				` + customerRefDescription + `
				Example usage:
				{
					"document": "12345678900",
					"account_type": "checking"
				}
				`,
			),
			mcp.WithOutputSchema[dto.BankStatementsResultDTO](),
			mcp.WithString(
				"account_type",
				mcp.Description("Optional filter by account type"),
				mcp.Enum("checking", "savings", "payment"),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_bank_statements", opts...)
}

func (s *Server) HandleGetBankStatements(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.BankStatementsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetBankStatements(ctx, ref, request.GetString("account_type", ""))
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetCashFlowAnalysisTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's cash flow analysis derived from the
				shared bank statements. Reports average monthly inflow and outflow, net
				cash flow, inflow volatility, days with a negative balance, and whether
				a recurring income was detected.
				` + customerRefDescription + `
				By default returns only the most recent analysis. Set "limit" to 0 to
				retrieve the full history.
				Example usage:
				{
					"customer_id": 123
				}
				`,
			),
			mcp.WithOutputSchema[dto.CashFlowAnalysesResultDTO](),
			mcp.WithInteger(
				"limit",
				mcp.Description("How many analyses to return, most recent first. 0 returns all"),
				mcp.DefaultNumber(1),
				mcp.Min(0),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_cash_flow_analysis", opts...)
}

func (s *Server) HandleGetCashFlowAnalysis(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.CashFlowAnalysesResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	limit := request.GetInt("limit", 1)
	if limit < 0 {
		return nil, fmt.Errorf("invalid limit: must be a non-negative integer")
	}

	result, restErr := s.openFinanceService.GetCashFlowAnalyses(ctx, ref, limit)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetRecurringTransactionsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's recurring incomes and fixed expenses
				identified from the shared transaction history. Each entry reports its
				type, category, amount, frequency, counterparty and whether it is still
				active.
				` + customerRefDescription + `
				Example usage:
				{
					"customer_id": 123,
					"transaction_type": "income"
				}
				`,
			),
			mcp.WithOutputSchema[dto.RecurringTransactionsResultDTO](),
			mcp.WithString(
				"transaction_type",
				mcp.Description("Optional filter by transaction type"),
				mcp.Enum("income", "expense"),
			),
			mcp.WithBoolean(
				"only_active",
				mcp.Description("Return only transactions still active. Defaults to true"),
				mcp.DefaultBool(true),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_recurring_transactions", opts...)
}

func (s *Server) HandleGetRecurringTransactions(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.RecurringTransactionsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetRecurringTransactions(
		ctx,
		ref,
		request.GetString("transaction_type", ""),
		request.GetBool("only_active", true),
	)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetDataSharingConsentsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's data sharing consents. Each consent
				reports the institution, its status (granted, revoked, expired,
				awaiting), the authorized scope, and the grant, expiry and revocation
				timestamps.
				A customer without an active consent has no Open Finance data available;
				treat the corresponding criteria as missing data rather than as zero.
				` + customerRefDescription + `
				Example usage:
				{
					"document": "12345678900"
				}
				`,
			),
			mcp.WithOutputSchema[dto.DataSharingConsentsResultDTO](),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_data_sharing_consents", opts...)
}

func (s *Server) HandleGetDataSharingConsents(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.DataSharingConsentsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetDataSharingConsents(ctx, ref)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}
```

- [ ] **Step 4: Registrar as quatro**

Em `tools/server.go`, acrescentar ao final de `registerTools`:

```go
	mcpSrv.AddTool(s.GetBankStatementsTool(), mcp.NewStructuredToolHandler(s.HandleGetBankStatements))
	mcpSrv.AddTool(s.GetCashFlowAnalysisTool(), mcp.NewStructuredToolHandler(s.HandleGetCashFlowAnalysis))
	mcpSrv.AddTool(s.GetRecurringTransactionsTool(), mcp.NewStructuredToolHandler(s.HandleGetRecurringTransactions))
	mcpSrv.AddTool(s.GetDataSharingConsentsTool(), mcp.NewStructuredToolHandler(s.HandleGetDataSharingConsents))
```

- [ ] **Step 5: Rodar os testes e confirmar que passam**

Run: `cd open-finance && go test ./... && go build ./...`
Expected: PASS em todos os pacotes.

- [ ] **Step 6: Verificar as sete ferramentas no servidor em execução**

```bash
cd open-finance && ENV=local MCP_PORT=8082 go run ./cmd/mcp &
sleep 3
curl -s -X POST http://localhost:8082/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | grep -o '"name":"[a-z_]*"'
kill %1
```

Expected: as sete linhas — `get_customer_by_id`, `get_customer_by_document`, `get_all_customers`, `get_bank_statements`, `get_cash_flow_analysis`, `get_recurring_transactions`, `get_data_sharing_consents`.

- [ ] **Step 7: Commit**

```bash
git add open-finance/internal/interfaces/mcp/
git commit -m "feat: ferramentas de extrato, fluxo de caixa, recorrências e consentimento"
```

---

### Task 7: Dados sintéticos alinhados ao birô

Substitui os dez clientes placeholder pelos dez do birô e gera os dados de Open Finance coerentes com cada perfil de crédito.

**Files:**
- Create: `open-finance/cmd/fixtures/generate.py`
- Modify: `open-finance/cmd/fixtures/fixtures/personal_informations.json`
- Modify: `open-finance/cmd/fixtures/fixtures/persons.json`
- Modify: `open-finance/cmd/fixtures/fixtures/addresses.json`
- Modify: `open-finance/cmd/fixtures/fixtures/person_addresses.json`
- Modify: `open-finance/cmd/fixtures/fixtures/person_documents.json`
- Modify: `open-finance/cmd/fixtures/fixtures/files.json`
- Modify: `open-finance/cmd/fixtures/fixtures/bank_account_profiles.json`
- Modify: `open-finance/cmd/fixtures/fixtures/bank_statements.json`
- Modify: `open-finance/cmd/fixtures/fixtures/cash_flow_analyses.json`
- Modify: `open-finance/cmd/fixtures/fixtures/recurring_transactions.json`
- Modify: `open-finance/cmd/fixtures/fixtures/data_sharing_consents.json`

**Interfaces:**
- Consumes: o esquema da Task 1; as fixtures do birô como fonte dos dados de identificação.
- Produces: dez clientes com os mesmos `id`, nomes e CPFs do birô, e as cinco fixtures de Open Finance povoadas.

- [ ] **Step 1: Copiar as fixtures de identificação do birô**

```bash
cd open-finance/cmd/fixtures/fixtures
for f in files.json addresses.json personal_informations.json person_addresses.json person_documents.json; do
  cp ../../../../bureau/cmd/fixtures/fixtures/$f .
done
```

- [ ] **Step 2: Escrever o gerador**

Criar `open-finance/cmd/fixtures/generate.py`. O script é determinístico — sem aleatoriedade — para que regerar as fixtures produza exatamente o mesmo arquivo:

```python
#!/usr/bin/env python3
"""Gera as fixtures de Open Finance a partir dos perfis definidos abaixo.

Cada perfil é derivado do quadro de crédito do mesmo cliente no birô, de modo
que as duas fontes contem uma história coerente. Rode a partir de
cmd/fixtures/: python3 generate.py
"""

import json
from datetime import datetime, timedelta
from pathlib import Path

OUT = Path(__file__).parent / "fixtures"
TS = "2026-08-31T00:00:00Z"

INSTITUTIONS = {
    "alfa": ("Banco Sintético Alfa", "11222333000181"),
    "beta": ("Banco Sintético Beta", "22333444000172"),
    "gama": ("Banco Sintético Gama", "33444555000163"),
    "delta": ("Banco Sintético Delta", "44555666000154"),
}

# person_id: perfil. "renda" é a renda declarada no birô, usada só para
# documentar a coerência entre as fontes; não vai para nenhuma fixture.
PROFILES = {
    1: {  # Felipe Pereira Santos — score 254, DTI 1,10, 0/7 pagamentos em dia
        "renda": 3690,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 2, "age": 34, "chk": True, "sav": False, "inv": False, "value": None},
        "cash": {"in": 3550.00, "out": 4300.00, "vol": 0.38, "neg_days": 22, "recurring": False},
        "accounts": [("alfa", "checking", 210.00)],
        "recurring": [
            ("expense", "rent", "Aluguel residencial", 1200.00, "Imobiliária Sintética", True),
            ("expense", "utility", "Energia elétrica", 320.00, "Concessionária Sintética", True),
            ("income", "salary", "Salário", 1750.00, "Empregador Sintético", False),
        ],
    },
    2: {  # Henrique Martins Barbosa — score 726, DTI 0,31, 6/6 em dia
        "renda": 7437,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 3, "age": 71, "chk": True, "sav": True, "inv": True, "value": 18500.00},
        "cash": {"in": 7600.00, "out": 6100.00, "vol": 0.08, "neg_days": 0, "recurring": True},
        "accounts": [("beta", "checking", 4200.00), ("alfa", "savings", 9800.00)],
        "recurring": [
            ("income", "salary", "Salário", 6950.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 1900.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Serviço de streaming", 89.00, "Streaming Sintético", True),
        ],
    },
    3: {  # Fernanda Costa Barbosa — score 448, 2 negativações ativas
        "renda": 2477,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 2, "age": 22, "chk": True, "sav": False, "inv": False, "value": None},
        "cash": {"in": 2380.00, "out": 2610.00, "vol": 0.31, "neg_days": 14, "recurring": True},
        "accounts": [("alfa", "checking", 640.00)],
        "recurring": [
            ("income", "salary", "Salário", 2150.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 950.00, "Imobiliária Sintética", True),
            ("expense", "utility", "Água e esgoto", 210.00, "Concessionária Sintética", True),
        ],
    },
    4: {  # Fernanda Rodrigues Pereira — score 349, alerta de fraude
        "renda": 3836,
        "consent": {"bank": "gama", "status": "granted"},
        "profile": {"rel": 3, "age": 15, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 4100.00, "out": 3980.00, "vol": 0.55, "neg_days": 9, "recurring": False},
        "accounts": [("gama", "checking", 880.00), ("beta", "savings", 1500.00)],
        "recurring": [
            ("expense", "rent", "Aluguel residencial", 1100.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de telefonia", 49.00, "Telecom Sintética", True),
        ],
    },
    5: {  # Gabriela Ribeiro Barbosa — score 809 com fluxo apertado: fontes discordam
        "renda": 6915,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 4, "age": 88, "chk": True, "sav": True, "inv": True, "value": 42000.00},
        "cash": {"in": 7100.00, "out": 6980.00, "vol": 0.19, "neg_days": 7, "recurring": True},
        "accounts": [("beta", "checking", 3100.00), ("gama", "savings", 5400.00)],
        "recurring": [
            ("income", "salary", "Salário", 6400.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2400.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento de veículo", 1850.00, "Banco Sintético Beta", True),
            ("expense", "subscription", "Serviço de streaming", 89.00, "Streaming Sintético", True),
        ],
    },
    6: {  # Henrique Almeida Ribeiro — renda indeterminável e consentimento revogado
        "renda": None,
        "consent": {"bank": "alfa", "status": "revoked"},
        "profile": None,
        "cash": None,
        "accounts": [],
        "recurring": [],
    },
    7: {  # Igor Souza Martins — score 931, DTI 0,21
        "renda": 16649,
        "consent": {"bank": "delta", "status": "granted"},
        "profile": {"rel": 4, "age": 132, "chk": True, "sav": True, "inv": True, "value": 310000.00},
        "cash": {"in": 17200.00, "out": 12100.00, "vol": 0.06, "neg_days": 0, "recurring": True},
        "accounts": [("delta", "checking", 12400.00), ("beta", "savings", 48000.00)],
        "recurring": [
            ("income", "salary", "Salário", 15800.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 3800.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de saúde", 129.00, "Operadora Sintética", True),
        ],
    },
    8: {  # Lucas Martins Souza — renda declarada 12.043 contra 4.200 detectados
        "renda": 12043,
        "consent": {"bank": "gama", "status": "granted"},
        "profile": {"rel": 2, "age": 41, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 5900.00, "out": 5750.00, "vol": 0.44, "neg_days": 11, "recurring": True},
        "accounts": [("gama", "checking", 1900.00)],
        "recurring": [
            ("income", "salary", "Salário", 4200.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 1750.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento imobiliário", 2100.00, "Banco Sintético Gama", True),
        ],
    },
    9: {  # Eduardo Barbosa Almeida — score 663, DTI 0,39
        "renda": 10529,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 3, "age": 57, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 10800.00, "out": 9500.00, "vol": 0.14, "neg_days": 3, "recurring": True},
        "accounts": [("beta", "checking", 5200.00), ("alfa", "savings", 7300.00)],
        "recurring": [
            ("income", "salary", "Salário", 9600.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2600.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento de veículo", 1500.00, "Banco Sintético Beta", True),
        ],
    },
    10: {  # Eduardo Ribeiro Ribeiro — score 800, DTI 0,26
        "renda": 7958,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 3, "age": 96, "chk": True, "sav": True, "inv": True, "value": 61000.00},
        "cash": {"in": 8300.00, "out": 6600.00, "vol": 0.09, "neg_days": 0, "recurring": True},
        "accounts": [("alfa", "checking", 6100.00), ("delta", "savings", 22000.00)],
        "recurring": [
            ("income", "salary", "Salário", 7500.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2100.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de telefonia", 45.00, "Telecom Sintética", True),
        ],
    },
}

# Os três meses cobertos pelos extratos, encerrando em 2026-08-31.
MONTHS = [
    ("2026-06-01T00:00:00Z", "2026-06-30T00:00:00Z"),
    ("2026-07-01T00:00:00Z", "2026-07-31T00:00:00Z"),
    ("2026-08-01T00:00:00Z", "2026-08-31T00:00:00Z"),
]


def base(row_id):
    return {"id": row_id, "created_at": TS, "updated_at": TS, "deleted_at": None}


def build():
    profiles, statements, analyses, recurrences, consents = [], [], [], [], []
    profile_by_person = {}
    next_id = {"profile": 1, "stmt": 1, "cash": 1, "rec": 1, "consent": 1}

    for person_id in sorted(PROFILES):
        p = PROFILES[person_id]

        if p["profile"]:
            pid = next_id["profile"]
            next_id["profile"] += 1
            profile_by_person[person_id] = pid
            profiles.append({
                **base(pid),
                "person_id": person_id,
                "profile_date": "2026-08-31T00:00:00Z",
                "banking_relationships": p["profile"]["rel"],
                "account_age_average": p["profile"]["age"],
                "has_checking_account": p["profile"]["chk"],
                "has_savings_account": p["profile"]["sav"],
                "has_investment_account": p["profile"]["inv"],
                "investments_value": p["profile"]["value"],
            })

        if p["cash"]:
            c = p["cash"]
            net = round(c["in"] - c["out"], 2)
            cid = next_id["cash"]
            next_id["cash"] += 1
            analyses.append({
                **base(cid),
                "person_id": person_id,
                "analysis_date": "2026-08-31T00:00:00Z",
                "period_days": 90,
                "average_monthly_inflow": c["in"],
                "average_monthly_outflow": c["out"],
                "net_cash_flow": net,
                "inflow_volatility": c["vol"],
                "negative_balance_days": c["neg_days"],
                "has_recurring_income": c["recurring"],
            })

        for bank, account_type, opening in p["accounts"]:
            name, cnpj = INSTITUTIONS[bank]
            balance = opening
            for period_start, period_end in MONTHS:
                if account_type == "checking":
                    credits, debits = p["cash"]["in"], p["cash"]["out"]
                    count = 62
                else:
                    # Poupança: movimentação enxuta, com aporte mensal modesto.
                    credits, debits = 400.00, 200.00
                    count = 4
                closing = round(balance + credits - debits, 2)
                sid = next_id["stmt"]
                next_id["stmt"] += 1
                statements.append({
                    **base(sid),
                    "person_id": person_id,
                    "institution": name,
                    "institution_document": cnpj,
                    "account_type": account_type,
                    "period_start": period_start,
                    "period_end": period_end,
                    "opening_balance": balance,
                    "closing_balance": closing,
                    "total_credits": credits,
                    "total_debits": debits,
                    "transaction_count": count,
                    "currency": "BRL",
                })
                balance = closing

        for kind, category, description, amount, counterparty, active in p["recurring"]:
            rid = next_id["rec"]
            next_id["rec"] += 1
            recurrences.append({
                **base(rid),
                "person_id": person_id,
                "transaction_type": kind,
                "category": category,
                "description": description,
                "amount": amount,
                "frequency": "monthly",
                "counterparty": counterparty,
                "first_detected_date": "2025-09-05T00:00:00Z",
                "last_occurrence_date": "2026-08-05T00:00:00Z" if active else "2026-03-05T00:00:00Z",
                "is_active": active,
            })

        c = p["consent"]
        name, _ = INSTITUTIONS[c["bank"]]
        cid = next_id["consent"]
        next_id["consent"] += 1
        revoked = c["status"] == "revoked"
        consents.append({
            **base(cid),
            "person_id": person_id,
            "consent_id": f"urn:openfinance:consent:{person_id:03d}",
            "institution": name,
            "status": c["status"],
            "scope": json.dumps(["ACCOUNTS_READ", "ACCOUNTS_BALANCES_READ", "RESOURCES_READ"]),
            "granted_at": "2026-01-15T00:00:00Z" if revoked else "2026-05-10T00:00:00Z",
            "expires_at": None if revoked else "2027-05-10T00:00:00Z",
            "revoked_at": "2026-04-02T00:00:00Z" if revoked else None,
        })

    return profile_by_person, {
        "bank_account_profiles.json": profiles,
        "bank_statements.json": statements,
        "cash_flow_analyses.json": analyses,
        "recurring_transactions.json": recurrences,
        "data_sharing_consents.json": consents,
    }


def rewrite_persons(profile_by_person):
    """Reescreve persons.json trocando os vínculos do birô pelo perfil bancário."""
    path = OUT / "persons.json"
    rows = json.loads(path.read_text())
    for row in rows:
        row.pop("credit_score_id", None)
        row.pop("financial_profile_id", None)
        row["bank_account_profile_id"] = profile_by_person.get(row["id"])
    path.write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")


def main():
    profile_by_person, files = build()
    for name, rows in files.items():
        (OUT / name).write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")
        print(f"{name}: {len(rows)} registros")
    rewrite_persons(profile_by_person)
    print("persons.json: vínculos reescritos")


if __name__ == "__main__":
    main()
```

- [ ] **Step 3: Copiar `persons.json` do birô e rodar o gerador**

```bash
cd open-finance/cmd/fixtures
cp ../../../bureau/cmd/fixtures/fixtures/persons.json fixtures/persons.json
python3 generate.py
```

Expected: cinco linhas de contagem — `bank_account_profiles.json: 9 registros`, `bank_statements.json: 45 registros`, `cash_flow_analyses.json: 9 registros`, `recurring_transactions.json: 27 registros`, `data_sharing_consents.json: 10 registros` — mais `persons.json: vínculos reescritos`.

- [ ] **Step 4: Carregar as fixtures no banco**

```bash
cd open-finance
make migration-down && make migration-up
go run ./cmd/fixtures -dir ./cmd/fixtures/fixtures -truncate
```

Expected: uma linha `inserted N rows into <tabela>` por fixture, sem erro. As tabelas `tokens` e `sessions` podem aparecer como `skipping` se estiverem vazias.

- [ ] **Step 5: Conferir a coerência dos dados carregados**

```bash
docker compose exec postgres psql -U "$DB_USERNAME" -d open-finance -c "
SELECT p.id, pi.full_name, pi.document,
       c.status AS consentimento,
       cf.net_cash_flow, cf.negative_balance_days,
       COALESCE(SUM(rt.amount) FILTER (WHERE rt.transaction_type = 'income' AND rt.is_active), 0) AS renda_recorrente
FROM persons p
JOIN personal_informations pi ON pi.id = p.personal_information_id
LEFT JOIN data_sharing_consents c ON c.person_id = p.id
LEFT JOIN cash_flow_analyses cf ON cf.person_id = p.id
LEFT JOIN recurring_transactions rt ON rt.person_id = p.id
GROUP BY p.id, pi.full_name, pi.document, c.status, cf.net_cash_flow, cf.negative_balance_days
ORDER BY p.id;"
```

Verificar, linha a linha:
- os dez nomes e CPFs são os do birô (Felipe Pereira Santos, CPF `83609091800`, na primeira linha)
- cliente 1 com `net_cash_flow` negativo e 22 dias de saldo negativo
- cliente 6 com consentimento `revoked` e as demais colunas nulas
- cliente 8 com renda recorrente `4200.00`, contra renda declarada de 12.043 no birô
- clientes 2, 7 e 10 com `net_cash_flow` positivo e 0 dias de saldo negativo

- [ ] **Step 6: Consultar uma ferramenta contra os dados reais**

```bash
cd open-finance && ENV=local MCP_PORT=8082 go run ./cmd/mcp &
sleep 3
curl -s -X POST http://localhost:8082/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_recurring_transactions","arguments":{"document":"18979021232","transaction_type":"income"}}}'
kill %1
```

Expected: resposta com `"customer_id": 8` e um item de `"amount": 4200`.

- [ ] **Step 7: Commit**

```bash
git add open-finance/cmd/fixtures/
git commit -m "feat: dados sintéticos de Open Finance alinhados aos dez clientes do birô"
```

---

### Task 8: Política de crédito v1.1

**Files:**
- Create: `docs/politica_credito_agente_v1.1.md`
- Modify: `bureau/docs/politica_credito_agente.md` (apenas a nota de versão no topo)

**Interfaces:**
- Consumes: os dados produzidos pela Task 7 e as ferramentas das Tasks 5 e 6.
- Produces: o documento de instruções que acompanha o prompt do agente na fase integrada.

- [ ] **Step 1: Marcar a v1.0 como política do experimento preliminar**

Substituir, no topo de `bureau/docs/politica_credito_agente.md`, a linha de versão por:

```markdown
> **Versão:** 1.0 · **Vigência:** 2026-07-02 · **Moeda:** BRL
> **Situação:** política do experimento preliminar, conduzido apenas com o
> servidor de birô de crédito. Para a fase integrada, com o servidor de Open
> Finance conectado, use `docs/politica_credito_agente_v1.1.md`.
```

- [ ] **Step 2: Escrever a v1.1**

Criar `docs/politica_credito_agente_v1.1.md` copiando integralmente a v1.0 e aplicando as alterações abaixo. Copiar antes de editar preserva as seções que não mudam (decisões possíveis, formato geral, exemplo resolvido).

**Cabeçalho:**

```markdown
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
```

**§1, acrescentar após o item 2:**

```markdown
2-A. Para analisar um cliente, recupere o cadastro do **birô** e os dados de
   **Open Finance**. As duas fontes expõem ferramentas de mesmo nome; identifique
   a origem pelo conector. Registre em `fontes_consultadas` quais servidores
   foram efetivamente consultados.
2-B. A ausência de consentimento ativo de compartilhamento **não** impede a
   análise. Os critérios de Open Finance passam a ser pontuados pelas regras de
   dado ausente da §8.
```

**§4, acrescentar como última regra da tabela:**

```markdown
| K8 | Divergência de renda entre fontes | Havendo consentimento com `status = "granted"` e não expirado **e** ao menos uma `recurring_transaction` com `transaction_type = "income"` e `is_active = true`, se a renda considerada (birô) > **2 ×** a soma dos `amount` dessas receitas | `ANALISE_MANUAL` |
```

Acrescentar logo abaixo da tabela:

```markdown
A dupla condição de K8 é necessária: sem consentimento ativo ou sem receita
recorrente detectada, a soma é zero e a regra dispararia para todo cliente sem
dados de Open Finance. Nesses casos, a ausência é tratada por C8, não por K8.
```

**§5, substituir a tabela do scorecard inteira por:**

```markdown
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
```

**§6, substituir o parágrafo de degradação por:**

```markdown
**Regra de degradação por dados ausentes:** se **2 ou mais** critérios do
scorecard tiverem sido pontuados por regra de nulo (C2, C3, C5, C7, C8 ou C9),
a decisão `APROVADO` é rebaixada para `APROVADO_COM_RESSALVAS`.
```

**§7, substituir a fórmula por:**

```markdown
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
```

**§8, acrescentar à tabela de consolidação:**

```markdown
| Sem consentimento `granted` e não expirado | C7 → 4 pontos, C8 → 3 pontos, C9 → 1 ponto | C7/C8/C9 |
| Consentimento ativo, sem `cash_flow_analyses` | C7 → 4 pontos, C9 → 1 ponto | C7/C9 |
| Consentimento ativo, sem receita recorrente | C8 → 0 pontos (ausência de ocorrência, não é dado faltante) | C8 |
| `inflow_volatility` nulo | Ignorar; não integra o scorecard | — |
```

**§9, no bloco JSON:** trocar `"politica_versao": "1.0"` por `"1.1"`, acrescentar
as três chaves de critério e o array de fontes:

```json
    "C7_fluxo_caixa":          { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C8_renda_recorrente":     { "valor_observado": null, "pontos": 0, "dado_ausente": false },
    "C9_dias_saldo_negativo":  { "valor_observado": null, "pontos": 0, "dado_ausente": false }
```

```json
  "fontes_consultadas": ["bureau-mcp", "open-finance-mcp"],
```

**§10:** substituir o exemplo resolvido, recalculando-o sob os pesos da v1.1:

```markdown
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
```

- [ ] **Step 3: Conferir que os pesos somam 100**

```bash
grep -oE '^\| \*\*C[0-9]\.[^|]*\| [0-9]+ ' docs/politica_credito_agente_v1.1.md \
  | grep -oE '[0-9]+ $' | paste -sd+ - | bc
```

Expected: `100`.

- [ ] **Step 4: Conferir a coerência interna do documento**

Ler a v1.1 do início ao fim confirmando que:
- nenhuma seção ainda cita os pesos antigos (35, 20, 20, 15, 5, 5)
- as nove chaves de critério da §9 batem com as nove linhas da tabela da §5
- K8 aparece na §4 e na §8, e em nenhuma outra regra
- o exemplo da §10 usa os pesos da v1.1

- [ ] **Step 5: Commit**

```bash
git add docs/politica_credito_agente_v1.1.md bureau/docs/politica_credito_agente.md
git commit -m "docs: política de crédito v1.1 com critérios de Open Finance"
```

---

### Task 9: Infraestrutura

Coloca o servidor no Compose local e no provisionamento de nuvem. Nada é publicado: os arquivos ficam versionados, prontos para o autor aplicar.

**Files:**
- Modify: `docker-compose.yaml`
- Modify: `infra/main.tf`
- Modify: `infra/variables.tf`
- Modify: `infra/secrets.tf`
- Modify: `infra/outputs.tf`
- Create: `.github/workflows/deploy-open-finance.yml`

**Interfaces:**
- Consumes: o `Dockerfile` já existente em `open-finance/`, com o `target: mcp`.
- Produces: serviço Compose `open-finance-mcp` em `localhost:8082`; recursos Terraform do serviço `open-finance-mcp`.

- [ ] **Step 1: Acrescentar o serviço ao Compose**

Em `docker-compose.yaml`, logo após o bloco `bureau-mcp`:

```yaml
  open-finance-mcp:
    build:
      context: ./open-finance
      dockerfile: Dockerfile
      target: mcp
    ports:
      - "8082:8080"
    networks:
      - interservice-network
      - postgres-network
    env_file:
      - ./open-finance/.env
    depends_on:
      postgres:
        condition: service_healthy
```

- [ ] **Step 2: Verificar que o serviço sobe**

```bash
docker compose up -d --build open-finance-mcp
sleep 5
curl -s -o /dev/null -w '%{http_code}\n' -X POST http://localhost:8082/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

Expected: `200`.

- [ ] **Step 3: Parametrizar o Terraform por serviço**

Em `infra/variables.tf`, substituir a variável `image` por um mapa e acrescentar o mapa de URLs de banco:

```hcl
variable "images" {
  description = "Imagem de contêiner por serviço MCP"
  type        = map(string)
}

variable "database_urls" {
  description = "DATABASE_URL por serviço MCP"
  type        = map(string)
  sensitive   = true
}
```

Em `infra/main.tf`, trocar os recursos de serviço único por versões iteradas. O Artifact Registry, a service account, o serviço Cloud Run e o `iam_member` passam a usar `for_each` sobre o mesmo conjunto:

```hcl
locals {
  mcp_services = toset(["bureau-mcp", "open-finance-mcp"])
}

resource "google_artifact_registry_repository" "repo" {
  for_each = local.mcp_services

  location      = var.region
  repository_id = each.key
  format        = "DOCKER"

  depends_on = [google_project_service.services]
}

resource "google_service_account" "mcp_sa" {
  for_each = local.mcp_services

  account_id   = "${each.key}-sa"
  display_name = "MCP Server: ${each.key}"
}

resource "google_project_iam_member" "secret_accessor" {
  for_each = local.mcp_services

  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.mcp_sa[each.key].email}"
}

resource "google_cloud_run_v2_service" "mcp" {
  for_each = local.mcp_services

  name     = each.key
  location = var.region

  template {
    service_account = google_service_account.mcp_sa[each.key].email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    containers {
      image = var.images[each.key]

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        cpu_idle = false
      }

      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_url[each.key].secret_id
            version = "latest"
          }
        }
      }
    }
  }

  depends_on = [
    google_project_service.services,
    google_project_iam_member.secret_accessor,
  ]
}

resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  for_each = local.mcp_services

  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.mcp[each.key].name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
```

Em `infra/secrets.tf`, aplicar o mesmo `for_each` ao secret:

```hcl
resource "google_secret_manager_secret" "database_url" {
  for_each = local.mcp_services

  secret_id = "${each.key}-database-url"

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "database_url" {
  for_each = local.mcp_services

  secret      = google_secret_manager_secret.database_url[each.key].id
  secret_data = var.database_urls[each.key]
}
```

O `secret_id` do birô permanece `bureau-mcp-database-url` — só o endereço no estado do Terraform muda, o recurso no GCP é o mesmo. Ele entra na migração de estado do Step 4 junto com os demais.

Em `infra/outputs.tf`, substituir os dois outputs, que hoje apontam para recursos de nome fixo:

```hcl
output "service_urls" {
  description = "URL pública de cada servidor MCP no Cloud Run"
  value       = { for k, s in google_cloud_run_v2_service.mcp : k => s.uri }
}

output "artifact_registry_repos" {
  description = "Endereço do repositório de imagens de cada servidor MCP"
  value = {
    for k in local.mcp_services : k => "${var.region}-docker.pkg.dev/${var.project_id}/${k}/${k}"
  }
}
```

Remover também a variável `mcp_auth_token` de `infra/variables.tf`: a autenticação foi retirada no commit `641e890` e nenhum recurso a consome.

- [ ] **Step 4: Validar o Terraform sem aplicar**

```bash
cd infra && terraform fmt -check && terraform validate
```

Expected: `Success! The configuration is valid.` Se `terraform validate` exigir inicialização, rodar `terraform init -backend=false` antes.

> O `for_each` renomeia os recursos do birô no estado (`google_cloud_run_v2_service.bureau_mcp` vira `google_cloud_run_v2_service.mcp["bureau-mcp"]`). Antes do primeiro `apply`, migrar o estado com `terraform state mv`, ou o Terraform destruirá e recriará o serviço do birô. Registrar isso no commit; a migração de estado fica a cargo do autor, junto com o deploy.

- [ ] **Step 5: Criar o workflow de deploy**

Criar `.github/workflows/deploy-open-finance.yml`, espelhando `deploy-bureau.yml`:

```yaml
name: Deploy open-finance-mcp

on:
  push:
    branches: [main]
    paths:
      - "open-finance/**"
      - "infra/**"
      - ".github/workflows/deploy-open-finance.yml"

env:
  PROJECT_ID: ${{ secrets.GCP_PROJECT_ID }}
  REGION: southamerica-east1
  IMAGE: southamerica-east1-docker.pkg.dev/${{ secrets.GCP_PROJECT_ID }}/open-finance-mcp/open-finance-mcp

jobs:
  deploy:
    name: Build & Deploy
    runs-on: ubuntu-latest
    permissions:
      contents: read
      id-token: write

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.GCP_WORKLOAD_IDENTITY_PROVIDER }}
          service_account: ${{ secrets.GCP_SERVICE_ACCOUNT }}

      - name: Set up gcloud
        uses: google-github-actions/setup-gcloud@v2

      - name: Configure Docker
        run: gcloud auth configure-docker ${{ env.REGION }}-docker.pkg.dev --quiet

      - name: Build
        run: |
          docker build \
            --target mcp \
            --cache-from ${{ env.IMAGE }}:latest \
            --build-arg BUILDKIT_INLINE_CACHE=1 \
            --tag ${{ env.IMAGE }}:${{ github.sha }} \
            --tag ${{ env.IMAGE }}:latest \
            ./open-finance

      - name: Push
        run: |
          docker push ${{ env.IMAGE }}:${{ github.sha }}
          docker push ${{ env.IMAGE }}:latest

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3

      - name: Terraform Init
        working-directory: infra/
        run: terraform init

      - name: Terraform Apply
        working-directory: infra/
        run: |
          terraform apply -auto-approve \
            -var="project_id=${{ env.PROJECT_ID }}" \
            -var="region=${{ env.REGION }}" \
            -var='images={"bureau-mcp":"southamerica-east1-docker.pkg.dev/${{ env.PROJECT_ID }}/bureau-mcp/bureau-mcp:latest","open-finance-mcp":"${{ env.IMAGE }}:${{ github.sha }}"}' \
            -var='database_urls={"bureau-mcp":"${{ secrets.DATABASE_URL }}","open-finance-mcp":"${{ secrets.OPEN_FINANCE_DATABASE_URL }}"}'
```

- [ ] **Step 6: Ajustar o workflow do birô ao novo formato de variáveis**

Em `.github/workflows/deploy-bureau.yml`, substituir as duas linhas `-var="image=…"` e `-var="database_url=…"` pelos mapas equivalentes, para que os dois workflows falem a mesma linguagem com o Terraform:

```yaml
            -var='images={"bureau-mcp":"${{ env.IMAGE }}:${{ github.sha }}","open-finance-mcp":"southamerica-east1-docker.pkg.dev/${{ env.PROJECT_ID }}/open-finance-mcp/open-finance-mcp:latest"}' \
            -var='database_urls={"bureau-mcp":"${{ secrets.DATABASE_URL }}","open-finance-mcp":"${{ secrets.OPEN_FINANCE_DATABASE_URL }}"}'
```

Remover a linha `-var="mcp_auth_token=${{ secrets.MCP_AUTH_TOKEN }}"`, já que a variável sai de `variables.tf` no Step 3.

- [ ] **Step 7: Rodar a suíte inteira uma última vez**

```bash
cd open-finance && go build ./... && go test ./... && go vet ./...
cd ../infra && terraform fmt -check && terraform validate
```

Expected: build limpo, testes passando, `terraform validate` bem-sucedido.

- [ ] **Step 8: Commit**

```bash
git add docker-compose.yaml infra/ .github/workflows/
git commit -m "chore: provisionamento local e de nuvem do servidor de Open Finance"
```

---

## Pendências para o autor

Fora do escopo deste plano, mas necessárias antes do próximo item do cronograma:

1. **Segredo `OPEN_FINANCE_DATABASE_URL`** no repositório do GitHub, apontando para a base Supabase do Open Finance.
2. **Migração do estado do Terraform** (`terraform state mv`) antes do primeiro `apply` com `for_each`, para não recriar o serviço do birô.
3. **Conectar o servidor ao Claude Web** como segundo conector, depois do deploy.
4. **Derivar o gabarito** dos dez clientes sob a v1.1 — atividade de setembro no cronograma.
5. **Atualizar o TCC**: a tabela de critérios do capítulo 3 (pesos da v1.1) e a descrição das sete ferramentas; o capítulo 6 quando houver resultados.
