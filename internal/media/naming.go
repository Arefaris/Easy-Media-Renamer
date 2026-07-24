package media

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"myproject/internal/provider"
)

var invalidSegment = regexp.MustCompile(`[<>:"\\|?*\x00-\x1f]`)
var invalidTokenValue = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
var reserved = regexp.MustCompile(`(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)
var tokenRE = regexp.MustCompile(`\{[a-zA-Z0-9]+\}`)
var conditionalRE = regexp.MustCompile(`\[([^\[\]]*)\]`)

type NamingData struct {
	Episode                                                                 provider.Episode
	Show                                                                    provider.Show
	Ext, Filename                                                           string
	Source, Group, CRC32                                                    string
	VideoFormat, VideoCodec, AudioCodec, Channels, HD, Resolution, Duration string
}

func Sanitize(name string) string { return sanitizeSegment(name) }
func sanitizeSegment(name string) string {
	name = invalidSegment.ReplaceAllString(name, "-")
	name = strings.TrimRight(strings.TrimSpace(name), ". ")
	if reserved.MatchString(name) {
		name = "_" + name
	}
	for len([]byte(name)) > 255 {
		_, n := utf8.DecodeLastRuneInString(name)
		name = name[:len(name)-n]
	}
	return strings.TrimRight(name, ". ")
}

func Render(tmpl string, ep provider.Episode, show provider.Show, ext string) string {
	return RenderTemplate(tmpl, NamingData{Episode: ep, Show: show, Ext: ext})
}
func RenderTemplate(tmpl string, d NamingData) string {
	d.Ext = strings.TrimPrefix(d.Ext, ".")
	hasExtToken := strings.Contains(tmpl, "{ext}")
	const escapedOpen, escapedClose = "\uE000", "\uE001"
	tmpl = strings.ReplaceAll(strings.ReplaceAll(tmpl, `\[`, escapedOpen), `\]`, escapedClose)
	values := templateValues(d)
	for token, value := range values {
		values[token] = sanitizeTokenValue(value)
	}
	for {
		next := conditionalRE.ReplaceAllStringFunc(tmpl, func(block string) string {
			inner := block[1 : len(block)-1]
			tokens := tokenRE.FindAllString(inner, -1)
			for _, tok := range tokens {
				if values[tok] == "" {
					return ""
				}
			}
			return replaceTokens(inner, values)
		})
		if next == tmpl {
			break
		}
		tmpl = next
	}
	tmpl = replaceTokens(tmpl, values)
	tmpl = strings.ReplaceAll(strings.ReplaceAll(tmpl, escapedOpen, "["), escapedClose, "]")
	parts := strings.Split(strings.ReplaceAll(tmpl, "\\", "/"), "/")
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if s := sanitizeSegment(part); s != "" && s != "." && s != ".." {
			clean = append(clean, s)
		}
	}
	name := filepath.Join(clean...)
	if !hasExtToken && d.Ext != "" && !strings.HasSuffix(strings.ToLower(name), "."+strings.ToLower(d.Ext)) {
		name += "." + d.Ext
	}
	return name
}
func sanitizeTokenValue(value string) string {
	return strings.TrimRight(strings.TrimSpace(invalidTokenValue.ReplaceAllString(value, "-")), ". ")
}
func templateValues(d NamingData) map[string]string {
	e, s := d.Episode.Number, d.Episode.Season
	abs := ""
	if d.Episode.Absolute > 0 {
		abs = strconv.Itoa(d.Episode.Absolute)
	}
	return map[string]string{
		"{n}": d.Show.Name, "{show}": d.Show.Name, "{t}": d.Episode.Title, "{title}": d.Episode.Title, "{y}": d.Show.Year, "{year}": d.Show.Year,
		"{s}": strconv.Itoa(s), "{e}": strconv.Itoa(e), "{season}": fmt.Sprintf("%02d", s), "{episode}": fmt.Sprintf("%02d", e), "{s00}": fmt.Sprintf("%02d", s), "{e00}": fmt.Sprintf("%02d", e), "{s00e00}": fmt.Sprintf("S%02dE%02d", s, e), "{sxe}": fmt.Sprintf("%dx%02d", s, e),
		"{airdate}": d.Episode.AirDate, "{abs}": abs, "{ext}": d.Ext, "{fn}": d.Filename, "{vf}": d.VideoFormat, "{vc}": d.VideoCodec, "{ac}": d.AudioCodec, "{channels}": d.Channels, "{hd}": d.HD, "{resolution}": d.Resolution, "{duration}": d.Duration, "{source}": d.Source, "{group}": d.Group, "{crc32}": d.CRC32,
	}
}
func replaceTokens(s string, v map[string]string) string {
	return tokenRE.ReplaceAllStringFunc(s, func(k string) string { return v[k] })
}
func NeedsMediaInfo(tmpl string) bool {
	for _, t := range []string{"{vf}", "{vc}", "{ac}", "{channels}", "{hd}", "{resolution}", "{duration}"} {
		if strings.Contains(tmpl, t) {
			return true
		}
	}
	return false
}
