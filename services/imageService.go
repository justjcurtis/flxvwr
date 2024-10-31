package services

import (
	"flxvwr/models"
	"flxvwr/utils"
	"image"
	"image/color"
	"math"
	"os"
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
	Brightness float64
	Contrast   float64
}

func NewImageService() *ImageService {
	is := &ImageService{}
	is.imagePaths = make(map[string]fyne.URI)
	is.current = 0
	is.playlist = make([]string, 0)
	is.Brightness = 1.0
	is.Contrast = 1.0
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
	w.SetContent(is.GetImageContainer(w, ps))
	is.Zoomable.Scale = oldScale
	is.Zoomable.OffsetX = oldOffsetX
	is.Zoomable.OffsetY = oldOffsetY
	is.Zoomable.Refresh()

	if restartDelay {
		ps.LastSet = time.Now()
	}
}

func (is *ImageService) GetImageContainer(w fyne.Window, ps *PlayerService) *fyne.Container {
	img := is.GetImageFromURI(is.GetCurrent())
	adjusted := is.AdjustBrightnessAndContrast(img)
	image := canvas.NewImageFromImage(adjusted)
	image.FillMode = canvas.ImageFillContain
	zoomable := models.NewZoomableImage(image)
	is.Zoomable = zoomable
	imgContainer := container.NewWithoutLayout(zoomable.Image)
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height
	imgContainer.Resize(fyne.NewSize(width, height))
	zoomable.Image.Resize(imgContainer.Size())
	imgContainer.Move(fyne.NewPos(0, 0))
	return container.NewStack(imgContainer)
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
