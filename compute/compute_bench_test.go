package compute_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"Diplom/compute"
)

func buildGridNoError(
	b *testing.B,
	funcEval func(compute.Point) (compute.Point, error),
	min, max compute.Point,
	epsilon float64,
	anchors []compute.Point,
	basis compute.BasisType,
	buildType compute.BuildType,
	maxLevel int64,
	maxNodes int64,
) *compute.AdaptiveSparseGrid {
	b.Helper()
	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, anchors, basis, buildType, maxLevel, maxNodes)
	if err != nil {
		b.Fatalf("failed to build grid: %v", err)
	}
	return grid
}

func makeSmoothFunction(dim int) func(compute.Point) (compute.Point, error) {
	return func(arg compute.Point) (compute.Point, error) {
		res := 0.0
		for i := 0; i < dim; i++ {
			freq := float64(i + 1)
			res += math.Sin(freq*arg[i]) * math.Cos(arg[i]/(freq+1.0))
		}
		return compute.Point{res}, nil
	}
}

func makeDomain(dim int, scale float64) (compute.Point, compute.Point) {
	min := make(compute.Point, dim)
	max := make(compute.Point, dim)
	for i := 0; i < dim; i++ {
		min[i] = 0.0
		max[i] = scale * math.Pi
	}
	return min, max
}

// ============================================================================
// БЕНЧМАРКИ
// ============================================================================

// Бенчмарк 1: зависимость времени построения от размерности и размера области.
func BenchmarkDomainScalingDifferentDimensions(b *testing.B) {
	scales := []float64{1.0, 2.0, 4.0}
	dimensions := []int{1, 2, 3}

	for _, dim := range dimensions {
		funcEval := makeSmoothFunction(dim)
		for _, scale := range scales {
			min, max := makeDomain(dim, scale)
			name := fmt.Sprintf("dim_%d_scale_%.0fpi", dim, scale)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
					if err != nil {
						b.Fatalf("failed to create grid: %v", err)
					}
				}
			})
		}
	}
}

// Бенчмарк 2: сравнение последовательного и параллельного построения для нетривиальной 2D функции.
func BenchmarkParallelVsSequentialBuild(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		value := math.Sin(3*x) + 0.35*math.Cos(25*y) + 0.2*math.Sin(7*x+2.5) - 0.15*math.Cos(11*y-1.0)
		return compute.Point{value}, nil
	}

	min := compute.Point{0.0, 0.0}
	max := compute.Point{2.0 * math.Pi, 2.0 * math.Pi}

	for _, basis := range []compute.BasisType{compute.BasisTypeLinear, compute.BasisTypeQuadratic} {
		basisName := "linear"
		if basis == compute.BasisTypeQuadratic {
			basisName = "quadratic"
		}

		b.Run("sequential_"+basisName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, nil, basis, compute.BuildTypeSequential, 0, 0)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})

		b.Run("parallel_"+basisName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, nil, basis, compute.BuildTypeParallel, 0, 0)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})
	}
}

// Бенчмарк 3: прямое сравнение линейного и квадратичного базисов при разных значениях epsilon.
func BenchmarkLinearVsQuadraticBasis(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		value := math.Sin(15*x) + 0.6*math.Cos(40*y) + 0.2*math.Sin(31*x+0.3) - 0.12*math.Cos(19*y-0.7) + math.Sin(x*y)
		return compute.Point{value}, nil
	}

	min := compute.Point{0.0, 0.0}
	max := compute.Point{math.Pi, math.Pi}
	epsValues := []float64{1e-2, 5e-3, 1e-3, 1e-4, 1e-5, 5e-6, 1e-6}

	for _, eps := range epsValues {
		b.Run("linear_eps_"+strconv.FormatFloat(eps, 'g', -1, 64), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, eps, nil, compute.BasisTypeLinear, compute.BuildTypeParallel, 0, 0)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})

		b.Run("quadratic_eps_"+strconv.FormatFloat(eps, 'g', -1, 64), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, eps, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})
	}
}

// Бенчмарк 4: влияние anchor points на функцию, которая обманывает базовый критерий адаптации.
func BenchmarkAnchorPointsImpact(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		return compute.Point{math.Sin(arg[0])}, nil
	}

	min := compute.Point{0.0}
	max := compute.Point{4.0 * math.Pi}
	anchors := []compute.Point{{math.Pi}, {3.0 * math.Pi}}

	b.Run("without_anchor_points", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
			if err != nil {
				b.Fatalf("failed to create grid: %v", err)
			}
		}
	})

	b.Run("with_anchor_points", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, anchors, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
			if err != nil {
				b.Fatalf("failed to create grid: %v", err)
			}
		}
	})
}

