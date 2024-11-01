package utils

import (
	"errors"
	"testing"
)

func TestKillAppInstances(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	originalGetOS := getOS
	defer func() { getOS = originalGetOS }()
	tests := []struct {
		name      string
		appName   string
		os        string
		mockErr   error
		expectErr bool
	}{
		{"Kill on Linux", "testapp", "linux", nil, false},
		{"Kill on Windows", "testapp.exe", "windows", nil, false},
		{"Kill on MacOS", "testapp", "darwin", nil, false},
		{"Unsupported OS", "testapp", "unknownOS", nil, true},
		{"Command Failure", "testapp", "linux", errors.New("command failed"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOS = func() string { return tt.os }

			execCommand = func(command string, args ...string) error {
				return tt.mockErr
			}

			err := KillAppInstances(tt.appName)
			if (err != nil) != tt.expectErr {
				t.Errorf("KillAppInstances() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
