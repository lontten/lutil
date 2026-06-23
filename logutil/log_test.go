package logutil

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type user struct {
	Name string
	Age  int
}

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestLog_struct(t *testing.T) {
	out := captureStdout(func() {
		Log(&user{Name: "john", Age: 22})
	})
	assert.Contains(t, out, `"Name": "john"`)
	assert.Contains(t, out, `"Age": 22`)
}

func TestLog_mapAndSlice(t *testing.T) {
	out := captureStdout(func() {
		Log(map[string]int{"a": 1})
		Log([]int{1, 2, 3})
	})
	assert.Contains(t, out, `"a": 1`)
	assert.Contains(t, out, "1")
}

func TestLog_unmarshalable(t *testing.T) {
	out := captureStdout(func() {
		Log(make(chan int))
	})
	assert.Contains(t, out, "json.MarshalIndent err:")
}

func TestLog_multipleArgs(t *testing.T) {
	req := require.New(t)
	out := captureStdout(func() {
		Log(1, "ok")
	})
	req.Contains(out, "1")
	req.Contains(out, `"ok"`)
}
