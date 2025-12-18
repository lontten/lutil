package structutil

import (
	"testing"

	"github.com/lontten/lcore/v2/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type TestUser1 struct {
	Uid   int64
	Name  string
	Age   int
	Money decimal.Decimal
	Time  types.LocalTime
	Date  types.LocalDate
}

func TestStruct2StringMap(t *testing.T) {
	as := assert.New(t)
	nowTime := types.NowTime()
	nowDate := types.NowDate()
	var u = TestUser1{
		Uid:   1,
		Name:  "lontten",
		Age:   18,
		Money: decimal.NewFromFloat(100.01),
		Time:  nowTime,
		Date:  nowDate,
	}
	m, err := Struct2StringMap(u)
	as.Nil(err)
	as.Equal("1", m["Uid"])
	as.Equal("lontten", m["Name"])
	as.Equal("18", m["Age"])
	as.Equal("100.01", m["Money"])
	as.Equal(nowTime.String(), m["Time"])
	as.Equal(nowDate.String(), m["Date"])
}
