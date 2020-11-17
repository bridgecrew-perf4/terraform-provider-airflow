terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source   = "hashicorp.com/xabinapal/airflow"
    }
  }
}

resource "airflow_connection" "test" {
  name     = "test_connection2"
  type     = "http"
  schema   = "httpss"
  host     = "localhost"
  port     = "8080"
  login    = "admin"
  password = "admin"
  extra = jsonencode({
    test_key = "value"
  })
}

output "r_connection_test" {
  value = airflow_connection.test
}

data "airflow_connection_types" "all" {
}

output "d_airflow_connection_types_all" {
  value = data.airflow_connection_types.all
}

data "airflow_connection_ids" "all" {
}

output "d_airflow_connection_ids_all" {
  value = data.airflow_connection_ids.all
}

data "airflow_connection_ids" "filtered" {
  filter {
    type   = "mysql"
    limit  = 1
    offset = 1
  }
}

output "d_connection_ids_filtered" {
  value = data.airflow_connection_ids.filtered
}

data "airflow_connection" "airflow_db" {
  id = "airflow_db"
}

output "d_connection_airflow_db" {
  value = data.airflow_connection.airflow_db
}
