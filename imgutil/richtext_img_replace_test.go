package imgutil

import (
	"testing"
)

var s = `
<p>酒、燕麦、芒果汁、番茄汁等系列重货产品，不是快递费真的贵，是真有点重。京东快递小哥说：成
`

func TestNewImageReplacer(t *testing.T) {

	replacer := NewImageReplacer(func(localPath string) (string, error) {
		return "abc", nil
	}, DownloadImageToTemp)
	text, err := replacer.ReplaceRichText(s)
	if err != nil {
		t.Error(err)
	}
	t.Log(text)
}
