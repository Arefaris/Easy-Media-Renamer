package match

import (
	"myproject/internal/provider"
	"testing"
)

func TestStrategies(t *testing.T) {
	eps := []provider.Episode{{Season: 1, Number: 1, Title: "One", AirDate: "2024-01-01", Absolute: 1}, {Season: 1, Number: 2, Title: "Two", AirDate: "2024-01-02", Absolute: 2}}
	tests := []struct {
		p    Parsed
		want string
	}{{Parsed{Season: 1, Episodes: []int{2}}, "season+episode"}, {Parsed{Season: -1, Absolute: 2}, "absolute"}, {Parsed{Season: -1, Absolute: -1, AirDate: "2024-01-01"}, "airdate"}, {Parsed{Season: -1, Absolute: -1, Episodes: []int{2}}, "episode-only"}, {Parsed{Season: -1, Absolute: -1}, "position"}}
	for _, tc := range tests {
		got := Match([]Parsed{tc.p}, eps)
		if got[0].Strategy != tc.want {
			t.Errorf("got %s want %s", got[0].Strategy, tc.want)
		}
	}
}
func TestStrategyPriority(t *testing.T) {
	eps := []provider.Episode{{Season: 1, Number: 1, Absolute: 2}, {Season: 1, Number: 2, Absolute: 1}}
	got := Match([]Parsed{{Season: 1, Episodes: []int{1}, Absolute: 1}}, eps)
	if got[0].EpisodeIndex != 0 || got[0].Strategy != "season+episode" {
		t.Fatalf("%#v", got)
	}
}

func TestPositionFallbackDoesNotReuseEpisode(t *testing.T) {
	eps := []provider.Episode{{Season: 1, Number: 1}, {Season: 1, Number: 2}}
	files := []Parsed{{Season: 1, Episodes: []int{2}}, {Season: -1, Absolute: -1}}
	got := Match(files, eps)
	if got[0].EpisodeIndex != 1 || got[1].EpisodeIndex != 0 {
		t.Fatalf("duplicate or missing fallback: %#v", got)
	}
}
