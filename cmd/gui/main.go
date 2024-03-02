package main

import (
	"gogl/gogl"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const multiplier int = 100
const width int = 16 * multiplier
const height int = 9 * multiplier

const red = 0xFF2020FF
const green = 0xFF20FF20
const blue = 0xFFFF2020
const gray = 0xFF202020

// TODO float32 math?
var dx float64 = 1
var dy float64 = 1
var step float64 = 5

var ad float64 = 1.0

var x1 float64 = float64(width / 2)
var y1 float64 = float64(height / 2)
var x2 float64 = float64(width/2 + width/4)
var y2 float64 = float64(0)
var x3 float64 = float64(width)
var y3 float64 = float64(height / 4)

func moveX(px float64) float64 {
	x := px + dx*step
	if x < 0 || x > float64(width) {
		dx *= -1
		return px
	}
	return x
}

func moveY(py float64) float64 {
	y := py + dy*step
	if y < 0 || y > float64(height) {
		dy *= -1
		return py
	}
	return y
}

func movePoint(px float64, py float64) (float64, float64) {
	return moveX(px), moveY(py)
}

func triangleCenter(x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64) (float64, float64) {
	return (x1 + x2 + x3) / 3.0, (y1 + y2 + y3) / 3.0
}

func rotatePoint(x float64, y float64, angle float64, cx float64, cy float64) (float64, float64) {
	dx := x - cx
	dy := y - cy

	cos := math.Cos(angle)
	sin := math.Sin(angle)
	xf := dx*cos - dy*sin + cx
	yf := dx*sin + dy*cos + cy

	if x < 0 || x > float64(width) {
		ad *= -1.0
		xf = x
	}
	if yf < 0 || yf > float64(height) {
		ad *= -1.0
		yf = y
	}

	return xf, yf
}

func render(c *gogl.Canvas, dt float64) {
	angle := 0.1 * ad * dt * math.Pi

	cx, cy := triangleCenter(x1, y1, x2, y2, x3, y3)
	x1, y1 = movePoint(x1, y1)
	x2, y2 = movePoint(x2, y2)
	x3, y3 = movePoint(x3, y3)

	x1, y1 = rotatePoint(x1, y1, angle, cx, cy)
	x2, y2 = rotatePoint(x2, y2, angle, cx, cy)
	x3, y3 = rotatePoint(x3, y3, angle, cx, cy)

	c.Fill(gray)

	c.FillTriangleRGB(int(x1), int(y1), int(x2), int(y2), int(x3), int(y3))
}

// TODO move to function to support multiple examples
func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"gogl",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(width),
		int32(height),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	canvas := gogl.NewCanvas(width, height)

	prevTicks := sdl.GetTicks()
	for {
		ticks := sdl.GetTicks()
		dt := float64(ticks-prevTicks) / 1000.0
		prevTicks = ticks

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
				break
			}
		}

		render(canvas, dt)

		texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
		if err != nil {
			panic(err)
		}

		pixels, _, err := texture.Lock(nil)
		copy(pixels, canvas.AsBytes())
		texture.Unlock()

		err = renderer.Copy(texture, nil, nil)
		if err != nil {
			panic(err)
		}

		renderer.Present()

		err = texture.Destroy()
		if err != nil {
			panic(err)
		}
	}
}
