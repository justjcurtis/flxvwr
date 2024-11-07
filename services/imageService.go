package services

import (
	"image"
	"os"
	"time"

	"github.com/justjcurtis/flxvwr/models"
	"github.com/justjcurtis/flxvwr/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/spf13/viper"
)

var AcceptedExtensions map[string]bool = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

var PlaylistExtensions map[string]bool = map[string]bool{
	".txt": true,
	".m3u": true,
}

type ImageService struct {
	imagePaths map[string]fyne.URI
	knownPaths map[string]fyne.URI
	current    int
	playlist   []string
	Zoomable   *models.ZoomableImage
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.imagePaths = make(map[string]fyne.URI)
	is.knownPaths = make(map[string]fyne.URI)
	is.current = 0
	is.playlist = make([]string, 0)
	return is
}

func (is *ImageService) AddURI(uri fyne.URI) []fyne.URI {
	ext := uri.Extension()
	if !utils.IsFile(uri) {
		return nil
	}
	if _, ok := AcceptedExtensions[ext]; ok {
		is.imagePaths[uri.Path()] = uri
		return nil
	} else if _, ok := PlaylistExtensions[ext]; ok {
		uris := utils.GetURIsFromLines(utils.ReadLines(uri))
		return uris
	}
	return nil
}

func (is *ImageService) RecalculatePlaylist() {
	currentPath := ""
	if len(is.playlist) > 0 {
		currentPath = is.playlist[is.current]
	}
	playlist := make([]string, 0, len(is.imagePaths))
	index := 0
	for k := range is.imagePaths {
		if k == currentPath {
			is.current = index
		}
		playlist = append(playlist, k)
		index++
	}
	shouldShuffle := viper.GetBool("shuffle")
	if shouldShuffle {
		playlist = utils.Shuffle(playlist, 3)
	}
	is.playlist = playlist
}

func (is *ImageService) handleImport(u fyne.URI) []fyne.URI {
	overflow := []fyne.URI{}
	if is.knownPaths[u.Path()] != nil {
		return overflow
	}
	is.knownPaths[u.Path()] = u
	if utils.IsDir(u) {
		paths := utils.RecurseDir(u)
		for _, p := range paths {
			if is.knownPaths[p.Path()] != nil {
				continue
			}
			is.knownPaths[p.Path()] = p
			overflow = append(overflow, is.AddURI(p)...)
		}
		return overflow
	}

	if utils.IsFile(u) {
		overflow = append(overflow, is.AddURI(u)...)
	}
	return overflow
}

func (is *ImageService) ImportImages(uri []fyne.URI) {
	overflow := []fyne.URI{}
	for _, u := range uri {
		overflow = append(overflow, is.handleImport(u)...)
	}
	for len(overflow) > 0 {
		newOverflow := []fyne.URI{}
		for _, u := range overflow {
			newOverflow = append(newOverflow, is.handleImport(u)...)
		}
		overflow = newOverflow
	}

	is.RecalculatePlaylist()
}

func (is *ImageService) GetCurrent() fyne.URI {
	if is.current < 0 {
		is.current = 0
	}
	if len(is.playlist) == 0 {
		return nil
	}
	if is.current >= len(is.playlist) {
		is.current = len(is.playlist) - 1
	}
	return is.imagePaths[is.playlist[is.current]]
}

func (is *ImageService) Next() fyne.URI {
	is.current++
	if is.current >= len(is.playlist) {
		is.current = 0
	}
	return is.GetCurrent()
}

func (is *ImageService) Previous() fyne.URI {
	is.current--
	if is.current < 0 {
		is.current = len(is.playlist) - 1
	}
	return is.GetCurrent()
}

func (is *ImageService) Clear() {
	is.imagePaths = make(map[string]fyne.URI)
	is.knownPaths = make(map[string]fyne.URI)
	is.current = 0
	is.playlist = make([]string, 0)
}

func (is *ImageService) Update(w fyne.Window, ps *PlayerService, restartDelay bool) {
	if len(is.playlist) == 0 {
		return
	}
	is.updateImageContainer(w)
	if restartDelay {
		ps.LastSet = time.Now()
		if is.Zoomable != nil {
			is.Zoomable.Reset()
		}
	}
}

func (is *ImageService) updateImageContainer(w fyne.Window) {
	// Create a new image from the current URI
	img := is.GetImageFromURI(is.GetCurrent())
	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain

	// Create / Update the ZoomableImage
	if is.Zoomable != nil {
		is.Zoomable.Image = image
		is.Zoomable.Refresh()
	} else {
		is.Zoomable = models.NewZoomableImage(image, w)
	}

	// Create a new container for the image
	imgContainer := container.NewWithoutLayout(is.Zoomable.Image)

	// Resize the image container to the window size
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height
	imgContainer.Resize(fyne.NewSize(width, height))
	is.Zoomable.Image.Resize(imgContainer.Size())
	imgContainer.Move(fyne.NewPos(0, 0))

	// Set the content of the window to the image container
	result := container.NewStack(imgContainer)
	w.SetContent(result)

	// Update sizes and positions post w.SetContent
	actualSize := result.Size()
	imgContainer.Resize(actualSize)
	is.Zoomable.Image.Resize(actualSize)
	imgContainer.Move(fyne.NewPos(0, 0))
}

func (is *ImageService) GetImageFromURI(uri fyne.URI) image.Image {
	file, err := os.Open(uri.Path())
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}
