# Servidor MCP de Registro Interno — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Tornar o servidor `internal-registry/` funcional e integrável como terceiro conector do protótipo, expondo cinco dimensões de dados internos da instituição e passando a política de crédito para a v1.2.

**Architecture:** Mesma arquitetura em camadas do `open-finance/` (entidades Gorm → repositórios por dimensão → serviço orquestrador → ferramentas MCP), replicada 1:1 para as dimensões `CustomerRelationship`, `ContractedProduct`, `InternalPaymentRecord`, `PreApprovedLimit` e `IncomeDeclaration`. Entidades e DTOs já existem; este plano fecha as lacunas L1–L7.

**Tech Stack:** Go 1.25, `gorm.io/gorm` v1.31.1 (API genérica `gorm.G[T]`), `gorm.io/driver/postgres` v1.6.0, `github.com/mark3labs/mcp-go`, `github.com/BrunoPolaski/go-rest-err/rest_err`, `golang-migrate` (CLI), Terraform (GCP Cloud Run), Docker Compose.

**Spec:** `docs/superpowers/specs/2026-09-01-internal-registry-mcp-server-design.md`

## Global Constraints

- Module path: `github.com/BrunoPolaski/internal-registry`. Todos os imports usam esse prefixo.
- Padrão de repositório: `gorm.G[T](db).Where(...).Find(ctx)`; erros de banco viram `rest_err.NewInternalServerError(msg).WithCause(err)`; ausência de pessoa vira `rest_err.NewNotFoundError`.
- Filtros opcionais: string vazia (`""`) = não filtrar; ponteiro `nil` = não filtrar; `bool` conforme parâmetro.
- Identificação do cliente nas tools por dimensão: `customer_id` (int) **OU** `document` (string), exatamente um. A recusa de referência ambígua/vazia vive no serviço (`resolveCustomer`), espelhando `open-finance/internal/services/open_finance_service.go`.
- Descrições das tools declaram **"Internal Registry:"** na primeira linha.
- Novas dependências: **nenhuma**. Só o que já está no `go.mod`.
- Comentários e mensagens em português, no estilo dos arquivos existentes.
- Nenhum `git push` e nenhum deploy. Commits locais a cada task.
- Porta local do serviço: `8083:8080` (8081 = bureau, 8082 = open-finance).
- Fixtures: arrays JSON, chaves snake_case, timestamps RFC3339, `deleted_at: null` em linhas ativas.
- Política v1.2: total do scorecard **exatamente 100**; bandas de decisão inalteradas.

**Referência de cópia:** o servidor `open-finance/` é o template exato. Cada task cita o arquivo análogo. Ao copiar, trocar o prefixo do módulo e os nomes de tipo/dimensão.

---

### Task 1: Reescrever as migrations (L1)

**Files:**
- Modify: `internal-registry/internal/infra/thirdparty/database/migrations/000001_init.up.sql` (reescrever por completo)
- Modify: `internal-registry/internal/infra/thirdparty/database/migrations/000001_init.down.sql` (reescrever por completo)
- Não tocar em `000002_tokens.{up,down}.sql`.

**Interfaces:**
- Produces: tabelas `customer_relationships`, `contracted_products`, `internal_payment_records`, `pre_approved_limits` e `persons` com `customer_relationship_id`. As dimensões consumidas pelos repositórios da Task 2.

- [ ] **Step 1: Reescrever `000001_init.up.sql`**

Manter as tabelas compartilhadas idênticas às atuais (`api_keys`, `files`, `personal_informations`, `addresses`, `person_addresses`, `person_documents`, `admins`, `analysts`, `users`, `sessions`) e a `income_declarations` (que já corresponde à entidade). **Remover** todas as tabelas herdadas do birô sem entidade (`credit_scores`, `financial_profiles`, `credit_accounts`, `credit_inquiries`, `debts`, `payment_histories`, `negative_records`, `employment_records`, `legal_records`, `compliance_checks`, `fraud_alerts`, `risk_assessments`, `person_relationships`, `data_sources`, `person_data_sources`). Substituir a `persons` e acrescentar as quatro novas tabelas:

```sql
CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    customer_relationship_id BIGINT,
    last_verified_at TIMESTAMP
);

CREATE TABLE customer_relationships (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    customer_since TIMESTAMP NOT NULL,
    relationship_months INTEGER NOT NULL DEFAULT 0,
    segment VARCHAR(50),
    branch VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    churn_risk VARCHAR(50),
    internal_score INTEGER
);
CREATE UNIQUE INDEX idx_customer_relationships_person_id ON customer_relationships(person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_relationships_months ON customer_relationships(relationship_months);
CREATE INDEX idx_customer_relationships_active ON customer_relationships(is_active);
CREATE INDEX idx_customer_relationships_score ON customer_relationships(internal_score);

CREATE TABLE contracted_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    product_type VARCHAR(50) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    contract_number VARCHAR(100),
    contracted_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    balance DOUBLE PRECISION,
    monthly_value DOUBLE PRECISION
);
CREATE INDEX idx_contracted_products_person_id ON contracted_products(person_id);
CREATE INDEX idx_contracted_products_type ON contracted_products(product_type);
CREATE INDEX idx_contracted_products_number ON contracted_products(contract_number);
CREATE INDEX idx_contracted_products_status ON contracted_products(status);

CREATE TABLE internal_payment_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    contracted_product_id BIGINT REFERENCES contracted_products(id),
    reference_month TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    payment_date TIMESTAMP,
    amount_due DOUBLE PRECISION NOT NULL,
    amount_paid DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) NOT NULL,
    days_late INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_internal_payment_records_person_id ON internal_payment_records(person_id);
CREATE INDEX idx_internal_payment_records_product_id ON internal_payment_records(contracted_product_id);
CREATE INDEX idx_internal_payment_records_ref_month ON internal_payment_records(reference_month);
CREATE INDEX idx_internal_payment_records_status ON internal_payment_records(status);
CREATE INDEX idx_internal_payment_records_days_late ON internal_payment_records(days_late);

CREATE TABLE pre_approved_limits (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    product_type VARCHAR(50) NOT NULL,
    approved_amount DOUBLE PRECISION NOT NULL,
    interest_rate DOUBLE PRECISION,
    calculated_date TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    policy_version VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX idx_pre_approved_limits_person_id ON pre_approved_limits(person_id);
CREATE INDEX idx_pre_approved_limits_type ON pre_approved_limits(product_type);
CREATE INDEX idx_pre_approved_limits_valid_until ON pre_approved_limits(valid_until);
CREATE INDEX idx_pre_approved_limits_active ON pre_approved_limits(is_active);
```

Ordem de criação: tabelas compartilhadas → `persons` → `customer_relationships` → `contracted_products` → `internal_payment_records` (depende de contracted_products) → `pre_approved_limits` → `income_declarations` (manter a definição atual, que referencia `persons` e `files`).

- [ ] **Step 2: Reescrever `000001_init.down.sql`**

Derrubar na ordem inversa das dependências:

```sql
DROP TABLE IF EXISTS pre_approved_limits;
DROP TABLE IF EXISTS internal_payment_records;
DROP TABLE IF EXISTS contracted_products;
DROP TABLE IF EXISTS income_declarations;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS customer_relationships;
DROP TABLE IF EXISTS analysts;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS person_documents;
DROP TABLE IF EXISTS person_addresses;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS personal_informations;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS api_keys;
```

(A `tokens` fica com o `000002_tokens.down.sql`; não repetir aqui.)

- [ ] **Step 3: Aplicar num Postgres real e conferir as tabelas**

O Postgres local já roda (compartilhado com bureau/open-finance) na porta 5433. Criar o banco e aplicar as migrations:

