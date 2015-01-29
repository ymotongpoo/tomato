package monitor

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (t *Thread) FetchDatData(w io.Writer) error {
	resp, err := HTTPGet(t.URL, true)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Body)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// ParseDatData parses data stored in r to Thread. r should be UTF-8 encoded.
// Specification is written here:
// http://info.2ch.net/index.php/Monazilla/develop/dat
func (t *Thread) ParseDatData(r io.Reader) error {
	br := bufio.NewReader(r)
	var err error
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}
		contents := strings.Split(line, "<>")
		if len(contents) < 5 {
			return fmt.Errorf("comment data is broken: %v", line)
		}
		dateIDBE := strings.Split(contents[2], " ")
		c := Comment{
			Name:    contents[0],
			Email:   contents[1],
			Date:    dateIDBE[0] + " " + dateIDBE[1],
			ID:      strings.TrimPrefix(dateIDBE[2], "ID:"),
			Content: contents[3],
		}
		t.Comments = append(t.Comments, c)
	}
	if err != io.EOF {
		return err
	}
	return nil
}
