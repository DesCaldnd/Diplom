module;

export module helper_types;

import util;
import vector_value;

export namespace Compute
{
    enum class BasisType
    {
        LINEAR,
        QUADRATIC
    };

    template <size_t DIM>
    using Index = VectorValue<size_t, DIM>;

    template <size_t DIM>
    struct GridKey
    {
        Index<DIM> level;
        Index<DIM> index;

        bool operator==(const GridKey &other) const = default;
    };

    enum class Direction
    {
        NONE = 0,
        LEFT = 1 << 0,
        RIGHT = 1 << 1,
        BOTH = LEFT | RIGHT
    };

    constexpr Direction operator|(Direction a, Direction b)
    {
        return static_cast<Direction>(static_cast<int>(a) | static_cast<int>(b));
    }

    constexpr Direction operator|=(Direction &a, Direction b)
    {
        return a = a | b;
    }

    constexpr Direction operator&(Direction a, Direction b)
    {
        return static_cast<Direction>(static_cast<int>(a) & static_cast<int>(b));
    }

    constexpr Direction operator&=(Direction &a, Direction b)
    {
        return a = a & b;
    }


    template <size_t DIM>
    class Basis
    {
    private:
        // Единая функция оценки 1D базиса
        static ScalarType eval_1d(ScalarType x, size_t level, size_t index, BasisType type)
        {
            if (level == 0)
            {
                if (index == 0)
                    return 1.0 - x;
                if (index == 1)
                    return x;
                return 0.0;
            }

            ScalarType h = static_cast<ScalarType>(1.0) / (1uz << level);
            ScalarType center = static_cast<ScalarType>(index) * h;
            ScalarType dist = std::abs(x - center);

            if (dist >= h)
                return 0.0;

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
            if (level == 0)
                return (index == 0) ? 0.0 : 1.0;
            return static_cast<ScalarType>(index) / (1uz << level);
        }

        static std::pair<ScalarType, ScalarType> get_affect_bounds(size_t level, size_t index)
        {
            ScalarType unit = static_cast<ScalarType>(1) / (1uz << level);
            ScalarType coord = static_cast<ScalarType>(index) / (1uz << level);
            return {coord - unit, coord + unit};
        }

        static VectorValue<Direction, DIM> get_affect_direction(const Point<DIM> &x, const Index<DIM> &level, const Index<DIM> &index)
        {
            VectorValue<Direction, DIM> result = Direction::NONE;
            for (size_t i = 0; i < DIM; ++i)
            {
                ScalarType coord = level[i] > 0 ? static_cast<ScalarType>(index[i]) / (1uz << level[i]) : index[i];
                result[i] = x[i] < coord ? Direction::LEFT : Direction::RIGHT;
            }

            return result;
        }
    };

    template <size_t IN_DIM, size_t OUT_DIM>
    struct Node
    {
        GridKey<IN_DIM> key;
        Point<IN_DIM> center_unit;
        Point<OUT_DIM> alpha;
        bool has_children = false;

        bool is_point_in_affect_zone(const Point<IN_DIM>& point)
        {
            for (size_t i = 0; i < IN_DIM; ++i)
            {
                if (key.level[i] != 0)
                {
                    auto bounds = Basis<IN_DIM>::get_affect_bounds(key.level[i], key.index[i]);
                    if (point[i] <= bounds.first || point[i] >= bounds.second)
                    {
                        return false;
                    }
                } else
                {
                    if (std::abs(point[i] - (key.index[i] == 0 ? 0 : 1)) > std::numeric_limits<ScalarType>::epsilon())
                    {
                        return false;
                    }
                }
            }
            return true;
        }
    };

    template <size_t IN_DIM, size_t OUT_DIM>
    struct EntryPoint
    {
        Node<IN_DIM, OUT_DIM> node;
        // нужно, чтобы считать границы, как сетки меньшей размерности
        size_t dimensions;
    };
}

export namespace std
{
    template <size_t DIM>
    struct hash<Compute::GridKey<DIM>>
    {
        size_t operator()(const Compute::GridKey<DIM> &k) const
        {
            size_t seed = 0;
            auto combine = [&](size_t v)
            { seed ^= v + 0x9e3779b9 + (seed << 6) + (seed >> 2); };
            for (size_t i = 0; i < DIM; ++i)
            {
                combine(k.level[i]);
                combine(k.index[i]);
            }
            return seed;
        }
    };
}