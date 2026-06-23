package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLooseNumber(t *testing.T) {
	as := assert.New(t)
	as.True(IsLooseNumber("1,234.56"))
	as.True(IsLooseNumber(" +12 "))
	as.False(IsLooseNumber("12a"))
	as.False(IsLooseNumber(""))
	as.False(IsNotLooseNumber("123"))
	as.True(IsNotLooseNumber("abc"))
}
