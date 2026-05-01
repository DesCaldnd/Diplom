package compute

import "math"

type ScalarType = float64

type StaticPoint struct {
	dim  int64
	data [maxSupportedDim]ScalarType
}

func NewStaticPoint(dim int64) StaticPoint {
	if dim < 0 {
		dim = 0
	}
	if dim > maxSupportedDim {
		dim = maxSupportedDim
	}
	return StaticPoint{dim: dim}
}

func NewStaticPointFromValues(values ...ScalarType) StaticPoint {
	res := NewStaticPoint(int64(len(values)))
	for i := int64(0); i < res.dim; i++ {
		res.data[i] = values[i]
	}
	return res
}

func NewStaticPointFromSlice(src []ScalarType) StaticPoint {
	res := NewStaticPoint(int64(len(src)))
	for i := int64(0); i < res.dim; i++ {
		res.data[i] = src[i]
	}
	return res
}

func (p StaticPoint) Dim() int64 {
	return p.dim
}

func (p *StaticPoint) SetDim(dim int64) {
	if dim < 0 {
		dim = 0
	}
	if dim > maxSupportedDim {
		dim = maxSupportedDim
	}
	p.dim = dim
}

func (p StaticPoint) At(i int64) ScalarType {
	return p.data[i]
}

func (p *StaticPoint) Set(i int64, value ScalarType) {
	p.data[i] = value
}

func (p StaticPoint) Slice() []ScalarType {
	res := make([]ScalarType, p.dim)
	for i := int64(0); i < p.dim; i++ {
		res[i] = p.data[i]
	}
	return res
}

func (p *StaticPoint) CopyFrom(other StaticPoint) {
	p.dim = other.dim
	p.data = other.data
}

func (p *StaticPoint) Add(other StaticPoint) {
	for i := int64(0); i < p.dim; i++ {
		p.data[i] += other.data[i]
	}
}

func (p *StaticPoint) Sub(other StaticPoint) {
	for i := int64(0); i < p.dim; i++ {
		p.data[i] -= other.data[i]
	}
}

func (p *StaticPoint) Mul(other StaticPoint) {
	for i := int64(0); i < p.dim; i++ {
		p.data[i] *= other.data[i]
	}
}

func (p *StaticPoint) MulScalar(s ScalarType) {
	for i := int64(0); i < p.dim; i++ {
		p.data[i] *= s
	}
}

func (p StaticPoint) LengthSquared() ScalarType {
	var sum ScalarType
	for i := int64(0); i < p.dim; i++ {
		sum += p.data[i] * p.data[i]
	}
	return sum
}

func (p StaticPoint) Length() ScalarType {
	return math.Sqrt(p.LengthSquared())
}
