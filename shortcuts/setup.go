package shortcuts

import (
	"flxvwr/services"
	"flxvwr/views"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"github.com/spf13/viper"
)

var keypressThreshold = 500 * time.Millisecond

type modifiers struct {
	super bool
	ctrl  bool
	shift bool
	alt   bool
}

func (m *modifiers) reset() {
	m.super = false
	m.ctrl = false
	m.shift = false
	m.alt = false
}
func (m *modifiers) getMods() string {
	mods := ""
	if m.super {
		mods += "Super+"
	}
	if m.ctrl {
		mods += "Ctrl+"
	}
	if m.shift {
		mods += "Shift+"
	}
	if m.alt {
		mods += "Alt+"
	}
	return mods
}

func SetupShortcuts(a fyne.App, w fyne.Window, is *services.ImageService, ps *services.PlayerService, ns *services.NotificationService) {
	mods := modifiers{}
	lastKeyTime := time.Now().Add(-keypressThreshold * 2)
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if time.Since(lastKeyTime) > keypressThreshold {
			mods.reset()
		}
		lastKeyTime = time.Now()
		if e.Name == "SuperLeft" || e.Name == "SuperRight" {
			mods.super = true
			return
		}
		if e.Name == "ControlLeft" || e.Name == "ControlRight" {
			mods.ctrl = true
			return
		}
		if e.Name == "LeftShift" || e.Name == "RightShift" {
			mods.shift = true
			return
		}
		if e.Name == "AltLeft" || e.Name == "AltRight" {
			mods.alt = true
			return
		}

		input := mods.getMods() + fmt.Sprint(e.Name)

		if input == "Escape" || input == "Q" {
			a.Quit()
		}
		if input == "C" {
			ns.SetNotification("Cleared")
			ps.Stop()
			is.Clear()
			w.SetContent(views.StartView(a))
		}
		if input == "Right" {
			isPlaying := ps.IsPlaying
			ps.Stop()
			is.Next()
			is.Update(w, ps, true)
			ps.PlayPause()
			if !isPlaying {
				ps.PlayPause()
			}
		}
		if input == "Left" {
			isPlaying := ps.IsPlaying
			ps.Stop()
			is.Previous()
			is.Update(w, ps, true)
			ps.PlayPause()
			if !isPlaying {
				ps.PlayPause()
			}
		}
		if input == "Space" {
			if ps.IsPlaying {
				ns.SetNotification("Paused")
			} else {
				ns.SetNotification("Playing")
			}
			ps.PlayPause()
		}
		if input == "S" {
			shuffle := viper.GetBool("shuffle")
			viper.Set("shuffle", !shuffle)
			if !shuffle {
				ns.SetNotification("Shuffle On")
			} else {
				ns.SetNotification("Shuffle Off")
			}
			is.RecalculatePlaylist()
		}
		if input == "Up" {
			nextDelay := viper.GetFloat64("delay") + 1
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if input == "Shift+Up" {
			nextDelay := viper.GetFloat64("delay") + 0.5
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if input == "Down" {
			nextDelay := viper.GetFloat64("delay") - 0.5
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			if nextDelay < 1 {
				nextDelay = 1
			}
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if input == "Shift+Down" {
			nextDelay := viper.GetFloat64("delay") - 1
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay, 'f', 1, 64))
			if nextDelay < 1 {
				nextDelay = 1
			}
			viper.Set("delay", nextDelay)
			ps.CurrentDelay = viper.GetDuration("delay") * time.Second
		}
		if input == "/" {
			settingsWindow := views.Settings(a)
			settingsWindow.Show()
		}
		if input == "R" {
			is.Zoomable.Reset()
			is.Brightness = 1
			is.Contrast = 1
			is.Update(w, ps, false)
		}
		if input == "K" {
			is.Zoomable.Move(0, 20)
		}
		if input == "Shift+K" {
			is.Zoomable.Move(0, 5)
		}
		if input == "H" {
			is.Zoomable.Move(20, 0)
		}
		if input == "Shift+H" {
			is.Zoomable.Move(5, 0)
		}
		if input == "J" {
			is.Zoomable.Move(0, -20)
		}
		if input == "Shift+J" {
			is.Zoomable.Move(0, -5)
		}
		if input == "L" {
			is.Zoomable.Move(-20, 0)
		}
		if input == "Shift+L" {
			is.Zoomable.Move(-5, 0)
		}
		if input == "=" || input == "+" {
			is.Zoomable.Zoom(0.5)
		}
		if input == "Shift+=" || input == "Shift++" {
			is.Zoomable.Zoom(0.2)
		}
		if input == "-" || input == "_" {
			is.Zoomable.Zoom(-0.5)
		}
		if input == "Shift+-" || input == "Shift+_" {
			is.Zoomable.Zoom(-0.2)
		}
		if input == "B" {
			is.Brightness -= 0.05
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "Shift+B" {
			is.Brightness -= 0.01
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "N" {
			is.Brightness += 0.05
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "Shift+N" {
			is.Brightness += 0.01
			ns.SetNotification("Brightness " + strconv.FormatFloat(is.Brightness*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "V" {
			is.Contrast -= 0.05
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "Shift+V" {
			is.Contrast -= 0.01
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "M" {
			is.Contrast += 0.05
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
		if input == "Shift+M" {
			is.Contrast += 0.01
			ns.SetNotification("Contrast " + strconv.FormatFloat(is.Contrast*100, 'f', 0, 64) + "%")
			is.Update(w, ps, false)
		}
	})
}
