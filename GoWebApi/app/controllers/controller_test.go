package controllers

import (
	"strings"
	"testing"
)

// TestValidateNickname 覆盖按字符计数的边界:
// 中文为多字节字符, 若用 len() 按字节统计会误判, 此测试验证按字符统计正确。
func TestValidateNickname(t *testing.T) {
	tests := []struct {
		name    string
		nickname string
		wantErr bool
	}{
		{"空昵称", "", true},
		{"1个英文字符", "a", true},
		{"1个中文字符", "好", true},
		{"2个英文字符", "ab", false},
		{"2个中文字符", "你好", false},
		{"20个中文字符", strings.Repeat("好", 20), false},
		{"20个英文字符", strings.Repeat("a", 20), false},
		{"21个中文字符", strings.Repeat("好", 21), true},
		{"21个英文字符", strings.Repeat("a", 21), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateNickname(tt.nickname); (err != nil) != tt.wantErr {
				t.Fatalf("validateNickname(%q) err = %v, wantErr = %v", tt.nickname, err, tt.wantErr)
			}
		})
	}
}
