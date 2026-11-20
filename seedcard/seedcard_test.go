package seedcard

import "testing"

func TestShippedConf(t *testing.T) {
	// Fails the build if text.conf has been edited into something invalid.
	if _, err := Parse(textConf); err != nil {
		t.Fatal(err)
	}
}

func TestParse(t *testing.T) {
	ok := "version_mark = sc # comment\ntitle = Bifb Treasury X\nseed_qr = No\n"
	s, err := Parse(ok)
	if err != nil {
		t.Fatal(err)
	}
	if s.VersionMark != "SC" || s.Title != "BIFB TREASURY X" || s.SeedQR {
		t.Fatalf("got %+v", s)
	}
	bad := map[string]string{
		"unsupported character": "version_mark = SC!\ntitle =\nseed_qr = no",
		"too long mark":            "version_mark = SEEDC\ntitle =\nseed_qr = no",
		"empty mark":               "version_mark =\ntitle =\nseed_qr = no",
		"long title":               "version_mark = SC\ntitle = ABCDEFGHIJKLMNOPQRS\nseed_qr = no",
		"bad bool":                 "version_mark = SC\ntitle =\nseed_qr = maybe",
		"unknown key":              "version_mark = SC\ntitle =\nseed_qr = no\nurl = X",
		"missing key":              "version_mark = SC\nseed_qr = no",
		"duplicate key":            "version_mark = SC\nversion_mark = SD\ntitle =\nseed_qr = no",
		"no equals":                "version_mark SC\ntitle =\nseed_qr = no",
		"accented":                 "version_mark = SÉ\ntitle =\nseed_qr = no",
	}
	for name, conf := range bad {
		if _, err := Parse(conf); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
