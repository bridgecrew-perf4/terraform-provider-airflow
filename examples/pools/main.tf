terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source   = "hashicorp.com/xabinapal/airflow"
    }
  }
}

resource "airflow_pool" "test" {
  name  = "test_pool"
  slots = 1
}

output "r_pool_test" {
  value = airflow_pool.test
}

data "airflow_pool_ids" "all" {
}

output "d_airflow_pool_ids_all" {
  value = data.airflow_pool_ids.all
}

data "airflow_pool_ids" "filtered" {
  filter {
    limit  = 1
    offset = 1
  }
}

output "d_pool_ids_filtered" {
  value = data.airflow_pool_ids.filtered
}

data "airflow_pool" "default_pool" {
  id = "default_pool"
}

output "d_pool_default_pool" {
  value = data.airflow_pool.default_pool
}
