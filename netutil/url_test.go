package netutil

import "testing"

func TestSanitizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				url: "0.2+0.3.png",
			},
			want: "0.2_0.3.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SafeURL(tt.args.url); got != tt.want {
				t.Errorf("SafeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
