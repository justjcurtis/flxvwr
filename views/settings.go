package views

import (
	"errors"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/justjcurtis/flxvwr/services"
)

func floatValidator(s string) error {
	if s == "" {
		return errors.New("Empty field")
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return errors.New("Invalid number")
	} else if n, _ := strconv.ParseFloat(s, 64); n < 1 {
		return errors.New("Number must be greater than 1")
	}
	return nil
}

func Settings(a fyne.App, cs *services.ConfigService) fyne.Window {
	settingsWindow := a.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(200, 150))
	settingsWindow.SetFixedSize(true)
	delayConfig := cs.GetDelay().Seconds()
	shuffleConfig := cs.GetShuffle()

	delayEntry := widget.NewEntry()
	delayEntry.Validator = floatValidator
	delayEntry.SetText(strconv.FormatFloat(delayConfig, 'f', -1, 64))
	delayEntry.OnChanged = func(s string) {
		if err := delayEntry.Validate(); err != nil {
			return
		}
		delay, _ := strconv.ParseFloat(s, 64)
		delayDuration := time.Duration(delay) * time.Second
		cs.SetDelay(delayDuration)
	}

	shuffleCheck := widget.NewCheck("Shuffle", func(bool) {})
	shuffleCheck.SetChecked(shuffleConfig)
	shuffleCheck.OnChanged = func(checked bool) {
		cs.SetShuffle(checked)
	}

	settingsWindow.SetContent(container.NewVBox(
		widget.NewLabel("Settings"),
		container.NewGridWithColumns(2,
			widget.NewLabel("Delay"),
			delayEntry,
			widget.NewLabel("Shuffle"),
			shuffleCheck,
		),
		widget.NewButton("Close", func() {
			settingsWindow.Close()
		}),
	))

	settingsWindow.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if e.Name == "Escape" {
			settingsWindow.Close()
		}
		if e.Name == "Return" {
			settingsWindow.Close()
		}
		if e.Name == "Q" {
			settingsWindow.Close()
		}
		if e.Name == "/" {
			settingsWindow.Close()
		}
	})

	return settingsWindow
}
