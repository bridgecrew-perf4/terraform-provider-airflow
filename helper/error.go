package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/xabinapal/terraform-provider-airflow/api"
)

func GetResponseDiag(obj interface{}) *diag.Diagnostic {
	res := parseResponse(obj)

	if res.HTTPResponse == nil {
		return &diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "0 - Unknown error",
			Detail:   "HTTPResponse field is nil",
		}
	}

	status := res.HTTPResponse.StatusCode
	if status == 200 {
		return nil
	}

	apiErr, err := parseErrorResponse(res)
	if err != nil {
		return &diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%d - Unknown error", status),
			Detail:   err.Error(),
		}
	}

	var detail string
	if apiErr.Detail != nil {
		detail = *apiErr.Detail
	}

	return &diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("%d - %s", status, apiErr.Title),
		Detail:   detail,
	}
}

type response struct {
	HTTPResponse *http.Response
	Body         []byte
}

func parseResponse(res interface{}) *response {
	obj := reflect.ValueOf(res).Elem()

	resObj := obj.FieldByName("HTTPResponse")
	if !resObj.IsValid() {
		panic("Invalid response object")
	}

	bodyObj := obj.FieldByName("Body")
	if !bodyObj.IsValid() {
		panic("Invalid response object")
	}

	return &response{
		HTTPResponse: resObj.Interface().(*http.Response),
		Body:         bodyObj.Interface().([]byte),
	}
}

func parseErrorResponse(res *response) (*api.Error, error) {
	if !strings.Contains(res.HTTPResponse.Header.Get("Content-Type"), "json") {
		return nil, fmt.Errorf("HTTPResponse field is not JSON")
	}

	if res.Body == nil {
		return nil, fmt.Errorf("Body field is nil")
	}

	var dest api.Error
	if err := json.Unmarshal(res.Body, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
