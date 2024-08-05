package textiles

import (
	"github.com/memmaker/go/foundation"
	"image/color"
)

type TextIcon struct {
	Char       rune
	Foreground color.RGBA
	Background color.RGBA
}

func (t TextIcon) HasBackground() bool {
	return t.Background.B == 255
}

func NewTextIconFromNamedColorChar(ncc foundation.NamedColorChar, palette ColorPalette) TextIcon {
	return TextIcon{
		Char:       ncc.Char,
		Foreground: palette.Get(ncc.Fg),
		Background: palette.Get(ncc.Bg),
	}
}
