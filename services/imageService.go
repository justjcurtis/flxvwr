package services

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/justjcurtis/flxvwr/models"
	"github.com/justjcurtis/flxvwr/utils"
	"github.com/spf13/viper"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	imagePaths      map[string]fyne.URI
	knownPaths      map[string]fyne.URI
	imageStates     map[string]string
	current         int
	Zoomable        *models.ZoomableImage
	currentPlaylist int
	playlists       [][]string
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.imagePaths = make(map[string]fyne.URI)
	is.knownPaths = make(map[string]fyne.URI)
	is.imageStates = make(map[string]string)
	is.current = 0
	is.currentPlaylist = 0
	is.playlists = make([][]string, 10)
	for i := 0; i < 10; i++ {
		is.playlists[i] = make([]string, 0)
	}
	return is
}

func (is *ImageService) AddURI(uri fyne.URI) []fyne.URI {
	ext := uri.Extension()
	if !utils.IsFile(uri) {
		return nil
	}
	if _, ok := AcceptedExtensions[ext]; ok {
		is.imagePaths[uri.Path()] = uri
		is.playlists[is.currentPlaylist] = append(is.playlists[is.currentPlaylist], uri.Path())
		return nil
	} else if _, ok := PlaylistExtensions[ext]; ok {
		uris := utils.GetURIsFromLines(utils.ReadLines(uri))
		return uris
	}
	return nil
}

func (is *ImageService) GetCurrentPlaylist() []string {
	return is.playlists[is.currentPlaylist]
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
	is.RecalculateCurrentPlaylist()
}

func (is *ImageService) SetPlaylist(playlist int) {
	is.currentPlaylist = playlist
	is.current = 0
}

func (is *ImageService) RecalculateCurrentPlaylist() {
	currentPath := ""
	if len(is.playlists[is.currentPlaylist]) > 0 {
		currentPath = is.CurrentPath()
	}
	shouldShuffle := viper.GetBool("shuffle")
	if shouldShuffle {
		utils.Shuffle(is.playlists[is.currentPlaylist], 3)
	} else {
		utils.SortStrings(is.playlists[is.currentPlaylist])
	}
	if currentPath != "" {
		for i, p := range is.playlists[is.currentPlaylist] {
			if p == currentPath {
				is.current = i
				break
			}
		}
	}
}

func (is *ImageService) CurrentPath() string {
	return is.playlists[is.currentPlaylist][is.current]
}

func (is *ImageService) cacheZoomableImage() {
	if is.Zoomable != nil {
		if is.Zoomable.HasChanged() {
			is.imageStates[is.CurrentPath()] = is.Zoomable.ToString()
			return
		}
		is.imageStates[is.CurrentPath()] = ""
	}
}

func (is *ImageService) GetCurrent() fyne.URI {
	if is.current < 0 {
		is.current = 0
	}
	if len(is.GetCurrentPlaylist()) == 0 {
		return nil
	}
	if is.current >= len(is.GetCurrentPlaylist()) {
		is.current = len(is.GetCurrentPlaylist()) - 1
	}
	return is.imagePaths[is.CurrentPath()]
}

func (is *ImageService) Next() fyne.URI {
	if len(is.GetCurrentPlaylist()) == 0 {
		return nil
	}

	if is.Zoomable != nil {
		is.cacheZoomableImage()
	}
	is.current++
	if is.current >= len(is.GetCurrentPlaylist()) {
		is.current = 0
	}
	return is.GetCurrent()
}

func (is *ImageService) Previous() fyne.URI {
	if len(is.GetCurrentPlaylist()) == 0 {
		return nil
	}
	if is.Zoomable != nil {
		is.cacheZoomableImage()
	}

	is.current--
	if is.current < 0 {
		is.current = len(is.GetCurrentPlaylist()) - 1
	}
	return is.GetCurrent()
}

func (is *ImageService) Clear() {
	is.imagePaths = make(map[string]fyne.URI)
	is.knownPaths = make(map[string]fyne.URI)
	is.imageStates = make(map[string]string)
	is.current = 0
	is.currentPlaylist = 0
	for i := 0; i < 10; i++ {
		is.playlists[i] = make([]string, 0)
	}
}

func (is *ImageService) restoreZoomableState() {
	if is.Zoomable != nil {
		if is.imageStates[is.CurrentPath()] != "" {
			is.Zoomable.Set(is.imageStates[is.CurrentPath()])
			if is.Zoomable.Brightness != 1.0 || is.Zoomable.Contrast != 1.0 {
				db := (1.0 - is.Zoomable.Brightness) * -1
				dc := (1.0 - is.Zoomable.Contrast) * -1
				is.Zoomable.AdjustBrightnessAndContrast(db, dc)
			}
			if is.Zoomable.Rotation != 0 {
				fmt.Println(is.Zoomable.Rotation)
				is.Zoomable.Rotate(is.Zoomable.Rotation)
			}
			is.Zoomable.Set(is.imageStates[is.CurrentPath()])
			fmt.Println(is.imageStates[is.CurrentPath()])
			is.Zoomable.Refresh()
		} else {
			is.Zoomable.Reset()
		}
	}
}

func (is *ImageService) Update(w fyne.Window, ps *PlayerService, restartDelay bool) {
	if len(is.GetCurrentPlaylist()) == 0 {
		return
	}
	is.updateImageContainer(w)
	if restartDelay {
		ps.LastSet = time.Now()
		if is.Zoomable != nil {
			if is.imageStates[is.CurrentPath()] != "" {
				is.restoreZoomableState()
			} else {
				is.Zoomable.Reset()
			}
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
		log.Println(uri.Path(), err)
		is.Next()
		return is.GetImageFromURI(is.GetCurrent())
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(uri.Path(), err)
		is.Next()
		return is.GetImageFromURI(is.GetCurrent())
	}
	return img
}
