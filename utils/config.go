package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var getOS = func() string {
	return runtime.GOOS
}

var getLinuxConfigDir = func() (string, error) {
	return os.UserConfigDir()
}

var mkdirAll = func(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func GetConfigPath() (string, error) {
	var configPath string
	o := getOS()
	switch o {
	case "windows":
		configPath = filepath.Join(os.Getenv("APPDATA"), "flxvwr")
	case "linux":
		configDir, err := getLinuxConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config dir: %w", err)
		}
		configPath = filepath.Join(configDir, "flxvwr")
	case "darwin":
		configPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "flxvwr")
	default:
		return "", os.ErrNotExist
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := mkdirAll(configPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create config dir: %w", err)
		}
	}

	return configPath, nil
}
