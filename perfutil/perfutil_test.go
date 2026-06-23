package perfutil

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerfTime(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	p := NewPerfTime("test")
	req.NotNil(p)

	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	p.Reset()
	as.True(time.Since(before) >= 10*time.Millisecond)

	old := os.Stdout
	r, w, err := os.Pipe()
	req.NoError(err)
	os.Stdout = w

	p.Mark("step1")
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	req.NoError(err)
	as.Contains(buf.String(), "[PERF] test step1:")
}

func TestNewPerfTimeJoinName(t *testing.T) {
	p := NewPerfTime("a", "b")
	assert.NotNil(t, p)
}
