package request

import (
	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
)

func getMappingByAction(requestParams *v1alpha1.RequestParameters, action Action) (*v1alpha1.Mapping, bool) {
	switch action {
	case CREATE:
		return requestParams.Mappings.Create, true
	case GET:
		return requestParams.Mappings.Get, true
	case UPDATE:
		return requestParams.Mappings.Update, true
	case DELETE:
		return requestParams.Mappings.Delete, true
	}
	return nil, false
}
