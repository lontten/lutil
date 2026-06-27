package fuzzutil

import (
	"reflect"
	"testing"

	"github.com/lontten/lutil/logutil"
)

// 郑州市管城区 → 管城回族区（省略民族后缀）
func TestVocabulary_Match_AddressRegression_bug1(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"山西省", "大同市", "城区"},
		NamePath{"河南省", "郑州市", "管城回族区"},
	)
	got := vocab.Match("郑州市管城区", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"河南省", "郑州市", "管城回族区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

// 广州经济技术开发区长地址 → 黄埔区
func TestVocabulary_Match_AddressRegression_bug2(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"山西省", "临汾市", "古县"},
		NamePath{"广东省", "广州市", "黄埔区"},
	)
	got := vocab.Match("广州经济技术开发区永和经济区永顺大道西路2号", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"广东省", "广州市", "黄埔区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

// 新郑市华南城… 不误匹配南城县（全名链尾）
func TestVocabulary_Match_AddressRegression_bug3(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江西省", "抚州市", "南城县"},
		NamePath{"河南省", "郑州市", "新郑市"},
	)
	got := vocab.Match("新郑市华南城7B-1-524号", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"河南省", "郑州市", "新郑市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

// 同 bug3，词表链尾为南城（无「县」后缀）
func TestVocabulary_Match_AddressRegression_bug4(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江西省", "抚州市", "南城"},
		NamePath{"河南省", "郑州市", "新郑市"},
	)
	got := vocab.Match("新郑市华南城7B-1-524号", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"河南省", "郑州市", "新郑市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

// 新郑华南城…（无「市」）不误匹配南城
func TestVocabulary_Match_AddressRegression_bug5(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江西省", "抚州市", "南城"},
		NamePath{"河南省", "郑州市", "新郑"},
	)
	got := vocab.Match("新郑华南城7B-1-524号", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"河南省", "郑州市", "新郑"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestVocabulary_Match_AddressRegression_bug6(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"中国"},
		NamePath{"中国", "四川"},
		NamePath{"中国", "四川", "服饰"},
	)
	got := vocab.Match("中国四", MatchOpts().WithChinaAdminAddress())
	logutil.Log(got)
}
func TestVocabulary_Match_AddressRegression_bug7(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江西省", "萍乡市", "莲花县"},
		NamePath{"河南省", "郑州市", "中原区"},
	)
	got := vocab.Match("郑州高新技术产业开发区莲花街316号6号楼608室", MatchOpts().WithChinaAdminAddress())
	want := NamePath{"河南省", "郑州市", "中原区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}
