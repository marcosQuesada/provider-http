package requestgen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"

	"github.com/crossplane-contrib/provider-http/apis/request/v1alpha1"
	"github.com/crossplane-contrib/provider-http/internal/controller/request/requestprocessing"
	json_util "github.com/crossplane-contrib/provider-http/internal/json"
	"github.com/crossplane-contrib/provider-http/internal/utils"

	"golang.org/x/exp/maps"
)

type RequestDetails struct {
	Url     string
	Body    string // @TODO: What about this ¿?
	Headers map[string][]string
}

// GenerateRequestDetails generates request details.
func GenerateRequestDetails(methodMapping v1alpha1.Mapping, forProvider v1alpha1.RequestParameters, response v1alpha1.Response) (RequestDetails, error, bool) {
	jqObject := generateRequestObject(forProvider, response)
	url, err := generateURL(methodMapping.URL, jqObject)
	if err != nil {
		return RequestDetails{}, err, false
	}

	if !utils.IsUrlValid(url) {
		return RequestDetails{}, errors.Errorf(utils.ErrInvalidURL, url), false
	}

	body, err := generateBody(string(methodMapping.Body.Raw), jqObject)
	if err != nil {
		return RequestDetails{}, err, false
	}

	// @TODO: Clean
	// bodyObject := forProvider.Payload.Body
	// if len(bodyObject.Raw) > 0 && len(body) > 0 {
	//	body, err = mergeBody([]byte(body), bodyObject.Raw)
	//	if err != nil {
	//		return RequestDetails{}, fmt.Errorf("unable to merge bodies, %w", err), false
	//	}
	// }

	headers, err := generateHeaders(coalesceHeaders(methodMapping.Headers, forProvider.Headers), jqObject)
	if err != nil {
		return RequestDetails{}, err, false
	}

	return RequestDetails{Body: body, Url: url, Headers: headers}, nil, true
}

// nolint:unused
func mergeBody(src, patch []byte) (string, error) {
	current := map[string]interface{}{}
	if err := json.Unmarshal(src, &current); err != nil {
		return "", fmt.Errorf("unable to unmarshall, %w", err)
	}
	bom := map[string]interface{}{}
	if err := json.Unmarshal(patch, &bom); err != nil {
		return "", fmt.Errorf("unable to unmarshall, %w", err)
	}

	if err := mergo.Merge(&current, bom, mergo.WithOverwriteWithEmptyValue); err != nil {
		return "", fmt.Errorf("unable to merge, %w", err)
	}
	raw, err := json.Marshal(bom)
	if err != nil {
		return "", fmt.Errorf("unable to marshall, %w", err)
	}
	return string(raw), nil
}

// generateRequestObject creates a JSON-compatible map from the specified Request's ForProvider and Response fields.
// It merges the two maps, converts JSON strings to nested maps, and returns the resulting map.
func generateRequestObject(forProvider v1alpha1.RequestParameters, response v1alpha1.Response) map[string]interface{} {
	baseMap, _ := json_util.StructToMap(forProvider)
	statusMap, _ := json_util.StructToMap(map[string]interface{}{
		"response": response,
	})

	maps.Copy(baseMap, statusMap)
	json_util.ConvertJSONStringsToMaps(&baseMap)

	return baseMap
}

func IsRequestValid(requestDetails RequestDetails) bool {
	return (!strings.Contains(fmt.Sprint(requestDetails), "null")) && (requestDetails.Url != "")
}

// coalesceHeaders returns the non-nil headers, or the default headers if both are nil.
func coalesceHeaders(mappingHeaders, defaultHeaders map[string][]string) map[string][]string {
	if mappingHeaders != nil {
		return mappingHeaders
	}
	return defaultHeaders
}

// generateURL applies a JQ filter to generate a URL.
func generateURL(urlJQFilter string, jqObject map[string]interface{}) (string, error) {
	getURL, err := requestprocessing.ApplyJQOnStr(urlJQFilter, jqObject)
	if err != nil {
		return "", err
	}

	return getURL, nil
}

// generateBody applies a mapping body to generate the request body.
func generateBody(mappingBody string, jqObject map[string]interface{}) (string, error) {
	if mappingBody == "" {
		return "", nil
	}

	jqQuery, err := FromJSONMapToJQuery(requestprocessing.ConvertStringToJQQuery(mappingBody))
	if err != nil {
		return "", fmt.Errorf("FromJSONMapToJQuery error %w", err)
	}

	body, err := requestprocessing.ApplyJQOnStr(jqQuery, jqObject)
	if err != nil {
		return "", fmt.Errorf("ApplyJQOnStr error %w", err)
	}

	return body, nil
}

// generateHeaders applies JQ queries to generate headers.
func generateHeaders(headers map[string][]string, jqObject map[string]interface{}) (map[string][]string, error) {
	generatedHeaders, err := requestprocessing.ApplyJQOnMapStrings(headers, jqObject)
	if err != nil {
		return nil, err
	}

	return generatedHeaders, nil
}

// FromJSONMapToJQuery convert json map[string]any values to JQ query
func FromJSONMapToJQuery(in string) (string, error) {
	tmp := map[string]any{}
	if err := json.Unmarshal([]byte(in), &tmp); err != nil {
		return "", fmt.Errorf("unable to UnMarshall to map, error %w", err)
	}

	var sb strings.Builder
	sb.WriteString("{")

	i := 0
	for k, v := range tmp {
		i++
		sb.WriteString("\"" + k + "\"" + ":")

		parse(v, &sb)
		if i != len(tmp) {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")

	return sb.String(), nil
}

func parse(val any, res *strings.Builder) {
	switch v := val.(type) {
	case int, int64, float64:
		res.WriteString(fmt.Sprintf("%d, ", v))
	case string:
		if strings.HasPrefix(v, ".") {
			res.WriteString(v)
			return
		}
		res.WriteString("\"" + v + "\"")
	case []any:
		res.WriteString("[")
		for i, a := range v {
			parse(a, res)
			if i+1 != len(v) {
				res.WriteString(", ")
			}
		}
		res.WriteString("]")
	case map[string]any:
		res.WriteString("{")
		i := 0
		for kv, vv := range v {
			i++
			res.WriteString("\"" + kv + "\"" + ": ")
			parse(vv, res)
			if i != len(v) {
				res.WriteString(", ")
			}
		}
		res.WriteString("}")
	}
}
