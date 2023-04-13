---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neon_role Resource - terraform-provider-neon"
subcategory: ""
description: |-
  Neon role.
---

# neon_role (Resource)

Neon role.

## Example Usage

```terraform
resource "neon_role" "example" {
  name       = "sally"
  branch_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branch_id` (String) Branch the role belongs to.
- `name` (String) Name of the role.
- `project_id` (String) Project the role belongs to.

### Read-Only

- `id` (String) Identifier of the role.
- `password` (String, Sensitive) Password of the role.

## Import

Import is supported using the following syntax:

```shell
terraform import neon_role.example silent-wood-306223:br-mute-rain-788791:sally
```