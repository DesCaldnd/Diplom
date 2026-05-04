module;

#include <grpcpp/grpcpp.h>
#include <proto/grid.pb.h>
#include <proto/grid.grpc.pb.h>
#include <exprtk.hpp>

export module server;

import util;
import compute;

class Formula final
{
public:
    Formula(const std::string& formula_x, const std::string& formula_y, Compute::ScalarType dt, bool is_diffur) : dt_(dt), is_diffur_(is_diffur)
    {
        symbol_table_.add_variable("x", x_);
        symbol_table_.add_variable("y", y_);
        if (is_diffur_)
        {
            symbol_table_.add_variable("t", t_);
        }
        symbol_table_.add_constants();

        expression_x_.register_symbol_table(symbol_table_);

        exprtk::parser<Compute::ScalarType> parser;
        if (!parser.compile(formula_x, expression_x_))
        {
            throw std::runtime_error("Formula failed: x");
        }

        expression_y_.register_symbol_table(symbol_table_);

        if (!parser.compile(formula_y, expression_y_))
        {
            throw std::runtime_error("Formula failed: y");
        }
    }

    void set_diapason(Compute::ScalarType start_t, Compute::ScalarType end_t)
    {
        t_begin_ = start_t;
        t_end_ = end_t;
    }

    Compute::Point<2> operator()(Compute::Point<2> arg) const
    {
        std::lock_guard lock(mutex_);
        Compute::ScalarType x = arg[0], y = arg[1];

        if (is_diffur_)
        {
            for (t_ = t_begin_; t_ <= t_end_; t_ += dt_)
            {
                auto get_derivatives = [&] (Compute::ScalarType cur_x, Compute::ScalarType cur_y)
                {
                    x_ = cur_x;
                    y_ = cur_y;
                    return std::make_pair(expression_x_.value(), expression_y_.value());
                };

                auto [dx1, dy1] = get_derivatives(x, y);
                auto [dx2, dy2] = get_derivatives(x + dx1 * dt_ / 2, y + dy1 * dt_ / 2);
                auto [dx3, dy3] = get_derivatives(x + dx2 * dt_ / 2, y + dy2 * dt_ / 2);
                auto [dx4, dy4] = get_derivatives(x + dx3 * dt_, y + dy3);

                x += (dt_ / 6) * (dx1 + dx2 * 2 + dx3 * 3 + dx4);
                y += (dt_ / 6) * (dy1 + dy2 * 2 + dy3 * 3 + dy4);
            }
        } else
        {
            x_ = x;
            y_ = y;
            x = expression_x_.value();
            y = expression_y_.value();
        }

        return {x, y};
    }

private:
    mutable Compute::ScalarType x_;
    mutable Compute::ScalarType y_;
    mutable Compute::ScalarType t_;
    mutable std::mutex mutex_;
    bool is_diffur_;
    Compute::ScalarType t_begin_ = 0;
    Compute::ScalarType t_end_ = 1;
    Compute::ScalarType dt_;
    exprtk::symbol_table<Compute::ScalarType> symbol_table_;
    exprtk::expression<Compute::ScalarType> expression_x_;
    exprtk::expression<Compute::ScalarType> expression_y_;

    static constexpr Compute::ScalarType T_INTERVAL = 1;
};

export class GridServer : public grid::GridService::Service
{
    grpc::Status GetGrid2D(grpc::ServerContext* context, const grid::Grid2DRequest* request, grid::Grid2D* response) override
    {
        std::atomic_flag cancellation_flag;
        cancellation_flag.clear();

        auto cancellation_future = std::async(std::launch::async, [&]()
        {
            while (true)
            {
                if (context->IsCancelled() || cancellation_flag.test())
                {
                    cancellation_flag.test_and_set();
                    break;
                }
                std::this_thread::yield();
            }
        });

        try
        {
            Compute::Point<2> min{request->min().x(), request->min().y()};
            Compute::Point<2> max{request->max().x(), request->max().y()};
            Formula formula(request->formula_x(), request->formula_y(), request->eps(), request->formula_type() == grid::FormulaType::DIFFUR);
            formula.set_diapason(0, std::min(1.0, request->step()));
            Compute::BasisType basis_type = request->basis_type() == grid::BasisType::LINEAR ? Compute::BasisType::LINEAR : Compute::BasisType::QUADRATIC;
            Compute::BuildType build_type = request->build_type() == grid::BuildType::PARALLEL ? Compute::BuildType::PARALLEL : Compute::BuildType::SEQUENTIAL;
            std::vector<Compute::Point<2>> anchor_points;
            anchor_points.reserve(request->anchor_points_size());
            for (auto& anchor : request->anchor_points())
            {
                anchor_points.push_back({anchor.x(), anchor.y()});
            }

            Compute::AdaptiveSparseGrid result(formula, min, max, request->eps(), anchor_points, basis_type, build_type, request->max_level(), request->max_nodes_in_dim(), std::ref(cancellation_flag));

            for (Compute::ScalarType current_step = 1; current_step < request->step(); ++current_step)
            {
                formula.set_diapason(current_step, std::min(current_step + 1, request->step()));
                result = result.make_next_iteration(formula, request->eps(), anchor_points, basis_type, build_type, request->max_level(), request->max_nodes_in_dim(), std::ref(cancellation_flag));
            }

            result.to_pb_2D(response);
            cancellation_flag.test_and_set();
            cancellation_future.wait();

            std::cout << "Message size in mb: " << response->ByteSizeLong() / 1024.0 / 1024.0 << std::endl;
            return grpc::Status::OK;
        } catch (const std::runtime_error& e)
        {
            cancellation_flag.test_and_set();
            cancellation_future.wait();
            return {grpc::StatusCode::INVALID_ARGUMENT, e.what()};
        }

    }
};