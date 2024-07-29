package textiles

import (
    "fmt"
    "github.com/memmaker/go/fxtools"
    "github.com/memmaker/go/recfile"
    "image/color"
    "io"
    "os"
    "strings"
)

type NamedColor struct {
    Name  string
    Color color.RGBA
}
type ColorPalette struct {
    names  map[string]int // name -> index
    colors []NamedColor
}

func (c ColorPalette) Has(name string) bool {
    _, ok := c.names[name]
    return ok
}

func (c ColorPalette) Get(name string) color.RGBA {
    name = strings.ToLower(name)
    if foundColor, ok := c.names[name]; ok {
        return c.GetByIndex(foundColor)
    }
    return color.RGBA{}
}
func (c ColorPalette) Count() int {
    return len(c.colors)
}
func (c ColorPalette) WithColorChanged(index int, rgba color.RGBA) ColorPalette {
    newColors := make([]NamedColor, len(c.colors))
    copy(newColors, c.colors)
    newColors[index].Color = rgba
    return NewPaletteFromNamedColors(newColors)
}

func (c ColorPalette) WithColorRenamed(oldName, newName string) ColorPalette {
    newPalette := ColorPalette{make(map[string]int), make([]NamedColor, len(c.colors))}
    for name, colorIndex := range c.names {
        if name == oldName {
            newPalette.names[newName] = colorIndex
        } else {
            newPalette.names[name] = colorIndex
        }
        newPalette.colors[colorIndex] = c.colors[colorIndex]
    }
    return newPalette
}

func (c ColorPalette) GetByIndex(index int) color.RGBA {
    return c.colors[index].Color
}

func (c ColorPalette) GetNamedColorByIndex(index int) NamedColor {
    return c.colors[index]
}

func (c ColorPalette) AsNamedColors() []NamedColor {
    return c.colors
}
func (c ColorPalette) IsValidIndex(index int) bool {
    return index >= 0 && index < len(c.colors)
}

func (c ColorPalette) ToWriter(file io.StringWriter) error {
    colorRecord := recfile.Record{}
    for _, namedColor := range c.AsNamedColors() {
        colorRecord = append(colorRecord, recfile.Field{Name: namedColor.Name, Value: colorToString(namedColor.Color)})
    }
    return recfile.Write(file, []recfile.Record{colorRecord})
}

func (c ColorPalette) ToFile(filename string) error {
    outfile := fxtools.MustCreate(filename)
    defer outfile.Close()
    return c.ToWriter(outfile)
}

func colorToString(paletteColor color.RGBA) string {
    return fmt.Sprintf("%d | %d | %d", paletteColor.R, paletteColor.G, paletteColor.B)
}

func NewPaletteFromFileOrDefault(filename string) ColorPalette {
    paletteFile, openErr := os.Open(filename)
    if openErr != nil {
        defer paletteFile.Close()
    }
    return ReadPaletteFileOrDefault(paletteFile)
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

func NewPaletteFromNamedColors(colors []NamedColor) ColorPalette {
    names := make(map[string]int)
    for i, namedColor := range colors {
        names[namedColor.Name] = i
    }
    return ColorPalette{names, colors}
}

func NewDefaultPalette() ColorPalette {
    return NewPaletteFromNamedColors([]NamedColor{
        {"black", color.RGBA{0, 0, 0, 255}},
        {"white", color.RGBA{255, 255, 255, 255}},
        {"red", color.RGBA{255, 0, 0, 255}},
        {"green", color.RGBA{0, 255, 0, 255}},
        {"blue", color.RGBA{0, 0, 255, 255}},
    })
}
func recordToPalette(record recfile.Record) ColorPalette {
    names := make(map[string]int)
    colors := make([]NamedColor, len(record))
    for _, field := range record {
        colorName := strings.ToLower(field.Name) // case insensitive
        colorValue := field.AsRGB("|")
        names[colorName] = len(colors)
        colors = append(colors, NamedColor{colorName, colorValue})
    }
    return ColorPalette{names, colors}
}
