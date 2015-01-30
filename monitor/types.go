package monitor

import "time"

// ErrorManager holds
type ErrorManager struct {
	board  *Board
	thread *Thread
}

// Manager controls
type Manager struct {
	boards []*Board
	errBCh chan *Board
	errTCh chan *Thread
}

// Board is a set of threads.
type Board struct {
	Title      string
	URL        string
	Threadlist []*Thread
}

// Thread
type Thread struct {
	Title       string
	URL         string
	Comments    []Comment
	LastUpdated string
	Board       *Board
}

// Comment
type Comment struct {
	Name    string
	Email   string
	Date    time.Time
	ID      string
	BE      string
	Content string
}
