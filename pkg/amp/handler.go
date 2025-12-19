// GitHub Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Amp is a web framework made using the Go 1.22 Mux.
// Please ensure you are using Go 1.22, minimum, when using Amp.
package amp

import (
	"log/slog"
	"net/http"
)

// Amp Handler.
// Uses the *amp.Ctx.
// Returns an error, will slog the error if one occurs during execution.
// Centralised error handling within the Mux.
type Handler func(ctx *Ctx) error

// Unwrap an amp.Handler into a net/http HandlerFunc.
func (h Handler) Unwrap(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx(w, r)
		err := fn(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Take a net/http HandlerFunc and wrap it into an amp.Handler.
// Allows you to write http.HandlerFunc but still use them with amp.
func Wrap(fn http.HandlerFunc) Handler {
	return func(ctx *Ctx) error {
		fn(ctx.writer, ctx.request)
		return nil
	}
}
