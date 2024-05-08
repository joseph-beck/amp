package amp

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newCtx() Ctx {
	return Ctx{
		values:   make(map[string]any),
		valuesMu: sync.Mutex{},
	}
}

func TestCtxSet(t *testing.T) {
	ctx := newCtx()

	ctx.Set("key", "value")
	assert.Equal(t, ctx.values["key"], "value")
}

func TestCtxGet(t *testing.T) {
	ctx := newCtx()

	ctx.Set("key", "value")
	val, err := ctx.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	_, err = ctx.Get("no-key")
	assert.Error(t, err)
}
