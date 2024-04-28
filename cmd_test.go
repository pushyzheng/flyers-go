package flyers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShellCmd(t *testing.T) {
	cmdResult, err := NewShellCmd(context.Background(), "echo 'hello world'").Run()
	assert.Nil(t, err)

	err2 := cmdResult.ScanStdout(func(line string) {
		assert.Equal(t, "hello world", line)
	})
	assert.Nil(t, err2)
}
