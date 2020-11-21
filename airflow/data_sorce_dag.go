package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaDataSourceDag = "airflow_dag"

const (
	mkDataSourceDagName        = "name"
	mkDataSourceDagDescription = "description"
	mkDataSourceDagPaused      = "paused"
	mkDataSourceDagSubDag      = "subdag"
	mkDataSourceDagRootDagId   = "root_dag"
	mkDataSourceDagStartDate   = "start_date"
	mkDataSourceDagTimezone    = "timezone"
	mkDataSourceDagCatchup     = "catchup"
	mkDataSourceDagConcurrency = "concurrency"
	mkDataSourceDagDefaultView = "default_view"
	mkDataSourceDagOrientation = "orientation"
)

func dataSourceDag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDagRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceDagName: {
				Type:     schema.TypeString,
				Required: true,
			},
			mkDataSourceDagDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			mkDataSourceDagPaused: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			mkDataSourceDagSubDag: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			mkDataSourceDagRootDagId: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			mkDataSourceDagStartDate: {
				Type:     schema.TypeString,
				Computed: true,
			},
			mkDataSourceDagTimezone: {
				Type:     schema.TypeString,
				Computed: true,
			},
			mkDataSourceDagCatchup: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			mkDataSourceDagConcurrency: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			mkDataSourceDagDefaultView: {
				Type:     schema.TypeString,
				Computed: true,
			},
			mkDataSourceDagOrientation: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDagRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	name := d.Get(mkDataSourceDagName).(string)
	res, err := c.GetDagDetailsWithResponse(
		ctx,
		api.DAGID(name),
	)

	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	_ = d.Set(mkDataSourceDagName, res.JSON200.DagId)
	_ = d.Set(mkDataSourceDagDescription, res.JSON200.Description)
	_ = d.Set(mkDataSourceDagPaused, res.JSON200.IsPaused)
	_ = d.Set(mkDataSourceDagSubDag, res.JSON200.IsSubdag)
	_ = d.Set(mkDataSourceDagRootDagId, res.JSON200.RootDagId)
	_ = d.Set(mkDataSourceDagStartDate, res.JSON200.StartDate)
	_ = d.Set(mkDataSourceDagTimezone, res.JSON200.Timezone)
	_ = d.Set(mkDataSourceDagCatchup, res.JSON200.Catchup)
	_ = d.Set(mkDataSourceDagConcurrency, res.JSON200.Concurrency)
	_ = d.Set(mkDataSourceDagDefaultView, res.JSON200.DefaultView)
	_ = d.Set(mkDataSourceDagOrientation, res.JSON200.Orientation)

	d.SetId(*res.JSON200.DagId)

	return diags
}
