package util

import "testing"

// TestIsPassword 密码格式: 字母、数字、_、@、$，长度 >= 6
func TestIsPassword(t *testing.T) {
	valid := []string{
		"abc123",
		"123456",
		"A_@$b9",
		"a1b2c3d4",
	}
	for _, p := range valid {
		if !IsPassword(p) {
			t.Errorf("IsPassword(%q) = false, 期望 true", p)
		}
	}

	invalid := []string{
		"",
		"a",
		"12345",    // 长度不足
		"abc def",  // 包含空格
		"a-b=c",    // 包含非法字符
		"汉字密码",     // 非 ASCII 字符
		"中文abc123", // 混合非 ASCII 字符
	}
	for _, p := range invalid {
		if IsPassword(p) {
			t.Errorf("IsPassword(%q) = true, 期望 false", p)
		}
	}
}