```bash
cd /home/bruno/College/mcp-servers
docker compose up -d postgres
docker compose exec -T postgres psql -U root -c 'CREATE DATABASE "internal-registry";' || true
cd internal-registry
make migration-up   # usa golang-migrate com a URL do .env; garanta DB_NAME=internal-registry no .env local
```

Expected: `migration-up` aplica `000001` e `000002` sem erro.

- [ ] **Step 4: Verificar o esquema**

```bash
docker compose -f /home/bruno/College/mcp-servers/docker-compose.yaml exec -T postgres \
  psql -U root -d internal-registry -c '\dt'
```

Expected: aparecem `customer_relationships`, `contracted_products`, `internal_payment_records`, `pre_approved_limits`, `income_declarations`, `persons` e as compartilhadas; **não** aparecem `credit_scores`, `fraud_alerts`, `debts`, etc.

- [ ] **Step 5: Commit**

```bash
cd /home/bruno/College/mcp-servers
git add internal-registry/internal/infra/thirdparty/database/migrations/000001_init.up.sql \
        internal-registry/internal/infra/thirdparty/database/migrations/000001_init.down.sql
git commit -m "feat(internal-registry): migrations das cinco dimensoes internas"
```

---

### Task 2: Camada de dados — result DTOs, repositórios por dimensão e factory (L5)

**Files:**
- Create: `internal-registry/internal/infra/controllers/dto/internal_registry_result_dto.go`
- Create: `internal-registry/internal/infra/repositories/interfaces/customer_relationship_repository.go`
- Create: `internal-registry/internal/infra/repositories/interfaces/contracted_product_repository.go`
- Create: `internal-registry/internal/infra/repositories/interfaces/internal_payment_record_repository.go`
- Create: `internal-registry/internal/infra/repositories/interfaces/pre_approved_limit_repository.go`
- Create: `internal-registry/internal/infra/repositories/interfaces/income_declaration_repository.go`
- Create: `internal-registry/internal/infra/repositories/gorm_customer_relationship_repository.go`
- Create: `internal-registry/internal/infra/repositories/gorm_contracted_product_repository.go`
- Create: `internal-registry/internal/infra/repositories/gorm_internal_payment_record_repository.go`
- Create: `internal-registry/internal/infra/repositories/gorm_pre_approved_limit_repository.go`
- Create: `internal-registry/internal/infra/repositories/gorm_income_declaration_repository.go`
- Modify: `internal-registry/internal/infra/repositories/factory.go`
- Test: `internal-registry/internal/infra/repositories/contract_test.go`

**Interfaces:**
- Consumes: entidades da Task 1 e as já existentes.
- Produces (usadas pela Task 4):
  - `interfaces.CustomerRelationshipRepository.GetByPersonID(ctx, personID uint) (*entities.CustomerRelationship, *rest_err.RestErr)` — devolve `nil, nil` quando não há relacionamento (não-cliente).
  - `interfaces.ContractedProductRepository.GetByPersonID(ctx, personID uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr)`
  - `interfaces.InternalPaymentRecordRepository.GetByPersonID(ctx, personID uint, status string, productID *uint) ([]entities.InternalPaymentRecord, *rest_err.RestErr)`
  - `interfaces.PreApprovedLimitRepository.GetByPersonID(ctx, personID uint, onlyActive bool) ([]entities.PreApprovedLimit, *rest_err.RestErr)`
  - `interfaces.IncomeDeclarationRepository.GetByPersonID(ctx, personID uint, verifiedOnly bool) ([]entities.IncomeDeclaration, *rest_err.RestErr)`
  - Getters no factory: `CustomerRelationshipRepository()`, `ContractedProductRepository()`, `InternalPaymentRecordRepository()`, `PreApprovedLimitRepository()`, `IncomeDeclarationRepository()`.
  - DTOs de resultado: `CustomerRelationshipResultDTO`, `ContractedProductsResultDTO`, `InternalPaymentRecordsResultDTO`, `PreApprovedLimitsResultDTO`, `IncomeDeclarationsResultDTO`.

- [ ] **Step 1: Criar os result DTOs** (espelha `open-finance/.../dto/open_finance_result_dto.go`)

```go
package dto

// Embrulham as saídas das ferramentas por dimensão. O MCP valida o esquema
// contra um objeto de topo; os campos de identificação são metadado de
// rastreabilidade da consulta.

type CustomerRelationshipResultDTO struct {
	CustomerID   uint                     `json:"customer_id"`
	Document     string                   `json:"document"`
	Relationship *CustomerRelationshipDTO `json:"relationship"` // null quando não-cliente
}

type ContractedProductsResultDTO struct {
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []ContractedProductDTO `json:"items"`
}

type InternalPaymentRecordsResultDTO struct {
	CustomerID uint                       `json:"customer_id"`
	Document   string                     `json:"document"`
	Items      []InternalPaymentRecordDTO `json:"items"`
}

type PreApprovedLimitsResultDTO struct {
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []PreApprovedLimitDTO `json:"items"`
}

type IncomeDeclarationsResultDTO struct {
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []IncomeDeclarationDTO `json:"items"`
}
```

- [ ] **Step 2: Criar as cinco interfaces de repositório**

`customer_relationship_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type CustomerRelationshipRepository interface {
	// GetByPersonID devolve o relacionamento da pessoa, ou (nil, nil) quando ela
	// não é cliente da instituição.
	GetByPersonID(ctx context.Context, personID uint) (*entities.CustomerRelationship, *rest_err.RestErr)
}
```

`contracted_product_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type ContractedProductRepository interface {
	// GetByPersonID devolve os produtos contratados da pessoa. productType e
	// status vazios não filtram.
	GetByPersonID(ctx context.Context, personID uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr)
}
```

`internal_payment_record_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type InternalPaymentRecordRepository interface {
	// GetByPersonID devolve os registros de pagamento interno da pessoa, do mais
	// recente para o mais antigo. status vazio e productID nil não filtram.
	GetByPersonID(ctx context.Context, personID uint, status string, productID *uint) ([]entities.InternalPaymentRecord, *rest_err.RestErr)
}
```

`pre_approved_limit_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type PreApprovedLimitRepository interface {
	// GetByPersonID devolve os limites pré-aprovados da pessoa. onlyActive
	// restringe a is_active = true.
	GetByPersonID(ctx context.Context, personID uint, onlyActive bool) ([]entities.PreApprovedLimit, *rest_err.RestErr)
}
```

`income_declaration_repository.go`:

```go
package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type IncomeDeclarationRepository interface {
	// GetByPersonID devolve as declarações de renda da pessoa, da mais recente
	// para a mais antiga. verifiedOnly restringe a verified = true.
	GetByPersonID(ctx context.Context, personID uint, verifiedOnly bool) ([]entities.IncomeDeclaration, *rest_err.RestErr)
}
```

- [ ] **Step 3: Criar as cinco implementações Gorm** (espelham `gorm_recurring_transaction_repository.go`)

`gorm_customer_relationship_repository.go` (note o tratamento de ausência):

```go
package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormCustomerRelationshipRepository struct {
	db *gorm.DB
}

func NewGormCustomerRelationshipRepository(db *gorm.DB) interfaces.CustomerRelationshipRepository {
	return &gormCustomerRelationshipRepository{db: db}
}

func (g *gormCustomerRelationshipRepository) GetByPersonID(ctx context.Context, personID uint) (*entities.CustomerRelationship, *rest_err.RestErr) {
	rel, err := gorm.G[entities.CustomerRelationship](g.db).Where("person_id = ?", personID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // não-cliente: ausência é dado válido
		}
		return nil, rest_err.NewInternalServerError("error while fetching customer relationship").WithCause(err)
	}
	return &rel, nil
}
```

