package lcutils

import (
	"testing"
)

func TestGetFileName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				path: "/Users/go/src/github.com/utils/文件.go",
			},
			want: "文件.go",
		},
		{
			name: "test1",
			args: args{
				path: "/Users/go/src/github.com/utils/file.go",
			},
			want: "file.go",
		},
		{
			name: "test2",
			args: args{
				path: "/Users/go/src/github.com/utils/file.",
			},
			want: "file.",
		},
		{
			name: "test3",
			args: args{
				path: "/Users/go/src/github.com/utils/file",
			},
			want: "file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileName(tt.args.path); got != tt.want {
				t.Errorf("GetFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileNameNoSuffix(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				path: "/Users/go/src/github.com/utils/file.go",
			},
			want: "file",
		},
		{
			name: "test2",
			args: args{
				path: "/Users/go/src/github.com/utils/file.",
			},
			want: "file",
		},
		{
			name: "test3",
			args: args{
				path: "/Users/go/src/github.com/utils/file",
			},
			want: "file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileNameNoSuffix(tt.args.path); got != tt.want {
				t.Errorf("GetFileNameNoSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileSuffix(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				path: "/Users/go/src/github.com/utils/file.go",
			},
			want: "go",
		},
		{
			name: "test2",
			args: args{
				path: "/Users/go/src/github.com/utils/file.",
			},
			want: "",
		},
		{
			name: "test3",
			args: args{
				path: "/Users/go/src/github.com/utils/file",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileSuffix(tt.args.path); got != tt.want {
				t.Errorf("GetFileSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
