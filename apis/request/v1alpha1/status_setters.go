package v1alpha1

import (
	"time"
)

func (d *Request) SetStatusCode(statusCode int) {
	d.Status.Response.StatusCode = statusCode
}

func (d *Request) SetHeaders(headers map[string][]string) {
	d.Status.Response.Headers = headers
}

func (d *Request) SetBody(body string) {
	d.Status.Response.Body.Raw = []byte("{}")
	if len(body) > 0 {
		d.Status.Response.Body.Raw = []byte(body)
	}
}

func (d *Request) SetError(err error) {
	d.Status.Failed++
	if err != nil {
		d.Status.Error = err.Error()
	}
}

func (d *Request) ResetFailures() {
	d.Status.Failed = 0
	d.Status.Error = ""
}

func (d *Request) SetRequestDetails(url, method, body string, headers map[string][]string) {
	d.Status.RequestDetails.URL = url
	d.Status.RequestDetails.Headers = headers
	d.Status.RequestDetails.Method = method
	d.Status.RequestDetails.Body.Raw = []byte("{}")
	if len(body) > 0 {
		d.Status.RequestDetails.Body.Raw = []byte(body)
	}
}

func (d *Request) SetCache(statusCode int, headers map[string][]string, body string) {
	d.Status.Cache.Response.StatusCode = statusCode
	d.Status.Cache.Response.Headers = headers
	d.Status.Cache.Response.Body.Raw = []byte("{}")
	if len(body) > 0 {
		d.Status.Cache.Response.Body.Raw = []byte(body)
	}
	d.Status.Cache.LastUpdated = time.Now().UTC().Format(time.RFC3339)
}
