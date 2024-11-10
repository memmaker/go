package fxtools

import (
	"testing"
)

func TestSplitArgs(t *testing.T) {
	expression := "'Er Sagt\\'s so: \"Hallo\"', 123, \"Spencer's Schlüssel\",,1.00"
	want := Arguments{"hallo", "123", "Spencer"}
	result := splitArgs(expression, ',')
	for i, r := range result {
		if r != want[i] {
			t.Fatalf(`splitArgs(%q, ',') = %q, want %q`, expression, result, want)
		}
	}
}
