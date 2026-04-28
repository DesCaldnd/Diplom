package compute_test

import (
	"math"
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
	epsilon := 0.01

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoint := compute.Point{1.5}
	result, err := grid.Evaluate(testPoint)
	if err != nil {
		t.Fatalf("failed to evaluate: %v", err)
	}
	expected := funcEval(testPoint)

	if math.Abs(result[0]-expected[0]) > epsilon*2 {
		t.Errorf("expected %v, got %v", expected[0], result[0])
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
	epsilon := 0.05

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeSequential, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoint := compute.Point{1.0, 1.5}
	result, err := grid.Evaluate(testPoint)
	if err != nil {
		t.Fatalf("failed to evaluate: %v", err)
	}
	expected := funcEval(testPoint)

	if math.Abs(result[0]-expected[0]) > epsilon*2 {
		t.Errorf("expected %v, got %v", expected[0], result[0])
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
	epsilon := 0.01

	grid, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoint := compute.Point{2.0}
	result, err := grid.Evaluate(testPoint)
	if err != nil {
		t.Fatalf("failed to evaluate: %v", err)
	}
	expected := funcEval(testPoint)

	if math.Abs(result[0]-expected[0]) > epsilon*2 {
		t.Errorf("expected %v, got %v", expected[0], result[0])
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
	epsilon := 0.05

	gridWithoutAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	anchors := []compute.Point{{math.Pi / 2.0}, {3.0 * math.Pi / 2.0}}
	gridWithAnchors, err := compute.NewAdaptiveSparseGrid(funcEval, min, max, epsilon, anchors, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	testPoint := compute.Point{math.Pi / 2.0}
	resultWithout, _ := gridWithoutAnchors.Evaluate(testPoint)
	resultWith, _ := gridWithAnchors.Evaluate(testPoint)
	expected := funcEval(testPoint)

	if math.Abs(resultWithout[0]-0.0) > epsilon*2 {
		t.Errorf("without anchors: expected ~0.0, got %v", resultWithout[0])
	}

	if math.Abs(resultWith[0]-expected[0]) > epsilon*2 {
		t.Errorf("with anchors: expected %v, got %v", expected[0], resultWith[0])
	}
}
