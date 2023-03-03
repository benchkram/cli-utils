package cli

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/leebenson/conform"
	"github.com/logrusorgru/aurora"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
)

// handle global configuration through a config file, environment vars  cli parameters.

// Config the global config object
var GlobalConfig *config

func readGlobalConfig() {
	// Priority of configuration options
	// 1: CLI Parameters
	// 2: environment
	// 2: config.yaml
	// 3: defaults
	config, err := readConfig()
	if err != nil {
		panic(err.Error())
	}

	// Set config object for main package
	GlobalConfig = config
}

var defaultConfig = &config{
	Host: "",
	Port: "8080",

	Verbosity: 1,
	Pretty:    false,

	CPUProfile: false,
	MEMProfile: false,
}

// configInit must be called from the packages' init() func
func configInit() {
	cliFlags()
	bind()
	env()
}

// Create private data struct to hold config options.
// `mapstructure` => viper tags
// `struct` => fatih structs tag
type config struct {
	// Server config
	Host string `mapstructure:"host" structs:"host"`
	Port string `mapstructure:"port" structs:"port"`

	// Log
	Verbosity int  `mapstructure:"verbosity" structs:"verbosity"`
	Pretty    bool `mapstructure:"pretty" structs:"pretty"`

	// Profiling
	CPUProfile bool `mapstructure:"cpuprofile" structs:"cpuprofile"`
	MEMProfile bool `mapstructure:"memprofile" structs:"memprofile"`
}

// cliFlags defines cli parameters for all config options
func cliFlags() {
	// Keep cli parameters in sync with the config struct

	// Server
	rootCmd.PersistentFlags().String("host", defaultConfig.Host, "hostname to listen to")
	rootCmd.PersistentFlags().String("port", defaultConfig.Port, "port to listen to")

	// Log
	rootCmd.PersistentFlags().Int("verbosity", defaultConfig.Verbosity, "verbosity level from quiet to verbose (0-10)")
	rootCmd.PersistentFlags().Bool("pretty", defaultConfig.Pretty, "log pretty instead of json")

	// Profiling
	rootCmd.PersistentFlags().Bool("cpuprofile", defaultConfig.CPUProfile, "write cpu profile to file")
	rootCmd.PersistentFlags().Bool("memprofile", defaultConfig.MEMProfile, "write memory profile to file")
}

// bind will assign the environment variables to the cli parameters
func bind() {
	// Server
	_ = viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	// Log
	_ = viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))
	_ = viper.BindPFlag("pretty", rootCmd.PersistentFlags().Lookup("pretty"))

	// Profiling
	_ = viper.BindPFlag("cpuprofile", rootCmd.PersistentFlags().Lookup("cpuprofile"))
	_ = viper.BindPFlag("memprofile", rootCmd.PersistentFlags().Lookup("memprofile"))
}

// env create environment vars for all config options
func env() {
	// Typically we use capital letters for env vars

	// Server
	_ = viper.BindEnv("host", "HOST")
	_ = viper.BindEnv("port", "PORT")

	// Log
	_ = viper.BindEnv("verbosity", "VERBOSITY")
	_ = viper.BindEnv("pretty", "PRETTY")

	// Profiling
	_ = viper.BindEnv("cpuprofile", "SUBSVC_CPU_PROFILE")
	_ = viper.BindEnv("memprofile", "SUBSVC_MEM_PROFILE")
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
	cefaultsAsMap := structs.Map(defaultConfig)

	// Set defaults
	for key, value := range cefaultsAsMap {
		viper.SetDefault(key, value)
	}

	// Read config from file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		fmt.Printf("%s\n", aurora.Yellow("Could not find a config file"))
	default:
		return nil, fmt.Errorf("config file invalid: %s \n", err)
	}

	// Unmarshal config into struct
	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
