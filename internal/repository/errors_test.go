package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	err := NewNotFoundError("12345")

	assert.Error(t, err)
	assert.IsType(t, &NotFoundError{}, err)
	assert.Equal(t, "Order not found", err.Error())
	assert.Equal(t, "12345", err.(*NotFoundError).Number)
}

func TestDuplicateError(t *testing.T) {
	err := NewDuplicateError()

	assert.Error(t, err)
	assert.IsType(t, &DuplicateError{}, err)
	assert.Equal(t, "Login already exists", err.Error())
}

func TestShouldBePositiveError(t *testing.T) {
	err := NewShouldBePositiveError()

	assert.Error(t, err)
	assert.IsType(t, &ShouldBePositiveError{}, err)
	assert.Equal(t, "Sum should be positive", err.Error())
}
