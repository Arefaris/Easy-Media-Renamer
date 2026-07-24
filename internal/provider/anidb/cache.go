package anidb

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func cachePath() (string, error) {
	if d := os.Getenv("EMR_CACHE_DIR"); d != "" {
		return filepath.Join(d, "anidb-titles.dat"), nil
	}
	d, e := os.UserCacheDir()
	if e != nil {
		return "", e
	}
	return filepath.Join(d, "EasyMediaRenamer", "anidb-titles.dat"), nil
}
func (c *Client) titles(ctx context.Context) ([]title, error) {
	p, e := cachePath()
	if e != nil {
		return nil, e
	}
	if st, e := os.Stat(p); e == nil && time.Since(st.ModTime()) < 24*time.Hour {
		if f, e := os.Open(p); e == nil {
			defer f.Close()
			return parseTitles(f)
		}
	}
	titles, e := downloadTitles(ctx, c.HTTP)
	if e != nil {
		return nil, e
	}
	if e = os.MkdirAll(filepath.Dir(p), 0700); e == nil {
		if f, x := os.Create(p); x == nil {
			w := bufio.NewWriter(f)
			for _, t := range titles {
				_, _ = fmt.Fprintf(w, "%d|%d|%s|%s\n", t.AID, t.Type, t.Lang, t.Name)
			}
			_ = w.Flush()
			_ = f.Close()
		}
	}
	return titles, nil
}
func parseCachedLine(line string) (title, bool) {
	p := strings.SplitN(line, "|", 4)
	if len(p) != 4 {
		return title{}, false
	}
	a, e := strconv.Atoi(p[0])
	if e != nil {
		return title{}, false
	}
	t, e := strconv.Atoi(p[1])
	return title{a, t, p[2], p[3]}, e == nil
}
