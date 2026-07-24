package media

import (
	matcher "myproject/internal/match"
	"myproject/internal/provider"
	"os"
	"path/filepath"
	"strings"
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

func TestPlanAndApplyPreserveExtensionAfterNumericTitle(t *testing.T) {
	d := t.TempDir()
	t.Setenv("EMR_CACHE_DIR", filepath.Join(d, "cache"))
	mustWrite(t, filepath.Join(d, "a.mkv"), "video")
	ops, err := Plan(d, []string{"a.mkv"}, []provider.Episode{{Season: 1, Number: 1, Title: "Chapter 2.5"}}, provider.Show{Name: "Show"}, "{n} - {s00e00} - {t}")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Apply(ops); err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(ops[0].To) != ".mkv" {
		t.Fatalf("lost extension: %s", ops[0].To)
	}
	if _, err = os.Stat(ops[0].To); err != nil {
		t.Fatal(err)
	}
}

func TestPlanMatchedKeepsUnmatchedFileVisibleWithoutShifting(t *testing.T) {
	d := t.TempDir()
	files := []string{"Show.S01E03.PROPER.mkv", "Show.S01E03.mkv", "Show.S01E01.mkv"}
	for _, name := range files {
		mustWrite(t, filepath.Join(d, name), name)
	}
	eps := []provider.Episode{{Season: 1, Number: 1, Title: "First"}, {Season: 1, Number: 2, Title: "Second"}, {Season: 1, Number: 3, Title: "Third"}}
	pairs := []matcher.Pair{{FileIndex: 0, EpisodeIndex: 2, Confidence: .99, Strategy: "season+episode"}, {FileIndex: 1, EpisodeIndex: -1, Strategy: "unmatched"}, {FileIndex: 2, EpisodeIndex: 0, Confidence: .99, Strategy: "season+episode"}}
	ops, err := PlanMatched(d, files, eps, provider.Show{Name: "Show"}, "{n} - {s00e00} - {t}", pairs)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 3 || ops[1].Status != "unmatched" || !strings.Contains(ops[0].To, "S01E03") || !strings.Contains(ops[2].To, "S01E01") {
		t.Fatalf("%#v", ops)
	}
}

func TestPlanWithoutEpisodesKeepsEveryFileVisible(t *testing.T) {
	d := t.TempDir()
	files := []string{"one.mkv", "two.mkv"}
	ops, err := Plan(d, files, nil, provider.Show{Name: "Show"}, "{n}")
	if err != nil || len(ops) != len(files) || ops[0].Status != "unmatched" || ops[1].Status != "unmatched" {
		t.Fatalf("err=%v ops=%#v", err, ops)
	}
}

func TestPlanIgnoresUnreadableMediaInfo(t *testing.T) {
	d := t.TempDir()
	mustWrite(t, filepath.Join(d, "bad.mkv"), "not an mkv")
	ops, err := Plan(d, []string{"bad.mkv"}, []provider.Episode{{Season: 1, Number: 1, Title: "Pilot"}}, provider.Show{Name: "Show"}, "{n} - {s00e00}[ - {vf}]")
	if err != nil || len(ops) != 1 {
		t.Fatalf("err=%v ops=%#v", err, ops)
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

func TestCopyAndHardlinkCreateDirectories(t *testing.T) {
	for _, action := range []Action{ActionCopy, ActionHardlink} {
		t.Run(string(action), func(t *testing.T) {
			d := t.TempDir()
			t.Setenv("EMR_CACHE_DIR", filepath.Join(d, "cache"))
			from := filepath.Join(d, "source.mkv")
			to := filepath.Join(d, "nested", "target.mkv")
			mustWrite(t, from, "content")
			r, e := ApplyAction([]Op{{From: from, To: to, Status: "ok"}}, action)
			if e != nil || r.Renamed != 1 {
				t.Fatalf("%v %#v", e, r)
			}
			if string(mustRead(t, to)) != "content" {
				t.Fatal("bad target")
			}
			if _, e = os.Stat(from); e != nil {
				t.Fatal("source removed")
			}
		})
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
