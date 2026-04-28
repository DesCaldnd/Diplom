package compute_test

import (
	"math"
	"math/rand"
	"testing"

	"Diplom/compute"
)

// Простой интегратор методом Рунге-Кутты 4-го порядка
func rk4Step(f func(compute.Point, float64) compute.Point, x compute.Point, t float64, dt float64) compute.Point {
	k1 := f(x, t)
	k2 := f(x.Add(k1.MulScalar(dt / 2.0)), t+dt/2.0)
	k3 := f(x.Add(k2.MulScalar(dt / 2.0)), t+dt/2.0)
	k4 := f(x.Add(k3.MulScalar(dt)), t+dt)
	return x.Add(k1.Add(k2.MulScalar(2.0)).Add(k3.MulScalar(2.0)).Add(k4).MulScalar(dt / 6.0))
}

func integrateRk4(f func(compute.Point, float64) compute.Point, x0 compute.Point, t0 float64, t1 float64, steps int) compute.Point {
	dt := (t1 - t0) / float64(steps)
	x := x0
	t := t0
	for i := 0; i < steps; i++ {
		x = rk4Step(f, x, t, dt)
		t += dt
	}
	return x
}

// ============================================================================
// ТЕСТЫ
// ============================================================================

// Тест 1: Корректность работы с обычной одномерной функцией
// Проверяется, что интерполяция функции f(x) = x^2 работает корректно
// и ошибка не превышает заданный epsilon.
func TestSimple1DFunction(t *testing.T) {
	funcEval := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] * arg[0]}
	}

	min := compute.Point{-2.0}
	max := compute.Point{2.0}
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.Point{
		{-2.0}, {2.0}, {0.0}, {-1.5}, {1.5}, {0.123}, {-0.987},
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, compute.Point{min[0] + rand.Float64()*(max[0]-min[0])})
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected := funcEval(testPoint)

		if math.Abs(result[0]-expected[0]) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected[0], result[0])
		}
	}
}

// Тест 2: Корректность работы с многомерной функцией
// Проверяется интерполяция функции f(x, y) = x*y + sin(x)
// с различными параметрами (базис, тип сборки).
func TestMultidimFunction(t *testing.T) {
	funcEval := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0]*arg[1] + math.Sin(arg[0])}
	}

	min := compute.Point{0.0, 0.0}
	max := compute.Point{2.0, 2.0}
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeSequential, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.Point{
		{0.0, 0.0}, {2.0, 2.0}, {0.0, 2.0}, {2.0, 0.0},
		{1.0, 1.0}, {1.0, 1.5}, {0.5, 0.5}, {1.23, 0.87},
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, compute.Point{
			min[0] + rand.Float64()*(max[0]-min[0]),
			min[1] + rand.Float64()*(max[1]-min[1]),
		})
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected := funcEval(testPoint)

		if math.Abs(result[0]-expected[0]) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected[0], result[0])
		}
	}
}

// Тест 3: Корректность работы с дифференциальным уравнением
// Проверяется интерполяция решения дифференциального уравнения dx/dt = -x.
// Сравнивается с эталонным аналитическим решением.
func TestDifferentialEquation(t *testing.T) {
	diffEq := func(state compute.Point, t float64) compute.Point {
		return compute.Point{-state[0]}
	}

	funcEval := func(initialState compute.Point) compute.Point {
		return integrateRk4(diffEq, initialState, 0.0, 1.0, 100)
	}

	min := compute.Point{0.0}
	max := compute.Point{5.0}
	epsilon := 0.001

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.Point{
		{0.0}, {5.0}, {2.5}, {1.0}, {2.0}, {3.14}, {4.99},
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, compute.Point{min[0] + rand.Float64()*(max[0]-min[0])})
	}

	for _, testPoint := range testPoints {
		result, err := grid.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected := funcEval(testPoint)

		if math.Abs(result[0]-expected[0]) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected[0], result[0])
		}
	}
}

// Тест 4: Тест на anchor points
// Проверяется функция sin(x) на отрезке [0, 2*pi].
// Значения на краях и в центре равны 0. Без anchor points алгоритм
// может посчитать функцию тождественным нулем.
// С anchor points (например, pi/2 и 3*pi/2) алгоритм должен отработать корректно.
func TestAnchorPoints(t *testing.T) {
	funcEval := func(arg compute.Point) compute.Point {
		return compute.Point{math.Sin(arg[0])}
	}

	min := compute.Point{0.0}
	max := compute.Point{2.0 * math.Pi}
	epsilon := 0.001

	gridWithoutAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	anchors := []compute.Point{{math.Pi / 2.0}, {3.0 * math.Pi / 2.0}}
	gridWithAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, anchors, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoints := []compute.Point{
		{math.Pi / 2.0}, {math.Pi}, {3.0 * math.Pi / 2.0}, {math.Pi / 4.0}, {1.0}, {5.0},
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, compute.Point{min[0] + rand.Float64()*(max[0]-min[0])})
	}

	for _, testPoint := range testPoints {
		resultWithout, _ := gridWithoutAnchors.Evaluate(testPoint)
		resultWith, _ := gridWithAnchors.Evaluate(testPoint)
		expected := funcEval(testPoint)

		// Without anchors, it might just be a flat line at 0
		if math.Abs(resultWithout[0]-0.0) > epsilon*2 && math.Abs(resultWithout[0]-expected[0]) > epsilon*2 {
			// It's either 0 or correct, but usually 0
		}

		if math.Abs(resultWith[0]-expected[0]) > epsilon*2 {
			t.Errorf("with anchors at %v: expected %v, got %v", testPoint, expected[0], resultWith[0])
		}
	}
}

func TestMakeNextIteration(t *testing.T) {
	g := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] + 1.0}
	}

	f := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] * arg[0]}
	}

	min := compute.Point{0.0}
	max := compute.Point{2.0}
	epsilon := 0.001

	gridG, err := compute.NewAdaptiveSparseGrid(g, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create gridG: %v", err)
	}

	gridFG, err := gridG.MakeNextIteration(f, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create gridFG: %v", err)
	}

	testPoints := []compute.Point{
		{0.0}, {2.0}, {1.0}, {0.5}, {1.5}, {0.123}, {1.987},
	}
	for i := 0; i < 10; i++ {
		testPoints = append(testPoints, compute.Point{min[0] + rand.Float64()*(max[0]-min[0])})
	}

	for _, testPoint := range testPoints {
		result, err := gridFG.Evaluate(testPoint)
		if err != nil {
			t.Fatalf("failed to evaluate at %v: %v", testPoint, err)
		}
		expected := f(g(testPoint))

		if math.Abs(result[0]-expected[0]) > epsilon*2 {
			t.Errorf("at %v: expected %v, got %v", testPoint, expected[0], result[0])
		}
	}
}
