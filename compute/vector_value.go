package compute

import "math"

type ScalarType = float64

type Point []ScalarType

func (p Point) Clone() Point {
	res := make(Point, len(p))
	copy(res, p)
	return res
}

func (p Point) Add(other Point) {
	for i := range p {
		p[i] += other[i]
	}
}

func (p Point) Sub(other Point) {
	for i := range p {
		p[i] -= other[i]
	}
}

func (p Point) Mul(other Point) {
	for i := range p {
		p[i] *= other[i]
	}
}

func (p Point) MulScalar(s ScalarType) {
	for i := range p {
		p[i] *= s
	}
}

func (p Point) Length() ScalarType {
	var sum ScalarType
	for _, v := range p {
		sum += v * v
	}
	return math.Sqrt(sum)
}
