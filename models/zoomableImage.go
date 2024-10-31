package models

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type ZoomableImage struct {
	Image     *canvas.Image
	prevScale float32
	Scale     float32
	OffsetX   float32
	OffsetY   float32
}

func NewZoomableImage(image *canvas.Image) *ZoomableImage {
	return &ZoomableImage{
		Image:     image,
		prevScale: 1.0,
		Scale:     1.0,
		OffsetX:   0,
		OffsetY:   0,
	}
}

func (z *ZoomableImage) Reset() {
	z.Scale = 1.0
	z.OffsetX, z.OffsetY = 0, 0
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

	imageCenterX := (float32(z.Image.Size().Width) / 2) - z.OffsetX
	imageCenterY := (float32(z.Image.Size().Height) / 2) - z.OffsetY

	z.OffsetX -= deltaScale * imageCenterX
	z.OffsetY -= deltaScale * imageCenterY

	z.Refresh()
}

func (z *ZoomableImage) Move(dx, dy float32) {
	z.OffsetX += dx
	z.OffsetY += dy
	z.Refresh()
}
