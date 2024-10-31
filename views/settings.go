package views

import (
	"errors"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/viper"
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

func Settings(a fyne.App) fyne.Window {
	settingsWindow := a.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(200, 150))
	settingsWindow.SetFixedSize(true)
	delayConfig := viper.GetFloat64("delay")
	shuffleConfig := viper.GetBool("shuffle")

	delayEntry := widget.NewEntry()
	delayEntry.Validator = floatValidator
	delayEntry.SetText(strconv.FormatFloat(delayConfig, 'f', -1, 64))
	delayEntry.OnChanged = func(s string) {
		if err := delayEntry.Validate(); err != nil {
			return
		}
		delay, _ := strconv.ParseFloat(s, 64)
		viper.Set("delay", delay)
		if err := viper.WriteConfig(); err != nil {
			log.Println("Error writing config:", err)
		}
	}

	shuffleCheck := widget.NewCheck("Shuffle", func(bool) {})
	shuffleCheck.SetChecked(shuffleConfig)
	shuffleCheck.OnChanged = func(checked bool) {
		viper.Set("shuffle", checked)
		if err := viper.WriteConfig(); err != nil {
			log.Println("Error writing config:", err)
		}
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
	return settingsWindow
}
