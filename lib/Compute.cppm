//
// Created by ivanl on 01.03.2026.
//
module;

export module compute;

import util;
import function_traits;

namespace Compute
{
    //public prerequisites names
    export
    {
        enum class BasisType
        {
            LINEAR,
            QUADRATIC
        };
    }

    template<typename T, size_t DIM>
    struct VectorValue
    {
        std::array<T, DIM> coords;

        VectorValue() = default;

        VectorValue(const std::array<T, DIM> &coords) : coords(coords)
        {
        }

        VectorValue(std::array<T, DIM> &&coords) : coords(std::move(coords))
        {
        }

        VectorValue(const T &val)
        {
            coords.fill(val);
        }

        T &operator[](size_t i) { return coords[i]; }
        const T &operator[](size_t i) const { return coords[i]; }

        VectorValue &operator+=(const VectorValue &other)
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                coords[i] += other.coords[i];
            }
            return *this;
        }

        VectorValue operator+(const VectorValue &other) const
        {
            VectorValue res(*this);
            res += other;
            return res;
        }

        VectorValue &operator-=(const VectorValue &other)
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                coords[i] -= other.coords[i];
            }
            return *this;
        }

        VectorValue operator-(const VectorValue &other) const
        {
            VectorValue res(*this);
            res -= other;
            return res;
        }

        VectorValue &operator*=(const VectorValue &other) const
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                coords[i] *= other.coords[i];
            }
            return *this;
        }

        VectorValue operator*(const VectorValue &other) const
        {
            VectorValue res(*this);
            res *= other;
            return res;
        }

        VectorValue &operator*=(T other)
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                coords[i] *= other;
            }
            return *this;
        }

        VectorValue operator*(T other) const
        {
            VectorValue res(*this);
            res *= other;
            return res;
        }

        T length() const
        {
            return std::sqrt(std::inner_product(coords.begin(), coords.end(), coords.begin(), 0.0));
        }

        VectorValue<T, DIM - 1> without(size_t idx) const
        {
            if (idx >= DIM) throw std::out_of_range("Index out of range in method without");
            VectorValue<T, DIM - 1> result;
            for (size_t i = 0, current_idx = 0; i < DIM; ++i)
            {
                if (i == idx)
                    continue;

                result[current_idx++] = coords[i];
            }

            return result;
        }

        bool operator==(const VectorValue &other) const
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                if (coords[i] != other.coords[i]) return false;
            }
            return true;
        }
    };

    template<size_t DIM>
    using Index = VectorValue<size_t, DIM>;
    using ScalarType = double;
    template<size_t DIM>
    using Point = VectorValue<ScalarType, DIM>;

    template<size_t DIM>
    struct GridKey
    {
        Index<DIM> level;
        Index<DIM> index;

        bool operator==(const GridKey &other) const = default;
    };

    template<size_t DIM>
    struct GridKeyHash
    {
        size_t operator()(const GridKey<DIM> &k) const
        {
            size_t seed = 0;
            auto combine = [&](size_t v) { seed ^= v + 0x9e3779b9 + (seed << 6) + (seed >> 2); };
            for (size_t i = 0; i < DIM; ++i)
            {
                combine(k.level[i]);
                combine(k.index[i]);
            }
            return seed;
        }
    };

    template<size_t DIM>
    class Basis
    {
    private:
        // Единая функция оценки 1D базиса
        static ScalarType eval_1d(ScalarType x, size_t level, size_t index, BasisType type)
        {
            if (level == 0)
            {
                if (index == 0) return 1.0 - x; // Убывает от 0
                if (index == 1) return x; // Растет от 0
                return 0.0;
            }

            // --- ВНУТРЕННОСТИ (Level >= 1) ---
            // Локальные параболы / шляпы
            ScalarType h = 1.0 / (1 << level);
            ScalarType center = static_cast<ScalarType>(index) * h;
            ScalarType dist = std::abs(x - center);

            if (dist >= h) return 0.0;

            // Квадратичный базис (из диплома)
            ScalarType t = dist / h;

            ScalarType result = 0;

            switch (type)
            {
                case BasisType::LINEAR:
                    result = 1 - t;
                    break;
                case BasisType::QUADRATIC:
                    result = 1 - t * t;
                    break;
            }

            return result;
        }

    public:
        static ScalarType eval(const Point<DIM> &x, const Index<DIM> &level, const Index<DIM> &index, BasisType type)
        {
            ScalarType result = 1;
            for (size_t i = 0; i < DIM; ++i)
            {
                result *= eval_1d(x[i], level[i], index[i], type);
                if (result == 0)
                    return result;
            }

            return result;
        }

        static ScalarType get_coord(size_t level, size_t index)
        {
            if (level == 0) return (index == 0) ? 0.0 : 1.0;
            return static_cast<ScalarType>(index) / (1 << level);
        }
    };

    template<size_t IN_DIM, size_t OUT_DIM>
    struct Node
    {
        GridKey<IN_DIM> key;
        Point<IN_DIM> center_unit;
        Point<OUT_DIM> alpha;
        bool has_children = false;
    };

    template<size_t IN_DIM, size_t OUT_DIM>
    struct EntryPoint
    {
        Node<IN_DIM, OUT_DIM> node;
        // нужно, чтобы считать границы, как сетки меньшей размерности
        size_t dimensions;
    };

    export
    {
        template<size_t IN_DIM, size_t OUT_DIM>
        struct EvaluationNodeInfo
        {
            Node<IN_DIM, OUT_DIM> node;
            ScalarType basis;
            Point<OUT_DIM> diff;
        };

        template<function_traits::ArrayOp<ScalarType> FUNC>
        class AdaptiveSparseGrid
        {
            using FuncInfo = function_traits::Analyzer<FUNC>;
        public:
            static constexpr size_t IN_DIM = FuncInfo::IN;
            static constexpr size_t OUT_DIM = FuncInfo::OUT;

            using EntryPoint = EntryPoint<IN_DIM, OUT_DIM>;
            using GridKey = GridKey<IN_DIM>;
            using Node = Node<IN_DIM, OUT_DIM>;
            using GridKeyHash = GridKeyHash<IN_DIM>;
            using Argument = Point<IN_DIM>;
            using Answer = Point<OUT_DIM>;
            using Point = Point<IN_DIM>;
            using Basis = Basis<IN_DIM>;
            using NodeMap = std::unordered_map<GridKey, Node, GridKeyHash>;
            using NodeSet = std::unordered_set<GridKey, GridKeyHash>;
            using EvaluationNodeInfo = EvaluationNodeInfo<IN_DIM, OUT_DIM>;

        private:
            Point min_;
            Point max_;
            ScalarType epsilon_;
            FUNC func_;
            BasisType type_;

            std::vector<EntryPoint> entry_points_;
            NodeMap nodes_;

        public:
            AdaptiveSparseGrid(FUNC func, Point min, Point max, ScalarType epsilon,
                               BasisType type = BasisType::QUADRATIC) : min_(min), max_(max), epsilon_(epsilon),
                                                                        func_(func), type_(type)
            {
                entry_points_.reserve(pow_3(IN_DIM));
                build();
            }

            Answer evaluate(const Argument &x, std::optional<std::reference_wrapper<std::vector<EvaluationNodeInfo>>> meta = {}) const
            {
                check_evaluation_point(x);
                return evaluate_for_dim(to_unit(x), IN_DIM + 1, meta);
            }

            std::vector<EvaluationNodeInfo> dump_nodes()
            {
                std::vector<EvaluationNodeInfo> answer;
                for (auto&& pair : nodes_)
                {
                    EvaluationNodeInfo meta;
                    meta.node = pair.second;
                    meta.basis = 0;
                    meta.diff = 0;
                    answer.push_back(meta);
                }
                return answer;
            }

            std::vector<EvaluationNodeInfo> dump_entry_points()
            {
                std::vector<EvaluationNodeInfo> answer;
                for (auto&& entry : entry_points_)
                {
                    EvaluationNodeInfo meta;
                    meta.node = entry.node;
                    meta.basis = entry.dimensions;
                    meta.diff = 0;
                    answer.push_back(meta);
                }
                return answer;
            }

        private:
            Point to_real(const Point &unit) const
            {
                Point real;
                for (size_t i = 0; i < IN_DIM; ++i) real[i] = min_[i] + unit[i] * (max_[i] - min_[i]);
                return real;
            }

            Point to_unit(const Point &real) const
            {
                Point unit;
                for (size_t i = 0; i < IN_DIM; ++i) unit[i] = (real[i] - min_[i]) / (max_[i] - min_[i]);
                return unit;
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

            void build()
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
                            Node& node = nodes_[create_node(key, {}, i).value()];
                            build_grid(node, i);
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

            void build_grid(const Node &entry_point, size_t dimension)
            {
                std::queue<GridKey> node_queue;
                node_queue.push(entry_point.key);

                while (!node_queue.empty())
                {
                    Node& node = nodes_[node_queue.front()];
                    node_queue.pop();

                    if (node.alpha.length() <= epsilon_)
                        continue;

                    node.has_children = true;

                    for (size_t i = 0; i < IN_DIM; ++i)
                    {
                        if (node.key.level[i] == 0)
                            continue;

                        GridKey key_left = node.key, key_right = node.key;

                        key_left.level[i] = key_right.level[i] = node.key.level[i] + 1;

                        key_left.index[i] = 2 * node.key.index[i] - 1;
                        key_right.index[i] = 2 * node.key.index[i] + 1;

                        auto left_node = create_node(key_left, entry_point, dimension);
                        auto right_node = create_node(key_right, entry_point, dimension);

                        if (left_node.has_value())
                            node_queue.push(left_node.value());
                        if (right_node.has_value())
                            node_queue.push(right_node.value());
                    }
                }
            }

            std::optional<GridKey> create_node(const GridKey &key, std::optional<Node> entry_point, size_t dimension)
            {
                if (nodes_.contains(key))
                    return {};

                Node node;
                node.key = key;
                for (size_t i = 0; i < IN_DIM; ++i)
                    node.center_unit[i] = Basis::get_coord(key.level[i], key.index[i]);

                Argument real = to_real(node.center_unit);
                Answer etalon = func_(real.coords);
                Answer interp = evaluate_for_dim(node.center_unit, dimension);
                if (entry_point.has_value())
                {
                    NodeSet processed_nodes;
                    interp += evaluate_recursive(node.center_unit, entry_point.value(), processed_nodes);
                }

                node.alpha = etalon - interp;

                nodes_[key] = node;
                return key;
            }

            Answer evaluate_for_dim(const Argument &x, size_t max_grid_dim = IN_DIM + 1, std::optional<std::reference_wrapper<std::vector<EvaluationNodeInfo>>> meta = {}) const
            {
                Answer answer = 0;
                for (auto&& entry_point : entry_points_)
                {
                    if (entry_point.dimensions >= max_grid_dim)
                    {
                        break;
                    }
                    NodeSet processed_nodes;
                    answer += evaluate_recursive(x, entry_point.node, processed_nodes, meta);
                }
                return answer;
            }

            Answer evaluate_recursive(const Argument &x, const Node &node, NodeSet& processed_nodes, std::optional<std::reference_wrapper<std::vector<EvaluationNodeInfo>>> meta = {}) const
            {
                if (processed_nodes.contains(node.key))
                    return 0;

                processed_nodes.emplace(node.key);
                auto basis = Basis::eval(x, node.key.level, node.key.index, type_);

                if (std::abs(basis) < std::numeric_limits<ScalarType>::epsilon())
                {
                    return 0;
                }

                auto value = node.alpha * basis;

                if (meta.has_value())
                {
                    EvaluationNodeInfo meta_node;
                    meta_node.node = node;
                    meta_node.basis = basis;
                    meta_node.diff = value;
                    meta->get().push_back(meta_node);
                }

                if (node.has_children)
                {
                    for (size_t i = 0; i < IN_DIM; ++i)
                    {
                        if (node.key.level[i] != 0)
                        {
                            auto child = get_child_for_dim_and_arg(node, x, i);
                            if (child.has_value())
                            {
                                value += evaluate_recursive(x, child.value(), processed_nodes, meta);
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
                key.index[dimension] = 2 * key.index[dimension] + (
                                           x[dimension] < parent.center_unit[dimension] ? -1 : 1);

                auto child = nodes_.find(key);
                if (child != nodes_.end())
                    return child->second;
                else
                    return {};
            }
        };
    }
}
