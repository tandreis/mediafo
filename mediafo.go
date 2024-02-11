package mediafo

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveFiles moves files from source directory to destination directory
// and organizes them into subfolders based on creation year and month.
func MoveFiles(sourceDir string, destDir string) error {
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			createTime := info.ModTime()
			year := createTime.Year()
			month := createTime.Month()

			subfolder := filepath.Join(destDir, fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month))

			if err := os.MkdirAll(subfolder, os.ModePerm); err != nil {
				return err
			}

			destFile := filepath.Join(subfolder, info.Name())
			if err := os.Rename(path, destFile); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
