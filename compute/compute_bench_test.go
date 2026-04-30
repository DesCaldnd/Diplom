package compute_test

import (
	"math"
	"strconv"
	"testing"

	"Diplom/compute"
)

// ============================================================================
// БЕНЧМАРКИ
// ============================================================================

// Бенчмарк 1: Замеры простых функций разной размерности
// Проверяется время построения сетки для функций 1D, 2D и 3D.
func BenchmarkSimpleFunctionsDifferentDimensions(b *testing.B) {
	func1d := func(arg compute.Point) (compute.Point, error) { return compute.Point{math.Sin(arg[0])}, nil }
	func2d := func(arg compute.Point) (compute.Point, error) {
		return compute.Point{math.Sin(arg[0]) * math.Cos(arg[1])}, nil
	}
	func3d := func(arg compute.Point) (compute.Point, error) {
		return compute.Point{math.Sin(arg[0]) * math.Cos(arg[1]) * math.Sin(arg[2])}, nil
	}

	b.Run("1D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func1d, compute.Point{0.0}, compute.Point{math.Pi}, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})

	b.Run("2D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func2d, compute.Point{0.0, 0.0}, compute.Point{math.Pi, math.Pi}, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})

	b.Run("3D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func3d, compute.Point{0.0, 0.0, 0.0}, compute.Point{math.Pi, math.Pi, math.Pi}, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})
}

// Бенчмарк 2: Сравнения на одинаковых параметрах с разным build_type, basis_type
// Проверяется время построения сетки для 2D функции при SEQUENTIAL vs PARALLEL
// и LINEAR vs QUADRATIC базисах.
func BenchmarkBuildTypeAndBasisTypeComparison(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		return compute.Point{math.Sin(arg[0]*2.0) * math.Cos(arg[1]*2.0)}, nil
	}

	min := compute.Point{0.0, 0.0}
	max := compute.Point{1.0, 1.0}
	epsilon := 0.001

	b.Run("Sequential_Quadratic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeSequential, 0, 0)
		}
	})

	b.Run("Parallel_Quadratic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})

	b.Run("Sequential_Linear", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeLinear, compute.BuildTypeSequential, 0, 0)
		}
	})

	b.Run("Parallel_Linear", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeLinear, compute.BuildTypeParallel, 0, 0)
		}
	})
}

// Бенчмарк 3: Сравнения дифференциальных уравнений
// Проверяется два подхода к интерполяции решения диффура на отрезке времени [0, 2]:
// 1. t - это интервальная неопределенность (дополнительная размерность сетки).
// 2. Последовательное интегрирование при помощи make_next_iteration (t=1, затем t=2).
func BenchmarkDifferentialEquationApproaches(b *testing.B) {
	diffEq := func(state compute.Point, t float64) compute.Point {
		return compute.Point{-0.5 * state[0]}
	}

	tMaxValues := []float64{2.0, 10.0, 30.0, 50.0}

	for _, tMax := range tMaxValues {
		b.Run("Approach1_t_as_dimension_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			funcWithT := func(arg compute.Point) (compute.Point, error) {
				x0 := arg[0]
				t := arg[1]
				return integrateRk4(diffEq, compute.Point{x0}, 0.0, t, 50), nil
			}
			min := compute.Point{0.0, 0.0}
			max := compute.Point{10.0, tMax}
			for i := 0; i < b.N; i++ {
				compute.NewAdaptiveSparseGrid(funcWithT, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
			}
		})

		b.Run("Approach2_make_next_iteration_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			integrate1s := func(state compute.Point) (compute.Point, error) {
				return integrateRk4(diffEq, state, 0.0, 1.0, 25), nil
			}
			min := compute.Point{0.0}
			max := compute.Point{10.0}
			for i := 0; i < b.N; i++ {
				grid, _ := compute.NewAdaptiveSparseGrid(integrate1s, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				for step := 1; step < int(tMax); step++ {
					grid, _ = grid.MakeNextIteration(integrate1s, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				}
			}
		})
	}
}
