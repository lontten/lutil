package strutil

import (
	"fmt"
	"math"
	"strconv"
)
import "golang.org/x/exp/constraints"

// Num2Str 将任意整数类型的数值转换为字符串
func Num2Str[T constraints.Integer](num T) string {
	return fmt.Sprintf("%d", num)
}
func Num2StrP[T constraints.Integer](num T) *string {
	var s = Num2Str(num)
	return &s
}

// Str2Num 将字符串转换为int64数值，支持任意进制解析
func Str2Num[T constraints.Integer](str string) (T, error) {
	// 先拿到目标类型的“零值”用来后面类型开关
	var zero T

	// 根据 T 的类别走不同解析路径
	switch any(zero).(type) {
	case int, int8, int16, int32, int64: // 有符号
		v, err := strconv.ParseInt(str, 0, 64)
		if err != nil {
			return zero, fmt.Errorf("无法将字符串 %q 转换为有符号整数: %w", str, err)
		}
		// 越界检查
		if !inRangeSigned(v, zero) {
			return zero, fmt.Errorf("值 %d 超出 %T 范围", v, zero)
		}
		return T(v), nil

	case uint, uint8, uint16, uint32, uint64: // 无符号
		v, err := strconv.ParseUint(str, 0, 64)
		if err != nil {
			return zero, fmt.Errorf("无法将字符串 %q 转换为无符号整数: %w", str, err)
		}
		if !inRangeUnsigned(v, zero) {
			return zero, fmt.Errorf("值 %d 超出 %T 范围", v, zero)
		}
		return T(v), nil

	default: // 理论上进不来，constraints.Integer 已经限定
		return zero, fmt.Errorf("unsupported integer type: %T", zero)
	}
}

// ---------- 越界辅助函数 ----------
func inRangeSigned(v int64, t any) bool {
	switch t.(type) {
	case int8:
		return v >= math.MinInt8 && v <= math.MaxInt8
	case int16:
		return v >= math.MinInt16 && v <= math.MaxInt16
	case int32:
		return v >= math.MinInt32 && v <= math.MaxInt32
	case int64, int:
		return true // 64 位已经顶格
	}
	return false
}

func inRangeUnsigned(v uint64, t any) bool {
	switch t.(type) {
	case uint8:
		return v <= math.MaxUint8
	case uint16:
		return v <= math.MaxUint16
	case uint32:
		return v <= math.MaxUint32
	case uint64, uint:
		return true
	}
	return false
}

func Str2NumMust[T constraints.Integer](str string) T {
	num, err := Str2Num[T](str)
	if err != nil {
		panic(err)
	}
	return num
}

func Str2NumMustP[T constraints.Integer](str string) *T {
	num := Str2NumMust[T](str)
	return &num
}
