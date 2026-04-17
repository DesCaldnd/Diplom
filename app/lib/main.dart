import 'dart:async';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:grpc/grpc.dart';

import 'package:app/src/grid.dart';
import 'package:app/src/grpc_singleton.dart';
import 'package:app/src/generated/grid.pbgrpc.dart' as pb;

void main() {
  runApp(const GridApp());
}

class GridApp extends StatelessWidget {
  const GridApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Grid2D Evaluator',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: const GridScreen(),
    );
  }
}

class GridScreen extends StatefulWidget {
  const GridScreen({super.key});

  @override
  State<GridScreen> createState() => _GridScreenState();
}

class _GridScreenState extends State<GridScreen> {
  final TextEditingController _formXCtrl = TextEditingController();
  final TextEditingController _formYCtrl = TextEditingController();
  final TextEditingController _minXCtrl = TextEditingController();
  final TextEditingController _maxXCtrl = TextEditingController();
  final TextEditingController _minYCtrl = TextEditingController();
  final TextEditingController _maxYCtrl = TextEditingController();
  final TextEditingController _epsCtrl = TextEditingController();
  final TextEditingController _stepCtrl = TextEditingController(text: '1');

  pb.BasisType _basisType = pb.BasisType.QUADRATIC;
  pb.FormulaType _formulaType = pb.FormulaType.DIFFUR;
  pb.BuildType _buildType = pb.BuildType.PARALLEL;

  bool _isLoading = false;
  Grid2D? _resultGrid;
  double _origMinX = 0;
  double _origMaxX = 0;
  double _origMinY = 0;
  double _origMaxY = 0;

  @override
  void dispose() {
    _formXCtrl.dispose();
    _formYCtrl.dispose();
    _minXCtrl.dispose();
    _maxXCtrl.dispose();
    _minYCtrl.dispose();
    _maxYCtrl.dispose();
    _epsCtrl.dispose();
    _stepCtrl.dispose();
    super.dispose();
  }

