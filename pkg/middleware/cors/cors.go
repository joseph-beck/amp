package cors

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

type cors struct {
	// unexported allowedOrigins, converted from array to string.
	allowedOrigins string

	// unexported allowedOriginFimc, converts AllowOriginFunc or AllowOriginRequestFunc to this.
	allowOriginFunc func(request *http.Request, origin string) bool

	// unexported allowedMethods, converted from array to string.
	allowedMethods string

	// unexported allowedHeaders, converted from array to string.
	allowedHeaders string

	// unexported exposedHeaders, converted from array to string.
	exposedHeaders string

	// unexported maxAge, converted from int to string.
	maxAge string

	// unexported allowCredentials.
	allowCredentials bool

	// unexported allowPrivateNetwork.
	allowPrivateNetwork bool

	// unexorted debug.
	debug bool
}

func New(args ...Config) amp.Handler {
	cfg := Default()

	if len(args) > 0 {
		cfg = args[0]
	}

	cors := cors{
		allowCredentials:    cfg.AllowCredentials,
		allowPrivateNetwork: cfg.AllowPrivateNetwork,
		debug:               cfg.Debug,
	}

	if cfg.AllowedOrigins != nil {
		cors.allowedOrigins = strings.Join(cfg.AllowedOrigins, ", ")
	}

	if cfg.AllowOriginFunc != nil {
		cors.allowOriginFunc = func(request *http.Request, origin string) bool {
			return cfg.AllowOriginFunc(origin)
		}
	}

	if cfg.AllowOriginRequestFunc != nil {
		cors.allowOriginFunc = cfg.AllowOriginRequestFunc
	}

	if cfg.AllowedMethods != nil && len(cfg.AllowedMethods) > 0 {
		cors.allowedMethods = strings.Join(cfg.AllowedMethods, ", ")
	}

	if cfg.AllowedHeaders != nil && len(cfg.AllowedHeaders) > 0 {
		cors.allowedHeaders = strings.Join(cfg.AllowedHeaders, ", ")
	}

	if cfg.ExposedHeaders != nil && len(cfg.ExposedHeaders) > 0 {
		cors.exposedHeaders = strings.Join(cfg.ExposedHeaders, ", ")
	}

	if cfg.MaxAge > 0 {
		cors.maxAge = strconv.Itoa(cfg.MaxAge)
	} else if cfg.MaxAge <= 0 {
		cors.maxAge = "0"
	}

	return func(ctx *amp.Ctx) error {
		if cors.allowOriginFunc != nil {
			res := cors.allowOriginFunc(ctx.Request(), ctx.Origin())
			if res {
				return nil
			}

			ctx.AbortWithStatus(status.Forbidden)
			return nil
		}

		if cors.allowedOrigins != "" {
			ctx.Header("Access-Control-Allow-Origin", cors.allowedOrigins)
		}

		if cors.allowedMethods != "" {
			ctx.Header("Access-Control-Allow-Methods", cors.allowedMethods)
		}

		if cors.allowedHeaders != "" {
			ctx.Header("Access-Control-Allow-Headers", cors.allowedHeaders)
		}

		if cors.exposedHeaders != "" {
			ctx.Header("Access-Control-Expose-Headers", cors.exposedHeaders)
		}

		if cors.allowCredentials {
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		if cors.allowPrivateNetwork {
			ctx.Header("Access-Control-Allow-Private-Network", "true")
		}

		ctx.Header("Access-Control-Max-Age", cors.maxAge)

		if cors.debug {
			slog.Info(fmt.Sprintf("CORS %s %s %s", ctx.Origin(), ctx.Method(), ctx.Path()))
		}

		return nil
	}
}
