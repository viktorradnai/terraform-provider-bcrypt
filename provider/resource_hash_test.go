package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"regexp"
	"testing"
)


func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBcryptHashExists("bcrypt_hash.test_item"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_item", "cleartext", "test"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_item", "description", "hello"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_item", "tags.#", "2"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_item", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_item", "tags.1477001604", "tag2"),
				),
			},
		},
	})
}

func TestAccItem_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBcryptHashExists("bcrypt_hash.test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cleartext", "test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "description", "hello"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "tags.#", "2"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_update", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_update", "tags.1477001604", "tag2"),
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBcryptHashExists("bcrypt_hash.test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "cleartext", "test_update"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "description", "updated description"),
					resource.TestCheckResourceAttr(
						"bcrypt_hash.test_update", "tags.#", "2"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_update", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("bcrypt_hash.test_update", "tags.1477001604", "tag2"),
				),
			},
		},
	})
}

func TestAccItem_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemMultiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBcryptHashExists("bcrypt_hash.test_item"),
					testAccCheckBcryptHashExists("bcrypt_hash.another_item"),
				),
			},
		},
	})
}

var whiteSpaceRegex = regexp.MustCompile("name cannot contain whitespace")

func TestAccItem_WhitespaceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckItemWhitespace(),
				ExpectError: whiteSpaceRegex,
			},
		},
	})
}

func testAccCheckItemDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bcrypt_hash" {
			continue
		}

		var err error

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckBcryptHashExists(resource string) resource.TestCheckFunc {
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

func testAccCheckItemBasic() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_item" {
  cleartext   = "hello"
}
`)
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_update" {
  cleartext = "hello_pre_update"
}
`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_update" {
  cleartext = "hello_post_update"
}
`)
}

func testAccCheckItemMultiple() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_item" {
  cleartext   = "test"
}

resource "bcrypt_hash" "another_item" {
	cleartext   = "another_test"
}
`)
}

func testAccCheckItemWhitespace() string {
	return fmt.Sprintf(`
resource "bcrypt_hash" "test_item" {
	cleartext   = "test with whitespace"
}
`)
}
