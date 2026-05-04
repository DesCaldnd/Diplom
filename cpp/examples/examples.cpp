#include <iostream>
#include <cmath>
#include <vector>

import compute;
import util;

// Простой интегратор методом Рунге-Кутты 4-го порядка
template<size_t DIM, typename FUNC>
Compute::Point<DIM> rk4_step(const FUNC& f, Compute::Point<DIM> x, double t, double dt) {
    auto k1 = f(x, t);
    auto k2 = f(x + k1 * (dt / 2.0), t + dt / 2.0);
    auto k3 = f(x + k2 * (dt / 2.0), t + dt / 2.0);
    auto k4 = f(x + k3 * dt, t + dt);
    return x + (k1 + k2 * 2.0 + k3 * 2.0 + k4) * (dt / 6.0);
}

template<size_t DIM, typename FUNC>
Compute::Point<DIM> integrate_rk4(const FUNC& f, Compute::Point<DIM> x0, double t0, double t1, int steps) {
    double dt = (t1 - t0) / steps;
    Compute::Point<DIM> x = x0;
    double t = t0;
    for (int i = 0; i < steps; ++i) {
        x = rk4_step<DIM>(f, x, t, dt);
        t += dt;
    }
    return x;
}

// 1. Простой вводный пример
// Проверяется базовая работа алгоритма на простой одномерной функции.
// Функция: f(x) = x^2
void example_simple() {
    std::cout << "--- Example 1: Simple 1D Function ---" << std::endl;
    auto func = [](Compute::Point<1> arg) -> Compute::Point<1> {
        return {arg[0] * arg[0]};
    };

    Compute::Point<1> min = {-2.0};
    Compute::Point<1> max = {2.0};
    
    Compute::AdaptiveSparseGrid grid(func, min, max, 0.01);
    
    Compute::Point<1> test_point = {1.5};
    auto result = grid.evaluate(test_point);
    std::cout << "f(1.5) = " << result[0] << " (expected: 2.25)" << std::endl;
}

// 2. Пример многомерной относительно простой функции
// Проверяется работа алгоритма на функции от 3 переменных.
// Функция: f(x, y, z) = x*y + sin(z)
void example_multidim() {
    std::cout << "\n--- Example 2: Multidimensional Function (3D) ---" << std::endl;
    auto func = [](Compute::Point<3> arg) -> Compute::Point<1> {
        return {arg[0] * arg[1] + std::sin(arg[2])};
    };

    Compute::Point<3> min = {0.0, 0.0, 0.0};
    Compute::Point<3> max = {2.0, 2.0, 3.14159};
    
    Compute::AdaptiveSparseGrid grid(func, min, max, 0.05);
    
    Compute::Point<3> test_point = {1.0, 1.5, 3.14159 / 2.0};
    auto result = grid.evaluate(test_point);
    std::cout << "f(1.0, 1.5, pi/2) = " << result[0] << " (expected: 2.5)" << std::endl;
}

// 3. Пример дифференциального уравнения (разных размерностей)
// Проверяется интерполяция решения системы дифференциальных уравнений (модель Лотки-Вольтерры).
// dx/dt = alpha*x - beta*x*y
// dy/dt = delta*x*y - gamma*y
// Интерполируем значения (x(T), y(T)) в зависимости от начальных условий (x0, y0).
void example_ode() {
    std::cout << "\n--- Example 3: Differential Equation (2D) ---" << std::endl;
    
    auto diff_eq = [](Compute::Point<2> state, double t) -> Compute::Point<2> {
        double alpha = 2.0 / 3.0;
        double beta = 4.0 / 3.0;
        double delta = 1.0;
        double gamma = 1.0;
        double x = state[0];
        double y = state[1];
        return {alpha * x - beta * x * y, delta * x * y - gamma * y};
    };

    // Функция, которая берет начальные условия и возвращает состояние в момент T=2.0
    auto func = [&](Compute::Point<2> initial_state) -> Compute::Point<2> {
        return integrate_rk4<2>(diff_eq, initial_state, 0.0, 2.0, 100);
    };

    Compute::Point<2> min = {0.5, 0.5};
    Compute::Point<2> max = {2.0, 2.0};
    
    Compute::AdaptiveSparseGrid grid(func, min, max, 0.05);
    
    Compute::Point<2> test_point = {1.0, 1.0};
    auto result = grid.evaluate(test_point);
    auto expected = func(test_point);
    std::cout << "State at T=2.0 for initial (1.0, 1.0):" << std::endl;
    std::cout << "Interpolated: x=" << result[0] << ", y=" << result[1] << std::endl;
    std::cout << "Expected:     x=" << expected[0] << ", y=" << expected[1] << std::endl;
}

