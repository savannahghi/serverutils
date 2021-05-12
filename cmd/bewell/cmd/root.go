package cmd

import (
	"github.com/spf13/cobra"
	service "gitlab.slade360emr.com/go/base/cmd/bewell/cmd/service"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bewell",
	Short: "bewell CLI",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initCLI)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.AddCommand(service.ServiceCmd)

}

// initCLI called when the bewell command is initialized
func initCLI() {}
