package main

import (
	"fmt"
	"log"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/middleware/cors"
	"github.com/joseph-beck/amp/pkg/status"
)

func middleware() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("middleware")

		return nil // can also use ctx.Next()
	}
}

func getHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("handler")

		ctx.Status(status.OK)
		return nil
	}
}

func postHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("Test")

		var i interface{}
		err := ctx.BindJSON(&i)
		if err != nil {
			return err
		}

		fmt.Println(i)

		return ctx.RenderJSON(status.OK, i)
	}
}

func main() {
	a := amp.New()

	a.Use(cors.New(cors.Default()))

	a.Get("/", getHandler(), middleware())

	a.Post("/", postHandler(), middleware())

	log.Fatalln(a.ListenAndServe())
}