// 4. Пример композиции функций при помощи make_next_iteration
// Проверяется создание композиции функций f(g(x)).
// g(x) = x + 1, f(x) = x^2. Композиция: f(g(x)) = (x+1)^2.
void example_composition() {
    std::cout << "\n--- Example 4: Function Composition with make_next_iteration ---" << std::endl;
    
    auto g = [](Compute::Point<1> arg) -> Compute::Point<1> {
        return {arg[0] + 1.0};
    };
    
    auto f = [](Compute::Point<1> arg) -> Compute::Point<1> {
        return {arg[0] * arg[0]};
    };

    Compute::Point<1> min = {0.0};
    Compute::Point<1> max = {2.0};
    
    // Сначала строим сетку для g(x)
    Compute::AdaptiveSparseGrid grid_g(g, min, max, 0.01);
    
    // Затем применяем f к результату g(x)
    auto grid_fg = grid_g.make_next_iteration(f, 0.01);
    
    Compute::Point<1> test_point = {1.0};
    auto result = grid_fg.evaluate(test_point);
    std::cout << "f(g(1.0)) = " << result[0] << " (expected: 4.0)" << std::endl;
}

// 5. Пример дифференциального уравнения с make_next_iteration
// Проверяется последовательное интегрирование дифференциального уравнения.
// На каждом шаге уравнение интегрируется от t до t+1.
// dx/dt = -0.5 * x. Аналитическое решение: x(t) = x0 * exp(-0.5 * t).
void example_ode_composition() {
    std::cout << "\n--- Example 5: ODE with make_next_iteration ---" << std::endl;
    
    auto diff_eq = [](Compute::Point<1> state, double t) -> Compute::Point<1> {
        return {-0.5 * state[0]};
    };

    // Функция интегрирования на 1 секунду (от t до t+1)
    // Поскольку уравнение автономное, t0 можно всегда брать 0.
    auto integrate_1s = [&](Compute::Point<1> state) -> Compute::Point<1> {
        return integrate_rk4<1>(diff_eq, state, 0.0, 1.0, 50);
    };

    Compute::Point<1> min = {0.0};
    Compute::Point<1> max = {10.0};
    
    // Сетка для t=1
    Compute::AdaptiveSparseGrid grid_t1(integrate_1s, min, max, 0.01);
    
    // Сетка для t=2 (интегрируем еще на 1 секунду)
    auto grid_t2 = grid_t1.make_next_iteration(integrate_1s, 0.01);
    
    // Сетка для t=3
    auto grid_t3 = grid_t2.make_next_iteration(integrate_1s, 0.01);
    
    Compute::Point<1> test_point = {8.0};
    auto result_t1 = grid_t1.evaluate(test_point);
    auto result_t2 = grid_t2.evaluate(test_point);
    auto result_t3 = grid_t3.evaluate(test_point);
    
    std::cout << "Initial state: 8.0" << std::endl;
    std::cout << "State at t=1: " << result_t1[0] << " (expected: " << 8.0 * std::exp(-0.5 * 1.0) << ")" << std::endl;
    std::cout << "State at t=2: " << result_t2[0] << " (expected: " << 8.0 * std::exp(-0.5 * 2.0) << ")" << std::endl;
    std::cout << "State at t=3: " << result_t3[0] << " (expected: " << 8.0 * std::exp(-0.5 * 3.0) << ")" << std::endl;
}

int main() {
    example_simple();
    example_multidim();
    example_ode();
    example_composition();
    example_ode_composition();
    return 0;
}
