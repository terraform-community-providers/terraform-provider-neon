resource "neon_project" "example" {
  name      = "something"
  region_id = "aws-us-west-2"

  allowed_ips {
    ips                     = ["192.168.1.1", "10.0.0.0/24"]
    protected_branches_only = true
  }
}
