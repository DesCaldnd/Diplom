//
// Created by ivanl on 04.03.2026.
//

export module function_traits;

import util;

export namespace function_traits{
    template<typename T> struct Analyzer : Analyzer<decltype(&std::remove_cvref_t<T>::operator())> {};

    // 1. Конечная цель (By Value): сигнатура R(Arg), где Arg и R — std::array
    template<typename S, size_t In, size_t Out>
    struct Analyzer<std::array<S, Out>(std::array<S, In>)> {
        using Scalar = S;
        static constexpr size_t IN = In;
        static constexpr size_t OUT = Out;
    };

    // 2. Конечная цель (By Const Reference): поддержка (const std::array<>&)
    template<typename S, size_t In, size_t Out>
    struct Analyzer<std::array<S, Out>(const std::array<S, In>&)>
        : Analyzer<std::array<S, Out>(std::array<S, In>)> {};

    // 3. Снятие "оберток" указателей на функции и noexcept (C++17+)
    template<typename R, typename A> struct Analyzer<R(*)(A)>          : Analyzer<R(A)> {};
    template<typename R, typename A> struct Analyzer<R(*)(A) noexcept> : Analyzer<R(A)> {};

    // 4. Снятие "оберток" методов класса (для Лямбд и Функторов)
    // Покрываем 4 комбинации: const/mutable и noexcept/обычный
    template<typename C, typename R, typename A> struct Analyzer<R(C::*)(A)>                : Analyzer<R(A)> {};
    template<typename C, typename R, typename A> struct Analyzer<R(C::*)(A) const>          : Analyzer<R(A)> {};
    template<typename C, typename R, typename A> struct Analyzer<R(C::*)(A) noexcept>       : Analyzer<R(A)> {};
    template<typename C, typename R, typename A> struct Analyzer<R(C::*)(A) const noexcept> : Analyzer<R(A)> {};


    // --- 2. Концепт (Валидатор интерфейса) ---
    // Проверяет, что Analyzer смог "докопаться" до типа Scalar
    template<typename Func, typename T>
    concept ArrayOp = requires { typename Analyzer<Func>::Scalar; } && std::same_as<T, typename Analyzer<Func>::Scalar>;
}