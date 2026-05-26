//
// Created by ivanl on 01.03.2026.
//
module;

#include <iostream>
#include <array>
#include <vector>
#include <numeric>
#include <cmath>
#include <valarray>
#include <unordered_map>
#include <unordered_set>
#include <optional>
#include <limits>
#include <stdexcept>
#include <cstddef>
#include <algorithm>
#include <queue>
#include <type_traits>
#include <initializer_list>
#include <ranges>
#include <future>
#include <functional>
#include <thread>
#include <fmt/core.h>
#include <chrono>

export module util;

export namespace std
{
    using std::cout;
    using std::cerr;
    using std::endl;
    using std::array;
    using std::vector;
    using std::move;
    using std::sqrt;
    using std::inner_product;
    using std::out_of_range;
    using std::weak_ordering;
    using std::abs;
    using std::sin;
    using std::cos;
    using std::pow;
    using std::max;
    using std::same_as;
    using std::unordered_map;
    using std::unordered_set;
    using std::optional;
    using std::numeric_limits;
    using std::runtime_error;
    using std::string;
    using std::operator>>;
    using std::operator<<;
    using std::initializer_list;
    namespace ranges
    {
        using std::ranges::next_permutation;
        using std::ranges::next_permutation_result;
        using std::ranges::range;
    }
    using std::queue;
    using std::remove_cvref_t;
    using std::reference_wrapper;
    using std::operator==;
    using std::operator|;
    using std::hash;
    using std::pair;
    using std::async;
    using std::launch;
    using std::future;
    using std::void_t;
    using std::mutex;
    using std::lock_guard;
    using std::ref;
    using std::atomic_flag;
    namespace this_thread
    {
        using std::this_thread::yield;
    }
    using std::swap;
    namespace chrono {
        using std::chrono::high_resolution_clock;
    }
}

export namespace std::inline literals::inline chrono_literals {
  // [time.duration.literals], suffixes for duration literals
  using std::literals::chrono_literals::operator""h;
  using std::literals::chrono_literals::operator""min;
  using std::literals::chrono_literals::operator""s;
  using std::literals::chrono_literals::operator""ms;
  using std::literals::chrono_literals::operator""us;
  using std::literals::chrono_literals::operator""ns;

  // [using std::literals::chrono_literals::.cal.day.nonmembers], non-member functions
  using std::literals::chrono_literals::operator""d;

  // [using std::literals::chrono_literals::.cal.year.nonmembers], non-member functions
  using std::literals::chrono_literals::operator""y;
} // namespace std::inline literals::inline chrono_literals

export namespace fmt
{
    using fmt::format;
}

export {
    using std::size_t;
}