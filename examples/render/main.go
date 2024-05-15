package main

import (
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

type Mock struct {
	Value string `json:"value"`
}

func getMock() amp.Handler {
	return func(ctx *amp.Ctx) error {
		val := Mock{Value: "hello"}
		return ctx.RenderJSON(status.OK, val)
	}
}

func main() {
	a := amp.New()

	a.Get("/mock", getMock())

	log.Fatalln(a.ListenAndServe())
}
