//
// Created by ivanl on 01.03.2026.
//
module;

#include <proto/grid.pb.h>

export module compute;

import util;
import function_traits;
export import vector_value;
import helper_types;

export namespace Compute
{
    using Compute::BasisType;

    enum class BuildType
    {
        SEQUENTIAL,
        PARALLEL,
    };

    template<size_t IN_DIM, size_t OUT_DIM>
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

        template<typename FUNC>
        class NextIterationWrapper
        {
            const FUNC &func_;
            const AdaptiveSparseGrid &grid_;

        public:
            NextIterationWrapper(const FUNC &func, const AdaptiveSparseGrid &grid) : func_(func), grid_(grid)
            {
            }

            auto operator()(const Argument &x) const
            {
                return func_(grid_.evaluate(x));
            }
        };

    public:
        template<function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        AdaptiveSparseGrid(const FUNC &func, Point min, Point max, ScalarType epsilon,
                           const std::vector<Argument> &anchor_points = {},
                           BasisType basis_type = BasisType::QUADRATIC, BuildType build_type = BuildType::PARALLEL,
                           size_t max_level = 0, size_t max_nodes_in_grid = 0,
                           std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag = {}) : min_(min),
            max_(max), basis_type_(basis_type)
        {
            check_constraints();
            entry_points_.reserve(pow_3(IN_DIM));
            build(func, epsilon, calculate_original_function_at_points(func, anchor_points), build_type, max_level,
                  max_nodes_in_grid,
                  cancellation_flag);
        }

        Answer evaluate(const Argument &x) const
        {
            check_evaluation_point(x);
            return evaluate_for_dim(to_unit(x), IN_DIM + 1);
        }

        template<function_traits::ArrayOpIn<ScalarType, OUT_DIM> FUNC>
        AdaptiveSparseGrid make_next_iteration(const FUNC &func, ScalarType epsilon,
                                               std::optional<std::vector<Argument> > anchor_points = {},
                                               std::optional<BasisType> basis_type = {},
                                               std::optional<BuildType> build_type = {}, size_t max_level = 0,
                                               size_t max_nodes_in_grid = 0,
                                               std::optional<std::reference_wrapper<std::atomic_flag> >
                                               cancellation_flag = {})
        {
            return AdaptiveSparseGrid(NextIterationWrapper<FUNC>(func, *this), min_, max_, epsilon,
                                      anchor_points.value_or(std::vector<Argument>{}), basis_type.value_or(basis_type_),
                                      build_type.value_or(BuildType::PARALLEL), max_level, max_nodes_in_grid,
                                      cancellation_flag);
        }

        void to_pb_2D(grid::Grid2D *grid)
        {
            grid->set_basis_type(
                basis_type_ == BasisType::LINEAR ? grid::BasisType::LINEAR : grid::BasisType::QUADRATIC);
            auto min = grid->mutable_min();
            min->set_x(min_[0]);
            min->set_y(min_[1]);
            auto max = grid->mutable_max();
            max->set_x(max_[0]);
            max->set_y(max_[1]);

            for (auto &[_, node]: nodes_)
            {
                auto pb_node = grid->add_nodes();
                node_to_pb(node, pb_node);
            }

            for (auto &entry: entry_points_)
            {
                auto pb_node = grid->add_entry_points();
                node_to_pb(entry.node, pb_node);
            }
        }

    private:
        static void node_to_pb(const Node &node, grid::Grid2D::Node2D *to)
        {
            to->set_has_children(node.has_children);
            auto alpha = to->mutable_alpha();
            alpha->set_x(node.alpha[0]);
            alpha->set_y(node.alpha[1]);

            auto center = to->mutable_center_unit();
            center->set_x(node.center_unit[0]);
            center->set_y(node.center_unit[1]);

            auto grid_key = to->mutable_key();
            auto level = grid_key->mutable_level();
            level->set_x(node.key.level[0]);
            level->set_y(node.key.level[1]);

            auto index = grid_key->mutable_index();
            index->set_x(node.key.index[0]);
            index->set_y(node.key.index[1]);
        }

        void check_constraints()
        {
            for (size_t i = 0; i < IN_DIM; ++i)
            {
                if (std::abs(max_[i] - min_[i]) < std::numeric_limits<ScalarType>::epsilon())
                {
                    throw std::runtime_error(
                        std::format("Invalid bounds at dimension {}, {} ~ {}", i, min_[i], max_[i]));
                } else if (max_[i] < min_[i])
                {
                    std::swap(min_[i], max_[i]);
                }
            }
        }

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

        template<function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        std::vector<std::pair<Argument, Answer> > calculate_original_function_at_points(
            const FUNC &func, const std::vector<Argument> &anchor_points)
        {
            auto r = anchor_points | std::ranges::views::transform([&, this](const Argument &arg)
            {
                return std::pair{to_unit(arg), func(arg)};
            });
            return {r.begin(), r.end()};
        }

