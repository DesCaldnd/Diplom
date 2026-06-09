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
	nodes := grid.Nodes()
	payload := exportedGridExample{
		Name:      name,
		Min:       min,
		Max:       max,
		Epsilon:   epsilon,
		Basis:     basisNameForExport(basis),
		BuildType: buildTypeNameForExport(buildType),
		Nodes:     nodes,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(outDir, name+".json")
	fmt.Printf("%s has %d nodes\n", name, len(nodes))
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
	epsValues := []float64{1e-2, 5e-3, 1e-3, 1e-4, 1e-5, 5e-6, 1e-6}

	for _, eps := range epsValues {
		linearGrid, err := compute.NewAdaptiveSparseGrid(comparisonFunc, comparisonMin, comparisonMax, eps, nil, compute.BasisTypeLinear, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := saveGridJSONExample(outDir, fmt.Sprintf("basis_linear_eps_%g", eps), comparisonMin, comparisonMax, eps, compute.BasisTypeLinear, buildType, linearGrid); err != nil {
			fmt.Println(err)
			return
		}

		quadraticGrid, err := compute.NewAdaptiveSparseGrid(comparisonFunc, comparisonMin, comparisonMax, eps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := saveGridJSONExample(outDir, fmt.Sprintf("basis_quadratic_eps_%g", eps), comparisonMin, comparisonMax, eps, compute.BasisTypeQuadratic, buildType, quadraticGrid); err != nil {
			fmt.Println(err)
			return
		}
	}

	epsValues = []float64{0.01, 0.005, 0.001}
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
	flutterSimpleGrid, err := compute.NewAdaptiveSparseGrid(flutterSimple, flutterSimpleMin, flutterSimpleMax, 0.00001, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "flutter_simple_eps_0.00001", flutterSimpleMin, flutterSimpleMax, 0.00001, compute.BasisTypeQuadratic, buildType, flutterSimpleGrid); err != nil {
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
	flutterDiffGrid, err := compute.NewAdaptiveSparseGrid(flutterDiff, flutterDiffMin, flutterDiffMax, 0.00001, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "flutter_differential_eps_0.00001", flutterDiffMin, flutterDiffMax, 0.00001, compute.BasisTypeQuadratic, buildType, flutterDiffGrid); err != nil {
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

	var diffurMaxTimes = []float64{2.0, 5.0, 10.0, 20.0}

	for _, maxTime := range diffurMaxTimes {
		diffEq2D := func(state compute.Point, t float64) compute.Point {
			x := state[0]
			y := state[1]
			return compute.Point{-y / (1 + math.Sqrt(x*x+y*y)), -x / (1 + math.Sqrt(x*x+y*y))}
		}
		diffMinIter := compute.Point{-1.0, 0.0}
		diffMaxIter := compute.Point{1.0, 1.0}
		diffMinInter := compute.Point{-1.0, 0.0, 0}
		diffMaxInter := compute.Point{1.0, 1.0, maxTime}
		diffEps := 0.001

		diffIntervalAtTMax := func(arg compute.Point) (compute.Point, error) {
			return integrateRk4(diffEq2D, arg[0:2], 0.0, arg[2], int(arg[2]/0.04)), nil
		}
		diffIntervalGrid, err := compute.NewAdaptiveSparseGrid(diffIntervalAtTMax, diffMinInter, diffMaxInter, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := saveGridJSONExample(outDir, fmt.Sprintf("diffur_interval_tmax_%d_eps_0.001", int(maxTime)), diffMinInter, diffMaxInter, diffEps, compute.BasisTypeQuadratic, buildType, diffIntervalGrid); err != nil {
			fmt.Println(err)
			return
		}

		diffIterStep := func(state compute.Point) (compute.Point, error) {
			return integrateRk4(diffEq2D, state, 0.0, 1.0, 25), nil
		}
		diffIterGrid, err := compute.NewAdaptiveSparseGrid(diffIterStep, diffMinIter, diffMaxIter, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		nodeCount := 0
		for step := 1; step < int(maxTime); step++ {
			diffIterGrid, err = diffIterGrid.MakeNextIteration(diffIterStep, diffEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
			if err != nil {
				fmt.Println(err)
				return
			}
			nodeCount += len(diffIterGrid.Nodes())
		}
		if err := saveGridJSONExample(outDir, fmt.Sprintf("diffur_iterative_tmax_%d_eps_0.001", int(maxTime)), diffMinIter, diffMaxIter, diffEps, compute.BasisTypeQuadratic, buildType, diffIterGrid); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("diffur_iterative_tmax_%d_eps_0.001 has overall nodes %d\n", int(maxTime), nodeCount)
	}

	threeDimFunc := func(arg compute.Point) (compute.Point, error) {
		x := arg[0]
		y := arg[1]
		z := arg[2]

		u := x/2.0 - 1.0
		v := y/2.0 - 1.0
		w := z/2.0 - 1.0
		r2 := u*u + v*v + w*w
		radial := math.Sqrt(r2 + 0.08)
		swirl := math.Atan2(v+0.35*math.Sin(2.4*w), u+0.35*math.Cos(2.1*w))
		tube := math.Exp(-5.5 * math.Pow(radial-0.58-0.12*math.Sin(3.0*w+2.0*u*v), 2))
		core := math.Exp(-2.0 * r2)
		interaction := math.Sin(9.0*(u*v+v*w+w*u) + 3.5*swirl)
		braid := math.Cos(7.0*(u-v)*w + 4.0*math.Sin(2.0*u+v*w))
		ridge := math.Exp(-7.0 * math.Pow(u+0.45*math.Sin(2.6*v+1.7*w), 2))
		value := tube*interaction + 0.65*core*braid + 0.45*ridge*math.Sin(8.0*v*w+3.0*u)
		return compute.Point{value}, nil
	}
	threeDimMin, threeDimMax := makeDomain(3, 4.0)
	threeDimEps := 0.001
	threeDimGrid, err := compute.NewAdaptiveSparseGrid(threeDimFunc, threeDimMin, threeDimMax, threeDimEps, nil, compute.BasisTypeQuadratic, buildType, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := saveGridJSONExample(outDir, "three_dim_eps_0.001", threeDimMin, threeDimMax, threeDimEps, compute.BasisTypeQuadratic, buildType, threeDimGrid); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("exported grid json files")
	// Output: exported grid json files
}
