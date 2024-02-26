package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEndpointResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEndpointResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "300"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_endpoint.test",
				ImportState:       true,
				ImportStateIdFunc: endpointImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccEndpointResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "300"),
				),
			},
			// Update and Read testing
			{
				Config: testAccEndpointResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "1"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "2"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_endpoint.test",
				ImportState:       true,
				ImportStateIdFunc: endpointImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccEndpointResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEndpointResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "1"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "2"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "3600"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_endpoint.test",
				ImportState:       true,
				ImportStateIdFunc: endpointImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccEndpointResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "1"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "2"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "3600"),
				),
			},
			// Update will null values
			{
				Config: testAccEndpointResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("neon_endpoint.test", "id", idRegex()),
					resource.TestCheckResourceAttr("neon_endpoint.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "project_id", "polished-snowflake-328957"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "type", "read_only"),
					resource.TestMatchResourceAttr("neon_endpoint.test", "host", hostRegex("us-east-2")),
					resource.TestCheckResourceAttr("neon_endpoint.test", "min_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "max_cu", "0.25"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "compute_provisioner", "k8s-neonvm"),
					resource.TestCheckResourceAttr("neon_endpoint.test", "suspend_timeout", "300"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_endpoint.test",
				ImportState:       true,
				ImportStateIdFunc: endpointImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEndpointResourceConfigDefault() string {
	return `
resource "neon_endpoint" "test" {
  branch_id = "br-patient-mode-718259"
  project_id = "polished-snowflake-328957"
}
`
}

func testAccEndpointResourceConfigNonDefault() string {
	return `
resource "neon_endpoint" "test" {
  branch_id = "br-patient-mode-718259"
  project_id = "polished-snowflake-328957"
  min_cu = 1
  max_cu = 2
  suspend_timeout = 3600
}
`
}

func endpointImportIdFunc(state *terraform.State) (string, error) {
	rawState, ok := state.RootModule().Resources["neon_endpoint.test"]

	if !ok {
		return "", fmt.Errorf("Resource Not found")
	}

	return fmt.Sprintf("%s:%s", rawState.Primary.Attributes["project_id"], rawState.Primary.Attributes["id"]), nil
}
