package tools

import "strconv"

func StringsToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}
