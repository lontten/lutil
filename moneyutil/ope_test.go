package moneyutil

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGrowthPercentage(t *testing.T) {
	as := assert.New(t)
	got := GrowthPercentage(120, 100)
	as.True(decimal.NewFromInt(20).Equal(got))
	as.Panics(func() { GrowthPercentage(120, 0) })
	as.True(decimal.NewFromInt(20).Equal(SafeGrowthPercentage(120, 100, -1)))
}

func TestGrowthRate(t *testing.T) {
	as := assert.New(t)
	got := GrowthRate(115, 100)
	as.True(decimal.RequireFromString("0.15").Equal(got))
	as.Panics(func() { GrowthRate(115, 0) })
}

func TestCompoundGrowthRate(t *testing.T) {
	as := assert.New(t)
	got := CompoundGrowthRate(121, 100, 2)
	as.True(got.GreaterThan(decimal.Zero))
	as.Panics(func() { CompoundGrowthRate(121, 0, 2) })
	as.Panics(func() { CompoundGrowthRate(121, 100, 0) })
}