`gorm_contracted_product_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormContractedProductRepository struct {
	db *gorm.DB
}

func NewGormContractedProductRepository(db *gorm.DB) interfaces.ContractedProductRepository {
	return &gormContractedProductRepository{db: db}
}

func (g *gormContractedProductRepository) GetByPersonID(ctx context.Context, personID uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr) {
	query := gorm.G[entities.ContractedProduct](g.db).Where("person_id = ?", personID)
	if productType != "" {
		query = query.Where("product_type = ?", productType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	products, err := query.Order("contracted_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching contracted products").WithCause(err)
	}
	return products, nil
}
```

`gorm_internal_payment_record_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormInternalPaymentRecordRepository struct {
	db *gorm.DB
}

func NewGormInternalPaymentRecordRepository(db *gorm.DB) interfaces.InternalPaymentRecordRepository {
	return &gormInternalPaymentRecordRepository{db: db}
}

func (g *gormInternalPaymentRecordRepository) GetByPersonID(ctx context.Context, personID uint, status string, productID *uint) ([]entities.InternalPaymentRecord, *rest_err.RestErr) {
	query := gorm.G[entities.InternalPaymentRecord](g.db).Where("person_id = ?", personID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if productID != nil {
		query = query.Where("contracted_product_id = ?", *productID)
	}

	records, err := query.Order("reference_month DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching internal payment records").WithCause(err)
	}
	return records, nil
}
```

`gorm_pre_approved_limit_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormPreApprovedLimitRepository struct {
	db *gorm.DB
}

func NewGormPreApprovedLimitRepository(db *gorm.DB) interfaces.PreApprovedLimitRepository {
	return &gormPreApprovedLimitRepository{db: db}
}

func (g *gormPreApprovedLimitRepository) GetByPersonID(ctx context.Context, personID uint, onlyActive bool) ([]entities.PreApprovedLimit, *rest_err.RestErr) {
	query := gorm.G[entities.PreApprovedLimit](g.db).Where("person_id = ?", personID)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	limits, err := query.Order("calculated_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching pre approved limits").WithCause(err)
	}
	return limits, nil
}
```

`gorm_income_declaration_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormIncomeDeclarationRepository struct {
	db *gorm.DB
}

func NewGormIncomeDeclarationRepository(db *gorm.DB) interfaces.IncomeDeclarationRepository {
	return &gormIncomeDeclarationRepository{db: db}
}

func (g *gormIncomeDeclarationRepository) GetByPersonID(ctx context.Context, personID uint, verifiedOnly bool) ([]entities.IncomeDeclaration, *rest_err.RestErr) {
	query := gorm.G[entities.IncomeDeclaration](g.db).Where("person_id = ?", personID)
	if verifiedOnly {
		query = query.Where("verified = ?", true)
	}

	declarations, err := query.Order("declaration_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching income declarations").WithCause(err)
	}
	return declarations, nil
}
```

- [ ] **Step 4: Registrar as cinco no `factory.go`**

Acrescentar os cinco campos ao struct `RepositoryFactory`, inicializá-los em `NewRepositoryFactory` com `NewGorm...Repository(tpf.DB())`, e criar os cinco getters, exatamente no padrão dos campos existentes. Exemplo do bloco a inserir no struct:

```go
	customerRelationshipRepository interfaces.CustomerRelationshipRepository
	contractedProductRepository    interfaces.ContractedProductRepository
	internalPaymentRecordRepository interfaces.InternalPaymentRecordRepository
	preApprovedLimitRepository     interfaces.PreApprovedLimitRepository
	incomeDeclarationRepository    interfaces.IncomeDeclarationRepository
```

No construtor:

```go
		customerRelationshipRepository:  NewGormCustomerRelationshipRepository(tpf.DB()),
		contractedProductRepository:     NewGormContractedProductRepository(tpf.DB()),
		internalPaymentRecordRepository: NewGormInternalPaymentRecordRepository(tpf.DB()),
		preApprovedLimitRepository:      NewGormPreApprovedLimitRepository(tpf.DB()),
		incomeDeclarationRepository:     NewGormIncomeDeclarationRepository(tpf.DB()),
```

Getters:

```go
func (f *RepositoryFactory) CustomerRelationshipRepository() interfaces.CustomerRelationshipRepository {
	return f.customerRelationshipRepository
}
func (f *RepositoryFactory) ContractedProductRepository() interfaces.ContractedProductRepository {
	return f.contractedProductRepository
}
func (f *RepositoryFactory) InternalPaymentRecordRepository() interfaces.InternalPaymentRecordRepository {
	return f.internalPaymentRecordRepository
}
func (f *RepositoryFactory) PreApprovedLimitRepository() interfaces.PreApprovedLimitRepository {
	return f.preApprovedLimitRepository
}
func (f *RepositoryFactory) IncomeDeclarationRepository() interfaces.IncomeDeclarationRepository {
	return f.incomeDeclarationRepository
}
```

- [ ] **Step 5: Escrever o teste de contrato** (compilação garante que as impls satisfazem as interfaces)

`contract_test.go`:

```go
package repositories

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
)

// Asserções de compilação: cada construtor devolve a interface esperada.
var (
	_ interfaces.CustomerRelationshipRepository   = (*gormCustomerRelationshipRepository)(nil)
	_ interfaces.ContractedProductRepository      = (*gormContractedProductRepository)(nil)
	_ interfaces.InternalPaymentRecordRepository  = (*gormInternalPaymentRecordRepository)(nil)
	_ interfaces.PreApprovedLimitRepository       = (*gormPreApprovedLimitRepository)(nil)
	_ interfaces.IncomeDeclarationRepository      = (*gormIncomeDeclarationRepository)(nil)
)

func TestRepositoryFactoryWiresDimensions(t *testing.T) {
	// Guarda contra getters esquecidos: o factory precisa expor os cinco.
	var f RepositoryFactory
	_ = f.CustomerRelationshipRepository
	_ = f.ContractedProductRepository
	_ = f.InternalPaymentRecordRepository
	_ = f.PreApprovedLimitRepository
	_ = f.IncomeDeclarationRepository
}
```

- [ ] **Step 6: Compilar e testar**

Run: `cd internal-registry && go build ./... && go test ./internal/infra/repositories/...`
Expected: build OK, teste PASS.

- [ ] **Step 7: Commit**

```bash
git add internal-registry/internal/infra/repositories internal-registry/internal/infra/controllers/dto/internal_registry_result_dto.go
git commit -m "feat(internal-registry): repositorios e result DTOs das cinco dimensoes"
```

---

### Task 3: `PersonSummaryDTO` + `PersonService.GetAllSummary` (L4, parte 1)

**Files:**
- Modify: `internal-registry/internal/infra/controllers/dto/person_dto.go`
- Modify: `internal-registry/internal/services/person_service.go`
- Test: `internal-registry/internal/infra/controllers/dto/person_summary_dto_test.go`

**Interfaces:**
- Produces (usadas pelas Tasks 4 e 5):
  - `dto.PersonSummaryDTO{ ID uint; Name string; Document string }` e `dto.NewPersonSummaryDTO(*entities.Person) *PersonSummaryDTO`.
  - `(*services.PersonService).GetAllSummary(ctx, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)`.

- [ ] **Step 1: Escrever o teste do `NewPersonSummaryDTO`**

`person_summary_dto_test.go`:

