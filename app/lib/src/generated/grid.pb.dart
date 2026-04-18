// This is a generated file - do not edit.
//
// Generated from grid.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'grid.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'grid.pbenum.dart';

class Point2D extends $pb.GeneratedMessage {
  factory Point2D({
    $core.double? x,
    $core.double? y,
  }) {
    final result = create();
    if (x != null) result.x = x;
    if (y != null) result.y = y;
    return result;
  }

  Point2D._();

  factory Point2D.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Point2D.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Point2D',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..aD(1, _omitFieldNames ? '' : 'x')
    ..aD(2, _omitFieldNames ? '' : 'y')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Point2D clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Point2D copyWith(void Function(Point2D) updates) =>
      super.copyWith((message) => updates(message as Point2D)) as Point2D;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Point2D create() => Point2D._();
  @$core.override
  Point2D createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Point2D getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Point2D>(create);
  static Point2D? _defaultInstance;

  @$pb.TagNumber(1)
  $core.double get x => $_getN(0);
  @$pb.TagNumber(1)
  set x($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(1)
  $core.bool hasX() => $_has(0);
  @$pb.TagNumber(1)
  void clearX() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get y => $_getN(1);
  @$pb.TagNumber(2)
  set y($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasY() => $_has(1);
  @$pb.TagNumber(2)
  void clearY() => $_clearField(2);
}

class Grid2DRequest extends $pb.GeneratedMessage {
  factory Grid2DRequest({
    $core.String? formulaX,
    $core.String? formulaY,
    Point2D? min,
    Point2D? max,
    $core.double? eps,
    BuildType? buildType,
    BasisType? basisType,
    $core.Iterable<Point2D>? anchorPoints,
    $core.double? step,
    FormulaType? formulaType,
    $core.int? maxLevel,
    $fixnum.Int64? maxNodesInDim,
  }) {
    final result = create();
    if (formulaX != null) result.formulaX = formulaX;
    if (formulaY != null) result.formulaY = formulaY;
    if (min != null) result.min = min;
    if (max != null) result.max = max;
    if (eps != null) result.eps = eps;
    if (buildType != null) result.buildType = buildType;
    if (basisType != null) result.basisType = basisType;
    if (anchorPoints != null) result.anchorPoints.addAll(anchorPoints);
    if (step != null) result.step = step;
    if (formulaType != null) result.formulaType = formulaType;
    if (maxLevel != null) result.maxLevel = maxLevel;
    if (maxNodesInDim != null) result.maxNodesInDim = maxNodesInDim;
    return result;
  }

  Grid2DRequest._();

  factory Grid2DRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Grid2DRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Grid2DRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'formulaX')
    ..aOS(2, _omitFieldNames ? '' : 'formulaY')
    ..aOM<Point2D>(3, _omitFieldNames ? '' : 'min', subBuilder: Point2D.create)
    ..aOM<Point2D>(4, _omitFieldNames ? '' : 'max', subBuilder: Point2D.create)
    ..aD(5, _omitFieldNames ? '' : 'eps')
    ..aE<BuildType>(6, _omitFieldNames ? '' : 'buildType',
        enumValues: BuildType.values)
    ..aE<BasisType>(7, _omitFieldNames ? '' : 'basisType',
        enumValues: BasisType.values)
    ..pPM<Point2D>(8, _omitFieldNames ? '' : 'anchorPoints',
        subBuilder: Point2D.create)
    ..aD(9, _omitFieldNames ? '' : 'step')
    ..aE<FormulaType>(10, _omitFieldNames ? '' : 'formulaType',
        enumValues: FormulaType.values)
    ..aI(11, _omitFieldNames ? '' : 'maxLevel')
    ..aInt64(12, _omitFieldNames ? '' : 'maxNodesInDim')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2DRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2DRequest copyWith(void Function(Grid2DRequest) updates) =>
      super.copyWith((message) => updates(message as Grid2DRequest))
          as Grid2DRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Grid2DRequest create() => Grid2DRequest._();
  @$core.override
  Grid2DRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Grid2DRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Grid2DRequest>(create);
  static Grid2DRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get formulaX => $_getSZ(0);
  @$pb.TagNumber(1)
  set formulaX($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasFormulaX() => $_has(0);
  @$pb.TagNumber(1)
  void clearFormulaX() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get formulaY => $_getSZ(1);
  @$pb.TagNumber(2)
  set formulaY($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasFormulaY() => $_has(1);
  @$pb.TagNumber(2)
  void clearFormulaY() => $_clearField(2);

  @$pb.TagNumber(3)
  Point2D get min => $_getN(2);
  @$pb.TagNumber(3)
  set min(Point2D value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasMin() => $_has(2);
  @$pb.TagNumber(3)
  void clearMin() => $_clearField(3);
  @$pb.TagNumber(3)
  Point2D ensureMin() => $_ensure(2);

  @$pb.TagNumber(4)
  Point2D get max => $_getN(3);
  @$pb.TagNumber(4)
  set max(Point2D value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasMax() => $_has(3);
  @$pb.TagNumber(4)
  void clearMax() => $_clearField(4);
  @$pb.TagNumber(4)
  Point2D ensureMax() => $_ensure(3);

  @$pb.TagNumber(5)
  $core.double get eps => $_getN(4);
  @$pb.TagNumber(5)
  set eps($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasEps() => $_has(4);
  @$pb.TagNumber(5)
  void clearEps() => $_clearField(5);

  @$pb.TagNumber(6)
  BuildType get buildType => $_getN(5);
  @$pb.TagNumber(6)
  set buildType(BuildType value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasBuildType() => $_has(5);
  @$pb.TagNumber(6)
  void clearBuildType() => $_clearField(6);

  @$pb.TagNumber(7)
  BasisType get basisType => $_getN(6);
  @$pb.TagNumber(7)
  set basisType(BasisType value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasBasisType() => $_has(6);
  @$pb.TagNumber(7)
  void clearBasisType() => $_clearField(7);

  @$pb.TagNumber(8)
  $pb.PbList<Point2D> get anchorPoints => $_getList(7);

  @$pb.TagNumber(9)
  $core.double get step => $_getN(8);
  @$pb.TagNumber(9)
  set step($core.double value) => $_setDouble(8, value);
  @$pb.TagNumber(9)
  $core.bool hasStep() => $_has(8);
  @$pb.TagNumber(9)
  void clearStep() => $_clearField(9);

  @$pb.TagNumber(10)
  FormulaType get formulaType => $_getN(9);
  @$pb.TagNumber(10)
  set formulaType(FormulaType value) => $_setField(10, value);
  @$pb.TagNumber(10)
  $core.bool hasFormulaType() => $_has(9);
  @$pb.TagNumber(10)
  void clearFormulaType() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.int get maxLevel => $_getIZ(10);
  @$pb.TagNumber(11)
  set maxLevel($core.int value) => $_setSignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasMaxLevel() => $_has(10);
  @$pb.TagNumber(11)
  void clearMaxLevel() => $_clearField(11);

  @$pb.TagNumber(12)
  $fixnum.Int64 get maxNodesInDim => $_getI64(11);
  @$pb.TagNumber(12)
  set maxNodesInDim($fixnum.Int64 value) => $_setInt64(11, value);
  @$pb.TagNumber(12)
  $core.bool hasMaxNodesInDim() => $_has(11);
  @$pb.TagNumber(12)
  void clearMaxNodesInDim() => $_clearField(12);
}

class Grid2D_Index2D extends $pb.GeneratedMessage {
  factory Grid2D_Index2D({
    $fixnum.Int64? x,
    $fixnum.Int64? y,
  }) {
    final result = create();
    if (x != null) result.x = x;
    if (y != null) result.y = y;
    return result;
  }

  Grid2D_Index2D._();

  factory Grid2D_Index2D.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Grid2D_Index2D.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Grid2D.Index2D',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'x', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'y', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_Index2D clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_Index2D copyWith(void Function(Grid2D_Index2D) updates) =>
      super.copyWith((message) => updates(message as Grid2D_Index2D))
          as Grid2D_Index2D;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Grid2D_Index2D create() => Grid2D_Index2D._();
  @$core.override
  Grid2D_Index2D createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Grid2D_Index2D getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Grid2D_Index2D>(create);
  static Grid2D_Index2D? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get x => $_getI64(0);
  @$pb.TagNumber(1)
  set x($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasX() => $_has(0);
  @$pb.TagNumber(1)
  void clearX() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get y => $_getI64(1);
  @$pb.TagNumber(2)
  set y($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasY() => $_has(1);
  @$pb.TagNumber(2)
  void clearY() => $_clearField(2);
}

class Grid2D_GridKey2D extends $pb.GeneratedMessage {
  factory Grid2D_GridKey2D({
    Grid2D_Index2D? index,
    Grid2D_Index2D? level,
  }) {
    final result = create();
    if (index != null) result.index = index;
    if (level != null) result.level = level;
    return result;
  }

  Grid2D_GridKey2D._();

  factory Grid2D_GridKey2D.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Grid2D_GridKey2D.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Grid2D.GridKey2D',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..aOM<Grid2D_Index2D>(1, _omitFieldNames ? '' : 'index',
        subBuilder: Grid2D_Index2D.create)
    ..aOM<Grid2D_Index2D>(2, _omitFieldNames ? '' : 'level',
        subBuilder: Grid2D_Index2D.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_GridKey2D clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_GridKey2D copyWith(void Function(Grid2D_GridKey2D) updates) =>
      super.copyWith((message) => updates(message as Grid2D_GridKey2D))
          as Grid2D_GridKey2D;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Grid2D_GridKey2D create() => Grid2D_GridKey2D._();
  @$core.override
  Grid2D_GridKey2D createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Grid2D_GridKey2D getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Grid2D_GridKey2D>(create);
  static Grid2D_GridKey2D? _defaultInstance;

  @$pb.TagNumber(1)
  Grid2D_Index2D get index => $_getN(0);
  @$pb.TagNumber(1)
  set index(Grid2D_Index2D value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasIndex() => $_has(0);
  @$pb.TagNumber(1)
  void clearIndex() => $_clearField(1);
  @$pb.TagNumber(1)
  Grid2D_Index2D ensureIndex() => $_ensure(0);

  @$pb.TagNumber(2)
  Grid2D_Index2D get level => $_getN(1);
  @$pb.TagNumber(2)
  set level(Grid2D_Index2D value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasLevel() => $_has(1);
  @$pb.TagNumber(2)
  void clearLevel() => $_clearField(2);
  @$pb.TagNumber(2)
  Grid2D_Index2D ensureLevel() => $_ensure(1);
}

class Grid2D_Node2D extends $pb.GeneratedMessage {
  factory Grid2D_Node2D({
    Grid2D_GridKey2D? key,
    Point2D? centerUnit,
    Point2D? alpha,
    $core.bool? hasChildren,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (centerUnit != null) result.centerUnit = centerUnit;
    if (alpha != null) result.alpha = alpha;
    if (hasChildren != null) result.hasChildren = hasChildren;
    return result;
  }

  Grid2D_Node2D._();

  factory Grid2D_Node2D.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Grid2D_Node2D.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Grid2D.Node2D',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..aOM<Grid2D_GridKey2D>(1, _omitFieldNames ? '' : 'key',
        subBuilder: Grid2D_GridKey2D.create)
    ..aOM<Point2D>(2, _omitFieldNames ? '' : 'centerUnit',
        subBuilder: Point2D.create)
    ..aOM<Point2D>(3, _omitFieldNames ? '' : 'alpha',
        subBuilder: Point2D.create)
    ..aOB(4, _omitFieldNames ? '' : 'hasChildren')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_Node2D clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D_Node2D copyWith(void Function(Grid2D_Node2D) updates) =>
      super.copyWith((message) => updates(message as Grid2D_Node2D))
          as Grid2D_Node2D;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Grid2D_Node2D create() => Grid2D_Node2D._();
  @$core.override
  Grid2D_Node2D createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Grid2D_Node2D getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Grid2D_Node2D>(create);
  static Grid2D_Node2D? _defaultInstance;

  @$pb.TagNumber(1)
  Grid2D_GridKey2D get key => $_getN(0);
  @$pb.TagNumber(1)
  set key(Grid2D_GridKey2D value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);
  @$pb.TagNumber(1)
  Grid2D_GridKey2D ensureKey() => $_ensure(0);

  @$pb.TagNumber(2)
  Point2D get centerUnit => $_getN(1);
  @$pb.TagNumber(2)
  set centerUnit(Point2D value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasCenterUnit() => $_has(1);
  @$pb.TagNumber(2)
  void clearCenterUnit() => $_clearField(2);
  @$pb.TagNumber(2)
  Point2D ensureCenterUnit() => $_ensure(1);

  @$pb.TagNumber(3)
  Point2D get alpha => $_getN(2);
  @$pb.TagNumber(3)
  set alpha(Point2D value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasAlpha() => $_has(2);
  @$pb.TagNumber(3)
  void clearAlpha() => $_clearField(3);
  @$pb.TagNumber(3)
  Point2D ensureAlpha() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.bool get hasChildren => $_getBF(3);
  @$pb.TagNumber(4)
  set hasChildren($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasHasChildren() => $_has(3);
  @$pb.TagNumber(4)
  void clearHasChildren() => $_clearField(4);
}

class Grid2D extends $pb.GeneratedMessage {
  factory Grid2D({
    $core.Iterable<Grid2D_Node2D>? nodes,
    $core.Iterable<Grid2D_Node2D>? entryPoints,
    Point2D? min,
    Point2D? max,
    BasisType? basisType,
  }) {
    final result = create();
    if (nodes != null) result.nodes.addAll(nodes);
    if (entryPoints != null) result.entryPoints.addAll(entryPoints);
    if (min != null) result.min = min;
    if (max != null) result.max = max;
    if (basisType != null) result.basisType = basisType;
    return result;
  }

  Grid2D._();

  factory Grid2D.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Grid2D.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Grid2D',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'grid'),
      createEmptyInstance: create)
    ..pPM<Grid2D_Node2D>(1, _omitFieldNames ? '' : 'nodes',
        subBuilder: Grid2D_Node2D.create)
    ..pPM<Grid2D_Node2D>(2, _omitFieldNames ? '' : 'entryPoints',
        subBuilder: Grid2D_Node2D.create)
    ..aOM<Point2D>(3, _omitFieldNames ? '' : 'min', subBuilder: Point2D.create)
    ..aOM<Point2D>(4, _omitFieldNames ? '' : 'max', subBuilder: Point2D.create)
    ..aE<BasisType>(5, _omitFieldNames ? '' : 'basisType',
        enumValues: BasisType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Grid2D copyWith(void Function(Grid2D) updates) =>
      super.copyWith((message) => updates(message as Grid2D)) as Grid2D;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Grid2D create() => Grid2D._();
  @$core.override
  Grid2D createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Grid2D getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Grid2D>(create);
  static Grid2D? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Grid2D_Node2D> get nodes => $_getList(0);

  @$pb.TagNumber(2)
  $pb.PbList<Grid2D_Node2D> get entryPoints => $_getList(1);

  @$pb.TagNumber(3)
  Point2D get min => $_getN(2);
  @$pb.TagNumber(3)
  set min(Point2D value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasMin() => $_has(2);
  @$pb.TagNumber(3)
  void clearMin() => $_clearField(3);
  @$pb.TagNumber(3)
  Point2D ensureMin() => $_ensure(2);

  @$pb.TagNumber(4)
  Point2D get max => $_getN(3);
  @$pb.TagNumber(4)
  set max(Point2D value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasMax() => $_has(3);
  @$pb.TagNumber(4)
  void clearMax() => $_clearField(4);
  @$pb.TagNumber(4)
  Point2D ensureMax() => $_ensure(3);

  @$pb.TagNumber(5)
  BasisType get basisType => $_getN(4);
  @$pb.TagNumber(5)
  set basisType(BasisType value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasBasisType() => $_has(4);
  @$pb.TagNumber(5)
  void clearBasisType() => $_clearField(5);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
