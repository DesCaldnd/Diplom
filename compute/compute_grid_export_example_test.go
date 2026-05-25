package compute_test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"Diplom/compute"
)

type exportedGridExample struct {
	Name      string             `json:"name"`
	Min       compute.Point      `json:"min"`
	Max       compute.Point      `json:"max"`
	Epsilon   float64            `json:"epsilon"`
	Basis     string             `json:"basis"`
	BuildType string             `json:"build_type"`
	Nodes     []compute.NodeInfo `json:"nodes"`
}

func basisNameForExport(basis compute.BasisType) string {
	switch basis {
	case compute.BasisTypeLinear:
		return "linear"
	case compute.BasisTypeQuadratic:
		return "quadratic"
	default:
		return fmt.Sprintf("basis_%d", basis)
	}
}

func buildTypeNameForExport(buildType compute.BuildType) string {
	switch buildType {
	case compute.BuildTypeSequential:
		return "sequential"
	case compute.BuildTypeParallel:
		return "parallel"
	default:
		return fmt.Sprintf("build_%d", buildType)
	}
}

func saveGridJSONExample(outDir string, name string, min, max compute.Point, epsilon float64, basis compute.BasisType, buildType compute.BuildType, grid *compute.AdaptiveSparseGrid) error {
	payload := exportedGridExample{
		Name:      name,
		Min:       min,
		Max:       max,
		Epsilon:   epsilon,
		Basis:     basisNameForExport(basis),
		BuildType: buildTypeNameForExport(buildType),
		Nodes:     grid.Nodes(),
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(outDir, name+".json")
	return os.WriteFile(path, data, 0o644)
}

func Example_export2DGridsForVisualization() {
	outDir := filepath.Join("..", "scripts", "grid_data")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Println(err)
		return
	}

	buildType := compute.BuildTypeSequential

	comparisonFunc := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		value := math.Sin(15*x) + 0.6*math.Cos(40*y) + 0.2*math.Sin(31*x+0.3) - 0.12*math.Cos(19*y-0.7) + math.Sin(x*y)
		return compute.Point{value}, nil
	}
	comparisonMin := compute.Point{0.0, 0.0}
	comparisonMax := compute.Point{math.Pi, math.Pi}

	linearGrid, err := compute.NewAdaptiveSparseGrid(comparisonFunc, comparisonMin, comparisonMax, 0.01, []compute.Point{{2, 2}, {2.2, 1}}, compute.BasisTypeLinear, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "basis_linear_eps_0.01", comparisonMin, comparisonMax, 0.01, compute.BasisTypeLinear, buildType, linearGrid); err != nil {
		fmt.Println(err)
		return
	}

	quadraticGrid, err := compute.NewAdaptiveSparseGrid(comparisonFunc, comparisonMin, comparisonMax, 0.01, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "basis_quadratic_eps_0.01", comparisonMin, comparisonMax, 0.01, compute.BasisTypeQuadratic, buildType, quadraticGrid); err != nil {
		fmt.Println(err)
		return
	}

	epsValues := []float64{0.01, 0.005, 0.001}
	for _, eps := range epsValues {
		grid, err := compute.NewAdaptiveSparseGrid(comparisonFunc, comparisonMin, comparisonMax, eps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := saveGridJSONExample(outDir, fmt.Sprintf("quadratic_eps_%g", eps), comparisonMin, comparisonMax, eps, compute.BasisTypeQuadratic, buildType, grid); err != nil {
			fmt.Println(err)
			return
		}
	}

	flutterSimple := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		u := x
		v := y * (math.Sin(2.5*x) + 0.35*math.Cos(6*y) + 0.2*math.Sin(5*x+0.4) - 0.15*math.Cos(9*y-0.2))
		return compute.Point{u, v}, nil
	}
	flutterSimpleMin := compute.Point{0.0, 0.0}
	flutterSimpleMax := compute.Point{4.0, 2.0}
	flutterSimpleGrid, err := compute.NewAdaptiveSparseGrid(flutterSimple, flutterSimpleMin, flutterSimpleMax, 0.001, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "flutter_simple_eps_0.001", flutterSimpleMin, flutterSimpleMax, 0.001, compute.BasisTypeQuadratic, buildType, flutterSimpleGrid); err != nil {
		fmt.Println(err)
		return
	}

	flutterDiff := func(arg compute.Point) (compute.Point, error) {
		x0 := arg[0]
		y0 := arg[1]
		state := compute.Point{x0, y0}
		time := 0.0
		dt := 0.02
		steps := 50
		for i := 0; i < steps; i++ {
			state[1] = state[1] + dt*(math.Sin(state[0])*state[1]+time)
			time += dt
		}
		return state, nil
	}
	flutterDiffMin := compute.Point{-2.0, 0.0}
	flutterDiffMax := compute.Point{5.0, 1.0}
	flutterDiffGrid, err := compute.NewAdaptiveSparseGrid(flutterDiff, flutterDiffMin, flutterDiffMax, 0.001, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "flutter_differential_eps_0.001", flutterDiffMin, flutterDiffMax, 0.001, compute.BasisTypeQuadratic, buildType, flutterDiffGrid); err != nil {
		fmt.Println(err)
		return
	}

	simpleDiff := func(arg compute.Point) (compute.Point, error) {
		x0 := arg[0]
		y0 := arg[1]
		state := compute.Point{x0, y0}
		time := 0.0
		maxTime := 1.0
		dt := 0.02
		steps := int(maxTime / dt)
		for i := 0; i < steps; i++ {
			state[0] = -state[1] / (1 + math.Sqrt(math.Pow(state[0], 2)+math.Pow(state[1], 2)))
			state[1] = -state[0] / (1 + math.Sqrt(math.Pow(state[0], 2)+math.Pow(state[1], 2)))
			time += dt
		}
		return state, nil
	}
	simpleDiffMin := compute.Point{-1.0, 0.0}
	simpleDiffMax := compute.Point{1.0, 1.0}
	simpleDiffGrid, err := compute.NewAdaptiveSparseGrid(simpleDiff, simpleDiffMin, simpleDiffMax, 0.001, []compute.Point{{0.7, 0.7}, {-0.5, 0.2}}, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "simple_differential_eps_0.001", simpleDiffMin, simpleDiffMax, 0.001, compute.BasisTypeQuadratic, buildType, simpleDiffGrid); err != nil {
		fmt.Println(err)
		return
	}

	diffEq2D := func(state compute.Point, t float64) compute.Point {
		x := state[0]
		y := state[1]
		return compute.Point{-y / (1 + math.Sqrt(x*x+y*y)), -x / (1 + math.Sqrt(x*x+y*y))}
	}
	diffMin := compute.Point{-1.0, 0.0}
	diffMax := compute.Point{1.0, 1.0}
	diffEps := 0.001
	diffTMax := 5.0

	diffIntervalAtTMax := func(arg compute.Point) (compute.Point, error) {
		return integrateRk4(diffEq2D, arg, 0.0, diffTMax, 50), nil
	}
	diffIntervalGrid, err := compute.NewAdaptiveSparseGrid(diffIntervalAtTMax, diffMin, diffMax, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "diffur_interval_tmax_5_eps_0.001", diffMin, diffMax, diffEps, compute.BasisTypeQuadratic, buildType, diffIntervalGrid); err != nil {
		fmt.Println(err)
		return
	}

	diffIterStep := func(state compute.Point) (compute.Point, error) {
		return integrateRk4(diffEq2D, state, 0.0, 1.0, 25), nil
	}
	diffIterGrid, err := compute.NewAdaptiveSparseGrid(diffIterStep, diffMin, diffMax, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for step := 1; step < int(diffTMax); step++ {
		diffIterGrid, err = diffIterGrid.MakeNextIteration(diffIterStep, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if err := saveGridJSONExample(outDir, "diffur_iterative_tmax_5_eps_0.001", diffMin, diffMax, diffEps, compute.BasisTypeQuadratic, buildType, diffIterGrid); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("exported grid json files")
	// Output: exported grid json files
}
