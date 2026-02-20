//go:build integration

package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testCache *CacheService

func TestMain(m *testing.M) {
	ctx := context.Background()

	redisContainer, err := tcRedis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start redis container: %v", err)
	}

	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate redis container: %v", err)
		}
	}()

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get redis connection string: %v", err)
	}

	opt, err := redis.ParseURL(connStr)
	if err != nil {
		log.Fatalf("failed to parse redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}

	testCache = &CacheService{
		client:  client,
		enabled: true,
		prefix:  "test:",
		ttl:     5 * time.Minute,
	}

	fmt.Println("Test Redis ready")
	os.Exit(m.Run())
}

func flushRedis(t *testing.T) {
	t.Helper()
	err := testCache.client.FlushAll(context.Background()).Err()
	require.NoError(t, err)
}

func TestCache_SetAndGet_Integration(t *testing.T) {
	flushRedis(t)
	t.Cleanup(func() { flushRedis(t) })

	ctx := context.Background()

	// Set a value
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	input := testData{Name: "hello", Value: 42}
	err := testCache.Set(ctx, "mykey", input)
	require.NoError(t, err)

	// Get the value back
	var output testData
	found, err := testCache.Get(ctx, "mykey", &output)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "hello", output.Name)
	assert.Equal(t, 42, output.Value)

	// Get with string type
	err = testCache.Set(ctx, "strkey", "simple-string")
	require.NoError(t, err)

	var strOutput string
	found, err = testCache.Get(ctx, "strkey", &strOutput)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "simple-string", strOutput)
}

func TestCache_Delete_Integration(t *testing.T) {
	flushRedis(t)
	t.Cleanup(func() { flushRedis(t) })

	ctx := context.Background()

	// Set a value
	err := testCache.Set(ctx, "deletekey", "value-to-delete")
	require.NoError(t, err)

	// Verify it exists
	var output string
	found, err := testCache.Get(ctx, "deletekey", &output)
	require.NoError(t, err)
	assert.True(t, found)

	// Delete it
	err = testCache.Delete(ctx, "deletekey")
	require.NoError(t, err)

	// Verify it no longer exists
	found, err = testCache.Get(ctx, "deletekey", &output)
	require.NoError(t, err)
	assert.False(t, found)
}

func TestCache_InvalidatePattern_Integration(t *testing.T) {
	flushRedis(t)
	t.Cleanup(func() { flushRedis(t) })

	ctx := context.Background()

	// Set multiple keys with a common prefix pattern
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("users:%d", i)
		err := testCache.Set(ctx, key, fmt.Sprintf("user-%d", i))
		require.NoError(t, err)
	}

	// Set a key with a different prefix
	err := testCache.Set(ctx, "brands:1", "Toyota")
	require.NoError(t, err)

	// Invalidate all user keys
	err = testCache.InvalidatePattern(ctx, "users:*")
	require.NoError(t, err)

	// Verify user keys are gone
	for i := 0; i < 5; i++ {
		var output string
		key := fmt.Sprintf("users:%d", i)
		found, err := testCache.Get(ctx, key, &output)
		require.NoError(t, err)
		assert.False(t, found, "key %s should have been invalidated", key)
	}

	// Verify brand key still exists
	var brandOutput string
	found, err := testCache.Get(ctx, "brands:1", &brandOutput)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "Toyota", brandOutput)
}

func TestCache_GetMiss_Integration(t *testing.T) {
	flushRedis(t)
	t.Cleanup(func() { flushRedis(t) })

	ctx := context.Background()

	// Get a key that does not exist
	var output string
	found, err := testCache.Get(ctx, "nonexistent-key", &output)
	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, "", output)
}
