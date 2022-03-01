package provider

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHash_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testHashPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: testHashCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testHashCheckBcryptHashImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testHashCheckBcryptHashExists("bcrypt_hash.test_import"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_import", "cleartext", ""),
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

func testHashCheckBcryptHashImporter_basic() string {
	return fmt.Sprintf(`
		resource "bcrypt_hash" "test_import" {
			cleartext = ""
		}
	`)
}
