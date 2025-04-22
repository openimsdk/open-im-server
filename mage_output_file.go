//go:build ignore
// +build ignore

package main

import (
	"context"
	_flag "flag"
	_fmt "fmt"
	_ioutil "io/ioutil"
	_log "log"
	"os"
	"os/signal"
	_filepath "path/filepath"
	_sort "sort"
	"strconv"
	_strings "strings"
	"syscall"
	_tabwriter "text/tabwriter"
	"time"
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose bool          // print out log statements
		List    bool          // print out a list of targets
		Help    bool          // print out help for a specific target
		Timeout time.Duration // set a timeout to running the targets
		Args    []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}
		b, err := strconv.ParseBool(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}
		d, err := time.ParseDuration(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := _flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		_fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], _filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}

	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}

	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func(i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := _strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
	var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return _fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {

		targets := map[string]string{
			"build*": "",
			"check":  "",
			"start":  "",
			"stop":   "",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		_sort.Strings(keys)

		_fmt.Println("Targets:")
		w := _tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			_fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
		if err == nil {
			_fmt.Println("\n* default target")
		}
		return err
	}

	var ctx context.Context
	ctxCancel := func() {}

	// by deferring in a closure, we let the cancel function get replaced
	// by the getContext function.
	defer func() {
		ctxCancel()
	}()

	getContext := func() (context.Context, func()) {
		if ctx == nil {
			if args.Timeout != 0 {
				ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
			} else {
				ctx, ctxCancel = context.WithCancel(context.Background())
			}
		}

		return ctx, ctxCancel
	}

	runTarget := func(logger *_log.Logger, fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT)
		select {
		case <-sigCh:
			logger.Println("cancelling mage targets, waiting up to 5 seconds for cleanup...")
			cancel()
			cleanupCh := time.After(5 * time.Second)

			select {
			// target exited by itself
			case err = <-d:
				return err
			// cleanup timeout exceeded
			case <-cleanupCh:
				return _fmt.Errorf("cleanup timeout exceeded")
			// second SIGINT received
			case <-sigCh:
				logger.Println("exiting mage")
				return _fmt.Errorf("exit forced")
			}
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			_fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			// we intentionally don't cancel the context here, because
			// the next target will need to run with the same context.
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *_log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	_log.SetFlags(0)
	if !args.Verbose {
		_log.SetOutput(_ioutil.Discard)
	}
	logger := _log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			_log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch _strings.ToLower(args.Args[0]) {
		case "build":

			_fmt.Print("Usage:\n\n\tmage build\n\n")
			var aliases []string
			if len(aliases) > 0 {
				_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
			}
			return
		case "check":

			_fmt.Print("Usage:\n\n\tmage check\n\n")
			var aliases []string
			if len(aliases) > 0 {
				_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
			}
			return
		case "start":

			_fmt.Print("Usage:\n\n\tmage start\n\n")
			var aliases []string
			if len(aliases) > 0 {
				_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
			}
			return
		case "stop":

			_fmt.Print("Usage:\n\n\tmage stop\n\n")
			var aliases []string
			if len(aliases) > 0 {
				_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
			}
			return
		default:
			logger.Printf("Unknown target: %q\n", args.Args[0])
			os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		ignoreDefault, _ := strconv.ParseBool(os.Getenv("MAGEFILE_IGNOREDEFAULT"))
		if ignoreDefault {
			if err := list(); err != nil {
				logger.Println("Error:", err)
				os.Exit(1)
			}
			return
		}

		wrapFn := func(ctx context.Context) error {
			Build()
			return nil
		}
		ret := runTarget(logger, wrapFn)
		handleError(logger, ret)
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch _strings.ToLower(target) {

		}

		switch _strings.ToLower(target) {

		case "build":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build")
			}

			wrapFn := func(ctx context.Context) error {
				Build()
				return nil
			}
			ret := runTarget(logger, wrapFn)
			handleError(logger, ret)
		case "check":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Check\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Check")
			}

			wrapFn := func(ctx context.Context) error {
				Check()
				return nil
			}
			ret := runTarget(logger, wrapFn)
			handleError(logger, ret)
		case "start":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Start\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Start")
			}

			wrapFn := func(ctx context.Context) error {
				Start()
				return nil
			}
			ret := runTarget(logger, wrapFn)
			handleError(logger, ret)
		case "stop":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Stop\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Stop")
			}

			wrapFn := func(ctx context.Context) error {
				Stop()
				return nil
			}
			ret := runTarget(logger, wrapFn)
			handleError(logger, ret)

		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}
