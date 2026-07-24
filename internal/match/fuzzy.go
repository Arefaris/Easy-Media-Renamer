package match

import (
	"regexp"
	"sort"
	"strings"
)

var articles = map[string]bool{"the": true, "a": true, "an": true}

func NormalizeTitle(s string) string {
	s = strings.ToLower(regexp.MustCompile(`[^\pL\pN]+`).ReplaceAllString(s, " "))
	out := []string{}
	for _, w := range strings.Fields(s) {
		if !articles[w] {
			out = append(out, w)
		}
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}
func Similarity(a, b string) float64 {
	a, b = NormalizeTitle(a), NormalizeTitle(b)
	if a == b {
		return 1
	}
	if a == "" || b == "" {
		return 0
	}
	d := levenshtein([]rune(a), []rune(b))
	m := len([]rune(a))
	if n := len([]rune(b)); n > m {
		m = n
	}
	return 1 - float64(d)/float64(m)
}
func levenshtein(a, b []rune) int {
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i, ca := range a {
		prev := row[0]
		row[0] = i + 1
		for j, cb := range b {
			old := row[j+1]
			cost := 0
			if ca != cb {
				cost = 1
			}
			row[j+1] = min(row[j+1]+1, row[j]+1, prev+cost)
			prev = old
		}
	}
	return row[len(b)]
}
