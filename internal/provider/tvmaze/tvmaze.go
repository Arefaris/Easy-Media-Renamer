package tvmaze

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"myproject/internal/provider"
)

type Client struct {
	HTTP            *http.Client
	IncludeSpecials bool
}

func New(includeSpecials bool) *Client {
	return &Client{HTTP: &http.Client{Timeout: 15 * time.Second}, IncludeSpecials: includeSpecials}
}
func (c *Client) Name() string     { return "TVmaze" }
func (c *Client) Configured() bool { return true }

func (c *Client) Search(ctx context.Context, query string) ([]provider.Show, error) {
	var raw []struct {
		Show struct {
			ID        int      `json:"id"`
			Name      string   `json:"name"`
			Premiered string   `json:"premiered"`
			Genres    []string `json:"genres"`
		} `json:"show"`
	}
	if err := c.get(ctx, "https://api.tvmaze.com/search/shows?q="+url.QueryEscape(query), &raw); err != nil {
		return nil, err
	}
	out := make([]provider.Show, 0, len(raw))
	for _, v := range raw {
		year := ""
		if len(v.Show.Premiered) >= 4 {
			year = v.Show.Premiered[:4]
		}
		out = append(out, provider.Show{ID: strconv.Itoa(v.Show.ID), Name: v.Show.Name, Year: year, Kind: "tv", Genres: v.Show.Genres})
	}
	return out, nil
}
func (c *Client) Episodes(ctx context.Context, id string) ([]provider.Episode, error) {
	var raw []struct {
		Name    string `json:"name"`
		Season  int    `json:"season"`
		Number  int    `json:"number"`
		Type    string `json:"type"`
		AirDate string `json:"airdate"`
	}
	if err := c.get(ctx, "https://api.tvmaze.com/shows/"+url.PathEscape(id)+"/episodes?specials=1", &raw); err != nil {
		return nil, err
	}
	out := []provider.Episode{}
	for _, e := range raw {
		special := e.Season == 0 || e.Type == "significant_special"
		if special && !c.IncludeSpecials {
			continue
		}
		out = append(out, provider.Episode{Season: e.Season, Number: e.Number, Title: e.Name, Special: special, AirDate: e.AirDate})
	}
	abs := 0
	for i := range out {
		if !out[i].Special {
			abs++
			out[i].Absolute = abs
		}
	}
	return out, nil
}
func (c *Client) get(ctx context.Context, u string, dst any) error {
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if e != nil {
		return e
	}
	req.Header.Set("User-Agent", "EasyMediaRenamer/1.0")
	resp, e := c.HTTP.Do(req)
	if e != nil {
		return fmt.Errorf("TVmaze request: %w", e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TVmaze returned %s", resp.Status)
	}
	if e = json.NewDecoder(resp.Body).Decode(dst); e != nil {
		return fmt.Errorf("decode TVmaze response: %w", e)
	}
	return nil
}
