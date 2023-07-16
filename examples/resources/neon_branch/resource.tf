resource "neon_branch" "example" {
  name       = "analytics"
  parent_id  = neon_project.example.branch.id
  project_id = neon_project.example.id
}
