package gogl

type (
	Matrix4[N Number] struct {
		M [4][4]N
	}

	Matrix4f = Matrix4[float64]
)

func NewIdentityMatrix4[N Number]() Matrix4[N] {
	m := Matrix4[N]{}

	m.M[0][0] = N(1)
	m.M[1][1] = N(1)
	m.M[2][2] = N(1)
	m.M[3][3] = N(1)

	return m
}


func (m *Matrix4[N]) Multiply(n *Matrix4[N]) *Matrix4[N] {
	r := Matrix4[N]{}

	// TODO unroll?
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				r.M[i][j] += m.M[i][k] * n.M[k][j]
			}
		}
	}

	return &r
}

func (m *Matrix4[N]) MultiplyVec3(v *Vec3[N]) *Vec3[N] {
	r := &Vec3[N]{}

	t :=  m.M[3][0]*v.X + m.M[3][1]*v.Y + m.M[3][2]*v.Z + m.M[3][3]

	r.X = (m.M[0][0]*v.X + m.M[0][1]*v.Y + m.M[0][2]*v.Z + m.M[0][3])/t
	r.Y = (m.M[1][0]*v.X + m.M[1][1]*v.Y + m.M[1][2]*v.Z + m.M[1][3])/t
	r.Z = (m.M[2][0]*v.X + m.M[2][1]*v.Y + m.M[2][2]*v.Z + m.M[2][3])/t

	return r
}