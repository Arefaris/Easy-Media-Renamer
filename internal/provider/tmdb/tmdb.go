package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"myproject/internal/provider"
)

type Client struct {
	Key             string
	Language        string
	HTTP            *http.Client
	IncludeSpecials bool
}

func New(key string, specials bool, languages ...string) *Client {
	language := "en-US"
	if len(languages) > 0 && strings.TrimSpace(languages[0]) != "" {
		language = languages[0]
	}
	return &Client{Key: key, Language: language, HTTP: &http.Client{Timeout: 20 * time.Second}, IncludeSpecials: specials}
}
func (c *Client) Name() string     { return "TMDB" }
func (c *Client) Configured() bool { return strings.TrimSpace(c.Key) != "" }
func (c *Client) Search(ctx context.Context, q string) ([]provider.Show, error) {
	if !c.Configured() {
		return nil, errors.New("TMDB requires an API key in Settings")
	}
	var raw struct {
		Results []struct {
			ID        int    `json:"id"`
			MediaType string `json:"media_type"`
			Name      string `json:"name"`
			Title     string `json:"title"`
			First     string `json:"first_air_date"`
			Release   string `json:"release_date"`
			GenreIDs  []int  `json:"genre_ids"`
		} `json:"results"`
	}
	if e := c.get(ctx, "/search/multi", url.Values{"query": {q}, "include_adult": {"false"}}, &raw); e != nil {
		return nil, e
	}
	out := []provider.Show{}
	for _, v := range raw.Results {
		if v.MediaType != "tv" && v.MediaType != "movie" {
			continue
		}
		name, date := v.Name, v.First
		if v.MediaType == "movie" {
			name, date = v.Title, v.Release
		}
		year := ""
		if len(date) >= 4 {
			year = date[:4]
		}
		out = append(out, provider.Show{ID: v.MediaType + ":" + strconv.Itoa(v.ID), Name: name, Year: year, Kind: v.MediaType})
	}
	return out, nil
}
func (c *Client) Episodes(ctx context.Context, showID string) ([]provider.Episode, error) {
	parts := strings.SplitN(showID, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid TMDB show ID")
	}
	id, e := strconv.Atoi(parts[1])
	if e != nil {
		return nil, e
	}
	if parts[0] == "movie" {
		var m struct {
			Title string `json:"title"`
		}
		if e = c.get(ctx, fmt.Sprintf("/movie/%d", id), nil, &m); e != nil {
			return nil, e
		}
		return []provider.Episode{{Season: 0, Number: 1, Title: m.Title, Special: false}}, nil
	}
	var details struct {
		Seasons []struct {
			Season int `json:"season_number"`
		} `json:"seasons"`
	}
	if e = c.get(ctx, fmt.Sprintf("/tv/%d", id), nil, &details); e != nil {
		return nil, e
	}
	seasons := []int{}
	for _, s := range details.Seasons {
		if s.Season == 0 && !c.IncludeSpecials {
			continue
		}
		seasons = append(seasons, s.Season)
	}
	type result struct {
		eps []provider.Episode
		err error
	}
	ch := make(chan result, len(seasons))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, s := range seasons {
		wg.Add(1)
		go func(season int) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				ch <- result{err: ctx.Err()}
				return
			}
			defer func() { <-sem }()
			var raw struct {
				Episodes []struct {
					Name    string `json:"name"`
					Number  int    `json:"episode_number"`
					Season  int    `json:"season_number"`
					AirDate string `json:"air_date"`
				} `json:"episodes"`
			}
			e := c.get(ctx, fmt.Sprintf("/tv/%d/season/%d", id, season), nil, &raw)
			r := result{err: e}
			for _, v := range raw.Episodes {
				r.eps = append(r.eps, provider.Episode{Season: v.Season, Number: v.Number, Title: v.Name, Special: v.Season == 0, AirDate: v.AirDate})
			}
			ch <- r
		}(s)
	}
	wg.Wait()
	close(ch)
	out := []provider.Episode{}
	for r := range ch {
		if r.err != nil {
			return nil, r.err
		}
		out = append(out, r.eps...)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Season == out[j].Season {
			return out[i].Number < out[j].Number
		}
		return out[i].Season < out[j].Season
	})
	abs := 0
	for i := range out {
		if !out[i].Special {
			abs++
			out[i].Absolute = abs
		}
	}
	return out, nil
}
func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", c.Key)
	params.Set("language", c.Language)
	u := "https://api.themoviedb.org/3" + path + "?" + params.Encode()
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if e != nil {
		return e
	}
	resp, e := c.HTTP.Do(req)
	if e != nil {
		return fmt.Errorf("TMDB request: %w", e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var er struct {
			Message string `json:"status_message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&er)
		return fmt.Errorf("TMDB returned %s: %s", resp.Status, er.Message)
	}
	if e = json.NewDecoder(resp.Body).Decode(dst); e != nil {
		return fmt.Errorf("decode TMDB response: %w", e)
	}
	return nil
}
