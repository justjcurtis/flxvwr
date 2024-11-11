package shortcuts

import (
	"fmt"
	"strconv"
	"time"

	"github.com/justjcurtis/flxvwr/services"
	"github.com/justjcurtis/flxvwr/utils"
	"github.com/justjcurtis/flxvwr/views"

	"fyne.io/fyne/v2"
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

func SetupShortcuts(a fyne.App, w fyne.Window, is *services.ImageService, ps *services.PlayerService, ns *services.NotificationService, cs *services.ConfigService) {
	mods := modifiers{}
	isShowingShortcuts := false
	wasPlaying := false
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

		if input == "Shift+/" || input == "F1" {
			if !isShowingShortcuts {
				if ps.IsPlaying {
					wasPlaying = true
					ns.SetNotification("Paused")
					ps.PlayPause()
				} else {
					wasPlaying = false
				}
				isShowingShortcuts = true
				w.SetContent(views.StartView(a))
				return
			}
			isShowingShortcuts = false
			is.Update(w, ps, false)
			if !ps.IsPlaying && wasPlaying {
				ns.SetNotification("Playing")
				ps.PlayPause()
			}
		}
		if input == "Escape" {
			utils.KillAppInstances("flxvwr")
			a.Quit()
		}
		if input == "Q" {
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
			shuffle := cs.GetShuffle()
			cs.SetShuffle(!shuffle)
			if !shuffle {
				ns.SetNotification("Shuffle On")
			} else {
				ns.SetNotification("Shuffle Off")
			}
		}
		if input == "Up" {
			nextDelay := cs.GetDelay() + 1*time.Second
			cs.SetDelay(nextDelay)
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay.Seconds(), 'f', 1, 64))
		}
		if input == "Shift+Up" {
			nextDelay := cs.GetDelay() + (500 * time.Millisecond)
			cs.SetDelay(nextDelay)
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay.Seconds(), 'f', 1, 64))
		}
		if input == "Down" {
			nextDelay := cs.GetDelay() - 1*time.Second
			if nextDelay < 500*time.Millisecond {
				nextDelay = 500 * time.Millisecond
			}
			cs.SetDelay(nextDelay)
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay.Seconds(), 'f', 1, 64))
		}
		if input == "Shift+Down" {
			nextDelay := cs.GetDelay() - (500 * time.Millisecond)
			if nextDelay < 500*time.Millisecond {
				nextDelay = 500 * time.Millisecond
			}
			cs.SetDelay(nextDelay)
			ns.SetNotification("Delay " + strconv.FormatFloat(nextDelay.Seconds(), 'f', 1, 64))
		}
		if input == "/" {
			if ps.IsPlaying {
				wasPlaying = true
				ps.PlayPause()
				ns.SetNotification("Paused")
			} else {
				wasPlaying = false
			}
			settingsWindow := views.Settings(a, cs)
			settingsWindow.Show()
		}
		if input == "R" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Reset()
			is.Update(w, ps, false)
		}
		if input == "K" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(0, -20)
		}
		if input == "Shift+K" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(0, -5)
		}
		if input == "H" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(-20, 0)
		}
		if input == "Shift+H" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(-5, 0)
		}
		if input == "J" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(0, 20)
		}
		if input == "Shift+J" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(0, 5)
		}
		if input == "L" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(20, 0)
		}
		if input == "Shift+L" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Move(5, 0)
		}
		if input == "=" || input == "+" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Zoom(0.5)
		}
		if input == "Shift+=" || input == "Shift++" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Zoom(0.2)
		}
		if input == "-" || input == "_" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Zoom(-0.5)
		}
		if input == "Shift+-" || input == "Shift+_" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Zoom(-0.2)
		}
		if input == "[" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Rotate(3)
		}
		if input == "]" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.Rotate(1)
		}
		if input == "B" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(-0.05, 0)
			ns.SetNotification("Brightness " + strconv.FormatFloat(float64(is.Zoomable.Brightness*100), 'f', 0, 32) + "%")
		}
		if input == "Shift+B" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(-0.01, 0)
			ns.SetNotification("Brightness " + strconv.FormatFloat(float64(is.Zoomable.Brightness*100), 'f', 0, 32) + "%")
		}
		if input == "N" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0.05, 0)
			ns.SetNotification("Brightness " + strconv.FormatFloat(float64(is.Zoomable.Brightness*100), 'f', 0, 32) + "%")
		}
		if input == "Shift+N" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0.01, 0)
			ns.SetNotification("Brightness " + strconv.FormatFloat(float64(is.Zoomable.Brightness*100), 'f', 0, 32) + "%")
		}
		if input == "V" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0, -0.05)
			ns.SetNotification("Contrast " + strconv.FormatFloat(float64(is.Zoomable.Contrast*100), 'f', 0, 32) + "%")
		}
		if input == "Shift+V" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0, 0.01)
			ns.SetNotification("Contrast " + strconv.FormatFloat(float64(is.Zoomable.Contrast*100), 'f', 0, 32) + "%")
		}
		if input == "M" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0, 0.05)
			ns.SetNotification("Contrast " + strconv.FormatFloat(float64(is.Zoomable.Contrast*100), 'f', 0, 32) + "%")
		}
		if input == "Shift+M" {
			if is.Zoomable == nil {
				return
			}
			is.Zoomable.AdjustBrightnessAndContrast(0, 0.01)
			ns.SetNotification("Contrast " + strconv.FormatFloat(float64(is.Zoomable.Contrast*100), 'f', 0, 32) + "%")
		}
		if input == "1" || input == "2" || input == "3" || input == "4" || input == "5" || input == "6" || input == "7" || input == "8" || input == "9" || input == "0" {
			index, _ := strconv.Atoi(input)
			is.SetCurrentPlaylist(index)
			is.Zoomable.Reset()
			is.Update(w, ps, true)
			ns.SetNotification("Playlist " + input)
		}
		if input == "Shift+1" || input == "Shift+2" || input == "Shift+3" || input == "Shift+4" || input == "Shift+5" || input == "Shift+6" || input == "Shift+7" || input == "Shift+8" || input == "Shift+9" || input == "Shift+0" {
			index, _ := strconv.Atoi(input[6:])
			is.AddCurrentToPlaylist(index)
			is.Zoomable.Reset()
			is.Update(w, ps, true)
			ns.SetNotification("Added to playlist " + input)
		}
		if input == "X" {
			is.RemoveCurrentFromPlaylist()
			is.Zoomable.Reset()
			is.Update(w, ps, true)
			ns.SetNotification("Removed")
		}
	})
}
