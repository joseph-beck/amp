package amp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	{
		g := Group("/test")

		assert.Equal(t, "/test", g.prefix)
		assert.Equal(t, 0, len(g.handlers))
		assert.Equal(t, 0, len(g.middleware))
	}

	{
		g := Group("/test", func(ctx *Ctx) error { return nil })

		assert.Equal(t, "/test", g.prefix)
		assert.Equal(t, 0, len(g.handlers))
		assert.Equal(t, 1, len(g.middleware))
	}
}

func TestGroupUse(t *testing.T) {
	g := Group("/test")

	g.Use(func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.middleware))
}

func TestGroupPrefix(t *testing.T) {
	g := Group("/test")

	assert.Equal(t, "/test", g.Prefix())
}

func TestGroupHandler(t *testing.T) {
	g := Group("/test")

	g.Handler("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "HANDLER", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
}

func TestGroupGet(t *testing.T) {
	g := Group("/test")

	g.Get("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "GET", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
}

func TestGroupPost(t *testing.T) {
	g := Group("/test")

	g.Post("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "POST", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
}

func TestGroupPut(t *testing.T) {
	g := Group("/test")

	g.Put("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "PUT", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
}

func TestGroupPatch(t *testing.T) {
	g := Group("/test")

	g.Patch("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "PATCH", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}

func TestGroupDelete(t *testing.T) {
	g := Group("/test")

	g.Delete("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "DELETE", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}

func TestGroupHead(t *testing.T) {
	g := Group("/test")

	g.Head("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "HEAD", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}

func TestGroupOptions(t *testing.T) {
	g := Group("/test")

	g.Options("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "OPTIONS", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}

func TestGroupConnect(t *testing.T) {
	g := Group("/test")

	g.Connect("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "CONNECT", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}

func TestGroupTrace(t *testing.T) {
	g := Group("/test")

	g.Trace("/hello", func(ctx *Ctx) error { return nil })

	assert.Equal(t, 1, len(g.handlers))
	assert.Equal(t, "TRACE", g.handlers[0].method)
	assert.Equal(t, "/hello", g.handlers[0].path)
	assert.Equal(t, "/hello", len(g.handlers[0].middleware))
}
