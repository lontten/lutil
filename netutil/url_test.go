package netutil

import "testing"

func TestSafeFileName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "percent and plus replaced",
			input:  "%0.2+0.3.png",
			maxLen: 0,
			want:   "_0.2_0.3.png",
		},
		{
			name:   "chinese with percent",
			input:  "酒精%.png",
			maxLen: 0,
			want:   "酒精_.png",
		},
		{
			name:   "reserved hash amp eq",
			input:  "a#b&c=d.txt",
			maxLen: 0,
			want:   "a_b_c_d.txt",
		},
		{
			name:   "basename from path",
			input:  "a/b/c.png",
			maxLen: 0,
			want:   "c.png",
		},
		{
			name:   "path traversal",
			input:  "../x.txt",
			maxLen: 0,
			want:   "x.txt",
		},
		{
			name:   "windows illegal and space",
			input:  "my file?.png",
			maxLen: 0,
			want:   "my_file_.png",
		},
		{
			name:   "illegal chars",
			input:  `name<>:"|*.txt`,
			maxLen: 0,
			want:   "name______.txt",
		},
		{
			name:   "truncate keeps extension",
			input:  "verylongfilename.png",
			maxLen: 10,
			want:   "verylo.png",
		},
		{
			name:   "empty becomes file",
			input:  "",
			maxLen: 0,
			want:   "file",
		},
		{
			name:   "dot becomes file",
			input:  ".",
			maxLen: 0,
			want:   "file",
		},
		{
			name:   "dotdot becomes file",
			input:  "..",
			maxLen: 0,
			want:   "file",
		},
		{
			name:   "whitespace only becomes file",
			input:  "   ",
			maxLen: 0,
			want:   "file",
		},
		{
			name:   "no truncate when maxLen 0",
			input:  "verylongfilename.png",
			maxLen: 0,
			want:   "verylongfilename.png",
		},
		{
			name:   "control char replaced",
			input:  "a\x00b.txt",
			maxLen: 0,
			want:   "a_b.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SafeFileName(tt.input, tt.maxLen); got != tt.want {
				t.Errorf("SafeFileName(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestSafeURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		size []int
		want string
	}{
		{
			name: "percent and plus replaced",
			url:  "%0.2+0.3.png",
			want: "_0.2_0.3.png",
		},
		{
			name: "chinese with percent",
			url:  "酒精%.png",
			want: "酒精_.png",
		},
		{
			name: "omit size no truncate",
			url:  "verylongfilename.png",
			want: "verylongfilename.png",
		},
		{
			name: "size 0 no truncate",
			url:  "verylongfilename.png",
			size: []int{0},
			want: "verylongfilename.png",
		},
		{
			name: "size truncates with ext",
			url:  "verylongfilename.png",
			size: []int{10},
			want: "verylo.png",
		},
		{
			name: "basename",
			url:  "a/b/c.png",
			want: "c.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			if len(tt.size) > 0 {
				got = SafeURL(tt.url, tt.size...)
			} else {
				got = SafeURL(tt.url)
			}
			if got != tt.want {
				t.Errorf("SafeURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
