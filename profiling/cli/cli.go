package cli

// CLi Parameters
var (
	cpuProfile     bool
	memProfile     bool
	cpuProfileFile string
	memProfileFile string
)

// init is called before main
func init() {
	// Profiling cli flags
	rootCmd.PersistentFlags().BoolVar(&cpuProfile, "cpu-profile", false, "write cpu profile to file")
	rootCmd.PersistentFlags().BoolVar(&memProfile, "mem-profile", false, "write memory profile to file")

	rootCmd.PersistentFlags().StringVar(&cpuProfileFile, "cpu-profile-file", "cpu.prof", "write cpu profile to file")
	rootCmd.PersistentFlags().StringVar(&memProfileFile, "mem-profile-file", "mem.prof", "write memory profile to file")
}

// Execute is the entry point for the cli
// called from main
func Execute() error {
	defer stopProfiling()
	return rootCmd.Execute()
}
