package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func hostRegex(region string) *regexp.Regexp {
	return regexp.MustCompile("^[-0-9a-z\\.]+\\." + region + "\\.aws\\.neon\\.tech$")
}

func TestAccProjectResourceDefaultForUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigDefaultForUser("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-blue-haze-97971912"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccProjectResourceConfigDefaultForUser("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-blue-haze-97971912"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// Update just name
			{
				Config: testAccProjectResourceConfigDefaultForUser("nu-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nu-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-blue-haze-97971912"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProjectResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigDefault("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccProjectResourceConfigDefault("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// Update just name
			{
				Config: testAccProjectResourceConfigDefault("nu-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nu-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "false"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigDefaultUpdate("nue-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nue-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "604800"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.0", "112.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.1", "122.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "master"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with null values so that it can be deleted
			{
				Config: testAccProjectResourceConfigDefault("nue-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nue-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "15"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProjectResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigNonDefault("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "14"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "604800"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.0", "112.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.1", "122.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "master"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccProjectResourceConfigNonDefault("todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "14"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "604800"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.0", "112.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.1", "122.0.0.1"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "master"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "true"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "3600"),
				),
			},
			// Update with null values
			{
				Config: testAccProjectResourceConfigNonDefaultUpdate("nue-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nue-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "14"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with explicit logical replication false
			{
				Config: testAccProjectResourceConfigNonDefaultLogicalReplication("nue-todo-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_project.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "name", "nue-todo-app"),
					resource.TestCheckResourceAttr("neon_project.test", "platform_id", "aws"),
					resource.TestCheckResourceAttr("neon_project.test", "region_id", "aws-us-east-2"),
					resource.TestCheckResourceAttr("neon_project.test", "org_id", "org-aged-sky-67916740"),
					resource.TestCheckResourceAttr("neon_project.test", "pg_version", "14"),
					resource.TestCheckResourceAttr("neon_project.test", "history_retention", "86400"),
					resource.TestCheckResourceAttr("neon_project.test", "logical_replication", "true"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.ips.#", "0"),
					resource.TestCheckResourceAttr("neon_project.test", "allowed_ips.protected_branches_only", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.id", idRegex()),
					resource.TestCheckResourceAttr("neon_project.test", "branch.name", "main"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.protected", "false"),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_project.test", "branch.endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_project.test", "branch.endpoint.suspend_timeout", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfigDefaultForUser(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name      = "%s"
  region_id = "aws-us-east-2"
}
`, name)
}

func testAccProjectResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name      = "%s"
  region_id = "aws-us-east-2"
  org_id    = "org-aged-sky-67916740"
}
`, name)
}

func testAccProjectResourceConfigDefaultUpdate(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name               = "%s"
  region_id          = "aws-us-east-2"
  history_retention  = 604800
  org_id             = "org-aged-sky-67916740"

  logical_replication = true

  allowed_ips = {
    ips = ["112.0.0.1", "122.0.0.1"]
    protected_branches_only = true
  }

  branch = {
    name      = "master"
    protected = true

    endpoint = {
      min_cu         = 1
      max_cu         = 2
      suspend_timeout = 3600
    }
  }
}
`, name)
}

func testAccProjectResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name               = "%s"
  region_id          = "aws-us-east-2"
  history_retention  = 604800
  org_id             = "org-aged-sky-67916740"

  pg_version = 14

  logical_replication = true

  allowed_ips = {
    ips = ["112.0.0.1", "122.0.0.1"]
    protected_branches_only = true
  }

  branch = {
    name      = "master"
    protected = true

    endpoint = {
      min_cu         = 1
      max_cu         = 2
      suspend_timeout = 3600
    }
  }
}
`, name)
}

func testAccProjectResourceConfigNonDefaultUpdate(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name = "%s"
  region_id = "aws-us-east-2"
  org_id = "org-aged-sky-67916740"

  pg_version = 14
}
`, name)
}

func testAccProjectResourceConfigNonDefaultLogicalReplication(name string) string {
	return fmt.Sprintf(`
resource "neon_project" "test" {
  name = "%s"
  region_id = "aws-us-east-2"
  org_id = "org-aged-sky-67916740"

  pg_version = 14

  logical_replication = false
}
`, name)
}
