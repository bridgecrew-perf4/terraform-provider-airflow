---
page_title: "Provider: Airflow"
subcategory: ""
description: |-
  Terraform provider for interacting with Airflow API.
---

# Airflow Provider

-> Visit the [Call APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorials for an interactive getting started experience.

The Airflow provider is used to interact with a fictional coffee-shop application, Airflow. This provider is meant to serve as an educational tool to show users how:
1. use providers to [create, read, update and delete (CRUD) resources](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers) using Terraform.
1. create a custom Terraform provider.

To learn how to re-create the Airflow provider, refer to the [Call APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorials.

Use the navigation to the left to read about the available resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "airflow" {
  username = "education"
  password = "test123"
}
```

## Schema

### Optional

- **username** (String, Optional) Username to authenticate to Airflow API
- **password** (String, Optional) Password to authenticate to Airflow API
- **host** (String, Optional) Airflow API address (defaults to `localhost:19090`)