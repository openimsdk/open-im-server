package mageutil

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CheckAndReportBinariesStatus checks the running status of all binary files and reports it.
func CheckAndReportBinariesStatus() {
	InitForSSC()
	err := CheckBinariesRunning()
	if err != nil {
		PrintRed("Some programs are not running properly:")
		PrintRedNoTimeStamp(err.Error())
		return
	}
	PrintGreen("All services are running normally.")
	PrintBlue("Display details of the ports listened to by the service:")
	PrintListenedPortsByBinaries()
}

// StopAndCheckBinaries stops all binary processes and checks if they have all stopped.
func StopAndCheckBinaries() {
	InitForSSC()
	KillExistBinaries()
	err := CheckBinariesStop()
	if err != nil {
		PrintRed("Some services have not been stopped, details are as follows:" + err.Error())
		return
	}
	PrintGreen("All services have been stopped")
}

// StartToolsAndServices starts the process for tools and services.
func StartToolsAndServices() {
	InitForSSC()
	PrintBlue("Starting tools primarily involves component verification and other preparatory tasks.")
	if err := StartTools(); err != nil {
		PrintRed("Some tools failed to start, details are as follows, abort start")
		PrintRedNoTimeStamp(err.Error())
		return
	}
	PrintGreen("All tools executed successfully")

	KillExistBinaries()
	err := CheckBinariesStop()
	if err != nil {
		PrintRed("Some services running, details are as follows, abort start " + err.Error())
		return
	}
	PrintBlue("Starting services involves multiple RPCs and APIs and may take some time. Please be patient")
	err = StartBinaries()
	if err != nil {
		PrintRed("Failed to start all binaries")
		PrintRedNoTimeStamp(err.Error())
		return
	}
	CheckAndReportBinariesStatus()
}

// CompileForPlatform Main compile function
func CompileForPlatform(platform string) {

	PrintBlue(fmt.Sprintf("Compiling cmd for %s...", platform))

	cmdCompiledDirs := compileDir(filepath.Join(rootDirPath, "cmd"), platformsOutputBase, platform)

	PrintBlue(fmt.Sprintf("Compiling tools for %s...", platform))
	toolsCompiledDirs := compileDir(filepath.Join(rootDirPath, "tools"), toolsOutputBase, platform)
	createStartConfigYML(cmdCompiledDirs, toolsCompiledDirs)

}

func createStartConfigYML(cmdDirs, toolsDirs []string) {
	configPath := filepath.Join(rootDirPath, "start-config.yml")

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		PrintBlue("start-config.yml already exists, skipping creation.")
		return
	}

	var content strings.Builder
	content.WriteString("serviceBinaries:\n")
	for _, dir := range cmdDirs {
		content.WriteString(fmt.Sprintf("  %s: 1\n", dir))
	}
	content.WriteString("toolBinaries:\n")
	for _, dir := range toolsDirs {
		content.WriteString(fmt.Sprintf("  - %s\n", dir))
	}
	content.WriteString("maxFileDescriptors: 10000\n")

	err := ioutil.WriteFile(configPath, []byte(content.String()), 0644)
	if err != nil {
		PrintRed("Failed to create start-config.yml: " + err.Error())
		return
	}
	PrintGreen("start-config.yml created successfully.")
}

func compileDir(sourceDir, outputBase, platform string) []string {
	var compiledDirs []string
	var mu sync.Mutex
	targetOS, targetArch := strings.Split(platform, "_")[0], strings.Split(platform, "_")[1]
	outputDir := filepath.Join(outputBase, targetOS, targetArch)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	errors := make(chan error, 1)
	sem := make(chan struct{}, 4)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Base(path) != "main.go" {
			return nil
		}

		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer wg.Done()
			defer func() { <-sem }()

			dir := filepath.Dir(path)
			dirName := filepath.Base(dir)
			outputFileName := dirName
			if targetOS == "windows" {
				outputFileName += ".exe"
			}

			PrintBlue(fmt.Sprintf("Compiling dir: %s for platform: %s binary: %s ...", dirName, platform, outputFileName))
			err := sh.RunWith(map[string]string{"GOOS": targetOS, "GOARCH": targetArch}, "go", "build", "-o", filepath.Join(outputDir, outputFileName), filepath.Join(dir, "main.go"))
			if err != nil {
				errors <- fmt.Errorf("failed to compile %s for %s: %v", dirName, platform, err)
				PrintRed("Compilation aborted. " + fmt.Sprintf("failed to compile %s for %s: %v", dirName, platform, err))
				os.Exit(1)
				return
			}
			PrintGreen(fmt.Sprintf("Successfully compiled. dir: %s for platform: %s binary: %s", dirName, platform, outputFileName))
			mu.Lock()
			compiledDirs = append(compiledDirs, dirName)
			mu.Unlock()
		}()

		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directories:", err)
		os.Exit(1)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	if err, ok := <-errors; ok {
		fmt.Println(err)
		os.Exit(1)
	}
	return compiledDirs
}

// compileDir compiles Go programs in a specified directory, appending .exe extension for output files on Windows platform
//func compileDir(sourceDir, outputBase, platform string) {
//	targetOS, targetArch := strings.Split(platform, "_")[0], strings.Split(platform, "_")[1]
//	outputDir := filepath.Join(outputBase, targetOS, targetArch)
//
//	if err := os.MkdirAll(outputDir, 0755); err != nil {
//		fmt.Printf("Failed to create directory %s: %v\n", outputDir, err)
//		os.Exit(1)
//	}
//
//	var wg sync.WaitGroup
//	errors := make(chan error, 1)
//
//	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		if info.IsDir() || filepath.Base(path) != "main.go" {
//			return nil
//		}
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//
//			dir := filepath.Dir(path)
//			dirName := filepath.Base(dir)
//			outputFileName := dirName
//			if targetOS == "windows" {
//				outputFileName += ".exe"
//			}
//
//			PrintBlue(fmt.Sprintf("Compiling dir: %s for platform: %s binary: %s ...", dirName, platform, outputFileName))
//			err := sh.RunWith(map[string]string{"GOOS": targetOS, "GOARCH": targetArch}, "go", "build", "-o", filepath.Join(outputDir, outputFileName), filepath.Join(dir, "main.go"))
//			if err != nil {
//				errors <- fmt.Errorf("failed to compile %s for %s: %v", dirName, platform, err)
//				PrintRed("Compilation aborted. " + fmt.Sprintf("failed to compile %s for %s: %v", dirName, platform, err))
//				os.Exit(1)
//				return
//			}
//			PrintGreen(fmt.Sprintf("Compiling dir: %s for platform: %s binary: %s ...", dirName, platform, outputFileName))
//		}()
//
//		return nil
//	})
//
//	if err != nil {
//		fmt.Println("Error walking through directories:", err)
//		os.Exit(1)
//	}
//
//	wg.Wait()
//	close(errors)
//
//	// Check for errors
//	if err, ok := <-errors; ok {
//		fmt.Println(err)
//		fmt.Println("Compilation aborted.")
//		os.Exit(1)
//	}
//}
