import 'package:app/src/generated/grid.pb.dart' as grid;
import 'package:fixnum/fixnum.dart';
import 'package:dart_numerics/dart_numerics.dart' as num;
import 'dart:collection';

enum BasisType {
  linear,
  quadratic
}

enum BuildType {
  parallel,
  sequential
}

class Point2D {
  double x;
  double y;

  Point2D(this.x, this.y);
  Point2D.fromPB(grid.Point2D point) : x = point.x, y = point.y;
}

class Index2D {
  Int64 x;
  Int64 y;

  Index2D(this.x, this.y);
  Index2D.copyWith(Index2D other) : x = other.x, y = other.y;
  Index2D.fromPB(grid.Grid2D_Index2D index) : x = index.x, y = index.y;
}

class GridKey2D {
  Index2D index;
  Index2D level;

  GridKey2D(this.index, this.level);
  GridKey2D.fromPB(grid.Grid2D_GridKey2D key) : index = Index2D.fromPB(key.index), level = Index2D.fromPB(key.level);
}

bool equal(GridKey2D a, GridKey2D b) {
  return a.index.x == b.index.x && a.index.y == b.index.y && a.level.x == b.level.x && a.level.y == b.level.y;
}

int hash(GridKey2D key) {
  return key.index.x.hashCode ^ (key.index.y.hashCode << 2) ^ (key.level.x.hashCode << 3) ^ (key.level.y.hashCode << 4);
}

class Node2D {
  GridKey2D key;
  Point2D centerUnit;
  Point2D alpha;
  bool hasChildren;

  Node2D(this.key, this.centerUnit, this.alpha, this.hasChildren);
  Node2D.fromPB(grid.Grid2D_Node2D node) : key = GridKey2D.fromPB(node.key), centerUnit = Point2D.fromPB(node.centerUnit), alpha = Point2D.fromPB(node.alpha), hasChildren = node.hasChildren;
}


class Grid2D {
  Point2D min;
  Point2D max;
  BasisType basisType;
  List<Node2D> entryPoints = [];
  HashMap<GridKey2D, Node2D> nodes = HashMap(equals: equal, hashCode: hash);

  Grid2D(this.min, this.max, this.basisType, this.entryPoints, this.nodes);
  Grid2D.fromPB(grid.Grid2D grid2D) : min = Point2D.fromPB(grid2D.min), max = Point2D.fromPB(grid2D.max), basisType = grid2D.basisType == grid.BasisType.LINEAR ? BasisType.linear : BasisType.quadratic
  {
    grid2D.nodes.forEach((node) => nodes[GridKey2D.fromPB(node.key)] = Node2D.fromPB(node));
    grid2D.entryPoints.forEach((entryPoint) => entryPoints.add(Node2D.fromPB(entryPoint)));
  }

  Point2D evaluate(Point2D point) {
    _checkEvaluationPoint(point);
    var arg = _toUnit(point);
    Point2D answer = Point2D(0, 0);

    for (var entryPoint in entryPoints) {
      HashSet<GridKey2D> processedNodes = HashSet();
      var additional = _evaluateRecursive(arg, entryPoint, processedNodes);
      answer = Point2D(answer.x + additional.x, answer.y + additional.y);
    }

    return answer;
  }

  void _checkEvaluationPoint(Point2D point)
  {
    if (point.x < min.x || point.x > max.x || point.y < min.y || point.y > max.y) {
      throw Exception('Evaluation point is out of bounds');
    }
  }

  Point2D _toUnit(Point2D point) {
    return Point2D((point.x - min.x) / (max.x - min.x), (point.y - min.y) / (max.y - min.y));
  }

  Point2D _evaluateRecursive(Point2D arg, Node2D node, HashSet<GridKey2D> processedNodes) {
    if (processedNodes.contains(node.key)) {
      return Point2D(0, 0);
    }

    processedNodes.add(node.key);

    var basis = _basis(arg, node.key, basisType);

    if (basis.abs() < num.epsilon) {
      return Point2D(0, 0);
    }

    var value = Point2D(node.alpha.x * basis, node.alpha.y * basis);

    var children = _getChildrenForArg(arg, node);
    for (var child in children) {
      var additional = _evaluateRecursive(arg, child, processedNodes);
      value = Point2D(value.x + additional.x, value.y + additional.y);
    }

    return value;
  }

  List<Node2D> _getChildrenForArg(Point2D arg, Node2D node) {
    List<Node2D> children = [];

    GridKey2D xKey = GridKey2D(Index2D.copyWith(node.key.index), Index2D(node.key.level.x + 1, node.key.level.y));
    GridKey2D yKey = GridKey2D(Index2D.copyWith(node.key.index), Index2D(node.key.level.x, node.key.level.y + 1));

    xKey.index.x = Int64(2) * xKey.index.x + (arg.x < node.centerUnit.x ? Int64(-1) : Int64(1));
    yKey.index.y = Int64(2) * yKey.index.y + (arg.y < node.centerUnit.y ? Int64(-1) : Int64(1));

    if (nodes.containsKey(xKey)) {
      children.add(nodes[xKey]!);
    }

    if (nodes.containsKey(yKey)) {
      children.add(nodes[yKey]!);
    }

    return children;
  }

  double _basis(Point2D arg, GridKey2D key, BasisType basisType) {
    return _basis1D(arg.x, key.level.x, key.index.x, basisType) * _basis1D(arg.y, key.level.y, key.index.y, basisType);
  }

  double _basis1D(double arg, Int64 level, Int64 index, BasisType basisType) {
    if (level == 0)
    {
      if (index == 0) {
        return 1.0 - arg;
      }
      if (index == 1) {
        return arg;
      }
      return 0.0;
    }

    double h = 1.0 / (1 << level.toInt());
    double center = (index.toDouble()) * h;
    double dist = (arg - center).abs();

    if (dist >= h) {
      return 0.0;
    }

    double t = dist / h;

    if (basisType == BasisType.linear) {
      return 1.0 - t;
    } else {
      return 1.0 - t * t;
    }
  }
}