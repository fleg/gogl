package gogl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	Face struct {
		Indicies        []int
		TextureIndicies []int
		NormalIndicies  []int
	}

	Model struct {
		Verticies []Vec3f
		Faces     []Face
		UVs       []Vec2f
		Normals   []Vec3f
	}
)

func NewModelFromFile(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := &Model{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "v ") {
			v := Vec3f{}
			fmt.Sscanf(line, "v %f %f %f", &v.X, &v.Y, &v.Z)
			m.Verticies = append(m.Verticies, v)
		} else if strings.HasPrefix(line, "f ") {
			f := Face{
				Indicies:        make([]int, 3),
				TextureIndicies: make([]int, 3),
				NormalIndicies:  make([]int, 3),
			}
			if strings.Contains(line, "/") {
				fmt.Sscanf(
					line,
					"f %d/%d/%d %d/%d/%d %d/%d/%d",
					&f.Indicies[0], &f.TextureIndicies[0], &f.NormalIndicies[0],
					&f.Indicies[1], &f.TextureIndicies[1], &f.NormalIndicies[1],
					&f.Indicies[2], &f.TextureIndicies[2], &f.NormalIndicies[2],
				)
			} else {
				fmt.Sscanf(line, "f %d %d %d", &f.Indicies[0], &f.Indicies[1], &f.Indicies[2])
			}

			for i := 0; i < 3; i++ {
				f.Indicies[i] -= 1
				f.TextureIndicies[i] -= 1
				f.NormalIndicies[i] -= 1
			}

			m.Faces = append(m.Faces, f)
		} else if strings.HasPrefix(line, "vt ") {
			uv := Vec2f{}
			fmt.Sscanf(line, "vt %f %f", &uv.X, &uv.Y)
			m.UVs = append(m.UVs, uv)
		} else if strings.HasPrefix(line, "vn ") {
			n := Vec3f{}
			fmt.Sscanf(line, "vn %f %f %f", &n.X, &n.Y, &n.Z)
			n.Normalize()
			m.Normals = append(m.Normals, n)
		}
	}

	return m, nil
}
