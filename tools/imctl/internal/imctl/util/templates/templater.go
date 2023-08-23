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

package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/util/term"
)

type FlagExposer interface {
	ExposeFlags(cmd *cobra.Command, flags ...string) FlagExposer
}

func ActsAsRootCommand(cmd *cobra.Command, filters []string, groups ...CommandGroup) FlagExposer {
	if cmd == nil {
		panic("nil root command")
	}
	templater := &templater{
		RootCmd:       cmd,
		UsageTemplate: MainUsageTemplate(),
		HelpTemplate:  MainHelpTemplate(),
		CommandGroups: groups,
		Filtered:      filters,
	}
	cmd.SetFlagErrorFunc(templater.FlagErrorFunc())
	cmd.SilenceUsage = true
	cmd.SetUsageFunc(templater.UsageFunc())
	cmd.SetHelpFunc(templater.HelpFunc())
	return templater
}

func UseOptionsTemplates(cmd *cobra.Command) {
	templater := &templater{
		UsageTemplate: OptionsUsageTemplate(),
		HelpTemplate:  OptionsHelpTemplate(),
	}
	cmd.SetUsageFunc(templater.UsageFunc())
	cmd.SetHelpFunc(templater.HelpFunc())
}

type templater struct {
	UsageTemplate string
	HelpTemplate  string
	RootCmd       *cobra.Command
	CommandGroups
	Filtered []string
}

func (t *templater) FlagErrorFunc(exposedFlags ...string) func(*cobra.Command, error) error {
	return func(c *cobra.Command, err error) error {
		c.SilenceUsage = true
		switch c.CalledAs() {
		case "options":
			return fmt.Errorf("%s\nrun '%s' without flags", err, c.CommandPath())
		default:
			return fmt.Errorf("%s\nsee '%s --help' for usage", err, c.CommandPath())
		}
	}
}

func (t *templater) ExposeFlags(cmd *cobra.Command, flags ...string) FlagExposer {
	cmd.SetUsageFunc(t.UsageFunc(flags...))
	return t
}

func (t *templater) HelpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		tt := template.New("help")
		tt.Funcs(t.templateFuncs())
		template.Must(tt.Parse(t.HelpTemplate))
		out := term.NewResponsiveWriter(c.OutOrStdout())
		err := tt.Execute(out, c)
		if err != nil {
			c.Println(err)
		}
	}
}

func (t *templater) UsageFunc(exposedFlags ...string) func(*cobra.Command) error {
	return func(c *cobra.Command) error {
		tt := template.New("usage")
		tt.Funcs(t.templateFuncs(exposedFlags...))
		template.Must(tt.Parse(t.UsageTemplate))
		out := term.NewResponsiveWriter(c.OutOrStderr())
		return tt.Execute(out, c)
	}
}

func (t *templater) templateFuncs(exposedFlags ...string) template.FuncMap {
	return template.FuncMap{
		"trim":                strings.TrimSpace,
		"trimRight":           func(s string) string { return strings.TrimRightFunc(s, unicode.IsSpace) },
		"trimLeft":            func(s string) string { return strings.TrimLeftFunc(s, unicode.IsSpace) },
		"gt":                  cobra.Gt,
		"eq":                  cobra.Eq,
		"rpad":                rpad,
		"appendIfNotPresent":  appendIfNotPresent,
		"flagsNotIntersected": flagsNotIntersected,
		"visibleFlags":        visibleFlags,
		"flagsUsages":         flagsUsages,
		"cmdGroups":           t.cmdGroups,
		"cmdGroupsString":     t.cmdGroupsString,
		"rootCmd":             t.rootCmdName,
		"isRootCmd":           t.isRootCmd,
		"optionsCmdFor":       t.optionsCmdFor,
		"usageLine":           t.usageLine,
		"exposed": func(c *cobra.Command) *flag.FlagSet {
			exposed := flag.NewFlagSet("exposed", flag.ContinueOnError)
			if len(exposedFlags) > 0 {
				for _, name := range exposedFlags {
					if flag := c.Flags().Lookup(name); flag != nil {
						exposed.AddFlag(flag)
					}
				}
			}
			return exposed
		},
	}
}

