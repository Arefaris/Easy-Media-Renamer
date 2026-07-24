package anidb

import (
	"strings"
	"testing"
)

func TestParseTitles(t *testing.T) {
	v, e := parseTitles(strings.NewReader("# comment\n1|1|x-jat|Cowboy Bebop\n1|4|en|Cowboy Bebop\n"))
	if e != nil || len(v) != 2 || v[1].Name != "Cowboy Bebop" {
		t.Fatalf("%v %#v", e, v)
	}
}
func TestParseEpNo(t *testing.T) {
	s, n, sp, ok := parseEpNo("12")
	if !ok || sp || s != 1 || n != 12 {
		t.Fatal(s, n, sp, ok)
	}
	s, n, sp, ok = parseEpNo("S2")
	if !ok || !sp || s != 0 || n != 2 {
		t.Fatal(s, n, sp, ok)
	}
}
