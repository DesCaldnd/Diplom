#!/usr/bin/env python3
import argparse
import json
import math
from pathlib import Path

from PIL import Image, ImageDraw


FACE_COLOR = (220, 220, 220, 90)
EDGE_COLOR = (70, 70, 70, 255)
POINT_COLOR = (0, 0, 0, 255)


CUBE_CORNERS = [
    (0, 0, 0),
    (1, 0, 0),
    (1, 1, 0),
    (0, 1, 0),
    (0, 0, 1),
    (1, 0, 1),
    (1, 1, 1),
    (0, 1, 1),
]

CUBE_EDGES = [
    (0, 1),
    (1, 2),
    (2, 3),
    (3, 0),
    (4, 5),
    (5, 6),
    (6, 7),
    (7, 4),
    (0, 4),
    (1, 5),
    (2, 6),
    (3, 7),
]

CUBE_FACES = [
    (0, 1, 2, 3),
    (4, 5, 6, 7),
    (0, 1, 5, 4),
    (1, 2, 6, 5),
    (2, 3, 7, 6),
    (3, 0, 4, 7),
]


def parse_args():
    parser = argparse.ArgumentParser(description="Render 3D grid nodes from exported JSON")
    parser.add_argument("input", help="Path to input JSON file")
    parser.add_argument("output", help="Path to output image file")
    parser.add_argument("--width", type=int, default=1200, help="Image width in pixels")
    parser.add_argument("--height", type=int, default=900, help="Image height in pixels")
    parser.add_argument("--padding", type=int, default=80, help="Padding around drawing area")
    parser.add_argument("--edge-width", type=int, default=3, help="Cube edge width")
    parser.add_argument("--point-radius", type=int, default=4, help="Node point radius")
    parser.add_argument("--azimuth", type=float, default=-38.0, help="Camera azimuth angle in degrees")
    parser.add_argument("--elevation", type=float, default=26.0, help="Camera elevation angle in degrees")
    parser.add_argument("--roll", type=float, default=0.0, help="Camera roll angle in degrees")
    return parser.parse_args()


def load_grid(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def validate_grid(data):
    mins = data["min"]
    maxs = data["max"]
    if len(mins) != 3 or len(maxs) != 3:
        raise ValueError("input grid must be three-dimensional: len(min) == len(max) == 3")
    return mins, maxs, data["nodes"]


def unit_to_real(unit, mins, maxs):
    return tuple(mins[i] + unit[i] * (maxs[i] - mins[i]) for i in range(3))


def real_to_centered(point, mins, maxs):
    spans = [maxs[i] - mins[i] for i in range(3)]
    max_span = max(spans) if max(spans) != 0 else 1.0
    return tuple((point[i] - (mins[i] + maxs[i]) / 2.0) / max_span for i in range(3))


def rotate_point(point, azimuth, elevation, roll):
    x, y, z = point

    ca = math.cos(azimuth)
    sa = math.sin(azimuth)
    x, y = x * ca - y * sa, x * sa + y * ca

    ce = math.cos(elevation)
    se = math.sin(elevation)
    y, z = y * ce - z * se, y * se + z * ce

    cr = math.cos(roll)
    sr = math.sin(roll)
    x, y = x * cr - y * sr, x * sr + y * cr

    return x, y, z


def project(point, scale, center):
    x, y, _ = point
    return center[0] + x * scale, center[1] - y * scale


def fit_projection(rotated_points, width, height, padding):
    xs = [p[0] for p in rotated_points]
    ys = [p[1] for p in rotated_points]
    span_x = max(xs) - min(xs)
    span_y = max(ys) - min(ys)
    scale_x = (width - 2 * padding) / span_x if span_x else 1.0
    scale_y = (height - 2 * padding) / span_y if span_y else 1.0
    scale = min(scale_x, scale_y)
    center = (
        width / 2.0 - (min(xs) + max(xs)) * scale / 2.0,
        height / 2.0 + (min(ys) + max(ys)) * scale / 2.0,
    )
    return scale, center


def main():
    args = parse_args()
    data = load_grid(Path(args.input))
    mins, maxs, nodes = validate_grid(data)

    azimuth = math.radians(args.azimuth)
    elevation = math.radians(args.elevation)
    roll = math.radians(args.roll)

    cube_real = [unit_to_real(corner, mins, maxs) for corner in CUBE_CORNERS]
    cube_centered = [real_to_centered(point, mins, maxs) for point in cube_real]
    cube_rotated = [rotate_point(point, azimuth, elevation, roll) for point in cube_centered]

    node_points = []
    for node in nodes:
        unit = node["CenterUnit"]
        if len(unit) != 3:
            raise ValueError("every node CenterUnit must be three-dimensional")
        real = unit_to_real(unit, mins, maxs)
        centered = real_to_centered(real, mins, maxs)
        node_points.append(rotate_point(centered, azimuth, elevation, roll))

    scale, center = fit_projection(cube_rotated + node_points, args.width, args.height, args.padding)
    cube_projected = [project(point, scale, center) for point in cube_rotated]

    image = Image.new("RGBA", (args.width, args.height), "white")
    draw = ImageDraw.Draw(image, "RGBA")

    faces_by_depth = sorted(
        CUBE_FACES,
        key=lambda face: sum(cube_rotated[index][2] for index in face) / len(face),
    )
    for face in faces_by_depth:
        polygon = [cube_projected[index] for index in face]
        draw.polygon(polygon, fill=FACE_COLOR)

    for start, end in CUBE_EDGES:
        draw.line([cube_projected[start], cube_projected[end]], fill=EDGE_COLOR, width=args.edge_width)

    radius = args.point_radius
    for rotated in sorted(node_points, key=lambda point: point[2]):
        px, py = project(rotated, scale, center)
        draw.ellipse(
            [px - radius, py - radius, px + radius, py + radius],
            fill=POINT_COLOR,
            outline=POINT_COLOR,
        )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    image.convert("RGB").save(output_path)


if __name__ == "__main__":
    main()
