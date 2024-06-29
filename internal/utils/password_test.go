package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	p, err := GenerateHash("password", nil)
	assert.NoError(t, err)

	hash, err := DecodeHash(p.Encode())
	assert.NoError(t, err)

	assert.Equal(t, hash, p)

	match, err := hash.Compare("password")
	assert.NoError(t, err)
	assert.Equal(t, match, true)

	match, err = hash.Compare("password123")
	assert.NoError(t, err)
	assert.Equal(t, match, false)
}
