package cors

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/joseph-beck/amp/pkg/amp"
)

func New(args ...Config) amp.Handler {
	var cfg Config
	if len(args) > 0 {
		cfg = args[0]
	} else {
		cfg = Config{
			AllowedOrigins:   []string{},
			AllowedMethods:   []string{},
			AllowedHeaders:   []string{},
			AllowCredentials: false,
			ExposeHeaders:    []string{},
			Debug:            false,
		}
	}

	return func(ctx *amp.Ctx) error {
		wrtr := ctx.Writer()
		rqst := ctx.Request()

		orgn := rqst.Header.Get("Origin")
		if orgn != "" {
			if contains(cfg.AllowedOrigins, "*") || contains(cfg.AllowedOrigins, orgn) {
				wrtr.Header().Add("Access-Control-Allow-Origin", orgn)

				if cfg.AllowCredentials {
					wrtr.Header().Add("Access-Control-Allow-Credentials", "true")
				}

				if len(cfg.ExposeHeaders) > 0 {
					wrtr.Header().Add("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
				}
			}
		}

		if rqst.Method == http.MethodOptions {
			if len(cfg.AllowedMethods) > 0 {
				wrtr.Header().Add("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			}

			if len(cfg.AllowedHeaders) > 0 {
				wrtr.Header().Add("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			}

			wrtr.WriteHeader(http.StatusOK)
			return nil
		}

		if cfg.Debug {
			slog.Info("CORS " + rqst.Method + " " + orgn)
		}

		return ctx.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
