---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neon_project Resource - terraform-provider-neon"
subcategory: ""
description: |-
  Neon project.
---

# neon_project (Resource)

Neon project.

## Example Usage

```terraform
resource "neon_project" "example" {
  name   = "something"
  region = "aws-us-west-2"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the project.
- `region_id` (String) Region of the project.

### Optional

- `branch` (Attributes) Primary branch settings of the project. (see [below for nested schema](#nestedatt--branch))
- `pg_version` (Number) PostgreSQL version of the project.

### Read-Only

- `id` (String) Identifier of the project.
- `platform_id` (String) Platform of the project.

<a id="nestedatt--branch"></a>
### Nested Schema for `branch`

Optional:

- `endpoint` (Attributes) Comput endpoint settings of the branch. (see [below for nested schema](#nestedatt--branch--endpoint))
- `name` (String) Name of the branch.

Read-Only:

- `id` (String) Identifier of the branch.

<a id="nestedatt--branch--endpoint"></a>
### Nested Schema for `branch.endpoint`

Optional:

- `max_cu` (Number) Maximum number of compute units for the endpoint. **Default** `0.25`.
- `min_cu` (Number) Minimum number of compute units for the endpoint. **Default** `0.25`.

Read-Only:

- `host` (String) Host of the endpoint.
- `id` (String) Identifier of the endpoint.

## Import

Import is supported using the following syntax:

```shell
terraform import neon_project.example silent-wood-306223
```