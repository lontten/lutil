package structutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type copySrc struct {
	Name string
	Age  int
	Skip string
}

type copyDst struct {
	Name string
	Age  int
	Skip string
}

type copyDstPartial struct {
	Name string
}

func TestCopyStruct(t *testing.T) {
	as := assert.New(t)
	src := &copySrc{Name: "a", Age: 18, Skip: "x"}
	dst := &copyDst{Name: "b", Age: 0, Skip: "y"}
	CopyStruct(src, dst, nil)
	as.Equal("a", dst.Name)
	as.Equal(18, dst.Age)
	as.Equal("x", dst.Skip)
}

func TestCopyStruct_exclude(t *testing.T) {
	as := assert.New(t)
	src := &copySrc{Name: "a", Age: 18, Skip: "x"}
	dst := &copyDst{Name: "b", Age: 0, Skip: "y"}
	CopyStruct(src, dst, []string{"Skip"})
	as.Equal("a", dst.Name)
	as.Equal(18, dst.Age)
	as.Equal("y", dst.Skip)
}

func TestCopyStruct_missingField(t *testing.T) {
	as := assert.New(t)
	src := &copySrc{Name: "a", Age: 18}
	dst := &copyDstPartial{Name: "b"}
	CopyStruct(src, dst, nil)
	as.Equal("a", dst.Name)
}
