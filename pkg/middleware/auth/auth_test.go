package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/joseph-beck/amp/pkg/amp"
	"github.com/joseph-beck/amp/pkg/status"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := amp.New()

	authFunc := func(ctx *amp.Ctx) bool {
		val, err := ctx.Param("name")
		if err != nil {
			return false
		}

		if val == "hello" {
			return true
		}

		return false
	}

	a.Get("/test/one/{name}", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		AuthFunc: authFunc,
		NoAccessFunc: nil,
		NoAccessCode: status.Unauthorized,
	}))

	request := httptest.NewRequest("GET", "/test/one/hello", nil)
	writer := httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	request = httptest.NewRequest("GET", "/test/one/invalid", nil)
	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.Unauthorized, writer.Result().StatusCode)

	a.Get("/test/two/{name}", func(ctx *amp.Ctx) error {
		ctx.Status(status.OK)
		return nil
	}, New(Config{
		AuthFunc: authFunc,
		NoAccessFunc: func(ctx *amp.Ctx) error {
			ctx.Status(status.BadRequest)
			return nil
		},
		NoAccessCode: status.Unauthorized,
	}))

	request = httptest.NewRequest("GET", "/test/two/hello", nil)
	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.OK, writer.Result().StatusCode)

	request = httptest.NewRequest("GET", "/test/two/invalid", nil)
	writer = httptest.NewRecorder()
	a.ServeHTTP(writer, request)
	assert.Equal(t, status.BadRequest, writer.Result().StatusCode)
}
