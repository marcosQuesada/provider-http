package request

import (
	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
)

func getMappingByAction(requestParams *v1alpha1.RequestParameters, action Action) (*v1alpha1.Mapping, bool) {
	for _, mapping := range requestParams.Mappings {
		if mapping.Action == string(action) {
			return &mapping, true
		}
	}
	return nil, false
}
