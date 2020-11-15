package helper

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/xabinapal/terraform-provider-airflow/api"
)

func GetResponseDiag(res interface{}) *diag.Diagnostic {
	obj := reflect.ValueOf(res).Elem()

	statusFn := obj.MethodByName("StatusCode")
	if !statusFn.IsValid() {
		panic("Invalid response object")
	}

	statusVal := statusFn.Call([]reflect.Value{})
	if len(statusVal) != 1 || !statusVal[0].IsValid() ||
		statusVal[0].Kind() != reflect.Int {
		panic("Invalid response object")
	}

	status := statusVal[0].Interface().(int)
	if status == 200 {
		return nil
	}

	errVal := obj.FieldByName(fmt.Sprintf("JSON%d", status))
	if !errVal.IsValid() {
		return &diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%d - Unknown error", status),
		}
	} else if errVal.Type() != reflect.TypeOf((*api.Error)(nil)) {
		panic("Invalid response object")
	}

	err := errVal.Interface().(*api.Error)

	var detail string
	if err.Detail != nil {
		detail = *err.Detail
	}

	return &diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("%d - %s", status, err.Title),
		Detail:   detail,
	}
}
