// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Auth is a middleware used for authorizing requests.
package auth

import (
	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

// Configure your authorization middleware.
type Config struct {
	// AuthFunc dictates whether a request will be allowed to continue or not.
	// If AuthFunc is set to nil, then there will be no authorization.
	// When using Default(), AuthFunc is nil.
	AuthFunc func(ctx *amp.Ctx) bool

	// NoAccessFunc determines what happens if no access is granted by AuthFunc, or false is returned.
	// When access is granted it simply moves to the next handler and this is not used.
	// If this is nil, then NoAccessCode is given and the Ctx aborted.
	// When using Default(), NoAccessFunc is nil.
	NoAccessFunc amp.Handler

	// What code would like requests to be given if they are not authorized?
	// When using Default(), NoAccessCode is status.Unauthorized or 401.
	NoAccessCode int
}

// You can pass your own custom authorization function through the arguments of this function.
// When creating a new auth middleware, an empty handler will be given if no auth func is given.
func Default(args ...func(ctx *amp.Ctx) bool) Config {
	var authFunc func(ctx *amp.Ctx) bool

	if len(args) > 0 {
		authFunc = args[0]
	}

	return Config{
		AuthFunc:     authFunc,
		NoAccessFunc: nil,
		NoAccessCode: status.Unauthorized,
	}
}
