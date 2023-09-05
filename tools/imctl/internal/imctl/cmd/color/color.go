// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package color print colors supported by the current terminal.
package color

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/openim-sigs/component-base/util/stringutil"
	"github.com/spf13/cobra"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
)

// ColorOptions is an options struct to support color subcommands.
type ColorOptions struct {
	Type    []string
	Example bool

	genericclioptions.IOStreams
}

var (
	colorLong = templates.LongDesc(`Print the colors supported by the current terminal.

Color lets you use colorized outputs in terms of ANSI Escape Codes in Go (Golang). 
It has support for Windows too! The API can be used in several ways, pick one that suits you.

Find more information at:
    https://github.com/fatih/color`)

	colorExample = templates.Examples(`
		# Print supported foreground and background colors
		imctl color

		# Print supported colors by type
		imctl color -t fg-hi

		# Print all supported colors
		imctl color -t all`)

	availableTypes = []string{"fg", "fg-hi", "bg", "bg-hi", "base", "all"}

	colorCodeExample = templates.Examples(`
### 1. Standard colors

// Print with default helper functions
color.Cyan("Prints text in cyan.")

// A newline will be appended automatically
color.Blue("Prints %s in blue.", "text")

// These are using the default foreground colors
color.Red("We have red")
color.Magenta("And many others ..")

### 2. Mix and reuse colors

// Create a new color object
c := color.New(color.FgCyan).Add(color.Underline)
c.Println("Prints cyan text with an underline.")

// Or just add them to New()
d := color.New(color.FgCyan, color.Bold)
d.Printf("This prints bold cyan %s\n", "too!.")

// Mix up foreground and background colors, create new mixes!
red := color.New(color.FgRed)

boldRed := red.Add(color.Bold)
boldRed.Println("This will print text in bold red.")

whiteBackground := red.Add(color.BgWhite)
whiteBackground.Println("Red text with white background.")

### 3. Use your own output (io.Writer)

// Use your own io.Writer output
color.New(color.FgBlue).Fprintln(myWriter, "blue color!")

blue := color.New(color.FgBlue)
blue.Fprint(writer, "This will print text in blue.")

### 4. Custom print functions (PrintFunc)

// Create a custom print function for convenience
red := color.New(color.FgRed).PrintfFunc()
red("Warning")
red("Error: %s", err)

// Mix up multiple attributes
notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
notice("Don't forget this...")

### 5. Custom fprint functions (FprintFunc)

blue := color.New(FgBlue).FprintfFunc()
blue(myWriter, "important notice: %s", stars)

// Mix up with multiple attributes
success := color.New(color.Bold, color.FgGreen).FprintlnFunc()
success(myWriter, "Don't forget this...")

### 6. Insert into noncolor strings (SprintFunc)

// Create SprintXxx functions to mix strings with other non-colorized strings:
yellow := color.New(color.FgYellow).SprintFunc()
red := color.New(color.FgRed).SprintFunc()
fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
fmt.Printf("This %s rocks!\n", info("package"))

// Use helper functions
fmt.Println("This", color.RedString("warning"), "should be not neglected.")
fmt.Printf("%v %v\n", color.GreenString("Info:"), "an important message.")

// Windows supported too! Just don't forget to change the output to color.Output
fmt.Fprintf(color.Output, "Windows support: %s", color.GreenString("PASS"))

### 7. Plug into existing code

// Use handy standard colors
color.Set(color.FgYellow)

fmt.Println("Existing text will now be in yellow")
fmt.Printf("This one %s\n", "too")

color.Unset() // Don't forget to unset

// You can mix up parameters
color.Set(color.FgMagenta, color.Bold)
defer color.Unset() // Use it in your function

fmt.Println("All text will now be bold magenta.")

### 8. Disable/Enable color
 
There might be a case where you want to explicitly disable/enable color output. the 
go-isatty package will automatically disable color output for non-tty output streams 
(for example if the output were piped directly to less)

Color has support to disable/enable colors both globally and for single color 
definitions. For example suppose you have a CLI app and a --no-color bool flag. You 
can easily disable the color output with:


var flagNoColor = flag.Bool("no-color", false, "Disable color output")

if *flagNoColor {
	color.NoColor = true // disables colorized output
}

It also has support for single color definitions (local). You can
disable/enable color output on the fly:

c := color.New(color.FgCyan)
c.Println("Prints cyan text")

c.DisableColor()
c.Println("This is printed without any color")

c.EnableColor()
c.Println("This prints again cyan...")`)
)

