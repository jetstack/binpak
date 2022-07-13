variable "project" {
  description = "GCP Project"
}
variable "region" {
  description = "GCP Region"
}
variable "domain_name" {
  description = "Binpak domain name"
  type = string
  default = "binpak.me"
}
variable "nginx_ingress_ip" {
  description = "Ingress IP address"
  type = string
}