package compute_test

import (
	"math"
	"strconv"
	"testing"

	"Diplom/compute"
)

func BenchmarkSimpleFunctionsDifferentDimensions(b *testing.B) {
	func1d := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(math.Sin(arg.At(0))), nil
	}
	func2d := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(math.Sin(arg.At(0)) * math.Cos(arg.At(1))), nil
	}
	func3d := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(math.Sin(arg.At(0)) * math.Cos(arg.At(1)) * math.Sin(arg.At(2))), nil
	}

	b.Run("1D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func1d, scalarPoint(0.0), scalarPoint(math.Pi), 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})

	b.Run("2D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func2d, scalarPoint(0.0, 0.0), scalarPoint(math.Pi, math.Pi), 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})

	b.Run("3D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compute.NewAdaptiveSparseGrid(func3d, scalarPoint(0.0, 0.0, 0.0), scalarPoint(math.Pi, math.Pi, math.Pi), 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
		}
	})
}

func BenchmarkBuildTypeAndBasisTypeComparison(b *testing.B) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(math.Sin(arg.At(0)*2.0) * math.Cos(arg.At(1)*2.0)), nil
	}

	min := scalarPoint(0.0, 0.0)
	max := scalarPoint(1.0, 1.0)
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

func BenchmarkDifferentialEquationApproaches(b *testing.B) {
	diffEq := func(state compute.StaticPoint, t float64) compute.StaticPoint {
		return scalarPoint(-0.5 * state.At(0))
	}

	tMaxValues := []float64{2.0, 10.0, 30.0, 50.0}

	for _, tMax := range tMaxValues {
		b.Run("Approach1_t_as_dimension_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			funcWithT := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
				x0 := arg.At(0)
				t := arg.At(1)
				return integrateRk4(diffEq, scalarPoint(x0), 0.0, t, 50), nil
			}
			min := scalarPoint(0.0, 0.0)
			max := scalarPoint(10.0, tMax)
			for i := 0; i < b.N; i++ {
				compute.NewAdaptiveSparseGrid(funcWithT, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
			}
		})

		b.Run("Approach2_make_next_iteration_tMax_"+strconv.Itoa(int(tMax)), func(b *testing.B) {
			integrate1s := func(state compute.StaticPoint) (compute.StaticPoint, error) {
				return integrateRk4(diffEq, state, 0.0, 1.0, 25), nil
			}
			min := scalarPoint(0.0)
			max := scalarPoint(10.0)
			for i := 0; i < b.N; i++ {
				grid, _ := compute.NewAdaptiveSparseGrid(integrate1s, min, max, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				for step := 1; step < int(tMax); step++ {
					grid, _ = grid.MakeNextIteration(integrate1s, 0.001, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
				}
			}
		})
	}
}
