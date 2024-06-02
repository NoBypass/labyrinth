from datetime import datetime

from PIL import Image, ImageDraw

from python.consts import multiplier, grid_size, offset
from python.shapes import Field, Direction, Grid

from random import choice


size = grid_size * multiplier


def main():
    current_time = datetime.now()
    img = Image.new('RGB', (size + offset * 2, size + offset * 2), color='white')
    d = ImageDraw.Draw(img)

    grid = Grid(grid_size)
    grid[0][0].border[Direction.LEFT].hide()
    grid[grid_size - 1][grid_size - 1].border[Direction.RIGHT].hide()

    been_to: dict[Field, None] = {}
    backtracks: list[Field] = []
    pointer = grid[0][0]

    while len(been_to) < grid_size ** 2:
        available: dict[Direction, Field] = {}

        for direction in Direction:
            potential_pointer = pointer.move(grid, direction)
            if potential_pointer and potential_pointer not in been_to:
                available[direction] = potential_pointer

        been_to[pointer] = None
        if not available.keys():
            pointer = backtracks.pop(0)
            continue
        elif len(available) > 1:
            backtracks.append(pointer)

        direction = choice(list(available.keys()))
        pointer.border[direction].hide()
        pointer = available[direction]

    for l in grid.lines:
        l.draw(d)

    img.save('labyrinth.png')
    print(f"Time taken: {(datetime.now() - current_time).microseconds / 1000}ms")


if __name__ == '__main__':
    main()
