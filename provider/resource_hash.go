package provider

import (
	"fmt"
	"context"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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


func createHash(ctx context.Context, clear string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clear), cost)

	if err != nil {
		tflog.Error(ctx, err.Error())
		return "", err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}


func compareHash(ctx context.Context, hash string, clear string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))

	if err != nil {
		tflog.Error(ctx, err.Error())
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
				ForceNew:    false,
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
		CreateContext: resourceCreateHash,
		ReadContext:   resourceReadHash,
		UpdateContext: resourceUpdateHash,
		DeleteContext: resourceDeleteHash,
		Exists: resourceExistsHash,
		Importer: &schema.ResourceImporter{
			// StateContext: resourceImportHash,
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}


func resourceCreateHash(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cost := d.Get("cost").(int)
	hash, err := createHash(ctx, d.Get("cleartext").(string), cost)

	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	d.SetId(hash)
	return nil
}


func resourceReadHash(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// hash := d.Id()
	// d.SetId(hash)

	return nil
}


func resourceUpdateHash(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if compareHash(ctx, d.Id(), d.Get("cleartext").(string)) {
		tflog.Info(ctx, "Cleartext unchanged")
		return nil
	}

	hash, err := createHash(ctx, d.Get("cleartext").(string), d.Get("cost").(int))

	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	d.SetId(hash)

	return nil
}


func resourceDeleteHash(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// apiClient := m.(*client.Client)

	// itemId := d.Id()

	// err := apiClient.DeleteItem(itemId)

	var err error
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}


func resourceExistsHash(d *schema.ResourceData, m interface{}) (bool, error) {
	return true, nil
}


func resourceImportHash(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, diag.Diagnostics) {
	var da []*schema.ResourceData

	hash := d.Id()
	d.Set("cleartext", hash)
	d.Set("cost", 10)

	da = append(da, d)
	return da, nil
}
