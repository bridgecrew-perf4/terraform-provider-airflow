terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source   = "hashicorp.com/xabinapal/airflow"
    }
  }
}

resource "airflow_variable" "test" {
  name  = "test_variable"
  value = "test_value"
}

output "r_variable_test" {
  value = airflow_variable.test
}

data "airflow_variable_ids" "all" {
}

output "d_airflow_variable_ids_all" {
  value = data.airflow_variable_ids.all
}

data "airflow_variable_ids" "filtered" {
  filter {
    limit  = 1
    offset = 1
  }
}

output "d_variable_ids_filtered" {
  value = data.airflow_variable_ids.filtered
}

data "airflow_variable" "test" {
  name = "test"
}

output "d_variable_test" {
  value = data.airflow_variable.test
}
