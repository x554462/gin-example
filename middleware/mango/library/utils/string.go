package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GetRandomString(l int, typo string) string {
	var str = ""
	if typo == "num" {
		str += "0123456789"
	} else if typo == "char" {
		str += "abcdefghijklmnopqrstuvwxyz"
	} else {
		str += "0123456789abcdefghijklmnopqrstuvwxyz"
	}
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func TitleCasedName(name string) string {
	newStr := make([]rune, 0)
	upNextChar := true
	name = strings.ToLower(name)
	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}
		newStr = append(newStr, chr)
	}
	return string(newStr)
}
