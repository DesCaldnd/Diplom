package compute_test

import (
	"fmt"
	"math"

	"Diplom/compute"
)

func ExampleAdaptiveSparseGrid_simple() {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) * arg.At(0)), nil
	}

	min := scalarPoint(-2.0)
	max := scalarPoint(2.0)

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := scalarPoint(1.5)
	result, _ := grid.Evaluate(testPoint)
	fmt.Printf("f(1.5) = %.2f (expected: 2.25)\n", result.At(0))
	// Output: f(1.5) = 2.25 (expected: 2.25)
}

func ExampleAdaptiveSparseGrid_multidim() {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0)*arg.At(1) + math.Sin(arg.At(2))), nil
	}

	min := scalarPoint(0.0, 0.0, 0.0)
	max := scalarPoint(2.0, 2.0, math.Pi)

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.05, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := scalarPoint(1.0, 1.5, math.Pi/2.0)
	result, _ := grid.Evaluate(testPoint)
	fmt.Printf("f(1.0, 1.5, pi/2) = %.1f (expected: 2.5)\n", result.At(0))
	// Output: f(1.0, 1.5, pi/2) = 2.5 (expected: 2.5)
}

func ExampleAdaptiveSparseGrid_ode() {
	diffEq := func(state compute.StaticPoint, t float64) compute.StaticPoint {
		alpha := 2.0 / 3.0
		beta := 4.0 / 3.0
		delta := 1.0
		gamma := 1.0
		x := state.At(0)
		y := state.At(1)
		return scalarPoint(alpha*x-beta*x*y, delta*x*y-gamma*y)
	}

	funcEval := func(initialState compute.StaticPoint) (compute.StaticPoint, error) {
		return integrateRk4(diffEq, initialState, 0.0, 2.0, 100), nil
	}

	min := scalarPoint(0.5, 0.5)
	max := scalarPoint(2.0, 2.0)

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.05, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := scalarPoint(1.0, 1.0)
	result, _ := grid.Evaluate(testPoint)
	expected, _ := funcEval(testPoint)

	fmt.Printf("Interpolated close to expected: %v\n", math.Abs(result.At(0)-expected.At(0)) < 0.1 && math.Abs(result.At(1)-expected.At(1)) < 0.1)
	// Output: Interpolated close to expected: true
}

func ExampleAdaptiveSparseGrid_composition() {
	g := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) + 1.0), nil
	}

	f := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) * arg.At(0)), nil
	}

	min := scalarPoint(0.0)
	max := scalarPoint(2.0)

	gridG, _ := compute.NewAdaptiveSparseGrid(g, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridFG, _ := gridG.MakeNextIteration(f, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := scalarPoint(1.0)
	result, _ := gridFG.Evaluate(testPoint)
	fmt.Printf("f(g(1.0)) = %.1f (expected: 4.0)\n", result.At(0))
	// Output: f(g(1.0)) = 4.0 (expected: 4.0)
}

func ExampleAdaptiveSparseGrid_odeComposition() {
	diffEq := func(state compute.StaticPoint, t float64) compute.StaticPoint {
		return scalarPoint(-0.5 * state.At(0))
	}

	integrate1s := func(state compute.StaticPoint) (compute.StaticPoint, error) {
		return integrateRk4(diffEq, state, 0.0, 1.0, 50), nil
	}

	min := scalarPoint(0.0)
	max := scalarPoint(10.0)

	gridT1, _ := compute.NewAdaptiveSparseGrid(integrate1s, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridT2, _ := gridT1.MakeNextIteration(integrate1s, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridT3, _ := gridT2.MakeNextIteration(integrate1s, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := scalarPoint(8.0)
	resultT1, _ := gridT1.Evaluate(testPoint)
	resultT2, _ := gridT2.Evaluate(testPoint)
	resultT3, _ := gridT3.Evaluate(testPoint)

	fmt.Printf("State at t=1 close to expected: %v\n", math.Abs(resultT1.At(0)-8.0*math.Exp(-0.5*1.0)) < 0.05)
	fmt.Printf("State at t=2 close to expected: %v\n", math.Abs(resultT2.At(0)-8.0*math.Exp(-0.5*2.0)) < 0.05)
	fmt.Printf("State at t=3 close to expected: %v\n", math.Abs(resultT3.At(0)-8.0*math.Exp(-0.5*3.0)) < 0.05)
	// Output:
	// State at t=1 close to expected: true
	// State at t=2 close to expected: true
	// State at t=3 close to expected: true
}
