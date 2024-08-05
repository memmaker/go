package textiles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/memmaker/go/core"
	"image/color"
)

type TextIcon struct {
	Char       rune
	Fg         color.RGBA
	Bg         color.RGBA
	Attributes tcell.AttrMask
}

func (t TextIcon) Reversed() TextIcon {
	return TextIcon{t.Char, t.Bg, t.Fg, t.Attributes}
}
func (t TextIcon) WithFg(newFg color.RGBA) TextIcon {
	return TextIcon{t.Char, newFg, t.Bg, t.Attributes}
}

func (t TextIcon) WithBg(newBg color.RGBA) TextIcon {
	return TextIcon{t.Char, t.Fg, newBg, t.Attributes}
}

func (t TextIcon) WithColors(fgColor color.RGBA, bgColor color.RGBA) TextIcon {
	return TextIcon{t.Char, fgColor, bgColor, t.Attributes}
}

func (t TextIcon) WithRune(r rune) TextIcon {
	return TextIcon{r, t.Fg, t.Bg, t.Attributes}
}

func (t TextIcon) WithItalic() TextIcon {
	t.Attributes |= tcell.AttrItalic
	return t
}

func (t TextIcon) WithBold() TextIcon {
	t.Attributes |= tcell.AttrBold
	return t
}

func (t TextIcon) HasBackground() bool {
	return t.Bg.B == 255
}

func NewTextIconFromNamedColorChar(ncc core.NamedColorChar, palette ColorPalette) TextIcon {
	return TextIcon{
		Char: ncc.Char,
		Fg:   palette.Get(ncc.Fg),
		Bg:   palette.Get(ncc.Bg),
	}
}
