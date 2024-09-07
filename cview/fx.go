package cview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/memmaker/go/geometry"
	"time"
)

type StyledRune struct {
	Icon  rune
	Style tcell.Style
}

func FadeFromBlack(app *Application, animDelay time.Duration, stepSize int32, forwardBreakingKey bool) {
	screen := app.GetScreen()

	app.Lock()
	defer app.Unlock()

	var breakingKey *tcell.EventKey

	// 1. Redraw the screen but DON'T show it
	app.root.Draw(screen)

	// 2. Copy the screen contents & replace with black
	w, h := screen.Size()
	screenCopy := make([]StyledRune, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			icon, _, style, _ := screen.GetContent(x, y)
			screenCopy[y*w+x] = StyledRune{Icon: icon, Style: style}
			screen.SetContent(x, y, ' ', nil, style.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack))
		}
	}

	// 3. Lighten the screen towards the original contents

outerLoop:
	for i := 0; i < 100; i++ {
		if !restoreScreen(screen, screenCopy, w, h, stepSize) {
			break outerLoop
		}
		screen.Show()
		var waited time.Duration
		for waited < animDelay {
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				if keyEvent, ok := ev.(*tcell.EventKey); ok {
					breakingKey = keyEvent
					break outerLoop
				}
			}
			time.Sleep(10 * time.Millisecond)
			waited += 10 * time.Millisecond
		}
	}
	if breakingKey != nil && forwardBreakingKey {
		app.QueueEvent(breakingKey)
	}

}

func restoreScreen(screen tcell.Screen, screenCopy []StyledRune, w int, h int, lightenAmount int32) bool {
	centerPos := geometry.Point{X: w / 2, Y: h / 2}
	maxDist := geometry.Distance(centerPos, geometry.Point{X: 0, Y: 0})
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dist := geometry.Distance(centerPos, geometry.Point{X: x, Y: y})
			percent := Clamp(0.2, 1.0, (float64(dist)/float64(maxDist))+0.5)
			locationNeedsMoreWork := restoreScreenLocation(screen, screenCopy, x, y, int32(float64(lightenAmount)*percent), w, h)
			if locationNeedsMoreWork {
				workLeft = true
			}
		}
	}
	return workLeft
}
func restoreScreenLocation(screen tcell.Screen, screenCopy []StyledRune, x int, y int, amount int32, w int, h int) bool {
	// Get current screen content and style
	_, _, curStyle, _ := screen.GetContent(x, y)
	fg, bg, _ := curStyle.Decompose()
	currFgRed, currFgGreen, currFgBlue := fg.RGB()
	currBgRed, currBgGreen, currBgBlue := bg.RGB()

	// Calculate the position in the screen copy
	pos := y*w + x
	copyRune := screenCopy[pos]

	// Decompose the original colors from the copy
	originalStyle := copyRune.Style
	originalFG, originalBG, _ := originalStyle.Decompose()

	orgFgRed, orgFgGreen, orgFgBlue := originalFG.RGB()
	orgBgRed, orgBgGreen, orgBgBlue := originalBG.RGB()

	// Calculate the new colors

	newFgRed := min(orgFgRed, currFgRed+amount)
	newFgGreen := min(orgFgGreen, currFgGreen+amount)
	newFgBlue := min(orgFgBlue, currFgBlue+amount)

	newBgRed := min(orgBgRed, currBgRed+amount)
	newBgGreen := min(orgBgGreen, currBgGreen+amount)
	newBgBlue := min(orgBgBlue, currBgBlue+amount)

	// Determine if any work is left (i.e., if the colors are not yet fully restored)
	hadWorkLeft := currFgRed < orgFgRed || currFgGreen < orgFgGreen || currFgBlue < orgFgBlue || currBgRed < orgBgRed || currBgGreen < orgBgGreen || currBgBlue < orgBgBlue

	// Set the updated content on the screen

	drawStyle := originalStyle.
		Background(tcell.NewRGBColor(int32(newBgRed), int32(newBgGreen), int32(newBgBlue))).
		Foreground(tcell.NewRGBColor(int32(newFgRed), int32(newFgGreen), int32(newFgBlue)))

	screen.SetContent(x, y, copyRune.Icon, nil, drawStyle)

	return hadWorkLeft
}

