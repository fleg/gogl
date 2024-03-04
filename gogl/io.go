package gogl

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func CanvasToPNG(c *Canvas, path string) error {
	img := image.NewNRGBA(image.Rect(0, 0, c.Width, c.Height))

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			p := c.GetPixel(x, y)
			img.Set(x, y, color.NRGBA{
				R: uint8((p >> 0) & 0xFF),
				G: uint8((p >> 8) & 0xFF),
				B: uint8((p >> 16) & 0xFF),
				A: uint8((p >> 24) & 0xFF),
			})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}

func CanvasFromPNG(path string) (*Canvas, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	c := NewCanvas(img.Bounds().Dx(), img.Bounds().Dy())

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c.PutPixel(x, y, MakeColor(int(r), int(g), int(b), int(a)))
		}
	}

	return c, nil
}
