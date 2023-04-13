resource "neon_role" "example" {
  name       = "sally"
  branch_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
