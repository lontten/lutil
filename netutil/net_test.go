package netutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type echoResp struct {
	Msg string `json:"msg"`
}

func TestGet_OK(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_ = json.NewEncoder(w).Encode(echoResp{Msg: "hello"})
	}))
	defer srv.Close()

	got, err := Get[echoResp](srv.URL)
	req.NoError(err)
	assert.Equal(t, "hello", got.Msg)
}

func TestGet_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Get[echoResp](srv.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestPostJson_OK(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.Header.Get("Content-Type"), "application/json")
		var body map[string]any
		req.NoError(json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "world", body["msg"])
		_ = json.NewEncoder(w).Encode(echoResp{Msg: "ok"})
	}))
	defer srv.Close()

	code, got, err := PostJson[echoResp](srv.URL, map[string]string{"msg": "world"})
	req.NoError(err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "ok", got.Msg)
}

func TestPostJson_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad"}`))
	}))
	defer srv.Close()

	_, _, err := PostJson[echoResp](srv.URL, map[string]string{"a": "b"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
}

func TestPostForm_OK(t *testing.T) {
	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		req.NoError(r.ParseForm())
		assert.Equal(t, "v1", r.Form.Get("k1"))
		_ = json.NewEncoder(w).Encode(echoResp{Msg: "form"})
	}))
	defer srv.Close()

	got, err := PostFormOk[echoResp](srv.URL, url.Values{"k1": {"v1"}})
	req.NoError(err)
	assert.Equal(t, "form", got.Msg)
}

func TestGet_Timeout(t *testing.T) {
	old := defaultClient
	defaultClient = &http.Client{Timeout: 50 * time.Millisecond}
	defer func() { defaultClient = old }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_ = json.NewEncoder(w).Encode(echoResp{Msg: "late"})
	}))
	defer srv.Close()

	_, err := Get[echoResp](srv.URL)
	require.Error(t, err)
}

func TestPostJsonByte_JSONVsBinary(t *testing.T) {
	req := require.New(t)

	t.Run("json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(echoResp{Msg: "err-shape"})
		}))
		defer srv.Close()

		code, raw, got, err := PostJsonByte[echoResp](srv.URL, nil)
		req.NoError(err)
		assert.Equal(t, http.StatusOK, code)
		assert.Nil(t, raw)
		assert.Equal(t, "err-shape", got.Msg)
	})

	t.Run("binary", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write([]byte{0x01, 0x02})
		}))
		defer srv.Close()

		code, raw, _, err := PostJsonByte[echoResp](srv.URL, map[string]string{})
		req.NoError(err)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, []byte{0x01, 0x02}, raw)
	})
}
