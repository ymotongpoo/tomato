package monitor

import (
	"io"
	"net/http"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/xmlpath.v2"
)

const (
	UserAgent = "Monazilla/1.00 (Tomato/0.0.1)"
)

var BBSMenu = []string{
	`http://menu.2ch.net/bbsmenu.html`,
	`http://www.zonubbs.net/bbsmenu.html`,
	`http://azlucky.s25.xrea.com/2chboard/bbsmenu.html`,
}

// FetchBBSMenu returns
func FetchBBSMenu() (io.Reader, error) {
	var resp *http.Response
	var err error
	for _, m := range BBSMenu {
		resp, err = http.Get(m)
		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func ParseBBSManu(r io.Reader) ([]Board, error) {
	rInUTF8 := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	root, err := xmlpath.ParseHTML(rInUTF8)
	if err != nil {
		return nil, err
	}
	// TODO(ymotongpoo): implement here.
	_ = root
	return nil, nil // mock
}
