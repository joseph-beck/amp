package limiter

import (
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := amp.New()

	a.Get("/test/one", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		SkipFunc:    nil,
		NextFunc:    nil,
		Limit:       1,
		Duration:    1 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       true,
	}))

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.Locked, writer.Result().StatusCode)

	a.Get("/test/two", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		SkipFunc: nil,
		NextFunc: func(ctx *amp.Ctx) error {
			ctx.Status(status.BadRequest)
			return nil
		},
		Limit:       1,
		Duration:    1 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       true,
	}))

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.BadRequest, writer.Result().StatusCode)

	a.Get("/test/three", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		SkipFunc: func(ctx *amp.Ctx) bool {
			log.Println(ctx.Origin())
			return ctx.Origin() == "example.com"
		},
		NextFunc: func(ctx *amp.Ctx) error {
			return nil
		},
		Limit:       1,
		Duration:    1 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       true,
	}))

	request = httptest.NewRequest("GET", "/test/three", nil)
	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	a.Get("/test/four", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		SkipFunc:    nil,
		NextFunc:    nil,
		Limit:       0,
		Duration:    1 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       true,
	}))

	a.Get("/test/five", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		SkipFunc:    nil,
		NextFunc:    nil,
		Limit:       1,
		Duration:    0 * time.Minute,
		LimitCode:   status.Locked,
		SkipFails:   false,
		SkipSuccess: false,
		Debug:       true,
	}))
}
