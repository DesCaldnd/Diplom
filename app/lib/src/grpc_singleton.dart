import 'package:grpc/grpc_or_grpcweb.dart';
import 'package:app/src/generated/grid.pbgrpc.dart';
import 'package:app/src/generated/grid.pb.dart';

class ComputeService {
  final GrpcOrGrpcWebClientChannel channel;

  late final GridServiceClient gridClient;

  static ComputeService? _service;

  ComputeService._internal() : channel = GrpcOrGrpcWebClientChannel.toSeparatePorts(
      host: "192.168.1.8",
      grpcPort: 9999,
      grpcTransportSecure: false,
      grpcWebPort: 8082,
      grpcWebTransportSecure: false) {
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