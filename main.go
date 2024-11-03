package main

import (
	"log"
	"time"

	"github.com/justjcurtis/flxvwr/services"
	"github.com/justjcurtis/flxvwr/shortcuts"
	"github.com/justjcurtis/flxvwr/utils"
	"github.com/justjcurtis/flxvwr/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ImageService *services.ImageService
var PlayerService *services.PlayerService
var NotificationService *services.NotificationService

func main() {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.SetDefault("app.name", "flxvwr")
	viper.SetDefault("app.version", "0.0.1")
	viper.SetDefault("shuffle", true)
	viper.SetDefault("delay", 10.0)

	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using defaults")
	}

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("flxvwr")

	w.Resize(fyne.NewSize(800, 600))

	ImageService = services.NewImageService()
	PlayerService := services.NewPlayerService()
	NotificationService := services.NewNotificationService(w)

	viper.OnConfigChange(func(e fsnotify.Event) {
		newDelay := viper.GetDuration("delay") * time.Second
		newShuffle := viper.GetBool("shuffle")
		if newDelay != PlayerService.CurrentDelay {
			PlayerService.CurrentDelay = newDelay
		}
		if newShuffle != viper.GetBool("shuffle") {
			ImageService.RecalculatePlaylist()
		}
	})
	viper.WatchConfig()

	shortcuts.SetupShortcuts(a, w, ImageService, PlayerService, NotificationService)

	ticker := time.NewTicker(100 * time.Millisecond)
	handleResize := func() {
		ImageService.Update(w, PlayerService, false)
	}
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

	w.SetOnDropped(func(pos fyne.Position, uri []fyne.URI) {
		ImageService.ImportImages(pos, uri)
		ImageService.Update(w, PlayerService, true)
		PlayerService.LastSet = time.Now()
		PlayerService.IsPlaying = true
	})
	w.ShowAndRun()
}
