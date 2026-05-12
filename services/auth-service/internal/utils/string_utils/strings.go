package string_utils

import (
	"strings"
)

func FirstCharToLowerCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	firstChar := str[:1]
	return strings.ToLower(firstChar) + str[1:]
}
