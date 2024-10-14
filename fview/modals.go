package fview

import (
    "github.com/gdamore/tcell/v2"
    "github.com/memmaker/go/cview"
)

type InputPrimitive interface {
    cview.Primitive
    SetRect(x, y, w, h int)
    GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey
    SetInputCapture(f func(event *tcell.EventKey) *tcell.EventKey)
}

type MenuItem struct {
    Name       string
    Action     func()
    CloseMenus bool
}

func OpenSimpleMenu(app *cview.Application, pages *cview.Panels, menuItems []MenuItem) *cview.List {
    list := cview.NewList()
    list.ShowSecondaryText(false)
    list.SetBorder(true)

    //u.applyListStyle(list)
    list.SetSelectedFunc(func(index int, listItem *cview.ListItem) {
        action := menuItems[index]
        if action.CloseMenus {
            closeModal(app, pages)
        }
        action.Action()
    })

    longestItem := setListItemsFromMenuItems(list, menuItems)
    makeCenteredModal(app, pages, list, longestItem, len(menuItems))
    return list
}
func closeModal(app *cview.Application, pages *cview.Panels) {
    pages.RemovePanel("modal")
    resetFocusToMain(app, pages)
}
func makeCenteredModal(app *cview.Application, pages *cview.Panels, modal InputPrimitive, w, h int) {
    w = w + 2
    h = h + 2
    screenW, screenH := app.GetScreen().Size()
    if h > screenH {
        h = screenH
        w = w + 1 // scrollbar
    }
    if w > screenW {
        w = screenW
    }
    x, y := (screenW-w)/2, (screenH-h)/2
    modal.SetRect(x, y, w, h)

    originalInputCapture := modal.GetInputCapture()
    modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyEscape {
            closeModal(app, pages)
        }
        if originalInputCapture != nil {
            return originalInputCapture(event)
        }
        return event
    })

    pages.AddPanel("modal", modal, false, true)
    lockFocusToPrimitive(app, modal)
}

func resetFocusToMain(app *cview.Application, pages *cview.Panels) {
    app.SetBeforeFocusFunc(nil)
    _, frontPanel := pages.GetFrontPanel()
    app.SetFocus(frontPanel)
    //app.SetBeforeFocusFunc(u.defaultFocusHandler)
}
func lockFocusToPrimitive(app *cview.Application, p cview.Primitive) {
    app.SetBeforeFocusFunc(nil)
    app.SetFocus(p)
    app.SetBeforeFocusFunc(func(p cview.Primitive) bool { return false })
}
func setListItemsFromMenuItems(list *cview.List, menuItems []MenuItem) int {
    list.Clear()
    longestItem := 0
    for index, a := range menuItems {
        action := a
        shortcut := ShortCutFromIndex(index)
        listItem := cview.NewListItem(action.Name)
        listItem.SetShortcut(shortcut)
        list.AddItem(listItem)
        itemLength := cview.TaggedStringWidth(action.Name) + 4
        longestItem = max(longestItem, itemLength)
    }
    return longestItem
}

func ShortCutFromIndex(itemIndex int) rune {
    shortcut := rune(97 + itemIndex)
    if itemIndex > 25 {
        shortcut = rune(39 + itemIndex)
    }
    return shortcut
}
