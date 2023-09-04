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

package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/marmotedu/errors"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu"
	restclient "github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/log"
)

const (
	// DefaultErrorExitCode defines the default exit code.
	DefaultErrorExitCode = 1
)

type debugError interface {
	DebugError() (msg string, args []interface{})
}

var fatalErrHandler = fatal

// BehaviorOnFatal allows you to override the default behavior when a fatal
// error occurs, which is to call os.Exit(code). You can pass 'panic' as a function
// here if you prefer the panic() over os.Exit(1).
func BehaviorOnFatal(f func(string, int)) {
	fatalErrHandler = f
}

// DefaultBehaviorOnFatal allows you to undo any previous override.  Useful in
// tests.
func DefaultBehaviorOnFatal() {
	fatalErrHandler = fatal
}

// fatal prints the message (if provided) and then exits.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = fmt.Errorf("exit")

// CheckErr prints a user-friendly error to STDERR and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
//
// This method is generic to the command in use and may be used by non-IAM
// commands.
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// CheckDiffErr prints a user-friendly error to STDERR and exits with a
// non-zero and non-one exit code. Unrecognized errors will be printed
// with an "error: " prefix.
//
// This method is meant specifically for `iamctl diff` and may be used
// by other commands.
func CheckDiffErr(err error) {
	checkErr(err, func(msg string, code int) {
		fatalErrHandler(msg, code+1)
	})
}

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an iamctl exit code.
func checkErr(err error, handleErr func(string, int)) {
	// unwrap aggregates of 1
	if agg, ok := err.(errors.Aggregate); ok && len(agg.Errors()) == 1 {
		err = agg.Errors()[0]
	}

	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	default:
		switch err := err.(type) {
		case errors.Aggregate:
			handleErr(MultipleErrors(``, err.Errors()), DefaultErrorExitCode)
		default: // for any other error type
			msg, ok := StandardErrorMessage(err)
			if !ok {
				msg = err.Error()
				if !strings.HasPrefix(msg, "error: ") {
					msg = fmt.Sprintf("error: %s", msg)
				}
			}
			handleErr(msg, DefaultErrorExitCode)
		}
	}
}

// StandardErrorMessage translates common errors into a human readable message, or returns
// false if the error is not one of the recognized types. It may also log extended information to klog.
//
// This method is generic to the command in use and may be used by non-IAM
// commands.
func StandardErrorMessage(err error) (string, bool) {
	if debugErr, ok := err.(debugError); ok {
		log.Infof(debugErr.DebugError())
	}
	if t, ok := err.(*url.Error); ok {
		log.Infof("Connection error: %s %s: %v", t.Op, t.URL, t.Err)
		if strings.Contains(t.Err.Error(), "connection refused") {
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf(
				"The connection to the server %s was refused - did you specify the right host or port?",
				host,
			), true
		}

		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	}
	return "", false
}

// MultilineError returns a string representing an error that splits sub errors into their own
// lines. The returned string will end with a newline.
func MultilineError(prefix string, err error) string {
	if agg, ok := err.(errors.Aggregate); ok {
		errs := errors.Flatten(agg).Errors()
		buf := &bytes.Buffer{}
		switch len(errs) {
		case 0:
			return fmt.Sprintf("%s%v\n", prefix, err)
		case 1:
			return fmt.Sprintf("%s%v\n", prefix, messageForError(errs[0]))
		default:
			fmt.Fprintln(buf, prefix)
			for _, err := range errs {
				fmt.Fprintf(buf, "* %v\n", messageForError(err))
			}
			return buf.String()
		}
	}
	return fmt.Sprintf("%s%s\n", prefix, err)
}

// MultipleErrors returns a newline delimited string containing
// the prefix and referenced errors in standard form.
func MultipleErrors(prefix string, errs []error) string {
	buf := &bytes.Buffer{}
	for _, err := range errs {
		fmt.Fprintf(buf, "%s%v\n", prefix, messageForError(err))
	}
	return buf.String()
}

// messageForError returns the string representing the error.
func messageForError(err error) string {
	msg, ok := StandardErrorMessage(err)
	if !ok {
		msg = err.Error()
	}
	return msg
}

// UsageErrorf returns error with command path.
func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

// IsFilenameSliceEmpty checkes where filenames and directory are both zero value.
func IsFilenameSliceEmpty(filenames []string, directory string) bool {
	return len(filenames) == 0 && directory == ""
}

