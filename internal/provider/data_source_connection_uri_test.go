package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectionURIDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionURIDataSourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.neon_connection_uri.test", "id", "polished-snowflake-328957:::budget-app:budget-app"),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-summer-hill-233691\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "pooled_uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-summer-hill-233691-pooler\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
				),
			},
			{
				Config: testAccConnectionURIDataSourceConfigSemiDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.neon_connection_uri.test", "id", "polished-snowflake-328957:br-proud-heart-a5e356v0::budget-app:budget-app"),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-weathered-truth-a5451m4z\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "pooled_uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-weathered-truth-a5451m4z-pooler\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
				),
			},
			{
				Config: testAccConnectionURIDataSourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.neon_connection_uri.test", "id", "polished-snowflake-328957:br-proud-heart-a5e356v0:ep-purple-sun-a58l9gwh:budget-app:budget-app"),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-purple-sun-a58l9gwh\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
					resource.TestMatchResourceAttr(
						"data.neon_connection_uri.test", "pooled_uri",
						regexp.MustCompile(
							`^postgresql://budget-app:[^@]*@ep-purple-sun-a58l9gwh-pooler\.us-east-2\.aws\.neon\.tech/budget-app\?sslmode=require$`,
						),
					),
				),
			},
		},
	})
}

func testAccConnectionURIDataSourceConfigDefault() string {
	return `
data "neon_connection_uri" "test" {
  project_id    = "polished-snowflake-328957"
  database_name = "budget-app"
  role_name     = "budget-app"
}
`
}

func testAccConnectionURIDataSourceConfigSemiDefault() string {
	return `
data "neon_connection_uri" "test" {
  project_id    = "polished-snowflake-328957"
  branch_id     = "br-proud-heart-a5e356v0"
  database_name = "budget-app"
  role_name     = "budget-app"
}
`
}

func testAccConnectionURIDataSourceConfigNonDefault() string {
	return `
data "neon_connection_uri" "test" {
  project_id    = "polished-snowflake-328957"
  branch_id     = "br-proud-heart-a5e356v0"
  endpoint_id   = "ep-purple-sun-a58l9gwh"
  database_name = "budget-app"
  role_name     = "budget-app"
}
`
}