```go
package dto

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonSummaryDTO(t *testing.T) {
	p := &entities.Person{
		Model: gorm.Model{ID: 7},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Felipe Pereira Santos",
			Document: valueobjects.Document("11122233344"),
		},
	}
	got := NewPersonSummaryDTO(p)
	if got.ID != 7 || got.Name != "Felipe Pereira Santos" || got.Document != "11122233344" {
		t.Fatalf("summary inesperado: %+v", got)
	}
}

func TestNewPersonSummaryDTO_NilPersonalInformation(t *testing.T) {
	got := NewPersonSummaryDTO(&entities.Person{Model: gorm.Model{ID: 1}})
	if got.ID != 1 || got.Name != "" || got.Document != "" {
		t.Fatalf("esperado apenas ID preenchido, veio %+v", got)
	}
}
```

- [ ] **Step 2: Rodar o teste e ver falhar**

Run: `cd internal-registry && go test ./internal/infra/controllers/dto/ -run PersonSummary`
Expected: FAIL (`NewPersonSummaryDTO` não existe).

- [ ] **Step 3: Adicionar `PersonSummaryDTO`** ao fim de `person_dto.go` (espelha open-finance)

```go
// PersonSummaryDTO é uma projeção enxuta de um cliente, usada pela listagem.
// Expõe apenas o necessário para identificá-lo, obrigando o chamador a buscar o
// cadastro completo por get_customer_by_id ou get_customer_by_document antes de
// qualquer análise.
type PersonSummaryDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

func NewPersonSummaryDTO(entity *entities.Person) *PersonSummaryDTO {
	dto := &PersonSummaryDTO{ID: entity.ID}
	if entity.PersonalInformation != nil {
		dto.Name = entity.PersonalInformation.FullName
		dto.Document = entity.PersonalInformation.Document.String()
	}
	return dto
}
```

- [ ] **Step 4: Adicionar `GetAllSummary`** a `person_service.go` (espelha open-finance)

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

- [ ] **Step 5: Rodar o teste e ver passar**

Run: `cd internal-registry && go build ./... && go test ./internal/infra/controllers/dto/ -run PersonSummary`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal-registry/internal/infra/controllers/dto/person_dto.go \
        internal-registry/internal/infra/controllers/dto/person_summary_dto_test.go \
        internal-registry/internal/services/person_service.go
git commit -m "feat(internal-registry): PersonSummaryDTO e GetAllSummary"
```

---

### Task 4: `InternalRegistryService` (orquestração + resolução de cliente)

**Files:**
- Create: `internal-registry/internal/services/internal_registry_service.go`
- Test: `internal-registry/internal/services/internal_registry_service_test.go`

**Interfaces:**
- Consumes: os cinco repositórios da Task 2 e os DTOs de resultado; `PersonRepository` (já existe).
- Produces (usadas pela Task 5):
  - `services.CustomerRef{ ID uint; Document string }`
  - `services.NewInternalRegistryService(rf *repositories.RepositoryFactory) *InternalRegistryService`
  - `(*InternalRegistryService).GetCustomerRelationship(ctx, ref) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr)`
  - `(*InternalRegistryService).GetContractedProducts(ctx, ref, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr)`
  - `(*InternalRegistryService).GetInternalPaymentRecords(ctx, ref, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr)`
  - `(*InternalRegistryService).GetPreApprovedLimits(ctx, ref, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr)`
  - `(*InternalRegistryService).GetIncomeDeclarations(ctx, ref, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr)`

- [ ] **Step 1: Escrever os testes com dublês** (espelha `open_finance_service_test.go`)

`internal_registry_service_test.go`:

```go
package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
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
func (f *fakePersonRepository) GetByDocument(_ context.Context, doc string) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byDocument[doc]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}
func (f *fakePersonRepository) GetAll(_ context.Context, _, _ int, _ map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	return nil, 0, nil
}
func (f *fakePersonRepository) Delete(_ context.Context, _ uint) *rest_err.RestErr { return nil }

type fakeContractedProductRepository struct {
	gotType, gotStatus string
	products           []entities.ContractedProduct
}

func (f *fakeContractedProductRepository) GetByPersonID(_ context.Context, _ uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr) {
	f.gotType, f.gotStatus = productType, status
	return f.products, nil
}

func personFixture() *entities.Person {
	return &entities.Person{
		Model:               gorm.Model{ID: 5},
		PersonalInformation: &entities.PersonalInformation{FullName: "Gabriela Ribeiro Barbosa", Document: valueobjects.Document("35509139404")},
	}
}

func TestResolveCustomer_Ambiguous(t *testing.T) {
	s := &InternalRegistryService{personRepository: &fakePersonRepository{}}
	_, err := s.resolveCustomer(context.Background(), CustomerRef{ID: 1, Document: "x"})
	if err == nil {
		t.Fatal("esperado erro para customer_id e document simultâneos")
	}
}

func TestResolveCustomer_Empty(t *testing.T) {
	s := &InternalRegistryService{personRepository: &fakePersonRepository{}}
	_, err := s.resolveCustomer(context.Background(), CustomerRef{})
	if err == nil {
		t.Fatal("esperado erro quando nenhum identificador é informado")
	}
}

