//go:build live

package anidb

import (
	"context"
	"testing"
	"time"
)

func TestLiveSearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	c := New("", 0, false)
	shows, err := c.Search(ctx, "Cowboy Bebop")
	if err != nil || len(shows) == 0 {
		t.Fatalf("search: %v (%d results)", err, len(shows))
	}
}
func TestLiveEpisodes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	eps, err := Default(false).Episodes(ctx, "23")
	if err != nil || len(eps) == 0 {
		t.Fatalf("episodes: %v (%d results)", err, len(eps))
	}
}
