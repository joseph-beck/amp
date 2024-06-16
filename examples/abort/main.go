package main

import (
	"fmt"
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/middleware/cors"
	"github.com/joseph-beck/amp/pkg/status"
)

func abort() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("aborting")

		ctx.AbortWithStatus(status.InternalServerError)

		return nil // ctx.Abort() does not work with ctx.Next().
	}
}

func getHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("handler")

		ctx.Status(status.OK)
		return nil
	}
}

func main() {
	a := amp.New()

	a.Use(cors.New())

	a.Get("/", getHandler())

	a.Get("/abort", getHandler(), abort())

	log.Fatalln(a.ListenAndServe())
}
