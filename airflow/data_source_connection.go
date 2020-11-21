package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourceConnection = "airflow_connection"

const (
	mkDataSourceConnectionName   = "name"
	mkDataSourceConnectionType   = "type"
	mkDataSourceConnectionSchema = "schema"
	mkDataSourceConnectionHost   = "host"
	mkDataSourceConnectionPort   = "port"
	mkDataSourceConnectionLogin  = "login"
)

func dataSourceConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceConnectionName: {
				Type:     schema.TypeString,
				Required: true,
			},
			mkDataSourceConnectionType: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionSchema: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionHost: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionPort: {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionLogin: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceConnectionRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	name := d.Get(mkDataSourceConnectionName).(string)
	res, err := c.GetConnectionWithResponse(
		ctx,
		api.ConnectionID(name),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourceConnectionName, res.JSON200.ConnectionId)
	_ = d.Set(mkDataSourceConnectionType, res.JSON200.ConnType)
	_ = d.Set(mkDataSourceConnectionHost, res.JSON200.Host)
	_ = d.Set(mkDataSourceConnectionSchema, res.JSON200.Schema)
	_ = d.Set(mkDataSourceConnectionLogin, res.JSON200.Login)
	_ = d.Set(mkDataSourceConnectionPort, res.JSON200.Port)

	d.SetId(*res.JSON200.ConnectionId)

	return diags
}
