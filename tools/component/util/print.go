package util

import (
	"fmt"
	"log"
	"os"
)

const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
)

// colorErrPrint prints formatted string in red to stderr
func ColorErrPrint(msg string) {
	// ANSI escape code for red text
	const redColor = "\033[31m"
	// ANSI escape code to reset color
	const resetColor = "\033[0m"
	msg = redColor + msg + resetColor
	// Print to stderr in red
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func ColorSuccessPrint(format string, a ...interface{}) {
	// ANSI escape code for green text is \033[32m
	// \033[0m resets the color
	fmt.Printf("\033[32m"+format+"\033[0m", a...)
}

func colorPrint(colorCode int, format string, a ...any) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", colorCode, fmt.Sprintf(format, a...))
}

func colorErrPrint(colorCode int, format string, a ...any) {
	log.Printf("\x1b[%dm%s\x1b[0m\n", colorCode, fmt.Sprintf(format, a...))
}

func ErrorPrint(s string) {
	colorErrPrint(colorRed, "%v", s)
}

func SuccessPrint(s string) {
	colorPrint(colorGreen, "%v", s)
}

func WarningPrint(s string) {
	colorPrint(colorYellow, "Warning: But %v", s)
}

func ErrStr(err error, str string) error {
	return fmt.Errorf("%v;%s", err, str)
}
