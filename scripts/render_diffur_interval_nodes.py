#!/usr/bin/env python3
import argparse
import json
import math
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


BACKGROUND_COLOR = (255, 255, 255)
FACE_COLOR = (232, 232, 232, 110)
EDGE_COLOR = (80, 80, 80, 255)
POINT_COLOR = (0, 0, 0, 255)
ARROW_COLOR = (40, 40, 40, 255)
LABEL_COLOR = (20, 20, 20, 255)

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
    parser = argparse.ArgumentParser(
        description="Render interval differential equation grid as a 3D time box"
    )
    parser.add_argument(
        "--input",
        default="scripts/grid_data/diffur_interval_tmax_20_eps_0.001.json",
        help="Input JSON file",
    )
    parser.add_argument(
        "--output",
        default="Text/img/diffur_interval_nodes.png",
        help="Output image path",
    )
    parser.add_argument("--width", type=int, default=1600, help="Image width")
    parser.add_argument("--height", type=int, default=800, help="Image height")
    parser.add_argument("--padding", type=int, default=30, help="Canvas padding")
    parser.add_argument("--azimuth", type=float, default=-18.0, help="Rotation of box faces toward the viewer")
    parser.add_argument("--elevation", type=float, default=-8.0, help="Tilt angle")
    parser.add_argument("--point-radius", type=int, default=2, help="Node radius")
    parser.add_argument("--edge-width", type=int, default=3, help="Edge width")
    parser.add_argument("--font-size", type=int, default=28, help="Label font size")
    parser.add_argument("--time-font-size", type=int, default=24, help="Time tick font size")
    parser.add_argument("--time-depth-scale", type=float, default=3.04, help="Relative visible depth of time axis")
    return parser.parse_args()


