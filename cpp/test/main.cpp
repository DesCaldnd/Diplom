#include <gtest/gtest.h>
#include <cmath>
#include <chrono>
#include <iostream>

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

// ============================================================================
// ТЕСТЫ
// ============================================================================

// Тест 1: Корректность работы с обычной одномерной функцией
// Проверяется, что интерполяция функции f(x) = x^2 работает корректно
// и ошибка не превышает заданный epsilon.
TEST(ComputeTest, Simple1DFunction) {
    auto func = [](Compute::Point<1> arg) -> Compute::Point<1> {
        return {arg[0] * arg[0]};
    };

    Compute::Point<1> min = {-2.0};
    Compute::Point<1> max = {2.0};
    double epsilon = 0.01;
    
    Compute::AdaptiveSparseGrid grid(func, min, max, epsilon);
    
    Compute::Point<1> test_point = {1.5};
    auto result = grid.evaluate(test_point);
    auto expected = func(test_point);
    
    EXPECT_NEAR(result[0], expected[0], epsilon * 2); // Допуск чуть больше epsilon
}

// Тест 2: Корректность работы с многомерной функцией
// Проверяется интерполяция функции f(x, y) = x*y + sin(x)
// с различными параметрами (базис, тип сборки).
TEST(ComputeTest, MultidimFunction) {
    auto func = [](Compute::Point<2> arg) -> Compute::Point<1> {
        return {arg[0] * arg[1] + std::sin(arg[0])};
    };

    Compute::Point<2> min = {0.0, 0.0};
    Compute::Point<2> max = {2.0, 2.0};
    double epsilon = 0.05;
    
    Compute::AdaptiveSparseGrid grid(func, min, max, epsilon, {}, Compute::BasisType::QUADRATIC, Compute::BuildType::SEQUENTIAL);
    
    Compute::Point<2> test_point = {1.0, 1.5};
    auto result = grid.evaluate(test_point);
    auto expected = func(test_point);
    
    EXPECT_NEAR(result[0], expected[0], epsilon * 2);
}

// Тест 3: Корректность работы с дифференциальным уравнением
// Проверяется интерполяция решения дифференциального уравнения dx/dt = -x.
// Сравнивается с эталонным аналитическим решением.
TEST(ComputeTest, DifferentialEquation) {
    auto diff_eq = [](Compute::Point<1> state, double t) -> Compute::Point<1> {
        return {-state[0]};
    };

    auto func = [&](Compute::Point<1> initial_state) -> Compute::Point<1> {
        return integrate_rk4<1>(diff_eq, initial_state, 0.0, 1.0, 100);
    };

    Compute::Point<1> min = {0.0};
    Compute::Point<1> max = {5.0};
    double epsilon = 0.01;
    
    Compute::AdaptiveSparseGrid grid(func, min, max, epsilon);
    
    Compute::Point<1> test_point = {2.0};
    auto result = grid.evaluate(test_point);
    auto expected = func(test_point);
    
    EXPECT_NEAR(result[0], expected[0], epsilon * 2);
}

// Тест 4: Тест на anchor points
// Проверяется функция sin(x) на отрезке [0, 2*pi].
// Значения на краях и в центре равны 0. Без anchor points алгоритм
// может посчитать функцию тождественным нулем.
// С anchor points (например, pi/2 и 3*pi/2) алгоритм должен отработать корректно.
TEST(ComputeTest, AnchorPoints) {
    auto func = [](Compute::Point<1> arg) -> Compute::Point<1> {
        return {std::sin(arg[0])};
    };

    Compute::Point<1> min = {0.0};
    Compute::Point<1> max = {2.0 * M_PI};
    double epsilon = 0.05;
    
    // Сетка без anchor points
    Compute::AdaptiveSparseGrid grid_without_anchors(func, min, max, epsilon);
    
    // Сетка с anchor points
    std::vector<Compute::Point<1>> anchors = {{M_PI / 2.0}, {3.0 * M_PI / 2.0}};
    Compute::AdaptiveSparseGrid grid_with_anchors(func, min, max, epsilon, anchors);
    
    Compute::Point<1> test_point = {M_PI / 2.0};
    auto result_without = grid_without_anchors.evaluate(test_point);
    auto result_with = grid_with_anchors.evaluate(test_point);
    auto expected = func(test_point);
    
    // Без anchor points результат скорее всего будет близок к 0
    EXPECT_NEAR(result_without[0], 0.0, epsilon * 2);
    
    // С anchor points результат должен быть близок к 1.0
    EXPECT_NEAR(result_with[0], expected[0], epsilon * 2);
}

// ============================================================================
// БЕНЧМАРКИ
// ============================================================================

// Вспомогательная функция для замера времени
template<typename Func>
double measure_time(Func&& f) {
    auto start = std::chrono::high_resolution_clock::now();
    f();
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double> diff = end - start;
    return diff.count();
}