// NewColorOptions returns an initialized ColorOptions instance.
func NewColorOptions(ioStreams genericclioptions.IOStreams) *ColorOptions {
	return &ColorOptions{
		Type:      []string{},
		Example:   false,
		IOStreams: ioStreams,
	}
}

// NewCmdColor returns new initialized instance of color sub command.
func NewCmdColor(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewColorOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "color",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Print colors supported by the current terminal",
		TraverseChildren:      true,
		Long:                  colorLong,
		Example:               colorExample,
		// NoArgs, ArbitraryArgs, OnlyValidArgs, MinimumArgs, MaximumArgs, ExactArgs, ExactArgs
		ValidArgs: []string{},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringSliceVarP(&o.Type, "type", "t", o.Type,
		fmt.Sprintf("Specify the color type, available types: `%s`.", strings.Join(availableTypes, ",")))
	cmd.Flags().BoolVar(&o.Example, "example", o.Example, "Print code to show how to use color package.")

	return cmd
}

// Complete completes all the required options.
func (o *ColorOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(o.Type) == 0 {
		o.Type = []string{"fg", "bg"}

		return nil
	}

	if stringutil.StringIn("all", o.Type) {
		o.Type = []string{"fg", "fg-hi", "bg", "bg-hi", "base"}
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ColorOptions) Validate(cmd *cobra.Command, args []string) error {
	for _, t := range o.Type {
		if !stringutil.StringIn(t, availableTypes) {
			return cmdutil.UsageErrorf(cmd, "--type must be in: %s", strings.Join(availableTypes, ","))
		}
	}

	return nil
}

// Run executes a color subcommand using the specified options.
func (o *ColorOptions) Run(args []string) error {
	if o.Example {
		fmt.Fprint(o.Out, colorCodeExample+"\n")

		return nil
	}

	var data [][]string

	// print table to stdout
	table := tablewriter.NewWriter(o.Out)

	// set table attributes
	table.SetAutoMergeCells(true)
	table.SetRowLine(false)
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Category", "Color Name", "Effect"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor})

	for _, t := range o.Type {
		switch t {
		// Foreground text colors
		case "fg":
			fg := [][]string{
				{"fg", "black", color.BlackString("color.BlackString")},
				{"fg", "red", color.RedString("color.RedString")},
				{"fg", "green", color.GreenString("color.GreenString")},
				{"fg", "yellow", color.YellowString("color.YellowString")},
				{"fg", "blue", color.BlueString("color.BlueString")},
				{"fg", "magenta", color.MagentaString("color.MagentaString")},
				{"fg", "Cyan", color.CyanString("color.CyanString")},
				{"fg", "white", color.WhiteString("color.WhiteString")},
			}
			data = append(data, fg...)
			// Foreground Hi-Intensity text colors
		case "fg-hi":
			fg := [][]string{
				{"fg-hi", "black", color.HiBlackString("color.HiBlackString")},
				{"fg-hi", "red", color.HiRedString("color.HiRedString")},
				{"fg-hi", "green", color.HiGreenString("color.HiGreenString")},
				{"fg-hi", "yellow", color.HiYellowString("color.HiYellowString")},
				{"fg-hi", "blue", color.HiBlueString("color.HiBlueString")},
				{"fg-hi", "magenta", color.HiMagentaString("color.HiMagentaString")},
				{"fg-hi", "Cyan", color.HiCyanString("color.HiCyanString")},
				{"fg-hi", "white", color.HiWhiteString("color.HiWhiteString")},
			}
			data = append(data, fg...)

			// Background text colors
		case "bg":
			bg := [][]string{
				{"bg", "black", color.New(color.FgWhite, color.BgBlack).SprintFunc()("color.BgBlack")},
				{"bg", "red", color.New(color.FgBlack, color.BgRed).SprintFunc()("color.BgRed")},
				{"bg", "greep", color.New(color.FgBlack, color.BgGreen).SprintFunc()("color.BgGreen")},
				{"bg", "yellow", color.New(color.FgWhite, color.BgYellow).SprintFunc()("color.BgYellow")},
				{"bg", "blue", color.New(color.FgWhite, color.BgBlue).SprintFunc()("color.BgBlue")},
				{"bg", "magenta", color.New(color.FgWhite, color.BgMagenta).SprintFunc()("color.BgMagenta")},
				{"bg", "cyan", color.New(color.FgWhite, color.BgCyan).SprintFunc()("color.BgCyan")},
				{"bg", "white", color.New(color.FgBlack, color.BgWhite).SprintFunc()("color.BgWhite")},
			}
			data = append(data, bg...)
			// Background Hi-Intensity text colors
		case "bg-hi":
			bghi := [][]string{
				{"bg-hi", "black", color.New(color.FgWhite, color.BgHiBlack).SprintFunc()("color.BgHiBlack")},
				{"bg-hi", "red", color.New(color.FgBlack, color.BgHiRed).SprintFunc()("color.BgHiRed")},
				{"bg-hi", "greep", color.New(color.FgBlack, color.BgHiGreen).SprintFunc()("color.BgHiGreen")},
				{"bg-hi", "yellow", color.New(color.FgWhite, color.BgHiYellow).SprintFunc()("color.BgHiYellow")},
				{"bg-hi", "blue", color.New(color.FgWhite, color.BgHiBlue).SprintFunc()("color.BgHiBlue")},
				{"bg-hi", "magenta", color.New(color.FgWhite, color.BgHiMagenta).SprintFunc()("color.BgHiMagenta")},
				{"bg-hi", "cyan", color.New(color.FgWhite, color.BgHiCyan).SprintFunc()("color.BgHiCyan")},
				{"bg-hi", "white", color.New(color.FgBlack, color.BgHiWhite).SprintFunc()("color.BgHiWhite")},
			}
			data = append(data, bghi...)
			// Base attributes
		case "base":
			base := [][]string{
				{"base", "black", color.New(color.FgGreen, color.Reset).SprintFunc()("color.Reset")},
				{"base", "red", color.New(color.FgGreen, color.Bold).SprintFunc()("color.Bold")},
				{"base", "greep", color.New(color.FgGreen, color.Faint).SprintFunc()("color.Faint")},
				{"base", "yellow", color.New(color.FgGreen, color.Italic).SprintFunc()("color.Italic")},
				{"base", "blue", color.New(color.FgGreen, color.Underline).SprintFunc()("color.Underline")},
				{"base", "magenta", color.New(color.FgGreen, color.BlinkSlow).SprintFunc()("color.BlinkSlow")},
				{"base", "cyan", color.New(color.FgGreen, color.BlinkRapid).SprintFunc()("color.BlinkRapid")},
				{"base", "white", color.New(color.FgGreen, color.ReverseVideo).SprintFunc()("color.ReverseVideo")},
				{"base", "white", color.New(color.FgGreen, color.Concealed).SprintFunc()("color.Concealed")},
				{"base", "white", color.New(color.FgGreen, color.CrossedOut).SprintFunc()("color.CrossedOut")},
			}
			data = append(data, base...)
		default:
		}
	}

	table.AppendBulk(data)
	table.Render()

	return nil
}
