provider "kubernetes" {
  config_path    = "~/.kube/config"
}

resource "google_service_account" "cert-manager-dns" {
  account_id   = "cert-manager-dns"
  display_name = "cert-manager DNS-01 Challenge SA"
}

resource "google_project_iam_member" "cert-manager-dns" {
  project = var.project
  role    = "roles/dns.admin"
  member  = "serviceAccount:${google_service_account.cert-manager-dns.email}"
}


resource "google_service_account_key" "cert-manager-dns-sa" {
  service_account_id = google_service_account.cert-manager-dns.name
}

resource "kubernetes_secret_v1" "cert-manager-dns-sa-key" {
  metadata {
    name = "cert-manager-dns-sa-key"
    namespace = "cert-manager"
  }

  data = {
    "key.json" = base64decode(google_service_account_key.cert-manager-dns-sa.private_key)
  }
}