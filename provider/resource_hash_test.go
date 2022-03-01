package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)


func TestHash_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testHashPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: testHashCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testHashCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					testHashCheckBcryptHashExists("bcrypt_hash.test_item"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_item", "cleartext", "hello"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_item", "cost", "10"),
				),
			},
		},
	})
}


func TestHash_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testHashPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: testHashCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testHashCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testHashCheckBcryptHashExists("bcrypt_hash.test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cleartext", "hello_pre_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cost", "10"),
				),
			},
			{
				Config: testHashCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testHashCheckBcryptHashExists("bcrypt_hash.test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cleartext", "hello_post_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cost", "11"),
				),
			},
		},
	})
}


func TestHash_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testHashPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: testHashCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testHashCheckItemMultiple(),
				Check: resource.ComposeTestCheckFunc(
					testHashCheckBcryptHashExists("bcrypt_hash.test_item"),
					testHashCheckBcryptHashExists("bcrypt_hash.another_item"),
				),
			},
		},
	})
}

func testHashCheckItemDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bcrypt_hash" {
			continue
		}
	}
	// add any checks for dangling resources created by Terraform
	return nil
}


func testHashCheckBcryptHashExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		// name := rs.Primary.ID

		var err error
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}


func testHashCheckItemBasic() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_item" {
  cleartext   = "hello"
}
`)
}


func testHashCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_update" {
  cleartext = "hello_pre_update"
}
`)
}


func testHashCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_update" {
  cleartext = "hello_post_update"
	cost      = 11
}
`)
}


func testHashCheckItemMultiple() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_item" {
  cleartext   = "test"
}


resource "bcrypt_hash" "another_item" {
	cleartext   = "another_test"
}
`)
}
