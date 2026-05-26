#!/usr/bin/env python3
import argparse
import math
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


BACKGROUND = "white"
AXIS_COLOR = (60, 60, 60)
GRID_COLOR = (230, 230, 230)
TRUE_COLOR = (32, 91, 170)
WITHOUT_ANCHORS_COLOR = (120, 120, 120)
WITH_ANCHORS_COLOR = (218, 84, 45)
INITIAL_POINT_COLOR = (0, 0, 0)
ANCHOR_POINT_COLOR = (180, 0, 0)
TEXT_COLOR = (30, 30, 30)


def parse_args():
    parser = argparse.ArgumentParser(description="Render anchor-points demonstration for sin(x)")
    parser.add_argument("output", help="Path to output image file")
    parser.add_argument("--width", type=int, default=1250, help="Image width in pixels")
    parser.add_argument("--height", type=int, default=1000, help="Image height in pixels")
    parser.add_argument("--padding", type=int, default=70, help="Padding around drawing area")
    parser.add_argument("--font-size", type=int, default=27, help="Label font size")
    parser.add_argument("--samples", type=int, default=600, help="Number of samples for the true sin(x) curve")
    return parser.parse_args()


def load_monospace_cyrillic(size):
    candidates = [
        "/System/Library/Fonts/Supplemental/Courier New.ttf",
        "/Library/Fonts/Courier New.ttf",
        "/System/Library/Fonts/Menlo.ttc",
        "/Library/Fonts/Menlo.ttc",
        "/System/Library/Fonts/SFNSMono.ttf",
        "/System/Library/Fonts/Monaco.ttf",
        "/opt/homebrew/share/fonts/dejavu/DejaVuSansMono.ttf",
        "/usr/local/share/fonts/dejavu/DejaVuSansMono.ttf",
        "/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
        "/Library/Fonts/Arial Unicode.ttf",
    ]
    for candidate in candidates:
        try:
            return ImageFont.truetype(candidate, size)
        except OSError:
            pass
    return ImageFont.load_default()


def to_pixel(x, y, x_min, x_max, y_min, y_max, left, top, right, bottom):
    px = left + (x - x_min) * (right - left) / (x_max - x_min)
    py = bottom - (y - y_min) * (bottom - top) / (y_max - y_min)
    return px, py


def draw_dashed_line(draw, start, end, fill, width=2, dash=15, gap=10):
    x1, y1 = start
    x2, y2 = end
    length = math.hypot(x2 - x1, y2 - y1)
    if length == 0:
        return
    dx = (x2 - x1) / length
    dy = (y2 - y1) / length
    current = 0.0
    while current < length:
        segment_end = min(current + dash, length)
        draw.line(
            [
                (x1 + dx * current, y1 + dy * current),
                (x1 + dx * segment_end, y1 + dy * segment_end),
            ],
            fill=fill,
            width=width,
        )
        current += dash + gap


def draw_polyline(draw, points, fill, width=3):
    if len(points) >= 2:
        draw.line(points, fill=fill, width=width, joint="curve")


def draw_point(draw, point, radius, fill, outline=None):
    x, y = point
    draw.ellipse([x - radius, y - radius, x + radius, y + radius], fill=fill, outline=outline or fill)


def draw_label(draw, xy, text, font, fill=TEXT_COLOR):
    draw.text(xy, text, fill=fill, font=font)


