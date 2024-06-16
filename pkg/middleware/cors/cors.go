package cors

import (
	"github.com/joseph-beck/amp/pkg/amp"
)

func New(args ...Config) amp.Handler {
	_ = Default()

	if len(args) > 0 {
		_ = args[0]
	}

	return func(ctx *amp.Ctx) error {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		return ctx.Next()
	}
}
