// This is a generated file - do not edit.
//
// Generated from grid.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'grid.pb.dart' as $0;

export 'grid.pb.dart';

@$pb.GrpcServiceName('grid.GridService')
class GridServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  GridServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Grid2D> getGrid2D(
    $0.Grid2DRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGrid2D, request, options: options);
  }

  // method descriptors

  static final _$getGrid2D = $grpc.ClientMethod<$0.Grid2DRequest, $0.Grid2D>(
      '/grid.GridService/GetGrid2D',
      ($0.Grid2DRequest value) => value.writeToBuffer(),
      $0.Grid2D.fromBuffer);
}

@$pb.GrpcServiceName('grid.GridService')
abstract class GridServiceBase extends $grpc.Service {
  $core.String get $name => 'grid.GridService';

  GridServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.Grid2DRequest, $0.Grid2D>(
        'GetGrid2D',
        getGrid2D_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Grid2DRequest.fromBuffer(value),
        ($0.Grid2D value) => value.writeToBuffer()));
  }

  $async.Future<$0.Grid2D> getGrid2D_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Grid2DRequest> $request) async {
    return getGrid2D($call, await $request);
  }

  $async.Future<$0.Grid2D> getGrid2D(
      $grpc.ServiceCall call, $0.Grid2DRequest request);
}
