package config

import "flag"

// Flags is the flag set that AddOptions adds to. Test authors should
// also use it instead of directly adding to the global command line.
var Flags = flag.NewFlagSet("", flag.ContinueOnError)

// CopyFlags ensures that all flags that are defined in the source flag
// set appear in the target flag set as if they had been defined there
// directly. From the flag package it inherits the behavior that there
// is a panic if the target already contains a flag from the source.
func CopyFlags(source *flag.FlagSet, target *flag.FlagSet) {
	source.VisitAll(func(flag *flag.Flag) {
		// We don't need to copy flag.DefValue. The original
		// default (from, say, flag.String) was stored in
		// the value and gets extracted by Var for the help
		// message.
		target.Var(flag.Value, flag.Name, flag.Usage)
	})
}
