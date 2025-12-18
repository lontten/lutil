package imgutil

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetImgSrcForHtml(html string) ([]string, error) {
	var imgSrcs []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return imgSrcs, fmt.Errorf("解析HTML失败: %v", err)
	}
	imgNodes := doc.Find("img")
	if imgNodes.Length() == 0 {
		return imgSrcs, nil
	}
	imgNodes.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			imgSrcs = append(imgSrcs, src)
		}
	})

	return imgSrcs, nil
}
