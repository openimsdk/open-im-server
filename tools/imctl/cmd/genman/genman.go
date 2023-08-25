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
	"os"

	"k8s.io/kubernetes/cmd/genutils"
)

func main() {
	// TODO use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/man/man1"
	module := ""
	if len(os.Args) == 3 {
		path = os.Args[1]
		module = os.Args[2]
	} else {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory] [module] \n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by command so the output is consistent,
	// regardless of where we run.
	os.Setenv("HOME", "/home/username")

	// openim-api
	// openim-cmdutils
	// openim-crontask
	// openim-msggateway
	// openim-msgtransfer
	// openim-push
	// openim-rpc-auth
	// openim-rpc-conversation
	// openim-rpc-friend
	// openim-rpc-group
	// openim-rpc-msg
	// openim-rpc-third
	// openim-rpc-user
	switch module {
	case "openim-api":
	}
}
