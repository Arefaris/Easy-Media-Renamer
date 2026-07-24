package media

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func ScanDirectory(dir string, extensions []string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("scan directory: %w", err)
	}
	allowed := make(map[string]bool, len(extensions))
	for _, e := range extensions {
		allowed[strings.ToLower(e)] = true
	}
	out := make([]File, 0)
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") || !allowed[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("inspect %s: %w", e.Name(), err)
		}
		out = append(out, File{Name: e.Name(), Path: filepath.Join(dir, e.Name()), Size: info.Size()})
	}
	NaturalSort(out, func(f File) string { return f.Name })
	return out, nil
}
