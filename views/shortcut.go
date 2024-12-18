package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var ShortcutMap [][]string = [][]string{
	{"Quit", "Q"},
	{"Quit All", "Esc"},
	{"Clear", "C"},
	{"Play/Pause", "Space"},
	{"Prev/Next", "Left/Right"},
	{"Delay +/-", "Up/Down"},
	{"Shuffle", "S"},
	{"Zoom +/-", "+/-"},
	{"Pan LRDU", "HLJK"},
	{"Rotate", "[ / ]"},
	{"Reset image", "R"},
	{"Brightness +/-", "B/N"},
	{"Contrast +/-", "V/M"},
	{"Smaller increments", "Shift+..."},
	{"Add current to playlist", "Shift+1 to 9"},
	{"Switch to playlist", "1 to 9"},
	{"Remove from playlist", "X"},
	{"Export/Save current playlist", "E"},
	{"Show/Hide shortcuts", "F1 | ?"},
	{"Settings", "/"},
}

func ShortcutKey() *fyne.Container {
	title := canvas.NewText("Shortcuts", color.White)
	title.TextSize = 20

	textElements := container.NewGridWithColumns(2)
	for _, row := range ShortcutMap {
		k, v := row[0], row[1]
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
