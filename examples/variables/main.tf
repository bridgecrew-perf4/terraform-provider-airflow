terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

data "airflow_variable_ids" "all" {
}

output "airflow_variable_ids_all" {
  value = data.airflow_variable_ids.all
}

data "airflow_variable_ids" "filtered" {
  filter {
    limit = 1
    offset = 1
  }
}

output "variable_ids_filtered" {
  value = data.airflow_variable_ids.filtered
}

data "airflow_variable" "test2" {
  id = "test2"
}

output "variable_test2" {
  value = data.airflow_variable.test2
}
