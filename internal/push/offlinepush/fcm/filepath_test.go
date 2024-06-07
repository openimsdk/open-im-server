package fcm

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	fmt.Println(filepath.Join("a/b/", "a.json"))
}
