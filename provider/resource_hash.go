package provider

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}


func createHash(clear string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clear), cost)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}


func compareHash(hash string, clear string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}


func resourceHash() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cleartext": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The string to hash",
				ForceNew:    true,
				Sensitive: 		true,
			},
			"cost": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The cost parameter for the bcrypt algorithm",
				Default:     10,
			},
			// "hash": {
			// 	Type:					schema.TypeString,
			// 	Computed: 		true,
			// 	Description:  "The hashed value",
			// },
		},
		Create: resourceCreateHash,
		Read:   resourceReadHash,
		Update: resourceUpdateHash,
		Delete: resourceDeleteHash,
		Exists: resourceExistsHash,
		Importer: &schema.ResourceImporter{
			State: resourceImportHash,
		},
	}
}


func resourceCreateHash(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// tfTags := d.Get("tags").(*schema.Set).List()
	// tags := make([]string, len(tfTags))
	// for i, tfTag := range tfTags {
	// 	tags[i] = tfTag.(string)
	// }

	// item := server.Item{
	// 	Name:        d.Get("name").(string),
	// 	Description: d.Get("description").(string),
	// 	Tags:        tags,
	// }

	// err := apiClient.NewItem(&item)

	cost := d.Get("cost").(int)
	hash, err := createHash(d.Get("cleartext").(string), cost)

	if err != nil {
		return err
	}

	d.SetId(hash)
	d.Set("hash", hash)
	d.Set("cost", cost)
	return nil
}


func resourceReadHash(d *schema.ResourceData, m interface{}) error {

	// hash := d.Id()
	// d.SetId(hash)

	return nil
}


func resourceUpdateHash(d *schema.ResourceData, m interface{}) error {

	if compareHash(d.Id(), d.Get("cleartext").(string)) {
		log.Println("Cleartext unchanged")
		return nil
	}

	hash, err := createHash(d.Get("cleartext").(string), d.Get("cost").(int))

	if err != nil {
		return err
	}

	d.SetId(hash)
	d.Set("hash", hash)

	return nil
}


func resourceDeleteHash(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// itemId := d.Id()

	// err := apiClient.DeleteItem(itemId)

	var err error
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}


func resourceExistsHash(d *schema.ResourceData, m interface{}) (bool, error) {
	// apiClient := m.(*client.Client)

	// itemId := d.Id()
	// _, err := apiClient.GetItem(itemId)
	var err error

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}


func resourceImportHash(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	var da []*schema.ResourceData

	hash := d.Id()
	d.Set("cleartext", hash)
	d.Set("cost", 10)

	da = append(da, d)
	return da, nil
}
