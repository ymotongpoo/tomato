package monitor

import (
	"os"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestParseDatData(t *testing.T) {
	file, err := os.Open("data/1422369026.dat")
	if err != nil {
		t.Error(err)
	}
	rInUTF8 := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
	th := &Thread{}

	err = th.ParseDatData(rInUTF8)
	if err != nil {
		t.Error(err)
	}
	for i, c := range th.Comments {
		t.Logf("[%v] %v %v %v ID:%v\n%v\n", i, c.Name, c.Email, c.Date, c.ID, c.Content)
	}
}
