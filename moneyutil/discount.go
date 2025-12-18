package moneyutil

import (
	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
)

// DiscountedPrice 计算打折后的价格
// discount: 折扣整数，如85表示85折，100表示不打折
func DiscountedPrice(originalPrice any, discount int) decimal.Decimal {
	if discount <= 0 {
		return decimal.Zero
	}

	decimalOriginal := types.ToDecimal(originalPrice)
	discountRate := decimal.NewFromInt(int64(discount)).Div(decimal.NewFromInt(100))
	return decimalOriginal.Mul(discountRate)
}

// DiscountedPriceRound 计算打折后的价格并四舍五入到指定精度
func DiscountedPriceRound(originalPrice any, discount int, precision int32) decimal.Decimal {
	price := DiscountedPrice(originalPrice, discount)
	return price.Round(precision)
}

// DiscountRatePercentage 计算优惠率
// 返回优惠百分比，如15表示15% off (85折)
func DiscountRatePercentage(currentPrice, originalPrice any) int {
	return 100 - DiscountRate(currentPrice, originalPrice)
}

// DiscountRate 计算折扣率
// 返回折扣整数，如85表示85折，100表示不打折
func DiscountRate(currentPrice, originalPrice any) int {
	decimalOriginal := types.ToDecimal(originalPrice)
	if decimalOriginal.IsZero() {
		return 100 // 如果原价为0，默认不打折
	}

	decimalCurrent := types.ToDecimal(currentPrice)
	rate := decimalCurrent.Div(decimalOriginal).Mul(decimal.NewFromInt(100))

	// 四舍五入到整数
	rounded := rate.Round(0)
	return int(rounded.IntPart())
}

// SafeDiscountedPrice 安全的打折价格计算
// 如果折扣不在有效范围内(0-100)，返回默认值
func SafeDiscountedPrice(originalPrice any, discount int, defaultValue any) decimal.Decimal {
	if discount < 0 || discount > 100 {
		return types.ToDecimal(defaultValue)
	}

	return DiscountedPrice(originalPrice, discount)
}

// IsValidDiscount 检查折扣是否有效
// 有效范围: 1-100
func IsValidDiscount(discount int) bool {
	return discount >= 1 && discount <= 100
}

// GetDiscountDescription 获取折扣描述
// 如 "85折", "100折(不打折)"
func GetDiscountDescription(discount int) string {
	if discount == 100 {
		return "不打折"
	}
	return decimal.NewFromInt(int64(discount)).String() + "折"
}

// DiscountAmount 计算折扣金额
// 返回原价 - 折后价
func DiscountAmount(originalPrice any, discount int) decimal.Decimal {
	decimalOriginal := types.ToDecimal(originalPrice)
	discountedPrice := DiscountedPrice(originalPrice, discount)
	return decimalOriginal.Sub(discountedPrice)
}

// GetDiscountTier 获取折扣档次
// 根据折扣值返回折扣档次描述
func GetDiscountTier(discount int) string {
	switch {
	case discount >= 90 && discount <= 100:
		return "轻微折扣"
	case discount >= 80 && discount < 90:
		return "中等折扣"
	case discount >= 70 && discount < 80:
		return "较大折扣"
	case discount >= 50 && discount < 70:
		return "大幅折扣"
	case discount < 50:
		return "深度折扣"
	default:
		return "未知档次"
	}
}

// BulkDiscount 计算批量折扣
// 根据购买数量计算折扣，数量越多折扣越大
func BulkDiscount(originalPrice any, quantity int, discountTiers map[int]int) decimal.Decimal {
	// 默认折扣档次
	defaultTiers := map[int]int{
		1:  100, // 1件不打折
		5:  95,  // 5件95折
		10: 90,  // 10件9折
		20: 85,  // 20件85折
		50: 80,  // 50件8折
	}

	// 如果提供了自定义折扣档次，使用自定义的
	if discountTiers != nil {
		defaultTiers = discountTiers
	}

	// 找到适用的折扣
	applicableDiscount := 100 // 默认不打折
	minQuantity := -1

	for tierQuantity, tierDiscount := range defaultTiers {
		if quantity >= tierQuantity && tierQuantity > minQuantity {
			minQuantity = tierQuantity
			applicableDiscount = tierDiscount
		}
	}

	return DiscountedPrice(originalPrice, applicableDiscount)
}
