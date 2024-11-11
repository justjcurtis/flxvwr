package models

import "time"

type Config struct {
	Delay   time.Duration
	Shuffle bool
}
