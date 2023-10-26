package main

import (
	"fmt"

	"github.com/fatih/color"
)

// 定义一个函数以打印重要的链接信息
func printLinks() {
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("OpenIM Github: %s\n", blue("https://github.com/OpenIMSDK/Open-IM-Server"))
	fmt.Printf("Slack Invitation: %s\n", blue("https://openimsdk.slack.com"))
}

func main() {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue, color.Bold)

	yellow.Println("Current module is still under development.")

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
	printLinks() // 调用函数以打印链接信息
}
