module;

export module vector_value;

import util;

export namespace Compute
{
    template <typename T, size_t DIM>
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

        VectorValue(const std::initializer_list<T>& values)
        {
            auto iter = values.begin();
            for(size_t i = 0; i < values.size(); ++i, ++iter)
            {
                coords[i] = *iter;
            }

            for(size_t i = values.size(); i < DIM; ++i)
            {
                coords[i] = 0;
            }
        }

        VectorValue(const T &val)
        {
            coords.fill(val);
        }

        void fill(const T &val)
        {
            coords.fill(val);
        }

        T &operator[](size_t i) { return coords[i]; }
        const T &operator[](size_t i) const { return coords[i]; }
        const

            VectorValue &
            operator+=(const VectorValue &other)
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

        VectorValue &operator*=(const VectorValue &other)
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

        VectorValue &operator|=(const VectorValue &other)
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                coords[i] |= other.coords[i];
            }
            return *this;
        }

        VectorValue operator|(const VectorValue &other) const
        {
            VectorValue res(*this);
            res |= other;
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

        bool operator==(const VectorValue &other) const
        {
            for (size_t i = 0; i < DIM; ++i)
            {
                if (coords[i] != other.coords[i])
                    return false;
            }
            return true;
        }
    };

    using ScalarType = double;
    template <size_t DIM>
    using Point = VectorValue<ScalarType, DIM>;
}

export namespace std
{
    template <typename T, size_t DIM>
    struct hash<Compute::VectorValue<T, DIM>>
    {
        size_t operator()(const Compute::VectorValue<T, DIM> &k) const
        {
            size_t seed = 0;
            auto combine = [&](size_t v)
            { seed ^= v + 0x9e3779b9 + (seed << 6) + (seed >> 2); };
            for (size_t i = 0; i < DIM; ++i)
            {
                combine(k[i]);
            }
            return seed;
        }
    };
}