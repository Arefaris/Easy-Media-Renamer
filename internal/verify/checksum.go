package verify

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Result struct {
	File     string `json:"file"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	OK       bool   `json:"ok"`
	Error    string `json:"error,omitempty"`
}

func Compute(path, algo string) (string, error) {
	var h hash.Hash
	switch strings.ToLower(algo) {
	case "crc32", "sfv":
		h = crc32.NewIEEE()
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	default:
		return "", fmt.Errorf("unsupported checksum %q", algo)
	}
	f, e := os.Open(path)
	if e != nil {
		return "", e
	}
	defer f.Close()
	if _, e = io.Copy(h, f); e != nil {
		return "", e
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}
func WriteSFV(output string, files []string) error {
	f, e := os.Create(output)
	if e != nil {
		return e
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, p := range files {
		sum, e := Compute(p, "crc32")
		if e != nil {
			return e
		}
		if _, e = fmt.Fprintf(w, "%s %s\n", filepath.Base(p), sum); e != nil {
			return e
		}
	}
	return w.Flush()
}
func VerifySFV(path string) ([]Result, error) { return verifyManifest(path, "crc32", false) }
func WriteManifest(output, algo string, files []string) error {
	f, e := os.Create(output)
	if e != nil {
		return e
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, p := range files {
		sum, e := Compute(p, algo)
		if e != nil {
			return e
		}
		if _, e = fmt.Fprintf(w, "%s *%s\n", strings.ToLower(sum), filepath.Base(p)); e != nil {
			return e
		}
	}
	return w.Flush()
}
func VerifyManifest(path, algo string) ([]Result, error) { return verifyManifest(path, algo, true) }
func verifyManifest(path, algo string, hashFirst bool) ([]Result, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	base := filepath.Dir(path)
	out := []Result{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		var name, expected string
		if hashFirst {
			p := strings.Fields(line)
			if len(p) < 2 {
				continue
			}
			expected = p[0]
			name = strings.TrimLeft(strings.Join(p[1:], " "), "*")
		} else {
			idx := strings.LastIndex(line, " ")
			if idx < 1 {
				continue
			}
			name = strings.TrimSpace(line[:idx])
			expected = strings.TrimSpace(line[idx+1:])
		}
		actual, x := Compute(filepath.Join(base, name), algo)
		r := Result{File: name, Expected: strings.ToUpper(expected), Actual: actual, OK: x == nil && strings.EqualFold(expected, actual)}
		if x != nil {
			r.Error = x.Error()
		}
		out = append(out, r)
	}
	return out, s.Err()
}

var embeddedCRC = regexp.MustCompile(`(?i)\[([A-F0-9]{8})\]`)

func VerifyEmbeddedCRC(path string) (Result, error) {
	m := embeddedCRC.FindStringSubmatch(filepath.Base(path))
	if m == nil {
		return Result{}, fmt.Errorf("filename has no embedded CRC32")
	}
	actual, e := Compute(path, "crc32")
	r := Result{File: filepath.Base(path), Expected: strings.ToUpper(m[1]), Actual: actual, OK: e == nil && strings.EqualFold(m[1], actual)}
	if e != nil {
		r.Error = e.Error()
	}
	return r, e
}
