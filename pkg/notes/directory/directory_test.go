package directory_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
	"testing"
)

func TestModelCreateNewDirectorySuccess(t *testing.T) {
	testModelCreateNewDirectorySuccess(t, nil, "root")
}

func TestModelCreateNewDirectoryFailure(t *testing.T) {
	testModelCreateNewDirectorySuccess(t,
		liberr.WithArgs(
			liberr.Operation("Directory.validate"),
			liberr.ValidationError,
			errors.New("directory name cannot be empty"),
		), "")
}

func testModelCreateNewDirectorySuccess(t *testing.T, expectedError error, name string) {
	_, err := directory.NewDirectory(name)
	assert.Equal(t, expectedError, err)
}
