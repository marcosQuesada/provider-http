package json

import (
	"bytes"
	encoder "encoding/json"
)

func Contains(container, containee map[string]interface{}) bool {
	for key, value := range containee {
		if containerValue, exists := container[key]; !exists || !deepEqual(value, containerValue) {
			return false
		}
	}
	return true
}

func IsJSONString(jsonStr string) bool {
	var js map[string]interface{}
	return encoder.Unmarshal([]byte(jsonStr), &js) == nil
}

func JsonStringToMap(jsonStr string) map[string]interface{} {
	var jsonData map[string]interface{}
	_ = encoder.Unmarshal([]byte(jsonStr), &jsonData)
	return jsonData
}

// Converts JSON strings within a map to maps for JSON data processing.
func ConvertJSONStringsToMaps(merged *map[string]interface{}) {
	for key, value := range *merged {

		switch valueToHandle := value.(type) {
		case string:
			if IsJSONString(valueToHandle) {
				mappedJSON := JsonStringToMap(valueToHandle)
				(*merged)[key] = mappedJSON
			}
		case map[string]interface{}:
			ConvertJSONStringsToMaps(&valueToHandle)
		case []interface{}:
			structToMap, _ := (StructToMap(valueToHandle))
			ConvertJSONStringsToMaps(&structToMap)
		}
	}
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := encoder.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = encoder.Unmarshal(data, &newMap) // Convert to a map
	return
}

func deepEqual(a, b interface{}) bool {
	aBytes, err := encoder.Marshal(a)
	if err != nil {
		return false
	}

	bBytes, err := encoder.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(aBytes, bBytes)
}
