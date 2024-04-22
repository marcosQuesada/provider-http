package request

import (
	"testing"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

var (
	testPostMapping = v1alpha1.Mapping{
		Method: "POST",
		Body:   "{ username: .payload.body.username, email: .payload.body.email }",
		URL:    ".payload.baseUrl",
	}

	testPutMapping = v1alpha1.Mapping{
		Method: "PUT",
		Body:   "{ username: \"john_doe_new_username\" }",
		URL:    "(.payload.baseUrl + \"/\" + .response.body.id)",
	}

	testGetMapping = v1alpha1.Mapping{
		Method: "GET",
		URL:    "(.payload.baseUrl + \"/\" + .response.body.id)",
	}

	testDeleteMapping = v1alpha1.Mapping{
		Method: "DELETE",
		URL:    "(.payload.baseUrl + \"/\" + .response.body.id)",
	}
)

func Test_getMappingByMethod(t *testing.T) {
	type args struct {
		requestParams *v1alpha1.RequestParameters
		method        Action
	}
	type want struct {
		mapping *v1alpha1.Mapping
		ok      bool
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"Fail": {
			args: args{
				requestParams: &v1alpha1.RequestParameters{
					Payload: v1alpha1.Payload{
						Body:    "{\"username\": \"john_doe\", \"email\": \"john.doe@example.com\"}",
						BaseUrl: "https://api.example.com/users",
					},
					Mappings: []v1alpha1.Mapping{
						testGetMapping,
						testPutMapping,
						testDeleteMapping,
					},
				},
				method: "POST",
			},
			want: want{
				mapping: nil,
				ok:      false,
			},
		},
		"Success": {
			args: args{
				requestParams: &v1alpha1.RequestParameters{
					Payload: v1alpha1.Payload{
						Body:    "{\"username\": \"john_doe\", \"email\": \"john.doe@example.com\"}",
						BaseUrl: "https://api.example.com/users",
					},
					Mappings: []v1alpha1.Mapping{
						testPostMapping,
						testGetMapping,
						testPutMapping,
						testDeleteMapping,
					},
				},
				method: "POST",
			},
			want: want{
				mapping: &testPostMapping,
				ok:      true,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Skip()
			got, ok := getMappingByAction(tc.args.requestParams, tc.args.method)
			if diff := cmp.Diff(tc.want.mapping, got); diff != "" {
				t.Fatalf("getMappingByAction(...): -want result, +got result: %s", diff)
			}

			if diff := cmp.Diff(tc.want.ok, ok); diff != "" {
				t.Fatalf("getMappingByAction(...): -want result, +got result: %s", diff)
			}
		})
	}
}
