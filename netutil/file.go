package netutil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
	"unicode"
)

// MaxDownloadSize 是 DownloadFileToLocal 的默认响应体上限。
const MaxDownloadSize = 500 * 1024 * 1024 // 500MB

// DownloadFileToLocal 下载 url 指向的文件到本地临时文件，返回文件路径。
// 使用 MaxDownloadSize（500MB）作为上限；需要更大或更小上限时用 DownloadFileToLocalLimit。
//
// 限制：仅接受 HTTP 200；响应体不超过 MaxDownloadSize。
// 安全：若 url 来自用户输入，可能被用于 SSRF（访问内网），调用方须自行校验/白名单。
func DownloadFileToLocal(rawURL string) (string, error) {
	return DownloadFileToLocalLimit(rawURL, MaxDownloadSize)
}

// DownloadFileToLocalLimit 按 maxBytes 限制响应体大小下载文件到本地临时文件。
// maxBytes 必须大于 0。
func DownloadFileToLocalLimit(rawURL string, maxBytes int64) (string, error) {
	if maxBytes <= 0 {
		return "", fmt.Errorf("maxBytes 必须大于 0")
	}

	resp, err := defaultClient.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	suffix := fileSuffixFromURL(rawURL)
	pattern := "download-*"
	if suffix != "" {
		pattern = "download-*." + suffix
	}

	out, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	closed := false
	defer func() {
		if !closed {
			_ = out.Close()
		}
	}()

	n, err := io.CopyN(out, resp.Body, maxBytes+1)
	if err != nil && err != io.EOF {
		_ = out.Close()
		closed = true
		_ = os.Remove(out.Name())
		return "", err
	}
	if n > maxBytes {
		_ = out.Close()
		closed = true
		_ = os.Remove(out.Name())
		return "", fmt.Errorf("文件超过大小限制: %d字节", maxBytes)
	}

	if err := out.Close(); err != nil {
		closed = true
		_ = os.Remove(out.Name())
		return "", err
	}
	closed = true
	return out.Name(), nil
}

// fileSuffixFromURL 从 URL 路径提取安全的文件扩展名（不含点），不可用时返回空串。
func fileSuffixFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	ext := path.Ext(u.Path)
	if ext == "" || ext == "." {
		return ""
	}
	ext = strings.TrimPrefix(ext, ".")
	if len(ext) > 16 {
		return ""
	}
	for _, r := range ext {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}
	}
	return ext
}

// CheckFileUrlCanDownload 检查 url 是否可下载，返回是否可用及错误信息。
func CheckFileUrlCanDownload(url string) (bool, string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	// 尝试 HEAD
	req, _ := http.NewRequest("HEAD", url, nil)
	req.Header.Set("User-Agent", "url-checker")
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode < 400 {
		resp.Body.Close()
		return true, ""
	}
	if resp != nil {
		resp.Body.Close()
	}

	// 尝试 GET + Range
	getReq, _ := http.NewRequest("GET", url, nil)
	getReq.Header.Set("User-Agent", "url-checker")
	getReq.Header.Set("Range", "bytes=0-0") // 只请求第1个字节
	getResp, err := client.Do(getReq)
	if err != nil {
		return false, fmt.Sprintf("GET 请求失败: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode >= 200 && getResp.StatusCode < 400 {
		b := make([]byte, 1)
		_, readErr := getResp.Body.Read(b)
		if readErr != nil && readErr != io.EOF {
			return false, fmt.Sprintf("读取失败: %v", readErr)
		}
		return true, ""
	}
	return false, fmt.Sprintf("GET 状态码: %d", getResp.StatusCode)
}
