package logutil

import (
	"encoding/json"
	"fmt"
)

func Log(v ...any) {
	for _, a := range v {
		bytes, err := json.Marshal(a)
		if err != nil {
			fmt.Println("json.Marshal err:", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
