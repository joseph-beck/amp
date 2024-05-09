package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

func getRoot() amp.Handler {
	return func(ctx *amp.Ctx) error {
		return ctx.Render(status.OK, "hello from root")
	}
}

func getError() amp.Handler {
	return func(ctx *amp.Ctx) error {
		return errors.New("error")
	}
}

func getParam() amp.Handler {
	return func(ctx *amp.Ctx) error {
		name, err := ctx.Param("name")
		if err != nil {
			return err
		}

		return ctx.Render(status.OK, fmt.Sprintf("hello %s", name))
	}
}

func main() {
	a := amp.New()

	a.Get("/", getRoot())

	a.Get("/name/{name}", getParam())

	a.Get("/error", getError())

	log.Fatalln(a.ListenAndServe())
}
