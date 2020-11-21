---
page_title: "Airflow: airflow_variable"
subcategory: ""
description: |-
  Provides an Airflow variable resource.
---

# Resource: `airflow_variable`

Provides an Airflow variable resource.

## Example Usage

```terraform
resource "airflow_variable" "test" {
  name  = "my-variable-name"
  value = "my-variable-value"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, Forces new resource) The unique identifier of the variable. The unique identifier of the pool. Must be less than or equal to 250 characters in length.
- `value` - (Optional) The value of the variable.

## Import

Airflow variable can be imported using the `name` argument, *e.g.*

```
$ terraform import airflow_variable.test my-variable-name
```