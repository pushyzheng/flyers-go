package flyers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	assert.True(t, NewFile("~/.ssh/id_rsa.pub").IsFile())
	assert.False(t, NewFile("~/.ssh/id_rsa.pub").IsDir())

	assert.True(t, NewFile("~/.ssh").IsDir())
	assert.False(t, NewFile("~/.ssh").IsFile())
}
