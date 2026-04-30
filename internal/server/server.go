package server

import (
	"context"
	"fmt"
	"math"
	"time"

	"Diplom/compute"
	pb "Diplom/gridProto"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grid_requests_total",
			Help: "Total number of requests to GetGrid2D",
		},
		[]string{"status"},
	)
	requestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "grid_request_duration_seconds",
			Help:    "Time spent processing GetGrid2D requests",
			Buckets: prometheus.DefBuckets,
		},
	)
)

type Formula struct {
	dt       compute.ScalarType
	isDiffur bool
	progX    *vm.Program
	progY    *vm.Program
	tBegin   compute.ScalarType
	tEnd     compute.ScalarType
}

func NewFormula(formulaX, formulaY string, dt compute.ScalarType, isDiffur bool) (*Formula, error) {
	if dt < 0.001 {
		dt = 0.001
	}
	if dt > 0.1 {
		dt = 0.1
	}

	env := map[string]interface{}{
		"x":     0.0,
		"y":     0.0,
		"t":     0.0,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"asin":  math.Asin,
		"acos":  math.Acos,
		"atan":  math.Atan,
		"sinh":  math.Sinh,
		"cosh":  math.Cosh,
		"tanh":  math.Tanh,
		"exp":   math.Exp,
		"log":   math.Log,
		"log10": math.Log10,
		"sqrt":  math.Sqrt,
		"abs":   math.Abs,
		"pow":   math.Pow,
	}

	progX, err := expr.Compile(formulaX, expr.Env(env))
	if err != nil {
		return nil, err
	}
	progY, err := expr.Compile(formulaY, expr.Env(env))
	if err != nil {
		return nil, err
	}

	return &Formula{
		dt:       dt,
		isDiffur: isDiffur,
		progX:    progX,
		progY:    progY,
	}, nil
}

func (f *Formula) SetDiapason(startT, endT compute.ScalarType) {
	f.tBegin = startT
	f.tEnd = endT
}