        void check_evaluation_point(const Argument &real) const
        {
            for (size_t i = 0; i < IN_DIM; ++i)
            {
                if (real[i] < min_[i] - std::numeric_limits<ScalarType>::epsilon() || real[i] > max_[i] +
                    std::numeric_limits<ScalarType>::epsilon())
                    throw std::runtime_error(std::format("Evaluation point out of bounds at dimension {}", i));
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

        template<function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        void build(const FUNC &func, ScalarType epsilon, const std::vector<std::pair<Argument, Answer> > &anchors,
                   BuildType build_type, size_t max_level, size_t max_nodes_in_grid,
                   std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag)
        {
            const auto launch_type = build_type == BuildType::PARALLEL ? std::launch::async : std::launch::deferred;
            for (size_t i = 0; i <= IN_DIM; ++i)
            {
                GridKey key;
                key.level = 0;
                key.index = 0;

                for (size_t j = 0; j < i; ++j)
                {
                    key.level[IN_DIM - j - 1] = 1;
                }

                std::vector<std::future<std::pair<Node, NodeMap> > > futures;

                do
                {
                    size_t index_permutations = 1zu << (IN_DIM - i);

                    for (size_t j = 0; j < index_permutations; ++j)
                    {
                        expand_indicies(key, j);
                        NodeMap tmp;
                        Node node = tmp[create_node(func, key, {}, i, tmp).value()];
                        futures.emplace_back(std::async(launch_type, &AdaptiveSparseGrid::build_grid<FUNC>,
                                                        this, std::ref(func), epsilon, std::ref(anchors), node, i,
                                                        max_level, max_nodes_in_grid,
                                                        cancellation_flag));
                    }
                } while (std::ranges::next_permutation(key.level.coords).found);

                for (auto &future: futures)
                {
                    future.wait();
                }

                for (auto &future: futures)
                {
                    auto [node, new_nodes] = future.get();
                    EntryPoint entry_point;
                    entry_point.dimensions = i;
                    entry_point.node = node;
                    entry_points_.push_back(entry_point);
                    nodes_.insert_range(new_nodes);
                }
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

        template<function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        std::pair<Node, NodeMap> build_grid(const FUNC &func, ScalarType epsilon,
                                            const std::vector<std::pair<Argument, Answer> > &anchors, Node entry_point,
                                            size_t dimension, size_t max_level, size_t max_nodes_in_grid,
                                            std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag)
        {
            std::queue<GridKey> node_queue;
            node_queue.push(entry_point.key);
            NodeMap new_nodes;
            auto &ref_node = new_nodes.emplace(entry_point.key, entry_point).first->second;
            size_t current_max_level = entry_point.key.level.max();

            while (!node_queue.empty())
            {
                if (cancellation_flag.has_value() && cancellation_flag.value().get().test())
                {
                    auto str = std::format("Interrupted at node count: {}, dimension: {}, max level: {}", new_nodes.size(), dimension, current_max_level);
                    std::cerr << str << std::endl;
                    throw std::runtime_error(str);
                }

                Node &node = new_nodes[node_queue.front()];
                node_queue.pop();

                bool can_continue = false;
                bool can_continue_force = node.key.level.max() >= sizeof(size_t) * 8;
                VectorValue<Direction, IN_DIM> directions(Direction::BOTH);

                if (can_continue_force || node.alpha.length() <= epsilon)
                {
                    directions.fill(Direction::NONE);
                    can_continue = true;
                }

                if (can_continue && !can_continue_force)
                {
                    for (auto &anchor: anchors)
                    {
                        if (node.is_point_in_affect_zone(anchor.first) && (
                                evaluate_for_dim_and_entry_point(anchor.first, dimension, entry_point, new_nodes) -
                                anchor.second).length() >= epsilon)
                        {
                            can_continue = false;
                            directions |= Basis::get_affect_direction(anchor.first, node.key.level, node.key.index);
                        }
                    }
                }

                if (can_continue || can_continue_force)
                {
                    continue;
                }

                node.has_children = true;

                for (size_t i = 0; i < IN_DIM; ++i)
                {
                    if (node.key.level[i] == 0 || (max_level > 0 && node.key.level[i] >= max_level))
                        continue;

                    GridKey key_left = node.key, key_right = node.key;

                    key_left.level[i] = key_right.level[i] = node.key.level[i] + 1;

                    if (key_left.level[i] > current_max_level)
                    {
                        current_max_level = key_left.level[i];
                    }

                    key_left.index[i] = 2 * node.key.index[i] - 1;
                    key_right.index[i] = 2 * node.key.index[i] + 1;

                    if ((Direction::LEFT & directions[i]) != Direction::NONE)
                    {
                        auto left_node = create_node(func, key_left, ref_node, dimension, new_nodes);
                        if (left_node.has_value())
                        {
                            node_queue.push(left_node.value());
                        }
                    }

                    if ((Direction::RIGHT & directions[i]) != Direction::NONE)
                    {
                        auto right_node = create_node(func, key_right, ref_node, dimension, new_nodes);

                        if (right_node.has_value())
                        {
                            node_queue.push(right_node.value());
                        }
                    }
                }

                if (max_nodes_in_grid > 0 && new_nodes.size() >= max_nodes_in_grid)
                {
                    break;
                }
            }

            return {ref_node, new_nodes};
        }

        Answer evaluate_for_dim_and_entry_point(const Argument &x, size_t max_grid_dim, std::optional<Node> entry_point,
                                                const NodeMap &additional_nodes = {})
        {
            Answer interp = evaluate_for_dim(x, max_grid_dim);
            if (entry_point.has_value())
            {
                NodeSet processed_nodes;
                interp += evaluate_recursive(x, entry_point.value(), processed_nodes, additional_nodes);
            }
            return interp;
        }

        template<function_traits::ArrayOpInOut<ScalarType, IN_DIM, OUT_DIM> FUNC>
        std::optional<GridKey> create_node(const FUNC &func, const GridKey &key, std::optional<Node> entry_point,
                                           size_t dimension, NodeMap &additional_nodes)
        {
            if (additional_nodes.contains(key))
                return {};

            Node node;
            node.key = key;
            for (size_t i = 0; i < IN_DIM; ++i)
                node.center_unit[i] = Basis::get_coord(key.level[i], key.index[i]);

            Argument real = to_real(node.center_unit);
            Answer etalon = func(real.coords);
            Answer interp =
                    evaluate_for_dim_and_entry_point(node.center_unit, dimension, entry_point, additional_nodes);

            node.alpha = etalon - interp;

            additional_nodes[key] = node;
            return key;
        }

        Answer evaluate_for_dim(const Argument &x, size_t max_grid_dim = IN_DIM + 1) const
        {
            Answer answer = 0;
            for (auto &&entry_point: entry_points_)
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

        Answer evaluate_recursive(const Argument &x, const Node &node, NodeSet &processed_nodes,
                                  const NodeMap &additional_nodes = {}) const
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
                        auto child = get_child_for_dim_and_arg(node, x, i, additional_nodes);
                        if (child.has_value())
                        {
                            value += evaluate_recursive(x, child.value(), processed_nodes, additional_nodes);
                        }
                    }
                }
            }
            return value;
        }

        std::optional<Node> get_child_for_dim_and_arg(const Node &parent, const Argument &x, size_t dimension,
                                                      const NodeMap &additional_nodes = {}) const
        {
            GridKey key = parent.key;
            key.level[dimension] += 1;
            key.index[dimension] = 2 * key.index[dimension] + (x[dimension] < parent.center_unit[dimension] ? -1 : 1);

            auto child = nodes_.find(key);
            if (child != nodes_.end())
                return child->second;
            child = additional_nodes.find(key);
            if (child != additional_nodes.end())
                return child->second;
            else
                return {};
        }
    };

