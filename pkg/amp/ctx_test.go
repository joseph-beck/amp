package amp

import (
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
