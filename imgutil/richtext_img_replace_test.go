package imgutil

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageReplacer_ReplaceRichText_NoRemoteImages(t *testing.T) {
	as := assert.New(t)
	called := false
	r := NewImageReplacer(
		func(localPath string) (string, error) {
			called = true
			return "", nil
		},
		func(remotePath string) (string, error) {
			called = true
			return "", nil
		},
	)

	out, err := r.ReplaceRichText(`<p>hello <img src="/local.png"></p>`)
	require.NoError(t, err)
	as.Contains(out, `/local.png`)
	as.False(called)
}

func TestImageReplacer_ReplaceRichText_EmptyImages(t *testing.T) {
	as := assert.New(t)
	r := NewImageReplacer(
		func(string) (string, error) { return "", nil },
		func(string) (string, error) { return "", nil },
	)

	html := `<p>no images</p>`
	out, err := r.ReplaceRichText(html)
	require.NoError(t, err)
	as.Equal(html, out)
}

func TestImageReplacer_ReplaceRichText_RemoteWithMocks(t *testing.T) {
	as := assert.New(t)

	r := NewImageReplacer(
		func(localPath string) (string, error) {
			return "https://cdn.example.com/uploaded.jpg", nil
		},
		func(remotePath string) (string, error) {
			f, err := os.CreateTemp("", "imgutil-test-*.jpg")
			require.NoError(t, err)
			require.NoError(t, f.Close())
			return f.Name(), nil
		},
		WithTimeout(2*time.Second),
		WithMaxRetries(0),
		WithConcurrent(2),
	)

	html := `<div><img src="https://example.com/a.jpg"><img src="data:image/png;base64,xxx"></div>`
	out, err := r.ReplaceRichText(html)
	require.NoError(t, err)
	as.Contains(out, "https://cdn.example.com/uploaded.jpg")
	as.Contains(out, "data:image/png;base64,xxx")
	as.False(strings.Contains(out, "https://example.com/a.jpg"))
}

func TestImageReplacer_isRemoteURL(t *testing.T) {
	as := assert.New(t)
	r := NewImageReplacer(nil, nil)

	as.True(r.isRemoteURL("https://example.com/a.png"))
	as.True(r.isRemoteURL("http://example.com/a.png"))
	as.False(r.isRemoteURL("/local/a.png"))
	as.False(r.isRemoteURL("data:image/png;base64,xx"))
	as.False(r.isRemoteURL(":bad"))
}

func TestImageReplacer_extractBodyContent(t *testing.T) {
	as := assert.New(t)
	r := NewImageReplacer(nil, nil)

	as.Equal("<p>hi</p>", r.extractBodyContent("<html><body><p>hi</p></body></html>"))
	as.Equal("<p>plain</p>", r.extractBodyContent("<p>plain</p>"))
}

func TestImageReplacer_Options(t *testing.T) {
	as := assert.New(t)
	r := NewImageReplacer(
		nil, nil,
		WithMaxSize(1024),
		WithTimeout(time.Second),
		WithConcurrent(3),
		WithMaxRetries(1),
		WithAllowedTypes([]string{"image/png"}),
	)
	as.Equal(int64(1024), r.MaxSize)
	as.Equal(time.Second, r.Timeout)
	as.Equal(3, r.Concurrent)
	as.Equal(1, r.MaxRetries)
	as.Equal([]string{"image/png"}, r.AllowedTypes)
}
