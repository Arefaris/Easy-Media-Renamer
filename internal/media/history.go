package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type HistoryEntry struct {
	ID        string    `json:"id"`
	Time      time.Time `json:"time"`
	Directory string    `json:"directory"`
	Action    Action    `json:"action"`
	Count     int       `json:"count"`
	Ops       []Op      `json:"ops"`
}

func historyDir() (string, error) {
	if d := os.Getenv("EMR_CACHE_DIR"); d != "" {
		return filepath.Join(d, "history"), nil
	}
	d, e := os.UserCacheDir()
	if e != nil {
		return "", e
	}
	return filepath.Join(d, "EasyMediaRenamer", "history"), nil
}
func writeHistory(ops []Op, action Action) error {
	d, e := historyDir()
	if e != nil {
		return e
	}
	if e = os.MkdirAll(d, 0700); e != nil {
		return e
	}
	now := time.Now().UTC()
	id := strconv.FormatInt(now.UnixNano(), 10)
	dir := ""
	if len(ops) > 0 {
		dir = filepath.Dir(ops[0].From)
	}
	entry := HistoryEntry{ID: id, Time: now, Directory: dir, Action: action, Count: len(ops), Ops: ops}
	b, e := json.MarshalIndent(entry, "", "  ")
	if e != nil {
		return e
	}
	if e = os.WriteFile(filepath.Join(d, id+".json"), b, 0600); e != nil {
		return e
	}
	return rotateHistory(d, 50)
}
func ListHistory() ([]HistoryEntry, error) {
	d, e := historyDir()
	if e != nil {
		return nil, e
	}
	entries, e := os.ReadDir(d)
	if errors.Is(e, os.ErrNotExist) {
		return []HistoryEntry{}, nil
	}
	if e != nil {
		return nil, e
	}
	out := []HistoryEntry{}
	for _, de := range entries {
		if de.IsDir() || filepath.Ext(de.Name()) != ".json" {
			continue
		}
		b, x := os.ReadFile(filepath.Join(d, de.Name()))
		if x != nil {
			continue
		}
		var h HistoryEntry
		if json.Unmarshal(b, &h) == nil {
			out = append(out, h)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.After(out[j].Time) })
	return out, nil
}
func RevertHistory(id string) (Result, error) {
	list, e := ListHistory()
	if e != nil {
		return Result{}, e
	}
	var h *HistoryEntry
	for i := range list {
		if list[i].ID == id {
			h = &list[i]
			break
		}
	}
	if h == nil {
		return Result{}, fmt.Errorf("history entry %s not found", id)
	}
	res := Result{Failed: []Op{}}
	switch h.Action {
	case ActionCopy, ActionHardlink, ActionSymlink:
		for i := len(h.Ops) - 1; i >= 0; i-- {
			op := h.Ops[i]
			if e := os.Remove(op.To); e != nil {
				op.Error = e.Error()
				res.Failed = append(res.Failed, op)
			} else {
				res.Renamed++
			}
		}
	default:
		rev := make([]Op, 0, len(h.Ops))
		for i := len(h.Ops) - 1; i >= 0; i-- {
			rev = append(rev, Op{From: h.Ops[i].To, To: h.Ops[i].From, Status: "ok"})
		}
		res, e = ApplyWithoutJournal(rev)
	}
	if e == nil && len(res.Failed) == 0 {
		d, _ := historyDir()
		_ = os.Remove(filepath.Join(d, id+".json"))
	}
	return res, e
}
func ClearHistory() error {
	d, e := historyDir()
	if e != nil {
		return e
	}
	entries, e := os.ReadDir(d)
	if errors.Is(e, os.ErrNotExist) {
		return nil
	}
	if e != nil {
		return e
	}
	for _, de := range entries {
		if !de.IsDir() && filepath.Ext(de.Name()) == ".json" {
			if e = os.Remove(filepath.Join(d, de.Name())); e != nil {
				return e
			}
		}
	}
	return nil
}
func rotateHistory(d string, keep int) error {
	list, e := ListHistory()
	if e != nil {
		return e
	}
	if len(list) <= keep {
		return nil
	}
	for _, h := range list[keep:] {
		if e = os.Remove(filepath.Join(d, h.ID+".json")); e != nil {
			return e
		}
	}
	return nil
}
