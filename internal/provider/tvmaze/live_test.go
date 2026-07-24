//go:build live

package tvmaze

import (
	"context"
	"testing"
	"time"
)

func TestLiveSearchAndEpisodes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := New(false)
	shows, err := c.Search(ctx, "Breaking Bad")
	if err != nil || len(shows) == 0 {
		t.Fatalf("search: %v (%d results)", err, len(shows))
	}
	eps, err := c.Episodes(ctx, shows[0].ID)
	if err != nil || len(eps) == 0 {
		t.Fatalf("episodes: %v (%d results)", err, len(eps))
	}
}
