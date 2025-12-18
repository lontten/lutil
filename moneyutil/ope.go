package moneyutil

import (
	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
)

// GrowthPercentage 计算增长百分比
// 公式: ((新值 - 旧值) / 旧值) * 100
// 如果旧值为0会panic
func GrowthPercentage(newValue, oldValue any) decimal.Decimal {
	decimalOldValue := types.ToDecimal(oldValue)
	if decimalOldValue.IsZero() {
		panic("old value cannot be zero for growth percentage calculation")
	}

	decimalNewValue := types.ToDecimal(newValue)
	difference := decimalNewValue.Sub(decimalOldValue)
	return difference.Div(decimalOldValue).Mul(decimal.NewFromInt(100))
}

// GrowthPercentageWithPrecision 计算增长百分比并四舍五入到指定精度
func GrowthPercentageWithPrecision(newValue, oldValue any, precision int32) decimal.Decimal {
	percentage := GrowthPercentage(newValue, oldValue)
	return percentage.Round(precision)
}

// SafeGrowthPercentage 安全的增长百分比计算，如果旧值为0返回默认值
func SafeGrowthPercentage(newValue, oldValue, defaultValue any) decimal.Decimal {
	decimalOldValue := types.ToDecimal(oldValue)
	if decimalOldValue.IsZero() {
		return types.ToDecimal(defaultValue)
	}

	decimalNewValue := types.ToDecimal(newValue)
	difference := decimalNewValue.Sub(decimalOldValue)
	return difference.Div(decimalOldValue).Mul(decimal.NewFromInt(100))
}

// GrowthRate 计算增长率，返回小数形式而不是百分比（例如0.15表示15%）
func GrowthRate(newValue, oldValue any) decimal.Decimal {
	decimalOldValue := types.ToDecimal(oldValue)
	if decimalOldValue.IsZero() {
		panic("old value cannot be zero for growth rate calculation")
	}

	decimalNewValue := types.ToDecimal(newValue)
	difference := decimalNewValue.Sub(decimalOldValue)
	return difference.Div(decimalOldValue)
}

// CompoundGrowthRate 计算复合增长率
// 公式: (新值/旧值)^(1/期数) - 1
func CompoundGrowthRate(newValue, oldValue any, periods int32) decimal.Decimal {
	decimalOldValue := types.ToDecimal(oldValue)
	if decimalOldValue.IsZero() {
		panic("old value cannot be zero for compound growth rate calculation")
	}
	if periods <= 0 {
		panic("periods must be positive")
	}

	decimalNewValue := types.ToDecimal(newValue)
	ratio := decimalNewValue.Div(decimalOldValue)
	exponent := decimal.NewFromInt32(1).Div(decimal.NewFromInt32(periods))
	return ratio.Pow(exponent).Sub(decimal.NewFromInt(1))
}
