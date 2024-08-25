package cview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/memmaker/go/fxtools"
	"github.com/memmaker/go/geometry"
	"time"
)

type StyledRune struct {
	Icon  rune
	Style tcell.Style
}

func FadeFromBlack(app *Application, animDelay time.Duration) {
	screen := app.GetScreen()

	app.Lock()
	defer app.Unlock()

	var breakingKey *tcell.EventKey

	// create a copy of the current screen

	// 1. Redraw the screen but DON'T show it
	app.root.Draw(screen)

	// 2. Copy the screen contents & replace with black
	blackStyle := tcell.StyleDefault.Background(tcell.ColorBlack)
	w, h := screen.Size()
	screenCopy := make([]StyledRune, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			icon, _, style, _ := screen.GetContent(x, y)
			screenCopy[y*w+x] = StyledRune{Icon: icon, Style: style}
			screen.SetContent(x, y, ' ', nil, blackStyle)
		}
	}

	// 3. Lighten the screen towards the original contents

outerLoop:
	for i := 0; i < 100; i++ {
		if !lightScreen(screen, screenCopy, w, h) {
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
	if breakingKey != nil {
		app.QueueEvent(breakingKey)
	}

}

func lightScreen(screen tcell.Screen, screenCopy []StyledRune, w int, h int) bool {
	lightenAmount := int32(10)
	centerPos := geometry.Point{X: w / 2, Y: h / 2}
	maxDist := geometry.Distance(centerPos, geometry.Point{X: 0, Y: 0})
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dist := geometry.Distance(centerPos, geometry.Point{X: x, Y: y})
			percent := fxtools.Clamp(0.2, 1.0, (float64(dist)/float64(maxDist))+0.5)
			workDone := lightScreenLocation(screen, screenCopy, x, y, int32(float64(lightenAmount)*percent), w, h)
			if workDone {
				workLeft = true
			}
		}
	}
	return workLeft
}
func lightScreenLocation(screen tcell.Screen, screenCopy []StyledRune, x int, y int, amount int32, w int, h int) bool {
	// Get current screen content and style
	_, _, style, _ := screen.GetContent(x, y)
	fg, bg, _ := style.Decompose()
	currFgRed, currFgGreen, currFgBlue := fg.RGB()
	currBgRed, currBgGreen, currBgBlue := bg.RGB()

	// Calculate the position in the screen copy
	pos := y*w + x
	copyRune := screenCopy[pos]

	// Decompose the original colors from the copy
	originalFG, originalBG, _ := copyRune.Style.Decompose()

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
	screen.SetContent(x, y, copyRune.Icon, nil, style.Background(tcell.NewRGBColor(int32(newBgRed), int32(newBgGreen), int32(newBgBlue))).Foreground(tcell.NewRGBColor(int32(newFgRed), int32(newFgGreen), int32(newFgBlue))))

	return hadWorkLeft
}

func FadeToBlack(app *Application, animDelay time.Duration) {
	screen := app.GetScreen()

	app.Lock()
	defer app.Unlock()

	var breakingKey *tcell.EventKey
outerLoop:
	for i := 0; i < 100; i++ {
		if !darkenScreen(screen) {
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
	if breakingKey != nil {
		app.QueueEvent(breakingKey)
	}
}

func darkenScreen(screen tcell.Screen) bool {
	darkenAmount := int32(10)
	w, h := screen.Size()
	centerPos := geometry.Point{X: w / 2, Y: h / 2}
	maxDist := geometry.Distance(centerPos, geometry.Point{X: 0, Y: 0})
	workLeft := false
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dist := geometry.Distance(centerPos, geometry.Point{X: x, Y: y})
			percent := fxtools.Clamp(0.2, 1.0, (float64(dist)/float64(maxDist))+0.5)
			workDone := darkenScreenLocation(screen, x, y, int32(float64(darkenAmount)*percent))
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
