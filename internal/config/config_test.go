package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	d := t.TempDir()
	t.Setenv("APPDATA", d)
	t.Setenv("TMDB_API_KEY", "")
	want := Default()
	want.TMDBAPIKey = "secret"
	want.NameTemplate = "{title}"
	if err := Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.TMDBAPIKey != want.TMDBAPIKey || got.NameTemplate != want.NameTemplate {
		t.Fatalf("got %#v", got)
	}
	p, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(p) != filepath.Join(d, appDir) {
		t.Fatalf("unexpected path %s", p)
	}
}
