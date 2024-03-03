package gogl

import (
	"unsafe"
)

type (
	Canvas struct {
		Width  int
		Height int

		Pixels []uint32
	}
)

func (canvas *Canvas) GetPixel(x int, y int) uint32 {
	if x >= canvas.Width || y >= canvas.Height || x < 0 || y < 0 {
		return 0
	}

	return canvas.Pixels[y*canvas.Width+x]
}

func (canvas *Canvas) PutPixel(x int, y int, color uint32) {
	if x >= canvas.Width || y >= canvas.Height || x < 0 || y < 0 {
		return
	}

	canvas.Pixels[y*canvas.Width+x] = color
}

func (canvas *Canvas) AsBytes() []byte {
	l := len(canvas.Pixels)

	if l == 0 {
		return nil
	}

	size := unsafe.Sizeof(canvas.Pixels[0])

	return unsafe.Slice((*byte)(unsafe.Pointer(&canvas.Pixels[0])), int(size)*l)
}

func (canvas *Canvas) Fill(color uint32) {
	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			canvas.PutPixel(x, y, color)
		}
	}
}

// TODO: improve it, this is not optimized code, first thing that came to my mind
func (canvas *Canvas) Line(x1 int, y1 int, x2 int, y2 int, color uint32) {
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	dx := x2 - x1
	dy := y2 - y1

	var ydir int
	if y2 > y1 {
		ydir = 1
	} else {
		ydir = -1
	}

	if dx == 0 {
		for y := y1; y != y2+ydir; y += ydir {
			canvas.PutPixel(x1, y, color)
		}
		return
	}

	yp := y1
	for x := x1; x <= x2; x++ {
		yn := dy*x/dx - dy*x1/dx + y1
		for y := yp; y != yn+ydir; y += ydir {
			canvas.PutPixel(x, y, color)
		}
		yp = yn
	}
}

func (canvas *Canvas) FillRect(left int, bottom int, w int, h int, color uint32) {
	for y := bottom; y <= bottom+h; y++ {
		for x := left; x <= left+w; x++ {
			canvas.PutPixel(x, y, color)
		}
	}
}

func (canvas *Canvas) FillCircle(cx int, cy int, r int, color uint32) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				canvas.PutPixel(x+cx, y+cy, color)
			}
		}
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func sign(a int) int {
	if a == 0 {
		return 0
	} else if a < 0 {
		return 1
	}

	return -1
}

func barycentric(x int, y int, x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) (int, int, int, int) {
	det := (x1-x3)*(y2-y3) - (x2-x3)*(y1-y3)
	u := (y2-y3)*(x-x3) + (x3-x2)*(y-y3)
	v := (y3-y1)*(x-x3) + (x1-x3)*(y-y3)
	w := det - u - v
	return u, v, w, det
}

func isBarycentricInside(u int, v int, w int, det int) bool {
	if det == 0 {
		return false
	}

	dets := sign(det)
	us := sign(u) * dets
	vs := sign(v) * dets
	ws := sign(w) * dets

	if us >= 0 && vs >= 0 && ws >= 0 {
		return true
	}

	return false
}

func (canvas *Canvas) triangleBbox(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) (int, int, int, int) {
	left := max(0, min(x1, min(x2, x3)))
	bottom := max(0, min(y1, min(y2, y3)))
	right := min(canvas.Width-1, max(x1, max(x2, x3)))
	up := min(canvas.Height-1, max(y1, max(y2, y3)))

	return left, bottom, right, up
}

func MakeColor(r int, g int, b int, a int) uint32 {
	return uint32(0xFF000000 | (b&0xFF)<<16 | (g&0xFF)<<8 | (r & 0xFF))
}

func SplitColor(c uint32) (int, int, int, int) {
	return int(c & 0xFF), int(c>>8) & 0xFF, int(c>>16) & 0xFF, int(c>>24) & 0xFF
}

func (canvas *Canvas) FillTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, color uint32) {
	// canvas.Line(x1, y1, x2, y2, 0xFF2020FF);
	// canvas.Line(x3, y3, x2, y2, 0xFF2020FF);
	// canvas.Line(x1, y1, x3, y3, 0xFF2020FF);

	left, bottom, right, up := canvas.triangleBbox(x1, y1, x2, y2, x3, y3)

	// canvas.Line(left, bottom, left + w, bottom, 0xFF2020FF);
	// canvas.Line(left, bottom, left, bottom + h, 0xFF2020FF);
	// canvas.Line(left + w, bottom, left + w, bottom + h, 0xFF2020FF);
	// canvas.Line(left, bottom + h, left + w, bottom + h, 0xFF2020FF);

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			if isBarycentricInside(barycentric(x, y, x1, y1, x2, y2, x3, y3)) {
				canvas.PutPixel(x, y, color)
			}
		}
	}
}

