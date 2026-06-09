#!/usr/bin/env python3
import argparse
import json
import math
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


BACKGROUND_COLOR = (255, 255, 255)
FACE_COLOR = (242, 242, 242)
EDGE_COLOR = (90, 90, 90)
POINT_COLOR = (0, 0, 0)
DASH_COLOR = (120, 120, 120)
ARROW_COLOR = (40, 40, 40)
LABEL_COLOR = (20, 20, 20)

KNOWN_TIMES = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]
REALIZED_TIMES = [2, 5, 10, 20]


def parse_args():
    parser = argparse.ArgumentParser(
        description="Render iterative differential equation grids as a stack of 2D slices"
    )
    parser.add_argument(
        "--input-dir",
        default="scripts/grid_data",
        help="Directory with exported JSON files",
    )
    parser.add_argument(
        "--output",
        default="Text/img/diffur_iterative_nodes.png",
        help="Output image path",
    )
    parser.add_argument("--width", type=int, default=1600, help="Image width")
    parser.add_argument("--height", type=int, default=630, help="Image height")
    parser.add_argument("--padding", type=int, default=80, help="Canvas padding")
    parser.add_argument("--azimuth", type=float, default=-20.0, help="Rotation of book faces toward the viewer")
    parser.add_argument("--elevation", type=float, default=-8.0, help="Tilt angle")
    parser.add_argument("--slice-gap", type=float, default=0.16, help="Gap between slices in normalized units")
    parser.add_argument("--box-height", type=float, default=0.5, help="Relative height of the box")
    parser.add_argument("--box-depth", type=float, default=0.5, help="Relative depth of the box")
    parser.add_argument("--point-radius", type=int, default=1, help="Node radius")
    parser.add_argument("--edge-width", type=int, default=3, help="Visible slice edge width")
    parser.add_argument("--dash-width", type=int, default=2, help="Dashed slice edge width")
    parser.add_argument("--font-size", type=int, default=28, help="Label font size")
    parser.add_argument("--time-font-size", type=int, default=24, help="Time tick font size")
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
    left = padding - min_x * scale
    top = padding + max_y * scale
    center = (left, top)
    return scale, center


def centered_xy(unit_xy):
    return unit_xy[0] - 0.5, unit_xy[1] - 0.5


def build_slice_geometry(times, slice_gap, data_by_time, box_height, box_depth):
    slices = []
    max_shift = (max(times) - min(times)) * slice_gap
    half_height = box_height / 2.0
    face_depth = box_depth / 2.0
    for time_value in times:
        x_shift = (time_value - min(times)) * slice_gap - max_shift / 2.0
        corners = [
            (x_shift, -half_height, -face_depth),
            (x_shift, -half_height, face_depth),
            (x_shift, half_height, face_depth),
            (x_shift, half_height, -face_depth),
        ]
        nodes = []
        if time_value in data_by_time:
            for node in data_by_time[time_value]["nodes"]:
                x, y = centered_xy(node["CenterUnit"])
                depth = x * box_depth
                nodes.append((x_shift, y * box_height, depth))
        slices.append({
            "time": time_value,
            "realized": time_value in data_by_time,
            "corners": corners,
            "nodes": nodes,
            "depth": x_shift,
        })
    return slices


def draw_dashed_line(draw, start, end, color, width, dash=12, gap=8):
    x1, y1 = start
    x2, y2 = end
    dx = x2 - x1
    dy = y2 - y1
    length = math.hypot(dx, dy)
    if length == 0:
        return
    step_x = dx / length
    step_y = dy / length
    position = 0.0
    while position < length:
        segment_end = min(position + dash, length)
        p1 = (x1 + step_x * position, y1 + step_y * position)
        p2 = (x1 + step_x * segment_end, y1 + step_y * segment_end)
        draw.line([p1, p2], fill=color, width=width)
        position += dash + gap


def draw_slice(draw, projected_corners, realized, edge_width, dash_width):
    polygon = list(projected_corners)
    if realized:
        draw.polygon(polygon, fill=FACE_COLOR)
        for idx in range(4):
            draw.line([polygon[idx], polygon[(idx + 1) % 4]], fill=EDGE_COLOR, width=edge_width)
    else:
        for idx in range(4):
            draw_dashed_line(
                draw,
                polygon[idx],
                polygon[(idx + 1) % 4],
                color=DASH_COLOR,
                width=dash_width,
            )


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
    input_dir = Path(args.input_dir)
    data_by_time = {}
    for time_value in REALIZED_TIMES:
        path = input_dir / f"diffur_iterative_tmax_{time_value}_eps_0.001.json"
        data_by_time[time_value] = load_json(path)

    slices = build_slice_geometry(KNOWN_TIMES, args.slice_gap, data_by_time, args.box_height, args.box_depth)
    azimuth = math.radians(args.azimuth)
    elevation = math.radians(args.elevation)

    rotated_points = []
    for slice_data in slices:
        rotated_points.extend(rotate_point(point, azimuth, elevation) for point in slice_data["corners"])
        rotated_points.extend(rotate_point(point, azimuth, elevation) for point in slice_data["nodes"])

    bottom_reserved = 150
    scale, center = fit_projection(rotated_points, args.width, args.height, args.padding, bottom_reserved)

    image = Image.new("RGB", (args.width, args.height), BACKGROUND_COLOR)
    draw = ImageDraw.Draw(image)
    label_font = load_font(args.font_size)
    tick_font = load_font(args.time_font_size)

    projected_slices = []
    for slice_data in slices:
        rotated_corners = [rotate_point(point, azimuth, elevation) for point in slice_data["corners"]]
        rotated_nodes = [rotate_point(point, azimuth, elevation) for point in slice_data["nodes"]]
        projected_slices.append(
            {
                "time": slice_data["time"],
                "realized": slice_data["realized"],
                "depth": slice_data["depth"],
                "corners": [project(point, scale, center) for point in rotated_corners],
                "nodes": [project(point, scale, center) for point in rotated_nodes],
            }
        )

    projected_slices.sort(key=lambda item: item["depth"])
    for slice_data in projected_slices:
        draw_slice(draw, slice_data["corners"], slice_data["realized"], args.edge_width, args.dash_width)
        if slice_data["realized"]:
            for px, py in slice_data["nodes"]:
                radius = args.point_radius
                draw.ellipse([px - radius, py - radius, px + radius, py + radius], fill=POINT_COLOR, outline=POINT_COLOR)

    arrow_y_shift = 95
    left_x = min(min(point[0] for point in slice_data["corners"]) for slice_data in projected_slices)
    right_x = max(max(point[0] for point in slice_data["corners"]) for slice_data in projected_slices)
    lowest_y = max(max(point[1] for point in slice_data["corners"]) for slice_data in projected_slices)
    arrow_start = (left_x, lowest_y + arrow_y_shift)
    arrow_end = (right_x, lowest_y + arrow_y_shift)
    draw_arrow(draw, arrow_start, arrow_end, ARROW_COLOR, 4)

    for slice_data in projected_slices:
        center_x = sum(point[0] for point in slice_data["corners"]) / 4.0
        center_y = arrow_start[1] + 10
        label = str(slice_data["time"])
        bbox = draw.textbbox((0, 0), label, font=tick_font)
        draw.text((center_x - (bbox[2] - bbox[0]) / 2.0, center_y), label, fill=LABEL_COLOR, font=tick_font)

    draw.text((arrow_end[0] + 18, arrow_end[1] - 18), "t", fill=LABEL_COLOR, font=label_font)

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    image.save(output_path)


if __name__ == "__main__":
    main()
