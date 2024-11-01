package services

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func setupConfig() {
	viper.SetDefault("delay", 5)
}

func TestNewPlayerService(t *testing.T) {
	setupConfig()
	ps := NewPlayerService()

	expectedDelay := viper.GetDuration("delay") * time.Second
	if ps.CurrentDelay != expectedDelay {
		t.Errorf("Expected CurrentDelay to be %v, got %v", expectedDelay, ps.CurrentDelay)
	}
	if ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be false, got %v", ps.IsPlaying)
	}
	if ps.offset != 0 {
		t.Errorf("Expected offset to be 0, got %v", ps.offset)
	}
}

func TestPlayPause(t *testing.T) {
	ps := NewPlayerService()

	ps.PlayPause()
	if !ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be true after first PlayPause, got %v", ps.IsPlaying)
	}
	if ps.offset != 0 {
		t.Errorf("Expected offset to be 0 after starting play, got %v", ps.offset)
	}

	time.Sleep(100 * time.Millisecond)

	ps.PlayPause()
	if ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be false after second PlayPause, got %v", ps.IsPlaying)
	}
	if ps.offset == 0 {
		t.Errorf("Expected offset to be non-zero after pause, got %v", ps.offset)
	}

	pausedOffset := ps.offset

	ps.PlayPause()
	if !ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be true after resuming, got %v", ps.IsPlaying)
	}
	if ps.offset != 0 {
		t.Errorf("Expected offset to be reset to 0 after resuming, got %v", ps.offset)
	}

	timeSinceLastSet := time.Since(ps.LastSet)
	if timeSinceLastSet-pausedOffset > 1*time.Millisecond {
		t.Errorf("Expected LastSet to account for offset after resuming, but it did not")
	}
}

func TestStop(t *testing.T) {
	ps := NewPlayerService()
	ps.PlayPause()

	if !ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be true before Stop, got %v", ps.IsPlaying)
	}

	ps.Stop()
	if ps.IsPlaying {
		t.Errorf("Expected IsPlaying to be false after Stop, got %v", ps.IsPlaying)
	}
	if ps.offset != 0 {
		t.Errorf("Expected offset to be reset to 0 after Stop, got %v", ps.offset)
	}
}
