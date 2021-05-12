package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// validateSchemaCmd represents the schema command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "creates a bewell microservice with boilerplate",
	Long:    ``, //long description
	Example: "",
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service create called")
	},
}

func init() {

}
