package main

import (
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
)

func main() {
	a := amp.New()

	log.Fatalln(a.ListenAndServe())
}
