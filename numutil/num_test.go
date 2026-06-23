package numutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound45(t *testing.T) {
	as := assert.New(t)
	as.Equal(int64(1), Round45(0.5))
	as.Equal(int64(1), Round45(0.51))
	as.Equal(int64(0), Round45(0.49))
	as.Equal(int64(-1), Round45(-0.5))
	as.Equal(int64(-1), Round45(-0.51))
	as.Equal(int64(0), Round45(-0.49))
	as.Equal(int64(2), Round45(1.5))
}

func TestRound45i(t *testing.T) {
	as := assert.New(t)
	as.Equal(1, Round45i(0.5))
	as.Equal(1, Round45i(0.51))
	as.Equal(0, Round45i(0.49))
	as.Equal(-1, Round45i(-0.5))
	as.Equal(2, Round45i(1.5))
}
