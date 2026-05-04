//
// Created by ivanl on 04.03.2026.
//

export module function_traits;

import util;
import vector_value;

export namespace function_traits
{
    template<typename T, typename = void>
    struct Analyzer
    {
    };

    // --- 1. Конечная цель (By Value) ---
    template<typename S, size_t In, size_t Out>
    struct Analyzer<Compute::VectorValue<S, Out>(Compute::VectorValue<S, In>)>
    {
        using Scalar = S;
        static constexpr size_t IN = In;
        static constexpr size_t OUT = Out;
    };

    // --- 2. Конечная цель (By Const Reference) ---
    template<typename S, size_t In, size_t Out>
    struct Analyzer<Compute::VectorValue<S, Out>(const Compute::VectorValue<S, In> &)>
            : Analyzer<Compute::VectorValue<S, Out>(Compute::VectorValue<S, In>)>
    {
    };

    // --- 3. Снятие "оберток" указателей на функции ---
    template<typename R, typename A>
    struct Analyzer<R(*)(A)> : Analyzer<R(A)>
    {
    };

    template<typename R, typename A>
    struct Analyzer<R(*)(A) noexcept> : Analyzer<R(A)>
    {
    };

    // --- 4. Снятие "оберток" методов класса ---
    template<typename C, typename R, typename A>
    struct Analyzer<R(C::*)(A)> : Analyzer<R(A)>
    {
    };

    template<typename C, typename R, typename A>
    struct Analyzer<R(C::*)(A) const> : Analyzer<R(A)>
    {
    };

    template<typename C, typename R, typename A>
    struct Analyzer<R(C::*)(A) noexcept> : Analyzer<R(A)>
    {
    };

    template<typename C, typename R, typename A>
    struct Analyzer<R(C::*)(A) const noexcept> : Analyzer<R(A)>
    {
    };

    // --- 5. Лямбды и Функторы (SFINAE специализация) ---
    // Если у типа T (например, лямбды) есть метод operator(),
    // эта специализация активируется и делегирует тип в пункт 4.
    template<typename T>
    struct Analyzer<T, std::void_t<decltype(&std::remove_cvref_t<T>::operator())> >
            : Analyzer<decltype(&std::remove_cvref_t<T>::operator())>
    {
    };


    // --- Концепты (Остаются такими же) ---
    template<typename Func, typename T, size_t IN_DIM>
    concept ArrayOpIn = requires { typename Analyzer<Func>::Scalar; } &&
                        std::same_as<T, typename Analyzer<Func>::Scalar> &&
                        Analyzer<Func>::IN == IN_DIM;

    template<typename Func, typename T, size_t IN_DIM, size_t OUT_DIM>
    concept ArrayOpInOut = ArrayOpIn<Func, T, IN_DIM> &&
                           Analyzer<Func>::OUT == OUT_DIM;
}
