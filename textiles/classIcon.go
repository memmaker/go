package textiles

import (
	"github.com/memmaker/go/foundation"
	"github.com/memmaker/go/recfile"
	"io"
	"strings"
)

type ClassIcon struct {
	Name string
	Icon foundation.NamedColorChar
}

func (c ClassIcon) WithForeground(rgba string) ClassIcon {
	c.Icon.Fg = rgba
	return c
}
func (c ClassIcon) WithBackground(rgba string) ClassIcon {
	c.Icon.Bg = rgba
	return c
}
func (c ClassIcon) WithIconChar(char rune) ClassIcon {
	c.Icon.Char = char
	return c
}
func (c ClassIcon) WithName(name string) ClassIcon {
	c.Name = name
	return c
}
func ReadClassIcons(reader io.Reader, palette ColorPalette) map[string]ClassIcon {
	records := recfile.Read(reader)
	categories := make(map[string]ClassIcon, len(records))
	for _, record := range records {
		category := NewClassIconFromRecord(record, palette)
		categories[strings.ToLower(category.Name)] = category
	}
	return categories
}

func NewClassIconFromRecord(record recfile.Record, palette ColorPalette) ClassIcon {
	category := ClassIcon{}
	for _, field := range record {
		switch field.Name {
		case "Name":
			category.Name = field.Value
		case "Icon":
			category.Icon.Char = []rune(field.Value)[0]
		case "Foreground":
			category.Icon.Fg = field.Value
		case "Background":
			category.Icon.Bg = field.Value
		}
	}
	return category
}

func WriteClassIcons(writer io.StringWriter, categories map[string]ClassIcon, palette ColorPalette) error {
	records := make([]recfile.Record, 0, len(categories))
	for _, category := range categories {
		records = append(records, ClassIconToRecord(category, palette))
	}
	return recfile.Write(writer, records)
}

func ClassIconToRecord(category ClassIcon, palette ColorPalette) recfile.Record {
	record := recfile.Record{}
	record = append(record, recfile.Field{Name: "Name", Value: category.Name})
	record = append(record, recfile.Field{Name: "Icon", Value: string(category.Icon.Char)})
	record = append(record, recfile.Field{Name: "Foreground", Value: category.Icon.Fg})
	record = append(record, recfile.Field{Name: "Background", Value: category.Icon.Bg})
	return record
}
