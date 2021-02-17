package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ReadRequestToTarget handles the process of sending requests to a Slade 360 API though the supplied client and reading back the response as bytes.
// Network and IO errors are reported back.
// `target` **MUST** be a pointer to a struct - the struct to which the JSON response should be decoded
func ReadRequestToTarget(apiClient Client, method string, path string, query string, content []byte, target interface{}) error {
	url := ComposeAPIURL(apiClient, path, query)
	return PerformRequest(apiClient, method, url, content, target)
}

// ReadAuthServerRequestToTarget handles the process of sending requests to a Slade 360 auth server URL through th supplied client and reading back the response as bytes.
// Network and IO errors are reported back.
// `target` **MUST** be a pointer to a struct - the struct to which the JSON response should be decoded
func ReadAuthServerRequestToTarget(client Client, method string, url string, s string, content []byte, target interface{}) error {
	return PerformRequest(client, method, url, content, target)
}

// PerformRequest implements the final common pathway for requests to EDI APIs (regardless of the server the request is going to)
func PerformRequest(httpClient Client, method string, url string, content []byte, target interface{}) error {
	var bodyContent io.Reader
	if content != nil {
		bodyContent = bytes.NewReader(content)
	}
	resp, err := httpClient.MakeRequest(method, url, bodyContent)
	defer CloseRespBody(resp)
	if err != nil {
		return fmt.Errorf("PerformRequest(method, url, reader): %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("json.NewDecoder(resp.Body).Decode(&target): %w", err)
	}
	return nil
}

// ReadWriteRequestToTarget handles the process of sending requests to a Slade 360 API though the supplied client, reading back the response as bytes
// and write and interface of the data.
// Network and IO errors are reported back.
// `target` **MUST** be a pointer to a struct - the struct to which the JSON response should be decoded
func ReadWriteRequestToTarget(apiClient Client, method string, path string, query string, content []byte, target interface{}) (interface{}, error) {
	url := ComposeAPIURL(apiClient, path, query)
	return PerformRequestWithTarget(apiClient, method, url, content, target)
}

// PerformRequestWithTarget implements the final common pathway for requests to slade 360 APIs (regardless of the server the request is going to)
// and returns the target with data
func PerformRequestWithTarget(httpClient Client, method string, url string, content []byte, target interface{}) (interface{}, error) {
	var bodyContent io.Reader
	if content != nil {
		bodyContent = bytes.NewReader(content)
	}
	resp, err := httpClient.MakeRequest(method, url, bodyContent)
	defer CloseRespBody(resp)
	if err != nil {
		return nil, fmt.Errorf("PerformRequest(method, url, reader): %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder(resp.Body).Decode(&target): %w", err)
	}

	return target, nil
}
