variable "project_id" {
  description = "ID do projeto GCP"
  type        = string
}

variable "region" {
  description = "Região do GCP"
  type        = string
  default     = "southamerica-east1"
}

variable "database_url" {
  description = "Connection string do Supabase. Ex: postgresql://postgres.[ref]:[senha]@aws-0-*.pooler.supabase.com:6543/postgres?sslmode=require&pgbouncer=true"
  type        = string
  sensitive   = true
}