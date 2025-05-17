// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Cors is a middleware used for cross-origin requests.
package cors

import (
	"net/http"
)

// Configure the CORS middleware.
// Used within the args of New().
// Can use Default() to build the default CORS config.
type Config struct {
	// A string slice that contains origins that a cross-domain request can be made from.
	// If "*" is used then all origins will be allowed.
	// You can use "*" as a wildcard within origins, but only one per origin.
	// When using the Default(), AllowedOrigins will be []string{"*"}.
	AllowedOrigins []string

	// A custom function that can validate an origin.
	// It takes the origin as an argument, which it can use to validate access.
	// If this is not nil, the values of AllowedOrigins are ignored.
	AllowOriginFunc func(origin string) bool

	// Similar to AllowOriginFunc, takes both the request and origin to determine access.
	// If this function is not nil the values of AllowedOrigins and AllowOriginFunc are ignored.
	AllowOriginRequestFunc func(request *http.Request, origin string) bool

	// A string slice that contains the methods a client can use with cross-domain requests.
	// When using Default(), AllowedMethods will be []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}.
	AllowedMethods []string

	// A string slice that contains the headers a client can use with cross-domain requests.
	// If the value "*" is used all headers will be allowed.
	// When using Default(), AllowedHeaders will be []string{"Content-Type"}.
	AllowedHeaders []string

	// A string slice of headers that are safe to expose to the API.
	// When using Default(), ExposedHeaders will be []string{}.
	ExposedHeaders []string

	// MaxAge is used to indicate how long, in seconds, a preflight request can be cached.
	// When the value is 0 there is no caching of the preflight request.
	// If you need to force a a MaxAge of 0, use a negative number such as -1.
	// When using Default(), MaxAge will be 0.
	MaxAge int

	// AllowCredentials is used to indicate whether the request can contain things such as cookies,
	// HTTP authentication, client side SSL, etc.
	// When using Default(), AllowCredentials will be true.
	AllowCredentials bool

	// AllowPrivateNetwork is used to indicate whether to accept a cross-origin request over a private network.
	// When using Default(), AllowPrivateNetwork will be false.
	AllowPrivateNetwork bool

	// Debug CORS requests, uses slog like Amp to provide debugging detail.
	// When using Default(), Debug is false.
	Debug bool
}

func Default() Config {
	return Config{
		AllowedOrigins:         []string{"*"},
		AllowOriginFunc:        nil,
		AllowOriginRequestFunc: nil,
		AllowedMethods:         []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:         []string{"Content-Type"},
		ExposedHeaders:         []string{},
		MaxAge:                 0,
		AllowCredentials:       true,
		AllowPrivateNetwork:    false,
		Debug:                  false,
	}
}
