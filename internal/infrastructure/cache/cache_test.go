package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newDisabledCache() *CacheService {
	return &CacheService{enabled: false}
}

func TestCacheDisabled_GetReturnsFalse(t *testing.T) {
	c := newDisabledCache()
	var dest string
	found, err := c.Get(context.Background(), "some-key", &dest)
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestCacheDisabled_SetIsNoop(t *testing.T) {
	c := newDisabledCache()
	err := c.Set(context.Background(), "some-key", "some-value")
	assert.NoError(t, err)
}

func TestCacheDisabled_DeleteIsNoop(t *testing.T) {
	c := newDisabledCache()
	err := c.Delete(context.Background(), "some-key")
	assert.NoError(t, err)
}

func TestCacheDisabled_InvalidatePatternIsNoop(t *testing.T) {
	c := newDisabledCache()
	err := c.InvalidatePattern(context.Background(), "users:*")
	assert.NoError(t, err)
}

func TestBuildKey_SinglePart(t *testing.T) {
	key := BuildKey("foo")
	assert.Equal(t, "foo", key)
}

func TestBuildKey_MultipleParts(t *testing.T) {
	key := BuildKey("users", 42)
	assert.Equal(t, "users:42", key)
}

func TestBuildKey_EmptyParts(t *testing.T) {
	key := BuildKey()
	assert.Equal(t, "", key)
}
