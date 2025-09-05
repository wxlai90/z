package z

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var copyFile = io.Copy
var mkdirAll = os.MkdirAll

func (z *Z) SaveUploadedFile(key string, dstPath string) error {
	file, _, err := z.FormFile(key)
	if err != nil {
		return fmt.Errorf("failed to get form file: %w", err)
	}
	defer file.Close()

	if dir := filepath.Dir(dstPath); dir != "." && dir != "" {
		if err := mkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	if _, err := copyFile(out, file); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
