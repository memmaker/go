package fxtools

import (
	"github.com/memmaker/go/recfile"
	"strconv"
	"strings"
)

type StringFlags struct {
	underlying    map[string]int
	changeHandler func(string, int)
}

func NewStringFlags() *StringFlags {
	return &StringFlags{underlying: make(map[string]int)}
}
func (sf *StringFlags) Get(key string) int {
	if val, ok := sf.underlying[key]; ok {
		return val
	}
	return 0
}
func (sf *StringFlags) SetChangeHandler(handler func(string, int)) {
	sf.changeHandler = handler
}
func (sf *StringFlags) Set(key string, val int) {
	sf.underlying[key] = val
	sf.onChange(key, val)
}

func (sf *StringFlags) HasFlag(key string) bool {
	return sf.Get(key) != 0
}

func (sf *StringFlags) SetFlag(key string) {
	sf.Set(key, 1)
	sf.onChange(key, 1)
}

func (sf *StringFlags) ClearFlag(key string) {
	delete(sf.underlying, key)
	sf.onChange(key, 0)
}

func (sf *StringFlags) ToStringArray() []string {
	if len(sf.underlying) == 0 {
		return []string{}
	}
	var rows []TableRow
	for key, val := range sf.underlying {
		if val != 0 {
			rows = append(rows, TableRow{Columns: []string{key, strconv.Itoa(val)}})
		}
	}
	return TableLayout(rows, []TextAlignment{AlignLeft, AlignRight})
}

func (sf *StringFlags) String() string {
	return strings.Join(sf.ToStringArray(), "\n")
}

func (sf *StringFlags) onChange(key string, val int) {
	if sf.changeHandler != nil {
		sf.changeHandler(key, val)
	}
}
func (sf *StringFlags) ToRecord() recfile.Record {
	record := recfile.Record{}
	for key, val := range sf.underlying {
		record = append(record, recfile.Field{Name: key, Value: recfile.IntStr(val)})
	}
	return record
}

func NewStringFlagsFromRecord(record recfile.Record) *StringFlags {
	sf := NewStringFlags()
	for _, field := range record {
		sf.Set(field.Name, recfile.StrInt(field.Value))
	}
	return sf
}
