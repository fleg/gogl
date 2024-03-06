package gogl

func MakeColor(r int, g int, b int, a int) uint32 {
	return uint32(0xFF000000 | (b&0xFF)<<16 | (g&0xFF)<<8 | (r & 0xFF))
}

func SplitColor(c uint32) (int, int, int, int) {
	return int(c & 0xFF), int(c>>8) & 0xFF, int(c>>16) & 0xFF, int(c>>24) & 0xFF
}
