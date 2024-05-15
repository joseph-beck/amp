package main

import (
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

type Mock struct {
	Value string `json:"value" toml:"value"`
}

func getJSON() amp.Handler {
	return func(ctx *amp.Ctx) error {
		val := Mock{Value: "json"}
		return ctx.RenderJSON(status.OK, val)
	}
}

func getTOML() amp.Handler {
	return func(ctx *amp.Ctx) error {
		val := Mock{Value: "json"}
		return ctx.RenderTOML(status.OK, val)
	}
}

func main() {
	a := amp.New()

	a.Get("/json", getJSON())
	a.Get("/toml", getTOML())

	log.Fatalln(a.ListenAndServe())
}
