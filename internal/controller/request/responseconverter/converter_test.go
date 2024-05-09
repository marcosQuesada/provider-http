package responseconverter

import (
	"testing"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	httpClient "github.com/crossplane-contrib/provider-http/internal/clients/http"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
)

var testHeaders = map[string][]string{
	"fruits":                {"apple", "banana", "orange"},
	"colors":                {"red", "green", "blue"},
	"countries":             {"USA", "UK", "India", "Germany"},
	"programming_languages": {"Go", "Python", "JavaScript"},
}

func Test_HttpResponseToV1alpha1Response(t *testing.T) {
	type args struct {
		httpResponse httpClient.HttpResponse
	}
	type want struct {
		result v1alpha1.Response
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"Success": {
			args: args{
				httpResponse: httpClient.HttpResponse{
					Body:       `{"email":"john.doe@example.com","name":"john_doe"}`,
					Headers:    testHeaders,
					StatusCode: 200,
				},
			},
			want: want{
				result: v1alpha1.Response{
					Body:       runtime.RawExtension{Raw: []byte(`{"email":"john.doe@example.com","name":"john_doe"}`)},
					Headers:    testHeaders,
					StatusCode: 200,
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := HttpResponseToV1alpha1Response(tc.args.httpResponse)
			if diff := cmp.Diff(tc.want.result, got); diff != "" {
				t.Fatalf("HttpResponseToV1alpha1Response(...): -want result, +got result: %s", diff)
			}
		})
	}

}
