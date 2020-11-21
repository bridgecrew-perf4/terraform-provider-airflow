terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source   = "hashicorp.com/xabinapal/airflow"
    }
  }
}

data "airflow_dag" "tutorial" {
  name = "tutorial"
}

output "d_dag_tutorial" {
  value = data.airflow_dag.tutorial
}
