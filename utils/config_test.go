package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetConfigPath_Windows(t *testing.T) {
	getOS = func() string { return "windows" }
	mkdirAll = func(path string, perm os.FileMode) error { return nil }
	defer func() { getOS = func() string { return runtime.GOOS } }()
	defer func() { mkdirAll = func(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) } }()

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
	mkdirAll = func(path string, perm os.FileMode) error { return nil }
	getLinuxConfigDir = func() (string, error) { return "/home/testuser/.config", nil }
	defer func() { getLinuxConfigDir = func() (string, error) { return os.UserConfigDir() } }()
	defer func() { getOS = func() string { return runtime.GOOS } }()
	defer func() { mkdirAll = func(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) } }()

	testConfigPath := "/home/testuser/.config"

	expectedPath := filepath.Join(testConfigPath, "flxvwr")
	result, err := GetConfigPath()

	os.RemoveAll(result)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_Darwin(t *testing.T) {
	getOS = func() string { return "darwin" }
	mkdirAll = func(path string, perm os.FileMode) error { return nil }
	defer func() { getOS = func() string { return runtime.GOOS } }()
	defer func() { mkdirAll = func(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) } }()

	testHomePath := "/Users/testuser"
	os.Setenv("HOME", testHomePath)
	defer os.Unsetenv("HOME")

	expectedPath := filepath.Join(testHomePath, "Library", "Application Support", "flxvwr")
	result, err := GetConfigPath()

	os.RemoveAll(result)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, result)
	}
}

func TestGetConfigPath_UnsupportedOS(t *testing.T) {
	getOS = func() string { return "unsupportedOS" }
	mkdirAll = func(path string, perm os.FileMode) error { return nil }
	defer func() { getOS = func() string { return runtime.GOOS } }()
	defer func() { mkdirAll = func(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) } }()

	result, err := GetConfigPath()

	os.RemoveAll(result)
	if err == nil {
		t.Fatalf("Expected error for unsupported OS, got nil")
	}
	if result != "" {
		t.Errorf("Expected empty result for unsupported OS, got %v", result)
	}
}
