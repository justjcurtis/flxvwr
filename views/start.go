package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func StartView(a fyne.App) *fyne.Container {
	hello := canvas.NewText("Drop images/directories here", color.White)
	hello.TextStyle.Bold = true
	hello.TextSize = 24
	return container.NewCenter(
		container.NewVBox(
			hello,
			container.NewPadded(
				ShortcutKey(),
			),
		),
	)
}