// GetFlagString returns the value of the given flag.
func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

// GetFlagStringSlice can be used to accept multiple argument with flag repetition (e.g. -f arg1,arg2 -f arg3 ...).
func GetFlagStringSlice(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

// GetFlagStringArray can be used to accept multiple argument with flag repetition (e.g. -f arg1 -f arg2 ...).
func GetFlagStringArray(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

// GetFlagBool returns the value of the given flag.
func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return b
}

// GetFlagInt returns the value of the given flag.
func GetFlagInt(cmd *cobra.Command, flag string) int {
	i, err := cmd.Flags().GetInt(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

// GetFlagInt32 returns the value of the given flag.
func GetFlagInt32(cmd *cobra.Command, flag string) int32 {
	i, err := cmd.Flags().GetInt32(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

// GetFlagInt64 returns the value of the given flag.
func GetFlagInt64(cmd *cobra.Command, flag string) int64 {
	i, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

// GetFlagDuration return the value of the given flag.
func GetFlagDuration(cmd *cobra.Command, flag string) time.Duration {
	d, err := cmd.Flags().GetDuration(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return d
}

// AddCleanFlags add clean flags.
func AddCleanFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("user", "u", "", "Specify the user name.")
	cmd.Flags().BoolP("erase", "c", false, "Erase the records from the db")
}

// ValidateOptions defines the validate options.
type ValidateOptions struct {
	EnableValidation bool
}

// IsSiblingCommandExists receives a pointer to a cobra command and a target string.
// Returns true if the target string is found in the list of sibling commands.
func IsSiblingCommandExists(cmd *cobra.Command, targetCmdName string) bool {
	for _, c := range cmd.Parent().Commands() {
		if c.Name() == targetCmdName {
			return true
		}
	}

	return false
}

// DefaultSubCommandRun prints a command's help string to the specified output if no
// arguments (sub-commands) are provided, or a usage error otherwise.
func DefaultSubCommandRun(out io.Writer) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		c.SetOutput(out)
		RequireNoArguments(c, args)
		c.Help()
		CheckErr(ErrExit)
	}
}

// RequireNoArguments exits with a usage error if extra arguments are provided.
func RequireNoArguments(c *cobra.Command, args []string) {
	if len(args) > 0 {
		CheckErr(UsageErrorf(c, "unknown command %q", strings.Join(args, " ")))
	}
}

// ManualStrip is used for dropping comments from a YAML file.
func ManualStrip(file []byte) []byte {
	stripped := []byte{}
	lines := bytes.Split(file, []byte("\n"))
	for i, line := range lines {
		if bytes.HasPrefix(bytes.TrimSpace(line), []byte("#")) {
			continue
		}
		stripped = append(stripped, line...)
		if i < len(lines)-1 {
			stripped = append(stripped, '\n')
		}
	}
	return stripped
}

// Warning write warning message to io.Writer.
func Warning(cmdErr io.Writer, newGeneratorName, oldGeneratorName string) {
	fmt.Fprintf(cmdErr, "WARNING: New generator %q specified, "+
		"but it isn't available. "+
		"Falling back to %q.\n",
		newGeneratorName,
		oldGeneratorName,
	)
}

// CombineRequestErr combines the http response error and error in errs array.
func CombineRequestErr(resp gorequest.Response, body string, errs []error) error {
	var e, sep string
	if len(errs) > 0 {
		for _, err := range errs {
			e = sep + err.Error()
			sep = "\n"
		}
		return errors.New(e)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(body)
	}

	return nil
}

func NewForConfigOrDie() *marmotedu.Clientset {
	clientConfig := &restclient.Config{
		Host:          viper.GetString("server.address"),
		BearerToken:   viper.GetString("user.token"),
		Username:      viper.GetString("user.username"),
		Password:      viper.GetString("user.password"),
		SecretID:      viper.GetString("user.secret-id"),
		SecretKey:     viper.GetString("user.secret-key"),
		Timeout:       viper.GetDuration("server.timeout"),
		MaxRetries:    viper.GetInt("server.max-retries"),
		RetryInterval: viper.GetDuration("server.retry-interval"),
	}

	return marmotedu.NewForConfigOrDie(clientConfig)
}

func TableWriterDefaultConfig(table *tablewriter.Table) *tablewriter.Table {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with two space
	table.SetNoWhiteSpace(true)

	return table
}
