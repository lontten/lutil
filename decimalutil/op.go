package decimalutil

import (
	"github.com/shopspring/decimal"
)

// Add 加法运算，支持多个参数连续相加
func Add(values ...decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	result := values[0]
	for i := 1; i < len(values); i++ {
		result = result.Add(values[i])
	}
	return result
}

// Sub 减法运算
// 第一个参数减去后面所有参数
func Sub(first decimal.Decimal, others ...decimal.Decimal) decimal.Decimal {
	result := first
	for _, value := range others {
		result = result.Sub(value)
	}
	return result
}

// Mul 乘法运算，支持多个参数连续相乘
func Mul(values ...decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	result := values[0]
	for i := 1; i < len(values); i++ {
		result = result.Mul(values[i])
	}
	return result
}

// Div 除法运算
// 第一个参数除以后面所有参数
// 如果除数为0会panic
func Div(first decimal.Decimal, others ...decimal.Decimal) decimal.Decimal {
	result := first
	for _, value := range others {
		if value.IsZero() {
			panic("division by zero")
		}
		result = result.Div(value)
	}
	return result
}

// DivRound 带精度的除法运算（四舍五入）
// 第一个参数除以后面所有参数，并按照指定精度四舍五入
func DivRound(first decimal.Decimal, precision int32, others ...decimal.Decimal) decimal.Decimal {
	result := first
	for _, value := range others {
		if value.IsZero() {
			panic("division by zero")
		}
		result = result.DivRound(value, precision)
	}
	return result
}

// SafeDiv 安全的除法运算，如果除数为0返回默认值而不是panic
func SafeDiv(first decimal.Decimal, defaultValue decimal.Decimal, others ...decimal.Decimal) decimal.Decimal {
	result := first
	for _, value := range others {
		if value.IsZero() {
			return defaultValue
		}
		result = result.Div(value)
	}
	return result
}

// Abs 绝对值
func Abs(d decimal.Decimal) decimal.Decimal {
	return d.Abs()
}

// Neg 取负数
func Neg(d decimal.Decimal) decimal.Decimal {
	return d.Neg()
}

// Compare 比较两个decimal的大小
// 返回 -1 如果 d1 < d2
// 返回 0 如果 d1 == d2
// 返回 1 如果 d1 > d2
func Compare(d1, d2 decimal.Decimal) int {
	return d1.Cmp(d2)
}

// Equal 判断两个decimal是否相等
func Equal(d1, d2 decimal.Decimal) bool {
	return d1.Equal(d2)
}

// GreaterThan 判断 d1 是否大于 d2
func GreaterThan(d1, d2 decimal.Decimal) bool {
	return d1.GreaterThan(d2)
}

// GreaterThanOrEqual 判断 d1 是否大于等于 d2
func GreaterThanOrEqual(d1, d2 decimal.Decimal) bool {
	return d1.GreaterThanOrEqual(d2)
}

// LessThan 判断 d1 是否小于 d2
func LessThan(d1, d2 decimal.Decimal) bool {
	return d1.LessThan(d2)
}

// LessThanOrEqual 判断 d1 是否小于等于 d2
func LessThanOrEqual(d1, d2 decimal.Decimal) bool {
	return d1.LessThanOrEqual(d2)
}

// Sum 计算多个decimal的和
func Sum(values ...decimal.Decimal) decimal.Decimal {
	return Add(values...)
}

// Average 计算多个decimal的平均值
func Average(values ...decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	sum := Sum(values...)
	return sum.Div(decimal.NewFromInt(int64(len(values))))
}

// Max 返回多个decimal中的最大值
func Max(values ...decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	max := values[0]
	for _, value := range values[1:] {
		if value.GreaterThan(max) {
			max = value
		}
	}
	return max
}

// Min 返回多个decimal中的最小值
func Min(values ...decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}
	min := values[0]
	for _, value := range values[1:] {
		if value.LessThan(min) {
			min = value
		}
	}
	return min
}

// Percentage 计算百分比 (value / total) * 100
func Percentage(value, total decimal.Decimal) decimal.Decimal {
	if total.IsZero() {
		panic("total cannot be zero for percentage calculation")
	}
	return value.Div(total).Mul(decimal.NewFromInt(100))
}

// Round 四舍五入到指定精度
func Round(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.Round(precision)
}

// RoundBank 银行家舍入法（四舍六入五成双）
func RoundBank(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.RoundBank(precision)
}

// RoundDown 向下舍入
func RoundDown(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.RoundDown(precision)
}

// RoundUp 向上舍入
func RoundUp(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.RoundUp(precision)
}

// RoundCash 现金舍入法（四舍五入到最接近的5的倍数）
func RoundCash(d decimal.Decimal, interval uint8) decimal.Decimal {
	return d.RoundCash(interval)
}
