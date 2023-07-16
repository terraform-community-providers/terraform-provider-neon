resource "neon_endpoint" "example" {
  branch_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
