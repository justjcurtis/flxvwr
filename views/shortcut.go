package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var ShortcutMap map[string]string = map[string]string{
	"Quit":             "Q/Escape",
	"Clear":            "C",
	"Play/Pause":       "Space",
	"Prev/Next":        "Left/Right",
	"Delay +/-":        "Up/Down",
	"Shuffle":          "S",
	"Zoom +/-":         "+/-",
	"Pan LRDU":         "HLJK",
	"Reset Zoom & Pan": "R",
	"Brightness +/-":   "B/N",
	"Contrast +/-":     "V/M",
	"Settings":         "/",
}

func ShortcutKey() *fyne.Container {
	title := canvas.NewText("Shortcuts", color.White)
	title.TextSize = 20

	textElements := container.NewGridWithColumns(2)
	for k, v := range ShortcutMap {
		title := canvas.NewText(k, color.White)
		title.TextStyle = fyne.TextStyle{Bold: true}
		text := canvas.NewText(v, color.White)
		textElements.Add(title)
		textElements.Add(text)
	}

	box := container.NewVBox(
		container.NewCenter(title),
		textElements,
	)
	return box
}
