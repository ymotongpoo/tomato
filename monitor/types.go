package monitor

import "time"

// Category is a set of boards.
type Category struct {
	Title     string
	Boardlist []*Board
}

// Board is a set of threads.
type Board struct {
	Title      string
	URL        string
	Threadlist []*Thread
	Category   *Category
}

// Thread
type Thread struct {
	Title       string
	URL         string
	Comments    []Comment
	LastUpdated time.Time
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
