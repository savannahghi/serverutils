package base

import (
	"encoding/json"
	"fmt"
)

// StructToMap converts an object (struct) to a map.
//
// WARNING: int inputs are converted to floats in the output map. This is an
// unintended consequence of converting through JSON.
//
// Deprecated: the unintended co-ercing of ints to floats means that this
// function should be used with caution. If possible, use a different utility.
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

// ChunkStringSlice chunks the supplied slice of strings into chunks of the
// indicated length. The last chunk may be smaller than the indicated length.
func ChunkStringSlice(items []string, chunkSize int) [][]string {
	chunks := [][]string{}
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
