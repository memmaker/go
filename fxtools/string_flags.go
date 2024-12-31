package fxtools

import (
    "cmp"
    "github.com/memmaker/go/recfile"
    "slices"
    "strconv"
    "strings"
)

type StringFlags struct {
    underlying    map[string]int
    changeHandler func(string, int)
    flagListeners map[string][]func(int)
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
func (sf *StringFlags) AddFlagListener(key string, listener func(int)) {
    if sf.flagListeners == nil {
        sf.flagListeners = make(map[string][]func(int))
    }
    sf.flagListeners[key] = append(sf.flagListeners[key], listener)
}
func (sf *StringFlags) Increment(key string) {
    sf.Set(key, sf.Get(key)+1)
}
func (sf *StringFlags) Decrement(key string) {
    sf.Set(key, sf.Get(key)-1)
}
func (sf *StringFlags) Set(key string, val int) {
    if val == 0 {
        sf.ClearFlag(key)
        return
    }
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
    if len(rows) == 0 {
        return []string{}
    }
    return TableLayout(rows, []TextAlignment{AlignLeft, AlignRight})
}

func (sf *StringFlags) String() string {
    array := sf.ToStringArray()
    slices.SortStableFunc(array, func(i, j string) int {
        return cmp.Compare(i, j)
    })
    return strings.Join(array, "\n")
}

func (sf *StringFlags) onChange(key string, val int) {
    if sf.changeHandler != nil {
        sf.changeHandler(key, val)
    }
    if listeners, ok := sf.flagListeners[key]; ok {
        for _, listener := range listeners {
            listener(val)
        }
    }
}
func (sf *StringFlags) ToRecord() []recfile.Record {
    result := make([]recfile.Record, 0, len(sf.underlying))
    for key, val := range sf.underlying {
        if key != "" && val != 0 {
            result = append(result, recfile.Record{
                recfile.Field{Name: "Key", Value: key},
                recfile.Field{Name: "Value", Value: recfile.IntStr(val)},
            })
        }
    }
    return result
}

func NewStringFlagsFromRecord(records []recfile.Record) *StringFlags {
    sf := NewStringFlags()
    for _, record := range records {
        var key string
        var val int
        for _, field := range record {
            if field.Name == "Key" {
                key = field.Value
            } else if field.Name == "Value" {
                val = recfile.StrInt(field.Value)
            }
        }
        if key != "" && val != 0 {
            sf.Set(key, val)
        }
    }

    return sf
}
