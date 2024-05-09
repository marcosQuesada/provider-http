package requestgen

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	"github.com/crossplane-contrib/provider-http/internal/controller/request/requestprocessing"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

var testHeaders = map[string][]string{
	"fruits":                {"apple", "banana", "orange"},
	"colors":                {"red", "green", "blue"},
	"countries":             {"USA", "UK", "India", "Germany"},
	"programming_languages": {"Go", "Python", "JavaScript"},
}

var testHeaders2 = map[string][]string{
	"countries": {"USA", "UK", "India", "Germany"},
}

var (
	testPostMapping = v1alpha1.Mapping{
		Method:  "POST",
		Body:    runtime.RawExtension{Raw: []byte(`{"username": ".payload.body.username", "email": ".payload.body.email"}`)},
		URL:     ".payload.baseUrl",
		Headers: testHeaders,
	}

	testPutMapping = v1alpha1.Mapping{
		Method: "PUT",
		Body:   runtime.RawExtension{Raw: []byte(`{"username": "john_doe_new_username" }`)},
		//Body:    runtime.RawExtension{Raw: []byte("{\"username\": \"\\x34\\john_doe_new_username\\x34\\\" }")},
		URL:     "(.payload.baseUrl + \"/\" + .response.body.id)",
		Headers: testHeaders,
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

var (
	testForProvider = v1alpha1.RequestParameters{
		Payload: v1alpha1.Payload{
			Body:    runtime.RawExtension{Raw: []byte(`{"username": "john_doe", "email": "john.doe@example.com"}`)},
			BaseUrl: "https://api.example.com/users",
		},
		Mappings: v1alpha1.Mappings{
			&testPostMapping,
			&testGetMapping,
			&testPutMapping,
			&testDeleteMapping,
		},
		Headers: map[string][]string{},
	}
)

func Test_GenerateRequestDetails(t *testing.T) {
	type args struct {
		methodMapping v1alpha1.Mapping
		forProvider   v1alpha1.RequestParameters
		response      v1alpha1.Response
		logger        logging.Logger
	}
	type want struct {
		requestDetails RequestDetails
		err            error
		ok             bool
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"SuccessPost": {
			args: args{
				methodMapping: testPostMapping,
				forProvider:   testForProvider,
				response:      v1alpha1.Response{},
				logger:        logging.NewNopLogger(),
			},
			want: want{
				requestDetails: RequestDetails{
					Url:     "https://api.example.com/users",
					Body:    `{"email":"john.doe@example.com","username":"john_doe"}`,
					Headers: testHeaders,
				},
				err: nil,
				ok:  true,
			},
		},
		"SuccessPut": {
			args: args{
				methodMapping: testPutMapping,
				forProvider:   testForProvider,
				response: v1alpha1.Response{
					StatusCode: 200,
					Body:       runtime.RawExtension{Raw: []byte(`{"id":"123","username":"john_doe"}`)},
					Headers:    testHeaders,
				},
				logger: logging.NewNopLogger(),
			},
			want: want{
				requestDetails: RequestDetails{
					Url:     "https://api.example.com/users/123",
					Body:    `{"username":"john_doe_new_username"}`,
					Headers: testHeaders,
				},
				err: nil,
				ok:  true,
			},
		},
		"SuccessDelete": {
			args: args{
				methodMapping: testDeleteMapping,
				forProvider:   testForProvider,
				response: v1alpha1.Response{
					StatusCode: 200,
					Body:       runtime.RawExtension{Raw: []byte(`{"id":"123","username":"john_doe"}`)},
					Headers:    testHeaders,
				},
				logger: logging.NewNopLogger(),
			},
			want: want{
				requestDetails: RequestDetails{
					Url:     "https://api.example.com/users/123",
					Headers: map[string][]string{},
				},
				err: nil,
				ok:  true,
			},
		},
		"SuccessGet": {
			args: args{
				methodMapping: testGetMapping,
				forProvider:   testForProvider,
				response: v1alpha1.Response{
					StatusCode: 200,
					Body:       runtime.RawExtension{Raw: []byte(`{"id":"123","username":"john_doe"}`)},
					Headers:    testHeaders,
				},
				logger: logging.NewNopLogger(),
			},
			want: want{
				requestDetails: RequestDetails{
					Url:     "https://api.example.com/users/123",
					Headers: map[string][]string{},
				},
				err: nil,
				ok:  true,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, gotErr, ok := GenerateRequestDetails(tc.args.methodMapping, tc.args.forProvider, tc.args.response)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("GenerateRequestDetails(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.ok, ok); diff != "" {
				t.Fatalf("GenerateRequestDetails(...): -want ok, +got ok: %s", diff)
			}

			if diff := cmp.Diff(tc.want.requestDetails, got); diff != "" {
				t.Errorf("GenerateRequestDetails(...): -want result, +got result: %s", diff)
			}
		})
	}

}

func Test_IsRequestValid(t *testing.T) {
	type args struct {
		requestDetails RequestDetails
	}
	type want struct {
		ok bool
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ValidRequestDetails": {
			args: args{
				requestDetails: RequestDetails{
					Body:    `{"id": "123", "username": "john_doe"}`,
					Url:     "https://example",
					Headers: nil,
				},
			},
			want: want{
				ok: true,
			},
		},
		"NonValidRequestDetails": {
			args: args{
				requestDetails: RequestDetails{
					Body:    "",
					Url:     "",
					Headers: nil,
				},
			},
			want: want{
				ok: false,
			},
		},
		"NonValidUrl": {
			args: args{
				requestDetails: RequestDetails{
					Body:    `{"id": "123", "username": "john_doe"}`,
					Url:     "",
					Headers: nil,
				},
			},
			want: want{
				ok: false,
			},
		},
		"NonValidBody": {
			args: args{
				requestDetails: RequestDetails{
					Body:    `{"id": "null", "username": "john_doe"}`,
					Url:     "https://example",
					Headers: nil,
				},
			},
			want: want{
				ok: false,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := IsRequestValid(tc.args.requestDetails)
			if diff := cmp.Diff(tc.want.ok, got); diff != "" {
				t.Fatalf("IsRequestValid(...): -want bool, +got bool: %s", diff)
			}
		})
	}

}

func Test_coalesceHeaders(t *testing.T) {
	type args struct {
		mappingHeaders,
		defaultHeaders map[string][]string
	}
	type want struct {
		headers map[string][]string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"NonNilMappingHeaders": {
			args: args{
				mappingHeaders: testHeaders,
				defaultHeaders: testHeaders2,
			},
			want: want{
				headers: testHeaders,
			},
		},
		"NilMappingHeaders": {
			args: args{
				mappingHeaders: nil,
				defaultHeaders: testHeaders2,
			},
			want: want{
				headers: testHeaders2,
			},
		},
		"NilDefaultHeaders": {
			args: args{
				mappingHeaders: testHeaders,
				defaultHeaders: nil,
			},
			want: want{
				headers: testHeaders,
			},
		},
		"NilHeaders": {
			args: args{
				mappingHeaders: nil,
				defaultHeaders: nil,
			},
			want: want{
				headers: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := coalesceHeaders(tc.args.mappingHeaders, tc.args.defaultHeaders)
			if diff := cmp.Diff(tc.want.headers, got); diff != "" {
				t.Fatalf("coalesceHeaders(...): -want headers, +got headers: %s", diff)
			}
		})
	}
}

func Test_generateRequestObject(t *testing.T) {
	type args struct {
		forProvider v1alpha1.RequestParameters
		response    v1alpha1.Response
	}
	type want struct {
		result map[string]interface{}
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"Success": {
			args: args{
				forProvider: testForProvider,
				response: v1alpha1.Response{
					StatusCode: 200,
					Body:       runtime.RawExtension{Raw: []byte(`{"id": "123"}`)},
					Headers:    nil,
				},
			},
			want: want{
				result: map[string]any{
					"mappings": []any{
						map[string]any{
							"body":   "{ username: .payload.body.username, email: .payload.body.email }",
							"method": "POST",
							"headers": map[string]any{
								"colors":                []any{"red", "green", "blue"},
								"countries":             []any{"USA", "UK", "India", "Germany"},
								"fruits":                []any{"apple", "banana", "orange"},
								"programming_languages": []any{"Go", "Python", "JavaScript"},
							},
							"url": ".payload.baseUrl",
						},
						map[string]any{
							"method": "GET",
							"url":    `(.payload.baseUrl + "/" + .response.body.id)`,
						},
						map[string]any{
							"body":   `{ username: "john_doe_new_username" }`,
							"method": "PUT",
							"headers": map[string]any{
								"colors":                []any{"red", "green", "blue"},
								"countries":             []any{"USA", "UK", "India", "Germany"},
								"fruits":                []any{"apple", "banana", "orange"},
								"programming_languages": []any{"Go", "Python", "JavaScript"},
							},
							"url": `(.payload.baseUrl + "/" + .response.body.id)`,
						},
						map[string]any{
							"method": "DELETE",
							"url":    `(.payload.baseUrl + "/" + .response.body.id)`,
						},
					},
					"payload": map[string]any{
						"baseUrl": "https://api.example.com/users",
						"body":    map[string]any{"email": "john.doe@example.com", "username": "john_doe"},
					},
					"response": map[string]any{
						"body":       map[string]any{"id": "123"},
						"statusCode": float64(200),
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Skip()
			got := generateRequestObject(tc.args.forProvider, tc.args.response)
			if diff := cmp.Diff(tc.want.result, got); diff != "" {
				t.Fatalf("generateRequestObject(...): -want result, +got result: %s", diff)
			}
		})
	}
}
func TestEmbeddedRawExtensionMarshal(t *testing.T) {
	type test struct {
		Ext runtime.RawExtension
	}

	extension := test{Ext: runtime.RawExtension{Raw: []byte(`{foo:{"foo":"bar"}`)}}
	data, err := json.Marshal(extension)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"Ext":{"foo":"bar"}}` {
		t.Errorf("unexpected data: %s", string(data))
	}
}

func TestOriginalBodyMappers(t *testing.T) {
	mappingBody := "{ username: .payload.body.username, email: .payload.body.email }"
	jqQuery := requestprocessing.ConvertStringToJQQuery(mappingBody)

	response := v1alpha1.Response{}
	jqObject := generateRequestObject(testForProvider, response)
	body, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
	require.NoError(t, err)

	require.Equal(t, `{"email":"john.doe@example.com","username":"john_doe"}`, body)
}

func TestOriginalBodyMappersFromStringKeysAndValuesAsPureJQ(t *testing.T) {
	mappingBody := `{"username": ".payload.body.username", "email": ".payload.body.email"}`

	mappingBody = strings.ReplaceAll(mappingBody, "\"", "") // @TODO: How to
	jqQuery := requestprocessing.ConvertStringToJQQuery(mappingBody)
	response := v1alpha1.Response{}
	jqObject := generateRequestObject(testForProvider, response)
	body, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
	require.NoError(t, err)

	require.Equal(t, `{"email":"john.doe@example.com","username":"john_doe"}`, body)
}

func TestOriginalBodyMappersFromStringKeysAndValuesAsMixedJQ(t *testing.T) {
	mappingBody := `{"username": "john_doe_new_username", "email": ".payload.body.email"}`

	mappingBody = strings.ReplaceAll(mappingBody, "\"", "") // @TODO: Initial workaround, not valid for non JQ queries
	jqQuery := requestprocessing.ConvertStringToJQQuery(mappingBody)
	response := v1alpha1.Response{}
	jqObject := generateRequestObject(testForProvider, response)
	_, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
	require.Error(t, err)
}

// Example of mixed JQ values and literals
func TestOriginalBodyMappersFromStringKeysAndValuesAsMixedJQWorkingWithLiterals(t *testing.T) {
	mappingBody := `{"username": "john_doe_new_username", "email": .payload.body.email}`

	jqQuery := requestprocessing.ConvertStringToJQQuery(mappingBody)
	response := v1alpha1.Response{}
	jqObject := generateRequestObject(testForProvider, response)
	body, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
	require.NoError(t, err)
	require.Equal(t, `{"email":"john.doe@example.com","username":"john_doe_new_username"}`, body)
}

func TestFromStringKeyValueToJQ(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name                            string
		args                            args
		want                            string
		wantFromStringKeyValueToJQError bool
		wantApplyJQOnStrError           bool
	}{
		{
			name: "RawExtensionSupportingJQFields",
			args: args{input: `{"username": "john_doe_new_username", "email": ".payload.body.email"}`},
			want: `{"email":"john.doe@example.com","username":"john_doe_new_username"}`,
		},
		{
			name: "RawExtensionSupportingNestedJQFields",
			args: args{input: `{ "username": "john_doe_new_username", "email": {"foo":".payload.body.email"} }`},
			want: `{"email":{"foo":"john.doe@example.com"},"username":"john_doe_new_username"}`,
		},
		{
			name: "RawExtensionSupportingNestedMultipleJQFields",
			args: args{input: `{ "username": "john_doe_new_username", "email": {"foo":".payload.body.email", "bar":".payload.body.username"} }`},
			want: `{"email":{"bar":"john_doe","foo":"john.doe@example.com"},"username":"john_doe_new_username"}`,
		},
		{
			name: "RawExtensionSupportingNestedWithSingleItemArrayJQFields",
			args: args{input: `{ "username": "john_doe_new_username", "email": [".payload.body.email"] }`},
			want: `{"email":["john.doe@example.com"],"username":"john_doe_new_username"}`,
		},
		{
			name: "RawExtensionSupportingNestedWithArrayJQFields",
			args: args{input: `{ "username": "john_doe_new_username", "email": [".payload.body.email",".payload.body.username"] }`},
			want: `{"email":["john.doe@example.com","john_doe"],"username":"john_doe_new_username"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromJSONMapToJQuery(tt.args.input)
			if tt.wantFromStringKeyValueToJQError != (err != nil) {
				t.Fatalf("unexpected wantFromStringKeyValueToJQ error, got %v", err)
			}

			jqQuery := requestprocessing.ConvertStringToJQQuery(result)
			response := v1alpha1.Response{}
			jqObject := generateRequestObject(testForProvider, response)
			got, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
			if tt.wantApplyJQOnStrError != (err != nil) {
				t.Fatalf("unexpected wantApplyJQOnStr error, got %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRequestDetails() got = %v, want %v", got, tt.want)
			}

		})
	}
}
