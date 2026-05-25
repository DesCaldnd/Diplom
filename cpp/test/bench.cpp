#include <fmt/format.h>
#include <gtest/gtest.h>

import compute;
import util;

using namespace std::chrono_literals;

// Простой интегратор методом Рунге-Кутты 4-го порядка
template <size_t DIM, typename FUNC>
Compute::Point<DIM> rk4_step(const FUNC &f, Compute::Point<DIM> x, double t,
                             double dt) {
  auto k1 = f(x, t);
  auto k2 = f(x + k1 * (dt / 2.0), t + dt / 2.0);
  auto k3 = f(x + k2 * (dt / 2.0), t + dt / 2.0);
  auto k4 = f(x + k3 * dt, t + dt);
  return x + (k1 + k2 * 2.0 + k3 * 2.0 + k4) * (dt / 6.0);
}

template <size_t DIM, typename FUNC>
Compute::Point<DIM> integrate_rk4(const FUNC &f, Compute::Point<DIM> x0,
                                  double t0, double t1, int steps) {
  double dt = (t1 - t0) / steps;
  Compute::Point<DIM> x = x0;
  double t = t0;
  for (int i = 0; i < steps; ++i) {
    x = rk4_step<DIM>(f, x, t, dt);
    t += dt;
  }
  return x;
}

// Вспомогательная функция для замера времени
template <typename Func> double measure_time(Func &&f) {
  auto start = std::chrono::high_resolution_clock::now();
  f();
  auto end = std::chrono::high_resolution_clock::now();
  std::chrono::duration<double> diff = end - start;
  return diff.count();
}

template <typename Func> double measure_time_n_times(size_t n, Func &&f) {
  double res = 0;
  for (size_t i = 0; i < n; i++) {
    res += measure_time(f);
  }

  return res / n;
}

template <typename Func>
void bench(const std::string &name, size_t n, Func &&f) {
  std::cout << "Benchmark: " << name
            << ", time: " << static_cast<size_t>(measure_time_n_times(n, f) * 1000 * 1000'000)
            << "ns" << std::endl;
}

constexpr Compute::ScalarType Pi =
    3.14159265358979323846264338327950288419716939937510582097494459;
constexpr Compute::BuildType DefaultBuildType = Compute::BuildType::SEQUENTIAL;

TEST(ComputeBenchmark, DomainScalingDifferentDimensions) {
  std::vector<Compute::ScalarType> scales{1, 2, 4};

  for (auto scale : scales) {
    Compute::Point<1> min = 0;
    Compute::Point<1> max = Pi * scale;

    bench(fmt::format("dim_1_scale_{}*pi", scale), 40, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(
          [](Compute::Point<1> arg) {
            return Compute::Point<1>{std::sin(arg[0]) * std::cos(arg[0] * arg[0] / 2) + std::tan(0.1 * arg[0])};
          },
          min, max, 0.001, {}, Compute::BasisType::QUADRATIC,
          DefaultBuildType);
    });
  }

  for (auto scale : scales) {
    Compute::Point<2> min = 0;
    Compute::Point<2> max = Pi * scale;

    bench(fmt::format("dim_2_scale_{}*pi", scale), 40, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(
          [](Compute::Point<2> arg) {
            return Compute::Point<1>{std::sin(arg[0]) * std::cos(arg[0] * arg[0] / 2) +  std::tan(0.1 * arg[0])+
                                     std::sin(arg[1] * 2) *
                                         std::cos(arg[0] * arg[1] / 3) + std::tan(0.1 * arg[1])};
          },
          min, max, 0.001, {}, Compute::BasisType::QUADRATIC,
          DefaultBuildType);
    });
  }

  for (auto scale : scales) {
    Compute::Point<3> min = 0;
    Compute::Point<3> max = Pi * scale;

    bench(fmt::format("dim_3_scale_{}*pi", scale), 50, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(
          [](Compute::Point<3> arg) {
            return Compute::Point<1>{
                std::sin(arg[0]) * std::cos(arg[0] * arg[0] / 2) + std::tan(0.1 * arg[0]) +
                std::sin(arg[1] * 2) * std::cos(arg[0] * arg[1] / 3) + std::tan(0.1 * arg[1]) +
                std::sin(arg[2] * 3) * std::cos(arg[0] * arg[2] / 4) + std::tan(0.1 * arg[2])};
          },
          min, max, 0.001, {}, Compute::BasisType::QUADRATIC,
          DefaultBuildType);
    });
  }
}

