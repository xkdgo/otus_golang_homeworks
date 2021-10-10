package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		for i := 0; i < 5; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		} // ["4":4 , "3":3, "2":2, "1":1, "0":0]
		c.Set(Key("5"), 5) // ["5":5, "4":4 , "3":3, "2":2, "1":1]
		val, ok := c.Get("0")
		require.False(t, ok)
		require.Nil(t, val)

		wasInCache := c.Set(Key("5"), "five") // ["5":"five", "4":4 , "3":3, "2":2, "1":1]
		require.True(t, wasInCache)

		val, ok = c.Get("5")
		require.True(t, ok)
		require.Equal(t, "five", val)

		c.Set(Key("4"), "four")
		c.Set(Key("3"), "three")
		c.Set(Key("2"), "two")
		c.Set(Key("1"), "one") // ["1": "one", "2":"two", "3":three, "4":"four", "5":"five"]
		c.Set(Key("VeryNewItem"), "SomeValue")
		val, ok = c.Get("VeryNewItem")
		require.True(t, ok)
		require.Equal(t, "SomeValue", val) // ["VeryNewItem": "SomeValue", "1": "one", "2":"two", "3":three, "4":"four"]
		c.Set(Key("VeryNewItem"), "SomeValue")
		// "5":"five" doesnt appear anymore
		val, ok = c.Get("5")
		require.Nil(t, val)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
