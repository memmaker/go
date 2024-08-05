package core

type NamedColorChar struct {
	Char rune
	Fg   string
	Bg   string
}

func (c NamedColorChar) HasBackground() bool {
	return c.Bg != ""
}

func (c NamedColorChar) WithBackground(bg string) NamedColorChar {
	return NamedColorChar{
		Char: c.Char,
		Fg:   c.Fg,
		Bg:   bg,
	}
}

func (c NamedColorChar) IsEmpty() bool {
	return c.Char == 0 && c.Fg == "" && c.Bg == ""
}