def load_json(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def load_font(size: int):
    candidates = [
        "/System/Library/Fonts/Supplemental/Consolas.ttf",
        "/Library/Fonts/Consolas.ttf",
        "/System/Library/Fonts/Menlo.ttc",
    ]
    for candidate in candidates:
        font_path = Path(candidate)
        if font_path.exists():
            return ImageFont.truetype(str(font_path), size=size)
    return ImageFont.load_default()


def rotate_point(point, azimuth, elevation):
    x, y, z = point
    ca = math.cos(azimuth)
    sa = math.sin(azimuth)
    x, z = x * ca - z * sa, x * sa + z * ca

    ce = math.cos(elevation)
    se = math.sin(elevation)
    y, z = y * ce - z * se, y * se + z * ce
    return x, y, z


def project(point, scale, center):
    x, y, _ = point
    return center[0] + x * scale, center[1] - y * scale


def fit_projection(rotated_points, width, height, padding, bottom_reserved):
    xs = [point[0] for point in rotated_points]
    ys = [point[1] for point in rotated_points]
    min_x, max_x = min(xs), max(xs)
    min_y, max_y = min(ys), max(ys)
    usable_height = height - 2 * padding - bottom_reserved
    span_x = max(max_x - min_x, 1e-9)
    span_y = max(max_y - min_y, 1e-9)
    scale_x = (width - 2 * padding) / span_x
    scale_y = usable_height / span_y
    scale = min(scale_x, scale_y)
    center = (
        padding - min_x * scale,
        padding + max_y * scale,
    )
    return scale, center


def normalize_point(point, mins, maxs, time_depth_scale):
    time_fraction = (point[2] - mins[2]) / (maxs[2] - mins[2])
    x = (time_fraction - 0.5) * time_depth_scale
    y = (point[1] - mins[1]) / (maxs[1] - mins[1]) - 0.5
    z = (point[0] - mins[0]) / (maxs[0] - mins[0]) - 0.5
    return x, y, z


def build_box_corners(time_depth_scale):
    half_time = time_depth_scale / 2.0
    front_offset = 0.0
    back_offset = 0.0
    return [
        (-half_time, -0.5, -0.5 + front_offset),
        (half_time, -0.5, -0.5 + back_offset),
        (half_time, 0.5, -0.5 + back_offset),
        (-half_time, 0.5, -0.5 + front_offset),
        (-half_time, -0.5, 0.5 + front_offset),
        (half_time, -0.5, 0.5 + back_offset),
        (half_time, 0.5, 0.5 + back_offset),
        (-half_time, 0.5, 0.5 + front_offset),
    ]


def draw_arrow(draw, start, end, color, width):
    draw.line([start, end], fill=color, width=width)
    dx = end[0] - start[0]
    dy = end[1] - start[1]
    angle = math.atan2(dy, dx)
    head_len = 18
    left = (
        end[0] - head_len * math.cos(angle - math.pi / 7),
        end[1] - head_len * math.sin(angle - math.pi / 7),
    )
    right = (
        end[0] - head_len * math.cos(angle + math.pi / 7),
        end[1] - head_len * math.sin(angle + math.pi / 7),
    )
    draw.polygon([end, left, right], fill=color)


def main():
    args = parse_args()
    data = load_json(Path(args.input))
    mins = data["min"]
    maxs = data["max"]
    nodes = data["nodes"]

    if len(mins) != 3 or len(maxs) != 3:
        raise ValueError("Interval visualization expects a 3D grid JSON")

    corners = build_box_corners(args.time_depth_scale)
    normalized_nodes = [
        normalize_point(node["CenterUnit"], [0.0, 0.0, 0.0], [1.0, 1.0, 1.0], args.time_depth_scale)
        for node in nodes
    ]

    azimuth = math.radians(args.azimuth)
    elevation = math.radians(args.elevation)
    rotated_corners = [rotate_point(point, azimuth, elevation) for point in corners]
    rotated_nodes = [rotate_point(point, azimuth, elevation) for point in normalized_nodes]

    bottom_reserved = 150
    scale, center = fit_projection(rotated_corners + rotated_nodes, args.width, args.height, args.padding, bottom_reserved)
    projected_corners = [project(point, scale, center) for point in rotated_corners]

    image = Image.new("RGBA", (args.width, args.height), BACKGROUND_COLOR + (255,))
    draw = ImageDraw.Draw(image, "RGBA")
    label_font = load_font(args.font_size)
    tick_font = load_font(args.time_font_size)

    faces_by_depth = sorted(
        CUBE_FACES,
        key=lambda face: sum(rotated_corners[index][2] for index in face) / len(face),
    )
    for face in faces_by_depth:
        polygon = [projected_corners[index] for index in face]
        draw.polygon(polygon, fill=FACE_COLOR)

    for start, end in CUBE_EDGES:
        draw.line([projected_corners[start], projected_corners[end]], fill=EDGE_COLOR, width=args.edge_width)

    radius = args.point_radius
    for rotated in sorted(rotated_nodes, key=lambda point: point[2]):
        px, py = project(rotated, scale, center)
        draw.ellipse([px - radius, py - radius, px + radius, py + radius], fill=POINT_COLOR, outline=POINT_COLOR)

    arrow_y_shift = 95
    left_x = min(point[0] for point in projected_corners)
    right_x = max(point[0] for point in projected_corners)
    lowest_y = max(point[1] for point in projected_corners)
    arrow_start = (left_x, lowest_y + arrow_y_shift)
    arrow_end = (right_x, lowest_y + arrow_y_shift)
    draw_arrow(draw, arrow_start, arrow_end, ARROW_COLOR, 4)

    for time_value in range(0, 21):
        fraction = time_value / 20.0
        tick_x = arrow_start[0] + (arrow_end[0] - arrow_start[0]) * fraction
        label = str(time_value)
        bbox = draw.textbbox((0, 0), label, font=tick_font)
        draw.text((tick_x - (bbox[2] - bbox[0]) / 2.0, arrow_start[1] + 10), label, fill=LABEL_COLOR, font=tick_font)

    draw.text((arrow_end[0] + 18, arrow_end[1] - 18), "t", fill=LABEL_COLOR, font=label_font)

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    image.convert("RGB").save(output_path)


if __name__ == "__main__":
    main()
