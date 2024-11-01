package utils

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestGetKey(t *testing.T) {
	tests := []struct {
		part      string
		modifier  fyne.KeyModifier
		key       fyne.KeyName
		expectErr bool
	}{
		{"CTRL", fyne.KeyModifierControl, "", false},
		{"SHIFT", fyne.KeyModifierShift, "", false},
		{"ALT", fyne.KeyModifierAlt, "", false},
		{"SUPER", fyne.KeyModifierSuper, "", false},
		{"ESC", 0, fyne.KeyEscape, false},
		{"ENTER", 0, fyne.KeyReturn, false},
		{"SPACE", 0, fyne.KeySpace, false},
		{"TAB", 0, fyne.KeyTab, false},
		{"BACKSPACE", 0, fyne.KeyBackspace, false},
		{"DELETE", 0, fyne.KeyDelete, false},
		{"UP", 0, fyne.KeyUp, false},
		{"DOWN", 0, fyne.KeyDown, false},
		{"LEFT", 0, fyne.KeyLeft, false},
		{"RIGHT", 0, fyne.KeyRight, false},
		{"F1", 0, fyne.KeyF1, false},
		{"F2", 0, fyne.KeyF2, false},
		{"A", 0, fyne.KeyA, false},
		{"B", 0, fyne.KeyB, false},
		{"C", 0, fyne.KeyC, false},
		{"0", 0, fyne.Key0, false},
		{"9", 0, fyne.Key9, false},
		{"UNKNOWN", 0, "", true}, // testing an unknown key
	}

	for _, tt := range tests {
		t.Run(tt.part, func(t *testing.T) {
			modifier, key, err := GetKey(tt.part)

			if (err != nil) != tt.expectErr {
				t.Errorf("GetKey() error = %v, expectErr %v", err, tt.expectErr)
			}
			if modifier != tt.modifier {
				t.Errorf("expected modifier %v, got %v", tt.modifier, modifier)
			}
			if key != tt.key {
				t.Errorf("expected key %v, got %v", tt.key, key)
			}
		})
	}
}
