package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenCreateAndValidate(t *testing.T) {
	const UserID = 1
	token, err := CreateToken(UserID)
	assert.NoError(t, err)

	bearStr := "Bearer " + token
	userID, err := GetUserID(bearStr)
	assert.NoError(t, err)
	assert.Equal(t, UserID, userID)
}
