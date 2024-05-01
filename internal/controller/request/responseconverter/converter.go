package responseconverter

import (
	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	httpClient "github.com/crossplane-contrib/provider-http/internal/clients/http"
	"k8s.io/apimachinery/pkg/runtime"
)

// Convert HttpResponse to Response
func HttpResponseToV1alpha1Response(httpResponse httpClient.HttpResponse) v1alpha1.Response {
	r := runtime.RawExtension{Raw: []byte("{}")}
	if len(httpResponse.Body) > 0 {
		r.Raw = []byte(httpResponse.Body)
	}
	return v1alpha1.Response{
		StatusCode: httpResponse.StatusCode,
		Body:       r,
		Headers:    httpResponse.Headers,
	}
}
