package cli

import (
	"fmt"
	"github.com/benchkram/cli-utils/extended/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var rootCmd = &cobra.Command{
	Use:   "our_app",
	Short: "cli to start example server & client",
	Long:  "cli to start example server & client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		onStopProfiling = profilingInit()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runRootJob()
		// On the most outside function we only log error
		if err != nil {
			fmt.Println(err)
		}
	},
}

// runRootJob is the actual job that is executed by the root command
func runRootJob() (err error) {
	// Print the config
	fmt.Printf(
		"Config: {\n\tCPUProfile: %t\n\tCPUProfileFile: %s\n\tMEMProfile: %t\n\tMEMProfileFile: %s\n}\n",
		cpuProfile, cpuProfileFile, memProfile, memProfileFile)

	// Create new app instance
	newApp := app.NewApplication()

	return newApp.Start()
}
