package compute_test

import (
	"math"
	"math/rand"
	"testing"

	"Diplom/compute"
)

func scalarPoint(values ...float64) compute.StaticPoint {
	return compute.NewStaticPointFromValues(values...)
}

func rk4Step(f func(compute.StaticPoint, float64) compute.StaticPoint, x compute.StaticPoint, t float64, dt float64) compute.StaticPoint {
	k1 := f(x, t)

	x1 := x
	k1Half := k1
	k1Half.MulScalar(dt / 2.0)
	x1.Add(k1Half)
	k2 := f(x1, t+dt/2.0)

	x2 := x
	k2Half := k2
	k2Half.MulScalar(dt / 2.0)
	x2.Add(k2Half)
	k3 := f(x2, t+dt/2.0)

	x3 := x
	k3Dt := k3
	k3Dt.MulScalar(dt)
	x3.Add(k3Dt)
	k4 := f(x3, t+dt)

	res := x
	k1Copy := k1
	k2Copy := k2
	k2Copy.MulScalar(2.0)
	k3Copy := k3
	k3Copy.MulScalar(2.0)

	k1Copy.Add(k2Copy)
	k1Copy.Add(k3Copy)
	k1Copy.Add(k4)
	k1Copy.MulScalar(dt / 6.0)

	res.Add(k1Copy)
	return res
}

func integrateRk4(f func(compute.StaticPoint, float64) compute.StaticPoint, x0 compute.StaticPoint, t0 float64, t1 float64, steps int) compute.StaticPoint {
	dt := (t1 - t0) / float64(steps)
	x := x0
	t := t0
	for i := 0; i < steps; i++ {
		x = rk4Step(f, x, t, dt)
		t += dt
	}
	return x
}

func TestSimple1DFunction(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) * arg.At(0)), nil
	}

	min := scalarPoint(-2.0)
	max := scalarPoint(2.0)
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(-2.0), scalarPoint(2.0), scalarPoint(0.0), scalarPoint(-1.5), scalarPoint(1.5), scalarPoint(0.123), scalarPoint(-0.987),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(min.At(0)+rand.Float64()*(max.At(0)-min.At(0))))
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected, _ := funcEval(testPoint)

		if math.Abs(result.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected.At(0), result.At(0))
		}
	}
}

func TestMultidimFunction(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0)*arg.At(1) + math.Sin(arg.At(0))), nil
	}

	min := scalarPoint(0.0, 0.0)
	max := scalarPoint(2.0, 2.0)
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeSequential, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(0.0, 0.0), scalarPoint(2.0, 2.0), scalarPoint(0.0, 2.0), scalarPoint(2.0, 0.0),
		scalarPoint(1.0, 1.0), scalarPoint(1.0, 1.5), scalarPoint(0.5, 0.5), scalarPoint(1.23, 0.87),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(
			min.At(0)+rand.Float64()*(max.At(0)-min.At(0)),
			min.At(1)+rand.Float64()*(max.At(1)-min.At(1)),
		))
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected, _ := funcEval(testPoint)

		if math.Abs(result.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected.At(0), result.At(0))
		}
	}
}

func TestDifferentialEquation(t *testing.T) {
	diffEq := func(state compute.StaticPoint, t float64) compute.StaticPoint {
		return scalarPoint(-state.At(0))
	}

	funcEval := func(initialState compute.StaticPoint) (compute.StaticPoint, error) {
		return integrateRk4(diffEq, initialState, 0.0, 1.0, 100), nil
	}

	min := scalarPoint(0.0)
	max := scalarPoint(5.0)
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(0.0), scalarPoint(5.0), scalarPoint(2.5), scalarPoint(1.0), scalarPoint(2.0), scalarPoint(3.14), scalarPoint(4.99),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(min.At(0)+rand.Float64()*(max.At(0)-min.At(0))))
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected, _ := funcEval(testPoint)

		if math.Abs(result.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected.At(0), result.At(0))
		}
	}
}

