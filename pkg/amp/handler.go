package amp

import "net/http"

type Handler func(c *Ctx)

func Wrap(fn http.HandlerFunc) Handler {
	return func(ctx *Ctx) {
		fn(ctx.writer, ctx.request)
	}
}
