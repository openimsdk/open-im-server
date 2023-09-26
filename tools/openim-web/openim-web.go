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
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
)

var (
	distPathFlag string
	portFlag     string
)

func init() {
	flag.StringVar(&distPathFlag, "distPath", "/app/dist", "Path to the distribution")
	flag.StringVar(&portFlag, "port", "11001", "Port to run the server on")
}

func main() {
	flag.Parse()

	distPath := getConfigValue("DIST_PATH", distPathFlag, "/app/dist")
	fs := http.FileServer(http.Dir(distPath))

	withGzip := gziphandler.GzipHandler(fs)

	http.Handle("/", withGzip)

	port := getConfigValue("PORT", portFlag, "11001")
	log.Printf("Server listening on port %s in %s...", port, distPath)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigValue(envKey, flagValue, fallback string) string {
	envVal := os.Getenv(envKey)
	if envVal != "" {
		return envVal
	}
	if flagValue != "" {
		return flagValue
	}
	return fallback
}
