variable "project_id" {
  description = "ID do projeto GCP"
  type        = string
}

variable "region" {
  description = "Região do GCP"
  type        = string
  default     = "us-central1"
}

variable "database_urls" {
  description = "DATABASE_URL por serviço MCP"
  type        = map(string)
  sensitive   = true
}

variable "images" {
  description = "Imagem de contêiner por serviço MCP"
  type        = map(string)
}
