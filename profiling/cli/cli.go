package cli

// init is called before main
func init() {
	// Initialize the config and panic on failure
	if err := configInit(); err != nil {
		panic(err.Error())
	}
}

// Execute is the entry point for the cli
// called from main
func Execute() error {
	defer stopProfiling()
	return rootCmd.Execute()
}
