package netutil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func DownloadFileToLocal(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var suffix string = ""
	split := strings.Split(url, ".")
	if len(split) > 0 {
		suffix = split[len(split)-1]
	}

	// 创建一个文件用于保存

	out, err := os.CreateTemp("", "*."+suffix)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return out.Name(), nil
}

// 检查一个文件url是否是有效链接（是否可以下载）
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
		// 可选：验证是否能读到数据
		b := make([]byte, 1)
		_, readErr := getResp.Body.Read(b)
		if readErr != nil && readErr != io.EOF {
			return false, fmt.Sprintf("读取失败: %v", readErr)
		}
		return true, ""
	}
	return false, fmt.Sprintf("GET 状态码: %d", getResp.StatusCode)
}
