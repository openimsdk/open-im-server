package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/pflag"
)

// profiling configuration variables
var (
	profileName   string = "none"                       // Name of the profile to capture.
	profileOutput string = "profile.pprof"              // File to write the profile data.
)

// addProfilingFlags registers profiling related flags to the given FlagSet.
func addProfilingFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&profileName,
		"profile",
		"none",
		"Type of profile to capture. Options: none, cpu, heap, goroutine, threadcreate, block, mutex",
	)
	flags.StringVar(&profileOutput, "profile-output", "profile.pprof", "File to write the profile data")
}

// initProfiling sets up profiling based on the user's choice. 
// If 'cpu' is selected, it starts the CPU profile. For block and mutex profiles, 
// sampling rates are set up.
func initProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		f, err := os.Create(profileOutput)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(f)
	case "block":
		runtime.SetBlockProfileRate(1)  // Sampling every block event
		return nil
	case "mutex":
		runtime.SetMutexProfileFraction(1)  // Sampling every mutex event
		return nil
	default:
		if profile := pprof.Lookup(profileName); profile == nil {
			return fmt.Errorf("unknown profile type: '%s'", profileName)
		}
		return nil
	}
}

// flushProfiling writes the profiling data to the specified file. 
// For heap profiles, it runs the GC before capturing the data. 
// It stops the CPU profile if it was started.
func flushProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
		return nil
	case "heap":
		runtime.GC() // Run garbage collection before writing heap profile
		fallthrough
	default:
		profile := pprof.Lookup(profileName)
		if profile == nil {
			return errors.New("invalid profile type")
		}
		f, err := os.Create(profileOutput)
		if err != nil {
			return err
		}
		return profile.WriteTo(f, 0)
	}
}
