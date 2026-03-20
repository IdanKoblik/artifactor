package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUIFlag(t *testing.T) {
	flag := UIFlag()
	assert.Equal(t, "--ui", flag.Cmd)
	assert.Equal(t, "ui", flag.Name)
	assert.NotNil(t, flag.Handle)
}

func TestUIFlag_Handle(t *testing.T) {
	flag := UIFlag()
	err := flag.Handle([]string{})
	assert.NoError(t, err)
}
