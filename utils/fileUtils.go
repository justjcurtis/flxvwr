package utils

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func IsDir(uri fyne.URI) bool {
	fileInfo, err := os.Stat(uri.Path())
	if err != nil {
		fmt.Println(err)
	}

	return fileInfo.IsDir()
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
