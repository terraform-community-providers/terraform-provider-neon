---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neon_connection_uri Data Source - terraform-provider-neon"
subcategory: ""
description: |-
  Retrieves normal and pooled connection URIs for a Neon database and role.
---

# neon_connection_uri (Data Source)

Retrieves normal and pooled connection URIs for a Neon database and role.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database_name` (String) Name of the database.
- `project_id` (String) Project to retrieve the connection URI for.
- `role_name` (String) Name of the role.

### Optional

- `branch_id` (String) Branch to retrieve the connection URI for. Defaults to the project's default branch.
- `endpoint_id` (String) Endpoint to retrieve the connection URI for. Defaults to the banch's read-write endpoint.

### Read-Only

- `id` (String) Unique identifier for the connection URI data source.
- `pooled_uri` (String) Pooled connection URI.
- `uri` (String) Normal connection URI.


