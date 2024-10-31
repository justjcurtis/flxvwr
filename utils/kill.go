package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func KillAppInstances(appName string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("pkill", appName)
	case "windows":
		cmd = exec.Command("taskkill", "/IM", appName+".exe", "/F")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
