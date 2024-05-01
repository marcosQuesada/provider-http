package v1alpha1

import (
	"strings"

	"github.com/crossplane-contrib/provider-http/internal/json"
)

func (d *DisposableRequest) SetStatusCode(statusCode int) {
	d.Status.Response.StatusCode = statusCode
}

func (d *DisposableRequest) SetHeaders(headers map[string][]string) {
	d.Status.Response.Headers = headers
}

func (d *DisposableRequest) SetBody(body string) {
	d.Status.Response.RawBody = body
	d.Status.Response.Body.Raw = []byte("{}")
	if len(body) > 0 && json.IsJSONString(body) {
		d.Status.Response.Body.Raw = []byte(body)
	}
}

func (d *DisposableRequest) SetSynced(synced bool) {
	d.Status.Synced = synced
	d.Status.Failed = 0
	d.Status.Error = ""
}

func (d *DisposableRequest) SetError(err error) {
	d.Status.Failed++
	d.Status.Synced = true
	if err != nil {
		d.Status.Error = err.Error()
	}
}

func (d *DisposableRequest) SetRequestDetails(url, method, body string, headers map[string][]string) {
	d.Status.RequestDetails.Body.Raw = []byte("{}")
	body = strings.Trim(body, " ")
	if len(body) > 0 {
		d.Status.RequestDetails.Body.Raw = []byte(body)
	}
	d.Status.RequestDetails.URL = url
	d.Status.RequestDetails.Headers = headers
	d.Status.RequestDetails.Method = method
}
