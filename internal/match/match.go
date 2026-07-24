package match

import (
	"myproject/internal/provider"
	"sort"
)

type Pair struct {
	FileIndex    int     `json:"file_index"`
	EpisodeIndex int     `json:"episode_index"`
	Confidence   float64 `json:"confidence"`
	Strategy     string  `json:"strategy"`
}

func Match(files []Parsed, episodes []provider.Episode) []Pair {
	out := make([]Pair, 0, len(files))
	used := map[int]bool{}
	oneSeason := singleSeason(episodes)
	for fi, f := range files {
		idx, strategy, conf := -1, "unmatched", 0.0
		if f.Season >= 0 && len(f.Episodes) > 0 {
			idx = find(episodes, func(e provider.Episode) bool { return e.Season == f.Season && e.Number == f.Episodes[0] })
			if idx >= 0 {
				strategy, conf = "season+episode", .99
			}
		}
		if idx < 0 && f.Absolute > 0 {
			idx = find(episodes, func(e provider.Episode) bool { return e.Absolute == f.Absolute })
			if idx < 0 {
				sorted := episodeOrder(episodes)
				if f.Absolute <= len(sorted) {
					idx = sorted[f.Absolute-1]
				}
			}
			if idx >= 0 {
				strategy, conf = "absolute", .9
			}
		}
		if idx < 0 && f.AirDate != "" {
			idx = find(episodes, func(e provider.Episode) bool { return e.AirDate == f.AirDate })
			if idx >= 0 {
				strategy, conf = "airdate", .96
			}
		}
		if idx < 0 && len(f.Episodes) > 0 && f.Season < 0 && oneSeason >= 0 {
			idx = find(episodes, func(e provider.Episode) bool { return e.Season == oneSeason && e.Number == f.Episodes[0] })
			if idx >= 0 {
				strategy, conf = "episode-only", .75
			}
		}
		if idx < 0 && fi < len(episodes) {
			if !used[fi] {
				idx = fi
			} else {
				idx = findIndex(episodes, func(i int) bool { return !used[i] })
			}
			if idx >= 0 {
				strategy, conf = "position", .35
			}
		}
		if idx >= 0 && used[idx] {
			idx = -1
			strategy = "unmatched"
			conf = 0
		}
		if idx >= 0 {
			used[idx] = true
		}
		out = append(out, Pair{fi, idx, conf, strategy})
	}
	return out
}
func findIndex(eps []provider.Episode, p func(int) bool) int {
	for i := range eps {
		if p(i) {
			return i
		}
	}
	return -1
}
func find(eps []provider.Episode, p func(provider.Episode) bool) int {
	for i, e := range eps {
		if p(e) {
			return i
		}
	}
	return -1
}
func singleSeason(eps []provider.Episode) int {
	s := -1
	for _, e := range eps {
		if e.Special {
			continue
		}
		if s < 0 {
			s = e.Season
		} else if s != e.Season {
			return -1
		}
	}
	return s
}
func episodeOrder(eps []provider.Episode) []int {
	ids := make([]int, len(eps))
	for i := range ids {
		ids[i] = i
	}
	sort.Slice(ids, func(i, j int) bool {
		a, b := eps[ids[i]], eps[ids[j]]
		if a.Season == b.Season {
			return a.Number < b.Number
		}
		return a.Season < b.Season
	})
	return ids
}
