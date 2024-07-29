package textiles

import "image/color"

type TextIcon struct {
	Char       rune
	Foreground color.RGBA
	Background color.RGBA
}

func (t TextIcon) HasBackground() bool {
	return t.Background.B == 255
}
