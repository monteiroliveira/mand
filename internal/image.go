package internal

import (
	"bytes"
	"fmt"
	_ "golang.org/x/image/webp" // Load webp encoding
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/jung-kurt/gofpdf"
)

type ImageManager struct {
	dpi int
}

func NewImageManager() *ImageManager {
	return &ImageManager{
		dpi: 94,
	}
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
	path = path + ".png"
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	if err = encoder.Encode(f, content); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// TODO: Need to enhance this
func (i *ImageManager) SavePdfInSystem(pages [][]byte, chapterName string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	tmp_file := []string{}

	defer func() error {
		for _, path := range tmp_file {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	}()

	for idx, page := range pages {
		path := fmt.Sprintf("%s-%d.jpeg", chapterName, idx)

		img, _, err := image.Decode(bytes.NewReader(page))
		if err != nil {
			return err
		}

		bounds := img.Bounds()
		widthPx := bounds.Dx()
		heightPx := bounds.Dy()

		width := float64(widthPx) * (25.4 / float64(i.dpi))
		height := float64(heightPx) * (25.4 / float64(i.dpi))

		// Set custom page size
		pdf.AddPageFormat("P", gofpdf.SizeType{
			Wd: width,
			Ht: height,
		})

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if err = jpeg.Encode(f, img, &jpeg.Options{Quality: 95}); err != nil {
			return err
		}
		tmp_file = append(tmp_file, path)

		pdf.Image(path, 0, 0, width, height, false, "", 0, "")
	}

	err := pdf.OutputFileAndClose(chapterName + ".pdf")
	if err != nil {
		fmt.Println("Fail to create pdf")
	}

	return nil
}
