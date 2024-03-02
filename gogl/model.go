package gogl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	Model struct {
		Verticies []Vec3f
		Faces     [][]int
	}
)

func NewModelFromFile(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := &Model{
		Verticies: make([]Vec3f, 0),
		Faces:     make([][]int, 0),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "v ") {
			var v Vec3f
			fmt.Sscanf(line, "v %f %f %f", &v.X, &v.Y, &v.Z)
			m.Verticies = append(m.Verticies, v)
		} else if strings.HasPrefix(line, "f ") {
			indicies := make([]int, 3)
			textures := make([]int, 3)
			nindicies := make([]int, 3)

			if strings.Contains(line, "/") {
				fmt.Sscanf(
					line,
					"f %d/%d/%d %d/%d/%d %d/%d/%d",
					&indicies[0], &textures[0], &nindicies[0],
					&indicies[1], &textures[1], &nindicies[1],
					&indicies[2], &textures[2], &nindicies[2],
				)
			} else {
				fmt.Sscanf(line, "f %d %d %d", &indicies[0], &indicies[1], &indicies[2])
			}

			for i := 0; i < 3; i++ {
				indicies[i] -= 1
				textures[i] -= 1
				nindicies[i] -= 1
			}

			m.Faces = append(m.Faces, indicies)
		}
	}

	return m, nil
}
