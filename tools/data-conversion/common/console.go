package common

import "fmt"

func ErrorPrint(s string) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", 31, s)
}

func SuccessPrint(s string) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", 32, s)
}

func WarningPrint(s string) {
	fmt.Printf("\x1b[%dmWarning: But %v\x1b[0m\n", 33, s)
}
