# ──────────────────────────────────────────
# Secrets no Secret Manager
# Os valores são passados via variáveis sensíveis
# (nunca commite senhas no .tfvars!)
# ──────────────────────────────────────────

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
