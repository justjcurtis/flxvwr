package main

import (
	"flxvwr/services"
	"flxvwr/shortcuts"
	"flxvwr/utils"
	"flxvwr/views"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ImageService *services.ImageService
var PlayerService *services.PlayerService
var NotificationService *services.NotificationService

func main() {
	configPath := utils.GetConfigPath()
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

	go func() {
		for {
			if PlayerService.IsPlaying {
				if time.Since(PlayerService.LastSet) >= PlayerService.CurrentDelay {
					ImageService.Next()
					ImageService.Update(w, PlayerService)
				}
			}
		}
	}()

	w.SetContent(views.StartView(a))

	w.SetOnDropped(func(pos fyne.Position, uri []fyne.URI) {
		ImageService.ImportImages(pos, uri)
		ImageService.Update(w, PlayerService)
		PlayerService.LastSet = time.Now()
		PlayerService.IsPlaying = true
	})
	w.ShowAndRun()
}
