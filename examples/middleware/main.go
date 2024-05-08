package main

import (
	"fmt"
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
)

func middleware() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("middleware")
		return nil
	}
}

func handler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("handler")
		return nil
	}
}

func main() {
	a := amp.New()

	a.Get("/", handler(), middleware())

	log.Fatalln(a.ListenAndServe())
}
