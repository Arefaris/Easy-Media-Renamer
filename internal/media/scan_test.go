package media

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectoryFiltersAndSorts(t *testing.T) {
	d := t.TempDir()
	for _, n := range []string{"ep10.mkv", "ep2.mkv", "ep1.mkv", "readme.txt", ".hidden.mkv"} {
		mustWrite(t, filepath.Join(d, n), n)
	}
	if err := os.Mkdir(filepath.Join(d, "sub.mkv"), 0700); err != nil {
		t.Fatal(err)
	}
	files, err := ScanDirectory(d, []string{".mkv"})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 || files[0].Name != "ep1.mkv" || files[1].Name != "ep2.mkv" || files[2].Name != "ep10.mkv" {
		t.Fatalf("%#v", files)
	}
}

func TestScanRecursiveDepthAndHiddenDirectories(t *testing.T) {
	d := t.TempDir()
	mustWrite(t, filepath.Join(d, "root.mkv"), "root")
	for _, rel := range []string{"Season 1/ep1.mkv", "Season 1/deep/ep2.mkv", ".hidden/nope.mkv"} {
		p := filepath.Join(d, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, p, rel)
	}
	files, err := Scan(d, []string{".mkv"}, true, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || files[1].RelPath != filepath.Join("Season 1", "ep1.mkv") {
		t.Fatalf("%#v", files)
	}
}
