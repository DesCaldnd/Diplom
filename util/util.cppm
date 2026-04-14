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

export module util;

export namespace std
{
    using std::cout;
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
        using std::ranges::operator|;
        using std::ranges::transform;
        using std::ranges::transform_view;
        using std::ranges::range;

        namespace views 
        {
            using std::ranges::views::transform;
        }
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
}

export {
    using std::size_t;
}