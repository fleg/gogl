package gogl

import "unsafe"

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

func (canvas *Canvas) FillRGBTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
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

func NewCanvas(w int, h int) *Canvas {
	return &Canvas{
		Width:  w,
		Height: h,
		Pixels: make([]uint32, w*h),
	}
}
