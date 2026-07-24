package media

import (
	"myproject/internal/provider"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	for _, tc := range []struct{ in, want string }{{`a<b>:c?.mkv`, `a-b--c-.mkv`}, {`CON`, `_CON`}, {`title. `, `title`}} {
		if got := Sanitize(tc.in); got != tc.want {
			t.Errorf("Sanitize(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
	if got := Sanitize(strings.Repeat("я", 200)); len([]byte(got)) > 255 {
		t.Fatalf("too long: %d", len([]byte(got)))
	}
}
func TestRender(t *testing.T) {
	got := Render("{show} - s{season}e{episode} - {title}", provider.Episode{Season: 1, Number: 2, Title: "Pilot"}, provider.Show{Name: "Demo"}, ".mkv")
	if got != "Demo - s01e02 - Pilot.mkv" {
		t.Fatal(got)
	}
}

func TestRenderTemplateConditionalAndDirectories(t *testing.T) {
	d := NamingData{
		Episode: provider.Episode{Season: 2, Number: 3, Title: "A: Title"},
		Show:    provider.Show{Name: "Demo", Year: "2024"}, Ext: ".mkv",
	}
	got := RenderTemplate("{n} ({y})/Season {s}/{n} - {s00e00} - {t}[ - {group}]", d)
	want := filepath.Join("Demo (2024)", "Season 2", "Demo - S02E03 - A- Title.mkv")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	d.Group = "GROUP"
	if got := RenderTemplate("{n}[ - {group}]", d); got != "Demo - GROUP.mkv" {
		t.Fatal(got)
	}
}

func TestRenderTemplateSanitizesEverySegment(t *testing.T) {
	d := NamingData{Episode: provider.Episode{Season: 1, Number: 1}, Show: provider.Show{Name: "CON"}, Ext: ".mkv"}
	got := RenderTemplate("{n}/AUX/{s00e00}", d)
	want := filepath.Join("_CON", "_AUX", "S01E01.mkv")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderTemplateAlwaysPreservesKnownExtension(t *testing.T) {
	d := NamingData{Episode: provider.Episode{Season: 1, Number: 1, Title: "Chapter 2.5"}, Show: provider.Show{Name: "Show"}, Ext: ".mkv"}
	if got := RenderTemplate("{n} - {s00e00} - {t}", d); got != "Show - S01E01 - Chapter 2.5.mkv" {
		t.Fatal(got)
	}
	if got := RenderTemplate("{n}.{ext}", d); got != "Show.mkv" {
		t.Fatal(got)
	}
}

func TestTokenSlashCannotCreateDirectory(t *testing.T) {
	d := NamingData{Episode: provider.Episode{Season: 1, Number: 1, Title: "Us/Them"}, Show: provider.Show{Name: "Show"}, Ext: ".mkv"}
	if got := RenderTemplate("{n} - {t}", d); got != "Show - Us-Them.mkv" {
		t.Fatal(got)
	}
}

func TestEscapedLiteralBrackets(t *testing.T) {
	d := NamingData{Episode: provider.Episode{Number: 5}, Show: provider.Show{Name: "Frieren"}, Group: "SubsPlease", Ext: ".mkv"}
	if got := RenderTemplate(`\[{group}\] {n} - {e00}`, d); got != "[SubsPlease] Frieren - 05.mkv" {
		t.Fatal(got)
	}
}
