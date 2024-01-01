package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "torque_user" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_user.test", "account_role", "Admin"),
					resource.TestCheckResourceAttr("data.torque_user.test", "display_first_name", "First"),
					resource.TestCheckResourceAttr("data.torque_user.test", "display_last_name", "Last"),
					resource.TestCheckResourceAttr("data.torque_user.test", "first_name", "first"),
					resource.TestCheckResourceAttr("data.torque_user.test", "has_access_to_all_spaces", "1"),
					resource.TestCheckResourceAttr("data.torque_user.test", "join_date", "2023-01-24T15:10:24.877995"),
					resource.TestCheckResourceAttr("data.torque_user.test", "last_name", "last"),
					resource.TestCheckResourceAttr("data.torque_user.test", "timezone", "America/Chicago"),
					resource.TestCheckResourceAttr("data.torque_user.test", "user_email", "my@email.com"),
					resource.TestCheckResourceAttr("data.torque_user.test", "user_type", "REGULAR"),
				),
			},
		},
	})
}
