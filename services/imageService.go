package services

import (
	"flxvwr/models"
	"flxvwr/utils"
	"time"

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

type ImageService struct {
	imagePaths map[string]fyne.URI
	current    int
	playlist   []string
	Zoomable   *models.ZoomableImage
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.imagePaths = make(map[string]fyne.URI)
	is.current = 0
	is.playlist = make([]string, 0)
	return is
}

func (is *ImageService) AddURI(uri fyne.URI) {
	ext := uri.Extension()
	if _, ok := AcceptedExtensions[ext]; ok {
		is.imagePaths[uri.Path()] = uri
		return
	}
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

func (is *ImageService) ImportImages(pos fyne.Position, uri []fyne.URI) {
	for _, u := range uri {
		if utils.IsDir(u) {
			paths := utils.RecurseDir(u)
			for _, p := range paths {
				is.AddURI(p)
			}
			continue
		}
		is.AddURI(u)
	}
	is.RecalculatePlaylist()
}

func (is *ImageService) GetCurrent() fyne.URI {
	if is.current < 0 {
		is.current = 0
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
	is.current = 0
	is.playlist = make([]string, 0)
}

func (is *ImageService) Update(w fyne.Window, ps *PlayerService) {
	w.SetContent(is.GetImageContainer(w, ps))
	ps.LastSet = time.Now()
}

func (is *ImageService) GetImageContainer(w fyne.Window, ps *PlayerService) *fyne.Container {
	image := canvas.NewImageFromURI(is.GetCurrent())
	image.FillMode = canvas.ImageFillContain
	zoomable := models.NewZoomableImage(image)
	is.Zoomable = zoomable
	imgContainer := container.NewWithoutLayout(zoomable.Image)
	imgContainer.Resize(fyne.NewSize(800, 600))
	zoomable.Image.Resize(imgContainer.Size()) // Resize the image to fit container initially
	imgContainer.Move(fyne.NewPos(0, 0))
	return container.NewStack(imgContainer)
}
