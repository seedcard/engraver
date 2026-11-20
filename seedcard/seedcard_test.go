package seedcard

import "testing"

func TestShippedConf(t *testing.T) {
	// Fails the build if text.conf has been edited into something invalid.
	if _, err := Parse(textConf); err != nil {
		t.Fatal(err)
	}
}

func TestParse(t *testing.T) {
	ok := "version_mark = sc # comment\ntitle = Bifb Treasury X\nkey_number = 1/1\nseed_qr = No\n"
	s, err := Parse(ok)
	if err != nil {
		t.Fatal(err)
	}
	if s.VersionMark != "SC" || s.Title != "BIFB TREASURY X" || s.SeedQR || s.KeyNumber != "1/1" {
		t.Fatalf("got %+v", s)
	}
	bad := map[string]string{
		"unsupported character": "version_mark = SC!\ntitle =\nkey_number = auto\nseed_qr = no",
		"too long mark":         "version_mark = SEEDC\ntitle =\nkey_number = auto\nseed_qr = no",
		"empty mark":            "version_mark =\ntitle =\nkey_number = auto\nseed_qr = no",
		"long title":            "version_mark = SC\ntitle = ABCDEFGHIJKLMNOPQRS\nkey_number = auto\nseed_qr = no",
		"bad bool":              "version_mark = SC\ntitle =\nkey_number = auto\nseed_qr = maybe",
		"unknown key":           "version_mark = SC\ntitle =\nkey_number = auto\nseed_qr = no\nurl = X",
		"missing key":           "version_mark = SC\nkey_number = auto\nseed_qr = no",
		"duplicate key":         "version_mark = SC\nversion_mark = SD\ntitle =\nkey_number = auto\nseed_qr = no",
		"no equals":             "version_mark SC\ntitle =\nkey_number = auto\nseed_qr = no",
		"long key number":       "version_mark = SC\ntitle =\nkey_number = 10/100\nseed_qr = no",
		"empty key number":      "version_mark = SC\ntitle =\nkey_number =\nseed_qr = no",
		"accented":              "version_mark = SÉ\ntitle =\nkey_number = auto\nseed_qr = no",
	}
	for name, conf := range bad {
		if _, err := Parse(conf); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestKeyNumberAuto(t *testing.T) {
	s, err := Parse("version_mark = SC\ntitle =\nkey_number = Auto\nseed_qr = no")
	if err != nil {
		t.Fatal(err)
	}
	if s.KeyNumber != "" {
		t.Fatalf("auto should leave KeyNumber empty, got %q", s.KeyNumber)
	}
}
