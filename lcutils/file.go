package lcutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewTempFileReturnPath 生成唯一临时文件名（返回路径和错误）
// suffix 可以是 .xxx 或者 xxx.xxx
// 建议有后缀名
func NewTempFileReturnPath(suffix string) (string, error) {
	// 在系统临时目录创建文件，文件名格式：temp_随机数+suffix
	file, err := os.CreateTemp("", "temp_*"+suffix)
	if err != nil {
		return "", err
	}
	path := file.Name()

	// 立即关闭文件（由调用方决定后续操作）
	if err := file.Close(); err != nil {
		os.Remove(path) // 关闭失败时清理文件
		return "", err
	}
	return path, nil
}

// NewTempFileReturnPathMust 强制生成临时文件名（失败时panic）
// suffix 可以是 .xxx 或者 xxx.xxx
// 建议有后缀名
func NewTempFileReturnPathMust(suffix string) string {
	path, err := NewTempFileReturnPath(suffix)
	if err != nil {
		panic(err)
	}
	return path
}

// NewTempFile 生成唯一临时文件（返回文件对象和错误）
// suffix 可以是 .xxx 或者 xxx.xxx
// 建议有后缀名
func NewTempFile(suffix string) (*os.File, error) {
	// 直接返回已打开的文件句柄
	return os.CreateTemp("", "temp_*"+suffix)
}

// NewTempFileMust 强制生成临时文件（失败时panic）
// suffix 可以是 .xxx 或者 xxx.xxx
// 建议有后缀名
func NewTempFileMust(suffix string) *os.File {
	file, err := NewTempFile(suffix)
	if err != nil {
		panic(err)
	}
	return file
}

// NewTempReturnDirName 生成唯一临时目录名（返回目录路径和错误）
func NewTempReturnDirName() (string, error) {
	dir, err := os.MkdirTemp("", "tempdir_")
	if err != nil {
		return "", err
	}
	return dir, nil
}

// NewTempReturnDirNameMust 生成唯一临时目录名（必须成功，失败时 panic）
func NewTempReturnDirNameMust() string {
	dir, err := NewTempReturnDirName()
	if err != nil {
		panic(err) // 若创建失败则触发 panic
	}
	return dir
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CopyTemplateToTempFileReturnPath 复制模板文件到临时文件，返回临时文件路径
func CopyTemplateToTempFileReturnPath(templatePath string) (string, error) {
	tmpFile, err := CopyTemplateToTempFile(templatePath)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// CopyTemplateToTempFile 复制模板文件到临时文件，返回临时文件路径
func CopyTemplateToTempFile(templatePath string) (tmpFile *os.File, err error) {
	// 读取模板文件内容
	content, err := os.ReadFile(templatePath)
	if err != nil {
		err = fmt.Errorf("读取模板文件失败: %w", err)
		return
	}

	// 获取原文件权限
	fileInfo, err := os.Stat(templatePath)
	if err != nil {
		err = fmt.Errorf("获取文件信息失败: %w", err)
		return
	}
	perm := fileInfo.Mode().Perm()

	// 构造临时文件名模式（保留扩展名）
	base := filepath.Base(templatePath)
	ext := filepath.Ext(base)
	pattern := strings.TrimSuffix(base, ext) + "*" + ext // 格式：原文件名*扩展名

	// 创建临时文件
	tmpFile, err = os.CreateTemp("", pattern)
	if err != nil {
		err = fmt.Errorf("创建临时文件失败: %w", err)
		return
	}
	defer tmpFile.Close() // 确保文件关闭

	// 写入内容
	if _, err = tmpFile.Write(content); err != nil {
		err = fmt.Errorf("写入临时文件失败: %w", err)
		return
	}

	// 设置权限（与原文件一致）
	if err = os.Chmod(tmpFile.Name(), perm); err != nil {
		err = fmt.Errorf("设置权限失败: %w", err)
		return
	}
	return
}

// GetFileName 获取文件名(包含扩展名)
func GetFileName(path string) string {
	return filepath.Base(path)
}

// GetFileNameNoSuffix 获取文件名(不包含扩展名)
func GetFileNameNoSuffix(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// GetFileSuffix 获取文件后缀(不带点)
func GetFileSuffix(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