TEST(ComputeBenchmark, ParallelVsSequentialBuild) {
  auto funcEval = [](Compute::Point<2> arg) {
    auto x = arg[0], y = arg[1];
    return Compute::Point<1>{std::sin(3 * x) * std::cos(25 * y) + 0.2 * std::sin(7 * x + 2.5) - 0.15 * std::cos(11 * y - 1.0)};
  };

  Compute::Point<2> min = 0, max = 2 * Pi;

  std::vector<Compute::BuildType> build_types{
                                              Compute::BuildType::SEQUENTIAL, Compute::BuildType::PARALLEL};
  std::vector<Compute::BasisType> basis_types{
                                              Compute::BasisType::LINEAR, Compute::BasisType::QUADRATIC};

  for (auto build_type : build_types) {
    for (auto basis_type : basis_types) {
      std::string basis_name, build_name;

      switch (build_type) {
      case Compute::BuildType::PARALLEL:
        build_name = "PARALLEL";
        break;
      default:
        build_name = "SEQUENTIAL";
      }
      switch (basis_type) {
      case Compute::BasisType::QUADRATIC:
        basis_name = "QUADRATIC";
        break;
      default:
        basis_name = "LINEAR";
      }

      bench(fmt::format("{}_{}", build_name, basis_name), 50, [&]() {
        volatile Compute::AdaptiveSparseGrid grid(funcEval, min, max, 0.001, {},
                                                  basis_type, build_type);
      });
    }
  }
}

TEST(ComputeBenchmark, LinearVsQuadraticBasis) {
  auto funcEval = [](Compute::Point<2> arg) {
    auto x = arg[0], y = arg[1];
    return Compute::Point<1>{std::sin(15 * x) + 0.6 * std::cos(40 * y) +
                             0.2 * std::sin(31 * x + 0.3) -
                             0.12 * std::cos(19 * y - 0.7) + std::sin(x * y)};
  };

  Compute::Point<2> min = 0, max = Pi;

  std::vector<Compute::ScalarType> eps_vals{1e-2, 5e-3, 1e-3, 1e-4, 1e-5, 5e-6, 1e-6};

  for (auto eps : eps_vals) {
    bench(fmt::format("linear_eps_{}", eps), 50, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(funcEval, min, max, eps, {},
                                                Compute::BasisType::LINEAR, DefaultBuildType);
    });
    bench(fmt::format("quadratic_eps_{}", eps), 50, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(funcEval, min, max, eps, {},
                                                Compute::BasisType::QUADRATIC, DefaultBuildType);
    });
  }
}

TEST(ComputeBenchmark, DifferentialEquationApproaches) {
  auto diff_eq = [](Compute::Point<2> state, Compute::ScalarType t) {
    auto x = state[0], y = state[1];
    return Compute::Point<2>{-y / (1 + std::sqrt(x * x + y * y)),
                             -x / (1 + std::sqrt(x * x + y * y))};
  };

  std::vector<size_t> t_max_values{2, 5, 10, 20};

  for (auto t_max : t_max_values) {
    auto func_with_t = [&](Compute::Point<3> arg) {
      auto x0 = arg[0], y0 = arg[1], t = arg[2];
      return integrate_rk4(diff_eq, Compute::Point<2>{x0, y0}, 0, t, 50);
    };

    auto integrate1s = [&](Compute::Point<2> state) {
      return integrate_rk4(diff_eq, Compute::Point<2>{state[0], state[1]}, 0, 1, 25);
    };

    bench(fmt::format("t_as_interval_uncertainty_tMax_{}", t_max), 40, [&]() {
      Compute::Point<3> min = {-1, 0, 0},
                        max = {1, 1, static_cast<Compute::ScalarType>(t_max)};
      volatile Compute::AdaptiveSparseGrid grid(
          func_with_t, min, max, 0.001, {}, Compute::BasisType::QUADRATIC,
          DefaultBuildType);
    });

    bench(fmt::format("iterative_make_next_iteration_tMax_{}", t_max), 40,
          [&]() {
            Compute::Point<2> min = {-1, 0}, max = {1, 1};
            Compute::AdaptiveSparseGrid grid(integrate1s, min, max, 0.001, {},
                                             Compute::BasisType::QUADRATIC,
                                             DefaultBuildType);

            for (size_t step = 1; step < t_max; ++step) {
              grid = grid.make_next_iteration(integrate1s, 0.001, {},
                                              Compute::BasisType::QUADRATIC,
                                              DefaultBuildType);
            }
          });
  }
}

