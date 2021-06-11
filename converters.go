package go_utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// StructToMap converts an object (struct) to a map.
//
// WARNING: int inputs are converted to floats in the output map. This is an
// unintended consequence of converting through JSON.
//
// In future, this should be deprecated.
func StructToMap(item interface{}) (map[string]interface{}, error) {
	bs, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal to JSON: %v", err)
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal from JSON to map: %v", err)
	}
	return res, nil
}

// MapInterfaceToMapString converts a map with interface{} values to one with string
// values.
//
// It is used to convert a GraphQL (gqlgen) input Map to a map of strings for APIs
// that need map[string]string.
func MapInterfaceToMapString(in map[string]interface{}) (map[string]string, error) {
	out := map[string]string{}
	for k, v := range in {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%v (%T) is not a string", v, v)
		}
		out[k] = s
	}
	return out, nil
}

// ConvertStringMap converts a map[string]string to a map[string]interface{}.
//
// This is done mostly in order to conform to the gqlgen Graphql Map scalar.
func ConvertStringMap(inp map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	if inp == nil {
		return out
	}
	for k, v := range inp {
		val := interface{}(v)
		out[k] = val
	}
	return out
}

// ConvertInterfaceMap converts a map[string]interface{} to a map[string]string.
//
// Any conversion errors are written out to the output map instead of being
// returned as error values.
//
// New code is discouraged from using this function.
func ConvertInterfaceMap(inp map[string]interface{}) map[string]string {
	out := make(map[string]string)
	if inp == nil {
		return out
	}
	for k, v := range inp {
		val, ok := v.(string)
		if !ok {
			val = fmt.Sprintf("invalid string value: %#v", v)
			if IsDebug() {
				log.Printf(
					"non string value in map[string]interface{} that is to be converted into map[string]string: %#v", v)
			}
		}
		out[k] = val
	}
	return out
}

// ChunkStringSlice chunks the supplied slice of strings into chunks of the
// indicated length. The last chunk may be smaller than the indicated length.
func ChunkStringSlice(items []string, chunkSize int) [][]string {
	chunks := [][]string{}
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
