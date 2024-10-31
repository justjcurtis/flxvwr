package utils

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
)

func GetKey(part string) (fyne.KeyModifier, fyne.KeyName, error) {
	switch part {
	case "CTRL":
		return fyne.KeyModifierControl, "", nil
	case "SHIFT":
		return fyne.KeyModifierShift, "", nil
	case "ALT":
		return fyne.KeyModifierAlt, "", nil
	case "SUPER":
		return fyne.KeyModifierSuper, "", nil
	case "ESC":
		return 0, fyne.KeyEscape, nil
	case "ENTER":
		return 0, fyne.KeyReturn, nil
	case "SPACE":
		return 0, fyne.KeySpace, nil
	case "TAB":
		return 0, fyne.KeyTab, nil
	case "BACKSPACE":
		return 0, fyne.KeyBackspace, nil
	case "DELETE":
		return 0, fyne.KeyDelete, nil
	case "UP":
		return 0, fyne.KeyUp, nil
	case "DOWN":
		return 0, fyne.KeyDown, nil
	case "LEFT":
		return 0, fyne.KeyLeft, nil
	case "RIGHT":
		return 0, fyne.KeyRight, nil
	case "F1":
		return 0, fyne.KeyF1, nil
	case "F2":
		return 0, fyne.KeyF2, nil
	case "F3":
		return 0, fyne.KeyF3, nil
	case "F4":
		return 0, fyne.KeyF4, nil
	case "F5":
		return 0, fyne.KeyF5, nil
	case "F6":
		return 0, fyne.KeyF6, nil
	case "F7":
		return 0, fyne.KeyF7, nil
	case "F8":
		return 0, fyne.KeyF8, nil
	case "F9":
		return 0, fyne.KeyF9, nil
	case "F10":
		return 0, fyne.KeyF10, nil
	case "F11":
		return 0, fyne.KeyF11, nil
	case "F12":
		return 0, fyne.KeyF12, nil
	case "A":
		return 0, fyne.KeyA, nil
	case "B":
		return 0, fyne.KeyB, nil
	case "C":
		return 0, fyne.KeyC, nil
	case "D":
		return 0, fyne.KeyD, nil
	case "E":
		return 0, fyne.KeyE, nil
	case "F":
		return 0, fyne.KeyF, nil
	case "G":
		return 0, fyne.KeyG, nil
	case "H":
		return 0, fyne.KeyH, nil
	case "I":
		return 0, fyne.KeyI, nil
	case "J":
		return 0, fyne.KeyJ, nil
	case "K":
		return 0, fyne.KeyK, nil
	case "L":
		return 0, fyne.KeyL, nil
	case "M":
		return 0, fyne.KeyM, nil
	case "N":
		return 0, fyne.KeyN, nil
	case "O":
		return 0, fyne.KeyO, nil
	case "P":
		return 0, fyne.KeyP, nil
	case "Q":
		return 0, fyne.KeyQ, nil
	case "R":
		return 0, fyne.KeyR, nil
	case "S":
		return 0, fyne.KeyS, nil
	case "T":
		return 0, fyne.KeyT, nil
	case "U":
		return 0, fyne.KeyU, nil
	case "V":
		return 0, fyne.KeyV, nil
	case "W":
		return 0, fyne.KeyW, nil
	case "X":
		return 0, fyne.KeyX, nil
	case "Y":
		return 0, fyne.KeyY, nil
	case "Z":
		return 0, fyne.KeyZ, nil
	case "0":
		return 0, fyne.Key0, nil
	case "1":
		return 0, fyne.Key1, nil
	case "2":
		return 0, fyne.Key2, nil
	case "3":
		return 0, fyne.Key3, nil
	case "4":
		return 0, fyne.Key4, nil
	case "5":
		return 0, fyne.Key5, nil
	case "6":
		return 0, fyne.Key6, nil
	case "7":
		return 0, fyne.Key7, nil
	case "8":
		return 0, fyne.Key8, nil
	case "9":
		return 0, fyne.Key9, nil
	default:
		log.Printf("unknown key: %s", part)
		return 0, "", fmt.Errorf("unknown key: %s", part)
	}
}
