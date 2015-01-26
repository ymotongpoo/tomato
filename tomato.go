package main

import (
	"fmt"

	"github.com/ymotongpoo/tomato/monitor"
)

func main() {
	r, err := monitor.FetchBBSMenu()
	if err != nil {
		panic(err)
	}
	boards, err := monitor.ParseBBSMenu(r)
	for i, b := range boards {
		fmt.Printf("%v: %v\t%v\n", i, b.Title, b.URL)
	}
}
