package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const (
	mkDataSourcePoolId          = "id"
	mkDataSourcePoolSlots       = "slots"
	mkDataSourcePoolOpenSlots   = "open_slots"
	mkDataSourcePoolQueuedSlots = "queued_slots"
	mkDataSourcePoolUsedSlots   = "used_slots"
)

func dataSourcePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePoolRead,
		Schema: map[string]*schema.Schema{
			mkDataSourcePoolId: {
				Type:     schema.TypeString,
				Required: true,
			},
			mkDataSourcePoolSlots: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			mkDataSourcePoolOpenSlots: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			mkDataSourcePoolQueuedSlots: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			mkDataSourcePoolUsedSlots: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePoolRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	poolId := d.Get(mkDataSourcePoolId).(string)
	res, err := c.GetPoolWithResponse(
		ctx,
		api.PoolName(poolId),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourcePoolSlots, res.JSON200.Slots)
	_ = d.Set(mkDataSourcePoolOpenSlots, res.JSON200.OpenSlots)
	_ = d.Set(mkDataSourcePoolQueuedSlots, res.JSON200.QueuedSlots)
	_ = d.Set(mkDataSourcePoolUsedSlots, res.JSON200.UsedSlots)

	d.SetId(poolId)

	return diags
}
