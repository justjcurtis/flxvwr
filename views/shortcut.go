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
	"Delay +/- 0.5":    "Up/Down",
	"Shuffle":          "S",
	"Settings":         "/",
	"Pan LRDU":         "HLJK",
	"Zoom +/-":         "+/-",
	"Reset Zoom & Pan": "R",
}

func ShortcutKey() *fyne.Container {
	title := canvas.NewText("Shortcuts", color.White)
	title.TextSize = 20

	quitTitle := canvas.NewText("Quit", color.White)
	quitTitle.TextStyle = fyne.TextStyle{Bold: true}
	quitShortcut := canvas.NewText("Q/Escape", color.White)

	clearTitle := canvas.NewText("Clear", color.White)
	clearTitle.TextStyle = fyne.TextStyle{Bold: true}
	clearShortcut := canvas.NewText("C", color.White)

	playpauseTitle := canvas.NewText("Play/Pause", color.White)
	playpauseTitle.TextStyle = fyne.TextStyle{Bold: true}
	playpauseShortcut := canvas.NewText("Space", color.White)

	prevnextTitle := canvas.NewText("Prev/Next", color.White)
	prevnextTitle.TextStyle = fyne.TextStyle{Bold: true}
	prevnextShortcut := canvas.NewText("Left/Right", color.White)

	delayTitle := canvas.NewText("Delay 0.5 +/-", color.White)
	delayTitle.TextStyle = fyne.TextStyle{Bold: true}
	delayShortcut := canvas.NewText("Up/Down", color.White)

	shuffleTitle := canvas.NewText("Shuffle", color.White)
	shuffleTitle.TextStyle = fyne.TextStyle{Bold: true}
	shuffleShortcut := canvas.NewText("S", color.White)

	settingsTitle := canvas.NewText("Settings", color.White)
	settingsTitle.TextStyle = fyne.TextStyle{Bold: true}
	settingsShortcut := canvas.NewText("/", color.White)

	box := container.NewVBox(
		container.NewCenter(title),
		container.NewGridWithColumns(2,
			quitTitle, quitShortcut,
			clearTitle, clearShortcut,
			playpauseTitle, playpauseShortcut,
			prevnextTitle, prevnextShortcut,
			delayTitle, delayShortcut,
			shuffleTitle, shuffleShortcut,
			settingsTitle, settingsShortcut,
		),
	)
	return box
}
