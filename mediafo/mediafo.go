package mediafo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type osInterface interface {
	MkdirAll(path string, perm os.FileMode) error
	Rename(oldpath, newpath string) error
}

type realOSInterface struct{}

func (o realOSInterface) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (o realOSInterface) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// MoveMedia moves media files from source directory to destination directory
// and organizes them into subfolders based on media creation year and month.
func MoveMedia(sourceDir string, destDir string) error {
	fs := os.DirFS(sourceDir)
	return moveMedia(fs, sourceDir, destDir, realOSInterface{})
}

func moveMedia(fileSystem fs.FS, sourceDir string, destDir string, osi osInterface) error {
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		createTime := info.ModTime()
		year := createTime.Year()
		month := createTime.Month()

		subfolder := filepath.Join(destDir, fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month))

		if err := osi.MkdirAll(subfolder, os.ModePerm); err != nil {
			return err
		}

		sourceFile := filepath.Join(sourceDir, path)
		destFile := filepath.Join(subfolder, d.Name())
		if err := osi.Rename(sourceFile, destFile); err != nil {
			fmt.Fprintf(os.Stderr, "rename: %v\n", err)
			return err
		}

		return err
	})

	return err
}
