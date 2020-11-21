---
page_title: "Airflow: airflow_pool"
subcategory: ""
description: |-
  Provides an Airflow pool resource.
---

# Resource: `airflow_pool`

Provides an Airflow pool resource.

## Example Usage

```terraform
resource "airflow_pool" "test" {
  name  = "my-pool-name"
  slots = 0
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, Forces new resource) The unique identifier of the pool. Must be less than or equal to 256 characters in length.
- `slots` - (Required) The number of available slots. Must be at least `0` or the special value `-1` to represent infinite slots.

## Import

Airflow pool can be imported using the `name` argument, *e.g.*

```
$ terraform import airflow_pool.test my-pool-name
```