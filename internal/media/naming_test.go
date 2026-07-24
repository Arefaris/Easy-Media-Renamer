package media

import (
	"myproject/internal/provider"
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
