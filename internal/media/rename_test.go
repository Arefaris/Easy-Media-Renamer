package media

import (
	"myproject/internal/provider"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanConflicts(t *testing.T) {
	d := t.TempDir()
	mustWrite(t, filepath.Join(d, "one.mkv"), "1")
	mustWrite(t, filepath.Join(d, "occupied.mkv"), "x")
	ops, e := Plan(d, []string{"one.mkv"}, []provider.Episode{{Season: 1, Number: 1, Title: "occupied"}}, provider.Show{}, "{title}")
	if e != nil || len(ops) != 1 || ops[0].Status != "conflict-existing" {
		t.Fatalf("%v %#v", e, ops)
	}
}
func TestApplyCycle(t *testing.T) {
	d := t.TempDir()
	t.Setenv("EMR_CACHE_DIR", filepath.Join(d, "cache"))
	a, b := filepath.Join(d, "a.mkv"), filepath.Join(d, "b.mkv")
	mustWrite(t, a, "A")
	mustWrite(t, b, "B")
	r, e := Apply([]Op{{From: a, To: b, Status: "ok"}, {From: b, To: a, Status: "ok"}})
	if e != nil || r.Renamed != 2 {
		t.Fatalf("%v %#v", e, r)
	}
	if string(mustRead(t, a)) != "B" || string(mustRead(t, b)) != "A" {
		t.Fatal("cycle contents incorrect")
	}
	r, e = Undo()
	if e != nil || r.Renamed != 2 {
		t.Fatalf("undo: %v %#v", e, r)
	}
	if string(mustRead(t, a)) != "A" || string(mustRead(t, b)) != "B" {
		t.Fatal("undo contents incorrect")
	}
}
func mustWrite(t *testing.T, p, s string) {
	t.Helper()
	if e := os.WriteFile(p, []byte(s), 0600); e != nil {
		t.Fatal(e)
	}
}
func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
