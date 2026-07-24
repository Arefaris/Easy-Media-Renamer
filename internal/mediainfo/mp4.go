package mediainfo

import (
	mp4 "github.com/abema/go-mp4"
	"os"
	"time"
)

func probeMP4(path string) (Info, error) {
	f, e := os.Open(path)
	if e != nil {
		return Info{}, e
	}
	defer f.Close()
	out := Info{}
	_, e = mp4.ReadBoxStructure(f, func(h *mp4.ReadHandle) (interface{}, error) {
		typ := h.BoxInfo.Type
		switch typ {
		case mp4.BoxTypeAvc1():
			out.VideoCodec = "x264"
		case mp4.BoxTypeHev1(), mp4.BoxTypeHvc1():
			out.VideoCodec = "x265"
		case mp4.BoxTypeAv01():
			out.VideoCodec = "AV1"
		case mp4.BoxTypeMp4a():
			out.AudioCodec = "AAC"
		case mp4.BoxTypeAC3():
			out.AudioCodec = "AC3"
		}
		if !h.BoxInfo.IsSupportedType() {
			return nil, nil
		}
		box, _, x := h.ReadPayload()
		if x != nil {
			return nil, x
		}
		switch b := box.(type) {
		case *mp4.Mvhd:
			if b.Timescale > 0 {
				out.Duration = time.Duration(float64(b.GetDuration()) / float64(b.Timescale) * float64(time.Second))
			}
		case *mp4.Tkhd:
			w, h := int(b.GetWidth()), int(b.GetHeight())
			if w > out.Width {
				out.Width = w
			}
			if h > out.Height {
				out.Height = h
			}
		case *mp4.VisualSampleEntry:
			if int(b.Width) > out.Width {
				out.Width = int(b.Width)
			}
			if int(b.Height) > out.Height {
				out.Height = int(b.Height)
			}
		case *mp4.AudioSampleEntry:
			if int(b.ChannelCount) > out.Channels {
				out.Channels = int(b.ChannelCount)
			}
		}
		return h.Expand()
	})
	return out, e
}
