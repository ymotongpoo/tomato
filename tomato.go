package main

import (
	"log"

	"github.com/ymotongpoo/tomato/monitor"
)

func main() {
	log.Println("**** process start")

	m, err := monitor.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	m.Start()
}
