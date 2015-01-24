package monitor

import "time"

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
	LastUpdated time.Time
	Board       *Board
}

// Comment
type Comment struct {
	Name    string
	Email   string
	Date    string
	ID      string
	BE      string
	Comment string
}
