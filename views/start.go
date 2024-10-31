package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func StartView(a fyne.App) *fyne.Container {
	hello := widget.NewLabel("Drop images/directories here")
	return container.NewCenter(
		container.NewVBox(
			hello,
			widget.NewButton("Settings", func() {
				settingsWindow := Settings(a)
				settingsWindow.Show()
			}),
			container.NewPadded(),
			ShortcutKey(),
		),
	)
}
