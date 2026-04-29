package server

import (
	"context"
	"fmt"
	"math"

	"Diplom/compute"
	pb "Diplom/gridProto"

	"github.com/Knetic/govaluate"
)

type Formula struct {
	dt       compute.ScalarType
	isDiffur bool
	exprX    *govaluate.EvaluableExpression
	exprY    *govaluate.EvaluableExpression
	tBegin   compute.ScalarType
	tEnd     compute.ScalarType
}

func NewFormula(formulaX, formulaY string, dt compute.ScalarType, isDiffur bool) (*Formula, error) {
	functions := map[string]govaluate.ExpressionFunction{
		"sin": func(args ...interface{}) (interface{}, error) {
			return math.Sin(args[0].(float64)), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			return math.Cos(args[0].(float64)), nil
		},
		"tan": func(args ...interface{}) (interface{}, error) {
			return math.Tan(args[0].(float64)), nil
		},
		"asin": func(args ...interface{}) (interface{}, error) {
			return math.Asin(args[0].(float64)), nil
		},
		"acos": func(args ...interface{}) (interface{}, error) {
			return math.Acos(args[0].(float64)), nil
		},
		"atan": func(args ...interface{}) (interface{}, error) {
			return math.Atan(args[0].(float64)), nil
		},
		"sinh": func(args ...interface{}) (interface{}, error) {
			return math.Sinh(args[0].(float64)), nil
		},
		"cosh": func(args ...interface{}) (interface{}, error) {
			return math.Cosh(args[0].(float64)), nil
		},
		"tanh": func(args ...interface{}) (interface{}, error) {
			return math.Tanh(args[0].(float64)), nil
		},
		"exp": func(args ...interface{}) (interface{}, error) {
			return math.Exp(args[0].(float64)), nil
		},
		"log": func(args ...interface{}) (interface{}, error) {
			return math.Log(args[0].(float64)), nil
		},
		"log10": func(args ...interface{}) (interface{}, error) {
			return math.Log10(args[0].(float64)), nil
		},
		"sqrt": func(args ...interface{}) (interface{}, error) {
			return math.Sqrt(args[0].(float64)), nil
		},
		"abs": func(args ...interface{}) (interface{}, error) {
			return math.Abs(args[0].(float64)), nil
		},
		"pow": func(args ...interface{}) (interface{}, error) {
			return math.Pow(args[0].(float64), args[1].(float64)), nil
		},
	}

	exprX, err := govaluate.NewEvaluableExpressionWithFunctions(formulaX, functions)
	if err != nil {
		return nil, err
	}
	exprY, err := govaluate.NewEvaluableExpressionWithFunctions(formulaY, functions)
	if err != nil {
		return nil, err
	}

	return &Formula{
		dt:       dt,
		isDiffur: isDiffur,
		exprX:    exprX,
		exprY:    exprY,
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

	parameters := make(map[string]interface{}, 8)

	toFloat64 := func(value interface{}) (float64, error) {
		result, ok := value.(float64)
		if !ok {
			return 0, fmt.Errorf("unexpected evaluation result type %T", value)
		}
		return result, nil
	}

	if f.isDiffur {
		getDerivatives := func(curX, curY float64, t float64) (float64, float64, error) {
			parameters["x"] = curX
			parameters["y"] = curY
			parameters["t"] = t

			resX, err := f.exprX.Evaluate(parameters)
			if err != nil {
				return 0, 0, err
			}
			resY, err := f.exprY.Evaluate(parameters)
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

		for t := f.tBegin; t <= f.tEnd; t += f.dt {
			dx1, dy1, err := getDerivatives(x, y, t)
			if err != nil {
				return zeroPoint, err
			}
			dx2, dy2, err := getDerivatives(x+dx1*f.dt/2, y+dy1*f.dt/2, t)
			if err != nil {
				return zeroPoint, err
			}
			dx3, dy3, err := getDerivatives(x+dx2*f.dt/2, y+dy2*f.dt/2, t)
			if err != nil {
				return zeroPoint, err
			}
			dx4, dy4, err := getDerivatives(x+dx3*f.dt, y+dy3*f.dt, t)
			if err != nil {
				return zeroPoint, err
			}

			x += (f.dt / 6) * (dx1 + dx2*2 + dx3*2 + dx4)
			y += (f.dt / 6) * (dy1 + dy2*2 + dy3*2 + dy4)
		}
	} else {
		parameters["x"] = x
		parameters["y"] = y
		resX, err := f.exprX.Evaluate(parameters)
		if err != nil {
			return zeroPoint, err
		}
		resY, err := f.exprY.Evaluate(parameters)
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
	min := compute.Point{req.Min.X, req.Min.Y}
	max := compute.Point{req.Max.X, req.Max.Y}

	formula, err := NewFormula(req.FormulaX, req.FormulaY, req.Eps, req.FormulaType == pb.FormulaType_DIFFUR)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}

	response := toPb2D(grid)
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
