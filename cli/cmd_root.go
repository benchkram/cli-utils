package cli

import (
	"cli_example/app"
	"cli_example/server"
	"fmt"
	"github.com/benchkram/errz"

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
		logInit(int8(GlobalConfig.Verbosity), GlobalConfig.Pretty)
		onStopProfiling = profilingInit(
			GlobalConfig.CPUProfile,
			GlobalConfig.MEMProfile,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runRootJob()
		// On the most outside function we only log error
		errz.Log(err)
	},
}

// runRootJob is the actual job that is executed by the root command
func runRootJob() (err error) {
	log.Info("root job started", "version", Version, "commit", CommitHash)

	// Create new app instance
	newApp := app.NewApplication(app.VersionInfo{
		Version: Version,
		Commit:  CommitHash,
	})

	// Create new server instance
	newServer := server.NewServer(
		fmt.Sprintf("%s:%s", GlobalConfig.Host, GlobalConfig.Port),
		newApp,
	)

	// Start server - blocking
	err = newServer.Start()

	return err
}
