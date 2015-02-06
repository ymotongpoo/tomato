package monitor

import (
	"fmt"
	"log"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/html/charset"
)

var (
	BoardCapacity  = 1000
	ThreadCapacity = 50000
)

func (e ErrorManager) Error() string {
	if e.board != nil {
		return fmt.Sprintf("Board %v is error. Retry URL %v.", e.board.Title, e.board.URL)
	}
	if e.thread != nil {
		return fmt.Sprintf("Thread %v is error. Retry URL %v.", e.thread.Title, e.thread.URL)
	}
	return "Someting is wrong. ErrorManager should hold either of board or thread."
}

// NewManager generates new Manager with datastore path.
func NewManager() (*Manager, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	datastore := filepath.Join(pwd, "files")
	return &Manager{
		datastore: datastore,
		boards:    nil,
		errBCh:    make(chan *Board, BoardCapacity),
		errTCh:    make(chan *Thread, ThreadCapacity),
	}, nil
}

func (m Manager) Start() {
	// Find datastore first and board data.
	p := filepath.Join(m.datastore, BBSMenuFile)
	var bbsmenu io.Reader
	var fetcherr error
	bbsmenu, err := os.Open(p)
	if err != nil {
		if err == os.ErrNotExist {
			file, err := os.Create(p)
			if err != nil {
				log.Fatal(err)
			}
			bbsmenu, fetcherr = FetchBBSMenu()
			if fetcherr != nil {
				log.Fatal(fetcherr)
			}
			if _, err := io.Copy(file, bbsmenu); err != nil {
				log.Fatal(err)
			}
		}
		log.Fatal(err)
	}

	r, err := charset.NewReader(bbsmenu, "text/html")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch board data if possible and update URL if there are.
	// TODO(ymotongpoo): confirm file timestamp and last updated header of the bbsmenu.
	m.boards, err = ParseBBSMenu(r)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(ymotongpoo): Load threadlist data from datastore and check timestamp.
	for _, b := range m.boards {
		tr, err := b.FetchThreadlist()
		if err != nil {
			log.Printf("error on fetching threadlist of %v: %v", b.URL, err)
			continue
		}
		trInUTF8, err := charset.NewReader(tr, "text/plain")
		if err != nil {
			log.Println(err)
		}
		err = b.ParseThreadlist(trInUTF8)
		if err != nil {
			log.Printf("error on parsing threadlist of %v: %v", b.URL, err)
			continue
		}
		for _, t := range b.Threadlist {
			fmt.Printf("%v: %v -> %v\n", b.Title, t.Title, t.URL)
		}
		time.Sleep(3 * time.Second)
	}

	// TODO(ymotongpoo): Fetch threadlist data and update subject.txt saved.
	// TODO(ymotongpoo): Load thread data from datastore.
	// TODO(ymotongpoo): Fetch thread data. Be sure to range update using if-modified-since header.
}
