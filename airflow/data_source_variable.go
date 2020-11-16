package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const (
	mkDataSourceVariableId    = "id"
	mkDataSourceVariableValue = "value"
)

func dataSourceVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVariableRead,
		Schema: map[string]*schema.Schema{
			mkDataSourcePoolId: {
				Type:     schema.TypeString,
				Required: true,
			},
			mkDataSourceVariableValue: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceVariableRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	variableId := d.Get(mkDataSourceVariableId).(string)
	res, err := c.GetVariableWithResponse(
		ctx,
		api.VariableKey(variableId),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourceVariableValue, res.JSON200.Value)

	d.SetId(variableId)

	return diags
}
