// This is a generated file - do not edit.
//
// Generated from grid.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use basisTypeDescriptor instead')
const BasisType$json = {
  '1': 'BasisType',
  '2': [
    {'1': 'LINEAR', '2': 0},
    {'1': 'QUADRATIC', '2': 1},
  ],
};

/// Descriptor for `BasisType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List basisTypeDescriptor = $convert
    .base64Decode('CglCYXNpc1R5cGUSCgoGTElORUFSEAASDQoJUVVBRFJBVElDEAE=');

@$core.Deprecated('Use buildTypeDescriptor instead')
const BuildType$json = {
  '1': 'BuildType',
  '2': [
    {'1': 'SEQUENTIAL', '2': 0},
    {'1': 'PARALLEL', '2': 1},
  ],
};

/// Descriptor for `BuildType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List buildTypeDescriptor = $convert
    .base64Decode('CglCdWlsZFR5cGUSDgoKU0VRVUVOVElBTBAAEgwKCFBBUkFMTEVMEAE=');

@$core.Deprecated('Use formulaTypeDescriptor instead')
const FormulaType$json = {
  '1': 'FormulaType',
  '2': [
    {'1': 'SIMPLE', '2': 0},
    {'1': 'DIFFUR', '2': 1},
  ],
};

/// Descriptor for `FormulaType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List formulaTypeDescriptor = $convert
    .base64Decode('CgtGb3JtdWxhVHlwZRIKCgZTSU1QTEUQABIKCgZESUZGVVIQAQ==');

@$core.Deprecated('Use point2DDescriptor instead')
const Point2D$json = {
  '1': 'Point2D',
  '2': [
    {'1': 'x', '3': 1, '4': 1, '5': 1, '10': 'x'},
    {'1': 'y', '3': 2, '4': 1, '5': 1, '10': 'y'},
  ],
};

/// Descriptor for `Point2D`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List point2DDescriptor = $convert
    .base64Decode('CgdQb2ludDJEEgwKAXgYASABKAFSAXgSDAoBeRgCIAEoAVIBeQ==');

@$core.Deprecated('Use grid2DRequestDescriptor instead')
const Grid2DRequest$json = {
  '1': 'Grid2DRequest',
  '2': [
    {'1': 'formula_x', '3': 1, '4': 1, '5': 9, '10': 'formulaX'},
    {'1': 'formula_y', '3': 2, '4': 1, '5': 9, '10': 'formulaY'},
    {'1': 'min', '3': 3, '4': 1, '5': 11, '6': '.grid.Point2D', '10': 'min'},
    {'1': 'max', '3': 4, '4': 1, '5': 11, '6': '.grid.Point2D', '10': 'max'},
    {'1': 'eps', '3': 5, '4': 1, '5': 1, '10': 'eps'},
    {
      '1': 'build_type',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.grid.BuildType',
      '10': 'buildType'
    },
    {
      '1': 'basis_type',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.grid.BasisType',
      '10': 'basisType'
    },
    {
      '1': 'anchor_points',
      '3': 8,
      '4': 3,
      '5': 11,
      '6': '.grid.Point2D',
      '10': 'anchorPoints'
    },
    {'1': 'step', '3': 9, '4': 1, '5': 13, '9': 0, '10': 'step', '17': true},
    {
      '1': 'formula_type',
      '3': 10,
      '4': 1,
      '5': 14,
      '6': '.grid.FormulaType',
      '10': 'formulaType'
    },
  ],
  '8': [
    {'1': '_step'},
  ],
};

/// Descriptor for `Grid2DRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grid2DRequestDescriptor = $convert.base64Decode(
    'Cg1HcmlkMkRSZXF1ZXN0EhsKCWZvcm11bGFfeBgBIAEoCVIIZm9ybXVsYVgSGwoJZm9ybXVsYV'
    '95GAIgASgJUghmb3JtdWxhWRIfCgNtaW4YAyABKAsyDS5ncmlkLlBvaW50MkRSA21pbhIfCgNt'
    'YXgYBCABKAsyDS5ncmlkLlBvaW50MkRSA21heBIQCgNlcHMYBSABKAFSA2VwcxIuCgpidWlsZF'
    '90eXBlGAYgASgOMg8uZ3JpZC5CdWlsZFR5cGVSCWJ1aWxkVHlwZRIuCgpiYXNpc190eXBlGAcg'
    'ASgOMg8uZ3JpZC5CYXNpc1R5cGVSCWJhc2lzVHlwZRIyCg1hbmNob3JfcG9pbnRzGAggAygLMg'
    '0uZ3JpZC5Qb2ludDJEUgxhbmNob3JQb2ludHMSFwoEc3RlcBgJIAEoDUgAUgRzdGVwiAEBEjQK'
    'DGZvcm11bGFfdHlwZRgKIAEoDjIRLmdyaWQuRm9ybXVsYVR5cGVSC2Zvcm11bGFUeXBlQgcKBV'
    '9zdGVw');

