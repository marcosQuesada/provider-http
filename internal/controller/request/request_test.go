package request

import (
	"context"
	"net/http"
	"testing"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	httpClient "github.com/crossplane-contrib/provider-http/internal/clients/http"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

var (
	errBoom = errors.New("boom")
)

const (
	providerName    = "http-test"
	testRequestName = "test-request"
	testNamespace   = "testns"
)

var (
	testForProvider = v1alpha1.RequestParameters{
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
	}
)

type httpRequestModifier func(request *v1alpha1.Request)

func httpRequest(rm ...httpRequestModifier) *v1alpha1.Request {
	r := &v1alpha1.Request{
		ObjectMeta: v1.ObjectMeta{
			Name:      testRequestName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.RequestSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerName,
				},
			},
			ForProvider: testForProvider,
		},
		Status: v1alpha1.RequestStatus{},
	}

	for _, m := range rm {
		m(r)
	}

	return r
}

type notHttpRequest struct {
	resource.Managed
}

type MockSendRequestFn func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error)

type MockHttpClient struct {
	MockSendRequest MockSendRequestFn
}

func (c *MockHttpClient) SendRequest(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
	return c.MockSendRequest(ctx, method, url, body, headers, skipTLSVerify)
}

type MockSetRequestStatusFn func() error

type MockResetFailuresFn func()

type MockInitFn func(ctx context.Context, cr *v1alpha1.Request, res httpClient.HttpResponse)

type MockStatusHandler struct {
	MockSetRequest    MockSetRequestStatusFn
	MockResetFailures MockResetFailuresFn
}

func (s *MockStatusHandler) ResetFailures() {
	s.MockResetFailures()
}

func (s *MockStatusHandler) SetRequestStatus(ctx context.Context, cr *v1alpha1.Request, res httpClient.HttpResponse, err error) error {
	return s.MockSetRequest()
}

func Test_httpExternal_Create(t *testing.T) {
	type args struct {
		http      httpClient.Client
		localKube client.Client
		mg        resource.Managed
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"NotRequestResource": {
			args: args{
				mg: notHttpRequest{},
			},
			want: want{
				err: errors.New(errNotRequest),
			},
		},
		"RequestFailed": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, errBoom
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: errors.Wrap(errBoom, errFailedToSendHttpRequest),
			},
		},
		"Success": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, nil
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockCreate:       test.NewMockCreateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		tc := tc // Create local copies of loop variables

		t.Run(name, func(t *testing.T) {
			t.Skip()
			e := &external{
				localKube: tc.args.localKube,
				logger:    logging.NewNopLogger(),
				http:      tc.args.http,
			}
			_, gotErr := e.Create(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Create(...): -want error, +got error: %s", diff)
			}
		})
	}
}

func Test_httpExternal_Update(t *testing.T) {
	type args struct {
		http      httpClient.Client
		localKube client.Client
		mg        resource.Managed
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"NotRequestResource": {
			args: args{
				mg: notHttpRequest{},
			},
			want: want{
				err: errors.New(errNotRequest),
			},
		},
		"RequestFailed": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, errBoom
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: errors.Wrap(errBoom, errFailedToSendHttpRequest),
			},
		},
		"Success": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, nil
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockCreate:       test.NewMockCreateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		tc := tc // Create local copies of loop variables

		t.Run(name, func(t *testing.T) {
			t.Skip()
			e := &external{
				localKube: tc.args.localKube,
				logger:    logging.NewNopLogger(),
				http:      tc.args.http,
			}
			_, gotErr := e.Update(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Update(...): -want error, +got error: %s", diff)
			}
		})
	}
}

func Test_httpExternal_Delete(t *testing.T) {
	type args struct {
		http      httpClient.Client
		localKube client.Client
		mg        resource.Managed
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"NotRequestResource": {
			args: args{
				mg: notHttpRequest{},
			},
			want: want{
				err: errors.New(errNotRequest),
			},
		},
		"RequestFailed": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, errBoom
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: errors.Wrap(errBoom, errFailedToSendHttpRequest),
			},
		},
		"Success": {
			args: args{
				http: &MockHttpClient{
					MockSendRequest: func(ctx context.Context, method string, url string, body string, headers map[string][]string, skipTLSVerify bool) (resp httpClient.HttpDetails, err error) {
						return httpClient.HttpDetails{}, nil
					},
				},
				localKube: &test.MockClient{
					MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
					MockCreate:       test.NewMockCreateFn(nil),
					MockGet:          test.NewMockGetFn(nil),
				},
				mg: httpRequest(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		tc := tc // Create local copies of loop variables

		t.Run(name, func(t *testing.T) {
			t.Skip()
			e := &external{
				localKube: tc.args.localKube,
				logger:    logging.NewNopLogger(),
				http:      tc.args.http,
			}
			gotErr := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Delete(...): -want error, +got error: %s", diff)
			}
		})
	}
}

func TestPatchFieldValueToObject(t *testing.T) {

	to := &v1alpha1.Request{
		TypeMeta: v1.TypeMeta{
			APIVersion: "http.crossplane.io/v1alpha1",
			Kind:       "Request",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "xxx-request",
			Namespace: "Default",
		},
		Spec: v1alpha1.RequestSpec{
			ForProvider: v1alpha1.RequestParameters{
				Mappings: v1alpha1.Mappings{
					Get: &v1alpha1.Mapping{
						Method: http.MethodGet,
						URL:    "http://127.0.0.1:8081/api/v1/pets",
					},
				},
				Payload: v1alpha1.Payload{
					BaseUrl: "http://127.0.0.1:8081/api/v1/pets",
					Body:    "{ \"id\": 1112, \"name\": \"fake-simple-name-2\", \"color\": \"simple-color-2\", \"price\": 123219, \"state\": \"foo-state-2\" }",
					BodyObject: runtime.RawExtension{
						Raw: []byte(`{"foo":"bar", "zoom":"xxxx"}`),
					},
				},
				Headers: map[string][]string{
					"Authorization": []string{"Basic BASE64_ENCODED_USER_CREDENTIALS"},
				},
				WaitTimeout:           nil,
				InsecureSkipTLSVerify: false,
			},
		},
		Status: v1alpha1.RequestStatus{},
	}

	path := "spec.forProvider.headers.Authorization[0]"
	value := "23213123"
	err := v1alpha1.PatchFieldValueToObject(path, value, to)
	require.NoError(t, err)

	path = "spec.forProvider.payload.body.foo"
	value = "2423423423"
	err = v1alpha1.PatchFieldValueToObject(path, value, to)
	require.NoError(t, err)

}

func TestMergeBodies(t *testing.T) {
	p1 := map[string]interface{}{
		"A": 3, "B": "note", "C": true,
	}
	p2 := map[string]interface{}{
		"B": "", "C": false,
	}
	if err := mergo.Merge(&p1, p2, mergo.WithOverwriteWithEmptyValue); err != nil {
		t.Error(err)
	}
}

/**
runtime.RawExtension{
			Object: &Foo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "sample.org/v2",
					Kind:       "Sample",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      strconv.Itoa(i),
					Namespace: "default",
					Labels: map[string]string{
						"label-key-1": "label-value-1",
					},
					Annotations: map[string]string{
						"annotations-key-1": "annotations-value-1",
					},
				},
				Spec: FooSpec{
					Flied: i,
				},
			},
*/
