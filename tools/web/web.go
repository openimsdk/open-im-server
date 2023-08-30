package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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
	http.Handle("/", fs)

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
