package netutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckFileUrlCanDownload(t *testing.T) {
	as := assert.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("x"))
	}))
	defer srv.Close()

	ok, msg := CheckFileUrlCanDownload(srv.URL + "/file.txt")
	as.True(ok)
	as.Empty(msg)
}

func TestDownloadFileToLocal(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("download"))
	}))
	defer srv.Close()

	path, err := DownloadFileToLocal(srv.URL + "/a.txt")
	req.NoError(err)
	defer os.Remove(path)
	assert.True(t, strings.HasSuffix(path, ".txt"))
	data, err := os.ReadFile(path)
	req.NoError(err)
	assert.Equal(t, []byte("download"), data)
}

func TestDownloadFileToLocal_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := DownloadFileToLocal(srv.URL + "/missing.bin")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestDownloadFileToLocalLimit_WithinLimit(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("small"))
	}))
	defer srv.Close()

	path, err := DownloadFileToLocalLimit(srv.URL+"/a.bin", 1024)
	req.NoError(err)
	defer os.Remove(path)
	data, err := os.ReadFile(path)
	req.NoError(err)
	assert.Equal(t, []byte("small"), data)
}

func TestDownloadFileToLocalLimit_SizeLimit(t *testing.T) {
	const limit int64 = 1024
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(make([]byte, limit+1))
	}))
	defer srv.Close()

	_, err := DownloadFileToLocalLimit(srv.URL+"/big.bin", limit)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "大小限制")
}

func TestDownloadFileToLocalLimit_InvalidMaxBytes(t *testing.T) {
	_, err := DownloadFileToLocalLimit("http://example.com/a.bin", 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maxBytes")

	_, err = DownloadFileToLocalLimit("http://example.com/a.bin", -1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maxBytes")
}

func TestDownloadFileToLocal_QueryDoesNotPolluteSuffix(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	path, err := DownloadFileToLocal(srv.URL + "/doc.pdf?token=abc.xyz")
	req.NoError(err)
	defer os.Remove(path)
	assert.Equal(t, ".pdf", filepath.Ext(path))
}

func TestFileSuffixFromURL(t *testing.T) {
	assert.Equal(t, "txt", fileSuffixFromURL("https://ex.com/a/b.txt"))
	assert.Equal(t, "pdf", fileSuffixFromURL("https://ex.com/doc.pdf?x=1.2"))
	assert.Equal(t, "", fileSuffixFromURL("https://ex.com/noext"))
	assert.Equal(t, "", fileSuffixFromURL("https://ex.com/a.bad!ext"))
	assert.Equal(t, "", fileSuffixFromURL("://bad"))
}
