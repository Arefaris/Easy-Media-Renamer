package media

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	matcher "myproject/internal/match"
	"myproject/internal/mediainfo"
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

type Action string

const (
	ActionRename   Action = "rename"
	ActionMove     Action = "move"
	ActionCopy     Action = "copy"
	ActionHardlink Action = "hardlink"
	ActionSymlink  Action = "symlink"
	ActionTest     Action = "test"
)

func Plan(dir string, files []string, episodes []provider.Episode, show provider.Show, tmpl string) ([]Op, error) {
	if len(files) == 0 {
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
		parsed := matcher.Parse(filepath.Base(files[i]))
		ext := filepath.Ext(files[i])
		data := NamingData{
			Episode: episodes[i], Show: show, Ext: ext,
			Filename: strings.TrimSuffix(filepath.Base(files[i]), ext),
			Source:   parsed.Source, Group: parsed.Group, CRC32: parsed.CRC32,
		}
		if NeedsMediaInfo(tmpl) {
			info, _ := mediainfo.Probe(from)
			data.VideoFormat = mediainfo.VideoFormat(info)
			data.Resolution = data.VideoFormat
			data.VideoCodec, data.AudioCodec = info.VideoCodec, info.AudioCodec
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
		to := filepath.Clean(filepath.Join(dir, RenderTemplate(tmpl, data)))
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
	for i := n; i < len(files); i++ {
		from := filepath.Clean(filepath.Join(dir, files[i]))
		ops = append(ops, Op{From: from, To: from, Status: "unmatched"})
	}
	return ops, nil
}

// PlanMatched preserves file-to-episode identity from automatic matching. Every
// source file appears in the result; unmatched files are visible and skipped.
func PlanMatched(dir string, files []string, episodes []provider.Episode, show provider.Show, tmpl string, pairs []matcher.Pair) ([]Op, error) {
	byFile := make(map[int]matcher.Pair, len(pairs))
	for _, pair := range pairs {
		if pair.FileIndex >= 0 && pair.FileIndex < len(files) {
			byFile[pair.FileIndex] = pair
		}
	}
	matchedFiles := []string{}
	matchedEpisodes := []provider.Episode{}
	positions := []int{}
	for i := range files {
		pair, ok := byFile[i]
		if !ok || pair.EpisodeIndex < 0 || pair.EpisodeIndex >= len(episodes) {
			continue
		}
		matchedFiles = append(matchedFiles, files[i])
		matchedEpisodes = append(matchedEpisodes, episodes[pair.EpisodeIndex])
		positions = append(positions, i)
	}
	planned, err := Plan(dir, matchedFiles, matchedEpisodes, show, tmpl)
	if err != nil {
		return nil, err
	}
	out := make([]Op, len(files))
	for i, file := range files {
		from := filepath.Clean(filepath.Join(dir, file))
		out[i] = Op{From: from, To: from, Status: "unmatched"}
	}
	for i, pos := range positions {
		out[pos] = planned[i]
	}
	return out, nil
}

func Apply(ops []Op) (Result, error) {
	return applyRename(ops, ActionRename, true)
}

func ApplyAction(ops []Op, action Action) (Result, error) {
	if action == ActionRename || action == ActionMove {
		return applyRename(ops, action, true)
	}
	if action == ActionTest {
		r := Result{Failed: []Op{}}
		for _, op := range ops {
			if op.Status == "ok" {
				r.Renamed++
			} else {
				r.Skipped++
			}
		}
		return r, nil
	}
	if action != ActionCopy && action != ActionHardlink && action != ActionSymlink {
		return Result{}, fmt.Errorf("unsupported action %q", action)
	}
	res := Result{Failed: []Op{}}
	applied := []Op{}
	targets := map[string]bool{}
	for _, op := range ops {
		if op.Status != "ok" {
			res.Skipped++
			continue
		}
		key := strings.ToLower(filepath.Clean(op.To))
		if targets[key] {
			op.Error = "duplicate target"
			res.Failed = append(res.Failed, op)
			continue
		}
		targets[key] = true
		if _, e := os.Stat(op.To); e == nil {
			op.Error = "target already exists"
			res.Failed = append(res.Failed, op)
			continue
		}
		if e := os.MkdirAll(filepath.Dir(op.To), 0755); e != nil {
			op.Error = e.Error()
			res.Failed = append(res.Failed, op)
			continue
		}
		var e error
		switch action {
		case ActionCopy:
			e = copyFile(op.From, op.To)
		case ActionHardlink:
			e = os.Link(op.From, op.To)
		case ActionSymlink:
			e = os.Symlink(op.From, op.To)
		}
		if e != nil {
			op.Error = e.Error()
			res.Failed = append(res.Failed, op)
		} else {
			res.Renamed++
			applied = append(applied, op)
		}
	}
	if len(applied) > 0 {
		if e := writeHistory(applied, action); e != nil {
			return res, e
		}
	}
	return res, nil
}

func applyRename(ops []Op, action Action, record bool) (Result, error) {
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
		if err := os.MkdirAll(filepath.Dir(s.op.To), 0755); err != nil {
			_ = os.Rename(s.temp, s.op.From)
			s.op.Error = err.Error()
			res.Failed = append(res.Failed, s.op)
			continue
		}
		if err := moveFile(s.temp, s.op.To); err != nil {
			_ = os.Rename(s.temp, s.op.From)
			s.op.Error = err.Error()
			res.Failed = append(res.Failed, s.op)
			continue
		}
		res.Renamed++
		applied = append(applied, s.op)
	}
	if record && len(applied) > 0 {
		if err := writeHistory(applied, action); err != nil {
			return res, fmt.Errorf("renamed files but failed to write undo journal: %w", err)
		}
	}
	return res, nil
}

func moveFile(from, to string) error {
	if e := os.Rename(from, to); e == nil {
		return nil
	}
	if e := copyFile(from, to); e != nil {
		return e
	}
	if e := os.Remove(from); e != nil {
		_ = os.Remove(to)
		return e
	}
	return nil
}
func copyFile(from, to string) error {
	src, e := os.Open(from)
	if e != nil {
		return e
	}
	defer src.Close()
	st, e := src.Stat()
	if e != nil {
		return e
	}
	dst, e := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_EXCL, st.Mode().Perm())
	if e != nil {
		return e
	}
	ok := false
	defer func() {
		_ = dst.Close()
		if !ok {
			_ = os.Remove(to)
		}
	}()
	n, e := io.Copy(dst, src)
	if e != nil {
		return e
	}
	if n != st.Size() {
		return fmt.Errorf("copy size mismatch: %d != %d", n, st.Size())
	}
	if e = dst.Sync(); e != nil {
		return e
	}
	ok = true
	return dst.Close()
}

func Undo() (Result, error) {
	h, e := ListHistory()
	if e != nil {
		return Result{}, e
	}
	if len(h) == 0 {
		return Result{}, errors.New("nothing to undo")
	}
	return RevertHistory(h[0].ID)
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
