package monitor

import (
	"io"
	"net/http"

	"gopkg.in/xmlpath.v2"
)

const (
	UserAgent   = "Monazilla/1.00 (Tomato/0.0.1)"
	SubjectFile = "subject.txt"
)

var BBSMenu = []string{
	`http://menu.2ch.net/bbsmenu.html`,
	`http://www.zonubbs.net/bbsmenu.html`,
	`http://azlucky.s25.xrea.com/2chboard/bbsmenu.html`,
}

var boardPath = xmlpath.MustCompile(`//font[@size="2"]/a`)

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

// ParseBBSMenu parase BBS Menu data stored in r in ShiftJIS.
func ParseBBSMenu(r io.Reader) ([]Board, error) {
	root, err := xmlpath.ParseHTML(r)
	if err != nil {
		return nil, err
	}
	iter := boardPath.Iter(root)

	boards := []Board{}
	alink := xmlpath.MustCompile(`@href`)
	atext := xmlpath.MustCompile(`text()`)
	for iter.Next() {
		n := iter.Node()
		b := Board{}
		if s, ok := alink.String(n); ok {
			b.URL = s
		}
		if s, ok := atext.String(n); ok {
			b.Title = s
		}
		boards = append(boards, b)
	}
	return boards, nil
}

// ParseThreadList read subject.txt and make a list of thread in b.
func (b Board) ParseThreadList(io.Reader) error {
	return nil // mock
}
