package moneyutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDiscountDescription(t *testing.T) {
	as := assert.New(t)
	description := GetDiscountDescription(0)
	as.Equal("0折", description)

	description = GetDiscountDescription(100)
	as.Equal("不打折", description)

	description = GetDiscountDescription(85)
	as.Equal("85折", description)

	description = GetDiscountDescription(7)
	as.Equal("7折", description)
}