func TestGetContractedProducts_ResolvesDocumentAndFilters(t *testing.T) {
	person := personFixture()
	prodRepo := &fakeContractedProductRepository{products: []entities.ContractedProduct{{ProductType: "credit_card", Status: "active"}}}
	s := &InternalRegistryService{
		personRepository:            &fakePersonRepository{byDocument: map[string]*entities.Person{"35509139404": person}},
		contractedProductRepository: prodRepo,
	}

	res, err := s.GetContractedProducts(context.Background(), CustomerRef{Document: "35509139404"}, "credit_card", "active")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.CustomerID != 5 || res.Document != "35509139404" {
		t.Fatalf("identificação errada no resultado: %+v", res)
	}
	if prodRepo.gotType != "credit_card" || prodRepo.gotStatus != "active" {
		t.Fatalf("filtros não repassados: %+v", prodRepo)
	}
	if len(res.Items) != 1 {
		t.Fatalf("esperado 1 item, veio %d", len(res.Items))
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd internal-registry && go test ./internal/services/ -run 'ResolveCustomer|GetContractedProducts'`
Expected: FAIL (`InternalRegistryService` não existe).

- [ ] **Step 3: Implementar o serviço** (espelha `open_finance_service.go`)

```go
package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
)

// CustomerRef identifica o cliente de uma consulta por dimensão. Exatamente um
// dos dois campos deve estar preenchido.
type CustomerRef struct {
	ID       uint
	Document string
}

type InternalRegistryService struct {
	personRepository                interfaces.PersonRepository
	customerRelationshipRepository  interfaces.CustomerRelationshipRepository
	contractedProductRepository     interfaces.ContractedProductRepository
	internalPaymentRecordRepository interfaces.InternalPaymentRecordRepository
	preApprovedLimitRepository      interfaces.PreApprovedLimitRepository
	incomeDeclarationRepository     interfaces.IncomeDeclarationRepository
}

func NewInternalRegistryService(rf *repositories.RepositoryFactory) *InternalRegistryService {
	return &InternalRegistryService{
		personRepository:                rf.PersonRepository(),
		customerRelationshipRepository:  rf.CustomerRelationshipRepository(),
		contractedProductRepository:     rf.ContractedProductRepository(),
		internalPaymentRecordRepository: rf.InternalPaymentRecordRepository(),
		preApprovedLimitRepository:      rf.PreApprovedLimitRepository(),
		incomeDeclarationRepository:     rf.IncomeDeclarationRepository(),
	}
}

func (s *InternalRegistryService) resolveCustomer(ctx context.Context, ref CustomerRef) (*entities.Person, *rest_err.RestErr) {
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

func (s *InternalRegistryService) GetCustomerRelationship(ctx context.Context, ref CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	rel, err := s.customerRelationshipRepository.GetByPersonID(ctx, person.ID)
	if err != nil {
		return nil, err
	}
	result := &dto.CustomerRelationshipResultDTO{CustomerID: person.ID, Document: documentOf(person)}
	if rel != nil {
		result.Relationship = dto.NewCustomerRelationshipDTO(rel)
	}
	return result, nil
}

func (s *InternalRegistryService) GetContractedProducts(ctx context.Context, ref CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	products, err := s.contractedProductRepository.GetByPersonID(ctx, person.ID, productType, status)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ContractedProductDTO, 0, len(products))
	for i := range products {
		items = append(items, *dto.NewContractedProductDTO(&products[i]))
	}
	return &dto.ContractedProductsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetInternalPaymentRecords(ctx context.Context, ref CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	records, err := s.internalPaymentRecordRepository.GetByPersonID(ctx, person.ID, status, productID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.InternalPaymentRecordDTO, 0, len(records))
	for i := range records {
		items = append(items, *dto.NewInternalPaymentRecordDTO(&records[i]))
	}
	return &dto.InternalPaymentRecordsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetPreApprovedLimits(ctx context.Context, ref CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	limits, err := s.preApprovedLimitRepository.GetByPersonID(ctx, person.ID, onlyActive)
	if err != nil {
		return nil, err
	}
	items := make([]dto.PreApprovedLimitDTO, 0, len(limits))
	for i := range limits {
		items = append(items, *dto.NewPreApprovedLimitDTO(&limits[i]))
	}
	return &dto.PreApprovedLimitsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetIncomeDeclarations(ctx context.Context, ref CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	declarations, err := s.incomeDeclarationRepository.GetByPersonID(ctx, person.ID, verifiedOnly)
	if err != nil {
		return nil, err
	}
	items := make([]dto.IncomeDeclarationDTO, 0, len(declarations))
	for i := range declarations {
		items = append(items, *dto.NewIncomeDeclarationDTO(&declarations[i]))
	}
	return &dto.IncomeDeclarationsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}
```

> Nota: se `NewIncomeDeclarationDTO` não existir em `income_declaration_dto.go`, confira o nome do construtor lá e ajuste a chamada; os demais (`NewCustomerRelationshipDTO`, `NewContractedProductDTO`, `NewInternalPaymentRecordDTO`, `NewPreApprovedLimitDTO`) já existem em `internal_registry_dto.go`.

- [ ] **Step 4: Rodar e ver passar**

Run: `cd internal-registry && go build ./... && go test ./internal/services/ -run 'ResolveCustomer|GetContractedProducts'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal-registry/internal/services/internal_registry_service.go \
        internal-registry/internal/services/internal_registry_service_test.go
git commit -m "feat(internal-registry): InternalRegistryService com resolucao de cliente"
```

---

### Task 5: Ferramentas MCP + fiação do servidor (L4, parte 2)

**Files:**
- Modify: `internal-registry/internal/interfaces/mcp/tools/person.go` (renomear tools consolidadas e apontar a listagem para `GetAllSummary`)
- Create: `internal-registry/internal/interfaces/mcp/tools/internal_registry.go` (5 tools por dimensão)
- Modify: `internal-registry/internal/interfaces/mcp/tools/server.go` (interfaces `PersonService`/`InternalRegistryService`, struct, `NewMCPServer`, `registerTools`)
- Modify: `internal-registry/internal/interfaces/mcp/mcp.go` (injetar `NewInternalRegistryService`)
- Test: `internal-registry/internal/interfaces/mcp/tools/internal_registry_test.go`
- Test: `internal-registry/internal/interfaces/mcp/tools/person_test.go`

**Interfaces:**
- Consumes: `InternalRegistryService` (Task 4), `PersonSummaryDTO`/`GetAllSummary` (Task 3), result DTOs (Task 2).
- Produces: as 8 tools registradas — `get_all_customers`, `get_customer_by_id`, `get_customer_by_document`, `get_customer_relationship`, `get_contracted_products`, `get_internal_payment_records`, `get_pre_approved_limits`, `get_income_declarations`.

- [ ] **Step 1: Reescrever `person.go`** para renomear as três tools e apontar a listagem para `GetAllSummary`

Espelhar `open-finance/.../tools/person.go`: `get_person_by_id` → `get_customer_by_id`, `get_person_by_document` → `get_customer_by_document`, `get_all_persons` → `get_all_customers`; a descrição começa com **"Internal Registry:"**; `GetAllPersonsTool` passa a `mcp.WithOutputSchema[dto.PaginatedResponse[dto.PersonSummaryDTO]]()`; `HandleGetAllPersons` chama `s.personService.GetAllSummary(...)`. Manter os nomes dos métodos Go (`GetPersonByIDTool`, `HandleGetPersonByID`, etc.) para não mexer no registro. Texto de descrição para `get_customer_by_id`:

```go
		"get_customer_by_id",
		mcp.WithDescription(
			`
			Internal Registry: get a customer's internal institutional data by their ID.
			Returns the consolidated internal record: relationship (tenure, segment,
			internal score, churn risk), contracted products, internal payment history,
			pre-approved limits and income declarations.
			This is the institution's OWN data; for the credit bureau score and for
			Open Finance banking data, use the respective servers.
			Example usage:
			{
				"id": 123
			}
			`,
		),
```

(Aplicar descrições análogas a `get_customer_by_document` e `get_all_customers`, seguindo o open-finance.)

- [ ] **Step 2: Criar `internal_registry.go`** com as 5 tools por dimensão

Reaproveitar o helper de identificação do open-finance. Cabeçalho e helpers:

```go
package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

const customerRefDescription = `Identify the customer by "customer_id" OR by "document" (CPF). Provide exactly one of them.`

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

func (s *Server) GetCustomerRelationshipTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's relationship with the institution:
			how long they have been a customer (customer_since, relationship_months),
			segment, branch, whether the relationship is active, the internal behavioral
			score and the churn risk.
			A customer with no relationship record is not a client of this institution;
			treat the relationship criteria as missing data rather than as zero.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.CustomerRelationshipResultDTO](),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_customer_relationship", opts...)
}

func (s *Server) HandleGetCustomerRelationship(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.CustomerRelationshipResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetCustomerRelationship(ctx, ref)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetContractedProductsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the products the customer has contracted with the
			institution (checking account, credit card, loan, insurance, investment),
			each with its status, balance and monthly value.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.ContractedProductsResultDTO](),
		mcp.WithString("product_type", mcp.Description("Optional filter by product type"),
			mcp.Enum("checking_account", "credit_card", "loan", "insurance", "investment")),
		mcp.WithString("status", mcp.Description("Optional filter by status"),
			mcp.Enum("active", "closed", "suspended")),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_contracted_products", opts...)
}

func (s *Server) HandleGetContractedProducts(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.ContractedProductsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetContractedProducts(ctx, ref,
		request.GetString("product_type", ""), request.GetString("status", ""))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetInternalPaymentRecordsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's internal payment history for products
			contracted with the institution. Each record reports the reference month,
			due date, payment date, amount due and paid, status (on_time, late, missed,
			partial) and days late.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.InternalPaymentRecordsResultDTO](),
		mcp.WithString("status", mcp.Description("Optional filter by status"),
			mcp.Enum("on_time", "late", "missed", "partial")),
		mcp.WithInteger("product_id", mcp.Description("Optional filter by contracted product id"), mcp.Min(1)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_internal_payment_records", opts...)
}

func (s *Server) HandleGetInternalPaymentRecords(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.InternalPaymentRecordsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	var productID *uint
	if pid := request.GetInt("product_id", 0); pid > 0 {
		v := uint(pid)
		productID = &v
	}
	result, restErr := s.internalRegistryService.GetInternalPaymentRecords(ctx, ref,
		request.GetString("status", ""), productID)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetPreApprovedLimitsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's pre-approved credit limits granted by
			the institution's internal policies. Each limit reports the product type,
			approved amount, interest rate, calculation date, validity and whether it is
			active.
			` + customerRefDescription + `
			By default returns only active limits. Set "only_active" to false for the
			full history.`),
		mcp.WithOutputSchema[dto.PreApprovedLimitsResultDTO](),
		mcp.WithBoolean("only_active", mcp.Description("Return only active limits. Defaults to true"), mcp.DefaultBool(true)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_pre_approved_limits", opts...)
}

func (s *Server) HandleGetPreApprovedLimits(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PreApprovedLimitsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetPreApprovedLimits(ctx, ref, request.GetBool("only_active", true))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetIncomeDeclarationsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the income the customer declared to the institution
			during onboarding/relationship. Each declaration reports the type, monthly
			and yearly amount, source, and whether it was verified.
			` + customerRefDescription + `
			Set "verified_only" to true to return only verified declarations.`),
		mcp.WithOutputSchema[dto.IncomeDeclarationsResultDTO](),
		mcp.WithBoolean("verified_only", mcp.Description("Return only verified declarations. Defaults to false"), mcp.DefaultBool(false)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_income_declarations", opts...)
}

func (s *Server) HandleGetIncomeDeclarations(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.IncomeDeclarationsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetIncomeDeclarations(ctx, ref, request.GetBool("verified_only", false))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}
```

- [ ] **Step 3: Atualizar `server.go`** — interfaces, struct, construtor e registro

Espelhar o `server.go` do open-finance. Declarar as interfaces `PersonService` e `InternalRegistryService` no pacote `tools` (para permitir dublês), acrescentar `internalRegistryService InternalRegistryService` ao struct `Server`, ao `NewMCPServer` e registrar as 5 novas tools em `registerTools`:

```go
type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)
}

