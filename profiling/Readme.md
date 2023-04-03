# Golang Profiling Cli Boilerplate Example

Simple example of how to profile a golang application which uses [spf13/cobra](github.com/spf13/cobra) for configuration.
Checkout [profiling.go](profiling.go) on how to start and stop the profiling.

## Usage
App can take following parameters:
- `--cpu-profile` - enables cpu profiling (default false)
- `--mem-profile` - enables memory profiling (default false)
- `--cpu-profile-file` - sets the cpu profile file name (default "cpu.prof")
- `--mem-profile-file` - sets the memory profile file name (default "mem.prof")

For profiling just add the parameters `--cpu-profile` and/or `--mem-profile` to the application call.
This will generate the profiling files `cpu.prof` and `mem.prof` which can be inspected with [github.com/google/pprof](https://github.com/google/pprof)

## Example

Create a cpu profile:
```shell
go run . --cpu-profile
go tool pprof cpu.prof
```

Create a memory profile:
```shell
go run . --mem-profile
go tool pprof mem.prof
```
