# Guia de Setup — bureau-mcp no Cloud Run

## Pré-requisitos

```bash
brew install google-cloud-sdk terraform
gcloud auth login
gcloud config set project SEU_PROJECT_ID
```

---

## 1. Primeiro deploy (uma única vez)

### 1.1 Crie o bucket de state do Terraform

```bash
gsutil mb -l us-central1 gs://tfstate-bureau-mcp
gsutil versioning set on gs://tfstate-bureau-mcp
```

### 1.2 Rode o Supabase migration

Cole o arquivo `migrations/001_mcp_tokens.sql` no **SQL Editor** do Supabase e execute.

### 1.3 Terraform apply inicial (cria toda a infra)

```bash
cd infra/

terraform init

terraform apply \
  -var="project_id=SEU_PROJECT_ID" \
  -var="database_url=postgresql://postgres.[ref]:[senha]@aws-0-*.pooler.supabase.com:6543/postgres?sslmode=require&pgbouncer=true" \
  -var="mcp_auth_token=$(openssl rand -hex 32)"
```

> Anote o `mcp_auth_token` gerado — você vai precisar para configurar o cliente MCP.

---

## 2. Configure o GitHub Actions (Workload Identity — sem chave JSON)

O Workload Identity é mais seguro que uma chave JSON porque não há segredo para vazar.

### 2.1 Crie o Workload Identity Pool

```bash
PROJECT_ID="SEU_PROJECT_ID"
GITHUB_ORG="seu-usuario-ou-org"
REPO="seu-repo"

# Cria o pool
gcloud iam workload-identity-pools create "github-actions" \
  --location="global" \
  --display-name="GitHub Actions"

# Cria o provider
gcloud iam workload-identity-pools providers create-oidc "github" \
  --location="global" \
  --workload-identity-pool="github-actions" \
  --display-name="GitHub" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository=='${GITHUB_ORG}/${REPO}'"
```

### 2.2 Crie a Service Account para o GitHub Actions

```bash
# SA que o Actions vai usar para fazer push e deploy
gcloud iam service-accounts create "github-actions-sa" \
  --display-name="GitHub Actions Deploy"

SA_EMAIL="github-actions-sa@${PROJECT_ID}.iam.gserviceaccount.com"

# Permissões necessárias
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:${SA_EMAIL}" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:${SA_EMAIL}" \
  --role="roles/run.developer"

# Permite o GitHub impersonar essa SA
WORKLOAD_POOL_ID=$(gcloud iam workload-identity-pools describe github-actions \
  --location=global --format='value(name)')

gcloud iam service-accounts add-iam-policy-binding "${SA_EMAIL}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WORKLOAD_POOL_ID}/attribute.repository/${GITHUB_ORG}/${REPO}"

# Anote este valor para o secret GCP_WORKLOAD_IDENTITY_PROVIDER:
gcloud iam workload-identity-pools providers describe github \
  --workload-identity-pool="github-actions" \
  --location="global" \
  --format='value(name)'
```

### 2.3 Adicione os Secrets no GitHub

No repositório: **Settings → Secrets and variables → Actions**

| Secret | Valor |
|--------|-------|
| `GCP_PROJECT_ID` | ID do seu projeto GCP |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Output do último comando acima |
| `GCP_SERVICE_ACCOUNT` | `github-actions-sa@SEU_PROJECT_ID.iam.gserviceaccount.com` |

---

## 3. Como a autenticação de tokens funciona no Go

No seu servidor MCP, ao invés de consultar Redis, consulte a tabela:

```go
func (s *Server) validateToken(ctx context.Context, rawToken string) (bool, error) {
    hash := sha256Hex(rawToken) // hash SHA-256 do token recebido

    var exists bool
    err := s.db.QueryRowContext(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM mcp_tokens
            WHERE token_hash = $1
              AND is_revoked = FALSE
              AND (expires_at IS NULL OR expires_at > NOW())
        )
    `, hash).Scan(&exists)
    if err != nil {
        return false, err
    }

    if exists {
        // Atualiza last_used_at de forma assíncrona (não bloqueia a request)
        go s.db.ExecContext(ctx,
            `UPDATE mcp_tokens SET last_used_at = NOW() WHERE token_hash = $1`, hash)
    }

    return exists, nil
}

func sha256Hex(s string) string {
    h := sha256.Sum256([]byte(s))
    return hex.EncodeToString(h[:])
}
```

---

## 4. Fluxo após o setup

```
git push origin main
        │
        ▼
  GitHub Actions
  ├── docker build --target mcp
  ├── docker push → Artifact Registry
  └── gcloud run services update (nova imagem, zero downtime)
        │
        ▼
  Cloud Run serve o MCP na URL do output do Terraform
```

---

## 5. Inserir um token de acesso

```sql
-- No Supabase SQL Editor
INSERT INTO mcp_tokens (token_hash, description, expires_at)
VALUES (
    encode(sha256('SEU_TOKEN_AQUI'::bytea), 'hex'),
    'Claude Desktop - local dev',
    NOW() + INTERVAL '1 year'
);
```

Ou via `psql` / script Go de administração.
