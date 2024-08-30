package fxtools

import (
	"math/rand"
	"regexp"
	"strconv"
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
	pattern := `(\d+)(?:\s*-\s*(\d+))?`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(s)
	var minVal, maxVal int
	if len(matches) == 3 {
		minVal = ParseInt(matches[1])
		maxVal = ParseInt(matches[2])
	} else {
		minVal = ParseInt(matches[1])
		maxVal = minVal
	}
	return NewInterval(minVal, maxVal)
}

func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