type InternalRegistryService interface {
	GetCustomerRelationship(ctx context.Context, ref services.CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr)
	GetContractedProducts(ctx context.Context, ref services.CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr)
	GetInternalPaymentRecords(ctx context.Context, ref services.CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr)
	GetPreApprovedLimits(ctx context.Context, ref services.CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr)
	GetIncomeDeclarations(ctx context.Context, ref services.CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr)
}
```

Trocar o campo `personService *services.PersonService` por `personService PersonService` e acrescentar `internalRegistryService InternalRegistryService`; ajustar a assinatura de `NewMCPServer` para receber `personService PersonService, internalRegistryService InternalRegistryService`. Acrescentar ao `registerTools`:

```go
	mcpSrv.AddTool(s.GetCustomerRelationshipTool(), mcp.NewStructuredToolHandler(s.HandleGetCustomerRelationship))
	mcpSrv.AddTool(s.GetContractedProductsTool(), mcp.NewStructuredToolHandler(s.HandleGetContractedProducts))
	mcpSrv.AddTool(s.GetInternalPaymentRecordsTool(), mcp.NewStructuredToolHandler(s.HandleGetInternalPaymentRecords))
	mcpSrv.AddTool(s.GetPreApprovedLimitsTool(), mcp.NewStructuredToolHandler(s.HandleGetPreApprovedLimits))
	mcpSrv.AddTool(s.GetIncomeDeclarationsTool(), mcp.NewStructuredToolHandler(s.HandleGetIncomeDeclarations))
```

Adicionar os imports necessários (`entities`, `dto`, `rest_err`, `services`) como no open-finance.

- [ ] **Step 4: Atualizar `mcp.go`** — injetar o novo serviço

```go
	s := tools.NewMCPServer(
		services.NewUserService(rf),
		services.NewAddressService(rf),
		services.NewAnalystService(rf),
		services.NewPersonService(rf),
		services.NewInternalRegistryService(rf),
	)
```

- [ ] **Step 5: Escrever os testes de handler** (espelha `open_finance_test.go` e `person_test.go`)

`internal_registry_test.go` — dublê do serviço e casos de tabela:

```go
package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeInternalRegistryService struct {
	gotRef        services.CustomerRef
	gotType       string
	gotStatus     string
	gotOnlyActive bool
	err           *rest_err.RestErr
}

