package monitor

import (
	"fmt"
	"os"
	"path/filepath"
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
		boards:    make([]*Board, BoardCapacity),
		errBCh:    make(chan *Board, BoardCapacity),
		errTCh:    make(chan *Thread, ThreadCapacity),
	}, nil
}

func (m Manager) Start() {
	// TODO(ymotongpoo): Find datastore first and board data.
	// TODO(ymotongpoo): Fetch board data if possible and update URL if there are.
	// TODO(ymotongpoo): Load threadlist data from datastore and check timestamp.
	// TODO(ymotongpoo): Fetch threadlist data and update subject.txt saved.
	// TODO(ymotongpoo): Load thread data from datastore.
	// TODO(ymotongpoo): Fetch thread data. Be sure to range update using if-modified-since header.
}