@$core.Deprecated('Use grid2DDescriptor instead')
const Grid2D$json = {
  '1': 'Grid2D',
  '2': [
    {
      '1': 'nodes',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.grid.Grid2D.Node2D',
      '10': 'nodes'
    },
    {
      '1': 'entry_points',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.grid.Grid2D.Node2D',
      '10': 'entryPoints'
    },
    {'1': 'min', '3': 3, '4': 1, '5': 11, '6': '.grid.Point2D', '10': 'min'},
    {'1': 'max', '3': 4, '4': 1, '5': 11, '6': '.grid.Point2D', '10': 'max'},
    {
      '1': 'basis_type',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.grid.BasisType',
      '10': 'basisType'
    },
  ],
  '3': [Grid2D_Index2D$json, Grid2D_GridKey2D$json, Grid2D_Node2D$json],
};

@$core.Deprecated('Use grid2DDescriptor instead')
const Grid2D_Index2D$json = {
  '1': 'Index2D',
  '2': [
    {'1': 'x', '3': 1, '4': 1, '5': 4, '10': 'x'},
    {'1': 'y', '3': 2, '4': 1, '5': 4, '10': 'y'},
  ],
};

@$core.Deprecated('Use grid2DDescriptor instead')
const Grid2D_GridKey2D$json = {
  '1': 'GridKey2D',
  '2': [
    {
      '1': 'index',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.grid.Grid2D.Index2D',
      '10': 'index'
    },
    {
      '1': 'level',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.grid.Grid2D.Index2D',
      '10': 'level'
    },
  ],
};

@$core.Deprecated('Use grid2DDescriptor instead')
const Grid2D_Node2D$json = {
  '1': 'Node2D',
  '2': [
    {
      '1': 'key',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.grid.Grid2D.GridKey2D',
      '10': 'key'
    },
    {
      '1': 'center_unit',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.grid.Point2D',
      '10': 'centerUnit'
    },
    {
      '1': 'alpha',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.grid.Point2D',
      '10': 'alpha'
    },
    {'1': 'has_children', '3': 4, '4': 1, '5': 8, '10': 'hasChildren'},
  ],
};

/// Descriptor for `Grid2D`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grid2DDescriptor = $convert.base64Decode(
    'CgZHcmlkMkQSKQoFbm9kZXMYASADKAsyEy5ncmlkLkdyaWQyRC5Ob2RlMkRSBW5vZGVzEjYKDG'
    'VudHJ5X3BvaW50cxgCIAMoCzITLmdyaWQuR3JpZDJELk5vZGUyRFILZW50cnlQb2ludHMSHwoD'
    'bWluGAMgASgLMg0uZ3JpZC5Qb2ludDJEUgNtaW4SHwoDbWF4GAQgASgLMg0uZ3JpZC5Qb2ludD'
    'JEUgNtYXgSLgoKYmFzaXNfdHlwZRgFIAEoDjIPLmdyaWQuQmFzaXNUeXBlUgliYXNpc1R5cGUa'
    'JQoHSW5kZXgyRBIMCgF4GAEgASgEUgF4EgwKAXkYAiABKARSAXkaYwoJR3JpZEtleTJEEioKBW'
    'luZGV4GAEgASgLMhQuZ3JpZC5HcmlkMkQuSW5kZXgyRFIFaW5kZXgSKgoFbGV2ZWwYAiABKAsy'
    'FC5ncmlkLkdyaWQyRC5JbmRleDJEUgVsZXZlbBqqAQoGTm9kZTJEEigKA2tleRgBIAEoCzIWLm'
    'dyaWQuR3JpZDJELkdyaWRLZXkyRFIDa2V5Ei4KC2NlbnRlcl91bml0GAIgASgLMg0uZ3JpZC5Q'
    'b2ludDJEUgpjZW50ZXJVbml0EiMKBWFscGhhGAMgASgLMg0uZ3JpZC5Qb2ludDJEUgVhbHBoYR'
    'IhCgxoYXNfY2hpbGRyZW4YBCABKAhSC2hhc0NoaWxkcmVu');
