package limiter

import (
	"sync"
	"time"
)

// item used to calculate if a key has been hit too many times and when it was first created.
type item struct {
	// the amount of hits the current item has.
	// used to compare to the rate limit.
	hits int

	// time since request, or the time the last request was made.
	// this is compared to the current time to see if it violates the rate limit with the duration.
	timeSinceRequest time.Time
}

// creates a new empty item, the hits start at 1, and the timeSinceRequest is now.
func newItem() item {
	return item{
		hits:             1,
		timeSinceRequest: time.Now(),
	}
}

// increments the item.
// increments the number of hits.
func (i *item) increment() {
	i.hits++
}

// checks if the items time since request has expired.
func (i *item) expired(duration time.Duration) bool {
	return time.Since(i.timeSinceRequest) > duration
}

// use to store a map[string]any of key for the rate limiter.
type store struct {
	// map string item, for storing all items of the store.
	// the key, either the origin or created by the key generator of type string.
	// the item being created or iterated on depending on the request.
	items map[string]item

	// mutex for the store.
	// prevents any data races by locking during access.
	mu sync.Mutex
}

// create a new store with an empty map and mutex.
func newStore() store {
	return store{
		items: make(map[string]item, 0),
		mu:    sync.Mutex{},
	}
}

// get an item from the items, map string item.
// will return new item if it does not exist.
func (s *store) get(key string) item {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.items[key]
	if !ok {
		return newItem()
	}

	return val
}

// add an item to the store.
// if it does not exist a new item is made.
// otherwise the item is incremented.
func (s *store) insert(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.items[key]
	if !ok {
		s.items[key] = newItem()

		return
	}

	val.increment()
	s.items[key] = val
}

// checks to see if the given key exists in the store.
func (s *store) exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[key]
	return ok
}

// checks to see if the item of the given key is currently rate limited.
// uses duration and limit comparisons.
func (s *store) limited(key string, duration time.Duration, limit int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.items[key]
	if !ok {
		return false
	}

	exp := val.expired(duration)
	if exp {
		return false
	}

	if val.hits >= limit {
		return true
	}

	return false
}

// remove a given key from the store.
func (s *store) remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}
