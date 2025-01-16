package util

import "strings"

// 用于混淆隐藏敏感信息的工具函数

// MaskLoginName 登录名做脱敏处理
func MaskLoginName(loginName string) string {
	// 判断 loginName 是邮箱还是手机
	if strings.Contains(loginName, "@") {
		return MaskEmail(loginName)
	}
	return MaskPhone(loginName)
}

// MaskPhone 隐去手机号中间 4 位地区码, 如 155****8888
func MaskPhone(phone string) string {
	if n := len(phone); n >= 8 {
		return phone[:n-8] + "****" + phone[n-4:]
	}
	return phone
}

// MaskEmail 隐藏邮箱 ID 的中间部分 , 如 zhang@go-mall.com ---> z***g@go-mall.com
func MaskEmail(address string) string {
	index := strings.LastIndex(address, "@")
	id := address[:index]
	domain := address[index:]

	if len(id) <= 1 {
		return address
	}
	switch len(id) {
	case 2:
		id = id[:1] + "*"
	case 3:
		id = id[:1] + "*" + id[2:]
	case 4:
		id = id[:1] + "**" + id[3:]
	default:
		masks := strings.Repeat("*", len(id)-4)
		id = id[:2] + masks + id[len(id)-2:]
	}
	return id + domain
}

// MaskRealName 隐去真实姓名中间部分, 如 张三 ---> 张* 赵丽颖--->赵*颖 欧阳娜娜--->欧**娜
func MaskRealName(realName string) string {
	runeRealName := []rune(realName)
	if n := len(runeRealName); n >= 2 {
		if n == 2 {
			return string(runeRealName[:1]) + "*"
		} else {
			count := n - 2
			newRealName := runeRealName[:1]
			for i := 0; i < count; i++ {
				newRealName = append(newRealName, '*')
			}
			return string(newRealName) + string(runeRealName[n-1:])
		}
	}
	return realName
}
