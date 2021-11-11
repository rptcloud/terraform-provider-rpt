package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	//"strconv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrderCreate,
		ReadContext:   resourceOrderRead,
		UpdateContext: resourceOrderUpdate,
		DeleteContext: resourceOrderDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"rating": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"photo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	city := d.Get("city").(string)
	state := d.Get("state").(string)
	rating := d.Get("rating").(string)
	photo := d.Get("photo").(string)
	description := d.Get("description").(string)

	values := map[string]string{"name": name, "city": city, "state": state, "rating": rating, "photo": photo, "description": description}
	jsonHotel, err := json.Marshal(values)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hotels/", os.Getenv("endpoint")+":8000/api/v1"), bytes.NewBuffer(jsonHotel))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	response := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%.0f", response["id"]))

	resourceOrderRead(ctx, d, m)

	return diags
}

func resourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	client := &http.Client{Timeout: 10 * time.Second}
	var diags diag.Diagnostics

	hotelIDFloat := d.Get("id").(string)
	hotelID := strings.Split(hotelIDFloat, ".")[0]

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hotels/"+hotelID+"/", os.Getenv("endpoint")+":8000/api/v1"), nil)
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

	if err := d.Set("name", hotel["name"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("city", hotel["city"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", hotel["state"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rating", fmt.Sprintf("%.1f", hotel["rating"])); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("photo", hotel["photo"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", hotel["description"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceOrderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	hotelIDFloat := d.Get("id").(string)
	hotelID := strings.Split(hotelIDFloat, ".")[0]

	needUpdate := false
	if d.HasChange("name") {
		needUpdate = true
	}
	if d.HasChange("state") {
		needUpdate = true
	}
	if d.HasChange("city") {
		needUpdate = true
	}
	if d.HasChange("photo") {
		needUpdate = true
	}
	if d.HasChange("description") {
		needUpdate = true
	}
	if d.HasChange("rating") {
		needUpdate = true
	}

	if needUpdate {
		name := d.Get("name").(string)
		city := d.Get("city").(string)
		state := d.Get("state").(string)
		rating := d.Get("rating").(string)
		photo := d.Get("photo").(string)
		description := d.Get("description").(string)

		values := map[string]string{"name": name, "city": city, "state": state, "rating": rating, "photo": photo, "description": description}
		jsonHotel, err := json.Marshal(values)
		if err != nil {
			return diag.FromErr(err)
		}
		req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/hotels/"+hotelID+"/", os.Getenv("endpoint")+":8000/api/v1"), bytes.NewBuffer(jsonHotel))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			return diag.FromErr(err)
		}

		r, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}
		defer r.Body.Close()

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceOrderRead(ctx, d, m)
}

func resourceOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := &http.Client{Timeout: 10 * time.Second}

	hotelIDFloat := d.Get("id").(string)
	hotelID := strings.Split(hotelIDFloat, ".")[0]

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hotels/"+hotelID+"/", os.Getenv("endpoint")+":8000/api/v1"), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags

}