// Бенчмарк 5: сравнение двух способов работы с дифференциальным уравнением.
// 1. t как интервальная неопределённость (дополнительная размерность).
// 2. Последовательное применение make_next_iteration.
func BenchmarkDifferentialEquationApproaches(b *testing.B) {
	diffEq := func(state compute.Point, t float64) compute.Point {
		x := state[0]
		y := state[1]
		return compute.Point{-y / (1 + math.Sqrt(x*x+y*y)), -x / (1 + math.Sqrt(x*x+y*y))}
	}

	tMaxValues := []float64{2.0, 5.0, 10.0, 20.0}

	for _, tMax := range tMaxValues {
		b.Run("t_as_interval_uncertainty_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			funcWithT := func(arg compute.Point) (compute.Point, error) {
				x0 := arg[0]
				y0 := arg[1]
				t := arg[2]
				return integrateRk4(diffEq, compute.Point{x0, y0}, 0.0, t, 50), nil
			}
			min := compute.Point{-1.0, 0.0, 0.0}
			max := compute.Point{1.0, 1.0, tMax}
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcWithT, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})

		b.Run("iterative_make_next_iteration_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			integrate1s := func(state compute.Point) (compute.Point, error) {
				return integrateRk4(diffEq, state, 0.0, 1.0, 25), nil
			}
			min := compute.Point{-1.0, 0.0}
			max := compute.Point{1.0, 1.0}
			for i := 0; i < b.N; i++ {
				grid, err := compute.NewAdaptiveSparseGrid(integrate1s, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				if err != nil {
					b.Fatalf("failed to create initial grid: %v", err)
				}
				for step := 1; step < int(tMax); step++ {
					grid, err = grid.MakeNextIteration(integrate1s, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
					if err != nil {
						b.Fatalf("failed to create next iteration: %v", err)
					}
				}
			}
		})
	}
}

// Бенчмарк 6: влияние жёсткого ограничения числа узлов на время построения.
func BenchmarkNodeLimitOptimization(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		value := math.Sin(5*x) + math.Cos(33*y) + 0.22*math.Sin(17*x+0.4) - 0.18*math.Cos(21*y-0.2)
		return compute.Point{value}, nil
	}

	min := compute.Point{0, 0}
	max := compute.Point{2.0 * math.Pi, 2.0 * math.Pi}
	limits := []int64{0, 200, 500, 1000}

	for _, limit := range limits {
		name := "limit_" + strconv.FormatInt(limit, 10)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, limit)
				if err != nil {
					b.Fatalf("failed to create grid: %v", err)
				}
			}
		})
	}
}

// Бенчмарк 7: сравнение стоимости вычисления значения по уже построенной сетке.
func BenchmarkEvaluationCostAfterBuild(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		value := math.Sin(15*x) + 0.6*math.Cos(40*y) + 0.2*math.Sin(31*x+0.3) - 0.12*math.Cos(19*y-0.7) + math.Sin(x*y)
		return compute.Point{value}, nil
	}

	min := compute.Point{0.0, 0.0}
	max := compute.Point{math.Pi, math.Pi}
	testPoint := compute.Point{0.73, 1.11}

	linearGrid := buildGridNoError(b, funcEval, min, max, 0.001, nil, compute.BasisTypeLinear, compute.BuildTypeParallel, 0, 0)
	quadraticGrid := buildGridNoError(b, funcEval, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	b.Run("evaluate_linear_grid", func(b *testing.B) {
		for i := 0; i < b.N*10000; i++ {
			_, err := linearGrid.Evaluate(testPoint)
			if err != nil {
				b.Fatalf("failed to evaluate grid: %v", err)
			}
		}
	})

	b.Run("evaluate_quadratic_grid", func(b *testing.B) {
		for i := 0; i < b.N*10000; i++ {
			_, err := quadraticGrid.Evaluate(testPoint)
			if err != nil {
				b.Fatalf("failed to evaluate grid: %v", err)
			}
		}
	})
}

// Бенчмарк 8: сравнение всех сочетаний BuildType и BasisType на 4-мерной функции.
func BenchmarkFourDimensionalBuildTypeBasisProduct(b *testing.B) {
	funcEval := func(arg compute.Point) (compute.Point, error) {
		x1 := arg[0]
		x2 := arg[1]
		x3 := arg[2]
		x4 := arg[3]
		value := math.Sin(30*x1)*math.Cos(22.2323*x2) + math.Sin(70*x3) - 0.15*math.Sin(x4) + math.Cos(x1*x3)
		return compute.Point{value}, nil
	}

	min := compute.Point{0.0, 0.0, 0.0, 0.0}
	max := compute.Point{math.Pi, math.Pi, math.Pi, math.Pi}
	epsilon := 0.0000005

	cases := []struct {
		buildType compute.BuildType
		basisType compute.BasisType
		name      string
	}{
		{compute.BuildTypeSequential, compute.BasisTypeLinear, "sequential_linear"},
		{compute.BuildTypeSequential, compute.BasisTypeQuadratic, "sequential_quadratic"},
		{compute.BuildTypeParallel, compute.BasisTypeLinear, "parallel_linear"},
		{compute.BuildTypeParallel, compute.BasisTypeQuadratic, "parallel_quadratic"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, tc.basisType, tc.buildType, 0, 0)
				if err != nil {
					b.Fatalf("failed to create 4D grid: %v", err)
				}
			}
		})
	}
}
