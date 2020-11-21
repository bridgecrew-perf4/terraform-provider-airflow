package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xabinapal/terraform-provider-airflow/api"
	"github.com/xabinapal/terraform-provider-airflow/helper"
)

const schemaResourceConnection = "airflow_connection"

const (
	mkResourceConnectionName     = "name"
	mkResourceConnectionType     = "type"
	mkResourceConnectionSchema   = "schema"
	mkResourceConnectionHost     = "host"
	mkResourceConnectionPort     = "port"
	mkResourceConnectionLogin    = "login"
	mkResourceConnectionPassword = "password"
	mkResourceConnectionExtra    = "extra"
)

const (
	fieldConnectionType     = "conn_type"
	fieldConnectionSchema   = "schema"
	fieldConnectionHost     = "host"
	fieldConnectionPort     = "port"
	fieldConnectionLogin    = "login"
	fieldConnectionPassword = "password"
	fieldConnectionExtra    = "extra"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
		Schema: map[string]*schema.Schema{
			mkResourceConnectionName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 250),
				),
			},
			mkResourceConnectionType: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringInSlice(getConnectionTypesAsList(), false),
				),
			},
			mkResourceConnectionSchema: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 500),
				),
			},
			mkResourceConnectionHost: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 500),
				),
			},
			mkResourceConnectionPort: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			mkResourceConnectionLogin: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 500),
				),
			},
			mkResourceConnectionPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 5000),
				),
			},
			mkResourceConnectionExtra: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: helper.ValidateDiagFunc(
					validation.StringLenBetween(0, 5000),
				),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	name := d.Get(mkResourceConnectionName).(string)
	connType := d.Get(mkResourceConnectionType).(string)
	connSchema := d.Get(mkResourceConnectionSchema).(string)
	connHost := d.Get(mkResourceConnectionHost).(string)
	connPort := d.Get(mkResourceConnectionPort).(int)
	connLogin := d.Get(mkResourceConnectionLogin).(string)
	connPassword := d.Get(mkResourceConnectionPassword).(string)
	connExtra := d.Get(mkResourceConnectionExtra).(string)

	body := api.PostConnectionJSONRequestBody{
		ConnectionCollectionItem: api.ConnectionCollectionItem{
			ConnectionId: &name,
			ConnType:     &connType,
			Schema:       &connSchema,
			Host:         &connHost,
			Port:         &connPort,
			Login:        &connLogin,
		},
		Password: &connPassword,
		Extra:    &connExtra,
	}

	res, err := c.PostConnectionWithResponse(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	} else if d := helper.GetResponseDiag(res); d != nil {
		diags = append(diags, *d)
		return diags
	}

	d.SetId(name)

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	connectionId := d.Id()
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

	if err := d.Set(mkResourceConnectionName, res.JSON200.ConnectionId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(mkResourceConnectionType, res.JSON200.ConnType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(mkResourceConnectionSchema, res.JSON200.Schema); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(mkResourceConnectionHost, res.JSON200.Host); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(mkResourceConnectionPort, res.JSON200.Port); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(mkResourceConnectionLogin, res.JSON200.Login); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceConnectionUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	connectionId := d.Id()

	params := api.PatchConnectionParams{
		UpdateMask: new(api.UpdateMask),
	}

	body := api.PatchConnectionJSONRequestBody{
		ConnectionCollectionItem: api.ConnectionCollectionItem{
			ConnectionId: &connectionId,
		},
	}

	if d.HasChange(mkResourceConnectionType) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionType)
		connType := d.Get(mkResourceConnectionType).(string)
		body.ConnType = &connType
	}

	if d.HasChange(mkResourceConnectionSchema) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionSchema)
		connSchema := d.Get(mkResourceConnectionSchema).(string)
		body.Schema = &connSchema
	}

	if d.HasChange(mkResourceConnectionHost) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionHost)
		connHost := d.Get(mkResourceConnectionHost).(string)
		body.Host = &connHost
	}

	if d.HasChange(mkResourceConnectionPort) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionPort)
		connPort := d.Get(mkResourceConnectionPort).(int)
		body.Port = &connPort
	}

	if d.HasChange(mkResourceConnectionLogin) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionLogin)
		connLogin := d.Get(mkResourceConnectionLogin).(string)
		body.Login = &connLogin
	}

	if d.HasChange(mkResourceConnectionPassword) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionPassword)
		connPassword := d.Get(mkResourceConnectionPassword).(string)
		body.Password = &connPassword
	}

	if d.HasChange(mkResourceConnectionExtra) {
		*params.UpdateMask = append(*params.UpdateMask, fieldConnectionExtra)
		connExtra := d.Get(mkResourceConnectionExtra).(string)
		body.Extra = &connExtra
	}

	res, err := c.PatchConnectionWithResponse(
		ctx,
		api.ConnectionID(connectionId),
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

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(api.ClientWithResponsesInterface)

	var diags diag.Diagnostics

	connectionId := d.Id()

	res, err := c.DeleteConnectionWithResponse(
		ctx,
		api.ConnectionID(connectionId),
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
