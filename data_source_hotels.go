package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHotels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHotelsRead,
		Schema: map[string]*schema.Schema{
			"hotel_num": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"rating": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"photo": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHotelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	fmt.Println("test!")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hotels/"+d.Get("hotel_num").(string)+"/", os.Getenv("endpoint")+":8000/api/v1"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	hotel := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&hotel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", hotel["id"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", hotel["name"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("city", hotel["city"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", hotel["state"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rating", hotel["rating"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("photo", hotel["photo"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", hotel["description"]); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(d.Get("hotel_num").(string))

	return diags
}
