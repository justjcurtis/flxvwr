package services

import (
	"image"
	"image/color"
	"math"
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
	Brightness float64
	Contrast   float64
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.imagePaths = make(map[string]fyne.URI)
	is.knownPaths = make(map[string]fyne.URI)
	is.current = 0
	is.playlist = make([]string, 0)
	is.Brightness = 1.0
	is.Contrast = 1.0
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

	var oldScale float32 = 1.0
	var oldOffsetX float32 = 0.0
	var oldOffsetY float32 = 0.0

	if is.Zoomable != nil {
		oldScale = is.Zoomable.Scale
		oldOffsetX = is.Zoomable.OffsetX
		oldOffsetY = is.Zoomable.OffsetY
	}
	is.UpdateImageContainer(w, ps)
	is.Zoomable.Scale = oldScale
	is.Zoomable.OffsetX = oldOffsetX
	is.Zoomable.OffsetY = oldOffsetY
	is.Zoomable.Refresh()

	if restartDelay {
		ps.LastSet = time.Now()
	}
}

func (is *ImageService) UpdateImageContainer(w fyne.Window, ps *PlayerService) {
	img := is.GetImageFromURI(is.GetCurrent())
	adjusted := is.AdjustBrightnessAndContrast(img)
	image := canvas.NewImageFromImage(adjusted)
	image.FillMode = canvas.ImageFillContain
	zoomable := models.NewZoomableImage(image)
	is.Zoomable = zoomable
	imgContainer := container.NewWithoutLayout(is.Zoomable.Image)
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height
	imgContainer.Resize(fyne.NewSize(width, height))
	is.Zoomable.Image.Resize(imgContainer.Size())
	imgContainer.Move(fyne.NewPos(0, 0))
	result := container.NewStack(imgContainer)
	w.SetContent(result)
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

func (is *ImageService) AdjustBrightnessAndContrast(img image.Image) image.Image {
	bounds := img.Bounds()
	adjusted := image.NewRGBA(bounds)

	// Loop through each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			// Adjust brightness
			r := float64(originalColor.R) * is.Brightness
			g := float64(originalColor.G) * is.Brightness
			b := float64(originalColor.B) * is.Brightness

			// Adjust contrast
			r = (r-128)*is.Contrast + 128
			g = (g-128)*is.Contrast + 128
			b = (b-128)*is.Contrast + 128

			// Clamp color values to [0, 255] range
			adjustedColor := color.RGBA{
				R: uint8(math.Min(math.Max(r, 0), 255)),
				G: uint8(math.Min(math.Max(g, 0), 255)),
				B: uint8(math.Min(math.Max(b, 0), 255)),
				A: originalColor.A,
			}
			adjusted.Set(x, y, adjustedColor)
		}
	}
	return adjusted
}

func (is *ImageService) Rotate(w fyne.Window, direction int) {
	if is.Zoomable == nil {
		return
	}
	img := is.Zoomable.Image.Image
	bounds := img.Bounds()
	var rotated *image.RGBA

	switch direction % 4 {
	case 1:
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(bounds.Max.Y-y-1, x, img.At(x, y))
			}
		}
	case 2:
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
			}
		}
	case 3:
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(y, bounds.Max.X-x-1, img.At(x, y))
			}
		}
	default:
		return
	}

	is.Zoomable.Image.Image = rotated
	is.Zoomable.Reset()
	is.Zoomable.Image.Resize(fyne.NewSize(float32(rotated.Bounds().Max.X), float32(rotated.Bounds().Max.Y)))
	is.Zoomable.Refresh()
	imgContainer := container.NewWithoutLayout(is.Zoomable.Image)
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height
	imgContainer.Resize(fyne.NewSize(width, height))
	is.Zoomable.Image.Resize(imgContainer.Size())
	imgContainer.Move(fyne.NewPos(0, 0))
	container.NewStack(imgContainer)
	w.SetContent(imgContainer)
}
