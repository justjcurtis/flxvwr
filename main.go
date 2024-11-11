package main

import (
	"fmt"
	"time"

	"github.com/justjcurtis/flxvwr/services"
	"github.com/justjcurtis/flxvwr/shortcuts"
	"github.com/justjcurtis/flxvwr/utils"
	"github.com/justjcurtis/flxvwr/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var ImageService *services.ImageService
var PlayerService *services.PlayerService
var NotificationService *services.NotificationService
var ConfigService *services.ConfigService

func main() {

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme()) // TODO: create custom dark theme as theme.DarkTheme() is deprecated
	w := a.NewWindow("flxvwr")
	fmt.Println("flxvwr starting...")
	w.Resize(fyne.NewSize(800, 600))

	fmt.Println("Creating services...")
	ConfigService := services.NewConfigService()
	fmt.Println("Subscribing to config updates...")
	ImageService = services.NewImageService()
	ConfigService.Subscribe(ImageService.HandleConfigUpdate)
	PlayerService := services.NewPlayerService()
	ConfigService.Subscribe(PlayerService.HandleConfigUpdate)
	NotificationService := services.NewNotificationService(w)
	fmt.Println("finishing services...")

	shortcuts.SetupShortcuts(a, w, ImageService, PlayerService, NotificationService, ConfigService)
	fmt.Println("Shortcuts setup...")

	ticker := time.NewTicker(100 * time.Millisecond)
	handleResize := utils.Debounce(func() {
		ImageService.Update(w, PlayerService, false)
	}, 100*time.Millisecond)

	currentWidth := w.Canvas().Size().Width
	currentHeight := w.Canvas().Size().Height
	go func() {
		for range ticker.C {
			if currentWidth != w.Canvas().Size().Width || currentHeight != w.Canvas().Size().Height {
				handleResize()
				currentWidth = w.Canvas().Size().Width
				currentHeight = w.Canvas().Size().Height
			}
			if PlayerService.IsPlaying {
				if time.Since(PlayerService.LastSet) >= PlayerService.CurrentDelay {
					ImageService.Next()
					ImageService.Update(w, PlayerService, true)
				}
			}
		}
	}()

	w.SetContent(views.StartView(a))

	args := utils.GetArgs()

	w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		ImageService.ImportImages(uris)
		if ImageService.GetCurrent() != nil {
			ImageService.Update(w, PlayerService, true)
			PlayerService.LastSet = time.Now()
			PlayerService.IsPlaying = true
		}
	})

	if len(args) > 0 {
		uris := utils.GetURIsFromLines(args)
		ImageService.ImportImages(uris)
		if ImageService.GetCurrent() != nil {
			ImageService.Update(w, PlayerService, true)
			PlayerService.LastSet = time.Now()
			PlayerService.IsPlaying = true
		}
	}

	w.ShowAndRun()
}
