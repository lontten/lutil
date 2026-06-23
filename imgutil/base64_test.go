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

func TestGetImgSrcForHtml(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	html := `<html><body><img src="a.png"/><img src="b.jpg"/></body></html>`
	srcs, err := GetImgSrcForHtml(html)
	req.NoError(err)
	as.Equal([]string{"a.png", "b.jpg"}, srcs)

	srcs, err = GetImgSrcForHtml("<p>no image</p>")
	req.NoError(err)
	as.Empty(srcs)
}
