package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_getCacheKey(t *testing.T) {
	repo := &AuthRepository{}
	key := repo.getCacheKey("test_token")
	assert.Equal(t, "token:test_token", key)
}

func TestAuthRepository_CreateToken_Success(t *testing.T) {
	t.Skip("Requires actual database connection")
}

func TestAuthRepository_FetchToken_Integration(t *testing.T) {
	t.Skip("Requires actual database connection")

	repo := &AuthRepository{
		CacheTTL: 5 * time.Minute,
	}

	token, err := repo.FetchToken("test_token")
	assert.NoError(t, err)
	assert.Nil(t, token)
}

func TestCacheKeyPrefix(t *testing.T) {
	assert.Equal(t, "token:", REDIS_TOKEN_PREFIX)
}

func TestNewAuthRepository(t *testing.T) {
	t.Run("creates repository with default CacheTTL", func(t *testing.T) {
		repo := NewAuthRepository(nil, nil)
		assert.NotNil(t, repo)
		assert.Equal(t, 5*time.Minute, repo.CacheTTL)
	})
}

func TestRedisCacheOperations(t *testing.T) {
	t.Run("cache key is correctly prefixed", func(t *testing.T) {
		repo := &AuthRepository{}
		token := "test-token-123"
		key := repo.getCacheKey(token)
		assert.Equal(t, REDIS_TOKEN_PREFIX+token, key)
	})

	t.Run("different tokens produce different cache keys", func(t *testing.T) {
		repo := &AuthRepository{}
		key1 := repo.getCacheKey("token1")
		key2 := repo.getCacheKey("token2")
		assert.NotEqual(t, key1, key2)
	})
}
