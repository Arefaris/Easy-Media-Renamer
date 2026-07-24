package mediainfo

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProbeUnsupportedContainerIsEmpty(t *testing.T) {
	info, err := Probe(filepath.Join(t.TempDir(), "video.avi"))
	if err != nil || info != (Info{}) {
		t.Fatalf("info=%#v err=%v", info, err)
	}
}

func TestProbeSyntheticMKV(t *testing.T) {
	// EBML header + Segment containing Info and one AVC video TrackEntry.
	raw, err := hex.DecodeString("1A45DFA39D4286810142F7810142F2810442F381084282886D6174726F736B61428781044285810218538067C51549A966922AD7B1830F4240448988408F4000000000001654AE6BA9AEA7D7810173C58101838101868F565F4D504547342F49534F2F415643E08AB0820780BA820438")
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "minimal.mkv")
	if err = os.WriteFile(p, raw, 0600); err != nil {
		t.Fatal(err)
	}
	info, err := Probe(p)
	if err != nil {
		t.Fatal(err)
	}
	if info.Width != 1920 || info.Height != 1080 || info.VideoCodec != "x264" || info.Duration != time.Second {
		t.Fatalf("%#v", info)
	}
}

func TestProbeSyntheticMP4(t *testing.T) {
	var payload bytes.Buffer
	payload.Write([]byte{0, 0, 0, 0}) // version + flags
	_ = binary.Write(&payload, binary.BigEndian, uint32(0))
	_ = binary.Write(&payload, binary.BigEndian, uint32(0))
	_ = binary.Write(&payload, binary.BigEndian, uint32(1000))
	_ = binary.Write(&payload, binary.BigEndian, uint32(2500))
	_ = binary.Write(&payload, binary.BigEndian, uint32(0x00010000))
	_ = binary.Write(&payload, binary.BigEndian, uint16(0x0100))
	payload.Write(make([]byte, 10))
	payload.Write([]byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x40, 0, 0, 0})
	payload.Write(make([]byte, 24))
	_ = binary.Write(&payload, binary.BigEndian, uint32(2))
	mvhd := mp4Box("mvhd", payload.Bytes())
	raw := mp4Box("moov", mvhd)
	p := filepath.Join(t.TempDir(), "minimal.mp4")
	if err := os.WriteFile(p, raw, 0600); err != nil {
		t.Fatal(err)
	}
	info, err := Probe(p)
	if err != nil {
		t.Fatal(err)
	}
	if info.Duration != 2500*time.Millisecond {
		t.Fatalf("%#v", info)
	}
}

func mp4Box(kind string, payload []byte) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.BigEndian, uint32(len(payload)+8))
	b.WriteString(kind)
	b.Write(payload)
	return b.Bytes()
}

func TestVideoFormat(t *testing.T) {
	for _, tc := range []struct {
		height int
		want   string
	}{{2160, "2160p"}, {1080, "1080p"}, {720, "720p"}, {480, "480p"}, {360, ""}} {
		if got := VideoFormat(Info{Height: tc.height}); got != tc.want {
			t.Errorf("%d: %q", tc.height, got)
		}
	}
}
