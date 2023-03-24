# Golang Extended Cli Boilerplate Example

Using spf13/viper to gather parameters from various sources can be tricky. You will probably run into the problem of having to redefine default parameters for each means of input. This example shows how to to load parameters from different sources without having to define them multiple times. The common priority for input methods is this:

1. CLI parameters
2. Environment variables
3. Config file
4. Default values

Checkout the [cli package](/base/cli) on how to setup this config.

Also included in this example is adding `Version` and `CommitHash` to your application through ld flags. Checkout the [Makefile](/extended/Makefile) on how to load them during build from git tags and git commit hash.

It is often useful to get performance insights from your application. That is when you want to add pprof to your application. When doing this from the start of your project you're more likely to find bottlenecks or memory leaks earlier. The rampup and -down of the profiling is handled in [/cli/profiling.go](/extended/cli/profiling.go)

## Usage 
Build the application

```shell
make build
```

Parameters can be passed into the application in these 4 ways:

config.yaml:
```yaml
param1: "param from config"
```

defaultConfig in [config.go](/base/cli/config.go):
```go
var defaultConfig = config{
    Param2: "param with default value",
}
```

Environment variable
```shell
PARAM1="param from env" ./our_app
```

Cli parameter
```shell
./our_app --param1="param from cli"
```

For profiling just add the parameters `cpuprofile` and/or `memprofile` to the application call. This will generate the profiling files `cpuprofile.pprof` and `memprofile.pprof` which can be inspected with [github.com/google/pprof](https://github.com/google/pprof)
