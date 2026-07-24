package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKnownVectors(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "data.txt")
	if e := os.WriteFile(p, []byte("123456789"), 0600); e != nil {
		t.Fatal(e)
	}
	want := map[string]string{"crc32": "CBF43926", "md5": "25F9E794323B453885F5181F1B624D0B", "sha1": "F7C3BC1D808E04732ADF679965CCC34CA7AE3441", "sha256": "15E2B0D3C33891EBB0F1EF609EC419420C20E320CE94C65FBC8C3312448EB225"}
	for a, w := range want {
		g, e := Compute(p, a)
		if e != nil || g != w {
			t.Errorf("%s %s %v", a, g, e)
		}
	}
}
func TestSFVAndEmbedded(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "file [CBF43926].txt")
	if e := os.WriteFile(p, []byte("123456789"), 0600); e != nil {
		t.Fatal(e)
	}
	r, e := VerifyEmbeddedCRC(p)
	if e != nil || !r.OK {
		t.Fatalf("%v %#v", e, r)
	}
	sfv := filepath.Join(d, "files.sfv")
	if e = WriteSFV(sfv, []string{p}); e != nil {
		t.Fatal(e)
	}
	rs, e := VerifySFV(sfv)
	if e != nil || len(rs) != 1 || !rs[0].OK {
		t.Fatalf("%v %#v", e, rs)
	}
}
