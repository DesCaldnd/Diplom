package compute

import (
	"context"
	"fmt"
	"math"
	"sync"
)

type BuildType int64

const (
	BuildTypeSequential BuildType = iota
	BuildTypeParallel
)

type AdaptiveSparseGrid struct {
	min         Point
	max         Point
	basisType   BasisType
	inDim       int64
	outDim      int64
	entryPoints []entryPoint
	nodes       map[string]node
}

func NewAdaptiveSparseGrid(
	funcEval func(Point) (Point, error),
	min, max Point,
	epsilon ScalarType,
	anchorPoints []Point,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	return NewAdaptiveSparseGridWithContext(context.Background(), funcEval, min, max, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func NewAdaptiveSparseGridWithContext(
	ctx context.Context,
	funcEval func(Point) (Point, error),
	min, max Point,
	epsilon ScalarType,
	anchorPoints []Point,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	inDim := int64(len(min))
	if int64(len(max)) != inDim {
		return nil, fmt.Errorf("min and max must have the same dimension (IN_DIM)")
	}
	for _, anchor := range anchorPoints {
		if int64(len(anchor)) != inDim {
			return nil, fmt.Errorf("anchor points must have the same dimension as min/max (IN_DIM)")
		}
	}

	sampleOut, err := funcEval(min)
	if err != nil {
		return nil, err
	}
	outDim := int64(len(sampleOut))

	grid := &AdaptiveSparseGrid{
		min:       min,
		max:       max,
		basisType: basisType,
		inDim:     inDim,
		outDim:    outDim,
		nodes:     make(map[string]node),
	}

	err = grid.checkConstraints()
	if err != nil {
		return nil, err
	}

	anchors, err := grid.calculateOriginalFunctionAtPoints(funcEval, anchorPoints)
	if err != nil {
		return nil, err
	}

	err = grid.build(ctx, funcEval, epsilon, anchors, buildType, maxLevel, maxNodesInGrid)
	if err != nil {
		return nil, err
	}

	return grid, nil
}

func (g *AdaptiveSparseGrid) checkConstraints() error {
	for i := int64(0); i < g.inDim; i++ {
		if math.Abs(g.max[i]-g.min[i]) < 1e-9 {
			return fmt.Errorf("invalid bounds at dimension %d, %f ~ %f", i, g.min[i], g.max[i])
		} else if g.max[i] < g.min[i] {
			g.min[i], g.max[i] = g.max[i], g.min[i]
		}
	}
	return nil
}

func (g *AdaptiveSparseGrid) toReal(unit Point) Point {
	real := make(Point, g.inDim)
	for i := int64(0); i < g.inDim; i++ {
		real[i] = g.min[i] + unit[i]*(g.max[i]-g.min[i])
	}
	return real
}

func (g *AdaptiveSparseGrid) toUnit(real Point) Point {
	unit := make(Point, g.inDim)
	for i := int64(0); i < g.inDim; i++ {
		unit[i] = (real[i] - g.min[i]) / (g.max[i] - g.min[i])
	}
	return unit
}

type anchorPair struct {
	arg Point
	ans Point
}

func (g *AdaptiveSparseGrid) calculateOriginalFunctionAtPoints(funcEval func(Point) (Point, error), anchorPoints []Point) ([]anchorPair, error) {
	result := make([]anchorPair, 0, len(anchorPoints))
	for _, arg := range anchorPoints {
		ans, err := funcEval(arg)
		if err != nil {
			return nil, err
		}
		result = append(result, anchorPair{
			arg: g.toUnit(arg),
			ans: ans,
		})
	}
	return result, nil
}

func (g *AdaptiveSparseGrid) checkEvaluationPoint(real Point) error {
	for i := int64(0); i < g.inDim; i++ {
		if real[i] < g.min[i]-1e-9 || real[i] > g.max[i]+1e-9 {
			return fmt.Errorf("evaluation point out of bounds at dimension %d", i)
		}
	}
	return nil
}

func (g *AdaptiveSparseGrid) Evaluate(x Point) (Point, error) {
	if int64(len(x)) != g.inDim {
		return nil, fmt.Errorf("input point must have dimension %d", g.inDim)
	}
	if err := g.checkEvaluationPoint(x); err != nil {
		return nil, err
	}
	return g.evaluateForDim(g.toUnit(x), g.inDim+1), nil
}

func (g *AdaptiveSparseGrid) MakeNextIteration(
	funcEval func(Point) (Point, error),
	epsilon ScalarType,
	anchorPoints []Point,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	return g.MakeNextIterationWithContext(context.Background(), funcEval, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func (g *AdaptiveSparseGrid) MakeNextIterationWithContext(
	ctx context.Context,
	funcEval func(Point) (Point, error),
	epsilon ScalarType,
	anchorPoints []Point,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	for _, anchor := range anchorPoints {
		if int64(len(anchor)) != g.inDim {
			return nil, fmt.Errorf("anchor points must have the same dimension as min/max (IN_DIM)")
		}
	}
	wrapper := func(x Point) (Point, error) {
		evalRes, err := g.Evaluate(x)
		if err != nil {
			return nil, err
		}
		return funcEval(evalRes)
	}
	return NewAdaptiveSparseGridWithContext(ctx, wrapper, g.min, g.max, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func (g *AdaptiveSparseGrid) BasisType() BasisType {
	return g.basisType
}

func (g *AdaptiveSparseGrid) Min() Point {
	return g.min
}

func (g *AdaptiveSparseGrid) Max() Point {
	return g.max
}

type NodeInfo struct {
	Level       []int64
	Index       []int64
	CenterUnit  Point
	Alpha       Point
	HasChildren bool
}

func (g *AdaptiveSparseGrid) Nodes() []NodeInfo {
	res := make([]NodeInfo, 0, len(g.nodes))
	for _, n := range g.nodes {
		res = append(res, NodeInfo{
			Level:       n.key.level,
			Index:       n.key.index,
			CenterUnit:  n.centerUnit,
			Alpha:       n.alpha,
			HasChildren: n.hasChildren,
		})
	}
	return res
}

func (g *AdaptiveSparseGrid) EntryPoints() []NodeInfo {
	res := make([]NodeInfo, 0, len(g.entryPoints))
	for _, ep := range g.entryPoints {
		res = append(res, NodeInfo{
			Level:       ep.node.key.level,
			Index:       ep.node.key.index,
			CenterUnit:  ep.node.centerUnit,
			Alpha:       ep.node.alpha,
			HasChildren: ep.node.hasChildren,
		})
	}
	return res
}

func nextPermutation(arr []int64) bool {
	k := -1
	for i := len(arr) - 2; i >= 0; i-- {
		if arr[i] < arr[i+1] {
			k = i
			break
		}
	}
	if k == -1 {
		return false
	}
	l := -1
	for i := len(arr) - 1; i > k; i-- {
		if arr[k] < arr[i] {
			l = i
			break
		}
	}
	arr[k], arr[l] = arr[l], arr[k]
	for i, j := k+1, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return true
}

func expandIndices(key *gridKey, permutation int64, inDim int64) {
	idx := int64(0)
	for i := int64(0); i < inDim; i++ {
		if key.level[i] == 1 {
			key.index[i] = 1
		} else {
			if (permutation & (1 << idx)) != 0 {
				key.index[i] = 1
			} else {
				key.index[i] = 0
			}
			idx++
		}
	}
}

type buildResult struct {
	node     node
	newNodes map[string]node
	err      error
}

func (g *AdaptiveSparseGrid) build(
	ctx context.Context,
	funcEval func(Point) (Point, error),
	epsilon ScalarType,
	anchors []anchorPair,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) error {
	for i := int64(0); i <= g.inDim; i++ {
		key := gridKey{
			level: make([]int64, g.inDim),
			index: make([]int64, g.inDim),
		}

		for j := int64(0); j < i; j++ {
			key.level[g.inDim-j-1] = 1
		}

		var results []buildResult
		var wg sync.WaitGroup
		var mu sync.Mutex

		for {
			indexPermutations := int64(1) << (g.inDim - i)

			for j := int64(0); j < indexPermutations; j++ {
				currentKey := gridKey{
					level: make([]int64, g.inDim),
					index: make([]int64, g.inDim),
				}
				copy(currentKey.level, key.level)

				expandIndices(&currentKey, j, g.inDim)

				if buildType == BuildTypeParallel {
					wg.Add(1)
					go func(k gridKey) {
						defer wg.Done()
						res := g.buildGridTask(ctx, funcEval, epsilon, anchors, k, i, maxLevel, maxNodesInGrid)
						mu.Lock()
						results = append(results, res)
						mu.Unlock()
					}(currentKey)
				} else {
					res := g.buildGridTask(ctx, funcEval, epsilon, anchors, currentKey, i, maxLevel, maxNodesInGrid)
					results = append(results, res)
				}
			}

			if !nextPermutation(key.level) {
				break
			}
		}

		if buildType == BuildTypeParallel {
			wg.Wait()
		}

		for _, res := range results {
			if res.err != nil {
				return res.err
			}
			ep := entryPoint{
				node:       res.node,
				dimensions: i,
			}
			g.entryPoints = append(g.entryPoints, ep)
			for k, v := range res.newNodes {
				g.nodes[k] = v
			}
		}
	}
	return nil
}

func (g *AdaptiveSparseGrid) buildGridTask(
	ctx context.Context,
	funcEval func(Point) (Point, error),
	epsilon ScalarType,
	anchors []anchorPair,
	key gridKey,
	dimension int64,
	maxLevel int64,
	maxNodesInGrid int64,
) buildResult {
	tmp := make(map[string]node)
	n, err := g.createNode(funcEval, key, nil, dimension, tmp)
	if err != nil {
		return buildResult{err: err}
	}
	if n == nil {
		return buildResult{err: fmt.Errorf("failed to create initial node")}
	}

	startNode := tmp[n.String()]

	resNode, newNodes, err := g.buildGrid(ctx, funcEval, epsilon, anchors, startNode, dimension, maxLevel, maxNodesInGrid)
	return buildResult{
		node:     resNode,
		newNodes: newNodes,
		err:      err,
	}
}

func (g *AdaptiveSparseGrid) buildGrid(
	ctx context.Context,
	funcEval func(Point) (Point, error),
	epsilon ScalarType,
	anchors []anchorPair,
	entryPoint node,
	dimension int64,
	maxLevel int64,
	maxNodesInGrid int64,
) (node, map[string]node, error) {
	nodeQueue := []gridKey{entryPoint.key}
	newNodes := make(map[string]node)
	newNodes[entryPoint.key.String()] = entryPoint

	entryPointKeyString := entryPoint.key.String()

	currentMaxLevel := int64(0)
	for i := int64(0); i < g.inDim; i++ {
		if entryPoint.key.level[i] > currentMaxLevel {
			currentMaxLevel = entryPoint.key.level[i]
		}
	}

	for len(nodeQueue) > 0 {
		select {
		case <-ctx.Done():
			return node{}, nil, fmt.Errorf("interrupted at node count: %d, dimension: %d, max level: %d", len(newNodes), dimension, currentMaxLevel)
		default:
		}

		currentKey := nodeQueue[0]
		nodeQueue = nodeQueue[1:]

		currentNode := newNodes[currentKey.String()]
		activeEntryNode := newNodes[entryPointKeyString]

		canContinue := false

		maxLvl := int64(0)
		for i := int64(0); i < g.inDim; i++ {
			if currentNode.key.level[i] > maxLvl {
				maxLvl = currentNode.key.level[i]
			}
		}

		canContinueForce := maxLvl >= 64
		directions := make([]direction, g.inDim)
		for i := range directions {
			directions[i] = directionBoth
		}

		if canContinueForce || currentNode.alpha.Length() <= epsilon {
			for i := range directions {
				directions[i] = directionNone
			}
			canContinue = true
		}

		if canContinue && !canContinueForce {
			for _, anchor := range anchors {
				if currentNode.isPointInAffectZone(anchor.arg, g.inDim) {
					evalRes := g.evaluateForDimAndEntryPoint(anchor.arg, dimension, &activeEntryNode, newNodes)
					evalRes.Sub(anchor.ans)
					if evalRes.Length() >= epsilon {
						canContinue = false
						affectDirs := getAffectDirection(anchor.arg, currentNode.key.level, currentNode.key.index, g.inDim)
						for i := int64(0); i < g.inDim; i++ {
							directions[i] |= affectDirs[i]
						}
					}
				}
			}
		}

		if canContinue || canContinueForce {
			continue
		}

		currentNode.hasChildren = true
		newNodes[currentKey.String()] = currentNode

		for i := int64(0); i < g.inDim; i++ {
			if currentNode.key.level[i] == 0 || (maxLevel > 0 && currentNode.key.level[i] >= maxLevel) {
				continue
			}

			keyLeft := gridKey{
				level: make([]int64, g.inDim),
				index: make([]int64, g.inDim),
			}
			copy(keyLeft.level, currentNode.key.level)
			copy(keyLeft.index, currentNode.key.index)

			keyRight := gridKey{
				level: make([]int64, g.inDim),
				index: make([]int64, g.inDim),
			}
			copy(keyRight.level, currentNode.key.level)
			copy(keyRight.index, currentNode.key.index)

			keyLeft.level[i] = currentNode.key.level[i] + 1
			keyRight.level[i] = currentNode.key.level[i] + 1

			if keyLeft.level[i] > currentMaxLevel {
				currentMaxLevel = keyLeft.level[i]
			}

			keyLeft.index[i] = 2*currentNode.key.index[i] - 1
			keyRight.index[i] = 2*currentNode.key.index[i] + 1

			if (directionLeft & directions[i]) != directionNone {
				leftNodeKey, err := g.createNode(funcEval, keyLeft, &activeEntryNode, dimension, newNodes)
				if err != nil {
					return node{}, nil, err
				}
				if leftNodeKey != nil {
					nodeQueue = append(nodeQueue, *leftNodeKey)
				}
			}

			if (directionRight & directions[i]) != directionNone {
				rightNodeKey, err := g.createNode(funcEval, keyRight, &activeEntryNode, dimension, newNodes)
				if err != nil {
					return node{}, nil, err
				}
				if rightNodeKey != nil {
					nodeQueue = append(nodeQueue, *rightNodeKey)
				}
			}
		}

		if maxNodesInGrid > 0 && int64(len(newNodes)) >= maxNodesInGrid {
			break
		}
	}

	return newNodes[entryPointKeyString], newNodes, nil
}

func (g *AdaptiveSparseGrid) evaluateForDimAndEntryPoint(x Point, maxGridDim int64, entryPoint *node, additionalNodes map[string]node) Point {
	interp := g.evaluateForDim(x, maxGridDim)
	if entryPoint != nil {
		processedNodes := make(map[string]struct{})
		interp.Add(g.evaluateRecursive(x, *entryPoint, processedNodes, additionalNodes))
	}
	return interp
}

func (g *AdaptiveSparseGrid) createNode(
	funcEval func(Point) (Point, error),
	key gridKey,
	entryPoint *node,
	dimension int64,
	additionalNodes map[string]node,
) (*gridKey, error) {
	if _, exists := additionalNodes[key.String()]; exists {
		return nil, nil
	}

	n := node{
		key:        key,
		centerUnit: make(Point, g.inDim),
	}
	for i := int64(0); i < g.inDim; i++ {
		n.centerUnit[i] = getCoord(key.level[i], key.index[i])
	}

	realArg := g.toReal(n.centerUnit)
	etalon, err := funcEval(realArg)
	if err != nil {
		return nil, err
	}

	if int64(len(etalon)) != g.outDim {
		return nil, fmt.Errorf("function returned point of dimension %d, expected %d (OUT_DIM)", len(etalon), g.outDim)
	}

	interp := g.evaluateForDimAndEntryPoint(n.centerUnit, dimension, entryPoint, additionalNodes)

	etalon.Sub(interp)
	n.alpha = etalon
	additionalNodes[key.String()] = n

	return &key, nil
}

func (g *AdaptiveSparseGrid) evaluateForDim(x Point, maxGridDim int64) Point {
	answer := make(Point, g.outDim)
	processedNodes := make(map[string]struct{})
	for _, ep := range g.entryPoints {
		if ep.dimensions >= maxGridDim {
			break
		}
		clear(processedNodes)
		answer.Add(g.evaluateRecursive(x, ep.node, processedNodes, nil))
	}
	return answer
}

func (g *AdaptiveSparseGrid) evaluateRecursive(x Point, n node, processedNodes map[string]struct{}, additionalNodes map[string]node) Point {
	if _, exists := processedNodes[n.key.String()]; exists {
		return make(Point, g.outDim)
	}

	processedNodes[n.key.String()] = struct{}{}
	basis := evalBasis(x, n.key.level, n.key.index, g.basisType, g.inDim)

	if math.Abs(basis) < 1e-9 {
		return make(Point, g.outDim)
	}

	value := n.alpha.Clone()
	value.MulScalar(basis)

	if n.hasChildren {
		for i := int64(0); i < g.inDim; i++ {
			if n.key.level[i] != 0 {
				child, exists := g.getChildForDimAndArg(n, x, i, additionalNodes)
				if exists {
					value.Add(g.evaluateRecursive(x, child, processedNodes, additionalNodes))
				}
			}
		}
	}
	return value
}

func (g *AdaptiveSparseGrid) getChildForDimAndArg(parent node, x Point, dimension int64, additionalNodes map[string]node) (node, bool) {
	key := gridKey{
		level: make([]int64, g.inDim),
		index: make([]int64, g.inDim),
	}
	copy(key.level, parent.key.level)
	copy(key.index, parent.key.index)

	key.level[dimension] += 1

	dir := int64(1)
	if x[dimension] < parent.centerUnit[dimension] {
		dir = -1
	}
	key.index[dimension] = 2*key.index[dimension] + dir

	if child, exists := g.nodes[key.String()]; exists {
		return child, true
	}
	if additionalNodes != nil {
		if child, exists := additionalNodes[key.String()]; exists {
			return child, true
		}
	}
	return node{}, false
}
