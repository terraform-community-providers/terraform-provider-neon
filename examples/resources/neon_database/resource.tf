resource "neon_database" "example" {
  name       = "sally"
  owner_name = neon_role.example.role.name
  branch_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
