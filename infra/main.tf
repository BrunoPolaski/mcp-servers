terraform {
  required_version = ">= 1.6"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  backend "gcs" {
    bucket = "tfstate-bureau-mc"
    prefix = "cloud-run"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

locals {
  mcp_services = toset(["bureau-mcp", "open-finance-mcp"])
}

# ──────────────────────────────────────────
# APIs necessárias
# ──────────────────────────────────────────
resource "google_project_service" "services" {
  for_each = toset([
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
  ])
  service            = each.key
  disable_on_destroy = false
}

# ──────────────────────────────────────────
# Artifact Registry
# ──────────────────────────────────────────
resource "google_artifact_registry_repository" "repo" {
  for_each = local.mcp_services

  location      = var.region
  repository_id = each.key
  format        = "DOCKER"

  depends_on = [google_project_service.services]
}

# ──────────────────────────────────────────
# Service Account (princípio do menor privilégio)
# ──────────────────────────────────────────
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

# ──────────────────────────────────────────
# Cloud Run — um serviço por conector MCP
# ──────────────────────────────────────────
resource "google_cloud_run_v2_service" "mcp" {
  for_each = local.mcp_services

  name     = each.key
  location = var.region

  template {
    service_account = google_service_account.mcp_sa[each.key].email

    scaling {
      min_instance_count = 0 # escala a zero quando ocioso (economiza $)
      max_instance_count = 5
    }

    containers {
      # A tag :latest é substituída pelo GitHub Actions a cada deploy
      image = var.images[each.key]

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        # CPU só é alocada durante requisições (reduz custo)
        cpu_idle = false
      }

      # ── Variáveis de ambiente vindas do Secret Manager ──
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
