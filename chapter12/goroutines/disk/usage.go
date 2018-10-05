package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

type SizeInfo struct {
	Name string
	Size int64
}

func FormatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%4dB", size)
	}

	sizeFraction := float64(size) / 1024
	if sizeFraction < 1024 {
		return fmt.Sprintf("% 7.2fKB", sizeFraction)
	}

	sizeFraction = sizeFraction / 1024
	if sizeFraction < 1024 {
		return fmt.Sprintf("% 7.2fMB", sizeFraction)
	}

	sizeFraction = sizeFraction / 1024
	return fmt.Sprintf("% 7.2fGB", sizeFraction)
}

func DirSize(path string, result chan SizeInfo) {
	var size int64

	filepath.Walk(path, func(_ string, file os.FileInfo, err error) error {
		if err == nil {
			size += file.Size()
		}

		return nil
	})

	result <- SizeInfo{filepath.Base(path), size}
}
