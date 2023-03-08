# Cli boilerplate for simple and complex applications


A boilerplate is a pre-built template or framework that can be used as a starting point for new software projects.
There are several benefits to using a boilerplate for new projects, including:

- **Time-saving**: One of the most significant benefits of using a boilerplate is that it saves time. 
Rather than starting from scratch, developers can use an existing codebase that has already been tested, optimized, and documented. 
This can significantly speed up the development process, allowing developers to focus on implementing specific features rather than building the foundation of the application.

- **Consistency**: Using a boilerplate can help ensure consistency across projects.
By using the same structure, coding conventions, and design patterns, developers can maintain a consistent codebase and ensure that all projects are built in a similar way.

- **Best practices**: Boilerplates often incorporate best practices for development, such as security measures, 
code quality standards, and performance optimization techniques. By using a boilerplate, developers can leverage these
best practices without having to research and implement them themselves.

- **Community support**: Many boilerplates have large communities of developers who contribute to their development and 
offer support. This can be a valuable resource for developers who are new to a particular technology or who encounter 
issues while working on their projects.

- **Faster onboarding**: When new developers join a team or take over a project, a boilerplate can help them quickly 
get up to speed on the codebase and development processes. This can help reduce ramp-up time and ensure that new
developers are productive quickly.

Overall, using a boilerplate can save time, promote consistency, incorporate best practices, offer community support,
and facilitate faster onboarding of new team members. These benefits can lead to more efficient and effective software
development, better code quality, and a more consistent user experience across projects.

The two main parts every CLI application should have are:
- **Config**: The configuration of the application.
- **Profiling**: The cpu and memory profiling of the application.

## Config
Every CLI application should provide multiple ways to take parameters to ensure flexibility and ease-of-use for users.

Command-line arguments are a common way for users to pass options and values to an application when it starts. It allows users to provide inputs quickly and can be used to automate the application.

Environment variables are used to set values that can be accessed by an application at runtime. This approach is helpful for providing sensitive information such as passwords, tokens or keys, which should not be hardcoded into an application.

Config files provide a way to configure the application settings, including defaults and overrides. This approach is helpful when users want to save their preferences or when multiple users share the same configuration.

Providing multiple ways to pass parameters allows users to choose the most appropriate method for their use case, and provides the application with the flexibility to adapt to various scenarios.

## Profiling
Memory and CPU profiling are important techniques for identifying performance bottlenecks and optimizing the performance 
of software applications. Here are some key reasons why every application needs some kind of memory and CPU profiling:

- **Identifying memory leaks**: Memory profiling can help identify these leaks and enable developers to fix them.
- **Optimizing CPU usage**: CPU profiling helps identify code that takes a lot of CPU time to execute, which can then be 
optimized or parallelized to improve performance.
- **Finding hotspots**: Profiling can identify "hotspots" in the code, where the program spends a significant amount of
time, allowing developers to optimize these areas for better performance.
- **Improving scalability**: Profiling can help identify areas where the application's performance is limited by the 
available hardware resources, enabling developers to optimize the code or infrastructure for better scalability.

Overall, memory and CPU profiling are essential tools for identifying performance issues, improving the efficiency of 
the application, and enhancing the user experience.

## Where to start

When creating a new CLI application, it's helpful to start with a boilerplate that provides a foundation for the application.
To be reusable across multiple projects, the boilerplate should be designed to be flexible and customizable, and should
include the most common features that are needed for most applications.

### Don't overload the root of your project
The root of your project should only contain the most important files.
From my experience it is the best to only have your main.go file in the root of your project.
Any other logic should be in a subpackage.

Further the main.go file should not do more than loading the config and starting the application.
Also all cli related logic should be in a subpackage. Definition of parameters and different run options should all be
be placed into that package. That means that the main function should only contain the following:

