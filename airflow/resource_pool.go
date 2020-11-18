package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaResourcePool = "airflow_pool"

const (
	mkResourcePoolName  = "name"
	mkResourcePoolSlots = "slots"
)

const (
	fieldPoolName  = "name"
	fieldPoolSlots = "slots"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolCreate,
		ReadContext:   resourcePoolRead,
		UpdateContext: resourcePoolUpdate,
		DeleteContext: resourcePoolDelete,
		Schema: map[string]*schema.Schema{
			mkResourcePoolName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			mkResourcePoolSlots: {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.IntAtLeast(-1),
				),
			},
		},
	}
}

func resourcePoolCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	name := d.Get(mkResourcePoolName).(string)
	slots := d.Get(mkResourcePoolSlots).(int)

	body := api.PostPoolJSONRequestBody{
		Name:  &name,
		Slots: &slots,
	}

	res, err := c.PostPoolWithResponse(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	d.SetId(name)

	return resourcePoolRead(ctx, d, m)
}

func resourcePoolRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	poolId := d.Id()
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

	if err := d.Set(mkResourcePoolSlots, res.JSON200.Slots); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePoolUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	poolId := d.Id()

	params := api.PatchPoolParams{
		UpdateMask: new(api.UpdateMask),
	}

	body := api.PatchPoolJSONRequestBody{
		Name: &poolId,
	}

	if d.HasChange(mkResourcePoolSlots) {
		*params.UpdateMask = append(*params.UpdateMask, fieldPoolSlots)
		slots := d.Get(mkResourcePoolSlots).(int)
		body.Slots = &slots
	}

	res, err := c.PatchPoolWithResponse(
		ctx,
		api.PoolName(poolId),
		&params,
		body,
	)

	if err != nil {
		return diag.FromErr(err)
	} else if res.StatusCode() != 204 {
		if d := helper.GetResponseDiag(res); d != nil {
			diags = append(diags, *d)
			return diags
		}
	}

	return resourcePoolRead(ctx, d, m)
}

func resourcePoolDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	poolId := d.Id()

	res, err := c.DeletePoolWithResponse(
		ctx,
		api.PoolName(poolId),
	)
	if err != nil {
		return diag.FromErr(err)
	} else if res.StatusCode() != 204 {
		if d := helper.GetResponseDiag(res); d != nil {
			diags = append(diags, *d)
			return diags
		}
	}

	d.SetId("")

	return diags
}
