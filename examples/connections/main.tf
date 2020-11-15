terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

data "airflow_connection_ids" "all" {
}

output "airflow_connection_ids_all" {
  value = data.airflow_connection_ids.all
}

data "airflow_connection_ids" "filtered" {
  filter {
    type = "mysql"
    limit = 1
    offset = 1
  }
}

output "connection_ids_filtered" {
  value = data.airflow_connection_ids.filtered
}

data "airflow_connection" "airflow_db" {
  id = "airflow_db"
}

output "connection_airflow_db" {
  value = data.airflow_connection.airflow_db
}
