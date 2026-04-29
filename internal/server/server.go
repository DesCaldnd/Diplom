package server

import (
	"context"

	"Diplom/compute"
	pb "Diplom/gridProto"
	"github.com/Knetic/govaluate"
)

import "math"

type Formula struct {
	dt       compute.ScalarType
	isDiffur bool
	exprX    *govaluate.EvaluableExpression
	exprY    *govaluate.EvaluableExpression
	tBegin   compute.ScalarType
	tEnd     compute.ScalarType
}

func NewFormula(formulaX, formulaY string, dt compute.ScalarType, isDiffur bool) (*Formula, error) {
	exprX, err := govaluate.NewEvaluableExpression(formulaX)
	if err != nil {
		return nil, err
	}
	exprY, err := govaluate.NewEvaluableExpression(formulaY)
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
	x := arg[0]
	y := arg[1]

	parameters := make(map[string]interface{}, 8)

	if f.isDiffur {
		for t := f.tBegin; t <= f.tEnd; t += f.dt {
			getDerivatives := func(curX, curY float64) (float64, float64) {
				parameters["x"] = curX
				parameters["y"] = curY
				parameters["t"] = t
				resX, _ := f.exprX.Evaluate(parameters)
				resY, _ := f.exprY.Evaluate(parameters)
				return resX.(float64), resY.(float64)
			}

			dx1, dy1 := getDerivatives(x, y)
			dx2, dy2 := getDerivatives(x+dx1*f.dt/2, y+dy1*f.dt/2)
			dx3, dy3 := getDerivatives(x+dx2*f.dt/2, y+dy2*f.dt/2)
			dx4, dy4 := getDerivatives(x+dx3*f.dt, y+dy3*f.dt)

			x += (f.dt / 6) * (dx1 + dx2*2 + dx3*2 + dx4)
			y += (f.dt / 6) * (dy1 + dy2*2 + dy3*2 + dy4)
		}
	} else {
		parameters["x"] = x
		parameters["y"] = y
		resX, err := f.exprX.Evaluate(parameters)
		if err != nil {
			return nil, err
		}
		resY, err := f.exprY.Evaluate(parameters)
		if err != nil {
			return nil, err
		}
		x = resX.(float64)
		y = resY.(float64)
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
		func(p compute.Point) compute.Point {
			res, _ := formula.Evaluate(p)
			return res
		},
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
			func(p compute.Point) compute.Point {
				res, _ := formula.Evaluate(p)
				return res
			},
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
