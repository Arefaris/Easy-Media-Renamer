package anidb

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type title struct {
	AID        int
	Type       int
	Lang, Name string
}

func downloadTitles(ctx context.Context, httpClient *http.Client) ([]title, error) {
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, "https://anidb.net/api/anime-titles.dat.gz", nil)
	if e != nil {
		return nil, e
	}
	req.Header.Set("User-Agent", "EasyMediaRenamer/1.0")
	resp, e := httpClient.Do(req)
	if e != nil {
		return nil, fmt.Errorf("download AniDB titles: %w", e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("AniDB titles returned %s", resp.Status)
	}
	gz, e := gzip.NewReader(resp.Body)
	if e != nil {
		return nil, e
	}
	defer gz.Close()
	return parseTitles(gz)
}
func parseTitles(r io.Reader) ([]title, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 64*1024), 1024*1024)
	out := []title{}
	for s.Scan() {
		line := s.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p := strings.SplitN(line, "|", 4)
		if len(p) != 4 {
			continue
		}
		aid, e1 := strconv.Atoi(p[0])
		typ, e2 := strconv.Atoi(p[1])
		if e1 != nil || e2 != nil {
			continue
		}
		out = append(out, title{aid, typ, p[2], p[3]})
	}
	return out, s.Err()
}
