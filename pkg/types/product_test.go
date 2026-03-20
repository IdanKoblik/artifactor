package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTokens(t *testing.T) {
	t.Run("hashes all token keys", func(t *testing.T) {
		p := &Product{
			Name: "testproduct",
			Tokens: map[string]TokenPermissions{
				"plaintoken": {Download: true},
			},
		}
		p.HashTokens()

		assert.NotContains(t, p.Tokens, "plaintoken")
		assert.Len(t, p.Tokens, 1)
	})

	t.Run("preserves permissions after hashing", func(t *testing.T) {
		perms := TokenPermissions{Download: true, Upload: true, Maintainer: false, Delete: false}
		p := &Product{
			Name:   "testproduct",
			Tokens: map[string]TokenPermissions{"mytoken": perms},
		}
		p.HashTokens()

		for _, v := range p.Tokens {
			assert.Equal(t, perms, v)
		}
	})

	t.Run("empty tokens stays empty", func(t *testing.T) {
		p := &Product{
			Name:   "testproduct",
			Tokens: map[string]TokenPermissions{},
		}
		p.HashTokens()
		assert.Empty(t, p.Tokens)
	})
}
