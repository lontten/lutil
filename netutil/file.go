package netutil

import (
	"io"
	"net/http"
	"os"
	"strings"
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
