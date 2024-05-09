package amp

import (
	"net/http/httptest"
	"testing"

	"github.com/joseph-beck/amp/pkg/status"
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

func TestCtxPath(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		path := ctx.Path()
		assert.Equal(t, "/test", path)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxMethod(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		path := ctx.Method()
		assert.Equal(t, "GET", path)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
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

func TestCtxQuery(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		val, err := ctx.Query("query")
		assert.NoError(t, err)
		assert.Equal(t, "value", val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one?query=value", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		val, err := ctx.Query("query", "value")
		assert.NoError(t, err)
		assert.Equal(t, "value", val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/three", func(ctx *Ctx) error {
		val, err := ctx.Query("query")
		assert.Error(t, err)
		assert.Equal(t, "", val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/three", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxQueryInt(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		val, err := ctx.QueryInt("query")
		assert.NoError(t, err)
		assert.Equal(t, 1, val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one?query=1", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		val, err := ctx.QueryInt("query", 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/three", func(ctx *Ctx) error {
		val, err := ctx.QueryInt("query")
		assert.Error(t, err)
		assert.Equal(t, 0, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/three", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxQueryFloat(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		val, err := ctx.QueryFloat("query")
		assert.NoError(t, err)
		assert.Equal(t, float64(1), val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one?query=1.0", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		val, err := ctx.QueryFloat("query", float64(1))
		assert.NoError(t, err)
		assert.Equal(t, float64(1), val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/three", func(ctx *Ctx) error {
		val, err := ctx.QueryFloat("query")
		assert.Error(t, err)
		assert.Equal(t, float64(0), val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/three", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxQueryBool(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		val, err := ctx.QueryBool("query")
		assert.NoError(t, err)
		assert.Equal(t, true, val)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one?query=true", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		val, err := ctx.QueryBool("query", true)
		assert.NoError(t, err)
		assert.Equal(t, true, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/three", func(ctx *Ctx) error {
		val, err := ctx.QueryBool("query")
		assert.Error(t, err)
		assert.Equal(t, false, val)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/three", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxStatus(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		ctx.Status(status.OK)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxWrite(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		_, err := ctx.Write("write")
		assert.NoError(t, err)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxWriteBytes(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		_, err := ctx.WriteBytes([]byte("write"))
		assert.NoError(t, err)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRender(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.Render(status.OK, "write")
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRenderBytes(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderBytes(status.OK, []byte("write"))
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}
