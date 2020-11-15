package helper

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValidateDiagFunc(
	fn schema.SchemaValidateFunc,
) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		warnings, errs := fn(i, fmt.Sprintf("%+v", path))
		for _, warning := range warnings {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       warning,
				Detail:        warning,
				AttributePath: path,
			})
		}
		for _, err := range errs {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				Detail:        err.Error(),
				AttributePath: path,
			})
		}

		return diags
	}
}
