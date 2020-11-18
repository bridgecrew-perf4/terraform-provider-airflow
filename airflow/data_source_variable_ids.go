package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourceVariableIds = "airflow_variable_ids"

const (
	mkDataSourceVariableIdsFilter       = "filter"
	mkDataSourceVariableIdsFilterLimit  = "limit"
	mkDataSourceVariableIdsFilterOffset = "offset"
	mkDataSourceVariableIdsIds          = "ids"
)

func dataSourceVariableIds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVariableIdsRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceVariableIdsFilter: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVariableIdsFilterLimit: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
						mkDataSourceVariableIdsFilterOffset: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
					},
				},
			},
			mkDataSourceVariableIdsIds: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVariableIdsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(api.ClientWithResponsesInterface)

	var limit int
	var offset int

	filter := d.Get(mkDataSourceVariableIdsFilter).([]interface{})
	if len(filter) == 1 {
		f := filter[0].(map[string]interface{})

		if v, ok := f[mkDataSourceVariableIdsFilterLimit]; ok {
			limit = v.(int)
		}

		if v, ok := f[mkDataSourceVariableIdsFilterOffset]; ok {
			offset = v.(int)
		}
	}

	var params api.GetVariablesParams
	if limit != 0 {
		tmp := api.PageLimit(limit)
		params.Limit = &tmp
	}

	if offset != 0 {
		tmp := api.PageOffset(offset)
		params.Offset = &tmp
	}

	res, err := c.GetVariablesWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	var variables []string
	if res.JSON200.Variables != nil {
		for _, rVar := range *res.JSON200.Variables {
			variables = append(variables, *rVar.Key)
		}
	}

	_ = d.Set(mkDataSourceVariableIdsIds, variables)

	d.SetId("airflow")

	return diags
}
