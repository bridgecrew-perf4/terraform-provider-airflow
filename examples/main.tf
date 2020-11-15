terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

provider "airflow" {
  endpoint = "http://localhost:8080/api/v1"
  username = "admin"
  password = "admin"
}

module "connections" {
  source = "./connections"
}

output "connections" {
  value = module.connections
}
