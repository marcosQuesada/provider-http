/*
Copyright 2024 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package request

import (
	"context"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	apisv1alpha1 "github.com/crossplane-contrib/provider-http/apis/v1alpha1"
	httpClient "github.com/crossplane-contrib/provider-http/internal/clients/http"
	"github.com/crossplane-contrib/provider-http/internal/controller/request/requestgen"
	"github.com/crossplane-contrib/provider-http/internal/controller/request/statushandler"
	"github.com/crossplane-contrib/provider-http/internal/utils"
)

const (
	errNotRequest                   = "managed resource is not a Request custom resource"
	errTrackPCUsage                 = "cannot track ProviderConfig usage"
	errNewHttpClient                = "cannot create new Http client"
	errProviderNotRetrieved         = "provider could not be retrieved"
	errFailedToSendHttpRequest      = "something went wrong"
	errFailedToCheckIfUpToDate      = "failed to check if request is up to date"
	errFailedToUpdateStatusFailures = "failed to reset status failures counter"
	errFailedUpdateStatusConditions = "failed updating status conditions"
	errMappingNotFound              = "%s mapping doesn't exist in request, skipping operation"
	errGetReferencedResource        = "cannot get referenced resource"
	errPatchFromReferencedResource  = "cannot patch from referenced resource"
	errResolveResourceReferences    = "cannot resolve resource references"
)

var (
	ErrMappingNotFound = errors.New(errMappingNotFound)
)

// Setup adds a controller that reconciles Request managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, timeout time.Duration) error {
	name := managed.ControllerName(v1alpha1.RequestGroupKind)
	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.RequestGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			logger:          o.Logger,
			kube:            mgr.GetClient(),
			usage:           resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newHttpClientFn: httpClient.NewClient,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithTimeout(timeout),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Request{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	logger          logging.Logger
	kube            client.Client
	usage           resource.Tracker
	newHttpClientFn func(log logging.Logger, timeout time.Duration) (httpClient.Client, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Request)
	if !ok {
		return nil, errors.New(errNotRequest)
	}

	l := c.logger.WithValues("request", cr.Name)

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	n := types.NamespacedName{Name: cr.GetProviderConfigReference().Name}
	if err := c.kube.Get(ctx, n, pc); err != nil {
		return nil, errors.Wrap(err, errProviderNotRetrieved)
	}

	h, err := c.newHttpClientFn(l, utils.WaitTimeout(cr.Spec.ForProvider.WaitTimeout))
	if err != nil {
		return nil, errors.Wrap(err, errNewHttpClient)
	}

	return &external{
		localKube: c.kube,
		logger:    l,
		http:      h,
	}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	localKube client.Client
	logger    logging.Logger
	http      httpClient.Client
}

// nolint:gocyclo
func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Request)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRequest)
	}

	if err := c.resolveReferencies(ctx, cr); err != nil {

		return managed.ExternalObservation{}, errors.Wrap(err, errResolveResourceReferences)
	}

	observeRequestDetails, err := c.isUpToDate(ctx, cr)
	if err != nil && err.Error() == errObjectNotFound {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	if err != nil && errors.Is(err, ErrMappingNotFound) {
		cr.Status.SetConditions(xpv1.Available())
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errFailedToCheckIfUpToDate)
	}

	// Get the latest version of the resource before updating
	if err := c.localKube.Get(ctx, types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to get the latest version of the resource")
	}

	statusHandler, err := statushandler.NewStatusHandler(ctx, cr.Status.RequestDetails.Action, cr, observeRequestDetails.Details, observeRequestDetails.ResponseError, c.localKube, c.logger)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	synced := observeRequestDetails.Synced
	if synced {
		statusHandler.ResetFailures()
	}

	cr.Status.SetConditions(xpv1.Available())
	err = statusHandler.SetRequestStatus()
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, " failed updating status")
	}

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  synced,
		ConnectionDetails: nil,
	}, nil
}

func (c *external) deployAction(ctx context.Context, cr *v1alpha1.Request, action Action) error {
	mapping, ok := getMappingByAction(&cr.Spec.ForProvider, action)
	if !ok {
		c.logger.Info(errMappingNotFound, action)
		return nil
	}

	requestDetails, err := generateValidRequestDetails(cr, mapping)
	if err != nil {
		return err
	}

	details, err := c.http.SendRequest(ctx, mapping.Method, requestDetails.Url, requestDetails.Body, requestDetails.Headers, cr.Spec.ForProvider.InsecureSkipTLSVerify)

	statusHandler, err := statushandler.NewStatusHandler(ctx, string(action), cr, details, err, c.localKube, c.logger)
	if err != nil {
		return err
	}

	return statusHandler.SetRequestStatus()
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Request)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRequest)
	}

	return managed.ExternalCreation{}, errors.Wrap(c.deployAction(ctx, cr, CREATE), errFailedToSendHttpRequest)
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Request)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRequest)
	}

	return managed.ExternalUpdate{}, errors.Wrap(c.deployAction(ctx, cr, UPDATE), errFailedToSendHttpRequest)
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Request)
	if !ok {
		return errors.New(errNotRequest)
	}

	return errors.Wrap(c.deployAction(ctx, cr, DELETE), errFailedToSendHttpRequest)
}

// generateValidRequestDetails generates valid request details based on the given Request resource and Mapping configuration.
// It first attempts to generate request details using the HTTP response stored in the Request's status. If the generated
// details are valid, the function returns them. If not, it falls back to using the cached response in the Request's status
// and attempts to generate request details again. The function returns the generated request details or an error if the
// generation process fails.
func generateValidRequestDetails(cr *v1alpha1.Request, mapping *v1alpha1.Mapping) (requestgen.RequestDetails, error) {
	requestDetails, _, ok := requestgen.GenerateRequestDetails(*mapping, cr.Spec.ForProvider, cr.Status.Response)
	if requestgen.IsRequestValid(requestDetails) && ok {
		return requestDetails, nil
	}

	requestDetails, err, _ := requestgen.GenerateRequestDetails(*mapping, cr.Spec.ForProvider, cr.Status.Cache.Response)
	if err != nil {
		return requestgen.RequestDetails{}, err
	}

	return requestDetails, nil
}

func getReferenceInfo(ref v1alpha1.Reference) (string, string, string, string) {
	var apiVersion, kind, namespace, name string

	if ref.PatchesFrom != nil {
		// Reference information defined in PatchesFrom
		apiVersion = ref.PatchesFrom.APIVersion
		kind = ref.PatchesFrom.Kind
		namespace = ref.PatchesFrom.Namespace
		name = ref.PatchesFrom.Name
	} else if ref.DependsOn != nil {
		// Reference information defined in DependsOn
		apiVersion = ref.DependsOn.APIVersion
		kind = ref.DependsOn.Kind
		namespace = ref.DependsOn.Namespace
		name = ref.DependsOn.Name
	}

	return apiVersion, kind, namespace, name
}

// resolveReferencies resolves references for the current Object. If it fails to
// resolve some reference, e.g.: due to reference not ready, it will then return
// error and requeue to wait for resolving it next time.
func (c *external) resolveReferencies(ctx context.Context, obj *v1alpha1.Request) error {
	c.logger.Debug("Resolving referencies.")

	// Loop through references to resolve each referenced resource
	for _, ref := range obj.Spec.References {
		if ref.DependsOn == nil && ref.PatchesFrom == nil {
			continue
		}

		refAPIVersion, refKind, refNamespace, refName := getReferenceInfo(ref)
		res := &unstructured.Unstructured{}
		res.SetAPIVersion(refAPIVersion)
		res.SetKind(refKind)
		// Try to get referenced resource
		err := c.localKube.Get(ctx, client.ObjectKey{
			Namespace: refNamespace,
			Name:      refName,
		}, res)

		if err != nil {
			return errors.Wrap(err, errGetReferencedResource)
		}

		// @TODO: Assert Condition

		// Patch fields if any
		if ref.PatchesFrom != nil && ref.PatchesFrom.FieldPath != nil {
			if err := ref.ApplyFromFieldPathPatch(res, obj); err != nil {
				return errors.Wrap(err, errPatchFromReferencedResource)
			}
		}
	}

	return nil
}
