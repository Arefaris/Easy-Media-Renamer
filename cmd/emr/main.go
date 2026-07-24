package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myproject/internal/config"
	matcher "myproject/internal/match"
	"myproject/internal/media"
	"myproject/internal/provider"
	"myproject/internal/provider/anidb"
	"myproject/internal/provider/tmdb"
	"myproject/internal/provider/tvmaze"
	"myproject/internal/verify"
)

func main() { os.Exit(run()) }
func run() int {
	renameMode := flag.Bool("rename", false, "match and rename media files")
	check := flag.Bool("check", false, "compute checksums")
	list := flag.Bool("list", false, "list matching media files")
	recursive := flag.Bool("r", false, "scan recursively")
	db := flag.String("db", "TVmaze", "metadata provider: TVmaze, TMDB, AniDB")
	format := flag.String("format", "", "naming template")
	action := flag.String("action", "rename", "rename, move, copy, hardlink, symlink or test")
	filter := flag.String("filter", "", "comma-separated extensions")
	dry := flag.Bool("dry-run", false, "show operations without writing")
	lang := flag.String("lang", "en-US", "metadata language")
	algo := flag.String("algo", "sha256", "checksum algorithm")
	flag.Parse()
	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	cfg, e := config.Load()
	if e != nil {
		return fail(e)
	}
	exts := cfg.MediaExtensions
	if *filter != "" {
		exts = nil
		for _, x := range strings.Split(*filter, ",") {
			x = strings.TrimSpace(x)
			if !strings.HasPrefix(x, ".") {
				x = "." + x
			}
			exts = append(exts, x)
		}
	}
	files, e := media.Scan(dir, exts, *recursive, cfg.MaxDepth)
	if e != nil {
		return fail(e)
	}
	if *list {
		for _, f := range files {
			fmt.Printf("%10d  %s\n", f.Size, f.RelPath)
		}
		return 0
	}
	if *check {
		bad := false
		for _, f := range files {
			sum, x := verify.Compute(f.Path, *algo)
			if x != nil {
				fmt.Fprintf(os.Stderr, "ERROR %s: %v\n", f.RelPath, x)
				bad = true
			} else {
				fmt.Printf("%s  %s\n", sum, f.RelPath)
			}
		}
		if bad {
			return 2
		}
		return 0
	}
	if !*renameMode {
		flag.Usage()
		return 2
	}
	if len(files) == 0 {
		return fail(fmt.Errorf("no media files found"))
	}
	reg := provider.NewRegistry(tvmaze.New(cfg.IncludeSpecials), tmdb.New(cfg.TMDBAPIKey, cfg.IncludeSpecials, *lang), anidb.Default(cfg.IncludeSpecials))
	p, ok := reg.Get(*db)
	if !ok {
		return fail(fmt.Errorf("unknown provider %q", *db))
	}
	if !p.Configured() {
		return fail(fmt.Errorf("%s is not configured", *db))
	}
	parsed := make([]matcher.Parsed, len(files))
	query := ""
	for i, f := range files {
		parsed[i] = matcher.Parse(f.Name)
		if query == "" && parsed[i].Show != "" {
			query = parsed[i].Show
		}
	}
	if query == "" {
		return fail(fmt.Errorf("could not infer show name"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	shows, e := p.Search(ctx, query)
	if e != nil {
		return fail(e)
	}
	if len(shows) == 0 {
		return fail(fmt.Errorf("no show found for %q", query))
	}
	show := shows[0]
	best := 0.0
	for _, s := range shows {
		if score := matcher.Similarity(query, s.Name); score > best {
			best = score
			show = s
		}
	}
	eps, e := p.Episodes(ctx, show.ID)
	if e != nil {
		return fail(e)
	}
	pairs := matcher.Match(parsed, eps)
	names := []string{}
	ordered := []provider.Episode{}
	for _, pair := range pairs {
		if pair.EpisodeIndex < 0 {
			continue
		}
		names = append(names, files[pair.FileIndex].RelPath)
		ordered = append(ordered, eps[pair.EpisodeIndex])
		fmt.Printf("%-45s -> S%02dE%02d  %-15s %.0f%%\n", files[pair.FileIndex].RelPath, eps[pair.EpisodeIndex].Season, eps[pair.EpisodeIndex].Number, pair.Strategy, pair.Confidence*100)
	}
	tmpl := *format
	if tmpl == "" {
		tmpl = cfg.NameTemplate
	}
	ops, e := media.Plan(dir, names, ordered, show, tmpl)
	if e != nil {
		return fail(e)
	}
	for _, op := range ops {
		fmt.Printf("%s -> %s [%s]\n", filepath.Base(op.From), op.To, op.Status)
	}
	act := media.Action(*action)
	if *dry {
		act = media.ActionTest
	}
	res, e := media.ApplyAction(ops, act)
	if e != nil {
		return fail(e)
	}
	fmt.Printf("completed=%d skipped=%d failed=%d\n", res.Renamed, res.Skipped, len(res.Failed))
	if len(res.Failed) > 0 {
		return 3
	}
	return 0
}
func fail(e error) int { fmt.Fprintln(os.Stderr, "error:", e); return 1 }
