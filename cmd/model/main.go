package main

import (
	"gogl/gogl"
	"math"
)

// x and y are offsets on X and Y axis
func NewViewPort(x int, y int, w int, h int) *gogl.Matrix4f {
	vp := gogl.NewIdentityMatrix4[float64]()

		// {w/2, 0,   0, x+w/2},
		// {0,  -h/2, 0, y+h/2},
		// {0,   0,   1, 0},
		// {0,   0,   0, 1}
	vp.M[0][3] = float64(x + w/2)
	vp.M[1][3] = float64(y + h/2)
	vp.M[0][0] = float64(w/2)
	vp.M[1][1] = -float64(h/2) // invert Y axis to flip image horizonataly

	// vp.M[2][2] = float64(depth)/2
	// vp.M[2][3] = float64(depth)/2

	return vp
}

func LookAt(eye *gogl.Vec3f, center *gogl.Vec3f, up *gogl.Vec3f) *gogl.Matrix4f {
	minv := gogl.NewIdentityMatrix4[float64]()
	tr := gogl.NewIdentityMatrix4[float64]()

	z := eye.Subtract(center).Normalize()
	x := up.CrossProduct(z).Normalize()
	y := z.CrossProduct(x).Normalize()

	minv.M[0][0] = x.X
	minv.M[1][0] = y.X
	minv.M[2][0] = z.X

	minv.M[0][1] = x.Y
	minv.M[1][1] = y.Y
	minv.M[2][1] = z.Y

	minv.M[0][2] = x.Z
	minv.M[1][2] = y.Z
	minv.M[2][2] = z.Z

	tr.M[0][3] = -center.X
	tr.M[1][3] = -center.Y
	tr.M[2][3] = -center.Z

	return minv.Multiply(tr)
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

	head, err := gogl.NewModelFromFile("./assets/african_head.obj")
	if err != nil {
		panic(err)
	}
	if err := head.LoadTexture("./assets/african_head_diffuse.png"); err != nil {
		panic(err)
	}

	floor, err := gogl.NewModelFromFile("./assets/floor.obj")
	if err != nil {
		panic(err)
	}
	if err := floor.LoadTexture("./assets/floor_diffuse.png"); err != nil {
		panic(err)
	}


	models := []*gogl.Model{head, floor}

	light := gogl.Vec3f{X: 1.0, Y: 1.0, Z: 1.0}
	eye := gogl.Vec3f{X: 1.0, Y: 1.0, Z: 3.0}
	center := gogl.Vec3f{X: 0.0, Y: 0.0, Z: 0.0}
	up := gogl.Vec3f{X: 0.0, Y: 1.0, Z: 0.0}

	light.Normalize()

	projection := gogl.NewIdentityMatrix4[float64]()
	projection.M[3][2] = -1/eye.Subtract(&center).Length()

	viewport := NewViewPort(width / 8, height / 8, 3*width/4, 3*height/4)

	modelView := LookAt(&eye, &center, &up)

	zb := make([]float64, width*height)
	for i := 0; i < width*height; i++ {
		zb[i] = -math.MaxFloat64
	}

	for k := 0; k < len(models); k++ {
		m := models[k]

		for i := 0; i < len(m.Faces); i++ {
			face := m.Faces[i]
			screen := make([]gogl.Vec2i, 3)
			world := make([]gogl.Vec3f, 3)
			zs := make([]float64, 3)
			ws := make([]float64, 3)
			uvs := make([]gogl.Vec2f, 3)
			ns := make([]gogl.Vec3f, 3)

			for j := 0; j < 3; j++ {
				v := m.Verticies[face.Indicies[j]]

				s := viewport.Multiply(projection).Multiply(modelView).MultiplyVec4(v.ToVec4())

				world[j] = v
				screen[j].X = int(s.X/s.W+.5)
				screen[j].Y = int(s.Y/s.W+.5)

				zs[j] = s.Z/s.W
				ws[j] = s.W
				uvs[j] = m.UVs[face.TextureIndicies[j]]
				ns[j] = m.Normals[face.NormalIndicies[j]]
			}

			canvas.FillTriangleNUVZ(
				screen[0].X, screen[0].Y,
				screen[1].X, screen[1].Y,
				screen[2].X, screen[2].Y,
				zs[0], zs[1], zs[2],
				ws[0], ws[1], ws[2],
				zb,
				uvs[0].X, uvs[0].Y,
				uvs[1].X, uvs[1].Y,
				uvs[2].X, uvs[2].Y,
				m.Texture,
				ns[0].X, ns[0].Y, ns[0].Z,
				ns[1].X, ns[1].Y, ns[1].Z,
				ns[2].X, ns[2].Y, ns[2].Z,
				&light,
			)
		}
	}

	gogl.CanvasToPNG(canvas, "model.png")
}
