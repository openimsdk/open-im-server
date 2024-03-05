// Copyright © 2024 OpenIM. All rights reserved.
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

package checker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/openimsdk/open-im-server/tools/codescan/config"
)

type CheckResult struct {
	FilePath string
	Lines    []int
}

func checkFileForChineseComments(filePath string) ([]CheckResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []CheckResult
	scanner := bufio.NewScanner(file)
	reg := regexp.MustCompile(`[\p{Han}]+`)
	lineNumber := 0

	var linesWithChinese []int
	for scanner.Scan() {
		lineNumber++
		if reg.FindString(scanner.Text()) != "" {
			linesWithChinese = append(linesWithChinese, lineNumber)
		}
	}

	if len(linesWithChinese) > 0 {
		results = append(results, CheckResult{
			FilePath: filePath,
			Lines:    linesWithChinese,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func WalkDirAndCheckComments(cfg config.Config) error {
	var allResults []CheckResult
	err := filepath.Walk(cfg.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, fileType := range cfg.FileTypes {
			if filepath.Ext(path) == fileType {
				results, err := checkFileForChineseComments(path)
				if err != nil {
					return err
				}
				if len(results) > 0 {
					allResults = append(allResults, results...)
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(allResults) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString("Files containing Chinese comments:\n")
		for _, result := range allResults {
			errMsg.WriteString(fmt.Sprintf("%s: Lines %v\n", result.FilePath, result.Lines))
		}
		return fmt.Errorf(errMsg.String())
	}

	return nil
}
