package main

import (
	"gogl/gogl"
)

func main() {
	const multiplier = 100
	const width = 16 * multiplier
	const height = 9 * multiplier
	const red = 0xFF2020FF
	const green = 0xFF20FF20
	const blue = 0xFFFF2020

	canvas := gogl.NewCanvas(width, height)
	renderer := gogl.NewRenderer2D(canvas)

	canvas.Fill(0xFF202020)
	renderer.Line(0, 0, 100, 500, red)
	renderer.Line(0, 100, 500, 0, red)
	renderer.Line(0, 0, 500, 100, red)
	renderer.Line(0, 500, 100, 0, red)

	renderer.Line(0, 550, 750, 560, green)
	renderer.Line(0, 560, 750, 550, green)
	renderer.Line(-100, 590, width+100, 590, green)

	renderer.Line(width/2, 0, width/2, height-1, blue)
	renderer.Line(0, height/2, width-1, height/2, blue)

	renderer.FillRect(450, 500, 24, 24, blue)
	renderer.FillRect(500, 500, 24, 24, blue)
	renderer.FillRect(600, 500, 400, 24, blue)

	renderer.FillCircle(512, 480, 12, blue)
	renderer.FillCircle(450+12, 500+12, 12, red)

	// canvas.FillTriangle(width/2 + 50, height/2 - 50, width/2+width/4, 25, width - 50, height/4, green)
	renderer.FillTriangleRGB(width/2+50, height/2-50, width/2+width/4, 25, width-50, height/4)

	gogl.CanvasToPNG(canvas, "image.png")
}
