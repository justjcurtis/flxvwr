package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigPath() string {
	var configPath string
	o := runtime.GOOS
	switch o {
	case "windows":
		configPath = filepath.Join(os.Getenv("APPDATA"), "flxvwr")
	case "linux":
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			configPath = filepath.Join(xdgConfig, "flxvwr")
		} else {
			configPath = "/etc/flxvwr"
		}
	case "darwin":
		configPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "flxvwr")
	default:
		log.Fatalf("Unsupported OS: %s", o)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, os.ModePerm)
	}

	return configPath
}
