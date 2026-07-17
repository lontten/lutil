package imgutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetImgSrcForHtml(t *testing.T) {
	as := assert.New(t)

	t.Run("extracts all src attributes", func(t *testing.T) {
		html := `<div><img src="https://a.com/1.png"><img src="/local/2.jpg"><p>no</p></div>`
		srcs, err := GetImgSrcForHtml(html)
		require.NoError(t, err)
		as.Equal([]string{"https://a.com/1.png", "/local/2.jpg"}, srcs)
	})

	t.Run("empty when no images", func(t *testing.T) {
		srcs, err := GetImgSrcForHtml(`<p>hello</p>`)
		require.NoError(t, err)
		as.Empty(srcs)
	})

	t.Run("skips img without src", func(t *testing.T) {
		srcs, err := GetImgSrcForHtml(`<img alt="x"><img src="ok.png">`)
		require.NoError(t, err)
		as.Equal([]string{"ok.png"}, srcs)
	})

	t.Run("empty html", func(t *testing.T) {
		srcs, err := GetImgSrcForHtml("")
		require.NoError(t, err)
		as.Empty(srcs)
	})
}
