package cview

import (
    "github.com/gdamore/tcell/v2"
    "github.com/memmaker/go/fxtools"
    "github.com/memmaker/go/geometry"
    "time"
)

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
