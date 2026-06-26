package fuzzutil

import (
	"slices"
	"testing"
)

func TestAdminSuffixAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"回族区", "管城回族区", "管城"},
		{"达斡尔族区", "梅里斯达斡尔族区", "梅里斯"},
		{"自治县", "孟村回族自治县", "孟村回族"},
		{"旗", "正蓝旗", "正蓝"},
		{"自治旗", "鄂伦春自治旗", "鄂伦春"},
		{"盟", "锡林郭勒盟", "锡林郭勒"},
		{"林区", "神农架林区", "神农架"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adminSuffixAliases(tt.input)
			if !slices.Contains(got, tt.contains) {
				t.Fatalf("adminSuffixAliases(%q) = %v, want to contain %q", tt.input, got, tt.contains)
			}
		})
	}
}
