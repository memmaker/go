package recfile

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"
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

func FloatStr(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func RGBStr(value color.RGBA) string {
	return fmt.Sprintf("%d,%d,%d", value.R, value.G, value.B)
}

func TimeStr(value time.Time) string {
	return value.Format(time.RFC3339)
}

func StrTime(value string) time.Time {
	parse, _ := time.Parse(time.RFC3339, value)
	return parse
}

func StrRGB(value string) color.RGBA {
	parts := strings.Split(value, ",")
	r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}