    template<size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, Point<IN_DIM> min, Point<IN_DIM> max, ScalarType epsilon,
                       const std::vector<Point<IN_DIM> > &anchor_points = {}, BasisType type = BasisType::QUADRATIC,
                       BuildType build_type = BuildType::PARALLEL, size_t max_level = 0, size_t max_nodes_in_dim = 0,
                       std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag = {}) ->
        AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template<size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, std::array<ScalarType, IN_DIM> min, std::array<ScalarType, IN_DIM> max,
                       ScalarType epsilon, const std::vector<Point<IN_DIM> > &anchor_points = {},
                       BasisType type = BasisType::QUADRATIC, BuildType build_type = BuildType::PARALLEL,
                       size_t max_level = 0, size_t max_nodes_in_dim = 0,
                       std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag = {}) ->
        AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template<size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, Point<IN_DIM> min, std::array<ScalarType, IN_DIM> max, ScalarType epsilon,
                       const std::vector<Point<IN_DIM> > &anchor_points = {}, BasisType type = BasisType::QUADRATIC,
                       BuildType build_type = BuildType::PARALLEL, size_t max_level = 0, size_t max_nodes_in_dim = 0,
                       std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag = {}) ->
        AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;

    template<size_t IN_DIM, typename FUNC>
    AdaptiveSparseGrid(FUNC func, std::array<ScalarType, IN_DIM> min, Point<IN_DIM> max, ScalarType epsilon,
                       const std::vector<Point<IN_DIM> > &anchor_points = {}, BasisType type = BasisType::QUADRATIC,
                       BuildType build_type = BuildType::PARALLEL, size_t max_level = 0, size_t max_nodes_in_dim = 0,
                       std::optional<std::reference_wrapper<std::atomic_flag> > cancellation_flag = {}) ->
        AdaptiveSparseGrid<IN_DIM, function_traits::Analyzer<FUNC>::OUT>;
}
