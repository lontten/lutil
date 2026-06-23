package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToJsonStr(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	s, err := ToJsonStr(map[string]int{"a": 1})
	req.NoError(err)
	as.Equal(`{"a":1}`, s)

	_, err = ToJsonStr(make(chan int))
	as.Error(err)
}

func TestToJsonStrPanic(t *testing.T) {
	as := assert.New(t)
	as.Equal(`"ok"`, ToJsonStrPanic("ok"))
	as.Panics(func() { ToJsonStrPanic(make(chan int)) })
}

func TestToJsonStrDefault(t *testing.T) {
	as := assert.New(t)
	as.Equal(`{"a":1}`, ToJsonStrDefault(map[string]int{"a": 1}))
	as.Equal("", ToJsonStrDefault(make(chan int)))
}

func TestToJsonStrP(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	p, err := ToJsonStrP(1)
	req.NoError(err)
	as.Equal(`1`, *p)

	_, err = ToJsonStrP(make(chan int))
	as.Error(err)
}

func TestToJsonStrPPanic(t *testing.T) {
	as := assert.New(t)
	p := ToJsonStrPPanic(true)
	as.Equal(`true`, *p)
	as.Panics(func() { ToJsonStrPPanic(make(chan int)) })
}

func TestToObj(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	v, err := ToObj[int](`42`)
	req.NoError(err)
	as.Equal(42, v)

	_, err = ToObj[int](`not-json`)
	as.Error(err)
}

func TestToObjPanic(t *testing.T) {
	as := assert.New(t)
	as.Equal(42, ToObjPanic[int](`42`))
	as.Panics(func() { ToObjPanic[int](`bad`) })
}

func TestToObjDefault(t *testing.T) {
	as := assert.New(t)
	list := ToObjDefault[[]string](`["a","b"]`)
	as.Equal([]string{"a", "b"}, list)
	as.Equal(0, ToObjDefault[int](`bad`))
}
