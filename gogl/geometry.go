package gogl

func TriangleBBox(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, width int, height int) (int, int, int, int) {
	left := max(0, min(x1, min(x2, x3)))
	bottom := max(0, min(y1, min(y2, y3)))
	right := min(width-1, max(x1, max(x2, x3)))
	up := min(height-1, max(y1, max(y2, y3)))

	return left, bottom, right, up
}