func TestAnchorPoints(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(math.Sin(arg.At(0))), nil
	}

	min := scalarPoint(0.0)
	max := scalarPoint(2.0 * math.Pi)
	epsilon := 0.001

	gridWithoutAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	anchors := []compute.StaticPoint{scalarPoint(math.Pi / 2.0), scalarPoint(3.0 * math.Pi / 2.0)}
	gridWithAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, anchors, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(math.Pi / 2.0), scalarPoint(math.Pi), scalarPoint(3.0 * math.Pi / 2.0), scalarPoint(math.Pi / 4.0), scalarPoint(1.0), scalarPoint(5.0),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(min.At(0)+rand.Float64()*(max.At(0)-min.At(0))))
	}

	for _, testPoint := range testPoints {
		resultWithout, _ := gridWithoutAnchors.Evaluate(testPoint)
		resultWith, _ := gridWithAnchors.Evaluate(testPoint)
		expected, _ := funcEval(testPoint)

		if math.Abs(resultWithout.At(0)-0.0) > epsilon*2 && math.Abs(resultWithout.At(0)-expected.At(0)) > epsilon*2 {
		}

		if math.Abs(resultWith.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("with anchors at %v: expected %v, got %v", testPoint, expected.At(0), resultWith.At(0))
		}
	}
}

func TestEvaluateInvalidDimension(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) + arg.At(1)), nil
	}

	grid, err := compute.NewAdaptiveSparseGrid(
		funcEval,
		scalarPoint(0.0, 0.0),
		scalarPoint(1.0, 1.0),
		0.001,
		nil,
		compute.BasisTypeQuadratic,
		compute.BuildTypeParallel,
		0,
		0,
	)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	_, err = grid.Evaluate(scalarPoint(0.5))
	if err == nil {
		t.Fatal("expected error for invalid input dimension, got nil")
	}
}

func TestEvaluateOutOfBounds(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) + arg.At(1)), nil
	}

	grid, err := compute.NewAdaptiveSparseGrid(
		funcEval,
		scalarPoint(0.0, 0.0),
		scalarPoint(1.0, 1.0),
		0.001,
		nil,
		compute.BasisTypeQuadratic,
		compute.BuildTypeParallel,
		0,
		0,
	)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	invalidPoints := []compute.StaticPoint{
		scalarPoint(-0.1, 0.5),
		scalarPoint(0.5, -0.1),
		scalarPoint(1.1, 0.5),
		scalarPoint(0.5, 1.1),
	}

	for _, point := range invalidPoints {
		_, err = grid.Evaluate(point)
		if err == nil {
			t.Fatalf("expected out of bounds error for point %v, got nil", point)
		}
	}
}

func TestMakeNextIteration(t *testing.T) {
	g := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) + 1.0), nil
	}

	f := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0) * arg.At(0)), nil
	}

	min := scalarPoint(0.0)
	max := scalarPoint(2.0)
	epsilon := 0.001

	gridG, err := compute.NewAdaptiveSparseGrid(g, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create gridG: %v", err)
	}

	gridFG, err := gridG.MakeNextIteration(f, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create gridFG: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(0.0), scalarPoint(2.0), scalarPoint(1.0), scalarPoint(0.5), scalarPoint(1.5), scalarPoint(0.123), scalarPoint(1.987),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(min.At(0)+rand.Float64()*(max.At(0)-min.At(0))))
	}

	for _, testPoint := range testPoints {
		result, err := gridFG.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		gResult, err := g(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate g at %v: %v", testPoint, err)
		}
		expected, err := f(gResult)
		if err != nil {
			t.Fatalf("failed to evaluate f at %v: %v", testPoint, err)
		}

		if math.Abs(result.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected.At(0), result.At(0))
		}
	}
}

func TestLinearBasis(t *testing.T) {
	funcEval := func(arg compute.StaticPoint) (compute.StaticPoint, error) {
		return scalarPoint(arg.At(0)*arg.At(1) + math.Sin(arg.At(0))), nil
	}

	min := scalarPoint(0.0, 0.0)
	max := scalarPoint(2.0, 2.0)
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeLinear, compute.BuildTypeSequential, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.StaticPoint{
		scalarPoint(0.0, 0.0), scalarPoint(2.0, 2.0), scalarPoint(0.0, 2.0), scalarPoint(2.0, 0.0),
		scalarPoint(1.0, 1.0), scalarPoint(1.0, 1.5), scalarPoint(0.5, 0.5), scalarPoint(1.23, 0.87),
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, scalarPoint(
			min.At(0)+rand.Float64()*(max.At(0)-min.At(0)),
			min.At(1)+rand.Float64()*(max.At(1)-min.At(1)),
		))
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected, _ := funcEval(testPoint)

		if math.Abs(result.At(0)-expected.At(0)) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected.At(0), result.At(0))
		}
	}
}
