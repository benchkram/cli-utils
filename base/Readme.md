# Golang Base Cli Boilerplate Example

Using spf13/viper to gather parameters from various sources can be tricky. You will probably run into the problem of having to redefine default parameters for each means of input. This example shows how to to load parameters from different sources without having to define them multiple times. The common priority for input methods is this:

1. CLI parameters
2. Environment variables
3. Config file
4. Default values

Checkout the [cli package](/base/cli) on how to setup this config.

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
PARAM1="param from env" ./cli_example
```

Cli parameter
```shell
./cli_example --param1="param from cli"
```
