package services

import (
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
	playlists       []map[string]fyne.URI
	knownPaths      map[string]fyne.URI
	imageStates     map[string]string
	currentIndex    int
	currentPlaylist int
	playlist        []string
	Zoomable        *models.ZoomableImage
	shuffle         bool
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.knownPaths = make(map[string]fyne.URI)
	is.imageStates = make(map[string]string)
	is.currentIndex = 0
	is.currentPlaylist = 0
	is.shuffle = viper.GetBool("shuffle")
	is.playlists = make([]map[string]fyne.URI, 10)
	for i := 0; i < 10; i++ {
		is.playlists[i] = make(map[string]fyne.URI)
	}
	return is
}

func (is *ImageService) HandleConfigUpdate(config models.Config) {
	if is.shuffle != config.Shuffle {
		is.shuffle = config.Shuffle
		is.RecalculateCurrentPlaylist()
	}
}

func (is *ImageService) AddURI(uri fyne.URI) []fyne.URI {
	ext := uri.Extension()
	if !utils.IsFile(uri) {
		return nil
	}
	if _, ok := AcceptedExtensions[ext]; ok {
		is.playlists[is.currentPlaylist][uri.Path()] = uri
		return nil
	} else if _, ok := PlaylistExtensions[ext]; ok {
		uris := utils.GetURIsFromLines(utils.ReadLines(uri))
		return uris
	}
	return nil
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
	is.currentIndex = 0
}

func (is *ImageService) RecalculateCurrentPlaylist() {
	currentPath := ""
	if len(is.playlist) > 0 {
		currentPath = is.playlist[is.currentIndex]
	}
	playlist := make([]string, len(is.playlists[is.currentPlaylist]))
	i := 0
	for _, p := range is.playlists[is.currentPlaylist] {
		playlist[i] = p.Path()
		i++
	}

	if is.shuffle {
		utils.Shuffle(playlist, 3)
	} else {
		utils.SortStrings(playlist)
	}
	if currentPath != "" {
		i = 0
		for _, p := range is.playlists[is.currentPlaylist] {
			if p.Path() == currentPath {
				is.currentIndex = i
				break
			}
			i++
		}
	}
	is.playlist = playlist
}

func (is *ImageService) cacheZoomableImage() {
	if is.Zoomable != nil {
		if is.Zoomable.HasChanged() {
			is.imageStates[is.playlist[is.currentIndex]] = is.Zoomable.ToString()
			return
		}
		is.imageStates[is.playlist[is.currentIndex]] = ""
	}
}

func (is *ImageService) GetCurrent() fyne.URI {
	if is.currentIndex < 0 {
		is.currentIndex = 0
	}
	if len(is.playlist) == 0 {
		return nil
	}
	if is.currentIndex >= len(is.playlist) {
		is.currentIndex = len(is.playlist) - 1
	}
	return is.playlists[is.currentPlaylist][is.playlist[is.currentIndex]]
}

func (is *ImageService) Next() fyne.URI {
	if len(is.playlist) == 0 {
		return nil
	}

	if is.Zoomable != nil {
		is.cacheZoomableImage()
	}

	is.currentIndex++
	if is.currentIndex >= len(is.playlist) {
		is.currentIndex = 0
	}
	return is.GetCurrent()
}

func (is *ImageService) Previous() fyne.URI {
	if len(is.playlist) == 0 {
		return nil
	}
	if is.Zoomable != nil {
		is.cacheZoomableImage()
	}

	is.currentIndex--
	if is.currentIndex < 0 {
		is.currentIndex = len(is.playlist) - 1
	}
	return is.GetCurrent()
}

func (is *ImageService) Clear() {
	is.knownPaths = make(map[string]fyne.URI)
	is.imageStates = make(map[string]string)
	is.currentIndex = 0
	is.currentPlaylist = 0
	is.Zoomable = nil
	for i := 0; i < 10; i++ {
		is.playlists[i] = make(map[string]fyne.URI)
	}
}

func (is *ImageService) AddCurrentToPlaylist(index int) {
	is.playlists[index][is.playlist[is.currentIndex]] = is.playlists[is.currentPlaylist][is.playlist[is.currentIndex]]
}

func (is *ImageService) RemoveCurrentFromPlaylist() {
	delete(is.playlists[is.currentPlaylist], is.playlist[is.currentIndex])
	is.playlist = append(is.playlist[:is.currentIndex], is.playlist[is.currentIndex+1:]...)
}

func (is *ImageService) SetCurrentPlaylist(index int) {
	is.currentPlaylist = index
	is.currentIndex = 0
	is.playlist = make([]string, 0)
	is.RecalculateCurrentPlaylist()
}

func (is *ImageService) restoreZoomableState() {
	if is.Zoomable != nil {
		if is.imageStates[is.playlist[is.currentIndex]] != "" {
			is.Zoomable.Set(is.imageStates[is.playlist[is.currentIndex]])
			if is.Zoomable.Brightness != 1.0 || is.Zoomable.Contrast != 1.0 {
				db := (1.0 - is.Zoomable.Brightness) * -1
				dc := (1.0 - is.Zoomable.Contrast) * -1
				is.Zoomable.AdjustBrightnessAndContrast(db, dc)
			}
			if is.Zoomable.Rotation != 0 {
				is.Zoomable.Rotate(is.Zoomable.Rotation)
			}
			is.Zoomable.Set(is.imageStates[is.playlist[is.currentIndex]])
			is.Zoomable.Refresh()
		} else {
			is.Zoomable.Reset()
		}
	}
}

func (is *ImageService) Update(w fyne.Window, ps *PlayerService, restartDelay bool) {
	if len(is.playlist) == 0 {
		return
	}
	is.updateImageContainer(w)
	if restartDelay {
		ps.LastSet = time.Now()
		if is.Zoomable != nil {
			if is.imageStates[is.playlist[is.currentIndex]] != "" {
				is.restoreZoomableState()
			} else {
				is.Zoomable.Reset()
			}
		}
	}
}

func (is *ImageService) updateImageContainer(w fyne.Window) {
	// Create a new image from the currentIndex URI
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
