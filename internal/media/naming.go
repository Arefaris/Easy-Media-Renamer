package media

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"myproject/internal/provider"
)

var invalid = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
var reserved = regexp.MustCompile(`(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)

func Sanitize(name string) string {
	name = invalid.ReplaceAllString(name, "-")
	name = strings.TrimRight(strings.TrimSpace(name), ". ")
	if reserved.MatchString(name) {
		name = "_" + name
	}
	for len([]byte(name)) > 255 {
		_, size := utf8.DecodeLastRuneInString(name)
		name = name[:len(name)-size]
	}
	return strings.TrimRight(name, ". ")
}

func Render(tmpl string, ep provider.Episode, show provider.Show, ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	r := strings.NewReplacer(
		"{show}", show.Name, "{title}", ep.Title, "{year}", show.Year,
		"{season}", fmt.Sprintf("%02d", ep.Season), "{episode}", fmt.Sprintf("%02d", ep.Number), "{ext}", ext,
	)
	name := Sanitize(r.Replace(tmpl))
	if filepath.Ext(name) == "" && ext != "" {
		name += "." + ext
	}
	return name
}
