package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	// The default template version
	defaultTemplateVersion = "v1.3.0"
)

func main() {
	// Define the URL to get the latest version
	// latestVersionURL := "https://github.com/openimsdk/chat/releases/latest"
	// latestVersion, err := getLatestVersion(latestVersionURL)
	// if err != nil {
	// 	fmt.Printf("Failed to get the latest version: %v\n", err)
	// 	return
	// }
	latestVersion := defaultTemplateVersion

	// Construct the download URL
	downloadURL := fmt.Sprintf("https://github.com/openimsdk/chat/releases/download/%s/chat_Linux_x86_64.tar.gz", latestVersion)

	// Set the installation directory
	installDir := "/tmp/chat"

	// Clear the installation directory before proceeding
	err := os.RemoveAll(installDir)
	if err != nil {
		fmt.Printf("Failed to clear installation directory: %v\n", err)
		return
	}

	// Create the installation directory
	err = os.MkdirAll(installDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create installation directory: %v\n", err)
		return
	}

	// Download and extract OpenIM Chat to the installation directory
	err = downloadAndExtract(downloadURL, installDir)
	if err != nil {
		fmt.Printf("Failed to download and extract OpenIM Chat: %v\n", err)
		return
	}

	// Create configuration file directory
	configDir := filepath.Join(installDir, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create configuration directory: %v\n", err)
		return
	}

	// Download configuration files
	configURL := "https://raw.githubusercontent.com/openimsdk/chat/main/config/config.yaml"
	err = downloadAndExtract(configURL, configDir)
	if err != nil {
		fmt.Printf("Failed to download and extract configuration files: %v\n", err)
		return
	}

	// Define the processes to be started
	cmds := []string{
		"admin-api",
		"admin-rpc",
		"chat-api",
		"chat-rpc",
	}

	// Start each process in a new goroutine
	for _, cmd := range cmds {
		go startProcess(filepath.Join(installDir, cmd))
	}

	// Block the main thread indefinitely
	select {}
}

// getLatestVersion fetches the latest version number from a given URL
func getLatestVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return defaultTemplateVersion, nil
	}

	// Extract the version number from the URL
	latestVersion := filepath.Base(location)
	return latestVersion, nil
}

// downloadAndExtract downloads a file from a URL and extracts it to a destination directory
func downloadAndExtract(url, destDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading file, HTTP status code: %d", resp.StatusCode)
	}

	// Create the destination directory
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// Define the path for the downloaded file
	filePath := filepath.Join(destDir, "downloaded_file.tar.gz")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the downloaded file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	// Extract the file
	cmd := exec.Command("tar", "xzvf", filePath, "-C", destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// startProcess starts a process and prints any errors encountered
func startProcess(cmdPath string) {
	cmd := exec.Command(cmdPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to start process %s: %v\n", cmdPath, err)
	}
}