TEST(ComputeBenchmark, NodeLimitOptimization) {
  auto func_eval = [](Compute::Point<2> arg) {
    auto x = arg[0], y = arg[1];
    return Compute::Point<1>{std::sin(5 * x) + std::cos(33 * y)};
  };

  Compute::Point<2> min = 0.75 * Pi, max = 2 * Pi;

  std::vector<size_t> node_limits{0, 200, 500, 1000};

  for (auto limit : node_limits) {
    bench(fmt::format("limit_{}", limit), 50, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(
          func_eval, min, max, 0.001, {}, Compute::BasisType::QUADRATIC,
          DefaultBuildType, 0, limit);
    });
  }
}

TEST(ComputeBenchmark, EvaluationCostAfterBuild) {
  auto func_eval = [](Compute::Point<2> arg) {
    auto x = arg[0], y = arg[1];
    return Compute::Point<1>{std::sin(15 * x) + 0.6 * std::cos(40 * y) +
                             0.2 * std::sin(31 * x + 0.3) -
                             0.12 * std::cos(19 * y - 0.7) + std::sin(x * y)};
  };

  Compute::Point<2> min = 0, max = Pi, test_point{0.73, 1.11};

  Compute::AdaptiveSparseGrid linear_grid(func_eval, min, max, 0.001, {},
                                          Compute::BasisType::LINEAR,
                                          Compute::BuildType::PARALLEL);
  Compute::AdaptiveSparseGrid quadratic_grid(func_eval, min, max, 0.001, {},
                                             Compute::BasisType::QUADRATIC,
                                             Compute::BuildType::PARALLEL);

  bench("evaluate_linear_grid", 10, [&]() {
    for (size_t i = 0; i < 10000; ++i) {
      volatile Compute::Point<1> point = linear_grid.evaluate(test_point);
    }
  });

  bench("evaluate_quadratic_grid", 15, [&]() {
    for (size_t i = 0; i < 10000; ++i) {
      volatile Compute::Point<1> point = quadratic_grid.evaluate(test_point);
    }
  });
}

TEST(ComputeBenchmark, FourDimensionalBuildTypeBasisProduct) {
  auto func_eval = [](Compute::Point<4> arg) {
    auto x1 = arg[0], x2 = arg[1], x3 = arg[2], x4 = arg[3];
    return Compute::Point<1>{std::sin(30 * x1) * std::cos(22.2323 * x2) +
                             std::sin(70 * x3) - 0.15 * std::sin(x4) +
                             std::cos(x1 * x3)};
  };

  struct test_case {
    Compute::BuildType build_type;
    Compute::BasisType basis_type;
    std::string name;
  };
  std::vector<test_case> test_cases = {
    {Compute::BuildType::PARALLEL, Compute::BasisType::LINEAR, "parallel_linear"},
    {Compute::BuildType::PARALLEL, Compute::BasisType::QUADRATIC, "parallel_quadratic"},
    {Compute::BuildType::SEQUENTIAL, Compute::BasisType::LINEAR, "sequential_linear"},
    {Compute::BuildType::SEQUENTIAL, Compute::BasisType::QUADRATIC, "sequential_quadratic"},
  };
  Compute::ScalarType epsilon = 0.0000005;
  Compute::Point<4> min = 0, max = Pi;

  for (auto test_case : test_cases)
  {
    bench(test_case.name, 10, [&]() {
      volatile Compute::AdaptiveSparseGrid grid(func_eval, min, max, epsilon, {}, test_case.basis_type, test_case.build_type);
    });
  }
}

int main(int argc, char **argv) {
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}