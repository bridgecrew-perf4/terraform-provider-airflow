package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourcePoolIds = "airflow_pool_ids"

const (
	mkDataSourcePoolIdsFilter       = "filter"
	mkDataSourcePoolIdsFilterLimit  = "limit"
	mkDataSourcePoolIdsFilterOffset = "offset"
	mkDataSourcePoolIdsIds          = "ids"
)

func dataSourcePoolIds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePoolIdsRead,
		Schema: map[string]*schema.Schema{
			mkDataSourcePoolIdsFilter: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourcePoolIdsFilterLimit: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
						mkDataSourcePoolIdsFilterOffset: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
					},
				},
			},
			mkDataSourcePoolIdsIds: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePoolIdsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(api.ClientWithResponsesInterface)

	var limit int
	var offset int

	filter := d.Get(mkDataSourcePoolIdsFilter).([]interface{})
	if len(filter) == 1 {
		f := filter[0].(map[string]interface{})

		if v, ok := f[mkDataSourcePoolIdsFilterLimit]; ok {
			limit = v.(int)
		}

		if v, ok := f[mkDataSourcePoolIdsFilterOffset]; ok {
			offset = v.(int)
		}
	}

	var params api.GetPoolsParams
	if limit != 0 {
		tmp := api.PageLimit(limit)
		params.Limit = &tmp
	}

	if offset != 0 {
		tmp := api.PageOffset(offset)
		params.Offset = &tmp
	}

	res, err := c.GetPoolsWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	var pools []string
	if res.JSON200.Pools != nil {
		for _, rPool := range *res.JSON200.Pools {
			pools = append(pools, *rPool.Name)
		}
	}

	_ = d.Set(mkDataSourcePoolIdsIds, pools)

	d.SetId("airflow")

	return diags
}
