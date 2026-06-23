// Package logutil 提供简单的日志输出工具。
package logutil

import (
	"encoding/json"
	"fmt"
)

// Log 将参数以 JSON 缩进格式打印到标准输出。
// 若某个参数无法序列化为 JSON，打印错误信息并停止后续输出。
func Log(v ...any) {
	for _, a := range v {
		bytes, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			fmt.Println("json.MarshalIndent err:", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
