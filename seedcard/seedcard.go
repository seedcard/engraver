// Package seedcard holds the SeedCard engraving text, read from text.conf.
//
// text.conf is embedded at build time, so an image always engraves the same
// text. Invalid settings make the program refuse to start (and make
// `go test ./...` fail), rather than engraving something unexpected.
package seedcard

import (
	_ "embed"
	"fmt"
	"strings"

	"seedhammer.com/font/constant"
)

//go:embed text.conf
var textConf string

// Settings are the user-editable parts of the engraving.
type Settings struct {
	// VersionMark replaces Seed Hammer's "V1" mark on the seed side.
	VersionMark string
	// Title, when not empty, replaces the wallet name on the seed side.
	Title string
	// SeedQR enables the SeedQR on the seed side.
	SeedQR bool
}

const (
	MaxVersionMark = 4
	MaxTitle       = 18 // matches backup.MaxTitleLen
)

// Text is the parsed, validated content of text.conf.
var Text = mustParse(textConf)

func mustParse(conf string) Settings {
	s, err := Parse(conf)
	if err != nil {
		panic(fmt.Sprintf("seedcard/text.conf: %v", err))
	}
	return s
}

// Parse reads "key = value" lines; '#' starts a comment.
func Parse(conf string) (Settings, error) {
	s := Settings{}
	seen := map[string]bool{}
	for n, line := range strings.Split(conf, "\n") {
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			return s, fmt.Errorf("line %d: expected \"name = value\"", n+1)
		}
		key = strings.ToLower(strings.TrimSpace(key))
		val = strings.TrimSpace(val)
		if seen[key] {
			return s, fmt.Errorf("line %d: %q set twice", n+1, key)
		}
		seen[key] = true
		switch key {
		case "version_mark":
			v, err := engravable(val)
			if err != nil {
				return s, fmt.Errorf("line %d: version_mark: %v", n+1, err)
			}
			if len(v) < 1 || len(v) > MaxVersionMark {
				return s, fmt.Errorf("line %d: version_mark must be 1 to %d characters", n+1, MaxVersionMark)
			}
			s.VersionMark = v
		case "title":
			v, err := engravable(val)
			if err != nil {
				return s, fmt.Errorf("line %d: title: %v", n+1, err)
			}
			if len(v) > MaxTitle {
				return s, fmt.Errorf("line %d: title must be at most %d characters", n+1, MaxTitle)
			}
			s.Title = v
		case "seed_qr":
			switch strings.ToLower(val) {
			case "yes":
				s.SeedQR = true
			case "no":
				s.SeedQR = false
			default:
				return s, fmt.Errorf("line %d: seed_qr must be yes or no", n+1)
			}
		default:
			return s, fmt.Errorf("line %d: unknown setting %q", n+1, key)
		}
	}
	for _, k := range []string{"version_mark", "title", "seed_qr"} {
		if !seen[k] {
			return s, fmt.Errorf("missing setting %q", k)
		}
	}
	return s, nil
}

// engravable upper-cases s and checks every character exists in the
// engraving font.
func engravable(s string) (string, error) {
	s = strings.ToUpper(s)
	for _, r := range s {
		if r == ' ' {
			continue
		}
		if _, _, ok := constant.Font.Decode(r); !ok {
			return "", fmt.Errorf("character %q can't be engraved", r)
		}
	}
	return s, nil
}
