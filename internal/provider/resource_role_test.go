package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func existRegex() *regexp.Regexp {
	return regexp.MustCompile("^.+$")
}

func TestAccRoleResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRoleResourceConfigDefault("sally"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neon_role.test", "name", "sally"),
					resource.TestMatchResourceAttr("neon_role.test", "password", existRegex()),
					resource.TestCheckResourceAttr("neon_role.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_role.test", "project_id", "polished-snowflake-328957"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_role.test",
				ImportState:       true,
				ImportStateId:     "polished-snowflake-328957:br-patient-mode-718259:sally",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccRoleResourceConfigDefault("sally"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neon_role.test", "name", "sally"),
					resource.TestMatchResourceAttr("neon_role.test", "password", existRegex()),
					resource.TestCheckResourceAttr("neon_role.test", "branch_id", "br-patient-mode-718259"),
					resource.TestCheckResourceAttr("neon_role.test", "project_id", "polished-snowflake-328957"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "neon_role.test",
				ImportState:       true,
				ImportStateId:     "polished-snowflake-328957:br-patient-mode-718259:sally",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRoleResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "neon_role" "test" {
  name = "%s"
  branch_id = "br-patient-mode-718259"
  project_id = "polished-snowflake-328957"
}
`, name)
}
