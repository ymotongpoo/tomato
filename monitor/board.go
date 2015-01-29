package monitor

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

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

var BoardURLException = map[string]bool{
	"http://www.2ch.net/":             true,
	"http://info.2ch.net/":            true,
	"http://irc.2ch.net:9090":         true,
	"http://search.2ch.net/":          true,
	"http://dig.2ch.net/":             true,
	"http://i.2ch.net/":               true,
	"http://www.2ch.net/kakolog.html": true,
	"http://notice.2ch.net/":          true,
	"http://be.2ch.net/":              true, // NOTE(ymotongpoo): need additional implmentation to login.
	"http://2ch.tora3.net/":           true,
}

var boardPath = xmlpath.MustCompile(`//font[@size="2"]/a`)

// HTTPGet call HTTP GET method to urlStr with custom header recommended to access 2ch.
func HTTPGet(urlStr string, gzipped bool) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	if gzipped {
		req.Header.Set("Accept-Encoding", "gzip")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return resp, err
		}
		resp.Body = r
	}
	return resp, err
}

// FetchBBSMenu returns
func FetchBBSMenu() (io.Reader, error) {
	var resp *http.Response
	var err error
	for _, m := range BBSMenu {
		resp, err = HTTPGet(m, true)
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

// ParseBBSMenu parase BBS Menu data stored in r. Data stored in r are expected to be UTF-8,
// so decode 2ch's default encoding (ShiftJIS) in advance.
func ParseBBSMenu(r io.Reader) ([]*Board, error) {
	root, err := xmlpath.ParseHTML(r)
	if err != nil {
		return nil, err
	}
	iter := boardPath.Iter(root)

	boards := []*Board{}
	alink := xmlpath.MustCompile(`@href`)
	atext := xmlpath.MustCompile(`text()`)
	for iter.Next() {
		n := iter.Node()
		b := Board{}
		if s, ok := alink.String(n); ok {
			if BoardURLException[s] {
				continue
			}
			b.URL = s
		}
		if s, ok := atext.String(n); ok {
			b.Title = s
		}
		boards = append(boards, &b)
	}
	return boards, nil
}

// FetchThreadlist returns
func (b *Board) FetchThreadlist() (io.Reader, error) {
	resp, err := HTTPGet(b.URL+SubjectFile, true)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ParseThreadlist read subject.txt and make a list of thread in b.
// Data stored in r is expected to be UTF-8.
func (b *Board) ParseThreadlist(r io.Reader) error {
	br := bufio.NewReader(r)
	var err error
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}
		pairs := strings.SplitN(line, "<>", 2)
		if len(pairs) != 2 {
			return fmt.Errorf("Wrongly formatted data: %v", line)
		}
		t := &Thread{
			Title: strings.TrimSpace(pairs[1]),
			URL:   b.URL + "/dat/" + pairs[0],
			Board: b,
		}
		b.Threadlist = append(b.Threadlist, t)
	}
	if err != io.EOF {
		return err
	}
	return nil
}
