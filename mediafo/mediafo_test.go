package mediafo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

var vfsMap = fstest.MapFS{
	"file1.png": {
		Data:    []byte("file 1"),
		ModTime: time.Date(2010, 10, 1, 0, 0, 0, 0, time.UTC),
		Mode:    0666,
	},
	filepath.Join("dir1", "file2.png"): {
		Data:    []byte("file 2"),
		ModTime: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		Mode:    0666,
	},
	filepath.Join("dir2", "file2.png"): {
		Data:    []byte("file 2"),
		ModTime: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		Mode:    0666,
	},
}

var vfsMapWant = fstest.MapFS{
	filepath.Join("2010", "10", "file1.png"): {
		Data:    []byte("file 1"),
		ModTime: time.Date(2010, 10, 1, 0, 0, 0, 0, time.UTC),
		Mode:    0666,
	},
	filepath.Join("2011", "01", "file2.png"): {
		Data:    []byte("file 2"),
		ModTime: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		Mode:    0666,
	},
}

var vfsMapPtr *fstest.MapFS

type mockOSInterface struct{}

func (o mockOSInterface) Rename(oldpath, newpath string) error {
	fmt.Printf("Mock Rename: %s -> %s\n", oldpath, newpath)
	v := (*vfsMapPtr)[oldpath]
	(*vfsMapPtr)[newpath] = v
	delete((*vfsMapPtr), oldpath)

	return nil
}

func (o mockOSInterface) MkdirAll(path string, perm os.FileMode) error {
	fmt.Printf("Mock MkdirAll: %s (%v)\n", path, perm)
	return nil
}

func prepareVFSMap() {
	vfsMapCopy := fstest.MapFS{}
	for k, v := range vfsMap {
		vfsMapCopy[k] = v
	}
	vfsMapPtr = &vfsMapCopy
}

func TestBasic(t *testing.T) {
	prepareVFSMap()

	moveMedia(*vfsMapPtr, ".", ".", mockOSInterface{})

	if diff := cmp.Diff(vfsMapWant, *vfsMapPtr); diff != "" {
		t.Errorf("moveMedia() mismatch (-want +got):\n%s", diff)
	}
}
