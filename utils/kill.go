package utils

import (
	"fmt"
	"os/exec"
)

var execCommand = func(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func KillAppInstances(appName string) error {
	var err error
	currentOS := getOS()
	switch currentOS {
	case "linux", "darwin":
		err = execCommand("pkill", appName)
	case "windows":
		err = execCommand("taskkill", "/IM", appName+".exe", "/F")
	default:
		return fmt.Errorf("unsupported platform: %s", currentOS)
	}

	if err != nil {
		return err
	}

	return nil
}
