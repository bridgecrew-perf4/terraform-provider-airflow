package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const (
	mkDataSourceConnectionId     = "id"
	mkDataSourceConnectionType   = "type"
	mkDataSourceConnectionHost   = "host"
	mkDataSourceConnectionSchema = "schema"
	mkDataSourceConnectionLogin  = "login"
	mkDataSourceConnectionPort   = "port"
)

func dataSourceConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceConnectionId: {
				Type:     schema.TypeString,
				Required: true,
			},
			mkDataSourceConnectionType: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionHost: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionSchema: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionLogin: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceConnectionPort: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

	connectionId := d.Get(mkDataSourceConnectionId).(string)
	res, err := c.GetConnectionWithResponse(
		ctx,
		api.ConnectionID(connectionId),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourceConnectionType, res.JSON200.ConnType)
	_ = d.Set(mkDataSourceConnectionHost, res.JSON200.Host)
	_ = d.Set(mkDataSourceConnectionSchema, res.JSON200.Schema)
	_ = d.Set(mkDataSourceConnectionLogin, res.JSON200.Login)
	_ = d.Set(mkDataSourceConnectionPort, res.JSON200.Port)

	d.SetId(connectionId)

	return diags
}
