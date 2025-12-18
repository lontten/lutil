package strutil

import "testing"

func TestCheckLandline(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				phoneNumber: "028-61555395-8038",
			},
			want: true,
		},
		{
			name: "test1",
			args: args{
				phoneNumber: "028-61555395-8038",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckLandline(tt.args.phoneNumber); got != tt.want {
				t.Errorf("CheckLandline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPhoneAll(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				phoneNumber: "028-61555395-8038",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPhoneAll(tt.args.phoneNumber); got != tt.want {
				t.Errorf("CheckPhoneAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
