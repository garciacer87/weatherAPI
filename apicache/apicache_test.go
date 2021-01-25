package apicache

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestCacheValue(t *testing.T) {
	c := New(1)
	c.SetValue("test_1", []byte(`{"message":"test"}`))

	v := c.GetValue("test_1")
	if v == nil {
		t.Errorf("Got: nil. Expected: %s", v)
	}
}

func TestCacheExpiration(t *testing.T) {
	mockCache := &apiCache{cache.New(500*time.Millisecond, 1*time.Second)}

	mockCache.SetValue("test_1", []byte(`{"message":"test"}`))

	time.Sleep(1 * time.Second)

	v := mockCache.GetValue("test_1")
	if v != nil {
		t.Errorf("Got: %s. Expected: nil", v)
	}
}
