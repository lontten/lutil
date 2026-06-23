package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasText(t *testing.T) {
	as := assert.New(t)
	as.True(HasText(" a "))
	as.False(HasText("   "))
	as.False(HasTextP(nil))
	s := "x"
	as.True(HasTextP(&s))
}

func TestStrContains(t *testing.T) {
	as := assert.New(t)
	as.True(StrContainsAll("hello world", "hello", "world"))
	as.False(StrContainsAll("hello", "hello", "x"))
	as.True(StrContainsAny("hello", "x", "ell"))
	as.False(StrContainsAny("hello", "x", "y"))
}

func TestStrRight(t *testing.T) {
	as := assert.New(t)
	as.Equal("world", FirstStrRight("hello world", "hello "))
	as.Equal("baz bar", LastStrRight("foo bar baz bar", "bar "))
	as.Equal("", FirstStrRight("abc", "z"))
}

func TestLowerFirst(t *testing.T) {
	as := assert.New(t)
	as.Equal("hello", LowerFirst("Hello"))
	as.Equal("", LowerFirst(""))
}
