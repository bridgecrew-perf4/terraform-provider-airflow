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

const (
	mkProviderEndpoint = "endpoint"
	mkProviderUsername = "username"
	mkProviderPassword = "password"
)

const (
	envProviderEndpoint = "AIRFLOW_ENDPOINT"
	envProviderUsername = "AIRFLOW_USERNAME"
	envProviderPassword = "AIRFLOW_PASSWORD"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			mkProviderEndpoint: {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc(
					envProviderEndpoint,
					nil,
				),
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.IsURLWithHTTPorHTTPS,
				),
			},
			mkProviderUsername: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(envProviderUsername, nil),
			},
			mkProviderPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(mkProviderPassword, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			schemaResourceVariable:   resourceVariable(),
			schemaResourceConnection: resourceConnection(),
			schemaResourcePool:       resourcePool(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			schemaDataSourceVariable:        dataSourceVariable(),
			schemaDataSourceVariableIds:     dataSourceVariableIds(),
			schemaDataSourceConnection:      dataSourceConnection(),
			schemaDataSourceConnectionIds:   dataSourceConnectionIds(),
			schemaDataSourceConnectionTypes: dataSourceConnectionTypes(),
			schemaDataSourcePool:            dataSourcePool(),
			schemaDataSourcePoolIds:         dataSourcePoolIds(),
			schemaDataSourceDag:             dataSourceDag(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(
	ctx context.Context,
	d *schema.ResourceData,
) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoint := d.Get(mkProviderEndpoint).(string)
	username := d.Get(mkProviderUsername).(string)
	password := d.Get(mkProviderPassword).(string)

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
