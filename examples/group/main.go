package main

import (
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
)

func main() {
	g := amp.Group("/group")
	g.Get("/hello", func(ctx *amp.Ctx) error {
		return ctx.Render(200, "hello world!")
	})

	a := amp.New()

	a.Group(g)

	log.Fatalln(a.ListenAndServe())
}
