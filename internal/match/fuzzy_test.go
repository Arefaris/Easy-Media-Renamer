package match

import "testing"

func TestSimilarity(t *testing.T) {
	for _, x := range [][2]string{{"The Office US", "The Office (US)"}, {"BREAKING BAD", "breaking bad"}, {"A Team", "Team"}} {
		if s := Similarity(x[0], x[1]); s < .9 {
			t.Errorf("%q %q = %f", x[0], x[1], s)
		}
	}
	if Similarity("Lost", "Friends") > .5 {
		t.Fatal("unrelated titles too similar")
	}
}
