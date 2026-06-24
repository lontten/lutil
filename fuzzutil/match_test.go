package fuzzutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchKindString(t *testing.T) {
	as := assert.New(t)
	as.Equal(MatchNone, MatchKind(0))
	as.NotEqual(MatchContain, MatchFuzzy)
}

func TestLike_edgeCases(t *testing.T) {
	as := assert.New(t)
	as.Equal("", Like("x", nil))
	as.Equal("", Like("深圳市", []string{"上海", "广州"}))
	as.Equal("深圳", Like("深圳市", []string{"深圳", "广州"}))
}
