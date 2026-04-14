//
// Created by ivanl on 01.03.2026.
//
import compute;
import util;

Compute::Point<1> func(Compute::Point<2> arg)
{
    return {std::sin(arg[0] * 5) + 7 * std::cos(arg[1] * 6)};
}

Compute::Point<1> func2(Compute::Point<1> arg)
{
    return {arg[0] * 2};
}

Compute::Point<1> func3(Compute::Point<1> arg)
{
    if (!(arg[0] > -2 && arg[0] < -1))
        return 0;
    else
        return std::abs(1.5 + arg[0]) - 0.5;
}

void print(double a, double b)
{
    std::cout << a << " " << b << std::endl;
}

int main() {
    // Compute::Point<2> min = {0, 0}, max = {1, 1};
    // Compute::AdaptiveSparseGrid grid(&func, min, max, 0.0000001);
    // Compute::AdaptiveSparseGrid grid2 = grid.make_next_iteration(&func2, 0.0000001);

    // Compute::Point<2> arg;
    // std::cout << "Eval:" << std::endl;

    // arg = {0, 0};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);

    // arg = {0, 0.6};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);

    // arg = {0.3, 0};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);

    // arg = {0.6, 0.3};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);

    // arg = {0.5, 0.5};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);

    // arg = {1, 1};
    // print(grid2.evaluate(arg)[0], func(arg)[0]);


    Compute::Point<1> min3 = {-10}, max3 = {10};
    Compute::AdaptiveSparseGrid grid3(&func3, min3, max3, 0.0000001, {{-1.5}});
    Compute::Point<1> arg2;
    std::cout << std::endl << std::endl;

    arg2 = {0};
    print(grid3.evaluate(arg2)[0], func3(arg2)[0]);
    arg2 = {-10};
    print(grid3.evaluate(arg2)[0], func3(arg2)[0]);
    arg2 = {-1.5};
    print(grid3.evaluate(arg2)[0], func3(arg2)[0]);
}