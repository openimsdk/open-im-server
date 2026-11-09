//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/openimsdk/gomake/mageutil"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/utils/datautil"
)

var Default = BuildAll

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

func BuildAll() error { return Build(nil, nil) }

// Build support specifical binary build.
//
// Example: `mage build openim-api openim-rpc-user seq`
func Build(services, tools *string) (err error) {
	defer mageutil.PrintErrPtr(&err)

	return mageutil.WithSpinnerR("Building binaries...", func() error {
		return mageutil.Build(mageutil.ParseArgList(services), mageutil.ParseArgList(tools), nil, nil)
	})
}

func BuildWithCustomConfig(services, tools *string) (err error) {
	defer mageutil.PrintErrPtr(&err)

	config := &mageutil.PathOptions{
		RootDir:   &customRootDir,
		OutputDir: &customOutputDir,
		SrcDir:    &customSrcDir,
		ToolsDir:  &customToolsDir,
	}

	return mageutil.WithSpinnerR("Building binaries with custom config...", func() error {
		return mageutil.Build(mageutil.ParseArgList(services), mageutil.ParseArgList(tools), config, nil)
	})
}

func Start() (err error) {
	defer mageutil.PrintErrPtr(&err)

	return mageutil.WithSpinnerR("Starting...", func() error {
		return mageutil.StartToolsAndServices(nil, nil, nil)
	})
}

func StartWithCustomConfig(services, tools *string) (err error) {
	defer mageutil.PrintErrPtr(&err)

	config := &mageutil.PathOptions{
		RootDir:   &customRootDir,
		OutputDir: &customOutputDir,
		ConfigDir: &customConfigDir,
	}

	return mageutil.WithSpinnerR("Starting with custom config...", func() error {
		return mageutil.StartToolsAndServices(mageutil.ParseArgList(tools), mageutil.ParseArgList(services), config)
	})
}

func Stop() (err error) {
	defer mageutil.PrintErrPtr(&err)
	return mageutil.WithSpinnerR("Stopping...", mageutil.StopAndCheckServices)
}

func Check() (err error) {
	defer mageutil.PrintErrPtr(&err)
	return mageutil.WithSpinnerR("Checking binaries...", mageutil.CheckAndReportServicesStatus)
}

func Export() (err error) {
	defer mageutil.PrintErrPtr(&err)

	mappingPaths, err := mageutil.GetDefaultExportMappingPaths([]string{
		"cmd",
		"internal",
		"pkg",
		"test",
		"tools",
		"**/*.go",
		"go.mod",
		"go.sum",
		"go.work",
	})
	if err != nil {
		return fmt.Errorf("GetDefaultExportMappingPaths failed: %v", err)
	}

	return mageutil.WithSpinnerR("Exporting...", func() error {
		return mageutil.ExportMageLauncherArchived(mappingPaths, &mageutil.ExportOptions{
			ProjectName: datautil.ToPtr(fmt.Sprintf("open-im-server_%s", version.Version)),
			BuildOpt: &mageutil.BuildOptions{
				Release:  datautil.ToPtr(true),
				Compress: datautil.ToPtr(true),
			},
		})
	})
}
