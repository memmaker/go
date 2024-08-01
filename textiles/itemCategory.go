package textiles

import (
    "github.com/memmaker/go/recfile"
    "io"
)

type ItemCategory struct {
    Name string
    Icon TextIcon
}

func ReadItemCategoriesFile(reader io.Reader, palette ColorPalette) []ItemCategory {
    records := recfile.Read(reader)
    categories := make([]ItemCategory, len(records))
    for i, record := range records {
        categories[i] = recordToCategory(record, palette)
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
