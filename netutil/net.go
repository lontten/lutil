package netutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Get[T any](targetURL string) (T, error) {
	statusCode, result, err := GetBase[T](targetURL)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func GetBase[T any](targetURL string) (int, T, error) {
	var result T
	resp, err := http.Get(targetURL)
	if err != nil {
		return 0, result, fmt.Errorf("发送 GET 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, result, fmt.Errorf("读取响应体失败: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return resp.StatusCode, result, nil
}
func PostJsonOk[T any](targetURL string, data any) (T, error) {
	statusCode, result, err := PostJson[T](targetURL, data)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func PostJson[T any](targetURL string, data any) (int, T, error) {
	var result T
	code, body, err := PostJsonNative(targetURL, data)
	if err != nil {
		return 0, result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return code, result, nil
}

func PostJsonNative(targetURL string, data any) (int, []byte, error) {
	var jsonBody []byte
	var err error

	if data == nil {
		jsonBody = []byte("{}")
	} else {
		jsonBody, err = json.Marshal(data)
		if err != nil {
			return 0, []byte{}, fmt.Errorf("请求体 JSON 序列化失败: %w", err)
		}
	}
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, []byte{}, fmt.Errorf("发送 PostJson 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, []byte{}, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte{}, fmt.Errorf("读取响应体失败: %w", err)
	}
	return resp.StatusCode, body, nil
}
func PostJsonByte[T any](targetURL string, data any) (int, []byte, T, error) {
	var result T
	var jsonBody []byte
	var err error

	if data == nil {
		jsonBody = []byte("{}")
	} else {
		jsonBody, err = json.Marshal(data)
		if err != nil {
			return 0, []byte{}, result, fmt.Errorf("请求体 JSON 序列化失败: %w", err)
		}
	}
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, []byte{}, result, fmt.Errorf("发送 PostJson 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, []byte{}, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte{}, result, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 5. 判断响应类型（核心逻辑）
	contentType := resp.Header.Get("Content-Type")
	// 情况1：返回JSON错误（微信接口失败时Content-Type为application/json）
	if strings.Contains(contentType, "application/json") {
		if err := json.Unmarshal(body, &result); err != nil {
			return 0, []byte{}, result, fmt.Errorf("解析微信错误响应失败: %v, 原始数据: %s", err, string(body))
		}
		return resp.StatusCode, nil, result, nil
	}

	return resp.StatusCode, body, result, nil
}

func PostJsonByteOk[T any](targetURL string, data any) ([]byte, T, error) {
	statusCode, result, t, err := PostJsonByte[T](targetURL, data)
	if err != nil {
		return result, t, err
	}
	if statusCode != http.StatusOK {
		return result, t, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, t, nil
}

func PostJsonNativeOk(targetURL string, data any) ([]byte, error) {
	statusCode, result, err := PostJsonNative(targetURL, data)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func PostFormOk[T any](targetURL string, data url.Values) (T, error) {
	statusCode, result, err := PostForm[T](targetURL, data)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func PostForm[T any](targetURL string, data url.Values) (int, T, error) {
	var result T

	if data == nil {
		data = url.Values{}
	}

	resp, err := http.PostForm(targetURL, data)
	if err != nil {
		return 0, result, fmt.Errorf("发送 PostForm 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, result, fmt.Errorf("读取响应体失败: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return resp.StatusCode, result, nil
}
