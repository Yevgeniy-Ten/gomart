package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndComparePasswords(t *testing.T) {
	tests := []string{
		"first password",
		"second password",
	}
	for _, tt := range tests {
		hashed, err := HashPassword(tt)
		assert.NoError(t, err)
		assert.True(t, ComparePasswords(hashed, tt))
	}
}
