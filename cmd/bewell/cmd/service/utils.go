package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gitlab.slade360emr.com/go/base"
	"moul.io/http2curl"
)

func readSchemaFile(schemaFile string) (string, error) {
	data, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		return "", fmt.Errorf("cannot read schema file %v with error: %v", schemaFile, err)
	}
	return string(data), nil
}

func readSchemaFilesInDirectory(dir, extension string) (string, error) {
	schemaExtension := "." + extension

	var schema string
	var files []string

	var readSchemaFileFunc filepath.WalkFunc = func(path string, info fs.FileInfo, err error) error {
		name := info.Name()
		if !info.IsDir() && strings.Contains(name, schemaExtension) {
			files = append(files, name)
			definitions, _ := readSchemaFile(path)

			schema += definitions + "\n"
		}

		return nil
	}

	err := filepath.Walk(dir, readSchemaFileFunc)
	if err != nil {
		return "", fmt.Errorf("cannot read schema directory %v with error: %v", dir, err)
	}

	if schema == "" {
		return "", fmt.Errorf("no schema files found in directory: %v", dir)
	}

	fmt.Printf("\nschema files: %v\n\n", files)

	return schema, nil
}

// Change based on where/how request will be done internally. Consult Macharia for full context
func makeRequest(url string, body interface{}) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize: %v", err)
	}

	payload := bytes.NewBuffer(encoded)

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request:%v", err)
	}

	// TODO: Remember Authentication stuff. Consult Macharia
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if base.IsDebug() {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)
	}

	return client.Do(req)

}

func schemaRegistryRequest(payload GraphqlSchemaPayload, url string) (*SchemaStatus, error) {
	resp, err := makeRequest(url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to make validation http request: %v", err)
	}
	defer resp.Body.Close()

	status := SchemaStatus{}
	response := RegistryResponse{}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %v", err)
	}

	status.Valid = response.Success
	status.Message = response.Message

	return &status, nil

}

//RegistryResponse ...
type RegistryResponse struct {
	Success bool              `json:"success,omitempty"`
	Message string            `json:"message,omitempty"`
	Details []ResponseDetails `json:"details,omitempty"`
}

// ResponseDetails ...
type ResponseDetails struct {
	Message string `json:"message,omitempty"`
}
