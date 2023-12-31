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
				Config: providerConfig + `data "torque_users" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_users.test", "account_role", "Admin"),
					resource.TestCheckResourceAttr("data.torque_users.test", "display_first_name", "First"),
					resource.TestCheckResourceAttr("data.torque_users.test", "display_last_name", "Last"),
					resource.TestCheckResourceAttr("data.torque_users.test", "first_name", "first"),
					resource.TestCheckResourceAttr("data.torque_users.test", "has_access_to_all_spaces", "1"),
					resource.TestCheckResourceAttr("data.torque_users.test", "join_date", "2023-01-24T15:10:24.877995"),
					resource.TestCheckResourceAttr("data.torque_users.test", "last_name", "last"),
					resource.TestCheckResourceAttr("data.torque_users.test", "timezone", "America/Chicago"),
					resource.TestCheckResourceAttr("data.torque_users.test", "user_email", "my@email.com"),
					resource.TestCheckResourceAttr("data.torque_users.test", "user_type", "REGULAR"),
				),
			},
		},
	})
}
