package models

import "testing"

// TestUserStateTextAndValid 用户状态: 文本描述与合法性判断
func TestUserStateTextAndValid(t *testing.T) {
	cases := map[UserState]string{
		UserStateBanned: "禁用",
		UserStateNormal: "正常",
		UserState(99):   "",
	}
	for state, want := range cases {
		if got := GetUserStateText(state); got != want {
			t.Errorf("GetUserStateText(%d) = %q, 期望 %q", state, got, want)
		}
	}
	if !IsUserStateValid(UserStateBanned) || !IsUserStateValid(UserStateNormal) {
		t.Error("合法状态应返回 true")
	}
	if IsUserStateValid(UserState(99)) {
		t.Error("非法状态应返回 false")
	}
}

// TestUserSexTextAndValid 用户性别: 文本描述与合法性判断
func TestUserSexTextAndValid(t *testing.T) {
	cases := map[UserSex]string{
		UserSexUnknown: "",
		UserSexMale:    "男",
		UserSexFemale:  "女",
		UserSex(9):     "",
	}
	for sex, want := range cases {
		if got := GetUserSexText(sex); got != want {
			t.Errorf("GetUserSexText(%d) = %q, 期望 %q", sex, got, want)
		}
	}
	// 合法性仅男/女为有效值 (未填写不算有效)
	if !IsUserSexValid(UserSexMale) || !IsUserSexValid(UserSexFemale) {
		t.Error("男/女应返回 true")
	}
	if IsUserSexValid(UserSexUnknown) || IsUserSexValid(UserSex(9)) {
		t.Error("未填写与非法值应返回 false")
	}
}

// TestUserTypeTextAndValid 用户类型: 文本描述与合法性判断
func TestUserTypeTextAndValid(t *testing.T) {
	cases := map[UserType]string{
		UserTypeOrdinary: "普通用户",
		UserType(9):      "",
	}
	for typ, want := range cases {
		if got := GetUserTypeText(typ); got != want {
			t.Errorf("GetUserTypeText(%d) = %q, 期望 %q", typ, got, want)
		}
	}
	if !IsUserTypeValid(UserTypeOrdinary) {
		t.Error("合法类型应返回 true")
	}
	if IsUserTypeValid(UserType(9)) {
		t.Error("非法类型应返回 false")
	}
}
