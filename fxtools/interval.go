package fxtools

import (
	"math/rand"
	"strconv"
	"strings"
)

type Interval struct {
	Min int
	Max int
}

func (i Interval) NotZero() bool {
	return !i.IsZero()
}

func (i Interval) IsZero() bool {
	return i.Min == 0 && i.Max == 0
}

func (i Interval) ExpectedValue() int {
	return (i.Min + i.Max) / 2
}
func (i Interval) Scaled(scale float64) Interval {
	return NewInterval(int(float64(i.Min)*scale), int(float64(i.Max)*scale))
}

func (i Interval) ShortString() string {
	if i.Min == i.Max {
		return strconv.Itoa(i.Min)
	}
	return strconv.Itoa(i.Min) + "-" + strconv.Itoa(i.Max)
}

func (i Interval) Roll() int {
	if i.Min == i.Max {
		return i.Min
	}
	return rand.Intn(i.Max-i.Min+1) + i.Min
}

func NewInterval(min, max int) Interval {
	if min > max {
		min, max = max, min
	}
	return Interval{min, max}
}

func ParseInterval(s string) Interval {
	// looks like "1-6" or possibly "1", "2 - 6"
	if !strings.Contains(s, "-") {
		s = strings.TrimSpace(strings.ReplaceAll(s, " ", ""))
		return NewInterval(ParseInt(s), ParseInt(s))
	}
	parts := strings.Split(s, "-")
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])
	if left == "" {
		// a negative number
		realVal := ParseInt(right) * -1
		return NewInterval(realVal, realVal)
	}
	return NewInterval(ParseInt(left), ParseInt(right))
}

func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
