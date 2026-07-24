package anidb

import (
	"compress/gzip"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
	ClientName      string
	Version         int
	IncludeSpecials bool
	HTTP            *http.Client
	mu              sync.Mutex
	lastRequest     time.Time
}

const (
	RegisteredClientName    = "supereasymedia"
	RegisteredClientVersion = 1
)

func Default(includeSpecials bool) *Client {
	return New(RegisteredClientName, RegisteredClientVersion, includeSpecials)
}

func New(name string, version int, specials bool) *Client {
	return &Client{ClientName: name, Version: version, IncludeSpecials: specials, HTTP: &http.Client{Timeout: 30 * time.Second}}
}
func (c *Client) Name() string     { return "AniDB" }
func (c *Client) Configured() bool { return strings.TrimSpace(c.ClientName) != "" && c.Version > 0 }
func (c *Client) Search(ctx context.Context, q string) ([]provider.Show, error) {
	titles, e := c.titles(ctx)
	if e != nil {
		return nil, e
	}
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return []provider.Show{}, nil
	}
	type choice struct {
		name string
		rank int
	}
	byID := map[int]choice{}
	for _, t := range titles {
		if !strings.Contains(strings.ToLower(t.Name), q) {
			continue
		}
		rank := 9
		if t.Type == 1 {
			rank = 0
		} else if t.Type == 4 && t.Lang == "en" {
			rank = 1
		} else if t.Type == 4 {
			rank = 2
		}
		old, ok := byID[t.AID]
		if !ok || rank < old.rank {
			byID[t.AID] = choice{t.Name, rank}
		}
	}
	ids := make([]int, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return strings.ToLower(byID[ids[i]].name) < strings.ToLower(byID[ids[j]].name) })
	if len(ids) > 100 {
		ids = ids[:100]
	}
	out := make([]provider.Show, 0, len(ids))
	for _, id := range ids {
		out = append(out, provider.Show{ID: strconv.Itoa(id), Name: byID[id].name, Kind: "anime"})
	}
	return out, nil
}

type animeXML struct {
	Episodes []struct {
		EpNo   string `xml:"epno"`
		Titles []struct {
			Lang string `xml:"lang,attr"`
			Text string `xml:",chardata"`
		} `xml:"title"`
	} `xml:"episodes>episode"`
}

func (c *Client) Episodes(ctx context.Context, id string) ([]provider.Episode, error) {
	if !c.Configured() {
		return nil, errors.New("AniDB requires a registered client name and version in Settings")
	}
	c.mu.Lock()
	wait := 2*time.Second - time.Since(c.lastRequest)
	if wait > 0 {
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			c.mu.Unlock()
			return nil, ctx.Err()
		}
	}
	c.lastRequest = time.Now()
	c.mu.Unlock()
	v := url.Values{"request": {"anime"}, "client": {c.ClientName}, "clientver": {strconv.Itoa(c.Version)}, "protover": {"1"}, "aid": {id}}
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, "http://api.anidb.net:9001/httpapi?"+v.Encode(), nil)
	if e != nil {
		return nil, e
	}
	resp, e := c.HTTP.Do(req)
	if e != nil {
		return nil, fmt.Errorf("AniDB request: %w", e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("AniDB returned %s", resp.Status)
	}
	var r io.Reader = resp.Body
	if !resp.Uncompressed && strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "gzip") {
		gz, e := gzip.NewReader(resp.Body)
		if e != nil {
			return nil, fmt.Errorf("decode AniDB gzip response: %w", e)
		}
		defer gz.Close()
		r = gz
	}
	b, e := io.ReadAll(r)
	if e != nil {
		return nil, e
	}
	if strings.HasPrefix(strings.TrimSpace(string(b)), "<error") {
		var apiError struct {
			Code string `xml:"code,attr"`
			Text string `xml:",chardata"`
		}
		if e = xml.Unmarshal(b, &apiError); e != nil {
			return nil, fmt.Errorf("decode AniDB error: %w", e)
		}
		return nil, fmt.Errorf("AniDB error %s: %s", apiError.Code, strings.TrimSpace(apiError.Text))
	}
	var raw animeXML
	if e = xml.Unmarshal(b, &raw); e != nil {
		return nil, fmt.Errorf("decode AniDB response: %w", e)
	}
	out := []provider.Episode{}
	for _, ep := range raw.Episodes {
		season, num, special, ok := parseEpNo(ep.EpNo)
		if !ok || special && !c.IncludeSpecials {
			continue
		}
		name := "Episode " + ep.EpNo
		for _, t := range ep.Titles {
			if t.Lang == "en" {
				name = t.Text
				break
			}
		}
		out = append(out, provider.Episode{Season: season, Number: num, Title: name, Special: special})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Season == out[j].Season {
			return out[i].Number < out[j].Number
		}
		return out[i].Season < out[j].Season
	})
	return out, nil
}
func parseEpNo(v string) (int, int, bool, bool) {
	v = strings.TrimSpace(v)
	if n, e := strconv.Atoi(v); e == nil {
		return 1, n, false, true
	}
	if len(v) < 2 {
		return 0, 0, false, false
	}
	n, e := strconv.Atoi(v[1:])
	if e != nil {
		return 0, 0, false, false
	}
	switch v[0] {
	case 'S':
		return 0, n, true, true
	case 'C', 'T', 'P', 'O':
		return 0, n, true, true
	}
	return 0, 0, false, false
}
