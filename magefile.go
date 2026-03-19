//go:build mage
// +build mage

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openimsdk/gomake/mageutil"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/utils/datautil"
)

var Default = Build

var Aliases = map[string]any{
	"buildcc": BuildWithCustomConfig,
	"startcc": StartWithCustomConfig,
}

var (
	customRootDir   = "."       // workDir in mage, default is "./"(project root directory)
	customSrcDir    = "cmd"     // source code directory, default is "cmd"
	customOutputDir = "_output" // output directory, default is "_output"
	customConfigDir = "config"  // configuration directory, default is "config"
	customToolsDir  = "tools"   // tools source code directory, default is "tools"
)

// Build support specifical binary build.
//
// Example: `mage build openim-api openim-rpc-user seq`
func Build() {
	flag.Parse()
	bin := flag.Args()
	if len(bin) != 0 {
		bin = bin[1:]
	}
	mageutil.WithSpinner("Building binaries...", func() { mageutil.Build(bin, nil, nil) })
}

func BuildWithCustomConfig() {
	flag.Parse()
	bin := flag.Args()
	if len(bin) != 0 {
		bin = bin[1:]
	}

	config := &mageutil.PathOptions{
		RootDir:   &customRootDir,
		OutputDir: &customOutputDir,
		SrcDir:    &customSrcDir,
		ToolsDir:  &customToolsDir,
	}

	mageutil.WithSpinner("Building binaries with custom config...", func() {
		mageutil.Build(bin, config, nil)
	})
}

func Start() {
	mageutil.InitForSSC()
	err := setMaxOpenFiles()
	if err != nil {
		mageutil.PrintRed("setMaxOpenFiles failed " + err.Error())
		os.Exit(1)
	}

	flag.Parse()
	bin := flag.Args()
	if len(bin) != 0 {
		bin = bin[1:]
	}

	mageutil.WithSpinner("Starting...", func() {
		mageutil.StartToolsAndServices(bin, nil)
	})
}

func StartWithCustomConfig() {
	mageutil.InitForSSC()
	err := setMaxOpenFiles()
	if err != nil {
		mageutil.PrintRed("setMaxOpenFiles failed " + err.Error())
		os.Exit(1)
	}

	flag.Parse()
	bin := flag.Args()
	if len(bin) != 0 {
		bin = bin[1:]
	}

	config := &mageutil.PathOptions{
		RootDir:   &customRootDir,
		OutputDir: &customOutputDir,
		ConfigDir: &customConfigDir,
	}

	mageutil.WithSpinner("Starting with custom config...", func() {
		mageutil.StartToolsAndServices(bin, config)
	})
}

func Stop() {
	mageutil.WithSpinner("Stopping...", mageutil.StopAndCheckBinaries)
}

func Check() {
	mageutil.WithSpinner("Checking binaries...", mageutil.CheckAndReportBinariesStatus)
}

func Export() {
	mappingPaths, err := mageutil.GetDefaultExportMappingPaths([]string{
		"cmd",
		"internal",
		"pkg",
		"test",
		"tools",
		"**/*.go",
		"go.mod",
		"go.work",
	})
	if err != nil {
		mageutil.PrintRed("GetDefaultExportMappingPaths failed " + err.Error())
		os.Exit(1)
	}

	mageutil.WithSpinner("Exporting...", func() {
		mageutil.ExportMageLauncherArchived(mappingPaths, &mageutil.ExportOptions{
			ProjectName: datautil.ToPtr(fmt.Sprintf("open-im-server_%s", version.Version)),
			BuildOpt: &mageutil.BuildOptions{
				Release:  datautil.ToPtr(true),
				Compress: datautil.ToPtr(true),
			},
		})
	})
}
