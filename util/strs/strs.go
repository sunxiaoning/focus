package strutil

import (
	"strconv"
	"strings"
)

func IsBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func IsValidMoney(str string) bool {
	num, err := strconv.ParseFloat(str, 64)
	return err == nil && num > 0.00
}
