// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Auth is a middleware used for authorizing requests.
package auth

import (
	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

// unexported auth struct, used to store details about our authorization.
type auth struct {
	// unexported authFunc function.
	// if this is nil, an empty handler is given instead, no auth.
	authFunc func(ctx *amp.Ctx) bool

	// unexported noAccessFunc function.
	noAccessFunc amp.Handler

	// unexported noAccessCode.
	noAccessCode int
}

// Create a new Authorization Middelware.
// If no config is given the Default() config is used,
// this will result in being given an empty handler.
func New(args ...Config) amp.Handler {
	cfg := Default()

	if len(args) > 0 {
		cfg = args[0]
	}

	auth := auth{}

	// if we have no auth func, return an empty handler.
	if cfg.AuthFunc == nil {
		return func(ctx *amp.Ctx) error {
			return nil
		}
	}

	auth.authFunc = cfg.AuthFunc

	// if we have an invalid access code, make it valid, otherwise set it.
	if cfg.NoAccessCode <= 0 {
		auth.noAccessCode = status.Unauthorized
	} else {
		auth.noAccessCode = cfg.NoAccessCode
	}

	// if we have a no access func, lets set it.
	if cfg.NoAccessFunc != nil {
		auth.noAccessFunc = cfg.NoAccessFunc
	}

	return func(ctx *amp.Ctx) error {
		// if we have no access, abort the ctx.
		if !auth.authFunc(ctx) {
			ctx.Abort()
			// if we have a no access func, lets use it.
			if auth.noAccessFunc != nil {
				return auth.noAccessFunc(ctx)
			}

			ctx.Status(auth.noAccessCode)
			return nil
		}

		// if we are authorized we can continue down the handler chain.
		return nil
	}
}
