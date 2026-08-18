package controllers

import (
	"strings"
	"testing"
)

// TestValidateNickname 覆盖按字符计数的边界:
// 中文为多字节字符, 若用 len() 按字节统计会误判, 此测试验证按字符统计正确。
func TestValidateNickname(t *testing.T) {
	tests := []struct {
		name     string
		nickname string
		wantErr  bool
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

// TestNormalizePagination 覆盖分页参数钳制边界 (与 Rust/TS/PHP 端行为一致):
// page<=0→1, page_size<=0→12, page_size>100→100
func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name               string
		page, pageSize     int
		wantPage, wantSize int
	}{
		{"缺省参数 (0,0)", 0, 0, 1, 12},
		{"page 为负", -3, 10, 1, 10},
		{"page_size 为负", 2, -1, 2, 12},
		{"正常分页", 2, 20, 2, 20},
		{"page_size 恰好下限", 3, 1, 3, 1},
		{"page_size 恰好上限", 1, 100, 1, 100},
		{"page_size 超上限", 1, 999, 1, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize := normalizePagination(tt.page, tt.pageSize)
			if gotPage != tt.wantPage || gotPageSize != tt.wantSize {
				t.Fatalf("normalizePagination(%d, %d) = (%d, %d), want (%d, %d)",
					tt.page, tt.pageSize, gotPage, gotPageSize, tt.wantPage, tt.wantSize)
			}
		})
	}
}

// TestValidateEmailAndPassword 覆盖邮箱格式与密码规则边界:
// 密码长度 (6~64 位) 与字符集由 util.IsPassword 的正则统一维护 (见 app/util/regexp.go),
// 此测试守护"规则只在一处定义"的约定: 若有人再次在两处重复维护长度判断, 改动此处即可发现不一致。
func TestValidateEmailAndPassword(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"正常邮箱+密码", "user@example.com", "abc123", false},
		{"正常邮箱+恰好6位密码", "user@example.com", "123456", false},
		{"正常邮箱+恰好64位密码", "user@example.com", strings.Repeat("a", 64), false},
		{"邮箱格式错误", "not-an-email", "abc123", true},
		{"空邮箱", "", "abc123", true},
		{"密码不足6位", "user@example.com", "12345", true},
		{"密码超过64位", "user@example.com", strings.Repeat("a", 65), true},
		{"空密码", "user@example.com", "", true},
		{"密码含空格", "user@example.com", "abc def", true},
		{"密码含非法字符", "user@example.com", "abc-def", true},
		{"密码含中文", "user@example.com", "密码abc123", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEmailAndPassword(tt.email, tt.password); (err != nil) != tt.wantErr {
				t.Fatalf("validateEmailAndPassword(%q, %q) err = %v, wantErr = %v",
					tt.email, tt.password, err, tt.wantErr)
			}
		})
	}
}
