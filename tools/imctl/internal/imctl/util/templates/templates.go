// Package templates provides template functions for working with templates.
package templates

import (
	"strings"
	"unicode"
)

const (
	// SectionVars is the help template section that declares variables to be used in the template.
	SectionVars = `{{$isRootCmd := isRootCmd .}}` +
		`{{$rootCmd := rootCmd .}}` +
		`{{$visibleFlags := visibleFlags (flagsNotIntersected .LocalFlags .PersistentFlags)}}` +
		`{{$explicitlyExposedFlags := exposed .}}` +
		`{{$optionsCmdFor := optionsCmdFor .}}` +
		`{{$usageLine := usageLine .}}`

	// SectionAliases is the help template section that displays command aliases.
	SectionAliases = `{{if gt .Aliases 0}}Aliases:
{{.NameAndAliases}}

{{end}}`

	// SectionExamples is the help template section that displays command examples.
	SectionExamples = `{{if .HasExample}}Examples:
{{trimRight .Example}}

{{end}}`

	// SectionSubcommands is the help template section that displays the command's subcommands.
	SectionSubcommands = `{{if .HasAvailableSubCommands}}{{cmdGroupsString .}}

{{end}}`

	// SectionFlags is the help template section that displays the command's flags.
	SectionFlags = `{{ if or $visibleFlags.HasFlags $explicitlyExposedFlags.HasFlags}}Options:
{{ if $visibleFlags.HasFlags}}{{trimRight (flagsUsages $visibleFlags)}}{{end}}{{ if $explicitlyExposedFlags.HasFlags}}{{ if $visibleFlags.HasFlags}}
{{end}}{{trimRight (flagsUsages $explicitlyExposedFlags)}}{{end}}

{{end}}`

	// SectionUsage is the help template section that displays the command's usage.
	SectionUsage = `{{if and .Runnable (ne .UseLine "") (ne .UseLine $rootCmd)}}Usage:
  {{$usageLine}}

{{end}}`

	// SectionTipsHelp is the help template section that displays the '--help' hint.
	SectionTipsHelp = `{{if .HasSubCommands}}Use "{{$rootCmd}} <command> --help" for more information about a given command.
{{end}}`

	// SectionTipsGlobalOptions is the help template section that displays the 'options' hint for displaying global
	// flags.
	SectionTipsGlobalOptions = `{{if $optionsCmdFor}}Use "{{$optionsCmdFor}}" for a list of global command-line options (applies to all commands).
{{end}}`
)

// MainHelpTemplate if the template for 'help' used by most commands.
func MainHelpTemplate() string {
	return `{{with or .Long .Short }}{{. | trim}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// MainUsageTemplate if the template for 'usage' used by most commands.
func MainUsageTemplate() string {
	sections := []string{
		"\n\n",
		SectionVars,
		SectionAliases,
		SectionExamples,
		SectionSubcommands,
		SectionFlags,
		SectionUsage,
		SectionTipsHelp,
		SectionTipsGlobalOptions,
	}
	return strings.TrimRightFunc(strings.Join(sections, ""), unicode.IsSpace)
}

// OptionsHelpTemplate if the template for 'help' used by the 'options' command.
func OptionsHelpTemplate() string {
	return ""
}

// OptionsUsageTemplate if the template for 'usage' used by the 'options' command.
func OptionsUsageTemplate() string {
	return `{{ if .HasInheritedFlags}}The following options can be passed to any command:

{{flagsUsages .InheritedFlags}}{{end}}`
}
