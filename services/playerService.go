package services

import (
	"time"

	"github.com/justjcurtis/flxvwr/models"
	"github.com/spf13/viper"
)

type PlayerService struct {
	IsPlaying    bool
	LastSet      time.Time
	CurrentDelay time.Duration
	offset       time.Duration
}

func NewPlayerService() *PlayerService {
	delay := viper.GetDuration("delay") * time.Second
	return &PlayerService{
		IsPlaying:    false,
		LastSet:      time.Now(),
		CurrentDelay: delay,
	}
}

func (ps *PlayerService) HandleConfigUpdate(config models.Config) {
	ps.CurrentDelay = config.Delay
}

func (ps *PlayerService) PlayPause() {
	now := time.Now()
	if ps.IsPlaying {
		ps.offset = now.Sub(ps.LastSet)
		ps.IsPlaying = false
		return
	}
	ps.IsPlaying = true
	ps.LastSet = now.Add(-ps.offset)
	ps.offset = 0
}

func (ps *PlayerService) Stop() {
	ps.IsPlaying = false
	ps.offset = 0
}
