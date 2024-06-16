package cors

import (
	"log/slog"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

func New(args ...Config) amp.Handler {
	return func(ctx *amp.Ctx) error {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if ctx.Method() == "OPTIONS" {
			return ctx.Render(status.OK, "")
		}

		slog.Info("CORS")

		return ctx.Next()
	}
}
