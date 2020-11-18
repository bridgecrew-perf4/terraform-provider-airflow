package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourceConnectionIds = "airflow_connection_ids"

const (
	mkDataSourceConnectionIdsFilter       = "filter"
	mkDataSourceConnectionIdsFilterType   = "type"
	mkDataSourceConnectionIdsFilterLimit  = "limit"
	mkDataSourceConnectionIdsFilterOffset = "offset"
	mkDataSourceConnectionIdsIds          = "ids"
)

func dataSourceConnectionIds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionIdsRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceConnectionIdsFilter: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceConnectionIdsFilterType: {
							Type:     schema.TypeString,
							Optional: true,
						},
						mkDataSourceConnectionIdsFilterLimit: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
						mkDataSourceConnectionIdsFilterOffset: {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: helper.ValidateDiagFunc(
								validation.IntAtLeast(0),
							),
						},
					},
				},
			},
			mkDataSourceConnectionIdsIds: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceConnectionIdsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(api.ClientWithResponsesInterface)

	var connType string
	var limit int
	var offset int

	filter := d.Get(mkDataSourceConnectionIdsFilter).([]interface{})
	if len(filter) == 1 {
		f := filter[0].(map[string]interface{})

		if v, ok := f[mkDataSourceConnectionIdsFilterType]; ok {
			connType = v.(string)
		}

		if v, ok := f[mkDataSourceConnectionIdsFilterLimit]; ok {
			limit = v.(int)
		}

		if v, ok := f[mkDataSourceConnectionIdsFilterOffset]; ok {
			offset = v.(int)
		}
	}

	var params api.GetConnectionsParams
	if connType == "" {
		if limit != 0 {
			tmp := api.PageLimit(limit)
			params.Limit = &tmp
		}

		if offset != 0 {
			tmp := api.PageOffset(offset)
			params.Offset = &tmp
		}
	}

	res, err := c.GetConnectionsWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	var connections []string
	if res.JSON200.Connections != nil {
		if connType == "" {
			for _, rConn := range *res.JSON200.Connections {
				connections = append(connections, *rConn.ConnectionId)
			}
		} else {
			for _, rConn := range *res.JSON200.Connections {
				if rConn.ConnType == nil || *rConn.ConnType != connType {
					continue
				}

				if offset > 0 {
					offset -= 1
					continue
				}

				connections = append(connections, *rConn.ConnectionId)

				if limit > 1 {
					limit -= 1
				} else if limit == 1 {
					break
				}
			}
		}
	}

	_ = d.Set(mkDataSourceConnectionIdsIds, connections)

	d.SetId("airflow")

	return diags
}
