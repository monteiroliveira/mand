package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func createTestPNG(t *testing.T, width, height int, c color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test PNG: %v", err)
	}
	return buf.Bytes()
}

func TestNewImageManager(t *testing.T) {
	im := NewImageManager()
	if im == nil {
		t.Fatal("expected non-nil ImageManager")
	}
	if im.dpi != 94 {
		t.Errorf("expected dpi 94, got %d", im.dpi)
	}
}

func TestConcatPages_SinglePage(t *testing.T) {
	im := NewImageManager()
	page := createTestPNG(t, 100, 200, color.White)

	result, err := im.ConcatPages([][]byte{page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 200 {
		t.Errorf("expected 100x200, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestConcatPages_MultiplePages(t *testing.T) {
	im := NewImageManager()
	page1 := createTestPNG(t, 100, 200, color.White)
	page2 := createTestPNG(t, 100, 150, color.White)

	result, err := im.ConcatPages([][]byte{page1, page2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 {
		t.Errorf("expected width 100, got %d", bounds.Dx())
	}
	if bounds.Dy() != 350 {
		t.Errorf("expected height 350, got %d", bounds.Dy())
	}
}

func TestConcatPages_DifferentWidths(t *testing.T) {
	im := NewImageManager()
	page1 := createTestPNG(t, 80, 100, color.White)
	page2 := createTestPNG(t, 120, 100, color.White)

	result, err := im.ConcatPages([][]byte{page1, page2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 120 {
		t.Errorf("expected width 120 (max), got %d", bounds.Dx())
	}
	if bounds.Dy() != 200 {
		t.Errorf("expected height 200, got %d", bounds.Dy())
	}
}

func TestConcatPages_InvalidImage(t *testing.T) {
	im := NewImageManager()
	_, err := im.ConcatPages([][]byte{[]byte("not an image")})
	if err == nil {
		t.Fatal("expected error for invalid image data")
	}
}

func TestHasTransparency_Opaque(t *testing.T) {
	im := NewImageManager()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	if im.HasTransparency(img) {
		t.Error("expected no transparency for fully opaque image")
	}
}

func TestHasTransparency_WithAlpha(t *testing.T) {
	im := NewImageManager()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	img.Set(5, 5, color.RGBA{0, 0, 0, 128})
	if !im.HasTransparency(img) {
		t.Error("expected transparency for image with alpha pixel")
	}
}

func TestSaveImageInSystem(t *testing.T) {
	im := NewImageManager()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.White)
		}
	}

	tmpDir := t.TempDir()
	path := tmpDir + "/test_output"

	err := im.SaveImageInSystem(img, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path + ".png"); os.IsNotExist(err) {
		t.Error("expected PNG file to be created")
	}
}

func TestSavePdfInSystem(t *testing.T) {
	im := NewImageManager()
	page := createTestPNG(t, 100, 200, color.White)

	tmpDir := t.TempDir()
	path := tmpDir + "/test_chapter"

	err := im.SavePdfInSystem([][]byte{page}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path + ".pdf"); os.IsNotExist(err) {
		t.Error("expected PDF file to be created")
	}
}
