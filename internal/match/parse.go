package match

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Parsed struct {
	Show     string  `json:"show"`
	Year     string  `json:"year"`
	Season   int     `json:"season"`
	Episodes []int   `json:"episodes"`
	Absolute int     `json:"absolute"`
	AirDate  string  `json:"airdate"`
	Source   string  `json:"source"`
	Group    string  `json:"group"`
	Res      string  `json:"res"`
	CodecTag string  `json:"codec_tag"`
	CRC32    string  `json:"crc32"`
	Conf     float64 `json:"confidence"`
}

var (
	multiRE   = regexp.MustCompile(`(?i)\bS(\d{1,2})E(\d{1,3})(?:(?:\s*-\s*E?)|E)(\d{1,3})(?:\b|$)`)
	sxeRE     = regexp.MustCompile(`(?i)\bS(\d{1,2})E(\d{1,3})\b`)
	xRE       = regexp.MustCompile(`(?i)\b(\d{1,2})x(\d{1,3})\b`)
	longRE    = regexp.MustCompile(`(?i)Season[ ._-]+(\d{1,2})[ ._-]+Episode[ ._-]+(\d{1,3})`)
	dateRE    = regexp.MustCompile(`\b((?:19|20)\d{2})[._-](\d{2})[._-](\d{2})\b`)
	animeRE   = regexp.MustCompile(`(?i)^(?:\[([^\]]+)\]\s*)?(.+?)[ ._-]+-[ ._-]*(\d{1,4})(?:v\d)?(?:\s|\(|\[|$)`)
	episodeRE = regexp.MustCompile(`(?i)\bEP?[ ._-]?(\d{1,3})\b`)
	compactRE = regexp.MustCompile(`(?:^|\D)(\d)(\d{2})(?:\D|$)`)
	yearRE    = regexp.MustCompile(`(?:^|[ ._(-])((?:19|20)\d{2})(?:$|[ ._)])`)
	resRE     = regexp.MustCompile(`(?i)\b(480p|576p|720p|1080[pi]|2160p|4k)\b`)
	sourceRE  = regexp.MustCompile(`(?i)\b(BluRay|BDRip|WEB[ ._-]?DL|WEBRip|WEB|HDTV|DVDRip|DVD|Remux)\b`)
	codecRE   = regexp.MustCompile(`(?i)\b(x264|x265|h[ ._-]?264|h[ ._-]?265|HEVC|AV1)\b`)
	crcRE     = regexp.MustCompile(`(?i)\[([A-F0-9]{8})\]`)
	groupRE   = regexp.MustCompile(`-([A-Za-z0-9][A-Za-z0-9._]*)$`)
)

