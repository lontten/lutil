package datetimeutil

import (
	"testing"

	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

func mustDT(y, m, d, h, min, s int) types.LocalDateTime {
	return types.LocalDateTimeOfYmdHms(y, m, d, h, min, s)
}

func TestMaxMin(t *testing.T) {
	as := assert.New(t)
	d1 := mustDT(2024, 1, 1, 0, 0, 0)
	d2 := mustDT(2024, 6, 1, 0, 0, 0)
	d3 := mustDT(2024, 12, 1, 0, 0, 0)

	as.Equal(d3, Max(d1, d2, d3))
	as.Equal(d1, Min(d1, d2, d3))
}

func TestMaxMin_panic(t *testing.T) {
	as := assert.New(t)
	as.Panics(func() { Max() })
	as.Panics(func() { Min() })
}

func TestMaxPMinP(t *testing.T) {
	as := assert.New(t)
	d1 := mustDT(2024, 1, 1, 0, 0, 0)
	d2 := mustDT(2024, 6, 1, 0, 0, 0)
	var nilDT *types.LocalDateTime

	as.Nil(MaxP(nilDT))
	as.Nil(MinP(nilDT))
	as.Equal(&d2, MaxP(&d1, nil, &d2))
	as.Equal(&d1, MinP(&d1, nil, &d2))
}

func TestMaxNowMinNow(t *testing.T) {
	as := assert.New(t)
	future := mustDT(2099, 1, 1, 0, 0, 0)
	past := mustDT(2000, 1, 1, 0, 0, 0)

	as.Equal(future, MaxNow(&past, nil, &future))
	as.Equal(past, MinNow(&past, nil, &future))
}

func TestIsNowRL(t *testing.T) {
	as := assert.New(t)
	var nilDT *types.LocalDateTime
	as.False(IsNowR(nilDT))
	as.False(IsNowL(nilDT))

	future := mustDT(2099, 1, 1, 0, 0, 0)
	past := mustDT(2000, 1, 1, 0, 0, 0)
	as.True(IsNowR(&future))
	as.True(IsNowL(&past))
}
