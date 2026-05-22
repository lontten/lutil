package logutil

import (
	"encoding/json"
	"fmt"
)

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
