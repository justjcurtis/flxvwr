package shortcuts

import (
	"flxvwr/services"
	"flxvwr/views"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"github.com/spf13/viper"
)

func SetupShortcuts(a fyne.App, w fyne.Window, is *services.ImageService, ps *services.PlayerService, ns *services.NotificationService) {
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if e.Name == "Escape" || e.Name == "Q" {
			a.Quit()
		}
		if e.Name == "C" {
			ns.SetNotification("Cleared")
			ps.Stop()
			is.Clear()
			w.SetContent(views.StartView(a))
		}
		if e.Name == "Right" {
			isPlaying := ps.IsPlaying
			ps.Stop()
			is.Next()
			is.Update(w, ps, true)
			ps.PlayPause()
			if !isPlaying {
				ps.PlayPause()
			}
		}
		if e.Name == "Left" {
			isPlaying := ps.IsPlaying
			ps.Stop()
			is.Previous()
			is.Update(w, ps, true)
			ps.PlayPause()
			if !isPlaying {
				ps.PlayPause()
			}
		}
		if e.Name == "Space" {
			if ps.IsPlaying {
				ns.SetNotification("Paused")
			} else {
				ns.SetNotification("Playing")
			}
			ps.PlayPause()
		}
		if e.Name == "S" {
			shuffle := viper.GetBool("shuffle")
			viper.Set("shuffle", !shuffle)
			if !shuffle {
				ns.SetNotification("Shuffle On")
			} else {
				ns.SetNotification("Shuffle Off")
			}
			is.RecalculatePlaylist()
		}
		if e.Name == "Up" {
			nextDelay := viper.GetFloat64("delay") + 0.5
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if e.Name == "Down" {
			nextDelay := viper.GetFloat64("delay") - 0.5
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			if nextDelay < 1 {
				nextDelay = 1
			}
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if e.Name == "/" {
			settingsWindow := views.Settings(a)
			settingsWindow.Show()
		}
		if e.Name == "R" {
			is.Zoomable.Reset()
			is.Brightness = 1
			is.Contrast = 1
			is.Update(w, ps, false)
		}
		if e.Name == "K" {
			is.Zoomable.Move(0, 10)
		}
		if e.Name == "H" {
			is.Zoomable.Move(10, 0)
		}
		if e.Name == "J" {
			is.Zoomable.Move(0, -10)
		}
		if e.Name == "L" {
			is.Zoomable.Move(-10, 0)
		}
		if e.Name == "=" || e.Name == "+" {
			is.Zoomable.Zoom(0.2)
		}
		if e.Name == "-" || e.Name == "_" {
			is.Zoomable.Zoom(-0.2)
		}
		if e.Name == "B" {
			is.Brightness -= 0.01
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness, 'f', 2, 64))
			is.Update(w, ps, false)
		}
		if e.Name == "N" {
			is.Brightness += 0.01
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness, 'f', 2, 64))
			is.Update(w, ps, false)
		}
		if e.Name == "V" {
			is.Contrast -= 0.01
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast, 'f', 2, 64))
			is.Update(w, ps, false)
		}
		if e.Name == "M" {
			is.Contrast += 0.01
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast, 'f', 2, 64))
			is.Update(w, ps, false)
		}
	})
}
