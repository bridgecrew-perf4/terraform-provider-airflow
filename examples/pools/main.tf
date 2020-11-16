terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

data "airflow_pool_ids" "all" {
}

output "airflow_pool_ids_all" {
  value = data.airflow_pool_ids.all
}

data "airflow_pool_ids" "filtered" {
  filter {
    limit = 1
    offset = 1
  }
}

output "pool_ids_filtered" {
  value = data.airflow_pool_ids.filtered
}

data "airflow_pool" "default_pool" {
  id = "default_pool"
}

output "pool_default_pool" {
  value = data.airflow_pool.default_pool
}
