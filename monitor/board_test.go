package monitor

import (
	"os"
	"testing"
)

func TestParseBBSMenu(t *testing.T) {
	file, err := os.Open("data/bbsmenu.html")
	if err != nil {
		t.Errorf("%v", err)
	}
	boards, err := ParseBBSMenu(file)
	if err != nil {
		t.Errorf("Error on parsing BBS Menu: %v", err)
	}
	for i, b := range boards {
		if b.Title == "" || b.URL == "" {
			t.Errorf("%vth content has empty content", i)
		}
	}
}
