package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider exports the actual provider.
// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RPT_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RPT_PASSWORD", nil),
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("endpoint", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rpt_hotel": resourceOrder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rpt_hotels": dataSourceHotels(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type Config struct {
	User     string
	Password string
	Endpoint string
}

/*
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
		Endpoint: d.Get("endpoint").(string),
	}
	return &config, nil
}
*/

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	//endpoint := d.Get("endpoint").(string)

	http_client := &http.Client{Timeout: 10 * time.Second}

	type Client struct {
		username string
		password string
		token    string
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/hotels/"+d.Get("hotel_num").(string)+"/", os.Getenv("endpoint")+":8000/api/v1"), nil)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		r, err := http_client.Do(req)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		defer r.Body.Close()

		c := Client{username, password, ""}
		return c, diags
	}

	c := Client{"", "", ""}
	return c, diags
}
