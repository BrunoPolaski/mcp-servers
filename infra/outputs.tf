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
