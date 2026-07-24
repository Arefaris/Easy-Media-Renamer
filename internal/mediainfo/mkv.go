package mediainfo

import (
	mkv "github.com/remko/go-mkvparse"
	"strings"
	"time"
)

type mkvHandler struct {
	mkv.DefaultHandler
	info          Info
	trackType     int64
	timecodeScale int64
}

func probeMKV(path string) (Info, error) {
	h := &mkvHandler{timecodeScale: 1_000_000}
	e := mkv.ParsePath(path, h)
	return h.info, e
}
func (h *mkvHandler) HandleMasterBegin(id mkv.ElementID, info mkv.ElementInfo) (bool, error) {
	if id == mkv.TrackEntryElement {
		h.trackType = 0
	}
	if id == mkv.ClusterElement {
		return false, nil
	}
	return true, nil
}
func (h *mkvHandler) HandleInteger(id mkv.ElementID, v int64, _ mkv.ElementInfo) error {
	switch id {
	case mkv.TrackTypeElement:
		h.trackType = v
	case mkv.TimecodeScaleElement:
		h.timecodeScale = v
	case mkv.PixelWidthElement:
		if h.trackType == 1 {
			h.info.Width = int(v)
		}
	case mkv.PixelHeightElement:
		if h.trackType == 1 {
			h.info.Height = int(v)
		}
	case mkv.ChannelsElement:
		if h.trackType == 2 {
			h.info.Channels = int(v)
		}
	case mkv.TransferCharacteristicsElement:
		if v == 16 || v == 18 {
			h.info.HDR = true
		}
	}
	return nil
}
func (h *mkvHandler) HandleFloat(id mkv.ElementID, v float64, _ mkv.ElementInfo) error {
	if id == mkv.DurationElement {
		h.info.Duration = time.Duration(v * float64(h.timecodeScale))
	}
	return nil
}
func (h *mkvHandler) HandleString(id mkv.ElementID, v string, _ mkv.ElementInfo) error {
	switch id {
	case mkv.CodecIDElement:
		if h.trackType == 1 {
			h.info.VideoCodec = mapCodec(v)
		} else if h.trackType == 2 && h.info.AudioCodec == "" {
			h.info.AudioCodec = mapCodec(v)
		}
	case mkv.LanguageElement, mkv.LanguageIETFElement:
		if h.trackType == 2 && h.info.AudioLang == "" {
			h.info.AudioLang = v
		}
	}
	return nil
}
func mapCodec(v string) string {
	s := strings.ToUpper(v)
	switch {
	case strings.Contains(s, "HEVC"), strings.Contains(s, "H265"):
		return "x265"
	case strings.Contains(s, "AVC"), strings.Contains(s, "H264"):
		return "x264"
	case strings.Contains(s, "AV1"):
		return "AV1"
	case strings.Contains(s, "DTS"):
		return "DTS"
	case strings.Contains(s, "AC3"), strings.Contains(s, "EAC3"):
		return "AC3"
	case strings.Contains(s, "AAC"):
		return "AAC"
	case strings.Contains(s, "FLAC"):
		return "FLAC"
	case strings.Contains(s, "OPUS"):
		return "Opus"
	}
	return v
}
