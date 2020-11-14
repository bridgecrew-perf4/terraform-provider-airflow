---
page_title: "order Resource - terraform-provider-airflow"
subcategory: ""
description: |-
  The order resource allows you to configure a Airflow order.
---

# Resource `airflow_order`

-> Visit the [Perform CRUD operations with Providers](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorial for an interactive getting started experience.

The order resource allows you to configure a Airflow order.

## Example Usage

```terraform
resource "airflow_order" "edu" {
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
```

## Argument Reference

- `items` - (Required) Items in a Airflow order. See [Order item](#order-item) below for details.

### Order item

Each order item contains a `coffee` object and a `quantity`.

- `coffee` - (Required) Represents a Airflow coffee object. See [Coffee](#coffee) below for details.
- `quantity` - (Required) The number of coffee in an order item.

### Coffee

- `id` - (Required) The Airflow coffee ID.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

### Coffee

- `image` - The coffee's image URL path.
- `name` - The coffee name.
- `price` - The coffee price.
- `teaser` - The coffee teaser.