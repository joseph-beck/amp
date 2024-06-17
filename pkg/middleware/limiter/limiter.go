package limiter

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
)

// unexported limiter struct, used to store our rate limiter settings privately.
type limiter struct {
	// unexported skip function.
	skip func(ctx *amp.Ctx) bool

	// unexported next handler.
	next amp.Handler

	// unexported keyGenerator.
	// if we are not given a key generator, a default one using the origin is created.
	keyGenerator func(ctx *amp.Ctx) string

	// unexported limit.
	// if the limit is 0, an empty handler is returned instead.
	limit int

	// unexported duration.
	// if the duration is 0, an empty handler is returned instead.
	duration time.Duration

	// unexported limitCode.
	// this will default to 423 if it is set to 0.
	limitCode int

	// unexported skipFails.
	skipFails bool

	// unexported skipSuccess.
	skipSuccess bool

	// unexported debug.
	debug bool

	// unexported store.
	// stores all requests and their corresponding hits and time since requests.
	store store
}

// Create a new rate limiter middleware.
// If this is given a config it will use that, otherwise Default() config is used.
func New(args ...Config) amp.Handler {
	cfg := Default()

	if len(args) > 0 {
		cfg = args[0]
	}

	limiter := limiter{
		skipFails:   cfg.SkipFails,
		skipSuccess: cfg.SkipSuccess,
		debug:       cfg.Debug,
		store:       newStore(),
	}

	// lets set skip if we have a skip func.
	if cfg.Skip != nil {
		limiter.skip = cfg.Skip
	}

	// lets set next if we have a next handler.
	if cfg.Next != nil {
		limiter.next = cfg.Next
	}

	// default key generator uses the origin as a key.
	if cfg.KeyGenerator != nil {
		limiter.keyGenerator = cfg.KeyGenerator
	} else {
		limiter.keyGenerator = func(ctx *amp.Ctx) string {
			return ctx.Origin()
		}
	}

	// return an empty handler if there is no rate limit on requests.
	if cfg.Limit <= 0 {
		return func(ctx *amp.Ctx) error {
			return nil
		}
	}
	limiter.limit = cfg.Limit

	// return an empty handler if there is no rate limit on duration.
	if cfg.Duration <= 0*time.Minute {
		return func(ctx *amp.Ctx) error {
			return nil
		}
	}
	limiter.duration = cfg.Duration

	// check to see if the limit code is valid otherwise set it do default.
	if cfg.LimitCode > 0 {
		limiter.limitCode = cfg.LimitCode
	} else {
		limiter.limitCode = status.Locked
	}

	return func(ctx *amp.Ctx) error {
		// if we have a skip function, lets check if the ctx applies the skip.
		if limiter.skip != nil {
			if limiter.skip(ctx) {
				return nil
			}
		}

		// get our key, if we have a key generator use that, otherwise we get it from the origin.
		key := func() string {
			if limiter.keyGenerator != nil {
				return limiter.keyGenerator(ctx)
			}

			return ctx.Origin()
		}()

		// is our current request rate limited?
		limited := limiter.store.limited(key, limiter.duration, limiter.limit)

		// if we are rate limited abort.
		if limited {
			ctx.Abort()

			// if we have a custom next function, use it.
			if limiter.next != nil {
				err := limiter.next(ctx)
				if err != nil {
					return err
				}

				return nil
			}

			// render rate limited msg, with our limit code and exit.
			return ctx.Render(limiter.limitCode, "Rate Limit Reached")
		}

		// iterate through our stack
		err := ctx.Next()
		if err != nil {
			return err
		}

		// iterate limiter if we are not skipping fails and the request has failed.
		if !limiter.skipFails && ctx.GetStatus() >= 400 {
			limiter.store.insert(key)
		}

		// iterate limiter if we are not skipping success and the request has succeeded.
		if !limiter.skipSuccess && ctx.GetStatus() < 400 {
			limiter.store.insert(key)
		}

		// give some info if we are using the debugger.
		if limiter.debug {
			item := limiter.store.get(key)
			slog.Info(fmt.Sprintf("LIMITER %s %v %d", key, item.timeSinceRequest, item.hits))
		}

		// if we are not rate limited lets just continue through the mux.
		return nil
	}
}
