package cli

import (
	"fmt"
	"github.com/benchkram/cli_utils/extended/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var rootCmd = &cobra.Command{
	Use:   "cli_example",
	Short: "cli to start example server & client",
	Long:  "cli to start example server & client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readGlobalConfig()
		onStopProfiling = profilingInit(
			GlobalConfig.CPUProfile,
			GlobalConfig.MEMProfile,
		)
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
	// Print global config
	GlobalConfig.Print()

	// Create new app instance
	newApp := app.NewApplication(app.VersionInfo{
		Version: Version,
		Commit:  CommitHash,
	})

	return newApp.Start()
}
