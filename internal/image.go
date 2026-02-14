package internal

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageManager struct{}

func NewImageManager() *ImageManager {
	return &ImageManager{}
}

func (i *ImageManager) ConcatPages(pages [][]byte) (*image.RGBA, error) {
	height := 0
	width := 0

	imgs := make([]image.Image, len(pages))

	for i, page := range pages {
		img, _, err := image.Decode(bytes.NewReader(page))
		if err != nil {
			return nil, err
		}

		bounds := img.Bounds()
		height += bounds.Dy()

		if bounds.Dx() > width {
			width = bounds.Dx()
		}

		imgs[i] = img
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	currentY := 0
	for _, img := range imgs {
		bounds := img.Bounds()
		rect := image.Rect(
			(width/2)-(bounds.Dx()/2),
			currentY,
			bounds.Dx()+(width/2),
			currentY+bounds.Dy(),
		)
		draw.Draw(dst, rect, img, bounds.Min, draw.Src)
		currentY += bounds.Dy()
	}

	return dst, nil
}

func (i *ImageManager) HasTransparency(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0xffff {
				return true
			}
		}
	}
	return false
}

func (i *ImageManager) SaveImageInSystem(content image.Image, path string) error {
	hasTransparency := false
	if i.HasTransparency(content) {
		path += ".png"
		hasTransparency = true
	} else {
		path += ".jpg"
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if hasTransparency {
		if err = png.Encode(f, content); err != nil {
			return err
		}
	} else {
		if err = jpeg.Encode(f, content, nil); err != nil {
			return err
		}
	}

	return nil
}
