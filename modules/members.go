package modules

import (
	_type "awesomeProject1/type"
	"unicode"
)

// IsMemberParamValid 判断参数是否合法
func IsMemberParamValid(NewMember _type.Member) bool {
	// Nickname 不小于 4 位 不超过 20 位
	if len(NewMember.Nickname) < 4 || len(NewMember.Nickname) > 20 {
		return false
	}
	// Username 长度不小于 8 位 不超过 20 位
	if len(NewMember.Username) < 8 || len(NewMember.Username) > 20 {
		return false
	}
	// UserType是否为枚举值
	switch NewMember.UserType {
	case _type.Admin:
		fallthrough
	case _type.Student:
		fallthrough
	case _type.Teacher:
		break
	default:
		return false
	}
	// Username 只支持大小写
	if !isLetter(NewMember.Username) {
		return false
	}
	// Password 同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	if !isPasswordValid(NewMember.Password) {
		return false
	}
	return true
}

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isPasswordValid(pwd string) bool {
	// 长度不少于 8 位 不超过 20 位
	if len(pwd) < 8 || len(pwd) > 20 {
		return false
	}
	// 包括大小写、数字
	for _, r := range pwd {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
