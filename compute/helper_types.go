package compute

import (
	"fmt"
	"math"
)

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
	level []int64
	index []int64
}

func (k gridKey) String() string {
	return fmt.Sprintf("%v|%v", k.level, k.index)
}

type node struct {
	key         gridKey
	centerUnit  Point
	alpha       Point
	hasChildren bool
}

func (n *node) isPointInAffectZone(point Point, inDim int64) bool {
	for i := int64(0); i < inDim; i++ {
		if n.key.level[i] != 0 {
			boundsFirst, boundsSecond := getAffectBounds(n.key.level[i], n.key.index[i])
			if point[i] <= boundsFirst || point[i] >= boundsSecond {
				return false
			}
		} else {
			expected := float64(0)
			if n.key.index[i] != 0 {
				expected = 1
			}
			if math.Abs(point[i]-expected) > 1e-9 {
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
	var result float64

	switch basisType {
	case BasisTypeLinear:
		result = 1.0 - t
	case BasisTypeQuadratic:
		result = 1.0 - t*t
	}

	return result
}

func evalBasis(x Point, level []int64, index []int64, basisType BasisType, inDim int64) float64 {
	result := 1.0
	for i := int64(0); i < inDim; i++ {
		result *= eval1d(x[i], level[i], index[i], basisType)
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

func getAffectBounds(level int64, index int64) (float64, float64) {
	unit := 1.0 / float64(int64(1)<<level)
	coord := float64(index) / float64(int64(1)<<level)
	return coord - unit, coord + unit
}

func getAffectDirection(x Point, level []int64, index []int64, inDim int64) []direction {
	result := make([]direction, inDim)
	for i := int64(0); i < inDim; i++ {
		var coord float64
		if level[i] > 0 {
			coord = float64(index[i]) / float64(int64(1)<<level[i])
		} else {
			coord = float64(index[i])
		}
		if x[i] < coord {
			result[i] = directionLeft
		} else {
			result[i] = directionRight
		}
	}
	return result
}
