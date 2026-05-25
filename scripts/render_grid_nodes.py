#!/usr/bin/env python3
import argparse
import json
from pathlib import Path

from PIL import Image, ImageDraw


def parse_args():
    parser = argparse.ArgumentParser(description="Render grid nodes from exported JSON")
    parser.add_argument("input", help="Path to input JSON file")
    parser.add_argument("output", help="Path to output image file")
    parser.add_argument("--width", type=int, default=1200, help="Image width in pixels")
    parser.add_argument("--height", type=int, default=800, help="Image height in pixels")
    parser.add_argument("--padding", type=int, default=40, help="Padding around drawing area")
    parser.add_argument("--border-width", type=int, default=3, help="Rectangle border width")
    parser.add_argument("--point-radius", type=int, default=4, help="Node point radius")
    return parser.parse_args()


def load_grid(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def to_pixel(value, left, right, start, end):
    if right - left == 0:
        return (start + end) / 2
    return start + (value - left) * (end - start) / (right - left)


def main():
    args = parse_args()
    data = load_grid(Path(args.input))

    min_x, min_y = data["min"]
    max_x, max_y = data["max"]
    nodes = data["nodes"]

    image = Image.new("RGB", (args.width, args.height), "white")
    draw = ImageDraw.Draw(image)

    left = args.padding
    top = args.padding
    right = args.width - args.padding
    bottom = args.height - args.padding

    draw.rectangle([left, top, right, bottom], outline="grey", width=args.border_width)

    radius = args.point_radius
    for node in nodes:
        x_unit, y_unit = node["CenterUnit"]
        x_real = min_x + x_unit * (max_x - min_x)
        y_real = min_y + y_unit * (max_y - min_y)

        px = to_pixel(x_real, min_x, max_x, left, right)
        py = to_pixel(y_real, min_y, max_y, bottom, top)

        draw.ellipse(
            [px - radius, py - radius, px + radius, py + radius],
            fill="black",
            outline="black",
        )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    image.save(output_path)


if __name__ == "__main__":
    main()
