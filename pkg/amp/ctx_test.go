package amp

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtxSet(t *testing.T) {
	ctx := newCtx(nil, nil)

	ctx.Set("key", "value")
	assert.Equal(t, ctx.values["key"], "value")
}

func TestCtxGet(t *testing.T) {
	ctx := newCtx(nil, nil)

	ctx.Set("key", "value")
	val, err := ctx.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	_, err = ctx.Get("no-key")
	assert.Error(t, err)
}

func TestCtxParam(t *testing.T) {
	amp := New()

	amp.Get("/test/one/{key}", func(ctx *Ctx) error {
		val, err := ctx.Param("key")
		assert.NoError(t, err)
		assert.Equal(t, "value", val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one/value", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxParamInt(t *testing.T) {
	amp := New()

	amp.Get("/test/one/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamInt("key")
		assert.NoError(t, err)
		assert.Equal(t, 1, val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one/1", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamInt("key")
		assert.Error(t, err)
		assert.Equal(t, 0, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two/value", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxParamFloat(t *testing.T) {
	amp := New()

	amp.Get("/test/one/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamFloat("key")
		assert.NoError(t, err)
		assert.Equal(t, float64(1), val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one/1.0", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamFloat("key")
		assert.Error(t, err)
		assert.Equal(t, float64(0), val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two/value", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxParamBool(t *testing.T) {
	amp := New()

	amp.Get("/test/one/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamBool("key")
		assert.NoError(t, err)
		assert.Equal(t, true, val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one/true", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two/{key}", func(ctx *Ctx) error {
		val, err := ctx.ParamBool("key")
		assert.Error(t, err)
		assert.Equal(t, false, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two/value", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}