func (f *Formula) Evaluate(arg compute.Point) (compute.Point, error) {
	zeroPoint := compute.Point{0, 0}
	if len(arg) < 2 {
		return zeroPoint, fmt.Errorf("expected at least 2 arguments, got %d", len(arg))
	}

	x := arg[0]
	y := arg[1]

	env := map[string]interface{}{
		"x":     x,
		"y":     y,
		"t":     0.0,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"asin":  math.Asin,
		"acos":  math.Acos,
		"atan":  math.Atan,
		"sinh":  math.Sinh,
		"cosh":  math.Cosh,
		"tanh":  math.Tanh,
		"exp":   math.Exp,
		"log":   math.Log,
		"log10": math.Log10,
		"sqrt":  math.Sqrt,
		"abs":   math.Abs,
		"pow":   math.Pow,
	}

	toFloat64 := func(value interface{}) (float64, error) {
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		default:
			return 0, fmt.Errorf("unexpected evaluation result type %T", value)
		}
	}

	if f.isDiffur {
		getDerivatives := func(curX, curY float64, t float64) (float64, float64, error) {
			env["x"] = curX
			env["y"] = curY
			env["t"] = t

			resX, err := vm.Run(f.progX, env)
			if err != nil {
				return 0, 0, err
			}
			resY, err := vm.Run(f.progY, env)
			if err != nil {
				return 0, 0, err
			}

			dx, err := toFloat64(resX)
			if err != nil {
				return 0, 0, err
			}
			dy, err := toFloat64(resY)
			if err != nil {
				return 0, 0, err
			}

			return dx, dy, nil
		}

		// Адаптивный шаг Рунге-Кутты-Фельберга (RK45)
		t := f.tBegin
		dt := f.dt
		tol := 1e-6 // Допустимая погрешность

		for t < f.tEnd {
			if t+dt > f.tEnd {
				dt = f.tEnd - t
			}

			// Коэффициенты для RK45 (Cash-Karp)
			k1x, k1y, err := getDerivatives(x, y, t)
			if err != nil {
				return zeroPoint, err
			}

			k2x, k2y, err := getDerivatives(x+dt*(1.0/5.0)*k1x, y+dt*(1.0/5.0)*k1y, t+dt*(1.0/5.0))
			if err != nil {
				return zeroPoint, err
			}

			k3x, k3y, err := getDerivatives(x+dt*(3.0/40.0)*k1x+dt*(9.0/40.0)*k2x, y+dt*(3.0/40.0)*k1y+dt*(9.0/40.0)*k2y, t+dt*(3.0/10.0))
			if err != nil {
				return zeroPoint, err
			}

			k4x, k4y, err := getDerivatives(x+dt*(3.0/10.0)*k1x-dt*(9.0/10.0)*k2x+dt*(6.0/5.0)*k3x, y+dt*(3.0/10.0)*k1y-dt*(9.0/10.0)*k2y+dt*(6.0/5.0)*k3y, t+dt*(3.0/5.0))
			if err != nil {
				return zeroPoint, err
			}

			k5x, k5y, err := getDerivatives(x-dt*(11.0/54.0)*k1x+dt*(5.0/2.0)*k2x-dt*(70.0/27.0)*k3x+dt*(35.0/27.0)*k4x, y-dt*(11.0/54.0)*k1y+dt*(5.0/2.0)*k2y-dt*(70.0/27.0)*k3y+dt*(35.0/27.0)*k4y, t+dt)
			if err != nil {
				return zeroPoint, err
			}

			k6x, k6y, err := getDerivatives(x+dt*(1631.0/55296.0)*k1x+dt*(175.0/512.0)*k2x+dt*(575.0/13824.0)*k3x+dt*(44275.0/110592.0)*k4x+dt*(253.0/4096.0)*k5x, y+dt*(1631.0/55296.0)*k1y+dt*(175.0/512.0)*k2y+dt*(575.0/13824.0)*k3y+dt*(44275.0/110592.0)*k4y+dt*(253.0/4096.0)*k5y, t+dt*(7.0/8.0))
			if err != nil {
				return zeroPoint, err
			}

			// Решение 5-го порядка
			x5 := x + dt*((37.0/378.0)*k1x+(250.0/621.0)*k3x+(125.0/594.0)*k4x+(512.0/1771.0)*k6x)
			y5 := y + dt*((37.0/378.0)*k1y+(250.0/621.0)*k3y+(125.0/594.0)*k4y+(512.0/1771.0)*k6y)

			// Решение 4-го порядка
			x4 := x + dt*((2825.0/27648.0)*k1x+(18575.0/48384.0)*k3x+(13525.0/55296.0)*k4x+(277.0/14336.0)*k5x+(1.0/4.0)*k6x)
			y4 := y + dt*((2825.0/27648.0)*k1y+(18575.0/48384.0)*k3y+(13525.0/55296.0)*k4y+(277.0/14336.0)*k5y+(1.0/4.0)*k6y)

			// Оценка ошибки
			errX := math.Abs(x5 - x4)
			errY := math.Abs(y5 - y4)
			maxErr := math.Max(errX, errY)

			if maxErr <= tol {
				// Шаг успешен
				t += dt
				x = x5
				y = y5
			}

			// Адаптация шага
			if maxErr == 0 {
				dt *= 2.0
			} else {
				dt *= 0.9 * math.Pow(tol/maxErr, 0.2)
			}

			// Ограничения на шаг
			if dt < 1e-5 {
				dt = 1e-5
			} else if dt > f.dt*10 {
				dt = f.dt * 10
			}
		}
	} else {
		resX, err := vm.Run(f.progX, env)
		if err != nil {
			return zeroPoint, err
		}
		resY, err := vm.Run(f.progY, env)
		if err != nil {
			return zeroPoint, err
		}
		x, err = toFloat64(resX)
		if err != nil {
			return zeroPoint, err
		}
		y, err = toFloat64(resY)
		if err != nil {
			return zeroPoint, err
		}
	}

	return compute.Point{x, y}, nil
}

type GridServer struct {
	pb.UnimplementedGridServiceServer
}

