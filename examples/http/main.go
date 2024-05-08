package main

import (
	"errors"
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
)

func makeErrorHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		return errors.New("error")
	}
}

func main() {
	a := amp.New()

	a.Get("/error", makeErrorHandler())

	log.Fatalln(a.ListenAndServe())
}
