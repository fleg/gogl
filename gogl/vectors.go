package gogl

import "math"

type (
	Number interface {
		int | float64
	}

	Vec2[N Number] struct {
		X N
		Y N
	}

	Vec3[N Number] struct {
		X N
		Y N
		Z N
	}

	Vec3f = Vec3[float64]
	Vec2f = Vec2[float64]

	Vec3i = Vec3[int]
	Vec2i = Vec2[int]
)

// a × b = {ay*bz - az*by; az*bx - ax*bz; ax*by - ay*bx}
func (v *Vec3[N]) CrossProduct(u *Vec3[N]) *Vec3[N] {
	return &Vec3[N]{
		X: v.Y*u.Z - v.Z*u.Y,
		Y: v.Z*u.X - v.X*u.Z,
		Z: v.X*u.Y - v.Y*u.X,
	}
}

func (v *Vec3[N]) Subtract(u *Vec3[N]) *Vec3[N] {
	return &Vec3[N]{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v *Vec3[N]) DotProduct(u *Vec3[N]) N {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v *Vec3[N]) Length() N {
	return N(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v *Vec3[N]) Normalize() {
	l := v.Length()
	v.X /= l
	v.Y /= l
	v.Z /= l
}
