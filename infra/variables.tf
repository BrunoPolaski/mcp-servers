variable "project_id" {
  description = "ID do projeto GCP"
  type        = string
}

variable "region" {
  description = "Região do GCP"
  type        = string
  default     = "us-central1"
}

variable "database_url" {
  description = "Connection string do Supabase. Ex: postgresql://postgres.[ref]:[senha]@aws-0-*.pooler.supabase.com:6543/postgres?sslmode=require&pgbouncer=true"
  type        = string
  sensitive   = true
}

variable "image" {
  description = "Imagem Docker completa com tag. Ex: southamerica-east1-docker.pkg.dev/proj/bureau-mcp/bureau-mcp:abc1234"
  type        = string
  default     = "us-docker.pkg.dev/cloudrun/container/hello:latest" # placeholder para o primeiro apply
}

variable "mcp_auth_token" {
  description = "Bearer token para autenticar clientes no MCP server"
  type        = string
  sensitive   = true
}
