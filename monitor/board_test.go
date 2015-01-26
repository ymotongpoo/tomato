package monitor

import (
	"os"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestParseBBSMenu(t *testing.T) {
	file, err := os.Open("data/bbsmenu.html")
	if err != nil {
		t.Errorf("%v", err)
	}
	rInUTF8 := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
	boards, err := ParseBBSMenu(rInUTF8)
	if err != nil {
		t.Errorf("Error on parsing BBS Menu: %v", err)
	}
	for i, b := range boards {
		t.Logf("%v: %v\t%v", i, b.Title, b.URL)
		if b.Title == "" || b.URL == "" {
			t.Errorf("%vth content has empty content", i)
		}
	}
}

func TestParseThreadlist(t *testing.T) {
	file, err := os.Open("data/subject.txt")
	if err != nil {
		t.Errorf("%v", err)
	}
	rInUTF8 := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
	b := &Board{
		Title: "ニュース速報",
		URL:   "http://hayabusa3.2ch.net/news/",
	}
	err = b.ParseThreadlist(rInUTF8)
	if err != nil {
		t.Errorf("ParseThreadlist had error: %v", err)
	}
	if len(b.Threadlist) == 0 {
		t.Errorf("No thread parsed.")
	}
	for i, thread := range b.Threadlist {
		t.Logf("%v: %v --> %v", i, thread.Title, thread.URL)
	}
}