func Parse(filename string) Parsed {
	p := Parsed{Season: -1, Absolute: -1}
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	animeName := ""
	if m := crcRE.FindStringSubmatch(base); m != nil {
		p.CRC32 = strings.ToUpper(m[1])
	}
	if m := resRE.FindStringSubmatch(base); m != nil {
		p.Res = normalizeResolution(m[1])
	}
	if m := sourceRE.FindStringSubmatch(base); m != nil {
		p.Source = normalizeSource(m[1])
	}
	if m := codecRE.FindStringSubmatch(base); m != nil {
		p.CodecTag = normalizeCodec(m[1])
	}
	if m := groupRE.FindStringSubmatch(base); m != nil && !isReleaseTag(m[1]) {
		p.Group = m[1]
	}
	if m := yearRE.FindStringSubmatch(base); m != nil {
		p.Year = m[1]
	}

	marker := len(base)
	setMarker := func(i int) {
		if i >= 0 && i < marker {
			marker = i
		}
	}
	if m := multiRE.FindStringSubmatchIndex(base); m != nil {
		p.Season = atoi(base[m[2]:m[3]])
		first := atoi(base[m[4]:m[5]])
		segment := base[m[0]:m[1]]
		nums := regexp.MustCompile(`(?i)E(\d{1,3})`).FindAllStringSubmatch(segment, -1)
		p.Episodes = []int{first}
		for _, n := range nums[1:] {
			p.Episodes = append(p.Episodes, atoi(n[1]))
		}
		if strings.Contains(segment, "-") && len(p.Episodes) == 2 && p.Episodes[1] > first+1 {
			end := p.Episodes[1]
			p.Episodes = p.Episodes[:1]
			for n := first + 1; n <= end; n++ {
				p.Episodes = append(p.Episodes, n)
			}
		}
		p.Conf = .99
		setMarker(m[0])
	} else if m := sxeRE.FindStringSubmatchIndex(base); m != nil {
		p.Season = atoi(base[m[2]:m[3]])
		p.Episodes = []int{atoi(base[m[4]:m[5]])}
		p.Conf = .98
		setMarker(m[0])
	} else if m := xRE.FindStringSubmatchIndex(base); m != nil {
		p.Season = atoi(base[m[2]:m[3]])
		p.Episodes = []int{atoi(base[m[4]:m[5]])}
		p.Conf = .95
		setMarker(m[0])
	} else if m := longRE.FindStringSubmatchIndex(base); m != nil {
		p.Season = atoi(base[m[2]:m[3]])
		p.Episodes = []int{atoi(base[m[4]:m[5]])}
		p.Conf = .96
		setMarker(m[0])
	} else if m := dateRE.FindStringSubmatchIndex(base); m != nil {
		p.AirDate = base[m[2]:m[3]] + "-" + base[m[4]:m[5]] + "-" + base[m[6]:m[7]]
		p.Conf = .96
		setMarker(m[0])
	} else if m := animeRE.FindStringSubmatchIndex(base); m != nil {
		p.Absolute = atoi(base[m[6]:m[7]])
		p.Conf = .9
		animeName = base[m[4]:m[5]]
		setMarker(m[5])
		if m[2] >= 0 {
			p.Group = base[m[2]:m[3]]
		}
	}
	if len(p.Episodes) == 0 && p.Absolute < 0 && p.AirDate == "" {
		if m := episodeRE.FindStringSubmatchIndex(base); m != nil {
			p.Episodes = []int{atoi(base[m[2]:m[3]])}
			p.Conf = .72
			setMarker(m[0])
		} else if m := compactRE.FindStringSubmatchIndex(base); m != nil {
			nextIsScanType := m[5] < len(base) && (base[m[5]] == 'p' || base[m[5]] == 'P' || base[m[5]] == 'i' || base[m[5]] == 'I')
			if !nextIsScanType {
				p.Season = atoi(base[m[2]:m[3]])
				p.Episodes = []int{atoi(base[m[4]:m[5]])}
				p.Conf = .58
				setMarker(m[2])
			}
		}
	}
	name := base[:marker]
	if animeName != "" {
		name = animeName
	}
	name = regexp.MustCompile(`^\[[^\]]+\]\s*`).ReplaceAllString(name, "")
	name = yearRE.ReplaceAllString(name, " ")
	name = releaseStrip(name)
	p.Show = cleanTitle(name)
	return p
}

func cleanTitle(v string) string {
	v = strings.Trim(regexp.MustCompile(`[._]+`).ReplaceAllString(v, " "), " -._")
	v = regexp.MustCompile(`\s+`).ReplaceAllString(v, " ")
	words := strings.Fields(v)
	for i, w := range words {
		parts := strings.Split(w, "-")
		for j, part := range parts {
			if part == strings.ToUpper(part) && len(part) <= 3 {
				continue
			}
			r := []rune(strings.ToLower(part))
			if len(r) > 0 {
				r[0] = unicode.ToUpper(r[0])
			}
			parts[j] = string(r)
		}
		words[i] = strings.Join(parts, "-")
	}
	return strings.Join(words, " ")
}
func releaseStrip(v string) string {
	for _, re := range []*regexp.Regexp{resRE, sourceRE, codecRE, crcRE} {
		v = re.ReplaceAllString(v, " ")
	}
	return v
}
func atoi(v string) int { n, _ := strconv.Atoi(v); return n }
func normalizeResolution(v string) string {
	v = strings.ToLower(v)
	if v == "4k" {
		return "2160p"
	}
	return v
}
func normalizeSource(v string) string {
	s := strings.ToLower(strings.NewReplacer("-", "", "_", "", " ", "").Replace(v))
	switch s {
	case "bluray", "bdrip":
		return "BluRay"
	case "webdl":
		return "WEB-DL"
	case "webrip":
		return "WEBRip"
	case "web":
		return "WEB"
	case "hdtv":
		return "HDTV"
	case "dvdrip", "dvd":
		return "DVDRip"
	case "remux":
		return "Remux"
	}
	return v
}
func normalizeCodec(v string) string {
	s := strings.ToLower(strings.NewReplacer("-", "", "_", "", " ", "").Replace(v))
	switch s {
	case "x264", "h264":
		return "x264"
	case "x265", "h265", "hevc":
		return "x265"
	case "av1":
		return "AV1"
	}
	return v
}
func isReleaseTag(v string) bool {
	return resRE.MatchString(v) || sourceRE.MatchString(v) || codecRE.MatchString(v)
}
