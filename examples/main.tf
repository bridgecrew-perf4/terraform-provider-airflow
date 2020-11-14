terraform {
  required_providers {
    airflow = {
      versions = ["0.1.0"]
      source = "hashicorp.com/xabinapal/airflow"
    }
  }
}

provider "airflow" {
  username = "dos"
  password = "test123"
}

module "psl" {
  source = "./coffee"

  coffee_name = "Packer Spiced Latte"
}

output "psl" {
  value = module.psl.coffee
}

data "airflow_ingredients" "psl" {
  coffee_id = values(module.psl.coffee)[0].id
}

# output "psl_i" {
#   value = data.airflow_ingredients.psl
# }

resource "airflow_order" "new" {
  items {
    coffee {
      id = 3
    }
    quantity = 2
  }
  items {
    coffee {
      id = 2
    }
    quantity = 2
  }
}

output "new_order" {
  value = airflow_order.new
}


data "airflow_order" "first" {
  id = 1
}

output "first_order" {
  value = data.airflow_order.first
}