func (t *templater) cmdGroups(c *cobra.Command, all []*cobra.Command) []CommandGroup {
	if len(t.CommandGroups) > 0 && c == t.RootCmd {
		all = filter(all, t.Filtered...)
		return AddAdditionalCommands(t.CommandGroups, "Other Commands:", all)
	}
	all = filter(all, "options")
	return []CommandGroup{
		{
			Message:  "Available Commands:",
			Commands: all,
		},
	}
}

func (t *templater) cmdGroupsString(c *cobra.Command) string {
	groups := []string{}
	for _, cmdGroup := range t.cmdGroups(c, c.Commands()) {
		cmds := []string{cmdGroup.Message}
		for _, cmd := range cmdGroup.Commands {
			if cmd.IsAvailableCommand() {
				cmds = append(cmds, "  "+rpad(cmd.Name(), cmd.NamePadding())+" "+cmd.Short)
			}
		}
		groups = append(groups, strings.Join(cmds, "\n"))
	}
	return strings.Join(groups, "\n\n")
}

func (t *templater) rootCmdName(c *cobra.Command) string {
	return t.rootCmd(c).CommandPath()
}

func (t *templater) isRootCmd(c *cobra.Command) bool {
	return t.rootCmd(c) == c
}

func (t *templater) parents(c *cobra.Command) []*cobra.Command {
	parents := []*cobra.Command{c}
	for current := c; !t.isRootCmd(current) && current.HasParent(); {
		current = current.Parent()
		parents = append(parents, current)
	}
	return parents
}

func (t *templater) rootCmd(c *cobra.Command) *cobra.Command {
	if c != nil && !c.HasParent() {
		return c
	}
	if t.RootCmd == nil {
		panic("nil root cmd")
	}
	return t.RootCmd
}

func (t *templater) optionsCmdFor(c *cobra.Command) string {
	if !c.Runnable() {
		return ""
	}
	rootCmdStructure := t.parents(c)
	for i := len(rootCmdStructure) - 1; i >= 0; i-- {
		cmd := rootCmdStructure[i]
		if _, _, err := cmd.Find([]string{"options"}); err == nil {
			return cmd.CommandPath() + " options"
		}
	}
	return ""
}

func (t *templater) usageLine(c *cobra.Command) string {
	usage := c.UseLine()
	suffix := "[options]"
	if c.HasFlags() && !strings.Contains(usage, suffix) {
		usage += " " + suffix
	}
	return usage
}

func flagsUsages(f *flag.FlagSet) string {
	x := new(bytes.Buffer)

	f.VisitAll(func(flag *flag.Flag) {
		if flag.Hidden {
			return
		}
		format := "--%s=%s: %s\n"

		if flag.Value.Type() == "string" {
			format = "--%s='%s': %s\n"
		}

		if len(flag.Shorthand) > 0 {
			format = "  -%s, " + format
		} else {
			format = "   %s   " + format
		}

		fmt.Fprintf(x, format, flag.Shorthand, flag.Name, flag.DefValue, flag.Usage)
	})

	return x.String()
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

func appendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) {
		return s
	}
	return s + " " + stringToAppend
}

func flagsNotIntersected(l *flag.FlagSet, r *flag.FlagSet) *flag.FlagSet {
	f := flag.NewFlagSet("notIntersected", flag.ContinueOnError)
	l.VisitAll(func(flag *flag.Flag) {
		if r.Lookup(flag.Name) == nil {
			f.AddFlag(flag)
		}
	})
	return f
}

func visibleFlags(l *flag.FlagSet) *flag.FlagSet {
	hidden := "help"
	f := flag.NewFlagSet("visible", flag.ContinueOnError)
	l.VisitAll(func(flag *flag.Flag) {
		if flag.Name != hidden {
			f.AddFlag(flag)
		}
	})
	return f
}

func filter(cmds []*cobra.Command, names ...string) []*cobra.Command {
	out := []*cobra.Command{}
	for _, c := range cmds {
		if c.Hidden {
			continue
		}
		skip := false
		for _, name := range names {
			if name == c.Name() {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		out = append(out, c)
	}
	return out
}
