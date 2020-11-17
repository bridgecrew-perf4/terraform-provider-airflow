package airflow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
)

var (
	ProviderName    string
	ProviderVersion string
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AIRFLOW_ENDPOINT",
					nil,
				),
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.IsURLWithHTTPorHTTPS,
				),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AIRFLOW_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AIRFLOW_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_connection": resourceConnection(),
			"airflow_variable":   resourceVariable(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"airflow_connection":       dataSourceConnection(),
			"airflow_connection_ids":   dataSourceConnectionIds(),
			"airflow_connection_types": dataSourceConnectionTypes(),
			"airflow_pool":             dataSourcePool(),
			"airflow_pool_ids":         dataSourcePoolIds(),
			"airflow_variable":         dataSourceVariable(),
			"airflow_variable_ids":     dataSourceVariableIds(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(
	ctx context.Context,
	d *schema.ResourceData,
) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoint := d.Get("endpoint").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var requestEditors []api.RequestEditorFn

	if provider, err := helper.NewUserAgentProvider(
		fmt.Sprintf("%s/%s", ProviderName, ProviderVersion),
	); err == nil {
		requestEditors = append(requestEditors, provider.Intercept)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Airflow OpenAPI client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	if (username != "") && (password != "") {
		if provider, err := securityprovider.NewSecurityProviderBasicAuth(
			username,
			password,
		); err == nil {
			requestEditors = append(requestEditors, provider.Intercept)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Airflow OpenAPI client",
				Detail:   err.Error(),
			})
			return nil, diags
		}
	}

	c, err := api.NewClientWithResponses(
		endpoint,
		api.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				for _, fn := range requestEditors {
					if err := fn(ctx, req); err != nil {
						return err
					}
				}

				return nil
			},
		),
	)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Airflow OpenAPI client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
