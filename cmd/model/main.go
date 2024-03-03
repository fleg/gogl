package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
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

func CanvasFromPNG(path string) (*gogl.Canvas, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	c := gogl.NewCanvas(img.Bounds().Dx(), img.Bounds().Dy())

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c.PutPixel(x, y, gogl.MakeColor(int(r), int(g), int(b), int(a)))
		}
	}

	return c, nil
}

func projectToScreen(x float32, y float32, w int, h int) (int, int) {
	return int((x + 1.0) / 2.0 * float32(w)), h - int((y+1.0)/2.0*float32(h))
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

	texture, err := CanvasFromPNG("./assets/african_head_diffuse.png")
	if err != nil {
		panic(err)
	}

	light := gogl.Vec3f{X: 0.0, Y: 0.0, Z: 1.0}
	zb := make([]float32, width*height)
	for i := 0; i < width*height; i++ {
		zb[i] = -math.MaxFloat32
	}

	for i := 0; i < len(m.Faces); i++ {
		face := m.Faces[i]
		screen := make([]gogl.Vec2i, 3)
		world := make([]gogl.Vec3f, 3)
		zs := make([]float32, 3)
		uvs := make([]gogl.Vec2f, 3)
		ns := make([]gogl.Vec3f, 3)

		for j := 0; j < 3; j++ {
			v := m.Verticies[face.Indicies[j]]

			world[j] = v
			screen[j].X, screen[j].Y = projectToScreen(v.X, v.Y, width, height)
			zs[j] = v.Z
			uvs[j] = m.UVs[face.TextureIndicies[j]]
			ns[j] = m.Normals[face.NormalIndicies[j]]
		}

		normal := world[1].Subtract(&world[0]).CrossProduct(world[2].Subtract(&world[0]))
		normal.Normalize()

		intensity := normal.DotProduct(&light)

		if intensity > 0.0 {
			canvas.FillTriangleNUVZ(
				screen[0].X, screen[0].Y,
				screen[1].X, screen[1].Y,
				screen[2].X, screen[2].Y,
				zs[0], zs[1], zs[2],
				zb,
				uvs[0].X, uvs[0].Y,
				uvs[1].X, uvs[1].Y,
				uvs[2].X, uvs[2].Y,
				texture,
				ns[0].X, ns[0].Y, ns[0].Z,
				ns[1].X, ns[1].Y, ns[1].Z,
				ns[2].X, ns[2].Y, ns[2].Z,
				&light,
			)
		}
	}

	CanvasToPNG(canvas, "model.png")
}
