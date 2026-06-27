package fuzzutil

import (
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestFunctionalZonesToNameAliases(t *testing.T) {
	zones := map[string]string{
		"郑州高新技术产业开发区": "中原区",
		"郑州高新区":        "中原区",
		"广州经济技术开发区":   "黄埔区",
	}
	got := functionalZonesToNameAliases(zones)
	want := map[string][]string{
		"中原区": {"郑州高新技术产业开发区", "郑州高新区"},
		"黄埔区": {"广州经济技术开发区"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDefaultFunctionalZoneAliases_containsRegressionEntries(t *testing.T) {
	m := DefaultFunctionalZoneAliases()
	for zone, district := range map[string]string{
		"郑州高新技术产业开发区": "中原区",
		"广州经济技术开发区":   "黄埔区",
	} {
		if m[zone] != district {
			t.Fatalf("%q = %q, want %q", zone, m[zone], district)
		}
	}
}

func TestNationalFunctionalZones_scaleAndQuality(t *testing.T) {
	n := NationalFunctionalZoneAliasCount()
	if n < 750 || n > 1300 {
		t.Fatalf("alias count = %d, want in [750, 1300]", n)
	}
	seen := make(map[string]string)
	for alias, district := range nationalFunctionalZones {
		if district == "" {
			t.Fatalf("empty district for alias %q", alias)
		}
		if utf8.RuneCountInString(alias) < 4 {
			t.Fatalf("alias too short: %q", alias)
		}
		if prev, ok := seen[alias]; ok {
			t.Fatalf("duplicate alias %q: %q vs %q", alias, prev, district)
		}
		seen[alias] = district
	}
}

func TestMergeFunctionalZoneMap(t *testing.T) {
	dst := map[string]string{"A": "1"}
	got := mergeFunctionalZoneMap(dst, map[string]string{"B": "2", "A": "3"})
	want := map[string]string{"A": "3", "B": "2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestWithChinaAdminAddress_includesFunctionalZones(t *testing.T) {
	opts := MatchOpts().WithChinaAdminAddress()
	if !opts.stripAdminSuffixes {
		t.Fatal("expected stripAdminSuffixes")
	}
	if len(opts.functionalZones) == 0 {
		t.Fatal("expected functional zones merged")
	}
	if opts.functionalZones["郑州高新技术产业开发区"] != "中原区" {
		t.Fatalf("functional zone not loaded: %+v", opts.functionalZones["郑州高新技术产业开发区"])
	}
}
