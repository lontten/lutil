package moneyutil

import (
	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lutil/decimalutil"
	"github.com/shopspring/decimal"
)

func toDecimals(values ...any) []decimal.Decimal {
	if len(values) == 0 {
		return nil
	}
	result := make([]decimal.Decimal, len(values))
	for i, v := range values {
		result[i] = types.ToDecimal(v)
	}
	return result
}

// Add 加法运算，支持多个参数连续相加
func Add(values ...any) decimal.Decimal {
	return decimalutil.Add(toDecimals(values...)...)
}

// Sub 减法运算
// 第一个参数减去后面所有参数
func Sub(first any, others ...any) decimal.Decimal {
	return decimalutil.Sub(types.ToDecimal(first), toDecimals(others...)...)
}

// Mul 乘法运算，支持多个参数连续相乘
func Mul(values ...any) decimal.Decimal {
	return decimalutil.Mul(toDecimals(values...)...)
}

// Div 除法运算
// 第一个参数除以后面所有参数
// 如果除数为0会panic
func Div(first any, others ...any) decimal.Decimal {
	return decimalutil.Div(types.ToDecimal(first), toDecimals(others...)...)
}

// DivRound 带精度的除法运算（四舍五入）
// 第一个参数除以后面所有参数，并按照指定精度四舍五入
func DivRound(first any, precision int32, others ...any) decimal.Decimal {
	return decimalutil.DivRound(types.ToDecimal(first), precision, toDecimals(others...)...)
}

// SafeDiv 安全的除法运算，如果除数为0返回默认值而不是panic
func SafeDiv(first any, defaultValue any, others ...any) decimal.Decimal {
	return decimalutil.SafeDiv(types.ToDecimal(first), types.ToDecimal(defaultValue), toDecimals(others...)...)
}

// Abs 绝对值
func Abs(d any) decimal.Decimal {
	return decimalutil.Abs(types.ToDecimal(d))
}

// Neg 取负数
func Neg(d any) decimal.Decimal {
	return decimalutil.Neg(types.ToDecimal(d))
}

// Compare 比较两个decimal的大小
// 返回 -1 如果 d1 < d2
// 返回 0 如果 d1 == d2
// 返回 1 如果 d1 > d2
func Compare(d1, d2 any) int {
	return decimalutil.Compare(types.ToDecimal(d1), types.ToDecimal(d2))
}

// Equal 判断两个decimal是否相等
func Equal(d1, d2 any) bool {
	return decimalutil.Equal(types.ToDecimal(d1), types.ToDecimal(d2))
}

// GreaterThan 判断 d1 是否大于 d2
func GreaterThan(d1, d2 any) bool {
	return decimalutil.GreaterThan(types.ToDecimal(d1), types.ToDecimal(d2))
}

// GreaterThanOrEqual 判断 d1 是否大于等于 d2
func GreaterThanOrEqual(d1, d2 any) bool {
	return decimalutil.GreaterThanOrEqual(types.ToDecimal(d1), types.ToDecimal(d2))
}

// LessThan 判断 d1 是否小于 d2
func LessThan(d1, d2 any) bool {
	return decimalutil.LessThan(types.ToDecimal(d1), types.ToDecimal(d2))
}

// LessThanOrEqual 判断 d1 是否小于等于 d2
func LessThanOrEqual(d1, d2 any) bool {
	return decimalutil.LessThanOrEqual(types.ToDecimal(d1), types.ToDecimal(d2))
}

// Sum 计算多个decimal的和
func Sum(values ...any) decimal.Decimal {
	return decimalutil.Sum(toDecimals(values...)...)
}

// Average 计算多个decimal的平均值
func Average(values ...any) decimal.Decimal {
	return decimalutil.Average(toDecimals(values...)...)
}

// Max 返回多个decimal中的最大值
func Max(values ...any) decimal.Decimal {
	return decimalutil.Max(toDecimals(values...)...)
}

// Min 返回多个decimal中的最小值
func Min(values ...any) decimal.Decimal {
	return decimalutil.Min(toDecimals(values...)...)
}

// Percentage 计算百分比 (value / total) * 100
func Percentage(value, total any) decimal.Decimal {
	return decimalutil.Percentage(types.ToDecimal(value), types.ToDecimal(total))
}

// Round 四舍五入到指定精度
func Round(d any, precision int32) decimal.Decimal {
	return decimalutil.Round(types.ToDecimal(d), precision)
}

// RoundBank 银行家舍入法（四舍六入五成双）
func RoundBank(d any, precision int32) decimal.Decimal {
	return decimalutil.RoundBank(types.ToDecimal(d), precision)
}

// RoundDown 向下舍入
func RoundDown(d any, precision int32) decimal.Decimal {
	return decimalutil.RoundDown(types.ToDecimal(d), precision)
}

// RoundUp 向上舍入
func RoundUp(d any, precision int32) decimal.Decimal {
	return decimalutil.RoundUp(types.ToDecimal(d), precision)
}

// RoundCash 现金舍入法（四舍五入到最接近的5的倍数）
func RoundCash(d any, interval uint8) decimal.Decimal {
	return decimalutil.RoundCash(types.ToDecimal(d), interval)
}
