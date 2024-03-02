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

func projectToScreen(x float64, y float64, w int, h int) (int, int) {
	return int((x + 1.0) / 2.0 * float64(w)), h - int((y+1.0)/2.0*float64(h))
}

func main() {
	const multiplier = 100
	const width = 12 * multiplier
	const height = 12 * multiplier
	const red = 0xFF2020FF
	const green = 0xFF20FF20
	const blue = 0xFFFF2020

	canvas := gogl.NewCanvas(width, height)

	canvas.Fill(0xFF000000)

	m, err := gogl.NewModelFromFile("./assets/african_head.obj")
	// m, err := gogl.NewModelFromFile("./assets/teapot.obj")
	if err != nil {
		panic(err)
	}

	light := gogl.Vec3f{X: 0.0, Y: 0.0, Z: -1.0}

	for i := 0; i < len(m.Faces); i++ {
		face := m.Faces[i]
		screen := make([]gogl.Vec2i, 3)
		world := make([]gogl.Vec3f, 3)

		for j := 0; j < 3; j++ {
			v := m.Verticies[face[j]]

			world[j] = v
			screen[j].X, screen[j].Y = projectToScreen(v.X, v.Y, width, height)
		}

		normal := world[2].Subtract(&world[0]).CrossProduct(world[1].Subtract(&world[0]))
		normal.Normalize()

		intensity := normal.DotProduct(&light)

		if intensity > 0 {
			canvas.FillTriangle(
				screen[0].X, screen[0].Y,
				screen[1].X, screen[1].Y,
				screen[2].X, screen[2].Y,
				gogl.MakeColor(int(intensity*255), int(intensity*255), int(intensity*255), 255),
			)
		}
	}

	CanvasToPNG(canvas, "model.png")
}
