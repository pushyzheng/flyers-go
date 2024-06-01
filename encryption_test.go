package flyers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encrypt(t *testing.T) {
	s, err := AESEncrypt("hello", "thisis32bitlongpassphraseimusing")
	assert.Nil(t, err)
	assert.NotEmpty(t, s)
	fmt.Println(s)
}

func Test_decrypt(t *testing.T) {
	s, err := AESDecrypt("z-Uc41z3Jdxz_20ElUvrs6antoza", "thisis32bitlongpassphraseimusing")
	assert.Nil(t, err)
	assert.Equal(t, "hello", s)
}
