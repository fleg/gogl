package gogl

type (
	Renderer2D struct {
		canvas *Canvas
	}
)

// TODO: improve it, this is not optimized code, first thing that came to my mind
func (r *Renderer2D) Line(x1 int, y1 int, x2 int, y2 int, color uint32) {
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
			r.canvas.PutPixel(x1, y, color)
		}
		return
	}

	yp := y1
	for x := x1; x <= x2; x++ {
		yn := dy*x/dx - dy*x1/dx + y1
		for y := yp; y != yn+ydir; y += ydir {
			r.canvas.PutPixel(x, y, color)
		}
		yp = yn
	}
}

func (r *Renderer2D) FillRect(left int, bottom int, w int, h int, color uint32) {
	for y := bottom; y <= bottom+h; y++ {
		for x := left; x <= left+w; x++ {
			r.canvas.PutPixel(x, y, color)
		}
	}
}

func (r *Renderer2D) FillCircle(cx int, cy int, radius int, color uint32) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				r.canvas.PutPixel(x+cx, y+cy, color)
			}
		}
	}
}

func sign(a int) int {
	if a == 0 {
		return 0
	} else if a < 0 {
		return 1
	}

	return -1
}

func barycentricInt(x int, y int, x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) (int, int, int, int) {
	det := (x1-x3)*(y2-y3) - (x2-x3)*(y1-y3)
	u := (y2-y3)*(x-x3) + (x3-x2)*(y-y3)
	v := (y3-y1)*(x-x3) + (x1-x3)*(y-y3)
	w := det - u - v
	return u, v, w, det
}

func isInsideBarycentricInt(u int, v int, w int, det int) bool {
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

func (r *Renderer2D) FillTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, color uint32) {
	left, bottom, right, up := TriangleBBox(x1, y1, x2, y2, x3, y3, r.canvas.Width, r.canvas.Height)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			if isInsideBarycentricInt(barycentricInt(x, y, x1, y1, x2, y2, x3, y3)) {
				r.canvas.PutPixel(x, y, color)
			}
		}
	}
}

func (r *Renderer2D) FillTriangleRGB(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
	left, bottom, right, up := TriangleBBox(x1, y1, x2, y2, x3, y3, r.canvas.Width, r.canvas.Height)

	for y := bottom; y <= up; y++ {
		for x := left; x <= right; x++ {
			u, v, w, det := barycentricInt(x, y, x1, y1, x2, y2, x3, y3)
			if isInsideBarycentricInt(u, v, w, det) {
				r.canvas.PutPixel(x, y, MakeColor(
					255*u/det,
					255*v/det,
					255*w/det,
					0xFF,
				))
			}
		}
	}
}

func NewRenderer2D(c *Canvas) *Renderer2D {
	return &Renderer2D{
		canvas: c,
	}
}
