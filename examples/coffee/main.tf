terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

variable "coffee_name" {
  type    = string
  default = "Vagrante espresso"
}

data "airflow_coffees" "all" {}

# Returns all coffees
output "all_coffees" {
  value = data.airflow_coffees.all.coffees
}

# Only returns packer spiced latte
output "coffee" {
  value = {
    for coffee in data.airflow_coffees.all.coffees :
    coffee.id => coffee
    if coffee.name == var.coffee_name
  }
}
