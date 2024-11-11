package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func Exists(uri fyne.URI) bool {
	_, err := os.Stat(uri.Path())
	if err != nil {
		return false
	}
	return true
}

func IsFile(uri fyne.URI) bool {
	fileInfo, err := os.Stat(uri.Path())
	if err != nil {
		log.Println(err)
	}
	return !fileInfo.IsDir() && Exists(uri)
}

func ReadLines(uri fyne.URI) []string {
	file, err := os.Open(uri.Path())
	if err != nil {
		log.Println("Failed to open file:", uri.Path(), err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", uri.Path(), err)
	}

	return lines
}

func GetURIFromLine(line string) (fyne.URI, error) {
	EXTRegex := `^#EXT`
	if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, EXTRegex) {
		return nil, fmt.Errorf("URI is a comment: %s", line)
	}
	if !strings.HasPrefix(line, "file://") {
		absPath, err := filepath.Abs(line)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %s %w", line, err)
		}
		line = "file://" + absPath
	}
	uri, err := storage.ParseURI(line)
	if err != nil {
		return nil, err
	}
	if !Exists(uri) {
		return nil, fmt.Errorf("URI does not exist: %s", uri)
	}
	return uri, nil
}

func GetURIsFromLines(lines []string) []fyne.URI {
	var uris []fyne.URI
	for _, line := range lines {
		uri, err := GetURIFromLine(line)
		if err != nil {
			if !strings.Contains(err.Error(), "#EXT") {
				log.Println(err)
			}
			continue
		}
		uris = append(uris, uri)
	}
	return uris
}

func IsDir(uri fyne.URI) bool {
	fileInfo, err := os.Stat(uri.Path())
	if err != nil {
		return false
	}

	return fileInfo.IsDir() && Exists(uri)
}

func GetChildren(uri fyne.URI) []fyne.URI {
	lister, err := storage.ListerForURI(uri)
	if err != nil {
		return nil
	}
	paths, err := lister.List()
	if err != nil {
		return nil
	}
	return paths
}

func RecurseDir(uri fyne.URI) []fyne.URI {
	if IsDir(uri) {
		uriPaths := GetChildren(uri)
		var paths []fyne.URI
		for _, u := range uriPaths {
			paths = append(paths, RecurseDir(u)...)
		}
		return paths
	} else {
		return []fyne.URI{uri}
	}
}
