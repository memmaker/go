package cview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/memmaker/go/recfile"
	"strconv"
	"strings"
)

type FieldType uint8

const (
	FieldTypeString FieldType = iota
	FieldTypeIntegerInput
	FieldTypeFloatInput
	FieldTypeIntegerSlider
	FieldTypeSelect
	FieldTypeBool
)

func FieldTypeFromString(s string) FieldType {
	s = strings.ToLower(s)
	switch s {
	case "string":
		return FieldTypeString
	case "int":
		return FieldTypeIntegerInput
	case "float":
		return FieldTypeFloatInput
	case "slider":
		return FieldTypeIntegerSlider
	case "select":
		return FieldTypeSelect
	case "bool":
		return FieldTypeBool
	}
	return FieldTypeString
}

type SelectOption struct {
	DisplayText string
	Value       string
}
type FormElementDescription struct {
	FieldName    string
	Label        string
	PrefillValue TypedString
	FieldType    FieldType
	Options      []SelectOption
	IntMax       int
}

type TypedString string

func TypedStringFromBool(b bool) TypedString {
	return TypedString(strconv.FormatBool(b))
}

func TypedStringFromInt(i int) TypedString {
	return TypedString(strconv.Itoa(i))
}

func TypedStringFromFloat(f float64) TypedString {
	return TypedString(strconv.FormatFloat(f, 'f', -1, 64))
}

func (t TypedString) AsInt() int {
	i, _ := strconv.Atoi(string(t))
	return i
}

func (t TypedString) AsFloat() float64 {
	i, _ := strconv.ParseFloat(string(t), 64)
	return i
}

func (t TypedString) AsBool() bool {
	i, _ := strconv.ParseBool(string(t))
	return i
}

func (t TypedString) AsString() string {
	return string(t)
}

func (t TypedString) AsRune() rune {
	return []rune(t)[0]
}

func EditRecord(app *Application, panels *Panels, rec recfile.Record, onConfirm func(recfile.Record)) {
	OpenModalEditor(app, panels, rec.String(), func(lines []string) {
		modifiedRecord := recfile.RecordFromSlice(lines)
		onConfirm(modifiedRecord)
	})
}

func OpenModalEditor(app *Application, panels *Panels, preFill string, onClose func(lines []string)) {
	closeModal := func() {
		panels.RemovePanel("modal")
		_, frontPanel := panels.GetFrontPanel()
		app.SetBeforeFocusFunc(nil)
		app.SetFocus(frontPanel)
	}

	modal := NewTextArea()
	modal.SetBorder(true)
	modal.SetWrap(false)
	modal.SetWordWrap(false)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			closeModal()
			onClose(strings.Split(modal.GetText(), "\n"))
		}
		return event
	})

	width, height := app.GetScreen().Size()
	x := width / 4
	y := height / 4
	w := width / 2
	h := height / 2
	modal.SetSize(0, 0)
	modal.SetRect(x, y, w, h)

	modal.SetText(preFill, false)
	//modal.SetOffset(0, x)

	panels.AddPanel("modal", modal, false, true)
	app.SetFocus(modal)
	app.SetBeforeFocusFunc(func(p Primitive) bool { return false })
}

func AskForString(app *Application, panels *Panels, prompt, prefill string, onConfirm func(entered string)) *Modal {
	return OpenModalForm(app, panels, []FormElementDescription{
		{
			FieldName:    "text",
			PrefillValue: TypedString(prefill),
			Label:        prompt,
			FieldType:    FieldTypeString,
		},
	}, func(values map[string]TypedString) {
		if text, exists := values["text"]; exists {
			onConfirm(text.AsString())
		}
	})
}

func OpenModalForm(app *Application, panels *Panels, elements []FormElementDescription, confirm func(map[string]TypedString)) *Modal {
	modal := NewModalForm(elements, confirm, func() {
		panels.RemovePanel("modal")
		_, frontPanel := panels.GetFrontPanel()
		// Reset focus changes
		app.SetBeforeFocusFunc(nil)
		app.SetFocus(frontPanel)
	})
	panels.AddPanel("modal", modal, false, true)
	app.SetBeforeFocusFunc(nil)

	modal.GetForm().SetFocus(0)
	app.SetFocus(modal.GetForm().GetFormItem(0))

	// deny any focus change
	app.SetBeforeFocusFunc(func(p Primitive) bool {
		if p == modal {
			return true
		}
		x, y, w, h := p.GetRect()
		if modal.InRect(x, y) && modal.InRect(x+w, y+h) {
			return true
		}
		return false
	})
	return modal
}

func NewModalForm(elements []FormElementDescription, confirm func(map[string]TypedString), close func()) *Modal {
	modal := NewModal()
	form := modal.GetForm()

	for _, element := range elements {
		switch element.FieldType {
		case FieldTypeString:
			value := ""
			if element.PrefillValue != "" {
				value = element.PrefillValue.AsString()
			}
			form.AddInputField(element.Label, value, 30, nil, nil)
		case FieldTypeIntegerInput:
			value := ""
			if element.PrefillValue != "" {
				value = element.PrefillValue.AsString()
			}
			form.AddInputField(element.Label, value, 30, InputFieldInteger, nil)
		case FieldTypeFloatInput:
			value := ""
			if element.PrefillValue != "" {
				value = element.PrefillValue.AsString()
			}
			form.AddInputField(element.Label, value, 30, InputFieldFloat, nil)
		case FieldTypeIntegerSlider:
			curValue := 0
			if element.PrefillValue != "" {
				curValue = element.PrefillValue.AsInt()
			}
			form.AddSlider(element.Label, curValue, element.IntMax, 1, nil)
		case FieldTypeSelect:
			initialOption := 0
			if element.PrefillValue != "" {
				for i, option := range element.Options {
					if option.Value == element.PrefillValue.AsString() {
						initialOption = i
						break
					}
				}
			}
			form.AddDropDown(element.Label, initialOption, nil, optionsFromStrings(element.Options))
		case FieldTypeBool:
			initialValue := false
			if element.PrefillValue != "" {
				initialValue = element.PrefillValue.AsBool()
			}
			form.AddCheckBox(element.Label, "", initialValue, nil)
		}
	}
	modal.AddButtons([]string{"Cancel", "Confirm"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Confirm" && confirm != nil {
			values := make(map[string]TypedString)
			for _, element := range elements {
				switch element.FieldType {
				case FieldTypeString, FieldTypeIntegerInput, FieldTypeFloatInput:
					values[element.FieldName] = TypedString(form.GetFormItemByLabel(element.Label).(*InputField).GetText())
				case FieldTypeIntegerSlider:
					values[element.FieldName] = TypedString(strconv.Itoa(form.GetFormItemByLabel(element.Label).(*Slider).GetProgress()))
				case FieldTypeSelect:
					optionIndex, _ := form.GetFormItemByLabel(element.Label).(*DropDown).GetCurrentOption()
					values[element.FieldName] = TypedString(element.Options[optionIndex].Value)
				case FieldTypeBool:
					values[element.FieldName] = TypedString(strconv.FormatBool(form.GetFormItemByLabel(element.Label).(*CheckBox).IsChecked()))
				}
			}
			confirm(values)
		}
		close()
	})
	return modal
}

func optionsFromStrings(options []SelectOption) []*DropDownOption {
	result := make([]*DropDownOption, len(options))
	for i, option := range options {
		dropDownOption := NewDropDownOption(option.DisplayText)
		dropDownOption.SetReference(option.Value)
		result[i] = dropDownOption
	}
	return result
}
