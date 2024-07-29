package textiles

import (
    "fmt"
    "github.com/memmaker/go/recfile"
    "image/color"
    "io"
    "strings"
)

type ColorPalette map[string]color.RGBA

func (c ColorPalette) Has(name string) bool {
    _, ok := c[name]
    return ok
}

func (c ColorPalette) Get(name string) color.RGBA {
    name = strings.ToLower(name)
    if foundColor, ok := c[name]; ok {
        return foundColor
    }
    return color.RGBA{}
}

func (c ColorPalette) WithColorRenamed(oldName, newName string) ColorPalette {
    newPalette := make(ColorPalette)
    for name, paletteColor := range c {
        if name == oldName {
            newPalette[newName] = paletteColor
        } else {
            newPalette[name] = paletteColor
        }
    }
    return newPalette
}

func WritePaletteFile(file io.StringWriter, palette ColorPalette) error {
    colorRecord := recfile.Record{}
    for name, paletteColor := range palette {
        colorRecord = append(colorRecord, recfile.Field{Name: name, Value: colorToString(paletteColor)})
    }
    return recfile.Write(file, []recfile.Record{colorRecord})
}

func colorToString(paletteColor color.RGBA) string {
    return fmt.Sprintf("%d | %d | %d", paletteColor.R, paletteColor.G, paletteColor.B)
}

func ReadPaletteFile(file io.Reader) ColorPalette {
    records := recfile.Read(file)
    return recordToPalette(records[0])
}
func ReadPaletteFileOrDefault(file io.Reader) ColorPalette {
    records := recfile.Read(file)
    if len(records) == 0 {
        return NewDefaultPalette()
    }
    return recordToPalette(records[0])
}

func NewDefaultPalette() ColorPalette {
    return ColorPalette{
        "black":   color.RGBA{0, 0, 0, 255},
        "red":     color.RGBA{255, 0, 0, 255},
        "green":   color.RGBA{0, 255, 0, 255},
        "yellow":  color.RGBA{255, 255, 0, 255},
        "blue":    color.RGBA{0, 0, 255, 255},
        "magenta": color.RGBA{255, 0, 255, 255},
        "cyan":    color.RGBA{0, 255, 255, 255},
        "white":   color.RGBA{255, 255, 255, 255},
    }
}
func recordToPalette(record recfile.Record) ColorPalette {
    colors := make(map[string]color.RGBA)
    for _, field := range record {
        colorName := strings.ToLower(field.Name) // case insensitive
        colorValue := field.AsRGB("|")
        colors[colorName] = colorValue
    }
    return colors
}