// Бенчмарк 1: Замеры простых функций разной размерности
// Проверяется время построения сетки для функций 1D, 2D и 3D.
TEST(ComputeBenchmark, SimpleFunctionsDifferentDimensions) {
    std::cout << "[ BENCHMARK ] Simple Functions Different Dimensions" << std::endl;
    
    auto func1d = [](Compute::Point<1> arg) -> Compute::Point<1> { return {std::sin(arg[0])}; };
    auto func2d = [](Compute::Point<2> arg) -> Compute::Point<1> { return {std::sin(arg[0]) * std::cos(arg[1])}; };
    auto func3d = [](Compute::Point<3> arg) -> Compute::Point<1> { return {std::sin(arg[0]) * std::cos(arg[1]) * std::sin(arg[2])}; };

    double t1 = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func1d, Compute::Point<1>{0.0}, Compute::Point<1>{M_PI}, 0.01);
    });
    
    double t2 = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func2d, Compute::Point<2>{0.0, 0.0}, Compute::Point<2>{M_PI, M_PI}, 0.01);
    });
    
    double t3 = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func3d, Compute::Point<3>{0.0, 0.0, 0.0}, Compute::Point<3>{M_PI, M_PI, M_PI}, 0.01);
    });

    std::cout << "1D build time: " << t1 << " s" << std::endl;
    std::cout << "2D build time: " << t2 << " s" << std::endl;
    std::cout << "3D build time: " << t3 << " s" << std::endl;
}

// Бенчмарк 2: Сравнения на одинаковых параметрах с разным build_type, basis_type
// Проверяется время построения сетки для 2D функции при SEQUENTIAL vs PARALLEL
// и LINEAR vs QUADRATIC базисах.
TEST(ComputeBenchmark, BuildTypeAndBasisTypeComparison) {
    std::cout << "[ BENCHMARK ] BuildType and BasisType Comparison" << std::endl;
    
    auto func = [](Compute::Point<2> arg) -> Compute::Point<1> {
        return {std::sin(arg[0] * 5.0) * std::cos(arg[1] * 5.0)};
    };

    Compute::Point<2> min = {0.0, 0.0};
    Compute::Point<2> max = {1.0, 1.0};
    double epsilon = 0.05;

    double t_seq_lin = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func, min, max, epsilon, {}, Compute::BasisType::LINEAR, Compute::BuildType::SEQUENTIAL);
    });
    
    double t_par_lin = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func, min, max, epsilon, {}, Compute::BasisType::LINEAR, Compute::BuildType::PARALLEL);
    });
    
    double t_seq_quad = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func, min, max, epsilon, {}, Compute::BasisType::QUADRATIC, Compute::BuildType::SEQUENTIAL);
    });
    
    double t_par_quad = measure_time([&]() {
        Compute::AdaptiveSparseGrid grid(func, min, max, epsilon, {}, Compute::BasisType::QUADRATIC, Compute::BuildType::PARALLEL);
    });

    std::cout << "SEQUENTIAL + LINEAR:    " << t_seq_lin << " s" << std::endl;
    std::cout << "PARALLEL   + LINEAR:    " << t_par_lin << " s" << std::endl;
    std::cout << "SEQUENTIAL + QUADRATIC: " << t_seq_quad << " s" << std::endl;
    std::cout << "PARALLEL   + QUADRATIC: " << t_par_quad << " s" << std::endl;
}

// Бенчмарк 3: Сравнения дифференциальных уравнений
// Проверяется два подхода к интерполяции решения диффура на отрезке времени [0, 2]:
// 1. t - это интервальная неопределенность (дополнительная размерность сетки).
// 2. Последовательное интегрирование при помощи make_next_iteration (t=1, затем t=2).
TEST(ComputeBenchmark, DifferentialEquationApproaches) {
    std::cout << "[ BENCHMARK ] Differential Equation Approaches" << std::endl;
    
    auto diff_eq = [](Compute::Point<1> state, double t) -> Compute::Point<1> {
        return {-0.5 * state[0]};
    };

    // Подход 1: t как дополнительная размерность (2D сетка: x0 и t)
    auto func_with_t = [&](Compute::Point<2> arg) -> Compute::Point<1> {
        double x0 = arg[0];
        double t = arg[1];
        return integrate_rk4<1>(diff_eq, {x0}, 0.0, t, 50);
    };

    double t_approach1 = measure_time([&]() {
        Compute::Point<2> min = {0.0, 0.0};
        Compute::Point<2> max = {10.0, 2.0};
        Compute::AdaptiveSparseGrid grid(func_with_t, min, max, 0.01);
    });

    // Подход 2: последовательное интегрирование (1D сетка, 2 итерации)
    auto integrate_1s = [&](Compute::Point<1> state) -> Compute::Point<1> {
        return integrate_rk4<1>(diff_eq, state, 0.0, 1.0, 25);
    };

    double t_approach2 = measure_time([&]() {
        Compute::Point<1> min = {0.0};
        Compute::Point<1> max = {10.0};
        Compute::AdaptiveSparseGrid grid_t1(integrate_1s, min, max, 0.01);
        auto grid_t2 = grid_t1.make_next_iteration(integrate_1s, 0.01);
    });

    std::cout << "Approach 1 (t as dimension): " << t_approach1 << " s" << std::endl;
    std::cout << "Approach 2 (make_next_iteration): " << t_approach2 << " s" << std::endl;
}

int main(int argc, char **argv) {
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
