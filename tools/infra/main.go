package main

import (
	"fmt"

	"github.com/fatih/color"
)

// Define a function to print important link information
func printLinks() {
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("OpenIM Github: %s\n", blue("https://github.com/OpenIMSDK/Open-IM-Server"))
	fmt.Printf("Slack Invitation: %s\n", blue("https://openimsdk.slack.com"))
	fmt.Printf("Follow Twitter: %s\n", blue("https://twitter.com/founder_im63606"))
}

func main() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue, color.Bold)

	yellow.Println("Please use the release branch or tag for production environments!")

	message := `
____                       _____  __  __ 
/ __ \                     |_   _||  \/  |
| |  | | _ __    ___  _ __    | |  | \  / |
| |  | || '_ \  / _ \| '_ \   | |  | |\/| |
| |__| || |_) ||  __/| | | | _| |_ | |  | |
\____/ | .__/  \___||_| |_||_____||_|  |_|
	   | |                                
	   |_|                                

Keep checking for updates!
`

	blue.Println(message)
	printLinks() // Call the function to print the link information
}
