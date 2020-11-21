---
page_title: "Airflow: airflow_connection"
subcategory: ""
description: |-
  Provides an Airflow connection resource.
---

# Resource: `airflow_connection`

Provides an Airflow connection resource.

~> *NOTE:* The Airflow API does not currently (as of v2.0.0b2) expose the `password` and `extra` arguments of any existing connection. Because of this, they are never updated while read to avoid permanent diffs between configuration and state files.

## Example Usage

```terraform
resource "airflow_connection" "test" {
  name     = "my-connection-name"
  type     = "http"
  schema   = "https"
  host     = "localhost"
  port     = 80
  login    = "username"
  password = "password"
  extra = jsonencode({
    ssl = true
  })
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, Forces new resource) The unique identifier of the connection. Must be less than or equal to 250 characters in length.
- `type` - (Optional) The connection type. Must be one of the values defined in the [`airflow_connection_types` datasource](/docs/data-sources/connection_types.md).
- `schema` - (Optional) The schema. Must be less than or equal to 500 characters in length.
- `host` - (Optional) The host. Must be less than or equal to 500 characters in length.
- `port` - (Optional) The port number.
- `login` - (Optional) The login. Must be less than or equal to 500 characters in length.
- `password` - (Optional, Sensitive) The password. Must be less than or equal to 5000 characters in length.
- `extra` - (Optional) Extra metadata. Non-standard data such as private/SSH keys can be saved here. JSON encoded object. Must be less than or equal to 5000 characters in length.

## Import

Airflow connection can be imported using the `name` argument, *e.g.*

```
$ terraform import airflow_connection.test my-connection-name
```