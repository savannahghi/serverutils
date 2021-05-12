package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// validateSchemaCmd is the command used to validate the graphql schema for a service
var validateSchemaCmd = &cobra.Command{
	Use:     "validate-schema",
	Short:   "It is used to validate a service's schema with gateway",
	Long:    ``, //long description
	Example: "",
	Hidden:  false,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("name")
		serviceURL, _ := cmd.Flags().GetString("url")
		version, _ := cmd.Flags().GetString("version")
		schemaDir, _ := cmd.Flags().GetString("dir")
		// schemaValidationURL := "https://postman-echo.com/post"
		schemaValidationURL, _ := cmd.Flags().GetString("registry-url")
		schemaExtensionName, _ := cmd.Flags().GetString("file-extension")

		service := Service{
			Name:    serviceName,
			URL:     serviceURL,
			Version: version,
		}

		status, err := validateSchema(service, schemaDir, schemaExtensionName, schemaValidationURL)
		if err != nil {
			fmt.Printf("error validating schema: %v \n", err)
			os.Exit(1)
		}

		if !status.Valid {
			fmt.Printf("Schema for %v version:%v is Invalid 🤡\nMessage: %v\n", service.Name, version, status.Message)
			os.Exit(1)
		}

		fmt.Printf("Schema for service: %v version: %v is Valid 🥳🎉\n", service, version)
		os.Exit(0)
	},
}

func init() {
	// Local flags
	validateSchemaCmd.Flags().String("dir", ".", "The directory containing the graphql schema files")
	validateSchemaCmd.Flags().String("file-extension", "graphql", "The extension for graphql schema files in directory")
	validateSchemaCmd.Flags().String("registry-url", "", "The schema registry url")
}

// GraphqlSchemaPayload is the payload made when making validation/push requests
type GraphqlSchemaPayload struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	TypeDefs string `json:"type_defs"`
}

// SchemaStatus holds Status and message(if any) for a schema request
type SchemaStatus struct {
	Valid   bool
	Message string
}

func validateSchema(service Service, dir, extension, baseURL string) (*SchemaStatus, error) {
	err := service.ValidateFields()
	if err != nil {
		return nil, err
	}

	schema, err := readSchemaFilesInDirectory(dir, extension)
	if err != nil {
		return nil, err
	}

	payload := GraphqlSchemaPayload{
		Name:     service.Name,
		URL:      service.URL,
		Version:  service.Version,
		TypeDefs: schema,
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u.Path = "schema/validate"
	validationURL := u.String()

	return schemaRegistryRequest(payload, validationURL)
}
