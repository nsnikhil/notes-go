package directory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"notes/pkg/liberr"
	"testing"
)

func TestModelCreateNewDirectorySuccess(t *testing.T) {
	testModelCreateNewDirectorySuccess(t, nil, "root")
}

func TestModelCreateNewDirectoryFailure(t *testing.T) {
	testModelCreateNewDirectorySuccess(t,
		liberr.WithArgs(
			liberr.Operation("directory.validate"),
			liberr.ValidationError,
			errors.New("directory name cannot be empty"),
		), "")
}

func testModelCreateNewDirectorySuccess(t *testing.T, expectedError error, name string) {
	_, err := newDirectory(name)
	assert.Equal(t, expectedError, err)
}
