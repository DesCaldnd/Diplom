package compute

import "math"

type ScalarType = float64

type Point []ScalarType

func (p Point) Add(other Point) Point {
	res := make(Point, len(p))
	for i := range p {
		res[i] = p[i] + other[i]
	}
	return res
}

func (p Point) Sub(other Point) Point {
	res := make(Point, len(p))
	for i := range p {
		res[i] = p[i] - other[i]
	}
	return res
}

func (p Point) Mul(other Point) Point {
	res := make(Point, len(p))
	for i := range p {
		res[i] = p[i] * other[i]
	}
	return res
}

func (p Point) MulScalar(s ScalarType) Point {
	res := make(Point, len(p))
	for i := range p {
		res[i] = p[i] * s
	}
	return res
}

func (p Point) Length() ScalarType {
	var sum ScalarType
	for _, v := range p {
		sum += v * v
	}
	return math.Sqrt(sum)
}
