package cli

// CLi Parameters
var (
	CPUProfile     bool
	MEMProfile     bool
	CPUProfileFile string
	MEMProfileFile string
)

// init is called before main
func init() {
	// Profiling cli flags
	rootCmd.PersistentFlags().BoolVar(&CPUProfile, "cpu-profile", false, "write cpu profile to file")
	rootCmd.PersistentFlags().BoolVar(&MEMProfile, "mem-profile", false, "write memory profile to file")

	rootCmd.PersistentFlags().StringVar(&CPUProfileFile, "cpu-profile-file", "cpu.prof", "write cpu profile to file")
	rootCmd.PersistentFlags().StringVar(&MEMProfileFile, "mem-profile-file", "mem.prof", "write memory profile to file")
}

// Execute is the entry point for the cli
// called from main
func Execute() error {
	defer stopProfiling()
	return rootCmd.Execute()
}
