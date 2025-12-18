package moneyutil

import (
	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
)

// Add 加法运算，支持多个参数连续相加
func Add(values ...any) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	result := types.ToDecimal(values[0])
	for i := 1; i < len(values); i++ {
		result = result.Add(types.ToDecimal(values[i]))
	}
	return result
}

// Sub 减法运算
// 第一个参数减去后面所有参数
func Sub(first any, others ...any) decimal.Decimal {
	result := types.ToDecimal(first)
	for _, value := range others {
		result = result.Sub(types.ToDecimal(value))
	}
	return result
}

// Mul 乘法运算，支持多个参数连续相乘
func Mul(values ...any) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	result := types.ToDecimal(values[0])
	for i := 1; i < len(values); i++ {
		result = result.Mul(types.ToDecimal(values[i]))
	}
	return result
}

// Div 除法运算
// 第一个参数除以后面所有参数
// 如果除数为0会panic
func Div(first any, others ...any) decimal.Decimal {
	result := types.ToDecimal(first)
	for _, value := range others {
		decimalValue := types.ToDecimal(value)
		if decimalValue.IsZero() {
			panic("division by zero")
		}
		result = result.Div(decimalValue)
	}
	return result
}

// DivRound 带精度的除法运算（四舍五入）
// 第一个参数除以后面所有参数，并按照指定精度四舍五入
func DivRound(first any, precision int32, others ...any) decimal.Decimal {
	result := types.ToDecimal(first)
	for _, value := range others {
		decimalValue := types.ToDecimal(value)
		if decimalValue.IsZero() {
			panic("division by zero")
		}
		result = result.DivRound(decimalValue, precision)
	}
	return result
}

// SafeDiv 安全的除法运算，如果除数为0返回默认值而不是panic
func SafeDiv(first any, defaultValue any, others ...any) decimal.Decimal {
	result := types.ToDecimal(first)
	defaultDecimal := types.ToDecimal(defaultValue)

	for _, value := range others {
		decimalValue := types.ToDecimal(value)
		if decimalValue.IsZero() {
			return defaultDecimal
		}
		result = result.Div(decimalValue)
	}
	return result
}

// Abs 绝对值
func Abs(d any) decimal.Decimal {
	return types.ToDecimal(d).Abs()
}

// Neg 取负数
func Neg(d any) decimal.Decimal {
	return types.ToDecimal(d).Neg()
}

// Compare 比较两个decimal的大小
// 返回 -1 如果 d1 < d2
// 返回 0 如果 d1 == d2
// 返回 1 如果 d1 > d2
func Compare(d1, d2 any) int {
	return types.ToDecimal(d1).Cmp(types.ToDecimal(d2))
}

// Equal 判断两个decimal是否相等
func Equal(d1, d2 any) bool {
	return types.ToDecimal(d1).Equal(types.ToDecimal(d2))
}

// GreaterThan 判断 d1 是否大于 d2
func GreaterThan(d1, d2 any) bool {
	return types.ToDecimal(d1).GreaterThan(types.ToDecimal(d2))
}

// GreaterThanOrEqual 判断 d1 是否大于等于 d2
func GreaterThanOrEqual(d1, d2 any) bool {
	return types.ToDecimal(d1).GreaterThanOrEqual(types.ToDecimal(d2))
}

// LessThan 判断 d1 是否小于 d2
func LessThan(d1, d2 any) bool {
	return types.ToDecimal(d1).LessThan(types.ToDecimal(d2))
}

// LessThanOrEqual 判断 d1 是否小于等于 d2
func LessThanOrEqual(d1, d2 any) bool {
	return types.ToDecimal(d1).LessThanOrEqual(types.ToDecimal(d2))
}

// Sum 计算多个decimal的和
func Sum(values ...any) decimal.Decimal {
	return Add(values...)
}

// Average 计算多个decimal的平均值
func Average(values ...any) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	sum := Sum(values...)
	return sum.Div(decimal.NewFromInt(int64(len(values))))
}

// Max 返回多个decimal中的最大值
func Max(values ...any) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	max := types.ToDecimal(values[0])
	for _, value := range values[1:] {
		decimalValue := types.ToDecimal(value)
		if decimalValue.GreaterThan(max) {
			max = decimalValue
		}
	}
	return max
}

// Min 返回多个decimal中的最小值
func Min(values ...any) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	min := types.ToDecimal(values[0])
	for _, value := range values[1:] {
		decimalValue := types.ToDecimal(value)
		if decimalValue.LessThan(min) {
			min = decimalValue
		}
	}
	return min
}

// Percentage 计算百分比 (value / total) * 100
func Percentage(value, total any) decimal.Decimal {
	decimalTotal := types.ToDecimal(total)
	if decimalTotal.IsZero() {
		panic("total cannot be zero for percentage calculation")
	}
	return types.ToDecimal(value).Div(decimalTotal).Mul(decimal.NewFromInt(100))
}

// Round 四舍五入到指定精度
func Round(d any, precision int32) decimal.Decimal {
	return types.ToDecimal(d).Round(precision)
}

// RoundBank 银行家舍入法（四舍六入五成双）
func RoundBank(d any, precision int32) decimal.Decimal {
	return types.ToDecimal(d).RoundBank(precision)
}

// RoundDown 向下舍入
func RoundDown(d any, precision int32) decimal.Decimal {
	return types.ToDecimal(d).RoundDown(precision)
}

// RoundUp 向上舍入
func RoundUp(d any, precision int32) decimal.Decimal {
	return types.ToDecimal(d).RoundUp(precision)
}

// RoundCash 现金舍入法（四舍五入到最接近的5的倍数）
func RoundCash(d any, interval uint8) decimal.Decimal {
	return types.ToDecimal(d).RoundCash(interval)
}