```go
func main() {
	// This main function is only responsible for calling the cli package
	// and handling errors returned by the cli package
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

### Defining the config
Since we probably always need arguments in one way or another let's first design a config struct that can be used to define
all the parameters we need. Therefor we create a config.go file in the cli package of our project and define a struct
similar to the following:

```go
type config struct {
	Param1 string `mapstructure:"param1" structs:"param1" env:"PARAM1"`
	Param2 int `mapstructure:"param2" structs:"param2" env:"PARAM2"`
}
```
As you can see the parameters are tagged with `mapstructure` and `structs` tags. This is needed to be able to use the
[spf13/viper](https://github.com/spf13/viper) package in combination with the [fatih/structs](github.com/fatih/structs) 
package to load the config from different sources.

This allows us already to define default values for our parameters in the config struct. For example if we want to define
some default values we can do it like this:

```go
var defaultConfig = config{
    Param1: "Hello World",
    Param2: 42,
}
```
### The root command
The root command is the command that is executed when no other command is specified. We'll be using [spf13/cobra](github.com/spf13/cobra)
to define our commands. For convenience I recommend to create a go file for each command. In our case we only have one
command which is the root command. So we create a root_cmd.go file in the cli package and define the root command like this:

```go
var rootCmd = &cobra.Command{
	Use:   "cli_example",
	Short: "cli to start example application
	Long:  "cli to start example application",
	Run: func(cmd *cobra.Command, args []string) {
		// here is where we start the actual application
	},
}
```
This will be the command we run in the `cli.Execute` function we called from the main.go file.
```go
func Execute() error {
    return rootCmd.Execute()
}
```
With that complete we can move on to initializing the config.

### Initializing the config
Now that we have defined our config struct we need to initialize it. This is done in the `init` function of the cli package.
That step is important because we want to make sure that the config is initialized before the main function is being run
and the cli package is used.

For initialization these three steps are needed:
1. Define the cli flags on the root command
2. Bind the same flags to viper
3. Bind the environment variables to viper

The first step is done by defining the flags on the root command. This is done by using the `PersistentFlags` method of
the root command. This is important because we want to make sure that the flags are available for all subcommands as well.
The flags are defined like this:

```go
func cliFlags() {
    rootCmd.PersistentFlags().String("param1", defaultConfig.Param1, "param1 description")
    rootCmd.PersistentFlags().Int("param2", defaultConfig.Param2, "param2 description")
}
```

Similar to defining the flags on cobra in step 1 the viper binding in step 2 and 3 need to be done for each flag.
This requires a call to `viper.BindPFlag` and `viper.BindEnv` for each flag. To avoid code duplication we can create
a combined call. The tags we added previously to the config struct are used to define the name of the flag and the name
of the environment variable. The following code snippet demonstrates how the binding can be accomplished:

```go
func bindFlagsAndEnv() {
    for _, field := range structs.Fields(&config{}) {
		key := field.Tag("structs")
		env := field.Tag("env")
        viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key))
        viper.BindEnv(key, env)
    }
}
```
This will prevent us from having to change the binding code if we add or remove parameters from the config struct and
save us some time.

Now that we have a function for each of the three steps we can combine them into a single function that initializes the
config. This function can then be called by the packages init function. If your cli package is not calling any other
initialization functions you can also call the function directly from the init function.

```go
func configInit() {
    cliFlags()
	bindFlagsAndEnv()
}

func init() {
	configInit()
	// other init functions
	// e.g.
	// rootCmd.AddCommand(subCmd)
}
```

### Loading the config
Now that we have initialized the config we need to load it. The very first step here is to set the defaults. Again we
can use the tags we added to the config struct to define the default values. That way we don't have to write a call to `SetDefault`
for every parameter. The following code snippet demonstrates how the defaults can be set:

```go
defaultsAsMap := structs.Map(defaultConfig)