func (s *GridServer) GetGrid2D(ctx context.Context, req *pb.Grid2DRequest) (*pb.Grid2D, error) {
	start := time.Now()
	defer func() {
		requestDuration.Observe(time.Since(start).Seconds())
	}()

	min := compute.Point{req.Min.X, req.Min.Y}
	max := compute.Point{req.Max.X, req.Max.Y}

	formula, err := NewFormula(req.FormulaX, req.FormulaY, req.Eps, req.FormulaType == pb.FormulaType_DIFFUR)
	if err != nil {
		requestsTotal.WithLabelValues("error").Inc()
		return nil, err
	}

	select {
	case <-ctx.Done():
		requestsTotal.WithLabelValues("canceled").Inc()
		return nil, ctx.Err()
	default:
	}

	step := 1.0
	if req.Step != nil {
		step = *req.Step
	}

	formula.SetDiapason(0, math.Min(1.0, step))

	basisType := compute.BasisTypeQuadratic
	if req.BasisType == pb.BasisType_LINEAR {
		basisType = compute.BasisTypeLinear
	}

	buildType := compute.BuildTypeParallel
	if req.BuildType == pb.BuildType_SEQUENTIAL {
		buildType = compute.BuildTypeSequential
	}

	var anchorPoints []compute.Point
	for _, anchor := range req.AnchorPoints {
		anchorPoints = append(anchorPoints, compute.Point{anchor.X, anchor.Y})
	}

	grid, err := compute.NewAdaptiveSparseGridWithContext(
		ctx,
		formula.Evaluate,
		min, max,
		req.Eps,
		anchorPoints,
		basisType,
		buildType,
		int64(req.MaxLevel),
		req.MaxNodesInDim,
	)
	if err != nil {
		if ctx.Err() != nil {
			requestsTotal.WithLabelValues("canceled").Inc()
		} else {
			requestsTotal.WithLabelValues("error").Inc()
		}
		return nil, err
	}

	for currentStep := 1.0; currentStep < step; currentStep++ {
		formula.SetDiapason(currentStep, math.Min(currentStep+1.0, step))
		grid, err = grid.MakeNextIterationWithContext(
			ctx,
			formula.Evaluate,
			req.Eps,
			anchorPoints,
			basisType,
			buildType,
			int64(req.MaxLevel),
			req.MaxNodesInDim,
		)
		if err != nil {
			if ctx.Err() != nil {
				requestsTotal.WithLabelValues("canceled").Inc()
			} else {
				requestsTotal.WithLabelValues("error").Inc()
			}
			return nil, err
		}
	}

	response := toPb2D(grid)
	requestsTotal.WithLabelValues("success").Inc()
	return response, nil
}

func toPb2D(g *compute.AdaptiveSparseGrid) *pb.Grid2D {
	grid := &pb.Grid2D{}
	if g.BasisType() == compute.BasisTypeLinear {
		grid.BasisType = pb.BasisType_LINEAR
	} else {
		grid.BasisType = pb.BasisType_QUADRATIC
	}

	min := g.Min()
	max := g.Max()
	grid.Min = &pb.Point2D{X: min[0], Y: min[1]}
	grid.Max = &pb.Point2D{X: max[0], Y: max[1]}

	for _, n := range g.Nodes() {
		grid.Nodes = append(grid.Nodes, nodeToPb(n))
	}

	for _, ep := range g.EntryPoints() {
		grid.EntryPoints = append(grid.EntryPoints, nodeToPb(ep))
	}

	return grid
}

func nodeToPb(n compute.NodeInfo) *pb.Grid2D_Node2D {
	return &pb.Grid2D_Node2D{
		HasChildren: n.HasChildren,
		Alpha:       &pb.Point2D{X: n.Alpha[0], Y: n.Alpha[1]},
		CenterUnit:  &pb.Point2D{X: n.CenterUnit[0], Y: n.CenterUnit[1]},
		Key: &pb.Grid2D_GridKey2D{
			Level: &pb.Grid2D_Index2D{X: uint64(n.Level[0]), Y: uint64(n.Level[1])},
			Index: &pb.Grid2D_Index2D{X: uint64(n.Index[0]), Y: uint64(n.Index[1])},
		},
	}
}