func screenAnim(app *Application, animDelay time.Duration, stepSize int32, forwardBreakingKey bool, animator func(screen tcell.Screen, stepSize int32) (workLeft bool)) {
	screen := app.GetScreen()

	app.Lock()
	defer app.Unlock()

	var breakingKey *tcell.EventKey
outerLoop:
	for i := 0; i < 100; i++ {
		if !animator(screen, stepSize) {
			break outerLoop
		}
		screen.Show()
		var waited time.Duration
		for waited < animDelay {
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				if keyEvent, ok := ev.(*tcell.EventKey); ok {
					breakingKey = keyEvent
					break outerLoop
				}
			}
			time.Sleep(10 * time.Millisecond)
			waited += 10 * time.Millisecond
		}
	}
	if breakingKey != nil && forwardBreakingKey {
		app.QueueEvent(breakingKey)
	}
}

func FadeToBlackCircular(app *Application, animDelay time.Duration, stepSize int32, forwardBreakingKey bool) {
	screenAnim(app, animDelay, stepSize, forwardBreakingKey, darkenScreenCircular)
}

func FadeToBlack(app *Application, animDelay time.Duration, stepSize int32, forwardBreakingKey bool) {
	screenAnim(app, animDelay, stepSize, forwardBreakingKey, darkenScreen)
}

func darkenScreenCircular(screen tcell.Screen, darkenAmount int32) bool {
	w, h := screen.Size()
	centerPos := geometry.Point{X: w / 2, Y: h / 2}
	maxDist := geometry.Distance(centerPos, geometry.Point{X: 0, Y: 0})
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dist := geometry.Distance(centerPos, geometry.Point{X: x, Y: y})
			percent := Clamp(0.2, 1.0, (float64(dist)/float64(maxDist))+0.5)
			workDone := darkenScreenLocation(screen, x, y, int32(float64(darkenAmount)*percent))
			if workDone {
				workLeft = true
			}
		}
	}
	return workLeft
}

func darkenScreen(screen tcell.Screen, darkenAmount int32) bool {
	w, h := screen.Size()
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			workDone := darkenScreenLocation(screen, x, y, int32(darkenAmount))
			if workDone {
				workLeft = true
			}
		}
	}
	return workLeft
}

func darkenScreenLocation(screen tcell.Screen, x int, y int, darkenAmount int32) bool {
	icon, _, style, _ := screen.GetContent(x, y)
	fg, bg, _ := style.Decompose()
	fR, fG, fB := fg.RGB()
	bR, bG, bB := bg.RGB()
	hadWorkLeft := fR > 0 || fG > 0 || fB > 0 || bR > 0 || bG > 0 || bB > 0
	newFG := tcell.NewRGBColor(max(0, fR-darkenAmount), max(0, fG-darkenAmount), max(0, fB-darkenAmount))
	newBG := tcell.NewRGBColor(max(0, bR-darkenAmount), max(0, bG-darkenAmount), max(0, bB-darkenAmount))
	screen.SetContent(x, y, icon, nil, style.Background(newBG).Foreground(newFG))
	return hadWorkLeft
}

func FadeToWhite(app *Application, animDelay time.Duration, stepSize int, forwardBreakingKey bool) (cancel func()) {
	screen := app.GetScreen()

	cancelled := false
	app.Lock()
	defer app.Unlock()

	var breakingKey *tcell.EventKey
outerLoop:
	for i := 0; i < 100; i++ {
		if !lightenScreen(screen, stepSize) {
			break outerLoop
		}
		screen.Show()
		var waited time.Duration
		for waited < animDelay {
			if cancelled {
				return
			}
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				if keyEvent, ok := ev.(*tcell.EventKey); ok {
					breakingKey = keyEvent
					break outerLoop
				}
			}
			time.Sleep(10 * time.Millisecond)
			waited += 10 * time.Millisecond
		}
	}
	if breakingKey != nil && forwardBreakingKey {
		app.QueueEvent(breakingKey)
	}

	return func() {
		cancelled = true
	}
}

func lightenScreen(screen tcell.Screen, size int) bool {
	w, h := screen.Size()
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			workDone := lightenScreenLocation(screen, x, y, int32(size))
			if workDone {
				workLeft = true
			}
		}
	}
	return workLeft
}

func lightenScreenLocation(screen tcell.Screen, x int, y int, amount int32) bool {
	icon, _, style, _ := screen.GetContent(x, y)
	fg, bg, _ := style.Decompose()
	fR, fG, fB := fg.RGB()
	bR, bG, bB := bg.RGB()
	hadWorkLeft := fR < 255 || fG < 255 || fB < 255 || bR < 255 || bG < 255 || bB < 255
	newFG := tcell.NewRGBColor(min(255, fR+amount), min(255, fG+amount), min(255, fB+amount))
	newBG := tcell.NewRGBColor(min(255, bR+amount), min(255, bG+amount), min(255, bB+amount))
	screen.SetContent(x, y, icon, nil, style.Background(newBG).Foreground(newFG))
	return hadWorkLeft
}

func Clamp(min, max, value float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
