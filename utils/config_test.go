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

func TestGetConfigPath_Linux_XDGConfigHome(t *testing.T) {
	getOS = func() string { return "linux" }
	defer func() { getOS = func() string { return runtime.GOOS } }()

	testXDGConfigPath := "/home/testuser/.config"
	os.Setenv("XDG_CONFIG_HOME", testXDGConfigPath)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	expectedPath := filepath.Join(testXDGConfigPath, "flxvwr")
	result, err := GetConfigPath()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_Linux_Default(t *testing.T) {
	getOS = func() string { return "linux" }
	defer func() { getOS = func() string { return runtime.GOOS } }()

	os.Unsetenv("XDG_CONFIG_HOME")

	expectedPath := "/etc/flxvwr"
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
