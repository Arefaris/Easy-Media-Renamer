package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const appDir = "EasyMediaRenamer"

type Config struct {
	TMDBAPIKey      string   `json:"tmdb_api_key"`
	NameTemplate    string   `json:"name_template"`
	MediaExtensions []string `json:"media_extensions"`
	IncludeSpecials bool     `json:"include_specials"`
	Recursive       bool     `json:"recursive"`
	MaxDepth        int      `json:"max_depth"`
	Action          string   `json:"action"`
	AutoMatch       bool     `json:"auto_match"`
	Preset          string   `json:"preset"`
	Language        string   `json:"language"`
}

func Default() Config {
	return Config{
		NameTemplate:    "{show} - s{season}e{episode} - {title}",
		MediaExtensions: []string{".mkv", ".mp4", ".avi", ".m4v", ".mov", ".wmv", ".flv", ".webm", ".ts", ".m2ts", ".mpg", ".mpeg", ".srt", ".ass", ".sub"},
		MaxDepth:        5, Action: "rename", AutoMatch: true, Preset: "Simple", Language: "en-US",
	}
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appDir, "config.json"), nil
}

func Load() (Config, error) {
	cfg := Default()
	p, err := Path()
	if err != nil {
		return cfg, err
	}
	b, err := os.ReadFile(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, err
	}
	if err == nil {
		if err := json.Unmarshal(b, &cfg); err != nil {
			return cfg, err
		}
	}
	if v := os.Getenv("TMDB_API_KEY"); v != "" {
		cfg.TMDBAPIKey = v
	}
	normalize(&cfg)
	return cfg, nil
}

func Save(cfg Config) error {
	normalize(&cfg)
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func normalize(cfg *Config) {
	if strings.TrimSpace(cfg.NameTemplate) == "" {
		cfg.NameTemplate = Default().NameTemplate
	}
	if len(cfg.MediaExtensions) == 0 {
		cfg.MediaExtensions = Default().MediaExtensions
	}
	for i, ext := range cfg.MediaExtensions {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext != "" && !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		cfg.MediaExtensions[i] = ext
	}
	if cfg.MaxDepth <= 0 {
		cfg.MaxDepth = 5
	}
	switch cfg.Action {
	case "rename", "move", "copy", "hardlink", "symlink", "test":
	default:
		cfg.Action = "rename"
	}
	if cfg.Preset == "" {
		cfg.Preset = "Simple"
	}
	if cfg.Language == "" {
		cfg.Language = "en-US"
	}
}
