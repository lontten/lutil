package netutil

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	data, err := os.ReadFile(path)
	req.NoError(err)
	assert.Equal(t, []byte("download"), data)
}