func (f *fakeInternalRegistryService) GetCustomerRelationship(_ context.Context, ref services.CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.CustomerRelationshipResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetContractedProducts(_ context.Context, ref services.CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotType, f.gotStatus = ref, productType, status
	if f.err != nil {
		return nil, f.err
	}
	return &dto.ContractedProductsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetInternalPaymentRecords(_ context.Context, ref services.CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotStatus = ref, status
	if f.err != nil {
		return nil, f.err
	}
	return &dto.InternalPaymentRecordsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetPreApprovedLimits(_ context.Context, ref services.CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotOnlyActive = ref, onlyActive
	if f.err != nil {
		return nil, f.err
	}
	return &dto.PreApprovedLimitsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetIncomeDeclarations(_ context.Context, ref services.CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.IncomeDeclarationsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func reqWith(args map[string]any) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Arguments = args
	return r
}

func TestHandleGetContractedProducts_PassesRefAndFilters(t *testing.T) {
	svc := &fakeInternalRegistryService{}
	s := &Server{internalRegistryService: svc}
	_, err := s.HandleGetContractedProducts(context.Background(),
		reqWith(map[string]any{"customer_id": float64(5), "product_type": "credit_card", "status": "active"}), mcp.CallToolParams{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if svc.gotRef.ID != 5 || svc.gotType != "credit_card" || svc.gotStatus != "active" {
		t.Fatalf("ref/filtros não repassados: %+v", svc)
	}
}

func TestHandleGetContractedProducts_NoIdentifier(t *testing.T) {
	s := &Server{internalRegistryService: &fakeInternalRegistryService{}}
	_, err := s.HandleGetContractedProducts(context.Background(), reqWith(map[string]any{}), mcp.CallToolParams{})
	if err == nil {
		t.Fatal("esperado erro quando nem customer_id nem document são informados")
	}
}

func TestHandleGetPreApprovedLimits_DefaultsOnlyActiveTrue(t *testing.T) {
	svc := &fakeInternalRegistryService{}
	s := &Server{internalRegistryService: svc}
	_, err := s.HandleGetPreApprovedLimits(context.Background(),
		reqWith(map[string]any{"document": "11122233344"}), mcp.CallToolParams{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !svc.gotOnlyActive || svc.gotRef.Document != "11122233344" {
		t.Fatalf("only_active default ou ref errados: %+v", svc)
	}
}
```

Atualizar `person_test.go`: copiar o `fakePersonService` do open-finance (implementa `GetById`, `GetByDocument`, `GetAllSummary`) e os testes de `HandleGetPersonByID`/`ByDocument`/`GetAllPersons`, trocando o prefixo do módulo. Os testes de listagem devem passar por `GetAllSummary`.

- [ ] **Step 6: Rodar e ver passar**

Run: `cd internal-registry && go build ./... && go test ./internal/interfaces/mcp/...`
Expected: build OK, tudo PASS.

- [ ] **Step 7: Commit**

```bash
git add internal-registry/internal/interfaces/mcp
git commit -m "feat(internal-registry): tools por dimensao e renomeacao para customer"
```

---

### Task 6: Teste de mapeamento entidade → DTO (§8)

**Files:**
- Test: `internal-registry/internal/infra/controllers/dto/person_dto_test.go`

**Interfaces:**
- Consumes: `dto.NewPersonDTO` e `ToEntity` (já existem).

- [ ] **Step 1: Escrever o teste de nil-safety e round-trip** (espelha `open-finance/.../person_dto_test.go`)

```go
package dto

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonDTO_NilAssociationsNoPanic(t *testing.T) {
	// Pessoa sem nenhuma dimensão interna (não-cliente): não pode dar panic e as
	// listas ficam vazias/omitidas.
	p := &entities.Person{
		Model:               gorm.Model{ID: 6},
		PersonalInformation: &entities.PersonalInformation{FullName: "Henrique Almeida Ribeiro", Document: valueobjects.Document("99988877766")},
	}
	got := NewPersonDTO(p)
	if got.ID != 6 || got.PersonalInformation == nil {
		t.Fatalf("DTO base inesperado: %+v", got)
	}
	if got.CustomerRelationship != nil {
		t.Errorf("esperado relationship nil para não-cliente")
	}
	if len(got.ContractedProducts) != 0 || len(got.InternalPaymentRecords) != 0 ||
		len(got.PreApprovedLimits) != 0 || len(got.IncomeDeclarations) != 0 {
		t.Errorf("esperado dimensões vazias, veio %+v", got)
	}
}

func TestNewPersonDTO_FullAssociations(t *testing.T) {
	score := 740
	p := &entities.Person{
		Model:                gorm.Model{ID: 2},
		PersonalInformation:  &entities.PersonalInformation{FullName: "Henrique Martins Barbosa", Document: valueobjects.Document("22233344455")},
		CustomerRelationship: &entities.CustomerRelationship{RelationshipMonths: 96, Segment: "retail", InternalScore: &score},
		ContractedProducts:   []entities.ContractedProduct{{ProductType: "credit_card", ProductName: "Cartão Gold", Status: "active"}},
		InternalPaymentRecords: []entities.InternalPaymentRecord{{Status: "on_time", AmountDue: 500, AmountPaid: 500}},
		PreApprovedLimits:    []entities.PreApprovedLimit{{ProductType: "personal_loan", ApprovedAmount: 20000, IsActive: true}},
	}
	got := NewPersonDTO(p)
	if got.CustomerRelationship == nil || got.CustomerRelationship.RelationshipMonths != 96 {
		t.Fatalf("relationship não mapeado: %+v", got.CustomerRelationship)
	}
	if len(got.ContractedProducts) != 1 || got.ContractedProducts[0].ProductType != "credit_card" {
		t.Fatalf("produtos não mapeados: %+v", got.ContractedProducts)
	}
	if len(got.PreApprovedLimits) != 1 || got.PreApprovedLimits[0].ApprovedAmount != 20000 {
		t.Fatalf("limites não mapeados: %+v", got.PreApprovedLimits)
	}
}
```

- [ ] **Step 2: Rodar**

Run: `cd internal-registry && go test ./internal/infra/controllers/dto/`
Expected: PASS (o mapeamento já existe; o teste caracteriza e trava o comportamento). Se algum campo divergir, ajustar o teste ao DTO real — não alterar o DTO.

- [ ] **Step 3: Commit**

```bash
git add internal-registry/internal/infra/controllers/dto/person_dto_test.go
git commit -m "test(internal-registry): mapeamento entidade->DTO de person"
```

---

### Task 7: Fixtures dos 10 clientes (L2, L3)

**Files:**
- Modify: `internal-registry/cmd/fixtures/fixtures/personal_informations.json`, `persons.json`, `addresses.json`, `person_addresses.json`, `person_documents.json` (copiar do birô, alinhar)
- Modify: `internal-registry/cmd/fixtures/fixtures/customer_relationships.json`, `contracted_products.json`, `internal_payment_records.json`, `pre_approved_limits.json`, `income_declarations.json` (popular)

**Interfaces:**
- Consumes: esquema da Task 1; o runner `cmd/fixtures` já tem a ordem de carga correta.

- [ ] **Step 1: Copiar as fixtures de identidade do birô**

Copiar de `bureau/cmd/fixtures/fixtures/` os arquivos `personal_informations.json`, `addresses.json`, `person_addresses.json`, `person_documents.json` **na íntegra** (mesmos `id`, nomes e CPFs dos 10 clientes). Para `persons.json`, copiar do birô mas **remover** `credit_score_id` e `financial_profile_id` e **acrescentar** `customer_relationship_id` (apontando para o `id` correspondente em `customer_relationships.json`; para o cliente 6, `customer_relationship_id: null`). Manter `personal_information_id` e `last_verified_at`.

- [ ] **Step 2: Popular as cinco dimensões conforme a tabela de perfis**

Cada objeto usa snake_case, timestamps RFC3339, `deleted_at: null`. `person_id` refere o `id` do cliente. Valores por cliente (os `id` 1–10 seguem a ordem do birô; nomes conferem com o §6.2 da spec):

| id | Cliente | rel_months / segment / is_active | internal_score / churn | pagamento interno (12m) | produtos | pré-aprovado (ativo) |
|----|---------|----------------------------------|------------------------|-------------------------|----------|----------------------|
| 1 | Felipe Pereira Santos | 14 / retail / true | 280 / high | 12 registros, 1 `missed` com `days_late=120` em produto ativo (**K9**), 4 `late` | 2 (checking_account, loan) | nenhum ativo |
| 2 | Henrique Martins Barbosa | 96 / retail / true | 740 / low | 12 `on_time` | 3 (checking, credit_card, investment) | credit_card R$ 25.000 |
| 3 | Fernanda Costa Barbosa | 30 / retail / true | 430 / medium | 12 reg, 3 `late`, 0 `missed` (~0,75 on_time) | 2 (checking, credit_card) | credit_card R$ 3.000 |
| 4 | Fernanda Rodrigues Pereira | 8 / retail / true | 360 / high | 12 reg, 2 `late` | 1 (checking) | nenhum ativo |
| 5 | Gabriela Ribeiro Barbosa | 60 / private / true | **380 / high** (contradiz o birô 809) | 12 reg, 3 `late` (~0,75) | 3 (checking, credit_card, investment) | credit_card R$ 8.000 |
| 6 | Henrique Almeida Ribeiro | — (sem relacionamento) | — | nenhum registro | nenhum | nenhum |
| 7 | Igor Souza Martins | 132 / private / true | 910 / low | 12 `on_time` | 4 (checking, credit_card, loan, investment) | personal_loan R$ 60.000 |
| 8 | Lucas Martins Souza | 40 / retail / true | 450 / medium | 12 reg, 2 `late` | 2 (checking, credit_card) | credit_card R$ 5.000 |
| 9 | Eduardo Barbosa Almeida | 54 / retail / true | 670 / low | 12 reg, 1 `late` (~0,92) | 3 (checking, credit_card, insurance) | credit_card R$ 15.000 |
| 10 | Eduardo Ribeiro Ribeiro | 72 / retail / true | 780 / low | 12 `on_time` | 3 (checking, credit_card, investment) | personal_loan R$ 30.000 |

Regras de preenchimento:
- `customer_relationships`: um por cliente exceto o 6. `customer_since` = data coerente com `relationship_months` a partir de 2026-09-01. `internal_score`/`churn_risk`/`segment`/`is_active` conforme a tabela; `branch` livre (ex.: `"0001 - Centro"`).
- `contracted_products`: `contracted_date` passada, `status` = `active` salvo indicação; `balance`/`monthly_value` coerentes (ex.: loan com `monthly_value`, checking com `balance`). Guardar os `id` para referenciar nos pagamentos.
- `internal_payment_records`: 12 registros mensais por cliente (com produto), `reference_month` de 2025-09 a 2026-08. `status` conforme a tabela; para `on_time`, `payment_date ≤ due_date` e `days_late=0`; para `late`, `days_late` entre 5 e 40; o `missed` do cliente 1 com `days_late=120`, `payment_date=null`, `amount_paid=0`, em produto com `status=active`.
- `pre_approved_limits`: um limite ativo por cliente elegível, `is_active=true`, `valid_until` futuro (ex.: 2026-12-31), `approved_amount` conforme a tabela, `interest_rate` plausível, `policy_version:"internal-v3"`.
- `income_declarations`: uma por cliente exceto o 6, `monthly_amount` = renda declarada do cliente no birô (`bureau` `financial_profiles.declared_monthly_income`), `income_type` coerente (`salary` na maioria), `verified:true` salvo para 1 ou 2 casos.

- [ ] **Step 3: Recriar o banco, migrar e carregar as fixtures**

```bash
cd /home/bruno/College/mcp-servers
docker compose exec -T postgres psql -U root -c 'DROP DATABASE IF EXISTS "internal-registry";'
docker compose exec -T postgres psql -U root -c 'CREATE DATABASE "internal-registry";'
cd internal-registry && make migration-up && go run ./cmd/fixtures -dir ./cmd/fixtures/fixtures -truncate
```

Expected: o runner carrega todas as tabelas sem erro de FK.

- [ ] **Step 4: Conferir contagens e o caso K9**

```bash
docker compose -f /home/bruno/College/mcp-servers/docker-compose.yaml exec -T postgres psql -U root -d internal-registry -c \
"SELECT (SELECT count(*) FROM persons) persons, (SELECT count(*) FROM customer_relationships) rel, (SELECT count(*) FROM internal_payment_records WHERE status='missed' AND days_late>90) k9;"
```

Expected: `persons=10`, `rel=9` (cliente 6 sem relacionamento), `k9=1`.

- [ ] **Step 5: Commit**

```bash
git add internal-registry/cmd/fixtures/fixtures
git commit -m "feat(internal-registry): fixtures dos 10 clientes alinhados ao biro"
```

---

### Task 8: Política de crédito v1.2

**Files:**
- Create: `docs/politica_credito_agente_v1.2.md`

**Interfaces:**
- Consumes: `docs/politica_credito_agente_v1.1.md` (base). Produz o documento que rege os três servidores.

- [ ] **Step 1: Derivar a v1.2 a partir da v1.1**

Copiar `docs/politica_credito_agente_v1.1.md` e aplicar exatamente as mudanças da §7 da spec:
- Cabeçalho: versão 1.2; citar os três servidores (birô, Open Finance, registro interno) e as tools novas.
- §1: instrução para consultar também o registro interno; ausência de relacionamento interno não impede a análise (regras de dado ausente).
- §4: acrescentar **K9 — Inadimplência interna** após K8, com o texto da §7.3 da spec (missed & days_late>90 em produto ativo → `REPROVADO`; missed & days_late≤90 nos últimos 6 meses → `ANALISE_MANUAL`).
- §5 (scorecard): substituir a tabela pela rebalanceada da §7.2 da spec (C1–C12, total 100), incluindo as faixas ajustadas e a origem de C10–C12. Acrescentar C10, C11, C12 e suas regras de nulo (não-cliente → C10=2, C11=2; internal_score nulo → C11=2; sem pagamentos internos 12m → C12=3).
- §6: a degradação por dados ausentes passa a considerar também C10, C11 e C12.
- §7 (limite): acrescentar o piso pela pré-aprovação vigente (somente `APROVADO`), com o texto da §7.4 da spec.
- §8: consolidar as novas regras de dado ausente na tabela.
- §9: acrescentar `C10_relacionamento`, `C11_score_interno`, `C12_pagamento_interno` ao bloco `criterios`; `politica_versao:"1.2"`; `fontes_consultadas` admite `internal-registry-mcp`.
- §5.x: registrar que `contracted_products` e `income_declarations` são contexto, não pontuados (spec §7.5).

- [ ] **Step 2: Conferir a soma do scorecard**

Revisar manualmente: 25+12+12+10+3+3+8+4+5 (C1–C9) = 82; +6+6+6 (C10–C12) = 100. Bandas da §6 inalteradas.

- [ ] **Step 3: Commit**

```bash
git add docs/politica_credito_agente_v1.2.md
git commit -m "docs: politica de credito v1.2 com criterios de registro interno"
```

---

### Task 9: Infraestrutura (L6)

**Files:**
- Modify: `init.sql`
- Modify: `docker-compose.yaml`
- Modify: `infra/main.tf`, `infra/variables.tf`, `infra/secrets.tf`
- Create: `.github/workflows/deploy-internal-registry.yml`

**Interfaces:**
- Consumes: o serviço `internal-registry` (Dockerfile `target: mcp`).

- [ ] **Step 1: `init.sql`** — acrescentar a linha

```sql
CREATE DATABASE "internal-registry";
```

- [ ] **Step 2: `docker-compose.yaml`** — acrescentar o serviço espelhando `open-finance-mcp`, com porta `8083:8080`, `build.context: ./internal-registry`, `target: mcp`, `env_file: ./internal-registry/.env`, mesmas redes e `depends_on: postgres`.

- [ ] **Step 3: `infra/main.tf`** — acrescentar `"internal-registry-mcp"` ao `toset` `local.mcp_services`.

- [ ] **Step 4: `infra/variables.tf` / `infra/secrets.tf`** — acrescentar a imagem do novo serviço ao mapa `var.images` e o secret `DATABASE_URL` próprio do registro interno, no mesmo padrão do open-finance.

- [ ] **Step 5: Workflow** — copiar `.github/workflows/deploy-open-finance.yml` para `deploy-internal-registry.yml`, trocando nomes/paths (`internal-registry/**` e `infra/**`) e o alvo de build/deploy.

- [ ] **Step 6: Validar (sem aplicar)**

```bash
cd /home/bruno/College/mcp-servers
docker compose config >/dev/null && echo "compose OK"
cd infra && terraform fmt -check && terraform validate
```

Expected: compose válido; `terraform validate` OK (se o backend exigir `init`, rodar `terraform init -backend=false` antes do `validate`).

- [ ] **Step 7: Commit** (sem push, sem deploy)

```bash
cd /home/bruno/College/mcp-servers
git add init.sql docker-compose.yaml infra/main.tf infra/variables.tf infra/secrets.tf .github/workflows/deploy-internal-registry.yml
git commit -m "chore: provisionamento local e de nuvem do servidor de registro interno"
```

---

## Self-Review

**Spec coverage:**
- §3 (esquema) → Task 1. §4 (repos/serviço/summary) → Tasks 2, 3, 4. §5 (8 tools) → Task 5. §6 (fixtures) → Task 7. §7 (política v1.2) → Task 8. §8 (testes) → Tasks 2, 3, 4, 5, 6. §9 (infra) → Task 9. §11 (ordem) → ordem das tasks. Sem lacunas.

**Type consistency:** `CustomerRef`, os cinco métodos `Get*` do serviço, as cinco interfaces de repositório e os result DTOs têm nomes e assinaturas idênticos onde são produzidos (Tasks 2/4) e consumidos (Tasks 4/5). `GetByPersonID` do relacionamento devolve `(*entities.CustomerRelationship, *rest_err.RestErr)` com `(nil,nil)` para não-cliente, tratado no serviço (Task 4) e refletido no DTO `Relationship *CustomerRelationshipDTO` (Task 2).

**Placeholder scan:** as fixtures da Task 7 são especificadas por tabela de valores concretos + regras de derivação determinísticas (copiar do birô / do `declared_monthly_income`); não há "preencher depois".
