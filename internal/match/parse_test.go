package match

import "testing"

func TestParseRealNames(t *testing.T) {
	tests := []struct {
		name            string
		season, ep, abs int
		date, show      string
	}{
		{"The.Office.US.S02E01.720p.HDTV.x264-GRP.mkv", 2, 1, -1, "", "The Office US"},
		{"Breaking.Bad.S01E05.1080p.BluRay.x264-GRP.mkv", 1, 5, -1, "", "Breaking Bad"},
		{"Show.s1e1.mkv", 1, 1, -1, "", "Show"}, {"Show.S10E22.mp4", 10, 22, -1, "", "Show"},
		{"Show 1x05.avi", 1, 5, -1, "", "Show"}, {"Show.01x01.HDTV.mkv", 1, 1, -1, "", "Show"},
		{"Show Season 1 Episode 3.mkv", 1, 3, -1, "", "Show"}, {"Show.Season.12.Episode.7.mkv", 12, 7, -1, "", "Show"},
		{"Daily.Show.2024.01.15.WEB.mkv", -1, 0, -1, "2024-01-15", "Daily Show"}, {"News-2023-12-31.mp4", -1, 0, -1, "2023-12-31", "News"},
		{"[SubsPlease] Frieren - 12 (1080p) [A1B2C3D4].mkv", -1, 0, 12, "", "Frieren"},
		{"[Erai-raws] Anime Name - 001 [720p].mkv", -1, 0, 1, "", "Anime Name"},
		{"Anime.Name.EP05.mkv", -1, 5, -1, "", "Anime Name"}, {"Anime Name Ep.07.mp4", -1, 7, -1, "", "Anime Name"},
		{"Series.E01.WEB-DL.mkv", -1, 1, -1, "", "Series"}, {"Series 101 HDTV.avi", 1, 1, -1, "", "Series"},
		{"Series.305.mkv", 3, 5, -1, "", "Series"}, {"Show.S01E01-E03.mkv", 1, 1, -1, "", "Show"},
		{"Show.S01E01E02.mkv", 1, 1, -1, "", "Show"}, {"Show.S02E09-E10.1080p.mkv", 2, 9, -1, "", "Show"},
		{"The.Expanse.S06E06.2160p.WEB-DL.x265-GRP.mkv", 6, 6, -1, "", "The Expanse"},
		{"Mr.Robot.S04E13.1080p.WEBRip.AV1.mkv", 4, 13, -1, "", "Mr Robot"},
		{"Doctor.Who.2005.S01E01.DVDRip.avi", 1, 1, -1, "", "Doctor Who"},
		{"Westworld.S01E01.Remux.2160p.HEVC.mkv", 1, 1, -1, "", "Westworld"},
		{"Fargo.S05E10.1080i.HDTV.mkv", 5, 10, -1, "", "Fargo"},
		{"Loki.S02E03.4K.WEB.mkv", 2, 3, -1, "", "Loki"}, {"Lost.S01E01.h264.mkv", 1, 1, -1, "", "Lost"},
		{"Foundation_Season_2_Episode_4.mp4", 2, 4, -1, "", "Foundation"},
		{"Jeopardy.2024-02-20.720p.mkv", -1, 0, -1, "2024-02-20", "Jeopardy"},
		{"Some.Show.S01E99.mkv", 1, 99, -1, "", "Some Show"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Parse(tc.name)
			if p.Season != tc.season || p.Absolute != tc.abs || p.AirDate != tc.date || p.Show != tc.show {
				t.Fatalf("got %#v", p)
			}
			if tc.ep > 0 && (len(p.Episodes) == 0 || p.Episodes[0] != tc.ep) {
				t.Fatalf("episodes: %#v", p)
			}
		})
	}
}
func TestParseMetadata(t *testing.T) {
	p := Parse("Show.S01E01.2160p.BluRay.HEVC-GRP.mkv")
	if p.Res != "2160p" || p.Source != "BluRay" || p.CodecTag != "x265" || p.Group != "GRP" {
		t.Fatalf("%#v", p)
	}
	p = Parse("[SubsPlease] Frieren - 12 (1080p) [A1B2C3D4].mkv")
	if p.CRC32 != "A1B2C3D4" || p.Group != "SubsPlease" {
		t.Fatalf("%#v", p)
	}
}
func TestMultiEpisode(t *testing.T) {
	p := Parse("Show.S01E01-E03.mkv")
	if len(p.Episodes) != 3 || p.Episodes[2] != 3 {
		t.Fatalf("%#v", p)
	}
}

func TestResolutionOnlyIsNotCompactEpisode(t *testing.T) {
	for _, name := range []string{"Show.720p.mkv", "Show.576p.mkv", "Some.Movie.Title.480p.WEB.mkv"} {
		p := Parse(name)
		if p.Season != -1 || len(p.Episodes) != 0 {
			t.Errorf("%s: %#v", name, p)
		}
	}
}

func TestHyphenatedTitleAndReleaseGroup(t *testing.T) {
	for _, tc := range []struct{ name, show, group string }{
		{"Spider-Man.S01E01.720p.HDTV-KILLERS.mkv", "Spider-Man", "KILLERS"},
		{"Law-and-Order.S05E12.HDTV-LOL.mkv", "Law-And-Order", "LOL"},
	} {
		p := Parse(tc.name)
		if p.Show != tc.show || p.Group != tc.group {
			t.Errorf("%s: %#v", tc.name, p)
		}
	}
}
