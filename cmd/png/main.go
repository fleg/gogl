package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"gogl/gogl"
)

func CanvasToPNG(c *gogl.Canvas, path string) error {
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

func main() {
	const width = 800
	const height = 600
	const red = 0xFF2020FF
	const green = 0xFF20FF20
	const blue = 0xFFFF2020

	canvas := gogl.NewCanvas(width, height)

	canvas.Fill(0xFF202020)
	canvas.Line(0, 0, 100, 500, red)
	canvas.Line(0, 100, 500, 0, red)
	canvas.Line(0, 0, 500, 100, red)
	canvas.Line(0, 500, 100, 0, red)

	canvas.Line(0, 550, 750, 560, green)
	canvas.Line(0, 560, 750, 550, green)
	canvas.Line(-100, 590, width+100, 590, green)

	canvas.Line(width/2, 0, width/2, height-1, blue)
	canvas.Line(0, height/2, width-1, height/2, blue)

	canvas.FillRect(450, 500, 24, 24, blue)
	canvas.FillRect(500, 500, 24, 24, blue)
	canvas.FillRect(600, 500, 400, 24, blue)

	canvas.FillCircle(512, 480, 12, blue)
	canvas.FillCircle(450+12, 500+12, 12, red)

	// canvas.FillTriangle(width/2 + 50, height/2 - 50, width/2+width/4, 25, width - 50, height/4, green)
	canvas.FillRGBTriangle(width/2+50, height/2-50, width/2+width/4, 25, width-50, height/4)

	CanvasToPNG(canvas, "image.png")
}
