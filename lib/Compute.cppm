//
// Created by ivanl on 01.03.2026.
//
// добавить причину для якоря
// асинхронка
module;

export module compute;

import util;
import function_traits;
export import vector_value;
import helper_types;

export namespace Compute
{
    using Compute::BasisType;

    template <size_t IN_DIM, size_t OUT_DIM>
    class AdaptiveSparseGrid
    {
    public:
        using EntryPoint = EntryPoint<IN_DIM, OUT_DIM>;
        using GridKey = GridKey<IN_DIM>;
        using Node = Node<IN_DIM, OUT_DIM>;
        using Argument = Point<IN_DIM>;
        using Answer = Point<OUT_DIM>;
        using Point = Point<IN_DIM>;
        using Basis = Basis<IN_DIM>;
        using NodeMap = std::unordered_map<GridKey, Node>;
        using NodeSet = std::unordered_set<GridKey>;

    private:
        Point min_;
        Point max_;
        BasisType basis_type_;

        std::vector<EntryPoint> entry_points_;
        NodeMap nodes_;

        template <typename FUNC>
        class NextIterationWrapper
        {
            const FUNC &func_;
            const AdaptiveSparseGrid &grid_;

        public:
            NextIterationWrapper(const FUNC &func, const AdaptiveSparseGrid &grid) : func_(func), grid_(grid) {}

            auto operator()(const Argument &x) const
            {
                return func_(grid_.evaluate(x));
            }
        };

    public:
        template <function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        AdaptiveSparseGrid(const FUNC &func, Point min, Point max, ScalarType epsilon, const std::vector<Argument>& anchor_points = {},
                           BasisType basis_type = BasisType::QUADRATIC) : min_(min), max_(max), basis_type_(basis_type)
        {
            entry_points_.reserve(pow_3(IN_DIM));
            build(func, epsilon, calculate_original_function_at_points(func, anchor_points));
        }

        Answer evaluate(const Argument &x) const
        {
            check_evaluation_point(x);
            return evaluate_for_dim(to_unit(x), IN_DIM + 1);
        }

        template <function_traits::ArrayOpIn<ScalarType, OUT_DIM> FUNC>
        AdaptiveSparseGrid make_next_iteration(const FUNC &func, ScalarType epsilon, std::optional<std::vector<Argument>> anchor_points = {}, std::optional<BasisType> basis_type = {})
        {
            return AdaptiveSparseGrid(NextIterationWrapper<FUNC>(func, *this), min_, max_, epsilon, anchor_points.value_or({}), basis_type.value_or(basis_type_));
        }

    private:
        Point to_real(const Point &unit) const
        {
            Point real;
            for (size_t i = 0; i < IN_DIM; ++i)
                real[i] = min_[i] + unit[i] * (max_[i] - min_[i]);
            return real;
        }

        Point to_unit(const Point &real) const
        {
            Point unit;
            for (size_t i = 0; i < IN_DIM; ++i)
                unit[i] = (real[i] - min_[i]) / (max_[i] - min_[i]);
            return unit;
        }

        template <function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        std::vector<std::pair<Argument, Answer>> calculate_original_function_at_points(const FUNC& func, const std::vector<Argument>& anchor_points)
        {
            auto r = anchor_points | std::ranges::views::transform([&, this](const Argument& arg) { return std::pair{to_unit(arg), func(arg)}; });
            return {r.begin(), r.end()};
        }

        void check_evaluation_point(const Argument &real) const
        {
            for (size_t i = 0; i < IN_DIM; ++i)
            {
                if (real[i] < min_[i] || real[i] > max_[i])
                    throw std::runtime_error("Evaluation point out of bounds");
            }
        }

        static size_t pow_3(size_t exp)
        {
            size_t result = 1;
            size_t base = 3;

            while (exp)
            {
                if (exp & 1u)
                {
                    result *= base;
                }
                base *= base;
                exp >>= 1;
            }

            return result;
        }

        template <function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        void build(const FUNC &func, ScalarType epsilon, const std::vector<std::pair<Argument, Answer>>& anchors)
        {
            for (size_t i = 0; i <= IN_DIM; ++i)
            {
                GridKey key;
                key.level = 0;
                key.index = 0;

                for (size_t j = 0; j < i; ++j)
                {
                    key.level[IN_DIM - j - 1] = 1;
                }

                do
                {
                    size_t index_permutations = 1zu << (IN_DIM - i);

                    for (size_t j = 0; j < index_permutations; ++j)
                    {
                        expand_indicies(key, j);
                        EntryPoint entry_point;
                        entry_point.dimensions = i;
                        Node &node = nodes_[create_node(func, key, {}, i).value()];
                        build_grid(func, epsilon, anchors, node, i);
                        entry_point.node = node;
                        entry_points_.push_back(entry_point);
                    }
                } while (std::ranges::next_permutation(key.level.coords).found);
            }
        }

        static void expand_indicies(GridKey &key, size_t permutation)
        {
            size_t index = 0;
            for (size_t i = 0; i < IN_DIM; ++i)
            {
                if (key.level[i] == 1)
                    key.index[i] = 1;
                else
                {
                    key.index[i] = (permutation & (1u << index)) ? 1 : 0;
                    ++index;
                }
            }
        }

