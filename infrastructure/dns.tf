resource "google_dns_managed_zone" "binpak-me" {
  name        = "binpak-me"
  dns_name    = "${var.domain_name}."
  description = "Example DNS zone"
  labels = {
    foo = "bar"
  }
}

resource "google_dns_record_set" "frontend" {
  name = "${google_dns_managed_zone.binpak-me.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = google_dns_managed_zone.binpak-me.name
  rrdatas = [var.nginx_ingress_ip]
}

output "nameservers" {
  value = google_dns_managed_zone.binpak-me.name_servers
}