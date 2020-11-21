package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourceVariable = "airflow_variable"

const (
	mkDataSourceVariableName  = "name"
	mkDataSourceVariableValue = "value"
)

func dataSourceVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVariableRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceVariableName: {
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

	name := d.Get(mkDataSourceVariableName).(string)
	res, err := c.GetVariableWithResponse(
		ctx,
		api.VariableKey(name),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourceVariableName, res.JSON200.Key)
	_ = d.Set(mkDataSourceVariableValue, res.JSON200.Value)

	d.SetId(*res.JSON200.Key)

	return diags
}
