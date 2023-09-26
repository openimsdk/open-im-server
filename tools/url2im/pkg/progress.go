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
