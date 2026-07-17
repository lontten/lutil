package imgutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBytesToBase64(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	s, err := BytesToBase64([]byte("hi"))
	req.NoError(err)
	as.Equal("aGk=", s)

	_, err = BytesToBase64(nil)
	as.Error(err)
}

func TestBytesToURLSafeBase64(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	s, err := BytesToURLSafeBase64([]byte("?"))
	req.NoError(err)
	as.NotEmpty(s)
}
