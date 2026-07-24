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
