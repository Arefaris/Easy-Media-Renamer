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
	"myproject/internal/media"
	"myproject/internal/provider"
	"myproject/internal/provider/anidb"
	"myproject/internal/provider/tmdb"
	"myproject/internal/provider/tvmaze"
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
	a.providers = provider.NewRegistry(tvmaze.New(a.cfg.IncludeSpecials), tmdb.New(a.cfg.TMDBAPIKey, a.cfg.IncludeSpecials), anidb.Default(a.cfg.IncludeSpecials))
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
	a.mu.RUnlock()
	return media.ScanDirectory(dir, ext)
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
func (a *App) ApplyRename(dir string, ops []media.Op) (media.Result, error) {
	for _, op := range ops {
		if !inside(dir, op.From) || !inside(dir, op.To) {
			return media.Result{}, errors.New("rename path is outside selected directory")
		}
	}
	return media.Apply(ops)
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
		if filepath.Base(f) != f || !inside(dir, filepath.Join(dir, f)) {
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
