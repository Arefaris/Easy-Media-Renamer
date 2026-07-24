package media

import (
	"reflect"
	"sort"
	"testing"
)

func TestNaturalLess(t *testing.T) {
	got := []string{"ep10.mkv", "ep2.mkv", "ep1.mkv"}
	sort.Slice(got, func(i, j int) bool { return NaturalLess(got[i], got[j]) })
	want := []string{"ep1.mkv", "ep2.mkv", "ep10.mkv"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
