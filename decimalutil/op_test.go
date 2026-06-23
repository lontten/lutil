package decimalutil

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func d(v string) decimal.Decimal {
	return decimal.RequireFromString(v)
}

func TestAddSubMul(t *testing.T) {
	as := assert.New(t)
	as.True(Equal(d("0"), Add()))
	as.True(Equal(d("6"), Add(d("1"), d("2"), d("3"))))
	as.True(Equal(d("4"), Sub(d("10"), d("3"), d("3"))))
	as.True(Equal(d("0"), Mul()))
	as.True(Equal(d("24"), Mul(d("2"), d("3"), d("4"))))
}

func TestDiv(t *testing.T) {
	as := assert.New(t)
	as.True(Equal(d("5"), Div(d("10"), d("2"))))
	as.Panics(func() { Div(d("10"), d("0")) })
	as.True(Equal(d("-1"), SafeDiv(d("10"), d("-1"), d("0"))))
}

func TestCompareAndRound(t *testing.T) {
	as := assert.New(t)
	as.Equal(1, Compare(d("3"), d("2")))
	as.True(GreaterThan(d("3"), d("2")))
	as.True(LessThan(d("2"), d("3")))
	as.True(Equal(d("1.005").Round(2), Round(d("1.005"), 2)))
}

func TestSumAverageMaxMin(t *testing.T) {
	as := assert.New(t)
	as.True(Equal(d("6"), Sum(d("1"), d("2"), d("3"))))
	as.True(Equal(d("2"), Average(d("1"), d("2"), d("3"))))
	as.True(Equal(d("0"), Average()))
	as.True(Equal(d("3"), Max(d("1"), d("3"), d("2"))))
	as.True(Equal(d("1"), Min(d("1"), d("3"), d("2"))))
}

func TestPercentage(t *testing.T) {
	as := assert.New(t)
	as.True(Equal(d("25"), Percentage(d("25"), d("100"))))
	as.Panics(func() { Percentage(d("1"), d("0")) })
}

func TestAbsNeg(t *testing.T) {
	as := assert.New(t)
	as.True(Equal(d("5"), Abs(d("-5"))))
	as.True(Equal(d("-5"), Neg(d("5"))))
}