def main():
    args = parse_args()

    x_min = -math.pi
    x_max = math.pi
    y_min = -1.25
    y_max = 1.25

    left = args.padding
    top = 135
    right = args.width - args.padding
    bottom = args.height - args.padding - 65

    image = Image.new("RGB", (args.width, args.height), BACKGROUND)
    draw = ImageDraw.Draw(image)
    font = load_monospace_cyrillic(args.font_size)

    for tick_x in [-math.pi, -math.pi / 2, 0.0, math.pi / 2, math.pi]:
        px, _ = to_pixel(tick_x, 0.0, x_min, x_max, y_min, y_max, left, top, right, bottom)
        draw.line([(px, top), (px, bottom)], fill=GRID_COLOR, width=1)
    for tick_y in [-1.0, -0.5, 0.0, 0.5, 1.0]:
        _, py = to_pixel(0.0, tick_y, x_min, x_max, y_min, y_max, left, top, right, bottom)
        draw.line([(left, py), (right, py)], fill=GRID_COLOR, width=1)

    x_axis_start = to_pixel(x_min, 0.0, x_min, x_max, y_min, y_max, left, top, right, bottom)
    x_axis_end = to_pixel(x_max, 0.0, x_min, x_max, y_min, y_max, left, top, right, bottom)
    y_axis_start = to_pixel(0.0, y_min, x_min, x_max, y_min, y_max, left, top, right, bottom)
    y_axis_end = to_pixel(0.0, y_max, x_min, x_max, y_min, y_max, left, top, right, bottom)
    draw.line([x_axis_start, x_axis_end], fill=AXIS_COLOR, width=2)
    draw.line([y_axis_start, y_axis_end], fill=AXIS_COLOR, width=2)

    true_curve = []
    for i in range(args.samples + 1):
        x = x_min + (x_max - x_min) * i / args.samples
        true_curve.append(to_pixel(x, math.sin(x), x_min, x_max, y_min, y_max, left, top, right, bottom))
    draw_polyline(draw, true_curve, TRUE_COLOR, width=4)

    initial_points = [(-math.pi, 0.0), (0.0, 0.0), (math.pi, 0.0)]
    initial_pixels = [to_pixel(x, y, x_min, x_max, y_min, y_max, left, top, right, bottom) for x, y in initial_points]
    draw_dashed_line(draw, initial_pixels[0], initial_pixels[-1], WITHOUT_ANCHORS_COLOR, width=6)

    anchor_points = [(-math.pi / 2, -1.0), (math.pi / 2, 1.0)]
    with_anchor_points = [initial_points[0], anchor_points[0], initial_points[1], anchor_points[1], initial_points[2]]
    with_anchor_pixels = [to_pixel(x, y, x_min, x_max, y_min, y_max, left, top, right, bottom) for x, y in with_anchor_points]
    draw_polyline(draw, with_anchor_pixels, WITH_ANCHORS_COLOR, width=3)

    for pixel in initial_pixels:
        draw_point(draw, pixel, 6, INITIAL_POINT_COLOR)
    for x, y in anchor_points:
        draw_point(draw, to_pixel(x, y, x_min, x_max, y_min, y_max, left, top, right, bottom), 7, ANCHOR_POINT_COLOR)

    tick_labels = [
        (-math.pi, "-pi"),
        (-math.pi / 2, "-pi/2"),
        (0.0, "0"),
        (math.pi / 2, "pi/2"),
        (math.pi, "pi"),
    ]
    for x, label in tick_labels:
        px, py = to_pixel(x, 0.0, x_min, x_max, y_min, y_max, left, top, right, bottom)
        draw.line([(px, py - 5), (px, py + 5)], fill=AXIS_COLOR, width=2)
        draw_label(draw, (px - 14, py + 12), label, font)
    for y, label in [(-1.0, "-1"), (1.0, "1")]:
        px, py = to_pixel(0.0, y, x_min, x_max, y_min, y_max, left, top, right, bottom)
        draw.line([(px - 5, py), (px + 5, py)], fill=AXIS_COLOR, width=2)
        draw_label(draw, (px + 10, py - 7), label, font)

    draw_label(draw, (left, 18), "f(x) = sin(x), x: [-pi, pi]", font)
    draw_label(draw, (left, 55), "Черные точки: узлы -pi, 0, pi дают ноль -> сначала получается прямая", font)
    draw_label(draw, (left, 92), "Красные точки: якоря запускают уточнение и восстанавливают синус", font)

    legend_x = left
    legend_y = bottom + 45
    draw.line([(legend_x, legend_y), (legend_x + 48, legend_y)], fill=TRUE_COLOR, width=4)
    draw_label(draw, (legend_x + 60, legend_y - 14), "изначальная функция sin(x)", font)
    draw_dashed_line(draw, (legend_x + 510, legend_y), (legend_x + 558, legend_y), WITHOUT_ANCHORS_COLOR, width=4)
    draw_label(draw, (legend_x + 570, legend_y - 14), "без якорей", font)
    draw.line([(legend_x + 760, legend_y), (legend_x + 808, legend_y)], fill=WITH_ANCHORS_COLOR, width=3)
    draw_label(draw, (legend_x + 820, legend_y - 14), "с якорями", font)

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    image.save(output_path)


if __name__ == "__main__":
    main()
