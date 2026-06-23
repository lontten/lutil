package fileutil

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "chinese_filename",
			args: args{
				path: "/Users/go/src/github.com/utils/文件.go",
			},
			want: "文件.go",
		},
		{
			name: "ascii_filename",
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

func TestCopyTemplateToTempFile(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	templateFile, err := os.CreateTemp("", "template_*.txt")
	req.NoError(err)
	templatePath := templateFile.Name()
	defer os.Remove(templatePath)

	content := []byte("template content")
	_, err = templateFile.Write(content)
	req.NoError(err)
	req.NoError(templateFile.Close())

	tmpFile, err := CopyTemplateToTempFile(templatePath)
	req.NoError(err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.Seek(0, io.SeekStart)
	req.NoError(err)

	got, err := io.ReadAll(tmpFile)
	req.NoError(err)
	as.Equal(content, got)
}

func TestCopyTemplateToTempFileReturnPath(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	templateFile, err := os.CreateTemp("", "template_*.txt")
	req.NoError(err)
	templatePath := templateFile.Name()
	defer os.Remove(templatePath)

	content := []byte("template content")
	_, err = templateFile.Write(content)
	req.NoError(err)
	req.NoError(templateFile.Close())

	tmpPath, err := CopyTemplateToTempFileReturnPath(templatePath)
	req.NoError(err)
	defer os.Remove(tmpPath)

	got, err := os.ReadFile(tmpPath)
	req.NoError(err)
	as.Equal(content, got)
}

func TestCopyFile(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	dir := t.TempDir()
	src := dir + string(os.PathSeparator) + "src.txt"
	dst := dir + string(os.PathSeparator) + "dst.txt"
	req.NoError(os.WriteFile(src, []byte("hello"), 0644))
	req.NoError(CopyFile(src, dst))
	got, err := os.ReadFile(dst)
	req.NoError(err)
	as.Equal([]byte("hello"), got)
}

func TestNewTempFileReturnPath(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	path, err := NewTempFileReturnPath(".txt")
	req.NoError(err)
	defer os.Remove(path)
	as.FileExists(path)
}

func TestNewTempFile(t *testing.T) {
	req := require.New(t)
	f, err := NewTempFile(".log")
	req.NoError(err)
	defer os.Remove(f.Name())
	defer f.Close()
	_, err = f.WriteString("x")
	req.NoError(err)
}

func TestNewTempReturnDirName(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)

	dir, err := NewTempReturnDirName()
	req.NoError(err)
	defer os.RemoveAll(dir)
	info, err := os.Stat(dir)
	req.NoError(err)
	as.True(info.IsDir())
}