// Set defaults
for key, value := range defaultsAsMap {
    viper.SetDefault(key, value)
}
```

At this point we can load the config. This is a very simple step. We just need to tell viper the name of the config file
and the path where it can be found. In our case the config file is called `config` and it is located in the current
working directory. Of course, you could also create a cli flag to define the path to the config file.

```go
viper.SetConfigName("config")
viper.AddConfigPath(".")
if err := viper.ReadInConfig(); err == nil {
fmt.Println("Using config file:", viper.ConfigFileUsed())
}
```

Last step is bind the config to our config struct. This is done by using the `Unmarshal` method of viper. This method
takes a pointer to a struct and will fill the struct with the values from the config and either return it or store it in
a global config variable. In our case we want to return the config struct so we can pass it to our application.

```go
// Unmarshal config into struct
c := &config{}
err := viper.Unmarshal(c)
if err != nil {
    return nil, err
}
return c, nil
```

And voi la, we have a config struct that is loaded from the config file, the environment and the cli flags. It is now
ready to be passed to our application.

### Let's do some profiling next
Now that we have a config struct that is loaded from different sources we can start profiling our application. For this
we'll be using the [pprof](https://golang.org/pkg/runtime/pprof/) package that is part of the standard library. The
package provides a set of functions that can be used to profile the application. 
Often the `pprof.StartCPUProfile` and `pprof.StopCPUProfile` functions are used right inside the main function. A better
approach is to use these functions from withing the cli package. This way we can also define for which commands we want
to profile the application. The profiling can still be run by passing e.g. the `--profile` flag to the cli or use more parameters
to define the profiling output paths. So there's no downside to not using it in the main function.

One difference to doing the profiling in the main function is that we need to `defer` the `StopCPUProfile` function and
also the call to `WriteHeapProfile` needs to be done after the application is stopped.  To accomplish this we can simply
return an `onStop` function from the call to initialize the profiling and make sure that it is only called once after the
application has stopped. Best place to do this is in the `Execute` function of the cli package.

```go
// onStopProfiling is called when the cli exits
// profilingOnce makes sure it's only called once
var onStopProfiling func()
var profilingOnce sync.Once

// stopProfiling triggers _stopProfiling.
// It's safe to be called multiple times.
func stopProfiling() {
    if onStopProfiling != nil {
    profilingOnce.Do(onStopProfiling)
}
```

The profilingInit function can be called from the `PreRun` function of the cobra commands. This way the profiling is
only started when the command is run. The `PreRun` function is called before the command is executed.
To be able to stop the profiling from outside the `profilingInit` function we need to return a `stop` function which we 
can define like this:

```go
// doOnStop is a list of functions to be called on stop
var doOnStop []func()
// stop calls all functions in doOnStop
func stop() {
    for _, f := range doOnStop {
		if f != nil {
            f()
        }
    }
}
```
The stop function is useful because it provides a way to cleanly stop or shutdown a program or process. In this particular
implementation, the stop function works by iterating over a list of functions called doOnStop, which contains all the 
necessary functions that need to be executed before the program or process can be stopped.

By calling stop, all the functions in the doOnStop list are executed in the order they were added, ensuring that all 
necessary cleanup tasks are performed before the program or process is terminated. This is important because it helps 
avoid leaving any resources or data in an inconsistent state or causing any potential issues with future runs of the program.

Additionally, the doOnStop list allows for greater flexibility in the program's design, as it allows developers to add or 
remove functions from the list depending on the program's requirements. This makes it easier to customize the program's 
behavior and ensure that all necessary cleanup tasks are performed.

This implementation allows us to conditionally add cpu and memory profiling and append the stop function to the doOnStop
list.

```go
if cpuProfile {
    // Create profiling file
    f, err := os.Create("cpu_profile_path")
    if err != nil {
        return nil, err
    }

    // Start profiling
    err = pprof.StartCPUProfile(f)
    if err != nil {
        return nil, err
    }

    // Add function to stop cpu profiling to doOnStop list
    doOnStop = append(doOnStop, func() {
        pprof.StopCPUProfile()
        _ = f.Close()
    })
}

if memProfile {
    // Create profiling file
    f, err := os.Create("mem_profile_path")
    if err != nil {
        return nil, err
    }
    
    // Add function to stop memory profiling to doOnStop list
    doOnStop = append(doOnStop, func() {
        _ = pprof.WriteHeapProfile(f)
        _ = f.Close()
    })
}
```
### Conclusion
In this post we have looked at how to load a config struct from different sources. We have also looked at how to profile
an application and how to stop the profiling. We looked at the benefits of having this code in a separate cli package.
The code for this post can be found [here](github.com/benchkram/cli_example).


# Left out in this post
- Version?
- Redact?
- Logging