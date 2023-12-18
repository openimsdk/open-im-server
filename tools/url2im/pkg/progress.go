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

package pkg

import (
	"bufio"
	"os"
	"strconv"

	"github.com/kelindar/bitmap"
)

func ReadProgress(path string) (*Progress, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Progress{}, nil
		}
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var upload bitmap.Bitmap
	for scanner.Scan() {
		index, err := strconv.Atoi(scanner.Text())
		if err != nil || index < 0 {
			continue
		}
		upload.Set(uint32(index))
	}
	return &Progress{upload: upload}, nil
}

type Progress struct {
	upload bitmap.Bitmap
}

func (p *Progress) IsUploaded(index int) bool {
	if p == nil {
		return false
	}
	return p.upload.Contains(uint32(index))
}
