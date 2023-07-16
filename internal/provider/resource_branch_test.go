package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBranchResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchResourceConfigDefault("analytics"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckNoResourceAttr("neon_branch.test", "endpoint"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccBranchResourceConfigDefault("analytics"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckNoResourceAttr("neon_branch.test", "endpoint"),
				),
			},
			// Update and Read testing
			{
				Config: testAccBranchResourceConfigNonDefault("nue-analytics", "br-patient-mode-718259"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "nue-analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with default endpoint values
			{
				Config: testAccBranchResourceConfigEndpointDefault("nue-analytics", "br-patient-mode-718259"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "nue-analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.compute_provisioner", "k8s-pod"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.suspend_timeout", "300"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBranchResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchResourceConfigNonDefault("analytics", "br-aged-band-448480"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-aged-band-448480"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccBranchResourceConfigNonDefault("analytics", "br-aged-band-448480"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-aged-band-448480"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.min_cu", "1"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.max_cu", "2"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with default endpoint values
			{
				Config: testAccBranchResourceConfigEndpointDefault("nue-analytics", "br-aged-band-448480"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "nue-analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-aged-band-448480"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.id", idRegex()),
					resource.TestMatchResourceAttr("neon_branch.test", "endpoint.host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.compute_provisioner", "k8s-pod"),
					resource.TestCheckResourceAttr("neon_branch.test", "endpoint.suspend_timeout", "300"),
				),
			},
			// Update with null values
			{
				Config: testAccBranchResourceConfigDefaultParent("analytics"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_branch.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_branch.test", "name", "analytics"),
					resource.TestCheckResourceAttr("neon_branch.test", "parent_id", "br-aged-band-448480"),
					resource.TestCheckResourceAttr("neon_branch.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckNoResourceAttr("neon_branch.test", "endpoint"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_branch.test",
				ImportState:       true,
				ImportStateIdFunc: branchImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBranchResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "neon_branch" "test" {
  name = "%s"
  project_id = "polished-snowflake-328957"
}
`, name)
}

func testAccBranchResourceConfigDefaultParent(name string) string {
	return fmt.Sprintf(`
resource "neon_branch" "test" {
  name = "%s"
  parent_id = "br-aged-band-448480"
  project_id = "polished-snowflake-328957"
}
`, name)
}

func testAccBranchResourceConfigEndpointDefault(name string, parentId string) string {
	return fmt.Sprintf(`
resource "neon_branch" "test" {
  name = "%s"
  parent_id = "%s"
  project_id = "polished-snowflake-328957"

  endpoint = {}
}
`, name, parentId)
}

func testAccBranchResourceConfigNonDefault(name string, parentId string) string {
	return fmt.Sprintf(`
resource "neon_branch" "test" {
  name = "%s"
  parent_id = "%s"
  project_id = "polished-snowflake-328957"

  endpoint = {
	min_cu = 1
	max_cu = 2
	suspend_timeout = 3600
  }
}
`, name, parentId)
}

func branchImportIdFunc(state *terraform.State) (string, error) {
	rawState, ok := state.RootModule().Resources["neon_branch.test"]

	if !ok {
		return "", fmt.Errorf("Resource Not found")
	}

	return fmt.Sprintf("%s:%s", rawState.Primary.Attributes["project_id"], rawState.Primary.Attributes["id"]), nil
}