func (canvas *Canvas) FillTriangleRGB(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
	left, bottom, right, up := canvas.triangleBbox(x1, y1, x2, y2, x3, y3)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			u, v, w, det := barycentric(x, y, x1, y1, x2, y2, x3, y3)
			if isBarycentricInside(u, v, w, det) {
				canvas.PutPixel(x, y, MakeColor(
					255*u/det,
					255*v/det,
					255*w/det,
					0xFF,
				))
			}
		}
	}
}

func (canvas *Canvas) FillTriangleZ(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, color uint32, z1 float64, z2 float64, z3 float64, zb []float64) {
	left, bottom, right, up := canvas.triangleBbox(x1, y1, x2, y2, x3, y3)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			u, v, w, det := barycentric(x, y, x1, y1, x2, y2, x3, y3)
			if isBarycentricInside(u, v, w, det) {
				// interpolate z value
				z := z1*float64(u)/float64(det) + z2*float64(v)/float64(det) + z3*float64(w)/float64(det)

				i := x + y*canvas.Width
				if zb[i] < z {
					zb[i] = z
					canvas.PutPixel(x, y, color)
				}
			}
		}
	}
}

func (canvas *Canvas) FillTriangleUVZ(
	x1 int, y1 int,
	x2 int, y2 int,
	x3 int, y3 int,
	z1 float64, z2 float64, z3 float64,
	zb []float64,
	u1 float64, v1 float64,
	u2 float64, v2 float64,
	u3 float64, v3 float64,
	intensity float64,
	texture *Canvas,
) {
	left, bottom, right, up := canvas.triangleBbox(x1, y1, x2, y2, x3, y3)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			u, v, w, det := barycentric(x, y, x1, y1, x2, y2, x3, y3)
			if isBarycentricInside(u, v, w, det) {
				// interpolate z value
				z := z1*float64(u)/float64(det) + z2*float64(v)/float64(det) + z3*float64(w)/float64(det)

				i := x + y*canvas.Width
				if zb[i] < z {
					zb[i] = z

					// interpolate texture coordinates
					tx := int((u1*float64(u)/float64(det)+u2*float64(v)/float64(det)+u3*float64(w)/float64(det))*float64(texture.Width) + .5)
					ty := texture.Height - int((v1*float64(u)/float64(det)+v2*float64(v)/float64(det)+v3*float64(w)/float64(det))*float64(texture.Height)+.5)

					color := texture.GetPixel(tx, ty)
					r, g, b, _ := SplitColor(color)
					canvas.PutPixel(x, y, MakeColor(
						int(float64(r)*intensity),
						int(float64(g)*intensity),
						int(float64(b)*intensity),
						255,
					))
				}
			}
		}
	}
}

func (canvas *Canvas) FillTriangleNUVZ(
	x1 int, y1 int,
	x2 int, y2 int,
	x3 int, y3 int,
	z1 float64, z2 float64, z3 float64,
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
	left, bottom, right, up := canvas.triangleBbox(x1, y1, x2, y2, x3, y3)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			u, v, w, det := barycentric(x, y, x1, y1, x2, y2, x3, y3)
			if isBarycentricInside(u, v, w, det) {
				uf := float64(u)/float64(det)
				vf := float64(v)/float64(det)
				wf := float64(w)/float64(det)

				// interpolate z value
				z := z1*uf + z2*vf + z3*wf

				i := x + y*canvas.Width
				if zb[i] < z {
					zb[i] = z

					// interpolate texture coordinates
					tx := int((u1*uf+u2*vf+u3*wf)*float64(texture.Width) + .5)
					ty := texture.Height - int((v1*uf+v2*vf+v3*wf)*float64(texture.Height)+.5)

					// interpolate normals
					nx := n1x*uf + n2x*vf + n3x*wf
					ny := n1y*uf + n2y*vf + n3y*wf
					nz := n1z*uf + n2z*vf + n3z*wf

					n := Vec3f{X: nx, Y: ny, Z: nz}
					n.Normalize()

					intensity := n.DotProduct(light)

					if intensity > 0 {
						color := texture.GetPixel(tx, ty)
						r, g, b, a := SplitColor(color)
						canvas.PutPixel(x, y, MakeColor(
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

func NewCanvas(w int, h int) *Canvas {
	return &Canvas{
		Width:  w,
		Height: h,
		Pixels: make([]uint32, w*h),
	}
}
