package cli

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/leebenson/conform"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
)

// handle global configuration through a config file, environment vars  cli parameters.

// GlobalConfig the global config object
var GlobalConfig *config

func readGlobalConfig() {
	// Priority of configuration options
	// 1: CLI Parameters
	// 2: environment
	// 3: config.yaml
	// 4: defaults
	config, err := readConfig()
	if err != nil {
		panic(err.Error())
	}

	// Set config object for main package
	GlobalConfig = config
}

var defaultConfig = &config{
	MEMProfile: false,
}

// configInit must be called from the packages' init() func
func configInit() error {
	cliFlags()
	return bindFlagsAndEnv()
}

// Create private data struct to hold config options.
// `mapstructure` => viper tags
// `struct` => fatih structs tag
// `env` => environment variable name
type config struct {
	// Profiling
	CPUProfile bool `mapstructure:"cpu-profile" structs:"cpu-profile" env:"CPU_PROFILE"`
	MEMProfile bool `mapstructure:"mem-profile" structs:"mem-profile" env:"MEM_PROFILE"`

	CPUProfileFile string `mapstructure:"cpu-profile-file" structs:"cpu-profile-file" env:"CPU_PROFILE_FILE"`
	MEMProfileFile string `mapstructure:"mem-profile-file" structs:"mem-profile-file" env:"MEM_PROFILE_FILE"`
}

// cliFlags defines cli parameters for all config options
func cliFlags() {
	// Keep cli parameters in sync with the config struct

	// Profiling
	rootCmd.PersistentFlags().Bool("cpu-profile", defaultConfig.CPUProfile, "write cpu profile to file")
	rootCmd.PersistentFlags().Bool("mem-profile", defaultConfig.MEMProfile, "write memory profile to file")

	rootCmd.PersistentFlags().String("cpu-profile-file", defaultConfig.CPUProfileFile, "write cpu profile to file")
	rootCmd.PersistentFlags().String("mem-profile-file", defaultConfig.MEMProfileFile, "write memory profile to file")
}

// bindFlagsAndEnv will assign the environment variables to the cli parameters
func bindFlagsAndEnv() (err error) {
	for _, field := range structs.Fields(&config{}) {
		// Get the struct tag values
		key := field.Tag("structs")
		env := field.Tag("env")

		// Bind cobra flags to viper
		err = viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key))
		if err != nil {
			return err
		}
		err = viper.BindEnv(key, env)
		if err != nil {
			return err
		}
	}
	return nil
}

// Print the config object
// but remove sensitive data
func (c *config) Print() {
	cp := *c
	_ = conform.Strings(&cp)
	litter.Dump(cp)
}

// String the config object
// but remove sensitive data
func (c *config) String() string {
	cp := *c
	_ = conform.Strings(&cp)
	return litter.Sdump(cp)
}

// readConfig a helper to read default from a default config object.
func readConfig() (*config, error) {
	// Create a map of the default config
	defaultsAsMap := structs.Map(defaultConfig)

	// Set defaults
	for key, value := range defaultsAsMap {
		viper.SetDefault(key, value)
	}

	// Read config from file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal config into struct
	c := &config{}
	err := viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
