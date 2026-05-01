package compute

import "math"

const maxSupportedDim = 8

type BasisType int64

const (
	BasisTypeLinear BasisType = iota
	BasisTypeQuadratic
)

type direction int64

const (
	directionNone  direction = 0
	directionLeft  direction = 1 << 0
	directionRight direction = 1 << 1
	directionBoth  direction = directionLeft | directionRight
)

type gridKey struct {
	level [maxSupportedDim]int64
	index [maxSupportedDim]int64
}

type node struct {
	key         gridKey
	centerUnit  StaticPoint
	leftBound   [maxSupportedDim]float64
	rightBound  [maxSupportedDim]float64
	alpha       StaticPoint
	hasChildren bool
	maxLevel    int64
}

func (n *node) isPointInAffectZone(point StaticPoint) bool {
	for i := int64(0); i < point.dim; i++ {
		if n.key.level[i] != 0 {
			if point.data[i] <= n.leftBound[i] || point.data[i] >= n.rightBound[i] {
				return false
			}
		} else {
			expected := float64(0)
			if n.key.index[i] != 0 {
				expected = 1
			}
			if math.Abs(point.data[i]-expected) > 1e-9 {
				return false
			}
		}
	}
	return true
}

type entryPoint struct {
	node       node
	dimensions int64
}

func eval1d(x float64, level int64, index int64, basisType BasisType) float64 {
	if level == 0 {
		if index == 0 {
			return 1.0 - x
		}
		if index == 1 {
			return x
		}
		return 0.0
	}

	h := 1.0 / float64(int64(1)<<level)
	center := float64(index) * h
	dist := math.Abs(x - center)

	if dist >= h {
		return 0.0
	}

	t := dist / h
	switch basisType {
	case BasisTypeLinear:
		return 1.0 - t
	case BasisTypeQuadratic:
		return 1.0 - t*t
	default:
		return 0.0
	}
}

func evalBasis(x StaticPoint, n node, basisType BasisType) float64 {
	result := 1.0
	for i := int64(0); i < x.dim; i++ {
		result *= eval1d(x.data[i], n.key.level[i], n.key.index[i], basisType)
		if result == 0 {
			return result
		}
	}
	return result
}

func getCoord(level int64, index int64) float64 {
	if level == 0 {
		if index == 0 {
			return 0.0
		}
		return 1.0
	}
	return float64(index) / float64(int64(1)<<level)
}

func fillAffectDirections(x StaticPoint, n node, dst *[maxSupportedDim]direction) {
	for i := int64(0); i < x.dim; i++ {
		if x.data[i] < n.centerUnit.data[i] {
			dst[i] = directionLeft
		} else {
			dst[i] = directionRight
		}
	}
}
