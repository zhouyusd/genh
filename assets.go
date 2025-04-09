package genh

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/imports"
)

type Assets struct {
	dirs  map[string]struct{}
	files map[string][]byte
}

func (a *Assets) Add(path string, content []byte) {
	if a.files == nil {
		a.files = make(map[string][]byte)
	}
	a.files[path] = content
}

func (a *Assets) AddDir(path string) {
	if a.dirs == nil {
		a.dirs = make(map[string]struct{})
	}
	a.dirs[path] = struct{}{}
}

// write files and dirs in the assets.
func (a *Assets) Write() error {
	for dir := range a.dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}
	for path, content := range a.files {
		if err := os.WriteFile(path, content, 0644); err != nil {
			return fmt.Errorf("write file %q: %w", path, err)
		}
	}
	return nil
}

// Format runs "goimports" on all assets.
func (a *Assets) Format() error {
	var wg errgroup.Group
	wg.SetLimit(runtime.GOMAXPROCS(0))
	for path, content := range a.files {
		path, content := path, content
		wg.Go(func() error {
			src, err := imports.Process(path, content, nil)
			if err != nil {
				return fmt.Errorf("format file %s: %w", path, err)
			}
			if err := os.WriteFile(path, src, 0644); err != nil {
				return fmt.Errorf("write file %s: %w", path, err)
			}
			return nil
		})
	}
	return wg.Wait()
}
