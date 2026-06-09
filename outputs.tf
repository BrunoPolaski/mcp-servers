output "mcp_url" {
  description = "URL pública do bureau-mcp no Cloud Run"
  value       = google_cloud_run_v2_service.bureau_mcp.uri
}

output "artifact_registry_repo" {
  description = "Endereço do repositório no Artifact Registry"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/bureau-mcp/bureau-mcp"
}
