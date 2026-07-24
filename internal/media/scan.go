package media

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	RelPath string `json:"rel_path"`
	Size    int64  `json:"size"`
}

func ScanDirectory(dir string, extensions []string) ([]File, error) {
	return Scan(dir, extensions, false, 0)
}
func Scan(dir string, extensions []string, recursive bool, maxDepth int) ([]File, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	allowed := map[string]bool{}
	for _, e := range extensions {
		allowed[strings.ToLower(e)] = true
	}
	out := []File{}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		if rel == "." {
			return nil
		}
		depth := strings.Count(rel, string(filepath.Separator)) + 1
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") || isSystemDir(d.Name()) || !recursive || maxDepth > 0 && depth > maxDepth {
				return filepath.SkipDir
			}
			return nil
		}
		if maxDepth > 0 && depth > maxDepth {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") || !allowed[strings.ToLower(filepath.Ext(d.Name()))] {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		out = append(out, File{Name: d.Name(), Path: path, RelPath: rel, Size: info.Size()})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan directory: %w", err)
	}
	NaturalSort(out, func(f File) string { return f.RelPath })
	return out, nil
}
func isSystemDir(name string) bool {
	switch strings.ToLower(name) {
	case "system volume information", "$recycle.bin", "node_modules", ".git":
		return true
	}
	return false
}
func EnsureDirectory(path string) error { return os.MkdirAll(path, 0755) }
