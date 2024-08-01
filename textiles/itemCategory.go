package textiles

import (
	"github.com/memmaker/go/recfile"
	"image/color"
	"io"
	"strconv"
	"strings"
)

type ItemCategory struct {
	Name string
	Icon TextIcon
}

func (c ItemCategory) WithForeground(rgba color.RGBA) ItemCategory {
	c.Icon.Foreground = rgba
	return c
}
func (c ItemCategory) WithBackground(rgba color.RGBA) ItemCategory {
	c.Icon.Background = rgba
	return c
}
func (c ItemCategory) WithIconChar(char rune) ItemCategory {
	c.Icon.Char = char
	return c
}
func (c ItemCategory) WithName(name string) ItemCategory {
	c.Name = name
	return c
}
func ReadItemCategoriesFile(reader io.Reader, palette ColorPalette) map[string]ItemCategory {
	records := recfile.Read(reader)
	categories := make(map[string]ItemCategory, len(records))
	for _, record := range records {
		category := recordToCategory(record, palette)
		categories[strings.ToLower(category.Name)] = category
	}
	return categories
}

func recordToCategory(record recfile.Record, palette ColorPalette) ItemCategory {
	category := ItemCategory{}
	for _, field := range record {
		switch field.Name {
		case "Name":
			category.Name = field.Value
		case "Icon":
			category.Icon.Char = []rune(field.Value)[0]
		case "Foreground":
			category.Icon.Foreground = palette.Get(field.Value)
		case "Background":
			category.Icon.Background = palette.Get(field.Value)
		}
	}
	return category
}

func WriteItemCategoriesFile(writer io.StringWriter, categories map[string]ItemCategory, palette ColorPalette) error {
	records := make([]recfile.Record, 0, len(categories))
	for _, category := range categories {
		records = append(records, categoryToRecord(category, palette))
	}
	return recfile.Write(writer, records)
}

func categoryToRecord(category ItemCategory, palette ColorPalette) recfile.Record {
	record := recfile.Record{}
	record = append(record, recfile.Field{Name: "Name", Value: category.Name})
	record = append(record, recfile.Field{Name: "Icon", Value: string(category.Icon.Char)})
	record = append(record, recfile.Field{Name: "Foreground", Value: strconv.Itoa(palette.GetIndexOfColor(category.Icon.Foreground))})
	record = append(record, recfile.Field{Name: "Background", Value: strconv.Itoa(palette.GetIndexOfColor(category.Icon.Background))})
	return record
}
