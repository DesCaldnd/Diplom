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

const buildGridContextCheckPeriod = 200

type AdaptiveSparseGrid struct {
	min         StaticPoint
	max         StaticPoint
	span        StaticPoint
	invSpan     StaticPoint
	basisType   BasisType
	inDim       int64
	outDim      int64
	entryPoints []entryPoint
	nodes       map[gridKey]node
}

func NewAdaptiveSparseGrid(
	funcEval func(StaticPoint) (StaticPoint, error),
	min, max StaticPoint,
	epsilon ScalarType,
	anchorPoints []StaticPoint,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	return NewAdaptiveSparseGridWithContext(context.Background(), funcEval, min, max, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func NewAdaptiveSparseGridWithContext(
	ctx context.Context,
	funcEval func(StaticPoint) (StaticPoint, error),
	min, max StaticPoint,
	epsilon ScalarType,
	anchorPoints []StaticPoint,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	inDim := min.dim
	if inDim > maxSupportedDim {
		return nil, fmt.Errorf("input dimension %d exceeds maximum supported dimension %d", inDim, maxSupportedDim)
	}
	if max.dim != inDim {
		return nil, fmt.Errorf("min and max must have the same dimension (IN_DIM)")
	}
	for _, anchor := range anchorPoints {
		if anchor.dim != inDim {
			return nil, fmt.Errorf("anchor points must have the same dimension as min/max (IN_DIM)")
		}
	}

	sampleOut, err := funcEval(min)
	if err != nil {
		return nil, err
	}
	outDim := sampleOut.dim

	grid := &AdaptiveSparseGrid{
		min:       min,
		max:       max,
		basisType: basisType,
		inDim:     inDim,
		outDim:    outDim,
		nodes:     make(map[gridKey]node),
	}

	err = grid.checkConstraints()
	if err != nil {
		return nil, err
	}
	grid.fillSpanCache()

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
		if math.Abs(g.max.data[i]-g.min.data[i]) < 1e-9 {
			return fmt.Errorf("invalid bounds at dimension %d, %f ~ %f", i, g.min.data[i], g.max.data[i])
		} else if g.max.data[i] < g.min.data[i] {
			g.min.data[i], g.max.data[i] = g.max.data[i], g.min.data[i]
		}
	}
	return nil
}

func (g *AdaptiveSparseGrid) fillSpanCache() {
	g.span.dim = g.inDim
	g.invSpan.dim = g.inDim
	for i := int64(0); i < g.inDim; i++ {
		span := g.max.data[i] - g.min.data[i]
		g.span.data[i] = span
		g.invSpan.data[i] = 1.0 / span
	}
}

func (g *AdaptiveSparseGrid) toReal(unit StaticPoint) StaticPoint {
	var real StaticPoint
	real.dim = g.inDim
	for i := int64(0); i < g.inDim; i++ {
		real.data[i] = g.min.data[i] + unit.data[i]*g.span.data[i]
	}
	return real
}

func (g *AdaptiveSparseGrid) toUnit(real StaticPoint) StaticPoint {
	var unit StaticPoint
	unit.dim = g.inDim
	for i := int64(0); i < g.inDim; i++ {
		unit.data[i] = (real.data[i] - g.min.data[i]) * g.invSpan.data[i]
	}
	return unit
}

type anchorPair struct {
	arg StaticPoint
	ans StaticPoint
}

func (g *AdaptiveSparseGrid) calculateOriginalFunctionAtPoints(funcEval func(StaticPoint) (StaticPoint, error), anchorPoints []StaticPoint) ([]anchorPair, error) {
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

func (g *AdaptiveSparseGrid) checkEvaluationPoint(real StaticPoint) error {
	for i := int64(0); i < g.inDim; i++ {
		if real.data[i] < g.min.data[i]-1e-9 || real.data[i] > g.max.data[i]+1e-9 {
			return fmt.Errorf("evaluation point out of bounds at dimension %d", i)
		}
	}
	return nil
}

func (g *AdaptiveSparseGrid) Evaluate(x StaticPoint) (StaticPoint, error) {
	if x.dim != g.inDim {
		return StaticPoint{}, fmt.Errorf("input point must have dimension %d", g.inDim)
	}
	if err := g.checkEvaluationPoint(x); err != nil {
		return StaticPoint{}, err
	}
	return g.evaluateForDim(g.toUnit(x), g.inDim+1), nil
}

func (g *AdaptiveSparseGrid) MakeNextIteration(
	funcEval func(StaticPoint) (StaticPoint, error),
	epsilon ScalarType,
	anchorPoints []StaticPoint,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	return g.MakeNextIterationWithContext(context.Background(), funcEval, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func (g *AdaptiveSparseGrid) MakeNextIterationWithContext(
	ctx context.Context,
	funcEval func(StaticPoint) (StaticPoint, error),
	epsilon ScalarType,
	anchorPoints []StaticPoint,
	basisType BasisType,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) (*AdaptiveSparseGrid, error) {
	for _, anchor := range anchorPoints {
		if anchor.dim != g.inDim {
			return nil, fmt.Errorf("anchor points must have the same dimension as min/max (IN_DIM)")
		}
	}
	wrapper := func(x StaticPoint) (StaticPoint, error) {
		evalRes, err := g.Evaluate(x)
		if err != nil {
			return StaticPoint{}, err
		}
		return funcEval(evalRes)
	}
	return NewAdaptiveSparseGridWithContext(ctx, wrapper, g.min, g.max, epsilon, anchorPoints, basisType, buildType, maxLevel, maxNodesInGrid)
}

func (g *AdaptiveSparseGrid) BasisType() BasisType {
	return g.basisType
}

func (g *AdaptiveSparseGrid) Min() StaticPoint {
	return g.min
}

func (g *AdaptiveSparseGrid) Max() StaticPoint {
	return g.max
}

type NodeInfo struct {
	Level       []int64
	Index       []int64
	CenterUnit  StaticPoint
	Alpha       StaticPoint
	HasChildren bool
}

func (g *AdaptiveSparseGrid) Nodes() []NodeInfo {
	res := make([]NodeInfo, 0, len(g.nodes))
	for _, n := range g.nodes {
		res = append(res, nodeInfoFromNode(n, g.inDim))
	}
	return res
}

func (g *AdaptiveSparseGrid) EntryPoints() []NodeInfo {
	res := make([]NodeInfo, 0, len(g.entryPoints))
	for _, ep := range g.entryPoints {
		res = append(res, nodeInfoFromNode(ep.node, g.inDim))
	}
	return res
}

func nodeInfoFromNode(n node, inDim int64) NodeInfo {
	level := make([]int64, inDim)
	index := make([]int64, inDim)
	for i := int64(0); i < inDim; i++ {
		level[i] = n.key.level[i]
		index[i] = n.key.index[i]
	}
	return NodeInfo{
		Level:       level,
		Index:       index,
		CenterUnit:  n.centerUnit,
		Alpha:       n.alpha,
		HasChildren: n.hasChildren,
	}
}

func nextPermutation(arr *[maxSupportedDim]int64, length int64) bool {
	k := -1
	for i := int(length) - 2; i >= 0; i-- {
		if arr[i] < arr[i+1] {
			k = i
			break
		}
	}
	if k == -1 {
		return false
	}
	l := -1
	for i := int(length) - 1; i > k; i-- {
		if arr[k] < arr[i] {
			l = i
			break
		}
	}
	arr[k], arr[l] = arr[l], arr[k]
	for i, j := k+1, int(length)-1; i < j; i, j = i+1, j-1 {
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
	newNodes map[gridKey]node
	err      error
}

func (g *AdaptiveSparseGrid) build(
	ctx context.Context,
	funcEval func(StaticPoint) (StaticPoint, error),
	epsilon ScalarType,
	anchors []anchorPair,
	buildType BuildType,
	maxLevel int64,
	maxNodesInGrid int64,
) error {
	for i := int64(0); i <= g.inDim; i++ {
		var key gridKey
		for j := int64(0); j < i; j++ {
			key.level[g.inDim-j-1] = 1
		}

		results := make([]buildResult, 0)
		if buildType == BuildTypeParallel {
			resultCh := make(chan buildResult)
			var wg sync.WaitGroup
			for {
				indexPermutations := int64(1) << (g.inDim - i)
				for j := int64(0); j < indexPermutations; j++ {
					currentKey := key
					expandIndices(&currentKey, j, g.inDim)
					wg.Add(1)
					go func(k gridKey, dimension int64) {
						defer wg.Done()
						resultCh <- g.buildGridTask(ctx, funcEval, epsilon, anchors, k, dimension, maxLevel, maxNodesInGrid)
					}(currentKey, i)
				}
				if !nextPermutation(&key.level, g.inDim) {
					break
				}
			}
			go func() {
				wg.Wait()
				close(resultCh)
			}()
			for res := range resultCh {
				results = append(results, res)
			}
		} else {
			for {
				indexPermutations := int64(1) << (g.inDim - i)
				for j := int64(0); j < indexPermutations; j++ {
					currentKey := key
					expandIndices(&currentKey, j, g.inDim)
					results = append(results, g.buildGridTask(ctx, funcEval, epsilon, anchors, currentKey, i, maxLevel, maxNodesInGrid))
				}
				if !nextPermutation(&key.level, g.inDim) {
					break
				}
			}
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
	funcEval func(StaticPoint) (StaticPoint, error),
	epsilon ScalarType,
	anchors []anchorPair,
	key gridKey,
	dimension int64,
	maxLevel int64,
	maxNodesInGrid int64,
) buildResult {
	tmp := make(map[gridKey]node)
	n, err := g.createNode(funcEval, key, nil, dimension, tmp)
	if err != nil {
		return buildResult{err: err}
	}
	if n == nil {
		return buildResult{err: fmt.Errorf("failed to create initial node")}
	}

	startNode := tmp[*n]
	resNode, newNodes, err := g.buildGrid(ctx, funcEval, epsilon, anchors, startNode, dimension, maxLevel, maxNodesInGrid)
	return buildResult{
		node:     resNode,
		newNodes: newNodes,
		err:      err,
	}
}

func (g *AdaptiveSparseGrid) buildGrid(
	ctx context.Context,
	funcEval func(StaticPoint) (StaticPoint, error),
	epsilon ScalarType,
	anchors []anchorPair,
	entryPoint node,
	dimension int64,
	maxLevel int64,
	maxNodesInGrid int64,
) (node, map[gridKey]node, error) {
	nodeQueue := []gridKey{entryPoint.key}
	newNodes := make(map[gridKey]node)
	newNodes[entryPoint.key] = entryPoint
	entryPointKey := entryPoint.key
	epsilonSquared := epsilon * epsilon

	currentMaxLevel := entryPoint.maxLevel
	var directions [maxSupportedDim]direction
	var affectDirs [maxSupportedDim]direction
	iterationsSinceContextCheck := 0

	for len(nodeQueue) > 0 {
		iterationsSinceContextCheck++
		if iterationsSinceContextCheck >= buildGridContextCheckPeriod {
			iterationsSinceContextCheck = 0
			select {
			case <-ctx.Done():
				return node{}, nil, fmt.Errorf("interrupted at node count: %d, dimension: %d, max level: %d", len(newNodes), dimension, currentMaxLevel)
			default:
			}
		}

		currentKey := nodeQueue[0]
		nodeQueue = nodeQueue[1:]

		currentNode := newNodes[currentKey]
		activeEntryNode := newNodes[entryPointKey]

		canContinue := false
		canContinueForce := currentNode.maxLevel >= 64
		for idx := range directions {
			directions[idx] = directionBoth
		}

		if canContinueForce || currentNode.alpha.LengthSquared() <= epsilonSquared {
			for idx := range directions {
				directions[idx] = directionNone
			}
			canContinue = true
		}

		if canContinue && !canContinueForce {
			for _, anchor := range anchors {
				if currentNode.isPointInAffectZone(anchor.arg) {
					evalRes := g.evaluateForDimAndEntryPoint(anchor.arg, dimension, &activeEntryNode, newNodes)
					evalRes.Sub(anchor.ans)
					if evalRes.LengthSquared() >= epsilonSquared {
						canContinue = false
						fillAffectDirections(anchor.arg, currentNode, &affectDirs)
						for idx := int64(0); idx < g.inDim; idx++ {
							directions[idx] |= affectDirs[idx]
						}
					}
				}
			}
		}

		if canContinue || canContinueForce {
			continue
		}

		currentNode.hasChildren = true
		newNodes[currentKey] = currentNode

		for i := int64(0); i < g.inDim; i++ {
			if currentNode.key.level[i] == 0 || (maxLevel > 0 && currentNode.key.level[i] >= maxLevel) {
				continue
			}

			keyLeft := currentNode.key
			keyRight := currentNode.key
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

	return newNodes[entryPointKey], newNodes, nil
}

func (g *AdaptiveSparseGrid) evaluateForDimAndEntryPoint(x StaticPoint, maxGridDim int64, entryPoint *node, additionalNodes map[gridKey]node) StaticPoint {
	interp := g.evaluateForDim(x, maxGridDim)
	if entryPoint != nil {
		processedNodes := make(map[gridKey]struct{}, len(additionalNodes)+1)
		g.evaluateRecursiveInto(&interp, x, *entryPoint, processedNodes, additionalNodes)
	}
	return interp
}

func (g *AdaptiveSparseGrid) createNode(
	funcEval func(StaticPoint) (StaticPoint, error),
	key gridKey,
	entryPoint *node,
	dimension int64,
	additionalNodes map[gridKey]node,
) (*gridKey, error) {
	if _, exists := additionalNodes[key]; exists {
		return nil, nil
	}

	n := node{
		key:        key,
		centerUnit: NewStaticPoint(g.inDim),
		alpha:      NewStaticPoint(g.outDim),
	}
	for i := int64(0); i < g.inDim; i++ {
		coord := getCoord(key.level[i], key.index[i])
		n.centerUnit.data[i] = coord
		if key.level[i] != 0 {
			unit := 1.0 / float64(int64(1)<<key.level[i])
			n.leftBound[i] = coord - unit
			n.rightBound[i] = coord + unit
		}
		if key.level[i] > n.maxLevel {
			n.maxLevel = key.level[i]
		}
	}

	etalon, err := funcEval(g.toReal(n.centerUnit))
	if err != nil {
		return nil, err
	}

	if etalon.dim != g.outDim {
		return nil, fmt.Errorf("function returned point of dimension %d, expected %d (OUT_DIM)", etalon.dim, g.outDim)
	}

	interp := g.evaluateForDimAndEntryPoint(n.centerUnit, dimension, entryPoint, additionalNodes)
	etalon.Sub(interp)
	n.alpha = etalon
	additionalNodes[key] = n

	return &key, nil
}

func (g *AdaptiveSparseGrid) evaluateForDim(x StaticPoint, maxGridDim int64) StaticPoint {
	answer := NewStaticPoint(g.outDim)
	for _, ep := range g.entryPoints {
		if ep.dimensions >= maxGridDim {
			break
		}
		processedNodes := make(map[gridKey]struct{})
		g.evaluateRecursiveInto(&answer, x, ep.node, processedNodes, nil)
	}
	return answer
}

func (g *AdaptiveSparseGrid) evaluateRecursiveInto(dst *StaticPoint, x StaticPoint, n node, processedNodes map[gridKey]struct{}, additionalNodes map[gridKey]node) {
	if _, exists := processedNodes[n.key]; exists {
		return
	}
	processedNodes[n.key] = struct{}{}

	basis := evalBasis(x, n, g.basisType)
	if math.Abs(basis) < 1e-9 {
		return
	}

	for i := int64(0); i < dst.dim; i++ {
		dst.data[i] += n.alpha.data[i] * basis
	}

	if n.hasChildren {
		for i := int64(0); i < g.inDim; i++ {
			if n.key.level[i] == 0 {
				continue
			}
			child, exists := g.getChildForDimAndArg(n, x, i, additionalNodes)
			if exists {
				g.evaluateRecursiveInto(dst, x, child, processedNodes, additionalNodes)
			}
		}
	}
}

func (g *AdaptiveSparseGrid) getChildForDimAndArg(parent node, x StaticPoint, dimension int64, additionalNodes map[gridKey]node) (node, bool) {
	key := parent.key
	key.level[dimension]++
	if x.data[dimension] < parent.centerUnit.data[dimension] {
		key.index[dimension] = 2*key.index[dimension] - 1
	} else {
		key.index[dimension] = 2*key.index[dimension] + 1
	}

	if child, exists := g.nodes[key]; exists {
		return child, true
	}
	if additionalNodes != nil {
		if child, exists := additionalNodes[key]; exists {
			return child, true
		}
	}
	return node{}, false
}
