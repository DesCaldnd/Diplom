//
// Created by ivanl on 01.03.2026.
//
import compute;
import util;

std::array<double, 1> func(std::array<double, 2> arg)
{
    return {std::sin(arg[0] * 5) + 7 * std::cos(arg[1] * 6)};
}

void print(double a, double b)
{
    std::cout << a << " " << b << std::endl;
}

template <typename T, size_t DIM>
void print_array(std::array<T, DIM> a)
{
    for (auto v : a)
    {
        std::cout << v << " ";
    }
}

template<size_t IN_DIM, size_t OUT_DIM>
void print_info(Compute::EvaluationNodeInfo<IN_DIM, OUT_DIM> meta)
{
    std::cout << "Level: ";
    print_array(meta.node.key.level.coords);
    std::cout << std::endl;
    std::cout << "Index: ";
    print_array(meta.node.key.index.coords);
    std::cout << std::endl;
    std::cout << "Alpha: ";
    print_array(meta.node.alpha.coords);
    std::cout << std::endl;
    std::cout << "Basis: " << meta.basis << std::endl;
    std::cout << "Diff:  ";
    print_array(meta.diff.coords);
    std::cout << std::endl << std::endl;
}

template<size_t IN_DIM, size_t OUT_DIM>
void print_meta(const std::vector<Compute::EvaluationNodeInfo<IN_DIM, OUT_DIM>>& meta)
{
    for (auto&& v : meta)
    {
        print_info(v);
    }
}

int main() {
    std::array<double, 2> min = {0, 0}, max = {1, 1};
    Compute::AdaptiveSparseGrid grid(&func, min, max, 0.0000001);

    std::array<double, 2> arg;
    // std::vector<Compute::EvaluationNodeInfo<2, 1>> node_meta = grid.dump_entry_points();
    //
    // print_meta(node_meta);

    // arg = {0, 0};
    // print(grid.evaluate(arg)[0], func(arg)[0]);

    std::vector<Compute::EvaluationNodeInfo<2, 1>> evaluate_meta;
    arg = {0, 0.6};
    std::cout << "Eval:" << std::endl;
    print(grid.evaluate(arg, evaluate_meta)[0], func(arg)[0]);
    // print_meta(evaluate_meta);

    arg = {0.3, 0};
    print(grid.evaluate(arg)[0], func(arg)[0]);

    arg = {0.6, 0.3};
    print(grid.evaluate(arg)[0], func(arg)[0]);

    arg = {0.5, 0.5};
    print(grid.evaluate(arg)[0], func(arg)[0]);

    arg = {1, 1};
    print(grid.evaluate(arg)[0], func(arg)[0]);
}