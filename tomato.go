package main

import (
	"fmt"
	"log"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/ymotongpoo/tomato/monitor"
)

func main() {
	log.Println("start")
	r, err := monitor.FetchBBSMenu()
	if err != nil {
		panic(err)
	}
	rInUTF8 := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	boards, err := monitor.ParseBBSMenu(rInUTF8)
	if err != nil {
		panic(err)
	}
	for i, b := range boards {
		fmt.Printf("%v: %v\t%v\n", i, b.Title, b.URL)
	}
	log.Println("end")
}
