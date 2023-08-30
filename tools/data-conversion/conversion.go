package main

import (
	"fmt"
	. "github.com/OpenIMSDK/Open-IM-Server/tools/conversion/common"
	message "github.com/OpenIMSDK/Open-IM-Server/tools/conversion/msg"
	"github.com/OpenIMSDK/Open-IM-Server/tools/conversion/mysql"
	"sync"
)

var wg sync.WaitGroup

func main() {
	fmt.Printf("start MySQL data conversion. \n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.UserConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.FriendConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.GroupConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.GroupMemberConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.BlacksConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.RequestConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.BlacksConversion()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mysql.ChatLogsConversion()
	}()
	wg.Wait()
	SuccessPrint(fmt.Sprintf("Successfully completed the MySQL conversion. \n"))

	fmt.Printf("start message conversion. \n")
	message.GetMessage()
	SuccessPrint(fmt.Sprintf("Successfully completed the message conversion. \n"))
}
