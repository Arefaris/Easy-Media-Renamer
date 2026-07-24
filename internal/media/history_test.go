package media

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRevertArbitraryHistory(t *testing.T) {
	d := t.TempDir()
	t.Setenv("EMR_CACHE_DIR", filepath.Join(d, "cache"))
	a := filepath.Join(d, "a.txt")
	b := filepath.Join(d, "b.txt")
	mustWrite(t, a, "a")
	if _, e := ApplyAction([]Op{{From: a, To: b, Status: "ok"}}, ActionCopy); e != nil {
		t.Fatal(e)
	}
	list, e := ListHistory()
	if e != nil || len(list) != 1 {
		t.Fatalf("%v %#v", e, list)
	}
	if _, e = RevertHistory(list[0].ID); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(b); !os.IsNotExist(e) {
		t.Fatal("copy target still exists")
	}
	if _, e = os.Stat(a); e != nil {
		t.Fatal("source missing")
	}
}
