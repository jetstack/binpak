terraform {
  backend "local" {
    path = "tfstate/terraform.tfstate"
  }
}

provider "google" {
  project     = var.project
  region      = var.region
}