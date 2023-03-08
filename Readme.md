# Golang Cli Boilerplate Example

This is a simple example of a golang application boilerplate using a cli package.
The cli package is a helper package to add run parameters to your application as well as logging and profiling.
It can be used to create a simple cli application or a more complex server application.

What every application needs:
- parameters or configuration
- some kind of output (logging)
- versioning
- (optional) profiling

In the following I will explain how to use a cli package with a simple server application.

## Base application
For the base application we will start with a `main.go` file that only calls a cobra command.
The cobra command will call the root command of our application.

```go
package main

import (
	"fmt"
	"os"

	"cli_example/cli"
)

// Version and CommitHash set in compile time through ldflags
// Will be passed down to cli package
var (
	Version    = ""
	CommitHash = ""
)

func main() {
	// This main function is only responsible for calling the cli package
	// and handling errors returned by the cli package

	// Pass down version and commit hash
	cli.Version = Version
	cli.CommitHash = CommitHash

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
As you can see we also have a version and commit hash variable which we pass down to the cli package.
The Version and CommitHash variables are set in compile time through ldflags.
That way the `main.go` file as well as the basic cli package is easily reusable.

## First Cobra Command
Now we create our first cobra command in cli package. This is the root command of our application.
We can add subcommands to this command later.
Let's create a `cmd_root.go` file in the cli folder with the following content:

```go
package cli

import "github.com/spf13/cobra"

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var rootCmd = &cobra.Command{
    Use:   "cli_example",
    Short: "cli to start example server & client",
    Long:  "cli to start example server & client",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // TODO: later in this post
    },
    Run: func(cmd *cobra.Command, args []string) {
        // This is where your application starts
    },
}
```

## Config
A lot of cli applications only come with a few parameters, but as your application grows you will need to add more parameters.
The cli package allows you to add parameters to your application in a simple way.
It supports configuration files, environment variables and command line parameters.
This allows usage in cli as well as in CI/CD pipelines with either config files or environment variables.

We use cobra and viper to help us call an application with different modes or parameters.
# TODO: add links to cobra and viper?
Let's create a `config.go` file in the cli folder with a simple struct to hold our configuration:

```go
package cli

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/leebenson/conform"
	"github.com/logrusorgru/aurora"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
)

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
```

Whenever we want to add a new parameter to our application we just add it to the config struct.
The only parameters we need to the server itself are host and port. The rest is for logging and profiling.
But it's always useful to have the possibility to alter logging or profiling behaviour without rebuilding your application.

What is probably most helpful is to have a default config struct.
This way we can use it to set default values for our parameters and to read a config file.

```go
// defaultConfig holds default values for all config options
var defaultConfig = config{
    // Server
    Host: "localhost",
    Port: "8080",
    
    // Log
    Verbosity: 5,
    Pretty:    false,
    
    // Profiling
    CPUProfile: false,
    MEMProfile: false,
}
```

Now we want to make sure we can assign values to our config struct from parameters, environment variables and config files.
First the parameters:

```go
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
```

Then we use viper to bind environment variables to these parameters.

```go
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
```

And last we use viper to bind config files to these parameters.

```go
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
```

Before entering the application the config needs to be initialized. This can be done with the packages `init()` function.
# TODO: maybe move init to description of cli.go
```go
// configInit must be called from the packages' init() func
func configInit() {
    cliFlags()
    bind()
    env()
}

// init is called before main
func init() {
    // A custom sanitizer to redact sensitive data by defining a struct tag= named "redact".
    conform.AddSanitizer("redact", func(_ string) string { return "*****" })
    
    configInit()
}
```

When printing the config we need to make sure we don't print any sensitive information like passwords or tokens.
That's where the `redact` function comes in handy.

```go
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
```



## Logging
The cli package is a good starting point to instantiate a logger for your application.
It supports logging to stdout, stderr in different formats and levels.
This can be especially useful when you want to log in a more human-readable format when using in you local cli and a more machine-readable format when using in a CI/CD pipeline.

```shell
# output on you cli
$ 2021-03-01T12:00:00.000Z [INFO]  cli_example: Hello World, someValue: 123, someOtherValue: 456
```

```shell
# output on your CI/CD pipeline
$ 2021-03-01T12:00:00.000Z [INFO]  cli_example: {"level":"info","message":"Hello World", "someValue": "123", "someOtherValue": "456"}
```

## Build
I like to build my applications with a version number and a git commit hash. That way I can easily identify the version of my application.
To do so I add the following build flags to my build command:

```shell
go build -ldflags="-s -w -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH}"
```

### Initialize logging
To initialize logging we need to create the `loggingInit()` function.
Let's create a new file called `log.go` and add the loggingInit function:

```go
package cli

