package monitor

import (
	"fmt"
	"log"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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
	err = os.Mkdir(datastore, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
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
		file, err := os.Create(p)
		if err != nil {
			log.Fatal("file creation:", err)
		}
		bbsmenu, fetcherr = FetchBBSMenu()
		if fetcherr != nil {
			log.Fatal(fetcherr)
		}
		if _, err := io.Copy(file, bbsmenu); err != nil {
			log.Fatal(err)
		}
	}

	r, err := charset.NewReader(bbsmenu, "text/html")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch board data if possible and update URL if there are.
	// TODO(ymotongpoo): confirm file timestamp and last updated header of the bbsmenu.
	m.boards, err = ParseBBSMenu(r)
	if err != nil {
		log.Fatal("Parse: ", err)
	}

	// TODO(ymotongpoo): Load threadlist data from datastore and check timestamp.
	for _, b := range m.boards {
		log.Printf("%v: %v", b.Title, b.URL)
		tr, err := b.FetchThreadlist()
		if err != nil {
			log.Printf("error on fetching threadlist of %v: %v", b.URL, err)
			continue
		}

		trInUTF8 := transform.NewReader(tr, japanese.ShiftJIS.NewDecoder()) // TODO(ymotongpoo): find more generic way.
		if err != nil {
			log.Println("NewReader:", err)
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
}
