package monitor

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// FetchDatData fetch dat file from specified URL and write its data to w.
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
		datetime := t.convert2chTime(dateIDBE[0] + " " + dateIDBE[1])
		c := Comment{
			Name:    contents[0],
			Email:   contents[1],
			Date:    datetime,
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

func (t *Thread) convert2chTime(datetime string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	if datetime == "あぼーん" {
		return time.Date(1900, 1, 1, 0, 0, 0, 0, loc)
	}
	i := strings.Index(datetime, "(")
	j := strings.Index(datetime, ")")
	parsed, _ := time.ParseInLocation("2006/01/02 15:04:05.00", datetime[:i]+" "+datetime[j+1:], loc)
	return parsed
}
