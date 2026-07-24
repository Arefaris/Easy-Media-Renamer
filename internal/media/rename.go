package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myproject/internal/provider"
)

type Op struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
type Result struct {
	Renamed int  `json:"renamed"`
	Skipped int  `json:"skipped"`
	Failed  []Op `json:"failed"`
}

func Plan(dir string, files []string, episodes []provider.Episode, show provider.Show, tmpl string) ([]Op, error) {
	if len(files) == 0 || len(episodes) == 0 {
		return []Op{}, nil
	}
	seen := map[string]bool{}
	sources := map[string]bool{}
	for _, f := range files {
		sources[strings.ToLower(filepath.Clean(filepath.Join(dir, f)))] = true
	}
	n := len(files)
	if len(episodes) < n {
		n = len(episodes)
	}
	ops := make([]Op, 0, n)
	for i := 0; i < n; i++ {
		from := filepath.Clean(filepath.Join(dir, files[i]))
		to := filepath.Clean(filepath.Join(dir, Render(tmpl, episodes[i], show, filepath.Ext(files[i]))))
		op := Op{From: from, To: to, Status: "ok"}
		key := strings.ToLower(to)
		if strings.EqualFold(from, to) {
			op.Status = "unchanged"
		} else if seen[key] {
			op.Status = "conflict-duplicate"
		} else if _, err := os.Stat(to); err == nil && !sources[key] {
			op.Status = "conflict-existing"
		}
		seen[key] = true
		ops = append(ops, op)
	}
	return ops, nil
}

func journalPath() (string, error) {
	if d := os.Getenv("EMR_CACHE_DIR"); d != "" {
		return filepath.Join(d, "undo", "last.json"), nil
	}
	d, e := os.UserCacheDir()
	if e != nil {
		return "", e
	}
	return filepath.Join(d, "EasyMediaRenamer", "undo", "last.json"), nil
}

func Apply(ops []Op) (Result, error) {
	res := Result{Failed: []Op{}}
	valid := make([]Op, 0, len(ops))
	sources := make(map[string]bool, len(ops))
	for _, op := range ops {
		if op.Status == "ok" {
			sources[strings.ToLower(filepath.Clean(op.From))] = true
		}
	}
	targets := make(map[string]bool, len(ops))
	for _, op := range ops {
		if op.Status != "ok" {
			res.Skipped++
			continue
		}
		toKey := strings.ToLower(filepath.Clean(op.To))
		if targets[toKey] {
			op.Error = "duplicate target"
			res.Failed = append(res.Failed, op)
			continue
		}
		targets[toKey] = true
		if _, err := os.Stat(op.To); err == nil && !sources[toKey] {
			op.Error = "target already exists"
			res.Failed = append(res.Failed, op)
			continue
		}
		valid = append(valid, op)
	}
	type staged struct {
		op   Op
		temp string
	}
	stages := make([]staged, 0, len(valid))
	for i, op := range valid {
		if _, err := os.Stat(op.From); err != nil {
			op.Error = err.Error()
			res.Failed = append(res.Failed, op)
			continue
		}
		tmp := filepath.Join(filepath.Dir(op.From), fmt.Sprintf(".emr-%d-%d.tmp", time.Now().UnixNano(), i))
		if err := os.Rename(op.From, tmp); err != nil {
			op.Error = err.Error()
			res.Failed = append(res.Failed, op)
			continue
		}
		stages = append(stages, staged{op, tmp})
	}
	applied := make([]Op, 0, len(stages))
	for _, s := range stages {
		if err := os.Rename(s.temp, s.op.To); err != nil {
			_ = os.Rename(s.temp, s.op.From)
			s.op.Error = err.Error()
			res.Failed = append(res.Failed, s.op)
			continue
		}
		res.Renamed++
		applied = append(applied, s.op)
	}
	if len(applied) > 0 {
		if err := writeJournal(applied); err != nil {
			return res, fmt.Errorf("renamed files but failed to write undo journal: %w", err)
		}
	}
	return res, nil
}

func writeJournal(ops []Op) error {
	p, e := journalPath()
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		return e
	}
	b, e := json.MarshalIndent(ops, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile(p, b, 0600)
}

func Undo() (Result, error) {
	p, e := journalPath()
	if e != nil {
		return Result{}, e
	}
	b, e := os.ReadFile(p)
	if errors.Is(e, os.ErrNotExist) {
		return Result{}, errors.New("nothing to undo")
	}
	if e != nil {
		return Result{}, e
	}
	var ops []Op
	if e = json.Unmarshal(b, &ops); e != nil {
		return Result{}, e
	}
	reverse := make([]Op, 0, len(ops))
	for i := len(ops) - 1; i >= 0; i-- {
		reverse = append(reverse, Op{From: ops[i].To, To: ops[i].From, Status: "ok"})
	}
	res, e := ApplyWithoutJournal(reverse)
	if e == nil && len(res.Failed) == 0 {
		_ = os.Remove(p)
	}
	return res, e
}

func ApplyWithoutJournal(ops []Op) (Result, error) {
	res := Result{Failed: []Op{}}
	sources := make(map[string]bool, len(ops))
	for _, op := range ops {
		sources[strings.ToLower(filepath.Clean(op.From))] = true
	}
	type staged struct {
		op   Op
		temp string
	}
	ss := []staged{}
	for i, op := range ops {
		if _, err := os.Stat(op.To); err == nil && !sources[strings.ToLower(filepath.Clean(op.To))] {
			op.Error = "target already exists"
			res.Failed = append(res.Failed, op)
			continue
		}
		tmp := fmt.Sprintf("%s.undo-%d", op.From, i)
		if e := os.Rename(op.From, tmp); e != nil {
			op.Error = e.Error()
			res.Failed = append(res.Failed, op)
			continue
		}
		ss = append(ss, staged{op, tmp})
	}
	for _, s := range ss {
		if e := os.Rename(s.temp, s.op.To); e != nil {
			_ = os.Rename(s.temp, s.op.From)
			s.op.Error = e.Error()
			res.Failed = append(res.Failed, s.op)
		} else {
			res.Renamed++
		}
	}
	return res, nil
}
