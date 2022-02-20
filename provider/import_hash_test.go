package provider

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlert_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBcryptHashImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBcryptHashExists("bcrypt_hash.test_import"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_import", "cleartext", "Hunter12"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_import", "cost", "10"),
				),
			},
			{
				ResourceName:      "bcrypt_hash.test_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBcryptHashImporter_basic() string {
	return fmt.Sprintf(`
		resource "bcrypt_hash" "test_import" {
			cleartext = "Hunter12"
		}
	`)
}
