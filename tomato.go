package main

import (
	"fmt"
	"log"
	"time"

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
	for _, b := range boards {
		tr, err := b.FetchThreadlist()
		if err != nil {
			log.Printf("error on fetching threadlist of %v: %v", b.URL, err)
			continue
		}
		trInUTF8 := transform.NewReader(tr, japanese.ShiftJIS.NewDecoder())
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
