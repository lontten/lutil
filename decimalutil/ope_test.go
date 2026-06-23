package decimalutil

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGrowthPercentage(t *testing.T) {
	as := assert.New(t)
	old := decimal.NewFromInt(100)
	newV := decimal.NewFromInt(120)
	as.True(Equal(decimal.NewFromInt(20), GrowthPercentage(newV, old)))
	as.Panics(func() { GrowthPercentage(newV, decimal.Zero) })
	as.True(Equal(decimal.NewFromInt(20), SafeGrowthPercentage(newV, old, decimal.NewFromInt(-1))))
}

func TestGrowthRate(t *testing.T) {
	as := assert.New(t)
	old := decimal.NewFromInt(100)
	newV := decimal.NewFromInt(115)
	as.True(Equal(decimal.RequireFromString("0.15"), GrowthRate(newV, old)))
	as.Panics(func() { GrowthRate(newV, decimal.Zero) })
}

func TestCompoundGrowthRate(t *testing.T) {
	as := assert.New(t)
	old := decimal.NewFromInt(100)
	newV := decimal.NewFromInt(121)
	rate := CompoundGrowthRate(newV, old, 2)
	as.True(rate.GreaterThan(decimal.Zero))
	as.Panics(func() { CompoundGrowthRate(newV, decimal.Zero, 2) })
	as.Panics(func() { CompoundGrowthRate(newV, old, 0) })
}

func TestGrowthPercentageWithPrecision(t *testing.T) {
	as := assert.New(t)
	old := decimal.NewFromInt(3)
	newV := decimal.NewFromInt(4)
	got := GrowthPercentageWithPrecision(newV, old, 2)
	as.True(got.GreaterThan(decimal.Zero))
}
