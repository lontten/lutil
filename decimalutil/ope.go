package decimalutil

import (
	"github.com/shopspring/decimal"
)

// GrowthPercentage 计算增长百分比
// 公式: ((新值 - 旧值) / 旧值) * 100
// 如果旧值为0会panic
func GrowthPercentage(newValue, oldValue decimal.Decimal) decimal.Decimal {
	if oldValue.IsZero() {
		panic("old value cannot be zero for growth percentage calculation")
	}

	difference := newValue.Sub(oldValue)
	return difference.Div(oldValue).Mul(decimal.NewFromInt(100))
}

// GrowthPercentageWithPrecision 计算增长百分比并四舍五入到指定精度
func GrowthPercentageWithPrecision(newValue, oldValue decimal.Decimal, precision int32) decimal.Decimal {
	percentage := GrowthPercentage(newValue, oldValue)
	return percentage.Round(precision)
}

// SafeGrowthPercentage 安全的增长百分比计算，如果旧值为0返回默认值
func SafeGrowthPercentage(newValue, oldValue, defaultValue decimal.Decimal) decimal.Decimal {
	if oldValue.IsZero() {
		return defaultValue
	}

	difference := newValue.Sub(oldValue)
	return difference.Div(oldValue).Mul(decimal.NewFromInt(100))
}

// GrowthRate 计算增长率，返回小数形式而不是百分比（例如0.15表示15%）
func GrowthRate(newValue, oldValue decimal.Decimal) decimal.Decimal {
	if oldValue.IsZero() {
		panic("old value cannot be zero for growth rate calculation")
	}

	difference := newValue.Sub(oldValue)
	return difference.Div(oldValue)
}

// CompoundGrowthRate 计算复合增长率
// 公式: (新值/旧值)^(1/期数) - 1
func CompoundGrowthRate(newValue, oldValue decimal.Decimal, periods int32) decimal.Decimal {
	if oldValue.IsZero() {
		panic("old value cannot be zero for compound growth rate calculation")
	}
	if periods <= 0 {
		panic("periods must be positive")
	}

	ratio := newValue.Div(oldValue)
	exponent := decimal.NewFromInt32(1).Div(decimal.NewFromInt32(periods))
	return ratio.Pow(exponent).Sub(decimal.NewFromInt(1))
}
