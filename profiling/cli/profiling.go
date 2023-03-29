package cli

import (
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
)

// profilingInit starts cpu and memory profiling if enabled.
// It returns a function to stop profiling.
func profilingInit() func() {
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

	if GlobalConfig.CPUProfile {
		fmt.Println("cpu profile enabled")

		// Create profiling file
		f, err := os.Create(GlobalConfig.CPUProfileFile)
		if err != nil {
			fmt.Println("could not create cpu profile file")
			return stop
		}

		// Start profiling
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Println("could not start cpu profiling")
			return stop
		}

		// Add function to stop cpu profiling to doOnStop list
		doOnStop = append(doOnStop, func() {
			pprof.StopCPUProfile()
			_ = f.Close()
			fmt.Println("cpu profile stopped")
		})
	}

	if GlobalConfig.MEMProfile {
		fmt.Println("memory profile enabled")

		// Create profiling file
		f, err := os.Create(GlobalConfig.MEMProfileFile)
		if err != nil {
			fmt.Println("could not create memory profile file")
			return stop
		}

		// Add function to stop memory profiling to doOnStop list
		doOnStop = append(doOnStop, func() {
			_ = pprof.WriteHeapProfile(f)
			_ = f.Close()
			fmt.Println("memory profile stopped")
		})
	}

	return stop
}

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
}
