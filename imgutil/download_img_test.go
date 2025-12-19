package imgutil

import "testing"

func TestKaa(t *testing.T) {

	build, err := NewBuilder().Build()
	if err != nil {
		t.Error(err)
	}
	img, err := build.DownloadImg("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(img)
}
