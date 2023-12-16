package counter

import (
	"testing"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/sync/errgroup"
)

func TestGet(t *testing.T) {
	windowLen := 1 * time.Second
	window := time.Now().Truncate(windowLen)
	key := "test"
	c := New(windowLen)
	t.Run("Zero value", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			want := 0
			got, err := c.Get(key, window)
			if err != nil {
				t.Error(err)
			}
			if got != want {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		}
	})

	t.Run("Get value", func(t *testing.T) {
		want := 1
		v := uint64(want)
		c.cache.Set(generateKey(key, window), &v, ttlcache.DefaultTTL)
		got, err := c.Get(key, window)
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("Get() = %v, want %v", got, want)
		}
	})
}

func TestIncrement(t *testing.T) {
	windowLen := 1 * time.Millisecond
	window := time.Now().Truncate(windowLen)
	key := "test"
	t.Run("Increment simply", func(t *testing.T) {
		c := New(windowLen)
		for i := 0; i < 5; i++ {
			want := i + 1
			if err := c.Increment(key, window); err != nil {
				t.Error(err)
			}
			got, err := c.Get(key, window)
			if err != nil {
				t.Error(err)
			}
			if got != want {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		}
	})

	t.Run("Increment in parallel", func(t *testing.T) {
		c := New(windowLen)
		want := 1000
		eg := new(errgroup.Group)
		for i := 0; i < want; i++ {
			eg.Go(func() error {
				return c.Increment(key, window)
			})
		}
		if err := eg.Wait(); err != nil {
			t.Error(err)
		}
		got, err := c.Get(key, window)
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("Get() = %v, want %v", got, want)
		}
	})
}
