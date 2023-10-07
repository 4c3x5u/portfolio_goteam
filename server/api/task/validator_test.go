//go:build utest

package task

import (
	"server/assert"
	"testing"
)

// TestTitleValidator tests the Validate method of TitleValidator to assert that
// it returns the correct error message based on task title it's given.
func TestTitleValidator(t *testing.T) {
	sut := NewTitleValidator()

	t.Run("TitleEmpty", func(t *testing.T) {
		err := sut.Validate("")
		if assertErr := assert.Equal(
			"Task title cannot be empty.", err.Error(),
		); assertErr != nil {
			t.Error(err)
		}
	})
}
