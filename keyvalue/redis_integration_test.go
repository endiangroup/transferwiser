// +build integration

package keyvalue

import (
	"testing"

	"github.com/endiangroup/transferwiser/core"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
)

func redisClient(t *testing.T) *redis.Client {
	redisAddr := core.GetConfig().RedisAddr

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	require.NoError(t, client.Ping().Err())
	require.NoError(t, client.FlushDb().Err())
	return client
}

func Test_RedisKeyValueStorage(t *testing.T) {
	kv := NewRedisKeyValue(redisClient(t))
	sharedKeyValueTests(t, kv)
}

func Test_RedisKeyValueNamespacedStorage(t *testing.T) {
	kv := NewRedisNamespacedKeyValue(redisClient(t), "test-ns")
	sharedKeyValueTests(t, kv)
}
