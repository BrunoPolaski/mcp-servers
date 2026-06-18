# ──────────────────────────────────────────
# Secrets no Secret Manager
# Os valores são passados via variáveis sensíveis
# (nunca commite senhas no .tfvars!)
# ──────────────────────────────────────────

resource "google_secret_manager_secret" "database_url" {
  secret_id = "bureau-mcp-database-url"
  replication {
    auto {}
  }
  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "database_url" {
  secret      = google_secret_manager_secret.database_url.id
  secret_data = var.database_url
}