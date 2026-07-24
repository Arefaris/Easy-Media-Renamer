package mediainfo

import (
	"path/filepath"
	"strings"
	"time"
)

type Info struct {
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	VideoCodec string        `json:"video_codec"`
	AudioCodec string        `json:"audio_codec"`
	Channels   int           `json:"channels"`
	HDR        bool          `json:"hdr"`
	Duration   time.Duration `json:"duration"`
	AudioLang  string        `json:"audio_lang"`
}

func Probe(path string) (Info, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mkv", ".webm":
		return probeMKV(path)
	case ".mp4", ".m4v", ".mov":
		return probeMP4(path)
	default:
		return Info{}, nil
	}
}
func VideoFormat(i Info) string {
	if i.Height > 0 {
		return formatHeight(i.Height)
	}
	return ""
}
func formatHeight(h int) string {
	switch {
	case h >= 2000:
		return "2160p"
	case h >= 1000:
		return "1080p"
	case h >= 700:
		return "720p"
	case h >= 560:
		return "576p"
	case h >= 470:
		return "480p"
	}
	return ""
}
