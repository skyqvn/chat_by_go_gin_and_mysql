package main

import "fmt"

// 验证密码是否合法的函数
func IsAValidChatGroupPassword(password string) error {
	charList := [...]byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'k', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'K', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_', '-'}
	if len(password) > 20 {
		return fmt.Errorf("密码长度不能超过二十个字符")
	}
	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于八个字符")
	}

	// 判断密码的每一个字符是否都在charList中
	// 如果检测到任何一个字符不在当中，则返回错误
	for i := 0; i < len(password); i++ {
		ok := false
		for j := 0; j < len(charList); j++ {
			if password[i] == charList[j] {
				ok = true
				break
			}
		}
		if ok == false {
			return fmt.Errorf("密码必须是由英文字母大小写、数字、\"_\"、\"-\"组成的不多于二十个字符、不少于八个字符的字符串")
		}
	}
	return nil
}
