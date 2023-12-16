package counter

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// Counter is a sliding window counter implemented with a TTL cache
type Counter struct {
	cache *ttlcache.Cache[string, *uint64]
	// capacity is the maximum number of items to store in the cache
	capacity uint64
	// disableAutoDeleteExpired disables the automatic deletion of expired items
	disableAutoDeleteExpired bool
}

type Option func(*Counter)

// WithCapacity sets the maximum number of items to store in the cache
func WithCapacity(capacity uint64) Option {
	return func(c *Counter) {
		c.capacity = capacity
	}
}

// DisableAutoDeleteExpired disables the automatic deletion of expired items
func DisableAutoDeleteExpired() Option {
	return func(c *Counter) {
		c.disableAutoDeleteExpired = true
	}
}

// NewCounter creates a new Counter
func New(ttl time.Duration, opts ...Option) *Counter {
	c := &Counter{}
	for _, opt := range opts {
		opt(c)
	}
	ttlOpts := []ttlcache.Option[string, *uint64]{
		ttlcache.WithTTL[string, *uint64](ttl),
	}
	if c.capacity > 0 {
		ttlOpts = append(ttlOpts, ttlcache.WithCapacity[string, *uint64](c.capacity))
	}
	cache := ttlcache.New[string, *uint64](ttlOpts...)
	c.cache = cache
	if !c.disableAutoDeleteExpired {
		go cache.Start()
	}
	return c
}

// Get returns the count for the given key and window
func (c *Counter) Get(key string, window time.Time) (int, error) { //nostyle:getters
	key = generateKey(key, window)
	i := c.cache.Get(key)
	if i == nil {
		return 0, nil
	}
	return int(*i.Value()), nil
}

// Increment increments the count for the given key and window
func (c *Counter) Increment(key string, currWindow time.Time) error {
	key = generateKey(key, currWindow)
	zero := uint64(0)
	i, _ := c.cache.GetOrSet(key, &zero)
	atomic.AddUint64(i.Value(), 1)
	return nil
}

func generateKey(key string, window time.Time) string {
	return fmt.Sprintf("%s-%d", key, window.Unix())
}

// DeleteExpired deletes expired items from the cache
func (c *Counter) DeleteExpired() {
	c.cache.DeleteExpired()
}
