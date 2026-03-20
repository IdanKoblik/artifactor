package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	t.Run("valid simple name", func(t *testing.T) {
		assert.NoError(t, ValidateName("myproduct"))
	})

	t.Run("valid name with dots dashes underscores", func(t *testing.T) {
		assert.NoError(t, ValidateName("my-product_1.0"))
	})

	t.Run("starts with digit", func(t *testing.T) {
		assert.NoError(t, ValidateName("1product"))
	})

	t.Run("empty string is invalid", func(t *testing.T) {
		assert.Error(t, ValidateName(""))
	})

	t.Run("starts with special char is invalid", func(t *testing.T) {
		assert.Error(t, ValidateName("-product"))
	})

	t.Run("contains slash is invalid", func(t *testing.T) {
		assert.Error(t, ValidateName("my/product"))
	})

	t.Run("contains space is invalid", func(t *testing.T) {
		assert.Error(t, ValidateName("my product"))
	})
}

func TestSafeFilename(t *testing.T) {
	t.Run("valid filename", func(t *testing.T) {
		result, err := SafeFilename("myfile.txt")
		assert.NoError(t, err)
		assert.Equal(t, "myfile.txt", result)
	})

	t.Run("strips path traversal", func(t *testing.T) {
		result, err := SafeFilename("../../etc/passwd")
		assert.NoError(t, err)
		assert.Equal(t, "passwd", result)
	})

	t.Run("dot is invalid", func(t *testing.T) {
		_, err := SafeFilename(".")
		assert.Error(t, err)
	})

	t.Run("double dot is invalid", func(t *testing.T) {
		_, err := SafeFilename("..")
		assert.Error(t, err)
	})

	t.Run("filename with directory prefix", func(t *testing.T) {
		result, err := SafeFilename("subdir/file.zip")
		assert.NoError(t, err)
		assert.Equal(t, "file.zip", result)
	})
}
