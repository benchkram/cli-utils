package cli

import (
	"os"
	"runtime/pprof"
	"sync"

	"github.com/benchkram/errz"
)

// File names for profiling output
const (
	_cpuprofile = "cpuprofile.pprof"
	_memprofile = "memprofile.pprof"
)

// profilingInit starts cpu and memory profiling if enabled.
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
		errz.Fatal(err)

		// Start profiling
		err = pprof.StartCPUProfile(f)
		errz.Log(err)

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
		errz.Fatal(err)

		// Add function to stop memory profiling to doOnStop list
		doOnStop = append(doOnStop, func() {
			_ = pprof.WriteHeapProfile(f)
			_ = f.Close()
			log.Info("memory profile stopped")
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
