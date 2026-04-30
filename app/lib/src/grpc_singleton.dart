import 'package:grpc/grpc_or_grpcweb.dart';
import 'package:app/src/generated/grid.pbgrpc.dart';

class ComputeService {
  final GrpcOrGrpcWebClientChannel channel;

  late final GridServiceClient gridClient;

  static ComputeService? _service;

  ComputeService._internal() : channel = GrpcOrGrpcWebClientChannel.toSeparateEndpoints(
      grpcHost: "diplom.funandchecks.ru/api",
      grpcPort: 443,
      grpcTransportSecure: true,
      grpcWebHost: "diplom.funandchecks.ru/envoy",
      grpcWebPort: 443,
      grpcWebTransportSecure: true) {
    gridClient = GridServiceClient(
      channel,
    );
  }

  factory ComputeService() {
    _service ??= ComputeService._internal();
    return _service!;
  }

  void shutdown() async => await channel.shutdown();
}