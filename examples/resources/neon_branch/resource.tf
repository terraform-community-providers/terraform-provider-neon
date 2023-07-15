resource "neon_branch" "example" {
  name       = "analytics"
  branch_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
