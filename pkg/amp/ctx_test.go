package amp

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joseph-beck/amp/pkg/binding"
	"github.com/joseph-beck/amp/pkg/status"
	"github.com/stretchr/testify/assert"
)

type Mock struct {
	Key string `json:"key" xml:"key" toml:"key" yaml:"key"`
}

func TestWriteContentType(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		err := ctx.RenderString(status.OK, "value")
		assert.NoError(t, err)
		assert.Equal(t, plainContentType, ctx.writer.Header()["Content-Type"])

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxWriter(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		w := ctx.Writer()
		assert.Equal(t, ctx.writer, w)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxSetWriter(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		w := httptest.NewRecorder()
		ctx.SetWriter(w)
		assert.Equal(t, ctx.writer, w)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRequest(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		r := ctx.Request()
		assert.Equal(t, ctx.request, r)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxSetRequest(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		r := httptest.NewRequest("GET", "/test", nil)
		ctx.SetRequest(r)
		assert.Equal(t, ctx.request, r)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

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
		method := ctx.Method()
		assert.Equal(t, "GET", method)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxHeader(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		header := "none"
		ctx.Header("X-Auth", header)
		assert.Equal(t, header, ctx.writer.Header().Get("X-Auth"))

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxOrigin(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		origin := ctx.Origin()
		assert.Equal(t, "example.com", origin)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxNext(t *testing.T) {
	amp := New()

	amp.Get(
		"/test",
		func(ctx *Ctx) error {
			idx := ctx.index
			assert.Equal(t, 1, idx)

			return nil
		},
		func(ctx *Ctx) error {
			return ctx.Next()
		},
	)

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

	amp.Get("/test", func(ctx *Ctx) error {
		ctx.Status(status.OK)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxGetStatus(t *testing.T) {
	amp := New()

	amp.Get("/test", func(ctx *Ctx) error {
		ctx.Status(status.OK)
		assert.Equal(t, status.OK, ctx.GetStatus())

		ctx.Status(status.BadRequest)
		assert.Equal(t, status.BadRequest, ctx.GetStatus())

		return nil
	})

	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxAborted(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		assert.Equal(t, false, ctx.Aborted())

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		ctx.aborted = true
		assert.Equal(t, true, ctx.Aborted())

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxAbort(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		ctx.Abort()
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxAbortWithStatus(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		ctx.AbortWithStatus(status.BadRequest)
		assert.Equal(t, true, ctx.aborted)
		assert.Equal(t, status.BadRequest, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxAbortWithError(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		ctx.AbortWithError(status.BadRequest, nil)
		assert.Equal(t, true, ctx.aborted)
		assert.Equal(t, status.BadRequest, ctx.status)

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

func TestCtxRenderString(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderString(status.OK, "write")
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		err := ctx.RenderString(status.OK, "write %d", 123)
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRenderJSON(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderJSON(status.OK, Mock{Key: "value"})
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRenderTOML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderTOML(status.OK, Mock{Key: "value"})
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRenderYAML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderYAML(status.OK, Mock{Key: "value"})
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxRenderXML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		err := ctx.RenderXML(status.OK, Mock{Key: "value"})
		assert.NoError(t, err)
		assert.Equal(t, status.OK, ctx.status)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", nil)
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxShouldBindWith(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindWith(&obj, binding.JSON)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`{"key": "value"}`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindWith(&obj, binding.JSON)
		assert.Error(t, err)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxMustBindWith(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.MustBindWith(&obj, binding.JSON)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`{"key": "value"}`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.MustBindWith(&obj, binding.JSON)
		assert.Error(t, err)
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxShouldBindJSON(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindJSON(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`{"key": "value"}`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindJSON(&obj)
		assert.Error(t, err)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxBindJSON(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindJSON(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`{"key": "value"}`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindJSON(&obj)
		assert.Error(t, err)
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", nil)
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxShouldBindTOML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindTOML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`key = "value"`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindTOML(&obj)
		assert.Error(t, err)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key = key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxBindTOML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindTOML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`key = "value"`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindTOML(&obj)
		assert.Error(t, err)
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key = key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxShouldBindYAML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindYAML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`key: value`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindYAML(&obj)
		assert.Error(t, err)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxBindYAML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindYAML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`key: value`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindYAML(&obj)
		assert.Error(t, err)
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxShouldBindXML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindXML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
			<key>value</key>
		</root>
	`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.ShouldBindXML(&obj)
		assert.Error(t, err)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}

func TestCtxBindXML(t *testing.T) {
	amp := New()

	amp.Get("/test/one", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindXML(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "value", obj.Key)

		return nil
	})

	request := httptest.NewRequest("GET", "/test/one", strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
   			<key>value</key>
		</root>
	`))
	writer := httptest.NewRecorder()
	amp.ServeHTTP(writer, request)

	amp.Get("/test/two", func(ctx *Ctx) error {
		var obj Mock
		err := ctx.BindXML(&obj)
		assert.Error(t, err)
		assert.Equal(t, true, ctx.aborted)

		return nil
	})

	request = httptest.NewRequest("GET", "/test/two", strings.NewReader(`key`))
	writer = httptest.NewRecorder()
	amp.ServeHTTP(writer, request)
}
