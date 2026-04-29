package compute_test

import (
	"fmt"
	"math"

	"Diplom/compute"
)

// 1. Простой вводный пример
// Проверяется базовая работа алгоритма на простой одномерной функции.
// Функция: f(x) = x^2
func ExampleAdaptiveSparseGrid_simple() {
	funcEval := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] * arg[0]}
	}

	min := compute.Point{-2.0}
	max := compute.Point{2.0}

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := compute.Point{1.5}
	result, _ := grid.Evaluate(testPoint)
	fmt.Printf("f(1.5) = %.2f (expected: 2.25)\n", result[0])
	// Output: f(1.5) = 2.25 (expected: 2.25)
}

// 2. Пример многомерной относительно простой функции
// Проверяется работа алгоритма на функции от 3 переменных.
// Функция: f(x, y, z) = x*y + sin(z)
func ExampleAdaptiveSparseGrid_multidim() {
	funcEval := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0]*arg[1] + math.Sin(arg[2])}
	}

	min := compute.Point{0.0, 0.0, 0.0}
	max := compute.Point{2.0, 2.0, math.Pi}

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.05, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := compute.Point{1.0, 1.5, math.Pi / 2.0}
	result, _ := grid.Evaluate(testPoint)
	fmt.Printf("f(1.0, 1.5, pi/2) = %.1f (expected: 2.5)\n", result[0])
	// Output: f(1.0, 1.5, pi/2) = 2.5 (expected: 2.5)
}

// 3. Пример дифференциального уравнения (разных размерностей)
// Проверяется интерполяция решения системы дифференциальных уравнений (модель Лотки-Вольтерры).
// dx/dt = alpha*x - beta*x*y
// dy/dt = delta*x*y - gamma*y
// Интерполируем значения (x(T), y(T)) в зависимости от начальных условий (x0, y0).
func ExampleAdaptiveSparseGrid_ode() {
	diffEq := func(state compute.Point, t float64) compute.Point {
		alpha := 2.0 / 3.0
		beta := 4.0 / 3.0
		delta := 1.0
		gamma := 1.0
		x := state[0]
		y := state[1]
		return compute.Point{alpha*x - beta*x*y, delta*x*y - gamma*y}
	}

	funcEval := func(initialState compute.Point) compute.Point {
		return integrateRk4(diffEq, initialState, 0.0, 2.0, 100)
	}

	min := compute.Point{0.5, 0.5}
	max := compute.Point{2.0, 2.0}

	grid, _ := compute.NewAdaptiveSparseGrid(funcEval, min, max, 0.05, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := compute.Point{1.0, 1.0}
	result, _ := grid.Evaluate(testPoint)
	expected := funcEval(testPoint)

	fmt.Printf("Interpolated close to expected: %v\n", math.Abs(result[0]-expected[0]) < 0.1 && math.Abs(result[1]-expected[1]) < 0.1)
	// Output: Interpolated close to expected: true
}

// 4. Пример композиции функций при помощи make_next_iteration
// Проверяется создание композиции функций f(g(x)).
// g(x) = x + 1, f(x) = x^2. Композиция: f(g(x)) = (x+1)^2.
func ExampleAdaptiveSparseGrid_composition() {
	g := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] + 1.0}
	}

	f := func(arg compute.Point) compute.Point {
		return compute.Point{arg[0] * arg[0]}
	}

	min := compute.Point{0.0}
	max := compute.Point{2.0}

	gridG, _ := compute.NewAdaptiveSparseGrid(g, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridFG, _ := gridG.MakeNextIteration(f, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := compute.Point{1.0}
	result, _ := gridFG.Evaluate(testPoint)
	fmt.Printf("f(g(1.0)) = %.1f (expected: 4.0)\n", result[0])
	// Output: f(g(1.0)) = 4.0 (expected: 4.0)
}

// 5. Пример дифференциального уравнения с make_next_iteration
// Проверяется последовательное интегрирование дифференциального уравнения.
// На каждом шаге уравнение интегрируется от t до t+1.
// dx/dt = -0.5 * x. Аналитическое решение: x(t) = x0 * exp(-0.5 * t).
func ExampleAdaptiveSparseGrid_odeComposition() {
	diffEq := func(state compute.Point, t float64) compute.Point {
		return compute.Point{-0.5 * state[0]}
	}

	integrate1s := func(state compute.Point) compute.Point {
		return integrateRk4(diffEq, state, 0.0, 1.0, 50)
	}

	min := compute.Point{0.0}
	max := compute.Point{10.0}

	gridT1, _ := compute.NewAdaptiveSparseGrid(integrate1s, min, max, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridT2, _ := gridT1.MakeNextIteration(integrate1s, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)
	gridT3, _ := gridT2.MakeNextIteration(integrate1s, 0.01, nil, compute.BasisTypeQuadratic, compute.BuildTypeParallel, 0, 0)

	testPoint := compute.Point{8.0}
	resultT1, _ := gridT1.Evaluate(testPoint)
	resultT2, _ := gridT2.Evaluate(testPoint)
	resultT3, _ := gridT3.Evaluate(testPoint)

	fmt.Printf("State at t=1 close to expected: %v\n", math.Abs(resultT1[0]-8.0*math.Exp(-0.5*1.0)) < 0.05)
	fmt.Printf("State at t=2 close to expected: %v\n", math.Abs(resultT2[0]-8.0*math.Exp(-0.5*2.0)) < 0.05)
	fmt.Printf("State at t=3 close to expected: %v\n", math.Abs(resultT3[0]-8.0*math.Exp(-0.5*3.0)) < 0.05)
	// Output:
	// State at t=1 close to expected: true
	// State at t=2 close to expected: true
	// State at t=3 close to expected: true
}
