package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"myproject/internal/config"
	matcher "myproject/internal/match"
	"myproject/internal/media"
	"myproject/internal/mediainfo"
	"myproject/internal/provider"
	"myproject/internal/provider/anidb"
	"myproject/internal/provider/tmdb"
	"myproject/internal/provider/tvmaze"
	"myproject/internal/verify"
)

type App struct {
	ctx       context.Context
	mu        sync.RWMutex
	cfg       config.Config
	providers *provider.Registry
}

func NewApp() *App {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}
	a := &App{cfg: cfg}
	a.rebuildProviders()
	return a
}
func (a *App) startup(ctx context.Context) { a.ctx = ctx }
func (a *App) rebuildProviders() {
	a.providers = provider.NewRegistry(tvmaze.New(a.cfg.IncludeSpecials), tmdb.New(a.cfg.TMDBAPIKey, a.cfg.IncludeSpecials, a.cfg.Language), anidb.Default(a.cfg.IncludeSpecials))
}
func (a *App) ListProviders() []provider.Info {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.providers.List()
}
func (a *App) provider(name string) (provider.Provider, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	p, ok := a.providers.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown provider %q", name)
	}
	if !p.Configured() {
		return nil, fmt.Errorf("%s requires configuration", name)
	}
	return p, nil
}
func (a *App) SearchShows(name, query string) ([]provider.Show, error) {
	if strings.TrimSpace(query) == "" {
		return []provider.Show{}, nil
	}
	p, e := a.provider(name)
	if e != nil {
		return nil, e
	}
	ctx, cancel := context.WithTimeout(a.context(), 30*time.Second)
	defer cancel()
	return p.Search(ctx, query)
}
func (a *App) ListEpisodes(name, id string) ([]provider.Episode, error) {
	p, e := a.provider(name)
	if e != nil {
		return nil, e
	}
	ctx, cancel := context.WithTimeout(a.context(), 60*time.Second)
	defer cancel()
	return p.Episodes(ctx, id)
}
func (a *App) OpenDirectoryDialog() (string, error) {
	if a.ctx == nil {
		return "", errors.New("application is not started")
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Select media directory"})
}
func (a *App) ScanDirectory(dir string) ([]media.File, error) {
	a.mu.RLock()
	ext := append([]string(nil), a.cfg.MediaExtensions...)
	recursive, maxDepth := a.cfg.Recursive, a.cfg.MaxDepth
	a.mu.RUnlock()
	return media.Scan(dir, ext, recursive, maxDepth)
}
func (a *App) ScanDirectoryWithOptions(dir string, recursive bool, maxDepth int) ([]media.File, error) {
	a.mu.RLock()
	ext := append([]string(nil), a.cfg.MediaExtensions...)
	a.mu.RUnlock()
	return media.Scan(dir, ext, recursive, maxDepth)
}
func (a *App) AutoMatch(files []media.File, eps []provider.Episode) []matcher.Pair {
	parsed := make([]matcher.Parsed, len(files))
	for i, f := range files {
		parsed[i] = matcher.Parse(f.Name)
	}
	return matcher.Match(parsed, eps)
}
func (a *App) ListPresets() []media.Preset { return media.Presets }
func (a *App) RenderNamePreview(file media.File, ep provider.Episode, show provider.Show) (string, error) {
	a.mu.RLock()
	tmpl := a.cfg.NameTemplate
	a.mu.RUnlock()
	return a.renderTemplatePreview(tmpl, file, ep, show)
}
func (a *App) RenderTemplatePreview(tmpl string, file media.File, ep provider.Episode, show provider.Show) (string, error) {
	return a.renderTemplatePreview(tmpl, file, ep, show)
}
func (a *App) renderTemplatePreview(tmpl string, file media.File, ep provider.Episode, show provider.Show) (string, error) {
	p := matcher.Parse(file.Name)
	data := media.NamingData{Episode: ep, Show: show, Ext: filepath.Ext(file.Name), Filename: strings.TrimSuffix(file.Name, filepath.Ext(file.Name)), Source: p.Source, Group: p.Group, CRC32: p.CRC32}
	if media.NeedsMediaInfo(tmpl) {
		info, e := mediainfo.Probe(file.Path)
		if e != nil {
			return "", e
		}
		data.VideoFormat = mediainfo.VideoFormat(info)
		data.Resolution = data.VideoFormat
		data.VideoCodec = info.VideoCodec
		data.AudioCodec = info.AudioCodec
		if info.Channels > 0 {
			data.Channels = fmt.Sprintf("%d", info.Channels)
		}
		if info.HDR {
			data.HD = "HDR"
		}
		if info.Duration > 0 {
			data.Duration = info.Duration.String()
		}
	}
	return media.RenderTemplate(tmpl, data), nil
}
func (a *App) ListHistory() ([]media.HistoryEntry, error)        { return media.ListHistory() }
func (a *App) RevertHistory(id string) (media.Result, error)     { return media.RevertHistory(id) }
func (a *App) ClearHistory() error                               { return media.ClearHistory() }
func (a *App) ComputeChecksum(path, algo string) (string, error) { return verify.Compute(path, algo) }
func (a *App) VerifyEmbeddedCRC(path string) (verify.Result, error) {
	return verify.VerifyEmbeddedCRC(path)
}
func (a *App) PreviewRename(dir string, files []string, eps []provider.Episode, show provider.Show) ([]media.Op, error) {
	a.mu.RLock()
	tmpl := a.cfg.NameTemplate
	a.mu.RUnlock()
	if e := validateFilesInDir(dir, files); e != nil {
		return nil, e
	}
	return media.Plan(dir, files, eps, show, tmpl)
}
func (a *App) PreviewMatched(dir string, files []media.File, eps []provider.Episode, pairs []matcher.Pair, show provider.Show) ([]media.Op, error) {
	a.mu.RLock()
	tmpl := a.cfg.NameTemplate
	a.mu.RUnlock()
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.RelPath
		if names[i] == "" {
			names[i] = file.Name
		}
	}
	if e := validateFilesInDir(dir, names); e != nil {
		return nil, e
	}
	return media.PlanMatched(dir, names, eps, show, tmpl, pairs)
}
func (a *App) ApplyRename(dir string, ops []media.Op) (media.Result, error) {
	for _, op := range ops {
		if !inside(dir, op.From) || !inside(dir, op.To) {
			return media.Result{}, errors.New("rename path is outside selected directory")
		}
	}
	a.mu.RLock()
	action := media.Action(a.cfg.Action)
	a.mu.RUnlock()
	return media.ApplyAction(ops, action)
}
func (a *App) UndoLastRename() (media.Result, error) { return media.Undo() }
func (a *App) GetConfig() (config.Config, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cfg, nil
}
func (a *App) SaveConfig(cfg config.Config) error {
	if err := config.Save(cfg); err != nil {
		return err
	}
	loaded, err := config.Load()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.cfg = loaded
	a.rebuildProviders()
	a.mu.Unlock()
	return nil
}
func (a *App) context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}
func validateFilesInDir(dir string, files []string) error {
	for _, f := range files {
		if filepath.IsAbs(f) || f == "" || !inside(dir, filepath.Join(dir, f)) {
			return fmt.Errorf("invalid file name %q", f)
		}
	}
	return nil
}
func inside(dir, path string) bool {
	base, e1 := filepath.Abs(dir)
	target, e2 := filepath.Abs(path)
	if e1 != nil || e2 != nil {
		return false
	}
	rel, e := filepath.Rel(base, target)
	return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
