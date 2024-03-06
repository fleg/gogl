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

func NewCanvas(w int, h int) *Canvas {
	return &Canvas{
		Width:  w,
		Height: h,
		Pixels: make([]uint32, w*h),
	}
}
