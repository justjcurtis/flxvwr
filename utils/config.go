package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

var getOS = func() string {
	return runtime.GOOS
}

func GetConfigPath() (string, error) {
	var configPath string
	o := getOS()
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
		return "", os.ErrNotExist
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, os.ModePerm)
	}

	return configPath, nil
}
