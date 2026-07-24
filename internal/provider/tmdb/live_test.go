//go:build live

package tmdb

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLiveSearchAndEpisodes(t *testing.T) {
	key := os.Getenv("TMDB_API_KEY")
	if key == "" {
		t.Skip("TMDB_API_KEY not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	c := New(key, false)
	shows, err := c.Search(ctx, "Breaking Bad")
	if err != nil || len(shows) == 0 {
		t.Fatalf("search: %v (%d results)", err, len(shows))
	}
	for _, s := range shows {
		if s.Kind == "tv" {
			eps, err := c.Episodes(ctx, s.ID)
			if err != nil || len(eps) == 0 {
				t.Fatalf("episodes: %v (%d results)", err, len(eps))
			}
			return
		}
	}
	t.Fatal("no TV result")
}
