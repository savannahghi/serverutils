package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// pushSchemaCmd is the command used to push/update the graphql schema for a service
var pushSchemaCmd = &cobra.Command{
	Use:     "push-schema",
	Short:   "It is used to create or update the registered schema for a service",
	Long:    ``, //long description
	Example: "",
	Hidden:  false,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("name")
		serviceURL, _ := cmd.Flags().GetString("url")
		version, _ := cmd.Flags().GetString("version")
		schemaDir, _ := cmd.Flags().GetString("dir")
		schemaExtensionName, _ := cmd.Flags().GetString("file-extension")
		schemaPushURL, _ := cmd.Flags().GetString("registry-url")

		service := Service{
			Name:    serviceName,
			URL:     serviceURL,
			Version: version,
		}

		status, err := publishSchema(service, schemaDir, schemaExtensionName, schemaPushURL)
		if err != nil {
			fmt.Printf("error pushing schema: %v \n", err)
			os.Exit(1)
		}

		if !status.Valid {
			fmt.Printf("Schema for %v version:%v has not been published 🤡\nMessage: %v\n", service.Name, version, status.Message)
			os.Exit(1)
		}

		fmt.Printf("Schema for service: %v version: %v successfully published 🥳🎉\n", service, version)
		os.Exit(0)
	},
}

func init() {
	// Local flags
	pushSchemaCmd.Flags().String("dir", ".", "The directory containing the graphql schema files")
	pushSchemaCmd.Flags().String("file-extension", "graphql", "The extension for graphql schema files in directory")
	pushSchemaCmd.Flags().String("registry-url", "", "The schema registry url")
}

func publishSchema(service Service, dir, extension, baseURL string) (*SchemaStatus, error) {
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

	u.Path = "schema/push"
	publishSchemaURL := u.String()

	return schemaRegistryRequest(payload, publishSchemaURL)
}