        template <function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        void build_grid(const FUNC &func, ScalarType epsilon, const std::vector<std::pair<Argument, Answer>>& anchors, const Node &entry_point, size_t dimension)
        {
            std::queue<GridKey> node_queue;
            node_queue.push(entry_point.key);

            while (!node_queue.empty())
            {
                Node &node = nodes_[node_queue.front()];
                node_queue.pop();

                if (node.alpha.length() <= epsilon)
                {
                    bool can_continue = true;

                    for (auto& anchor : anchors)
                    {
                        if (node.is_point_in_affect_zone(anchor.first) && (evaluate_for_dim_and_entry_point(node.center_unit, dimension, entry_point) - anchor.second).length() >= epsilon)
                        {
                            can_continue = false;
                            break;
                        }
                    }

                    if (can_continue)
                    {
                        continue;
                    }
                }

                node.has_children = true;

                for (size_t i = 0; i < IN_DIM; ++i)
                {
                    if (node.key.level[i] == 0)
                        continue;

                    GridKey key_left = node.key, key_right = node.key;

                    key_left.level[i] = key_right.level[i] = node.key.level[i] + 1;

                    key_left.index[i] = 2 * node.key.index[i] - 1;
                    key_right.index[i] = 2 * node.key.index[i] + 1;

                    auto left_node = create_node(func, key_left, entry_point, dimension);
                    auto right_node = create_node(func, key_right, entry_point, dimension);

                    if (left_node.has_value())
                        node_queue.push(left_node.value());
                    if (right_node.has_value())
                        node_queue.push(right_node.value());
                }
            }
        }

        Answer evaluate_for_dim_and_entry_point(const Argument &x, size_t max_grid_dim, std::optional<Node> entry_point)
        {
            Answer interp = evaluate_for_dim(x, max_grid_dim);
            if (entry_point.has_value())
            {
                NodeSet processed_nodes;
                interp += evaluate_recursive(x, entry_point.value(), processed_nodes);
            }
            return interp;
        }

        template <function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        std::optional<GridKey> create_node(const FUNC &func, const GridKey &key, std::optional<Node> entry_point, size_t dimension)
        {
            if (nodes_.contains(key))
                return {};

            Node node;
            node.key = key;
            for (size_t i = 0; i < IN_DIM; ++i)
                node.center_unit[i] = Basis::get_coord(key.level[i], key.index[i]);

            Argument real = to_real(node.center_unit);
            Answer etalon = func(real.coords);
            Answer interp = evaluate_for_dim_and_entry_point(node.center_unit, dimension, entry_point);

            node.alpha = etalon - interp;

            nodes_[key] = node;
            return key;
        }

        Answer evaluate_for_dim(const Argument &x, size_t max_grid_dim = IN_DIM + 1) const
        {
            Answer answer = 0;
            for (auto &&entry_point : entry_points_)
            {
                if (entry_point.dimensions >= max_grid_dim)
                {
                    break;
                }
                NodeSet processed_nodes;
                answer += evaluate_recursive(x, entry_point.node, processed_nodes);
            }
            return answer;
        }

        Answer evaluate_recursive(const Argument &x, const Node &node, NodeSet &processed_nodes) const
        {
            if (processed_nodes.contains(node.key))
                return 0;

            processed_nodes.emplace(node.key);
            auto basis = Basis::eval(x, node.key.level, node.key.index, basis_type_);

            if (std::abs(basis) < std::numeric_limits<ScalarType>::epsilon())
            {
                return 0;
            }

            auto value = node.alpha * basis;

            if (node.has_children)
            {
                for (size_t i = 0; i < IN_DIM; ++i)
                {
                    if (node.key.level[i] != 0)
                    {
                        auto child = get_child_for_dim_and_arg(node, x, i);
                        if (child.has_value())
                        {
                            value += evaluate_recursive(x, child.value(), processed_nodes);
                        }
                    }
                }
            }
            return value;
        }

        std::optional<Node> get_child_for_dim_and_arg(const Node &parent, const Argument &x, size_t dimension) const
        {
            GridKey key = parent.key;
            key.level[dimension] += 1;
            key.index[dimension] = 2 * key.index[dimension] + (x[dimension] < parent.center_unit[dimension] ? -1 : 1);

            auto child = nodes_.find(key);
            if (child != nodes_.end())
                return child->second;
            else
                return {};
        }
    };

    template <size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, Point<IN_DIM> min, Point<IN_DIM> max, ScalarType epsilon, const std::vector<Point<IN_DIM>>& anchor_points = {}, BasisType type = BasisType::QUADRATIC) -> AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template <size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, std::array<ScalarType, IN_DIM> min, std::array<ScalarType, IN_DIM> max, ScalarType epsilon, const std::vector<Point<IN_DIM>>& anchor_points = {}, BasisType type = BasisType::QUADRATIC) -> AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template <size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, Point<IN_DIM> min, std::array<ScalarType, IN_DIM> max, ScalarType epsilon, const std::vector<Point<IN_DIM>>& anchor_points = {}, BasisType type = BasisType::QUADRATIC) -> AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template <size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, std::array<ScalarType, IN_DIM> min, Point<IN_DIM> max, ScalarType epsilon, const std::vector<Point<IN_DIM>>& anchor_points = {}, BasisType type = BasisType::QUADRATIC) -> AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;
}
