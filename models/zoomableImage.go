package models

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type ZoomableImage struct {
	window     fyne.Window
	Image      *canvas.Image
	prevScale  float32
	Scale      float32
	OffsetX    float32
	OffsetY    float32
	Rotation   int
	Brightness float32
	Contrast   float32
}

func NewZoomableImage(image *canvas.Image, w fyne.Window) *ZoomableImage {
	return &ZoomableImage{
		window:     w,
		Image:      image,
		prevScale:  1.0,
		Scale:      1.0,
		OffsetX:    0.0,
		OffsetY:    0.0,
		Brightness: 1.0,
		Contrast:   1.0,
		Rotation:   0,
	}
}

func (z *ZoomableImage) HasChanged() bool {
	return z.prevScale != z.Scale || z.OffsetX != 0.0 || z.OffsetY != 0.0 || z.Brightness != 1.0 || z.Contrast != 1.0 || z.Rotation != 0
}

func (z *ZoomableImage) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%d", z.prevScale, z.Scale, z.OffsetX, z.OffsetY, z.Brightness, z.Contrast, z.Rotation)
}

func FromString(str string, image *canvas.Image, w fyne.Window) *ZoomableImage {
	var prevScale, scale, offsetX, offsetY float32
	var brightness, contrast float64
	var rotation int
	fmt.Sscanf(str, "%f,%f,%f,%f,%f,%f,%d", &prevScale, &scale, &offsetX, &offsetY, &brightness, &contrast, &rotation)
	return &ZoomableImage{
		window:    w,
		Image:     image,
		prevScale: prevScale,
		Scale:     scale,
		OffsetX:   offsetX,
		OffsetY:   offsetY,
		Rotation:  rotation,
	}
}

func (z *ZoomableImage) Set(str string) {
	fmt.Println(str)
	fmt.Sscanf(str, "%f,%f,%f,%f,%f,%f,%d", &z.prevScale, &z.Scale, &z.OffsetX, &z.OffsetY, &z.Brightness, &z.Contrast, &z.Rotation)
	z.Refresh()
}

func (z *ZoomableImage) Reset() {
	z.ResetBrightnessContrast()
	z.ResetZoomAndPan()
	z.ResetRotation()
	z.Refresh()
}

func (z *ZoomableImage) ResetRotation() {
	z.Rotation = 0
	z.Refresh()
}

func (z *ZoomableImage) ResetBrightnessContrast() {
	z.Brightness = 1.0
	z.Contrast = 1.0
	z.Refresh()
}

func (z *ZoomableImage) ResetZoomAndPan() {
	z.Scale = 1.0
	z.OffsetX, z.OffsetY = 0.0, 0.0
	z.prevScale = 1.0
	z.Refresh()
}

func (z *ZoomableImage) Refresh() {
	newWidth := float32(z.Image.Size().Width) * (1 / z.prevScale) * z.Scale
	newHeight := float32(z.Image.Size().Height) * (1 / z.prevScale) * z.Scale
	z.prevScale = z.Scale
	z.Image.Resize(fyne.NewSize(newWidth, newHeight))
	z.Image.Move(fyne.NewPos(z.OffsetX, z.OffsetY))
	z.Image.Refresh()
}

func (z *ZoomableImage) Zoom(dz float32) {
	var zoomFactor float32 = 0.1
	prevScale := z.Scale
	z.Scale += zoomFactor * dz
	if dz < 0 {
		if z.Scale < 0.1 {
			z.Scale = 0.1
		}
	}
	postScale := z.Scale
	deltaScale := postScale - prevScale

	imageCenterX := (float32(z.Image.Size().Width) / 2) + z.OffsetX
	imageCenterY := (float32(z.Image.Size().Height) / 2) + z.OffsetY

	z.OffsetX -= deltaScale * imageCenterX
	z.OffsetY -= deltaScale * imageCenterY

	z.Refresh()
}

func (z *ZoomableImage) Move(dx, dy float32) {
	z.OffsetX += dx
	z.OffsetY += dy
	z.Refresh()
}

func (z *ZoomableImage) Rotate(direction int) {
	z.Rotation += direction
	z.Rotation = (z.Rotation + 4) % 4
	fmt.Println(z.Rotation)
	img := z.Image.Image
	bounds := img.Bounds()
	var rotated *image.RGBA

	switch direction % 4 {
	case 1:
		// Rotate 90 degrees clockwise
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(bounds.Max.Y-y-1, x, img.At(x, y))
			}
		}
	case 2:
		// Rotate 180 degrees
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
			}
		}
	case 3:
		// Rotate 90 degrees counter-clockwise
		rotated = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rotated.Set(y, bounds.Max.X-x-1, img.At(x, y))
			}
		}
	default:
		return
	}

	z.Image.Image = rotated
	z.ResetZoomAndPan()
	z.Image.Resize(fyne.NewSize(float32(rotated.Bounds().Max.X), float32(rotated.Bounds().Max.Y)))
	z.Refresh()
	imgContainer := container.NewWithoutLayout(z.Image)
	width := z.window.Canvas().Size().Width
	height := z.window.Canvas().Size().Height
	imgContainer.Resize(fyne.NewSize(width, height))
	z.Image.Resize(imgContainer.Size())
	imgContainer.Move(fyne.NewPos(0, 0))
	result := container.NewStack(imgContainer)
	z.window.SetContent(result)
	actualSize := result.Size()
	imgContainer.Resize(actualSize)
	z.Image.Resize(actualSize)
	imgContainer.Move(fyne.NewPos(0, 0))

}

func (z *ZoomableImage) AdjustBrightnessAndContrast(db, dc float32) {
	bounds := z.Image.Image.Bounds()
	adjusted := image.NewRGBA(bounds)

	z.Brightness = z.Brightness + db
	z.Contrast = z.Contrast + dc

	// Loop through each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := color.RGBAModel.Convert(z.Image.Image.At(x, y)).(color.RGBA)

			// Adjust brightness
			r := float32(originalColor.R) * (1 + db)
			g := float32(originalColor.G) * (1 + db)
			b := float32(originalColor.B) * (1 + db)

			// Adjust contrast
			r = (r-128)*(1+dc) + 128
			g = (g-128)*(1+dc) + 128
			b = (b-128)*(1+dc) + 128

			// Clamp color values to [0, 255] range
			adjustedColor := color.RGBA{
				R: uint8(math.Min(math.Max(float64(r), 0), 255)),
				G: uint8(math.Min(math.Max(float64(g), 0), 255)),
				B: uint8(math.Min(math.Max(float64(b), 0), 255)),
				A: originalColor.A,
			}
			adjusted.Set(x, y, adjustedColor)
		}
	}
	z.Image.Image = adjusted
	z.Refresh()
}
