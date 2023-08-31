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
