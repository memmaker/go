package textiles

import (
	"github.com/memmaker/go/core"
	"github.com/memmaker/go/geometry"
	"github.com/memmaker/go/recfile"
	"io"
	"strings"
)

type IconRecord struct {
	Name string
	Icon core.NamedColorChar
	Meta recfile.Record
}

func (c IconRecord) WithIcon(icon core.NamedColorChar) IconRecord {
	c.Icon = icon
	return c
}
func (c IconRecord) WithForeground(rgba string) IconRecord {
	c.Icon.Fg = rgba
	return c
}
func (c IconRecord) WithBackground(rgba string) IconRecord {
	c.Icon.Bg = rgba
	return c
}
func (c IconRecord) WithIconChar(char rune) IconRecord {
	c.Icon.Char = char
	return c
}
func (c IconRecord) WithName(name string) IconRecord {
	c.Name = name
	return c
}
func (c IconRecord) ToRecordWithPosition(mapPos geometry.Point) recfile.Record {
	record := c.ToRecord()
	return record.WithKeyValue("Position", mapPos.Encode())
}

func (c IconRecord) String() string {
	return c.ToRecord().String()
}
func ReadIconRecordsIntoMap(reader io.Reader) map[string]IconRecord {
	records := recfile.Read(reader)
	categories := make(map[string]IconRecord, len(records))
	for _, record := range records {
		category := NewIconRecord(record)
		categories[category.Name] = category
	}
	return categories
}

func ReadIconRecords(reader io.Reader) []IconRecord {
	records := recfile.Read(reader)
	categories := make([]IconRecord, 0, len(records))
	for _, record := range records {
		category := NewIconRecord(record)
		categories = append(categories, category)
	}
	return categories
}

func WriteIconRecordMap(writer io.StringWriter, categories map[string]IconRecord) error {
	records := make([]recfile.Record, 0, len(categories))
	for _, category := range categories {
		records = append(records, category.ToRecord())
	}
	return recfile.Write(writer, records)
}

func WriteIconRecords(writer io.StringWriter, categories []IconRecord) error {
	records := make([]recfile.Record, 0, len(categories))
	for _, category := range categories {
		records = append(records, category.ToRecord())
	}
	return recfile.Write(writer, records)
}

func (c IconRecord) ToRecord() recfile.Record {
	record := recfile.Record{}
	record = append(record, recfile.Field{Name: "Name", Value: c.Name})
	record = append(record, recfile.Field{Name: "Icon", Value: string(c.Icon.Char)})
	record = append(record, recfile.Field{Name: "Foreground", Value: c.Icon.Fg})
	if c.Icon.HasBackground() {
		record = append(record, recfile.Field{Name: "Background", Value: c.Icon.Bg})
	}
	for _, field := range c.Meta {
		record = append(record, field)
	}
	return record
}

func NewIconRecord(record recfile.Record) IconRecord {
	category := IconRecord{}
	for _, field := range record {
		switch strings.ToLower(field.Name) {
		case "name":
			category.Name = field.Value
		case "icon":
			category.Icon.Char = []rune(field.Value)[0]
		case "foreground":
			category.Icon.Fg = field.Value
		case "background":
			category.Icon.Bg = field.Value
		default:
			category.Meta = append(category.Meta, field)
		}
	}
	return category
}
