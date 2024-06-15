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

		return ctx.Next()
	}
}

func getHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		fmt.Println("handler")

		return nil
	}
}

func postHandler() amp.Handler {
	return func(ctx *amp.Ctx) error {
		var i interface{}
		err := ctx.BindJSON(&i)
		if err != nil {
			return err
		}

		return ctx.RenderJSON(status.OK, i)
	}
}

func main() {
	a := amp.New()

	a.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"*", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{},
		Debug:            true,
	}))

	a.Get("/", getHandler(), middleware())

	a.Post("/", postHandler())

	log.Fatalln(a.ListenAndServe())
}
