package netutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIPAndVersion(t *testing.T) {
	as := assert.New(t)
	as.True(IsValidIP("127.0.0.1"))
	as.False(IsValidIP("not-ip"))
	as.Equal(4, IPVersion("192.168.0.1"))
	as.Equal(6, IPVersion("::1"))
	as.Equal(0, IPVersion("bad"))
}

func TestRealIP(t *testing.T) {
	as := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.1:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.2")
	ip := RealIP(req, DefaultConfig)
	as.Equal("203.0.113.1", ip)
}

func TestRealIP_trustedProxy(t *testing.T) {
	as := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.5, 10.0.0.2")
	ip := RealIP(req, DefaultConfig)
	as.Equal("198.51.100.5", ip)
}

func TestIPMiddlewareAndContext(t *testing.T) {
	as := assert.New(t)
	handler := IPMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		as.Equal("1.2.3.4", IPFromContext(r.Context()))
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:9999"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	as.Equal("", IPFromContext(context.Background()))
}

func TestRealIPSimple(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "8.8.8.8:53"
	assert.NotEmpty(t, RealIPSimple(req))
}
