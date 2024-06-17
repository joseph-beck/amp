package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	item := newItem()
	assert.Equal(t, 1, item.hits)
}

func TestItemIncrement(t *testing.T) {
	item := newItem()
	item.increment()
	assert.Equal(t, 2, item.hits)
}

func TestItemExpired(t *testing.T) {
	item := newItem()
	exp := item.expired(0 * time.Second)
	assert.True(t, exp)
	exp = item.expired(1 * time.Second)
	assert.False(t, exp)
}

func TestNewStore(t *testing.T) {
	store := newStore()
	assert.Equal(t, make(map[string]item), store.items)
}

func TestStoreGet(t *testing.T) {
	store := newStore()
	store.insert("1")
	assert.Equal(t, store.get("1").hits, 1)
}

func TestStoreInsert(t *testing.T) {
	store := newStore()
	store.insert("1")
	assert.True(t, store.exists("1"))
}

func TestStoreExists(t *testing.T) {
	store := newStore()
	assert.False(t, store.exists("1"))
	store.insert("1")
	assert.True(t, store.exists("1"))
}

func TestStoreLimited(t *testing.T) {
	store := newStore()
	store.insert("1")
	assert.True(t, store.limited("1", 1*time.Minute, 1))
	assert.False(t, store.limited("1", 1*time.Minute, 2))
	assert.False(t, store.limited("1", 0*time.Minute, 2))
	assert.False(t, store.limited("2", 1*time.Minute, 2))
}

func TestStoreRemove(t *testing.T) {
	store := newStore()
	store.insert("1")
	store.remove("1")
	assert.False(t, store.exists("1"))
}
