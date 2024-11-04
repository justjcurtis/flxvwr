package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetConfigPath_Windows(t *testing.T) {
	getOS = func() string { return "windows" }
	defer func() { getOS = func() string { return runtime.GOOS } }()

	testAppDataPath := "C:\\Users\\TestUser\\AppData\\Roaming"
	os.Setenv("APPDATA", testAppDataPath)
	defer os.Unsetenv("APPDATA")

	expectedPath := filepath.Join(testAppDataPath, "flxvwr")
	result, err := GetConfigPath()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_Linux(t *testing.T) {
	getOS = func() string { return "linux" }
	getLinuxConfigDir = func() (string, error) { return "/home/testuser/.config", nil }
	defer func() { getLinuxConfigDir = func() (string, error) { return os.UserConfigDir() } }()
	defer func() { getOS = func() string { return runtime.GOOS } }()

	testConfigPath := "/home/testuser/.config"

	expectedPath := filepath.Join(testConfigPath, "flxvwr")
	result, err := GetConfigPath()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_Darwin(t *testing.T) {
	getOS = func() string { return "darwin" }
	defer func() { getOS = func() string { return runtime.GOOS } }()

	testHomePath := "/Users/testuser"
	os.Setenv("HOME", testHomePath)
	defer os.Unsetenv("HOME")

	expectedPath := filepath.Join(testHomePath, "Library", "Application Support", "flxvwr")
	result, err := GetConfigPath()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_UnsupportedOS(t *testing.T) {
	getOS = func() string { return "unsupportedOS" }
	defer func() { getOS = func() string { return runtime.GOOS } }()

	result, err := GetConfigPath()

	if err == nil {
		t.Fatalf("Expected error for unsupported OS, got nil")
	}
	if result != "" {
		t.Errorf("Expected empty result for unsupported OS, got %v", result)
	}
}
