package imgutil

import (
	"fmt"
	"strings"
)

func GetImg(html string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %v", err)
	}

}