  Future<void> _calculate() async {
    if (_formXCtrl.text.isEmpty ||
        _formYCtrl.text.isEmpty ||
        _minXCtrl.text.isEmpty ||
        _maxXCtrl.text.isEmpty ||
        _minYCtrl.text.isEmpty ||
        _maxYCtrl.text.isEmpty ||
        _epsCtrl.text.isEmpty ||
        _stepCtrl.text.isEmpty) {
      _showError("All fields must be filled.");
      return;
    }

    final minX = double.tryParse(_minXCtrl.text);
    final maxX = double.tryParse(_maxXCtrl.text);
    final minY = double.tryParse(_minYCtrl.text);
    final maxY = double.tryParse(_maxYCtrl.text);
    final eps = double.tryParse(_epsCtrl.text);
    final step = int.tryParse(_stepCtrl.text);

    if (minX == null ||
        maxX == null ||
        minY == null ||
        maxY == null ||
        eps == null ||
        step == null) {
      _showError("Invalid numeric values.");
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      final request = pb.Grid2DRequest(
        formulaX: _formXCtrl.text,
        formulaY: _formYCtrl.text,
        min: pb.Point2D(x: minX, y: minY),
        max: pb.Point2D(x: maxX, y: maxY),
        eps: eps,
        buildType: _buildType,
        basisType: _basisType,
        step: step,
        formulaType: _formulaType,
      );

      final response = await ComputeService().gridClient.getGrid2D(
        request,
        options: CallOptions(timeout: const Duration(seconds: 20)),
      );

      setState(() {
        _resultGrid = Grid2D.fromPB(response);
        _origMinX = minX;
        _origMaxX = maxX;
        _origMinY = minY;
        _origMaxY = maxY;
      });
    } on GrpcError catch (e) {
      _showError("gRPC Error: ${e.message}");
    } on TimeoutException {
      _showError("Request timed out (20s).");
    } catch (e) {
      _showError("An unexpected error occurred: $e");
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
  }

  Widget _buildTextField(TextEditingController controller, String label,
      {String? hint, bool isNumeric = true}) {
    return Expanded(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 4.0, vertical: 8.0),
        child: TextField(
          controller: controller,
          decoration: InputDecoration(
            labelText: label,
            hintText: hint,
            border: const OutlineInputBorder(),
            isDense: true,
          ),
          keyboardType: isNumeric ? TextInputType.number : TextInputType.text,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Grid2D Evaluator')),
      body: Stack(
        children: [
          Column(
            children: [
              Expanded(
                flex: 5,
                child: Container(
                  width: double.infinity,
                  margin: const EdgeInsets.all(8.0),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.grey),
                    color: Colors.white,
                  ),
                  child: ClipRect(
                    child: CustomPaint(
                      painter: _resultGrid != null
                          ? GridPainter(
                        grid: _resultGrid!,
                        minX: _origMinX,
                        maxX: _origMaxX,
                        minY: _origMinY,
                        maxY: _origMaxY,
                      )
                          : null,
                    ),
                  ),
                ),
              ),
              Expanded(
                flex: 3,
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(8.0),
                  child: Column(
                    children: [
                      Row(
                        children: [
                          _buildTextField(_formXCtrl, 'Formula X',
                              hint: _formulaType == pb.FormulaType.DIFFUR ? 'dx(t)/dt=f1(x, y, t)' : 'x=f1(x, y)', isNumeric: false),
                          _buildTextField(_formYCtrl, 'Formula Y',
                              hint: _formulaType == pb.FormulaType.DIFFUR ? 'dy(t)/dt=f2(x, y, t)' : 'y=f2(x, y)', isNumeric: false),
                        ],
                      ),
                      Row(
                        children: [
                          _buildTextField(_minXCtrl, 'Min X'),
                          _buildTextField(_maxXCtrl, 'Max X'),
                          _buildTextField(_minYCtrl, 'Min Y'),
                          _buildTextField(_maxYCtrl, 'Max Y'),
                        ],
                      ),
                      Row(
                        children: [
                          _buildTextField(_epsCtrl, 'Epsilon'),
                          _buildTextField(_stepCtrl, 'Step'),
                        ],
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          DropdownButton<pb.BasisType>(
                            value: _basisType,
                            items: const [
                              DropdownMenuItem(
                                  value: pb.BasisType.LINEAR,
                                  child: Text('Linear')),
                              DropdownMenuItem(
                                  value: pb.BasisType.QUADRATIC,
                                  child: Text('Quadratic')),
                            ],
                            onChanged: (val) {
                              if (val != null) setState(() => _basisType = val);
                            },
                          ),
                          DropdownButton<pb.FormulaType>(
                            value: _formulaType,
                            items: const [
                              DropdownMenuItem(
                                  value: pb.FormulaType.SIMPLE,
                                  child: Text('Simple')),
                              DropdownMenuItem(
                                  value: pb.FormulaType.DIFFUR,
                                  child: Text('Differential')),
                            ],
                            onChanged: (val) {
                              if (val != null) setState(() => _formulaType = val);
                            },
                          ),
                          DropdownButton<pb.BuildType>(
                            value: _buildType,
                            items: const [
                              DropdownMenuItem(
                                  value: pb.BuildType.SEQUENTIAL,
                                  child: Text('Sequential')),
                              DropdownMenuItem(
                                  value: pb.BuildType.PARALLEL,
                                  child: Text('Parallel')),
                            ],
                            onChanged: (val) {
                              if (val != null) setState(() => _buildType = val);
                            },
                          ),
                          ElevatedButton(
                            onPressed: _isLoading ? null : _calculate,
                            child: const Padding(
                              padding: EdgeInsets.symmetric(
                                  horizontal: 24.0, vertical: 12.0),
                              child: Text('Calculate'),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
          if (_isLoading)
            Container(
              color: Colors.black54,
              child: const Center(
                child: CircularProgressIndicator(),
              ),
            ),
        ],
      ),
    );
  }
}

class GridPainter extends CustomPainter {
  final Grid2D grid;
  final double minX;
  final double maxX;
  final double minY;
  final double maxY;

  GridPainter({
    required this.grid,
    required this.minX,
    required this.maxX,
    required this.minY,
    required this.maxY,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (size.width <= 0 || size.height <= 0) return;

    const int segments = 100;
    List<List<Point2D>> mappedPoints = List.generate(
      segments + 1,
          (_) => List.generate(segments + 1, (_) => Point2D(0.0, 0.0)),
    );

    double boundMinTx = double.infinity;
    double boundMaxTx = double.negativeInfinity;
    double boundMinTy = double.infinity;
    double boundMaxTy = double.negativeInfinity;

    for (int i = 0; i <= segments; i++) {
      for (int j = 0; j <= segments; j++) {
        double u = i / segments;
        double v = j / segments;

        double origX = minX + u * (maxX - minX);
        double origY = minY + v * (maxY - minY);

        Point2D pIn = Point2D(origX, origY);
        Point2D pOut = grid.evaluate(pIn);

        mappedPoints[i][j] = pOut;

        if (pOut.x < boundMinTx) boundMinTx = pOut.x;
        if (pOut.x > boundMaxTx) boundMaxTx = pOut.x;
        if (pOut.y < boundMinTy) boundMinTy = pOut.y;
        if (pOut.y > boundMaxTy) boundMaxTy = pOut.y;
      }
    }

    double spanX = boundMaxTx - boundMinTx;
    double spanY = boundMaxTy - boundMinTy;
    if (spanX <= 0) spanX = 1e-5;
    if (spanY <= 0) spanY = 1e-5;

    const double padding = 20.0;
    double usableW = size.width - padding * 2;
    double usableH = size.height - padding * 2;

    if (usableW <= 0 || usableH <= 0) return;

    double scale = math.min(usableW / spanX, usableH / spanY);

    double drawnW = spanX * scale;
    double drawnH = spanY * scale;

    double offsetX = padding + (usableW - drawnW) / 2;
    double offsetY = padding + (usableH - drawnH) / 2;

    Offset toScreen(Point2D p) {
      double sx = (p.x - boundMinTx) * scale + offsetX;
      double sy = size.height - ((p.y - boundMinTy) * scale + offsetY);
      return Offset(sx, sy);
    }

    final gridPaint = Paint()
      ..color = Colors.grey.withOpacity(0.3)
      ..strokeWidth = 1
      ..style = PaintingStyle.stroke;

    for (int i = 0; i <= 10; i++) {
      double x = padding + (usableW / 10) * i;
      double y = padding + (usableH / 10) * i;
      canvas.drawLine(Offset(x, padding), Offset(x, size.height - padding), gridPaint);
      canvas.drawLine(Offset(padding, y), Offset(size.width - padding, y), gridPaint);
    }

    for (int i = 0; i < segments; i++) {
      for (int j = 0; j < segments; j++) {
        Offset p1 = toScreen(mappedPoints[i][j]);
        Offset p2 = toScreen(mappedPoints[i + 1][j]);
        Offset p3 = toScreen(mappedPoints[i + 1][j + 1]);
        Offset p4 = toScreen(mappedPoints[i][j + 1]);

        double u = i / segments;
        double v = j / segments;

        Color bottomColor = Color.lerp(Colors.cyan, Colors.green, u)!;
        Color topColor = Color.lerp(Colors.blue, Colors.red, u)!;
        Color cellColor = Color.lerp(bottomColor, topColor, v)!;

        Path path = Path()
          ..moveTo(p1.dx, p1.dy)
          ..lineTo(p2.dx, p2.dy)
          ..lineTo(p3.dx, p3.dy)
          ..lineTo(p4.dx, p4.dy)
          ..close();

        canvas.drawPath(path, Paint()..color = cellColor..style = PaintingStyle.fill);
      }
    }
  }

  @override
  bool shouldRepaint(covariant GridPainter oldDelegate) {
    return oldDelegate.grid != grid ||
        oldDelegate.minX != minX ||
        oldDelegate.maxX != maxX ||
        oldDelegate.minY != minY ||
        oldDelegate.maxY != maxY;
  }
}

