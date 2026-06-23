package moneyutil

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAddSubMul(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(6).Equal(Add(1, 2, 3)))
	as.True(decimal.NewFromInt(4).Equal(Sub(10, 3, 3)))
	as.True(decimal.Zero.Equal(Mul()))
	as.True(decimal.NewFromInt(24).Equal(Mul(2, 3, 4)))
}

func TestDiv(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(5).Equal(Div(10, 2)))
	as.Panics(func() { Div(10, 0) })
	as.True(decimal.NewFromInt(-1).Equal(SafeDiv(10, -1, 0)))
}

func TestSumAverageMaxMin(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(6).Equal(Sum(1, 2, 3)))
	as.True(decimal.NewFromInt(2).Equal(Average(1, 2, 3)))
	as.True(decimal.NewFromInt(3).Equal(Max(1, 3, 2)))
	as.True(decimal.NewFromInt(1).Equal(Min(1, 3, 2)))
}

func TestPercentage(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(25).Equal(Percentage(25, 100)))
	as.Panics(func() { Percentage(1, 0) })
}
