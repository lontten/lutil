package moneyutil

import (
	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lutil/decimalutil"
	"github.com/shopspring/decimal"
)

// GrowthPercentage 计算增长百分比
// 公式: ((新值 - 旧值) / 旧值) * 100
// 如果旧值为0会panic
func GrowthPercentage(newValue, oldValue any) decimal.Decimal {
	return decimalutil.GrowthPercentage(types.ToDecimal(newValue), types.ToDecimal(oldValue))
}

// GrowthPercentageWithPrecision 计算增长百分比并四舍五入到指定精度
func GrowthPercentageWithPrecision(newValue, oldValue any, precision int32) decimal.Decimal {
	return decimalutil.GrowthPercentageWithPrecision(types.ToDecimal(newValue), types.ToDecimal(oldValue), precision)
}

// SafeGrowthPercentage 安全的增长百分比计算，如果旧值为0返回默认值
func SafeGrowthPercentage(newValue, oldValue, defaultValue any) decimal.Decimal {
	return decimalutil.SafeGrowthPercentage(types.ToDecimal(newValue), types.ToDecimal(oldValue), types.ToDecimal(defaultValue))
}

// GrowthRate 计算增长率，返回小数形式而不是百分比（例如0.15表示15%）
func GrowthRate(newValue, oldValue any) decimal.Decimal {
	return decimalutil.GrowthRate(types.ToDecimal(newValue), types.ToDecimal(oldValue))
}

// CompoundGrowthRate 计算复合增长率
// 公式: (新值/旧值)^(1/期数) - 1
func CompoundGrowthRate(newValue, oldValue any, periods int32) decimal.Decimal {
	return decimalutil.CompoundGrowthRate(types.ToDecimal(newValue), types.ToDecimal(oldValue), periods)
}
