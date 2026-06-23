package dateutil

import (
	"testing"

	"github.com/lontten/lcore/v2/types"
	"github.com/stretchr/testify/assert"
)

func mustDate(y, m, d int) types.LocalDate {
	return types.LocalDateOfYmd(y, m, d)
}

func TestMaxMin(t *testing.T) {
	as := assert.New(t)
	d1 := mustDate(2024, 1, 1)
	d2 := mustDate(2024, 6, 1)
	d3 := mustDate(2024, 12, 1)

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
	d1 := mustDate(2024, 1, 1)
	d2 := mustDate(2024, 6, 1)
	var nilDate *types.LocalDate

	as.Nil(MaxP(nilDate))
	as.Nil(MinP(nilDate))

	as.Equal(&d2, MaxP(&d1, nil, &d2))
	as.Equal(&d1, MinP(&d1, nil, &d2))
}

func TestMaxNowMinNow(t *testing.T) {
	as := assert.New(t)
	future := mustDate(2099, 1, 1)
	past := mustDate(2000, 1, 1)

	as.Equal(future, MaxNow(&past, nil, &future))
	as.Equal(past, MinNow(&past, nil, &future))
}

func TestMaxNowPMinNowP(t *testing.T) {
	as := assert.New(t)
	future := mustDate(2099, 1, 1)
	p := MaxNowP(&future)
	as.NotNil(p)
	as.True(p.After(types.NowDate()))

	past := mustDate(2000, 1, 1)
	p2 := MinNowP(&past)
	as.NotNil(p2)
	as.True(p2.Before(types.NowDate()))
}

func TestIsNowRL(t *testing.T) {
	as := assert.New(t)
	var nilDate *types.LocalDate
	as.False(IsNowR(nilDate))
	as.False(IsNowL(nilDate))

	future := mustDate(2099, 1, 1)
	past := mustDate(2000, 1, 1)
	as.True(IsNowR(&future))
	as.True(IsNowL(&past))
}
