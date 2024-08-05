package textiles

import (
	"github.com/memmaker/go/recfile"
	"io"
)

type TextTile struct {
	Name          string
	Description   string
	Icon          TextIcon
	IsWalkable    bool
	IsTransparent bool
}

func ReadTilesFile(reader io.Reader, palette ColorPalette) []TextTile {
	records := recfile.Read(reader)
	tiles := make([]TextTile, len(records))
	for i, record := range records {
		tiles[i] = recordToTile(record, palette)
	}
	return tiles
}

func recordToTile(record recfile.Record, palette ColorPalette) TextTile {
	tile := TextTile{}
	for _, field := range record {
		switch field.Name {
		case "Name":
			tile.Name = field.Value
		case "Description":
			tile.Description = field.Value
		case "Char":
			tile.Icon.Char = []rune(field.Value)[0]
		case "Foreground":
			tile.Icon.Fg = palette.Get(field.Value)
		case "Background":
			tile.Icon.Bg = palette.Get(field.Value)
		case "IsWalkable":
			tile.IsWalkable = field.AsBool()
		case "IsTransparent":
			tile.IsTransparent = field.AsBool()
		}
	}
	return tile
}
