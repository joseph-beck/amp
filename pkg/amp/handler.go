package amp

import "net/http"

type Handler func(c *Ctx) error

func (h Handler) Unwrap() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
	}
}

func Wrap(fn http.HandlerFunc) Handler {
	return func(ctx *Ctx) error {
		fn(ctx.writer, ctx.request)
		return nil
	}
}
