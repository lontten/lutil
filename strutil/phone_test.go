package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLandline(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		want        bool
	}{
		{name: "landline with ext", phoneNumber: "028-61555395-8038", want: true},
		{name: "service number is not landline", phoneNumber: "400-6162020", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckLandline(tt.phoneNumber); got != tt.want {
				t.Errorf("CheckLandline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckServicePhone(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		want        bool
	}{
		{name: "400 with hyphen", phoneNumber: "400-6162020", want: true},
		{name: "400 without hyphen", phoneNumber: "4006162020", want: true},
		{name: "400 segmented", phoneNumber: "400-616-2020", want: true},
		{name: "800 with hyphen", phoneNumber: "800-1234567", want: true},
		{name: "too short", phoneNumber: "400-123", want: false},
		{name: "not service prefix", phoneNumber: "123-4567890", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckServicePhone(tt.phoneNumber); got != tt.want {
				t.Errorf("CheckServicePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPhoneAll(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		want        bool
	}{
		{name: "landline", phoneNumber: "028-61555395-8038", want: true},
		{name: "mobile", phoneNumber: "13812345678", want: true},
		{name: "400 with hyphen", phoneNumber: "400-6162020", want: true},
		{name: "400 without hyphen", phoneNumber: "4006162020", want: true},
		{name: "400 segmented", phoneNumber: "400-616-2020", want: true},
		{name: "800", phoneNumber: "800-1234567", want: true},
		{name: "short mobile", phoneNumber: "1883711386", want: false},
		{name: "short 400", phoneNumber: "400-123", want: false},
		{name: "non service hyphen", phoneNumber: "123-4567890", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPhoneAll(tt.phoneNumber); got != tt.want {
				t.Errorf("CheckPhoneAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPhoneAll2(t *testing.T) {
	as := assert.New(t)
	as.False(CheckPhone("1883711386"))
	as.False(CheckLandline("1883711386"))
	as.False(CheckPhoneAll("1883711386"))
	as.True(CheckServicePhone("400-6162020"))
	as.True(CheckPhoneAll("400-6162020"))
	as.False(CheckLandline("400-6162020"))
}
