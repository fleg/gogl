package gogl

type (
	Renderer3D struct {
		canvas *Canvas
	}
)

func barycentric(x int, y int, x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) (float64, float64, float64) {
	det := (x1-x3)*(y2-y3) - (x2-x3)*(y1-y3)
	if det == 0 {
		return -1.0, -1.0, -1.0
	}

	u := float64((y2-y3)*(x-x3)+(x3-x2)*(y-y3)) / float64(det)
	v := float64((y3-y1)*(x-x3)+(x1-x3)*(y-y3)) / float64(det)
	w := 1.0 - u - v
	return u, v, w
}

func (renderer *Renderer3D) FillTriangle(
	x1 int, y1 int,
	x2 int, y2 int,
	x3 int, y3 int,
	z1 float64, z2 float64, z3 float64,
	w1 float64, w2 float64, w3 float64,
	zb []float64,
	u1 float64, v1 float64,
	u2 float64, v2 float64,
	u3 float64, v3 float64,
	texture *Canvas,
	n1x float64, n1y float64, n1z float64,
	n2x float64, n2y float64, n2z float64,
	n3x float64, n3y float64, n3z float64,
	light *Vec3f,
) {
	left, bottom, right, up := TriangleBBox(x1, y1, x2, y2, x3, y3, renderer.canvas.Width, renderer.canvas.Height)

	for x := left; x <= right; x++ {
		for y := bottom; y <= up; y++ {
			u, v, w := barycentric(x, y, x1, y1, x2, y2, x3, y3)
			if u >= 0.0 && v >= 0.0 && w >= 0.0 {
				u /= w1
				v /= w2
				w /= w3

				sum := u + v + w

				u /= sum
				v /= sum
				w /= sum

				// interpolate z value
				z := z1*u + z2*v + z3*w

				i := x + y*renderer.canvas.Width
				if zb[i] < z {
					zb[i] = z

					// interpolate texture coordinates
					tx := int((u1*u+u2*v+u3*w)*float64(texture.Width) + .5)
					ty := texture.Height - int((v1*u+v2*v+v3*w)*float64(texture.Height)+.5)

					// interpolate normals
					nx := n1x*u + n2x*v + n3x*w
					ny := n1y*u + n2y*v + n3y*w
					nz := n1z*u + n2z*v + n3z*w

					n := Vec3f{X: nx, Y: ny, Z: nz}

					intensity := n.Normalize().DotProduct(light)

					if intensity > 0 {
						color := texture.GetPixel(tx, ty)
						r, g, b, a := SplitColor(color)
						renderer.canvas.PutPixel(x, y, MakeColor(
							int(float64(r)*intensity),
							int(float64(g)*intensity),
							int(float64(b)*intensity),
							int(float64(a)*intensity),
						))
					}
				}
			}
		}
	}
}

func NewRenderer3D(c *Canvas) *Renderer3D {
	return &Renderer3D{
		canvas: c,
	}
}
