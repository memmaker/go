package recfile

import (
	"strconv"
	"strings"
)

func IntStr(value int) string {
	return strconv.Itoa(value)
}
func Int64Str(value int64) string {
	return strconv.FormatInt(value, 10)
}

func StrInt(value string) int {
	atoi, _ := strconv.Atoi(value)
	return atoi
}
func Int32Str(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}
func BoolStr(value bool) string {
	return strconv.FormatBool(value)
}
func StringsStr(value []string) string {
	return strings.Join(value, "\n")
}
func StrBool(value string) bool {
	parseBool, _ := strconv.ParseBool(value)
	return parseBool
}
