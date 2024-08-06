package textiles

import (
    "encoding/binary"
    "github.com/memmaker/go/geometry"
    "github.com/memmaker/go/recfile"
    "io"
    "os"
)

type TextTile struct {
    Name          string
    Description   string
    Icon          TextIcon
    IsWalkable    bool
    IsTransparent bool
}

func (t TextTile) WithIcon(icon TextIcon) TextTile {
    t.Icon = icon
    return t
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
    var icon TextIcon
    for _, field := range record {
        switch field.Name {
        case "Name":
            tile.Name = field.Value
        case "Description":
            tile.Description = field.Value
        case "Char":
            icon.Char = []rune(field.Value)[0]
        case "Foreground":
            icon.Fg = palette.Get(field.Value)
        case "Background":
            icon.Bg = palette.Get(field.Value)
        case "IsWalkable":
            tile.IsWalkable = field.AsBool()
        case "IsTransparent":
            tile.IsTransparent = field.AsBool()
        }
    }
    return tile.WithIcon(icon)
}

func SaveTileMap16(tiles []int16, dimension geometry.Point, filename string) error {
    file, openErr := os.Create(filename)
    if openErr != nil {
        return openErr
    }
    defer file.Close()
    handleErr := func(err error) {
        if err != nil {
            panic(err)
        }
    }

    handleErr(binary.Write(file, binary.LittleEndian, int32(dimension.X)))
    handleErr(binary.Write(file, binary.LittleEndian, int32(dimension.Y)))
    for _, tile := range tiles {
        handleErr(binary.Write(file, binary.LittleEndian, tile))
    }
    return nil
}

func ReadTileMap16(filename string) (geometry.Point, []int16) {
    file, openErr := os.Open(filename)
    if openErr != nil {
        return geometry.Point{}, nil
    }
    defer file.Close()
    var dimensionX, dimensionY int32

    handleErr := func(err error) {
        if err != nil {
            panic(err)
        }
    }

    handleErr(binary.Read(file, binary.LittleEndian, &dimensionX))
    handleErr(binary.Read(file, binary.LittleEndian, &dimensionY))
    dimension := geometry.Point{X: int(dimensionX), Y: int(dimensionY)}

    tiles := make([]int16, dimension.X*dimension.Y)
    for i := 0; i < len(tiles); i++ {
        handleErr(binary.Read(file, binary.LittleEndian, &tiles[i]))
    }
    return dimension, tiles
}
