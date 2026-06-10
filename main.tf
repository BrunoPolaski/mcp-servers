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
  location      = var.region
  repository_id = "bureau-mcp"
  format        = "DOCKER"

  depends_on = [google_project_service.services]
}

# ──────────────────────────────────────────
# Service Account (princípio do menor privilégio)
# ──────────────────────────────────────────
resource "google_service_account" "mcp_sa" {
  account_id   = "bureau-mcp-sa"
  display_name = "Bureau MCP Server"
}

resource "google_project_iam_member" "secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.mcp_sa.email}"
}

# ──────────────────────────────────────────
# Cloud Run — bureau-mcp
# ──────────────────────────────────────────
resource "google_cloud_run_v2_service" "bureau_mcp" {
  name     = "bureau-mcp"
  location = var.region

  template {
    service_account = google_service_account.mcp_sa.email

    scaling {
      min_instance_count = 0   # escala a zero quando ocioso (economiza $)
      max_instance_count = 5
    }

    containers {
      # A tag :latest é substituída pelo GitHub Actions a cada deploy
      image = "${var.region}-docker.pkg.dev/${var.project_id}/bureau-mcp/bureau-mcp:latest"

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
            secret  = google_secret_manager_secret.database_url.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "MCP_AUTH_TOKEN"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.mcp_auth_token.secret_id
            version = "latest"
          }
        }
      }

      # Adicione outras env vars do seu .env aqui seguindo o mesmo padrão
    }
  }

  depends_on = [
    google_project_service.services,
    google_project_iam_member.secret_accessor,
  ]
}

# Acesso público — autenticação feita pelo MCP_AUTH_TOKEN na aplicação
resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.bureau_mcp.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
