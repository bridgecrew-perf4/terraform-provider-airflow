package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const (
	mkResourceVariableName  = "name"
	mkResourceVariableValue = "value"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableCreate,
		ReadContext:   resourceVariableRead,
		UpdateContext: resourceVariableUpdate,
		DeleteContext: resourceVariableDelete,
		Schema: map[string]*schema.Schema{
			mkResourceVariableName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			mkResourceVariableValue: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVariableCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	name := d.Get(mkResourceVariableName).(string)
	value := d.Get(mkResourceVariableValue).(string)

	body := api.PostVariablesJSONRequestBody{
		VariableCollectionItem: api.VariableCollectionItem{
			Key: &name,
		},
		Value: &value,
	}

	res, err := c.PostVariablesWithResponse(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	d.SetId(name)

	return resourceVariableRead(ctx, d, m)
}

func resourceVariableRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	variableId := d.Id()
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

	if err := d.Set(mkResourceVariableValue, res.JSON200.Value); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVariableUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	variableId := d.Id()
	value := d.Get(mkResourceVariableValue).(string)

	body := api.PatchVariableJSONRequestBody{
		VariableCollectionItem: api.VariableCollectionItem{
			Key: &variableId,
		},
		Value: &value,
	}

	res, err := c.PatchVariableWithResponse(
		ctx,
		api.VariableKey(variableId),
		&api.PatchVariableParams{},
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

	return resourceVariableRead(ctx, d, m)
}

func resourceVariableDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	variableId := d.Id()

	res, err := c.DeleteVariableWithResponse(ctx, api.VariableKey(variableId))
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
