// GitHub Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Limiter is a middleware used for rate limiting requests.
package limiter

import (
	"time"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

// Configure the Amp Limiter middleware.
type Config struct {
	// Using the amp.Ctx, the middleware can be skipped if this function returns true.
	// When using the Default(), SkipFunc will be nil.
	SkipFunc func(ctx *amp.Ctx) bool

	// Handler that if the rate limit is reached, will go to this Handler instead.
	// The Handler will be ran before being aborted by the middleware.
	// When using the Default(), NextFunc will be nil.
	NextFunc amp.Handler

	// Allows you to generate custom keys based on the value returned from this function.
	// The default key generator will use the host address of the request.
	// When using the Default(), KeyGeneratorFunc will be nil.
	KeyGeneratorFunc func(ctx *amp.Ctx) string

	// Maximum number of tries that can be done before being rate limited.
	// When the user is rate limited they will be sent a status.Locked or a 423.
	// When using the Default(), Limit will be 10.
	Limit int

	// The time in which a rate limit will last.
	// When using the Default(), Duration will be 1 * time.Minute.
	Duration time.Duration

	// Define the error code of the rate limit.
	// When using Default(), LimitCode will be status.Locked or 423.
	LimitCode int

	// If the request has a code that is >= 400, then it will not be counted towards the rate limiter.
	// When using the Default(), SkipFails is false.
	SkipFails bool

	// If the request has a code that is < 400 then it will not be counted towards the rate limiter.
	// When using the Default(), SkipSuccess is false.
	SkipSuccess bool

	// Allows us to debug our rate limitting, will print information such as the time since request, key, etc.
	// When using Default(), Debug is false.
	Debug bool
}

// Returns the default configuration for the rate limiter.
func Default() Config {
	return Config{
		SkipFunc:    nil,
		NextFunc:    nil,
		Limit:       10,
		Duration:    1 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       false,
	}
}