import (
    "os"
    
    "github.com/go-logr/logr"
    "github.com/go-logr/zapr"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// This is the logger we can pass into services
var log logr.Logger

// logInit initializes the logger
func logInit(verbosity int8, pretty bool) {
    
    // assure info logs start at 1
    if verbosity != 0 {
        verbosity = verbosity - 1
    }
    
    // First, define our level-handling logic.
    stdError := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= zapcore.ErrorLevel
    })
    stdOutput := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= zapcore.Level(-verbosity) && lvl < zapcore.ErrorLevel
    })
    
    // High-priority output should also go to standard error, and low-priority
    // output should also go to standard out.
    consoleDebugging := zapcore.Lock(os.Stdout)
    consoleErrors := zapcore.Lock(os.Stderr)
    
    consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
    if pretty {
        consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
    }
    
    // Join the outputs, encoders, and level-handling functions into
    // zapcore.Cores, then tee the four cores together.
    core := zapcore.NewTee(
        zapcore.NewCore(consoleEncoder, consoleErrors, stdError),
        zapcore.NewCore(consoleEncoder, consoleDebugging, stdOutput),
    )
    
    // multiline stacktrace option removed
    options := []zap.Option{} // zap.AddStacktrace(stdError),
    
    if pretty {
        options = append(options, zap.Development(), zap.AddStacktrace(stdError))
    }
    
    // From a zapcore.Core, it's easy to construct a Logger.
    logger := zap.New(core, options...)
    
    log = zapr.NewLogger(logger)
}
```

## Profiling
Profiling is a great way to find bottlenecks in your application.
It can be used to find memory leaks, cpu hogs or to find out where your application spends most of its time.
The cli package supports cpu and memory profiling.

### Initialize profiling
To initialize profiling we need to create the `profilingInit()` function.
This function will be called from the packages' `init()` function if profiling is enabled.
Let's create a new file called `profiling.go` and add the profilingInit function:

```go
package cli

import (
    "os"
    "runtime/pprof"
)

// File names for profiling output
const (
    _cpuprofile = "cpuprofile.pprof"
    _memprofile = "memprofile.pprof"
)

// profiling starts cpu and memory profiling if enabled.
// It returns a function to stop profiling.
func profilingInit(cpuProfile, memProfile bool) func() {
    // doOnStop is a list of functions to be called on stop
    var doOnStop []func()
    // stop calls all necessary functions to stop profiling
    stop := func() {
        for _, d := range doOnStop {
            if d != nil {
                d()
            }
        }
    }
    
    if cpuProfile {
        log.Info("cpu profile enabled")
    
        // Create profiling file
        f, err := os.Create(_cpuprofile)
		if err != nil {
			log.Error(err, "failed to create cpu profile file")
            return stop
		}
    
        // Start profiling
        err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Error(err, "failed to start cpu profile")
            return stop
		}
    
        // Add function to stop cpu profiling to doOnStop list
        doOnStop = append(doOnStop, func() {
            pprof.StopCPUProfile()
            _ = f.Close()
            log.Info("cpu profile stopped")
        })
    }
    
    if memProfile {
        log.Info("memory profile enabled")
    
        // Create profiling file
        f, err := os.Create(_memprofile)
		if err != nil {
			log.Error(err, "failed to create memory profile file")
            return stop
		}
    
        // Add function to stop memory profiling to doOnStop list
        doOnStop = append(doOnStop, func() {
            _ = pprof.WriteHeapProfile(f)
            _ = f.Close()
            log.Info("memory profile stopped")
        })
    }
    
    return stop
}
```

### Profiling with pprof
To profile your application you need to start it with the `--cpuprofile` or `--memprofile` flag (or both).

## 
