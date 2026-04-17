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

import 'package:protobuf/protobuf.dart' as $pb;

class BasisType extends $pb.ProtobufEnum {
  static const BasisType LINEAR =
      BasisType._(0, _omitEnumNames ? '' : 'LINEAR');
  static const BasisType QUADRATIC =
      BasisType._(1, _omitEnumNames ? '' : 'QUADRATIC');

  static const $core.List<BasisType> values = <BasisType>[
    LINEAR,
    QUADRATIC,
  ];

  static final $core.List<BasisType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BasisType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BasisType._(super.value, super.name);
}

class BuildType extends $pb.ProtobufEnum {
  static const BuildType SEQUENTIAL =
      BuildType._(0, _omitEnumNames ? '' : 'SEQUENTIAL');
  static const BuildType PARALLEL =
      BuildType._(1, _omitEnumNames ? '' : 'PARALLEL');

  static const $core.List<BuildType> values = <BuildType>[
    SEQUENTIAL,
    PARALLEL,
  ];

  static final $core.List<BuildType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BuildType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BuildType._(super.value, super.name);
}

class FormulaType extends $pb.ProtobufEnum {
  static const FormulaType SIMPLE =
      FormulaType._(0, _omitEnumNames ? '' : 'SIMPLE');
  static const FormulaType DIFFUR =
      FormulaType._(1, _omitEnumNames ? '' : 'DIFFUR');

  static const $core.List<FormulaType> values = <FormulaType>[
    SIMPLE,
    DIFFUR,
  ];

  static final $core.List<FormulaType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FormulaType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FormulaType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
